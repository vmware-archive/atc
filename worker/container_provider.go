package worker

import (
	"context"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/metric"
	"github.com/concourse/baggageclaim"
)

const creatingContainerRetryDelay = 1 * time.Second

func NewContainerProvider(
	gardenClient garden.Client,
	baggageclaimClient baggageclaim.Client,
	volumeClient VolumeClient,
	dbWorker db.Worker,
	clock clock.Clock,
	//TODO: less of all this junk..
	imageFactory ImageFactory,
	dbVolumeFactory db.VolumeFactory,
	dbTeamFactory db.TeamFactory,
	lockFactory lock.LockFactory,
) ContainerProvider {

	return &DbContainerProvider{
		DbVolumeFactory: dbVolumeFactory,
		DbTeamFactory:   dbTeamFactory,
		LockFactory:     lockFactory,
		Worker:          dbWorker,
		Clock:           clock,
		client: &containerProvider{
			gardenClient:       gardenClient,
			baggageclaimClient: baggageclaimClient,
			volumeClient:       volumeClient,
			imageFactory:       imageFactory,
			dbVolumeFactory:    dbVolumeFactory,
			dbTeamFactory:      dbTeamFactory,
			lockFactory:        lockFactory,
			httpProxyURL:       dbWorker.HTTPProxyURL(),
			httpsProxyURL:      dbWorker.HTTPSProxyURL(),
			noProxy:            dbWorker.NoProxy(),
			clock:              clock,
			worker:             dbWorker,
		},
	}
}

//go:generate counterfeiter . ContainerProvider

type ContainerProvider interface {
	FindCreatedContainerByHandle(
		logger lager.Logger,
		handle string,
		teamID int,
	) (Container, bool, error)

	FindOrCreateContainer(
		ctx context.Context,
		logger lager.Logger,
		owner db.ContainerOwner,
		delegate ImageFetchingDelegate,
		metadata db.ContainerMetadata,
		spec ContainerSpec,
		resourceTypes creds.VersionedResourceTypes,
	) (Container, error)
}

type containerProvider struct {
	gardenClient       garden.Client
	baggageclaimClient baggageclaim.Client
	volumeClient       VolumeClient
	imageFactory       ImageFactory

	dbVolumeFactory db.VolumeFactory
	dbTeamFactory   db.TeamFactory

	lockFactory lock.LockFactory

	worker        db.Worker
	httpProxyURL  string
	httpsProxyURL string
	noProxy       string

	clock clock.Clock
}

type containerClient interface {
	found(lager.Logger, db.CreatedContainer) (Container, error)
	finalize(lager.Logger, db.CreatingContainer) (Container, bool, error)
	fetchImage(
		logger lager.Logger,
		provider ContainerProvider,
		spec ContainerSpec,
		delegate ImageFetchingDelegate,
		resourceTypes creds.VersionedResourceTypes,
	) (Image, error)
	create(
		ctx context.Context,
		logger lager.Logger,
		provider ContainerProvider,
		creatingContainer db.CreatingContainer,
		spec ContainerSpec,
		image Image,
	) (Container, error)
	find(lager.Logger, db.CreatedContainer, []db.CreatedVolume) (Container, bool, error)
}

type DbContainerProvider struct {
	DbVolumeFactory db.VolumeFactory
	DbTeamFactory   db.TeamFactory
	LockFactory     lock.LockFactory
	Worker          db.Worker
	Clock           clock.Clock

	client containerClient
}

func (p *DbContainerProvider) FindOrCreateContainer(
	ctx context.Context,
	logger lager.Logger,
	owner db.ContainerOwner,
	delegate ImageFetchingDelegate,
	metadata db.ContainerMetadata,
	spec ContainerSpec,
	resourceTypes creds.VersionedResourceTypes,
) (Container, error) {
	for {
		creatingContainer, createdContainer, err := p.DbTeamFactory.GetByID(spec.TeamID).FindContainerOnWorker(
			p.Worker.Name(),
			owner,
		)
		if err != nil {
			logger.Error("failed-to-find-container-in-db", err)
			return nil, err
		}

		if createdContainer != nil {
			logger = logger.WithData(lager.Data{"container": createdContainer.Handle()})
			logger.Debug("found-created-container-in-db")
			return p.client.found(logger, createdContainer)
		}

		var image Image

		if creatingContainer != nil {
			container, done, err := p.client.finalize(logger, creatingContainer)
			if err != nil {
				return nil, err
			}
			if done {
				return container, nil
			}
		} else {
			image, err = p.client.fetchImage(logger, p, spec, delegate, resourceTypes)
			if err != nil {
				logger.Error("failed-to-get-image-for-container", err)
				return nil, err
			}

			creatingContainer, err = p.DbTeamFactory.GetByID(spec.TeamID).CreateContainer(
				p.Worker.Name(),
				owner,
				metadata,
			)
			if err != nil {
				logger.Error("failed-to-create-container-in-db", err)
				return nil, err
			}
		}

		if image == nil {
			image, err = p.client.fetchImage(logger, p, spec, delegate, resourceTypes)
			if err != nil {
				logger.Error("failed-to-get-image-for-container", err)
				return nil, err
			}
		}

		logger = logger.WithData(lager.Data{"container": creatingContainer.Handle()})

		logger.Debug("created-creating-container-in-db")

		lock, acquired, err := p.LockFactory.Acquire(logger, lock.NewContainerCreatingLockID(creatingContainer.ID()))
		if err != nil {
			logger.Error("failed-to-acquire-container-creating-lock", err)
			return nil, err
		}

		if !acquired {
			p.Clock.Sleep(creatingContainerRetryDelay)
			continue
		}

		defer lock.Release()

		return p.client.create(ctx, logger, p, creatingContainer, spec, image)
	}
}

func (p *DbContainerProvider) FindCreatedContainerByHandle(
	logger lager.Logger,
	handle string,
	teamID int,
) (Container, bool, error) {
	createdContainer, found, err := p.DbTeamFactory.GetByID(teamID).FindCreatedContainerByHandle(handle)
	if err != nil {
		logger.Error("failed-to-lookup-in-db", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	createdVolumes, err := p.DbVolumeFactory.FindVolumesForContainer(createdContainer)
	if err != nil {
		return nil, false, err
	}

	return p.client.find(logger, createdContainer, createdVolumes)
}

func (p *containerProvider) found(logger lager.Logger, createdContainer db.CreatedContainer) (Container, error) {
	gardenContainer, err := p.gardenClient.Lookup(createdContainer.Handle())
	if err != nil {
		logger.Error("failed-to-lookup-created-container-in-garden", err)
		return nil, err
	}

	return p.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
}

func (p *containerProvider) finalize(logger lager.Logger, creatingContainer db.CreatingContainer) (Container, bool, error) {
	logger = logger.WithData(lager.Data{"container": creatingContainer.Handle()})
	logger.Debug("found-creating-container-in-db")

	gardenContainer, err := p.gardenClient.Lookup(creatingContainer.Handle())
	if err != nil {
		if _, ok := err.(garden.ContainerNotFoundError); !ok {
			logger.Error("failed-to-lookup-creating-container-in-garden", err)
			return nil, false, err
		} else {
			return nil, false, nil
		}
	}

	createdContainer, err := creatingContainer.Created()
	if err != nil {
		logger.Error("failed-to-mark-container-as-created", err)

		_ = p.gardenClient.Destroy(creatingContainer.Handle())

		return nil, false, err
	}

	logger.Debug("created-container-in-db")

	container, err := p.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
	if err != nil {
		return nil, false, err
	}

	return container, true, nil
}

func (p *containerProvider) fetchImage(
	logger lager.Logger,
	provider ContainerProvider,
	spec ContainerSpec,
	delegate ImageFetchingDelegate,
	resourceTypes creds.VersionedResourceTypes,
) (Image, error) {
	worker := NewGardenWorker(
		p.gardenClient,
		p.baggageclaimClient,
		provider,
		p.volumeClient,
		p.worker,
		p.clock,
	)

	return p.imageFactory.GetImage(
		logger,
		worker,
		p.volumeClient,
		spec.ImageSpec,
		spec.TeamID,
		delegate,
		resourceTypes,
	)
}

func (p *containerProvider) create(
	ctx context.Context,
	logger lager.Logger,
	provider ContainerProvider,
	creatingContainer db.CreatingContainer,
	spec ContainerSpec,
	image Image,
) (Container, error) {

	fetchedImage, err := image.FetchForContainer(ctx, logger, creatingContainer)
	if err != nil {
		creatingContainer.Failed()
		logger.Error("failed-to-fetch-image-for-container", err)
		return nil, err
	}

	logger.Debug("creating-container-in-garden")

	gardenContainer, err := p.createGardenContainer(
		logger,
		creatingContainer,
		spec,
		fetchedImage,
		provider,
	)
	if err != nil {
		_, failedErr := creatingContainer.Failed()
		if failedErr != nil {
			logger.Error("failed-to-mark-container-as-failed", err)
		}
		metric.FailedContainers.Inc()

		logger.Error("failed-to-create-container-in-garden", err)
		return nil, err
	}

	metric.ContainersCreated.Inc()

	logger.Debug("created-container-in-garden")

	createdContainer, err := creatingContainer.Created()
	if err != nil {
		logger.Error("failed-to-mark-container-as-created", err)

		_ = p.gardenClient.Destroy(creatingContainer.Handle())

		return nil, err
	}

	logger.Debug("created-container-in-db")

	return p.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
}

func (p *containerProvider) find(
	logger lager.Logger,
	createdContainer db.CreatedContainer,
	createdVolumes []db.CreatedVolume,
) (Container, bool, error) {

	gardenContainer, err := p.gardenClient.Lookup(createdContainer.Handle())
	if err != nil {
		if _, ok := err.(garden.ContainerNotFoundError); ok {
			logger.Info("container-not-found")
			return nil, false, nil
		}

		logger.Error("failed-to-lookup-on-garden", err)
		return nil, false, err
	}

	container, err := newGardenWorkerContainer(
		logger,
		gardenContainer,
		createdContainer,
		createdVolumes,
		p.gardenClient,
		p.volumeClient,
		p.worker.Name(),
	)

	if err != nil {
		logger.Error("failed-to-construct-container", err)
		return nil, false, err
	}

	return container, true, nil
}

func (p *containerProvider) constructGardenWorkerContainer(
	logger lager.Logger,
	createdContainer db.CreatedContainer,
	gardenContainer garden.Container,
) (Container, error) {
	createdVolumes, err := p.dbVolumeFactory.FindVolumesForContainer(createdContainer)
	if err != nil {
		logger.Error("failed-to-find-container-volumes", err)
		return nil, err
	}

	return newGardenWorkerContainer(
		logger,
		gardenContainer,
		createdContainer,
		createdVolumes,
		p.gardenClient,
		p.volumeClient,
		p.worker.Name(),
	)
}

func (p *containerProvider) createGardenContainer(
	logger lager.Logger,
	creatingContainer db.CreatingContainer,
	spec ContainerSpec,
	fetchedImage FetchedImage,
	provider ContainerProvider,
) (garden.Container, error) {
	volumeMounts := []VolumeMount{}

	scratchVolume, err := p.volumeClient.FindOrCreateVolumeForContainer(
		logger,
		VolumeSpec{
			Strategy:   baggageclaim.EmptyStrategy{},
			Privileged: fetchedImage.Privileged,
		},
		creatingContainer,
		spec.TeamID,
		"/scratch",
	)
	if err != nil {
		return nil, err
	}

	volumeMounts = append(volumeMounts, VolumeMount{
		Volume:    scratchVolume,
		MountPath: "/scratch",
	})

	if spec.Dir != "" && !p.anyMountTo(spec.Dir, spec.Inputs) {
		workdirVolume, volumeErr := p.volumeClient.FindOrCreateVolumeForContainer(
			logger,
			VolumeSpec{
				Strategy:   baggageclaim.EmptyStrategy{},
				Privileged: fetchedImage.Privileged,
			},
			creatingContainer,
			spec.TeamID,
			spec.Dir,
		)
		if volumeErr != nil {
			return nil, volumeErr
		}

		volumeMounts = append(volumeMounts, VolumeMount{
			Volume:    workdirVolume,
			MountPath: spec.Dir,
		})
	}

	worker := NewGardenWorker(
		p.gardenClient,
		p.baggageclaimClient,
		provider,
		p.volumeClient,
		p.worker,
		p.clock,
	)

	for _, inputSource := range spec.Inputs {
		var inputVolume Volume

		localVolume, found, err := inputSource.Source().VolumeOn(worker)
		if err != nil {
			return nil, err
		}

		if found {
			inputVolume, err = p.volumeClient.FindOrCreateCOWVolumeForContainer(
				logger,
				VolumeSpec{
					Strategy:   localVolume.COWStrategy(),
					Privileged: fetchedImage.Privileged,
				},
				creatingContainer,
				localVolume,
				spec.TeamID,
				inputSource.DestinationPath(),
			)
			if err != nil {
				return nil, err
			}
		} else {
			inputVolume, err = p.volumeClient.FindOrCreateVolumeForContainer(
				logger,
				VolumeSpec{
					Strategy:   baggageclaim.EmptyStrategy{},
					Privileged: fetchedImage.Privileged,
				},
				creatingContainer,
				spec.TeamID,
				inputSource.DestinationPath(),
			)
			if err != nil {
				return nil, err
			}

			err = inputSource.Source().StreamTo(inputVolume)
			if err != nil {
				return nil, err
			}
		}

		volumeMounts = append(volumeMounts, VolumeMount{
			Volume:    inputVolume,
			MountPath: inputSource.DestinationPath(),
		})
	}

	for _, outputPath := range spec.Outputs {
		outVolume, volumeErr := p.volumeClient.FindOrCreateVolumeForContainer(
			logger,
			VolumeSpec{
				Strategy:   baggageclaim.EmptyStrategy{},
				Privileged: fetchedImage.Privileged,
			},
			creatingContainer,
			spec.TeamID,
			outputPath,
		)
		if volumeErr != nil {
			return nil, volumeErr
		}

		volumeMounts = append(volumeMounts, VolumeMount{
			Volume:    outVolume,
			MountPath: outputPath,
		})
	}

	bindMounts := []garden.BindMount{}

	for _, mount := range spec.BindMounts {
		bindMount, found, mountErr := mount.VolumeOn(worker)
		if mountErr != nil {
			return nil, mountErr
		}
		if found {
			bindMounts = append(bindMounts, bindMount)
		}
	}

	volumeHandleMounts := map[string]string{}
	for _, mount := range volumeMounts {
		bindMounts = append(bindMounts, garden.BindMount{
			SrcPath: mount.Volume.Path(),
			DstPath: mount.MountPath,
			Mode:    garden.BindMountModeRW,
		})
		volumeHandleMounts[mount.Volume.Handle()] = mount.MountPath
	}

	gardenProperties := garden.Properties{}

	if spec.User != "" {
		gardenProperties[userPropertyName] = spec.User
	} else {
		gardenProperties[userPropertyName] = fetchedImage.Metadata.User
	}

	env := append(fetchedImage.Metadata.Env, spec.Env...)

	if p.httpProxyURL != "" {
		env = append(env, fmt.Sprintf("http_proxy=%s", p.httpProxyURL))
	}

	if p.httpsProxyURL != "" {
		env = append(env, fmt.Sprintf("https_proxy=%s", p.httpsProxyURL))
	}

	if p.noProxy != "" {
		env = append(env, fmt.Sprintf("no_proxy=%s", p.noProxy))
	}

	return p.gardenClient.Create(garden.ContainerSpec{
		Handle:     creatingContainer.Handle(),
		RootFSPath: fetchedImage.URL,
		Privileged: fetchedImage.Privileged,
		BindMounts: bindMounts,
		Env:        env,
		Properties: gardenProperties,
	})
}

func (p *containerProvider) anyMountTo(path string, inputs []InputSource) bool {
	for _, input := range inputs {
		if input.DestinationPath() == path {
			return true
		}
	}

	return false
}
