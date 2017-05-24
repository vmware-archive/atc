package resource

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/worker"
)

//go:generate counterfeiter . FetchSourceProviderFactory

type FetchSourceProviderFactory interface {
	NewFetchSourceProvider(
		logger lager.Logger,
		session Session,
		metadata Metadata,
		tags atc.Tags,
		teamID int,
		resourceTypes atc.VersionedResourceTypes,
		resourceInstance ResourceInstance,
		resourceOptions ResourceOptions,
		imageFetchingDelegate worker.ImageFetchingDelegate,
	) FetchSourceProvider
}

//go:generate counterfeiter . FetchSourceProvider

type FetchSourceProvider interface {
	Get() (FetchSource, error)
}

//go:generate counterfeiter . FetchSource

type FetchSource interface {
	LockName() (string, error)
	FindInitialized() (VersionedSource, bool, error)
	Initialize(signals <-chan os.Signal, ready chan<- struct{}) (VersionedSource, error)
}

type fetchSourceProviderFactory struct {
	workerClient           worker.Client
	dbResourceCacheFactory dbng.ResourceCacheFactory
}

func NewFetchSourceProviderFactory(
	workerClient worker.Client,
	dbResourceCacheFactory dbng.ResourceCacheFactory,
) FetchSourceProviderFactory {
	return &fetchSourceProviderFactory{
		workerClient:           workerClient,
		dbResourceCacheFactory: dbResourceCacheFactory,
	}
}

func (f *fetchSourceProviderFactory) NewFetchSourceProvider(
	logger lager.Logger,
	session Session,
	metadata Metadata,
	tags atc.Tags,
	teamID int,
	resourceTypes atc.VersionedResourceTypes,
	resourceInstance ResourceInstance,
	resourceOptions ResourceOptions,
	imageFetchingDelegate worker.ImageFetchingDelegate,
) FetchSourceProvider {
	return &fetchSourceProvider{
		logger:                 logger,
		session:                session,
		metadata:               metadata,
		tags:                   tags,
		teamID:                 teamID,
		resourceTypes:          resourceTypes,
		resourceInstance:       resourceInstance,
		resourceOptions:        resourceOptions,
		imageFetchingDelegate:  imageFetchingDelegate,
		workerClient:           f.workerClient,
		dbResourceCacheFactory: f.dbResourceCacheFactory,
	}
}

type fetchSourceProvider struct {
	logger                 lager.Logger
	session                Session
	metadata               Metadata
	tags                   atc.Tags
	teamID                 int
	resourceTypes          atc.VersionedResourceTypes
	resourceInstance       ResourceInstance
	resourceOptions        ResourceOptions
	workerClient           worker.Client
	imageFetchingDelegate  worker.ImageFetchingDelegate
	dbResourceCacheFactory dbng.ResourceCacheFactory
}

func (f *fetchSourceProvider) Get() (FetchSource, error) {
	resourceSpec := worker.WorkerSpec{
		ResourceType: string(f.resourceOptions.ResourceType()),
		Tags:         f.tags,
		TeamID:       f.teamID,
	}

	chosenWorker, err := f.workerClient.Satisfying(f.logger.Session("fetch-source-provider"), resourceSpec, f.resourceTypes)
	if err != nil {
		f.logger.Error("no-workers-satisfying-spec", err)
		return nil, err
	}

	resourceCache, err := f.dbResourceCacheFactory.FindOrCreateResourceCache(
		f.logger,
		f.resourceInstance.ResourceUser(),
		string(f.resourceOptions.ResourceType()),
		f.resourceOptions.Version(),
		f.resourceOptions.Source(),
		f.resourceOptions.Params(),
		f.resourceTypes,
	)
	if err != nil {
		f.logger.Error("failed-to-get-resource-cache", err, lager.Data{"user": f.resourceInstance.ResourceUser()})
		return nil, err
	}

	f.logger.Debug("initializing-resource-instance-fetch-source", lager.Data{"resource-cache": resourceCache})

	return NewResourceInstanceFetchSource(
		f.logger,
		resourceCache,
		f.resourceInstance,
		chosenWorker,
		f.resourceOptions,
		f.resourceTypes,
		f.tags,
		f.teamID,
		f.session,
		f.metadata,
		f.imageFetchingDelegate,
		f.dbResourceCacheFactory,
	), nil
}
