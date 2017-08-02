package worker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/baggageclaim"
	"github.com/cppforlife/go-semi-semantic/version"
)

var ErrUnsupportedResourceType = errors.New("unsupported resource type")
var ErrIncompatiblePlatform = errors.New("incompatible platform")
var ErrMismatchedTags = errors.New("mismatched tags")
var ErrNoVolumeManager = errors.New("worker does not support volume management")
var ErrTeamMismatch = errors.New("mismatched team")
var ErrNotImplemented = errors.New("Not implemented")

type MalformedMetadataError struct {
	UnmarshalError error
}

func (err MalformedMetadataError) Error() string {
	return fmt.Sprintf("malformed image metadata: %s", err.UnmarshalError)
}

const ephemeralPropertyName = "concourse:ephemeral"
const volumePropertyName = "concourse:volumes"
const volumeMountsPropertyName = "concourse:volume-mounts"
const userPropertyName = "user"
const RawRootFSScheme = "raw"
const ImageMetadataFile = "metadata.json"

//go:generate counterfeiter . Worker

type Worker interface {
	Client

	ActiveContainers() int

	Description() string
	Name() string
	ResourceTypes() []atc.WorkerResourceType
	Tags() atc.Tags
	Uptime() time.Duration
	IsOwnedByTeam() bool
	IsVersionCompatible(lager.Logger, *version.Version) bool

	FindVolumeForResourceCache(logger lager.Logger, resourceCache *db.UsedResourceCache) (Volume, bool, error)
	FindVolumeForTaskCache(lager.Logger, int, int, string, string) (Volume, bool, error)

	GardenClient() garden.Client
	BaggageclaimClient() baggageclaim.Client
}

type gardenWorker struct {
	containerProviderFactory ContainerProviderFactory
	gardenClient             garden.Client

	volumeClient       VolumeClient
	baggageclaimClient baggageclaim.Client

	provider WorkerProvider

	clock clock.Clock

	activeContainers int
	resourceTypes    []atc.WorkerResourceType
	platform         string
	tags             atc.Tags
	teamID           int
	name             string
	startTime        int64
	version          *string
}

func NewGardenWorker(
	containerProviderFactory ContainerProviderFactory,
	volumeClient VolumeClient,
	provider WorkerProvider,
	clock clock.Clock,
	activeContainers int,
	resourceTypes []atc.WorkerResourceType,
	platform string,
	tags atc.Tags,
	teamID int,
	name string,
	startTime int64,
	version *string,
	gardenClient garden.Client,
	baggageclaimClient baggageclaim.Client,
) Worker {
	return &gardenWorker{
		containerProviderFactory: containerProviderFactory,
		gardenClient:             gardenClient,

		volumeClient:       volumeClient,
		baggageclaimClient: baggageclaimClient,

		provider:         provider,
		clock:            clock,
		activeContainers: activeContainers,
		resourceTypes:    resourceTypes,
		platform:         platform,
		tags:             tags,
		teamID:           teamID,
		name:             name,
		startTime:        startTime,
		version:          version,
	}
}

func (worker *gardenWorker) GardenClient() garden.Client {
	return worker.gardenClient
}

func (worker *gardenWorker) BaggageclaimClient() baggageclaim.Client {
	return worker.baggageclaimClient
}

func (worker *gardenWorker) IsVersionCompatible(logger lager.Logger, comparedVersion *version.Version) bool {
	if comparedVersion == nil {
		return true
	}

	logger = logger.Session("check-version", lager.Data{
		"want-worker-version": comparedVersion.String(),
		"have-worker-version": worker.version,
	})

	if worker.version == nil {
		logger.Info("empty-worker-version")
		return false
	}

	v, err := version.NewVersionFromString(*worker.version)
	if err != nil {
		logger.Error("failed-to-parse-version", err)
		return false
	}

	switch v.Release.Compare(comparedVersion.Release) {
	case 0:
		return true
	case -1:
		return false
	default:
		if v.Release.Components[0].Compare(comparedVersion.Release.Components[0]) == 0 {
			return true
		}

		return false
	}
}

func (worker *gardenWorker) FindResourceTypeByPath(path string) (atc.WorkerResourceType, bool) {
	for _, rt := range worker.resourceTypes {
		if path == rt.Image {
			return rt, true
		}
	}

	return atc.WorkerResourceType{}, false
}

func (worker *gardenWorker) FindVolumeForResourceCache(logger lager.Logger, resourceCache *db.UsedResourceCache) (Volume, bool, error) {
	return worker.volumeClient.FindVolumeForResourceCache(logger, resourceCache)
}

func (worker *gardenWorker) FindVolumeForTaskCache(logger lager.Logger, teamID int, jobID int, stepName string, path string) (Volume, bool, error) {
	return worker.volumeClient.FindVolumeForTaskCache(logger, teamID, jobID, stepName, path)
}

func (worker *gardenWorker) LookupVolume(logger lager.Logger, handle string) (Volume, bool, error) {
	return worker.volumeClient.LookupVolume(logger, handle)
}

func (worker *gardenWorker) FindOrCreateContainer(
	logger lager.Logger,
	cancel <-chan os.Signal,
	delegate ImageFetchingDelegate,
	user db.ResourceUser,
	owner db.ContainerOwner,
	metadata db.ContainerMetadata,
	spec ContainerSpec,
	resourceTypes creds.VersionedResourceTypes,
) (Container, error) {
	containerProvider := worker.containerProviderFactory.ContainerProviderFor(worker)

	return containerProvider.FindOrCreateContainer(
		logger,
		cancel,
		user,
		owner,
		delegate,
		metadata,
		spec,
		resourceTypes,
	)
}

func (worker *gardenWorker) FindContainerByHandle(logger lager.Logger, teamID int, handle string) (Container, bool, error) {
	containerProvider := worker.containerProviderFactory.ContainerProviderFor(worker)
	return containerProvider.FindCreatedContainerByHandle(logger, handle, teamID)
}

func (worker *gardenWorker) ActiveContainers() int {
	return worker.activeContainers
}

func (worker *gardenWorker) Satisfying(logger lager.Logger, spec WorkerSpec, resourceTypes creds.VersionedResourceTypes) (Worker, error) {
	if spec.TeamID != worker.teamID && worker.teamID != 0 {
		return nil, ErrTeamMismatch
	}

	if spec.ResourceType != "" {
		underlyingType := determineUnderlyingTypeName(spec.ResourceType, resourceTypes)

		matchedType := false
		for _, t := range worker.resourceTypes {
			if t.Type == underlyingType {
				matchedType = true
				break
			}
		}

		if !matchedType {
			return nil, ErrUnsupportedResourceType
		}
	}

	if spec.Platform != "" {
		if spec.Platform != worker.platform {
			return nil, ErrIncompatiblePlatform
		}
	}

	if !worker.tagsMatch(spec.Tags) {
		return nil, ErrMismatchedTags
	}

	return worker, nil
}

func determineUnderlyingTypeName(typeName string, resourceTypes creds.VersionedResourceTypes) string {
	resourceTypesMap := make(map[string]creds.VersionedResourceType)
	for _, resourceType := range resourceTypes {
		resourceTypesMap[resourceType.Name] = resourceType
	}
	underlyingTypeName := typeName
	underlyingType, ok := resourceTypesMap[underlyingTypeName]
	for ok {
		underlyingTypeName = underlyingType.Type
		underlyingType, ok = resourceTypesMap[underlyingTypeName]
		delete(resourceTypesMap, underlyingTypeName)
	}
	return underlyingTypeName
}

func (worker *gardenWorker) AllSatisfying(logger lager.Logger, spec WorkerSpec, resourceTypes creds.VersionedResourceTypes) ([]Worker, error) {
	return nil, ErrNotImplemented
}

func (worker *gardenWorker) RunningWorkers(logger lager.Logger) ([]Worker, error) {
	return nil, ErrNotImplemented
}

func (worker *gardenWorker) GetWorker(logger lager.Logger, name string) (Worker, error) {
	return nil, ErrNotImplemented
}

func (worker *gardenWorker) Description() string {
	messages := []string{
		fmt.Sprintf("platform '%s'", worker.platform),
	}

	for _, tag := range worker.tags {
		messages = append(messages, fmt.Sprintf("tag '%s'", tag))
	}

	return strings.Join(messages, ", ")
}

func (worker *gardenWorker) Name() string {
	return worker.name
}

func (worker *gardenWorker) ResourceTypes() []atc.WorkerResourceType {
	return worker.resourceTypes
}

func (worker *gardenWorker) Tags() atc.Tags {
	return worker.tags
}

func (worker *gardenWorker) IsOwnedByTeam() bool {
	return worker.teamID != 0
}

func (worker *gardenWorker) Uptime() time.Duration {
	return worker.clock.Since(time.Unix(worker.startTime, 0))
}

func (worker *gardenWorker) tagsMatch(tags []string) bool {
	if len(worker.tags) > 0 && len(tags) == 0 {
		tags = []string{"default"}
	}

insert_coin:
	for _, stag := range tags {
		for _, wtag := range worker.tags {
			if stag == wtag {
				continue insert_coin
			}
		}

		return false
	}

	return true
}

type artifactDestination struct {
	destination Volume
}

func (wad *artifactDestination) StreamIn(path string, tarStream io.Reader) error {
	return wad.destination.StreamIn(path, tarStream)
}
