package resource

import (
	"code.cloudfoundry.org/clock"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/worker"
)

//go:generate counterfeiter . FetcherFactory

type FetcherFactory interface {
	FetcherFor(workerClient worker.Client) Fetcher
}

func NewFetcherFactory(
	lockFactory lock.LockFactory,
	clock clock.Clock,
	dbResourceCacheFactory db.ResourceCacheFactory,
	orchestrator resourceOrchestrator,
) FetcherFactory {
	return &fetcherFactory{
		lockFactory: lockFactory,
		clock:       clock,
		dbResourceCacheFactory: dbResourceCacheFactory,
		orchestrator:           orchestrator,
	}
}

type fetcherFactory struct {
	lockFactory            lock.LockFactory
	clock                  clock.Clock
	dbResourceCacheFactory db.ResourceCacheFactory
	orchestrator           resourceOrchestrator
}

func (f *fetcherFactory) FetcherFor(workerClient worker.Client) Fetcher {
	return NewFetcher(
		f.clock,
		f.lockFactory,
		NewFetchSourceProviderFactory(workerClient, f.dbResourceCacheFactory, f.orchestrator),
	)
}
