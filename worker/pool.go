package worker

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
)

//go:generate counterfeiter . WorkerProvider

type WorkerProvider interface {
	RunningWorkers(lager.Logger) ([]Worker, error)

	AllSatisfying(logger lager.Logger, spec WorkerSpec, resourceTypes creds.VersionedResourceTypes) ([]Worker, error)

	FindWorkerForContainer(
		logger lager.Logger,
		teamID int,
		handle string,
	) (Worker, bool, error)

	FindWorkerForContainerByOwner(
		logger lager.Logger,
		teamID int,
		owner db.ContainerOwner,
	) (Worker, bool, error)

	NewGardenWorker(
		logger lager.Logger,
		tikTok clock.Clock,
		savedWorker db.Worker,
	) Worker
}

var (
	ErrNoWorkers = errors.New("no workers")
)

type NoCompatibleWorkersError struct {
	Spec    WorkerSpec
	Workers []Worker
}

func (err NoCompatibleWorkersError) Error() string {
	availableWorkers := ""
	for _, worker := range err.Workers {
		availableWorkers += "\n  - " + worker.Description()
	}

	return fmt.Sprintf(
		"no workers satisfying: %s\n\navailable workers: %s",
		err.Spec.Description(),
		availableWorkers,
	)
}

type pool struct {
	provider WorkerProvider

	rand     *rand.Rand
	strategy ContainerPlacementStrategy
}

func NewPool(provider WorkerProvider, strategy ContainerPlacementStrategy) Client {
	return &pool{
		provider: provider,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		strategy: strategy,
	}
}

func (pool *pool) Satisfying(logger lager.Logger, spec WorkerSpec, resourceTypes creds.VersionedResourceTypes) (Worker, error) {
	compatibleWorkers, err := pool.provider.AllSatisfying(logger, spec, resourceTypes)
	if err != nil {
		return nil, err
	}
	randomWorker := compatibleWorkers[pool.rand.Intn(len(compatibleWorkers))]
	return randomWorker, nil
}

func (pool *pool) FindOrCreateContainer(
	ctx context.Context,
	logger lager.Logger,
	delegate ImageFetchingDelegate,
	owner db.ContainerOwner,
	metadata db.ContainerMetadata,
	spec ContainerSpec,
	resourceTypes creds.VersionedResourceTypes,
) (Container, error) {
	worker, found, err := pool.provider.FindWorkerForContainerByOwner(
		logger.Session("find-worker"),
		spec.TeamID,
		owner,
	)
	if err != nil {
		return nil, err
	}

	if !found {
		compatibleWorkers, err := pool.provider.AllSatisfying(logger, spec.WorkerSpec(), resourceTypes)
		if err != nil {
			return nil, err
		}

		worker, err = pool.strategy.Choose(compatibleWorkers, spec)
		if err != nil {
			return nil, err
		}
	}

	return worker.FindOrCreateContainer(
		ctx,
		logger,
		delegate,
		owner,
		metadata,
		spec,
		resourceTypes,
	)
}

func (pool *pool) FindContainerByHandle(logger lager.Logger, teamID int, handle string) (Container, bool, error) {
	worker, found, err := pool.provider.FindWorkerForContainer(
		logger.Session("find-worker"),
		teamID,
		handle,
	)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	return worker.FindContainerByHandle(logger, teamID, handle)
}

func (*pool) LookupVolume(lager.Logger, string) (Volume, bool, error) {
	return nil, false, errors.New("LookupVolume not implemented for pool")
}
