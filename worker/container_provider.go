package worker

import (
	"context"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
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
		WorkerName:      dbWorker.Name(),
		Clock:           clock,
		Client: &gardenContainerClient{
			gardenClient:       gardenClient,
			baggageclaimClient: baggageclaimClient,
			volumeClient:       volumeClient,
			imageFactory:       imageFactory,
			dbVolumeFactory:    dbVolumeFactory,
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
		Owner db.ContainerOwner,
		delegate ImageFetchingDelegate,
		metadata db.ContainerMetadata,
		spec ContainerSpec,
		resourceTypes creds.VersionedResourceTypes,
	) (Container, error)
}

type ContainerClient interface {
	Found(lager.Logger, db.CreatedContainer) (Container, error)
	Finalize(lager.Logger, db.CreatingContainer) (Container, bool, error)
	FetchImage(
		logger lager.Logger,
		provider ContainerProvider,
		spec ContainerSpec,
		delegate ImageFetchingDelegate,
		resourceTypes creds.VersionedResourceTypes,
	) (Image, error)
	Create(
		ctx context.Context,
		logger lager.Logger,
		provider ContainerProvider,
		creatingContainer db.CreatingContainer,
		spec ContainerSpec,
		image Image,
	) (Container, error)
	Find(lager.Logger, db.CreatedContainer, []db.CreatedVolume) (Container, bool, error)
}

type DbContainerProvider struct {
	DbVolumeFactory db.VolumeFactory
	DbTeamFactory   db.TeamFactory
	LockFactory     lock.LockFactory
	WorkerName      string
	Clock           clock.Clock

	Client ContainerClient
}

func (p *DbContainerProvider) Find(
	spec ContainerSpec,
	owner db.ContainerOwner,
) (db.CreatingContainer, db.CreatedContainer, error) {
	return p.DbTeamFactory.GetByID(spec.TeamID).FindContainerOnWorker(
		p.WorkerName,
		owner,
	)
}

func (p *DbContainerProvider) Create(
	spec ContainerSpec,
	owner db.ContainerOwner,
	metadata db.ContainerMetadata,
) (db.CreatingContainer, error) {
	return p.DbTeamFactory.GetByID(spec.TeamID).CreateContainer(
		p.WorkerName,
		owner,
		metadata,
	)
}

func (p *DbContainerProvider) Lock(logger lager.Logger, id int) (lock.Lock, bool, error) {
	return p.LockFactory.Acquire(logger, lock.NewContainerCreatingLockID(id))
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
		creatingContainer, createdContainer, err := p.Find(spec, owner)
		if err != nil {
			logger.Error("failed-to-find-container-in-db", err)
			return nil, err
		}

		if createdContainer != nil {
			logger = logger.WithData(lager.Data{"container": createdContainer.Handle()})
			logger.Debug("found-created-container-in-db")
			return p.Client.Found(logger, createdContainer)
		}

		var image Image

		if creatingContainer != nil {
			container, done, err := p.Client.Finalize(logger, creatingContainer)
			if err != nil {
				return nil, err
			}
			if done {
				return container, nil
			}
		} else {
			image, err = p.Client.FetchImage(logger, p, spec, delegate, resourceTypes)
			if err != nil {
				logger.Error("failed-to-get-image-for-container", err)
				return nil, err
			}

			creatingContainer, err = p.Create(
				spec,
				owner,
				metadata,
			)
			if err != nil {
				logger.Error("failed-to-create-container-in-db", err)
				return nil, err
			}
		}

		if image == nil {
			image, err = p.Client.FetchImage(logger, p, spec, delegate, resourceTypes)
			if err != nil {
				logger.Error("failed-to-get-image-for-container", err)
				return nil, err
			}
		}

		logger = logger.WithData(lager.Data{"container": creatingContainer.Handle()})
		logger.Debug("created-creating-container-in-db")

		lock, acquired, err := p.Lock(logger, creatingContainer.ID())
		if err != nil {
			logger.Error("failed-to-acquire-container-creating-lock", err)
			return nil, err
		}

		if !acquired {
			p.Clock.Sleep(creatingContainerRetryDelay)
			continue
		}

		defer lock.Release()

		return p.Client.Create(ctx, logger, p, creatingContainer, spec, image)
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

	return p.Client.Find(logger, createdContainer, createdVolumes)
}
