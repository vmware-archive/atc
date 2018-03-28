package worker

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/metric"
	"github.com/concourse/baggageclaim"
)

type gardenContainerClient struct {
	gardenClient       garden.Client
	baggageclaimClient baggageclaim.Client
	volumeClient       VolumeClient
	imageFactory       ImageFactory

	dbVolumeFactory db.VolumeFactory
	worker          db.Worker

	clock clock.Clock
}

func (c *gardenContainerClient) Found(logger lager.Logger, createdContainer db.CreatedContainer) (Container, error) {
	gardenContainer, err := c.gardenClient.Lookup(createdContainer.Handle())
	if err != nil {
		logger.Error("failed-to-lookup-created-container-in-garden", err)
		return nil, err
	}

	return c.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
}

func (c *gardenContainerClient) Finalize(logger lager.Logger, creatingContainer db.CreatingContainer) (Container, bool, error) {
	logger = logger.WithData(lager.Data{"container": creatingContainer.Handle()})
	logger.Debug("found-creating-container-in-db")

	gardenContainer, err := c.gardenClient.Lookup(creatingContainer.Handle())
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

		_ = c.gardenClient.Destroy(creatingContainer.Handle())

		return nil, false, err
	}

	logger.Debug("created-container-in-db")

	container, err := c.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
	if err != nil {
		return nil, false, err
	}

	return container, true, nil
}

func (c *gardenContainerClient) FetchImage(
	logger lager.Logger,
	provider ContainerProvider,
	spec ContainerSpec,
	delegate ImageFetchingDelegate,
	resourceTypes creds.VersionedResourceTypes,
) (Image, error) {
	worker := NewGardenWorker(
		c.gardenClient,
		c.baggageclaimClient,
		provider,
		c.volumeClient,
		c.worker,
		c.clock,
	)

	return c.imageFactory.GetImage(
		logger,
		worker,
		c.volumeClient,
		spec.ImageSpec,
		spec.TeamID,
		delegate,
		resourceTypes,
	)
}

func (c *gardenContainerClient) Create(
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

	gardenContainer, err := c.createGardenContainer(
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

		_ = c.gardenClient.Destroy(creatingContainer.Handle())

		return nil, err
	}

	logger.Debug("created-container-in-db")

	return c.constructGardenWorkerContainer(
		logger,
		createdContainer,
		gardenContainer,
	)
}

func (c *gardenContainerClient) Find(
	logger lager.Logger,
	createdContainer db.CreatedContainer,
	createdVolumes []db.CreatedVolume,
) (Container, bool, error) {

	gardenContainer, err := c.gardenClient.Lookup(createdContainer.Handle())
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
		c.gardenClient,
		c.volumeClient,
		c.worker.Name(),
	)

	if err != nil {
		logger.Error("failed-to-construct-container", err)
		return nil, false, err
	}

	return container, true, nil
}

func (c *gardenContainerClient) constructGardenWorkerContainer(
	logger lager.Logger,
	createdContainer db.CreatedContainer,
	gardenContainer garden.Container,
) (Container, error) {
	createdVolumes, err := c.dbVolumeFactory.FindVolumesForContainer(createdContainer)
	if err != nil {
		logger.Error("failed-to-find-container-volumes", err)
		return nil, err
	}

	return newGardenWorkerContainer(
		logger,
		gardenContainer,
		createdContainer,
		createdVolumes,
		c.gardenClient,
		c.volumeClient,
		c.worker.Name(),
	)
}

func (c *gardenContainerClient) createGardenContainer(
	logger lager.Logger,
	creatingContainer db.CreatingContainer,
	spec ContainerSpec,
	fetchedImage FetchedImage,
	provider ContainerProvider,
) (garden.Container, error) {
	volumeMounts := []VolumeMount{}

	scratchVolume, err := c.volumeClient.FindOrCreateVolumeForContainer(
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

	if spec.Dir != "" && !c.anyMountTo(spec.Dir, spec.Inputs) {
		workdirVolume, volumeErr := c.volumeClient.FindOrCreateVolumeForContainer(
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
		c.gardenClient,
		c.baggageclaimClient,
		provider,
		c.volumeClient,
		c.worker,
		c.clock,
	)

	for _, inputSource := range spec.Inputs {
		var inputVolume Volume

		localVolume, found, err := inputSource.Source().VolumeOn(worker)
		if err != nil {
			return nil, err
		}

		if found {
			inputVolume, err = c.volumeClient.FindOrCreateCOWVolumeForContainer(
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
			inputVolume, err = c.volumeClient.FindOrCreateVolumeForContainer(
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
		outVolume, volumeErr := c.volumeClient.FindOrCreateVolumeForContainer(
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

	if c.worker.HTTPProxyURL() != "" {
		env = append(env, fmt.Sprintf("http_proxy=%s", c.worker.HTTPProxyURL()))
	}

	if c.worker.HTTPSProxyURL() != "" {
		env = append(env, fmt.Sprintf("https_proxy=%s", c.worker.HTTPSProxyURL()))
	}

	if c.worker.NoProxy() != "" {
		env = append(env, fmt.Sprintf("no_proxy=%s", c.worker.NoProxy()))
	}

	return c.gardenClient.Create(garden.ContainerSpec{
		Handle:     creatingContainer.Handle(),
		RootFSPath: fetchedImage.URL,
		Privileged: fetchedImage.Privileged,
		BindMounts: bindMounts,
		Env:        env,
		Properties: gardenProperties,
	})
}

func (c *gardenContainerClient) anyMountTo(path string, inputs []InputSource) bool {
	for _, input := range inputs {
		if input.DestinationPath() == path {
			return true
		}
	}

	return false
}
