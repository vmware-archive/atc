package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/baggageclaim"
)

var ErrUnsupportedResourceType = errors.New("unsupported resource type")
var ErrIncompatiblePlatform = errors.New("incompatible platform")
var ErrMismatchedTags = errors.New("mismatched tags")
var ErrNoVolumeManager = errors.New("worker does not support volume management")
var ErrTeamMismatch = errors.New("mismatched team")

type MalformedMetadataError struct {
	UnmarshalError error
}

func (err MalformedMetadataError) Error() string {
	return fmt.Sprintf("malformed image metadata: %s", err.UnmarshalError)
}

const containerKeepalive = 30 * time.Second
const ContainerTTL = 5 * time.Minute

const VolumeTTL = 5 * time.Minute

const ephemeralPropertyName = "concourse:ephemeral"
const volumePropertyName = "concourse:volumes"
const volumeMountsPropertyName = "concourse:volume-mounts"
const userPropertyName = "user"
const RawRootFSScheme = "raw"

//go:generate counterfeiter . Worker

type Worker interface {
	Client

	ActiveContainers() int

	Description() string
	Name() string
	Uptime() time.Duration
	IsOwnedByTeam() bool
}

//go:generate counterfeiter . GardenWorkerDB

type GardenWorkerDB interface {
	CreateContainer(container db.Container, ttl time.Duration, maxLifetime time.Duration, volumeHandles []string) (db.SavedContainer, error)
	GetContainer(handle string) (db.SavedContainer, bool, error)
	UpdateExpiresAtOnContainer(handle string, ttl time.Duration) error
	ReapContainer(string) error
	GetPipelineByID(pipelineID int) (db.SavedPipeline, error)
	InsertVolume(db.Volume) error
	SetVolumeTTLAndSizeInBytes(string, time.Duration, int64) error
	GetVolumeTTL(string) (time.Duration, bool, error)
	GetVolumesByIdentifier(db.VolumeIdentifier) ([]db.SavedVolume, error)
}

type gardenWorker struct {
	gardenClient       garden.Client
	baggageclaimClient baggageclaim.Client
	volumeClient       VolumeClient
	volumeFactory      VolumeFactory
	pipelineDBFactory  db.PipelineDBFactory
	imageFactory       ImageFactory

	db       GardenWorkerDB
	provider WorkerProvider

	clock clock.Clock

	activeContainers int
	resourceTypes    []atc.WorkerResourceType
	platform         string
	tags             atc.Tags
	teamID           int
	name             string
	startTime        int64
	httpProxyURL     string
	httpsProxyURL    string
	noProxy          string
}

func NewGardenWorker(
	gardenClient garden.Client,
	baggageclaimClient baggageclaim.Client,
	volumeClient VolumeClient,
	volumeFactory VolumeFactory,
	imageFactory ImageFactory,
	pipelineDBFactory db.PipelineDBFactory,
	db GardenWorkerDB,
	provider WorkerProvider,
	clock clock.Clock,
	activeContainers int,
	resourceTypes []atc.WorkerResourceType,
	platform string,
	tags atc.Tags,
	teamID int,
	name string,
	startTime int64,
	httpProxyURL string,
	httpsProxyURL string,
	noProxy string,
) Worker {
	return &gardenWorker{
		gardenClient:       gardenClient,
		baggageclaimClient: baggageclaimClient,
		volumeClient:       volumeClient,
		volumeFactory:      volumeFactory,
		imageFactory:       imageFactory,
		db:                 db,
		provider:           provider,
		clock:              clock,
		pipelineDBFactory:  pipelineDBFactory,
		activeContainers:   activeContainers,
		resourceTypes:      resourceTypes,
		platform:           platform,
		tags:               tags,
		teamID:             teamID,
		name:               name,
		startTime:          startTime,
		httpProxyURL:       httpProxyURL,
		httpsProxyURL:      httpsProxyURL,
		noProxy:            noProxy,
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

func (worker *gardenWorker) FindVolume(logger lager.Logger, volumeSpec VolumeSpec) (Volume, bool, error) {
	return worker.volumeClient.FindVolume(logger, volumeSpec)
}

func (worker *gardenWorker) CreateVolume(logger lager.Logger, volumeSpec VolumeSpec, teamID int) (Volume, error) {
	return worker.volumeClient.CreateVolume(logger, volumeSpec, teamID)
}

func (worker *gardenWorker) ListVolumes(logger lager.Logger, properties VolumeProperties) ([]Volume, error) {
	return worker.volumeClient.ListVolumes(logger, properties)
}

func (worker *gardenWorker) LookupVolume(logger lager.Logger, handle string) (Volume, bool, error) {
	return worker.volumeClient.LookupVolume(logger, handle)
}

func (worker *gardenWorker) getImage(
	logger lager.Logger,
	imageSpec ImageSpec,
	teamID int,
	cancel <-chan os.Signal,
	delegate ImageFetchingDelegate,
	id Identifier,
	metadata Metadata,
	resourceTypes atc.ResourceTypes,
) (Volume, ImageMetadata, atc.Version, string, error) {
	// convert custom resource type from pipeline config into image_resource
	updatedResourceTypes := resourceTypes
	imageResource := imageSpec.ImageResource
	for _, resourceType := range resourceTypes {
		if resourceType.Name == imageSpec.ResourceType {
			updatedResourceTypes = resourceTypes.Without(imageSpec.ResourceType)
			imageResource = &atc.ImageResource{
				Source: resourceType.Source,
				Type:   resourceType.Type,
			}
		}
	}

	imageVolume := imageSpec.ImageVolumeAndMetadata.Volume
	imageMetadataReader := imageSpec.ImageVolumeAndMetadata.MetadataReader
	var version atc.Version

	// 'image_resource:' in task
	if imageResource != nil {
		var err error
		image := worker.imageFactory.NewImage(
			logger.Session("image"),
			cancel,
			*imageResource,
			id,
			metadata,
			worker.tags,
			teamID,
			updatedResourceTypes,
			worker,
			delegate,
			imageSpec.Privileged,
		)

		imageVolume, imageMetadataReader, version, err = image.Fetch()
		if err != nil {
			return nil, ImageMetadata{}, nil, "", err
		}
	}

	// use image artifact from previous step in subsequent task within the same job
	if imageVolume != nil {
		metadata, err := loadMetadata(imageMetadataReader)
		if err != nil {
			return nil, ImageMetadata{}, nil, "", err
		}

		imageURL := url.URL{
			Scheme: RawRootFSScheme,
			Path:   path.Join(imageVolume.Path(), "rootfs"),
		}

		return imageVolume, metadata, version, imageURL.String(), nil
	}

	// built-in resource type specified in step
	if imageSpec.ResourceType != "" {
		rootFSURL, volume, resourceTypeVersion, err := worker.getBuiltInResourceTypeImage(logger, imageSpec.ResourceType, teamID)
		if err != nil {
			return nil, ImageMetadata{}, nil, "", err
		}

		return volume, ImageMetadata{}, resourceTypeVersion, rootFSURL, nil
	}

	// 'image:' in task
	return nil, ImageMetadata{}, nil, imageSpec.ImageURL, nil
}

func (worker *gardenWorker) ValidateResourceCheckVersion(container db.SavedContainer) (bool, error) {
	if container.Type != db.ContainerTypeCheck || container.CheckType == "" || container.ResourceTypeVersion == nil {
		return true, nil
	}

	if container.PipelineID > 0 {
		savedPipeline, err := worker.db.GetPipelineByID(container.PipelineID)
		if err != nil {
			return false, err
		}

		pipelineDB := worker.pipelineDBFactory.Build(savedPipeline)

		_, found, err := pipelineDB.GetResourceType(container.CheckType)
		if err != nil {
			return false, err
		}

		// this is custom resource type, do not validate version on worker
		if found {
			return true, nil
		}
	}

	for _, workerResourceType := range worker.resourceTypes {
		if container.CheckType == workerResourceType.Type && workerResourceType.Version == container.ResourceTypeVersion[container.CheckType] {
			return true, nil
		}
	}

	return false, nil
}

func (worker *gardenWorker) getBuiltInResourceTypeImage(
	logger lager.Logger,
	resourceTypeName string,
	teamID int,
) (string, Volume, atc.Version, error) {
	for _, t := range worker.resourceTypes {
		if t.Type == resourceTypeName {
			importVolumeSpec := VolumeSpec{
				Strategy: HostRootFSStrategy{
					Path:       t.Image,
					Version:    &t.Version,
					WorkerName: worker.Name(),
				},
				Privileged: true,
				Properties: VolumeProperties{},
				TTL:        0,
			}

			importVolume, found, err := worker.FindVolume(logger, importVolumeSpec)
			if err != nil {
				return "", nil, atc.Version{}, err
			}

			if !found {
				importVolume, err = worker.CreateVolume(logger, importVolumeSpec, 0)
				if err != nil {
					return "", nil, atc.Version{}, err
				}
			}

			defer importVolume.Release(nil)

			cowVolume, err := worker.CreateVolume(logger, VolumeSpec{
				Strategy: ContainerRootFSStrategy{
					Parent: importVolume,
				},
				Privileged: true,
				Properties: VolumeProperties{},
				TTL:        VolumeTTL,
			}, teamID)
			if err != nil {
				return "", nil, atc.Version{}, err
			}

			rootFSURL := url.URL{
				Scheme: RawRootFSScheme,
				Path:   cowVolume.Path(),
			}

			return rootFSURL.String(), cowVolume, atc.Version{resourceTypeName: t.Version}, nil
		}
	}

	return "", nil, atc.Version{}, ErrUnsupportedResourceType
}

func loadMetadata(tarReader io.ReadCloser) (ImageMetadata, error) {
	defer tarReader.Close()

	var imageMetadata ImageMetadata
	if err := json.NewDecoder(tarReader).Decode(&imageMetadata); err != nil {
		return ImageMetadata{}, MalformedMetadataError{
			UnmarshalError: err,
		}
	}

	return imageMetadata, nil
}

func (worker *gardenWorker) CreateContainer(
	logger lager.Logger,
	cancel <-chan os.Signal,
	delegate ImageFetchingDelegate,
	id Identifier,
	metadata Metadata,
	spec ContainerSpec,
	resourceTypes atc.ResourceTypes,
) (Container, error) {
	imageVolume, imageMetadata, resourceTypeVersion, imageURL, err := worker.getImage(
		logger,
		spec.ImageSpec,
		spec.TeamID,
		cancel,
		delegate,
		id,
		metadata,
		resourceTypes,
	)
	if err != nil {
		return nil, err
	}
	if imageVolume != nil {
		defer imageVolume.Release(nil)
	}

	volumeMounts := spec.Outputs
	for _, mount := range spec.Inputs {
		cowVolume, err := worker.CreateVolume(logger, VolumeSpec{
			Strategy: ContainerRootFSStrategy{
				Parent: mount.Volume,
			},
			Privileged: spec.ImageSpec.Privileged,
			TTL:        VolumeTTL,
		}, spec.TeamID)
		if err != nil {
			return nil, err
		}
		// release *after* container creation
		defer cowVolume.Release(nil)

		volumeMounts = append(volumeMounts, VolumeMount{
			Volume:    cowVolume,
			MountPath: mount.MountPath,
		})

		logger.Debug("created-cow-volume", lager.Data{
			"original-volume-handle": mount.Volume.Handle(),
			"cow-volume-handle":      cowVolume.Handle(),
		})
	}

	bindMounts := []garden.BindMount{}
	volumeHandles := []string{}
	volumeHandleMounts := map[string]string{}
	for _, mount := range volumeMounts {
		volumeHandles = append(volumeHandles, mount.Volume.Handle())
		bindMounts = append(bindMounts, garden.BindMount{
			SrcPath: mount.Volume.Path(),
			DstPath: mount.MountPath,
			Mode:    garden.BindMountModeRW,
		})
		volumeHandleMounts[mount.Volume.Handle()] = mount.MountPath
	}

	if imageVolume != nil {
		volumeHandles = append(volumeHandles, imageVolume.Handle())
	}

	gardenProperties := garden.Properties{userPropertyName: imageMetadata.User}
	if spec.User != "" {
		gardenProperties = garden.Properties{userPropertyName: spec.User}
	}

	if len(volumeHandles) > 0 {
		volumesJSON, err := json.Marshal(volumeHandles)
		if err != nil {
			return nil, err
		}

		gardenProperties[volumePropertyName] = string(volumesJSON)

		mountsJSON, err := json.Marshal(volumeHandleMounts)
		if err != nil {
			return nil, err
		}

		gardenProperties[volumeMountsPropertyName] = string(mountsJSON)
	}

	if spec.Ephemeral {
		gardenProperties[ephemeralPropertyName] = "true"
	}

	env := append(imageMetadata.Env, spec.Env...)

	if worker.httpProxyURL != "" {
		env = append(env, fmt.Sprintf("http_proxy=%s", worker.httpProxyURL))
	}

	if worker.httpsProxyURL != "" {
		env = append(env, fmt.Sprintf("https_proxy=%s", worker.httpsProxyURL))
	}

	if worker.noProxy != "" {
		env = append(env, fmt.Sprintf("no_proxy=%s", worker.noProxy))
	}

	gardenSpec := garden.ContainerSpec{
		BindMounts: bindMounts,
		Privileged: spec.ImageSpec.Privileged,
		Properties: gardenProperties,
		RootFSPath: imageURL,
		Env:        env,
	}

	gardenContainer, err := worker.gardenClient.Create(gardenSpec)
	if err != nil {
		return nil, err
	}

	metadata.WorkerName = worker.name
	metadata.Handle = gardenContainer.Handle()
	metadata.User = gardenSpec.Properties["user"]

	id.ResourceTypeVersion = resourceTypeVersion

	_, err = worker.db.CreateContainer(
		db.Container{
			ContainerIdentifier: db.ContainerIdentifier(id),
			ContainerMetadata:   db.ContainerMetadata(metadata),
		},
		ContainerTTL,
		worker.maxContainerLifetime(metadata),
		volumeHandles,
	)
	if err != nil {
		return nil, err
	}

	return newGardenWorkerContainer(
		logger,
		gardenContainer,
		worker.gardenClient,
		worker.baggageclaimClient,
		worker.db,
		worker.clock,
		worker.volumeFactory,
		worker.name,
	)
}

func (worker *gardenWorker) maxContainerLifetime(metadata Metadata) time.Duration {
	if metadata.Type == db.ContainerTypeCheck {
		uptime := worker.Uptime()
		switch {
		case uptime < 5*time.Minute:
			return 5 * time.Minute
		case uptime > 1*time.Hour:
			return 1 * time.Hour
		default:
			return uptime
		}
	}

	return time.Duration(0)
}

func (worker *gardenWorker) FindContainerForIdentifier(logger lager.Logger, id Identifier) (Container, bool, error) {
	containerInfo, found, err := worker.provider.FindContainerForIdentifier(id)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, found, nil
	}

	valid, err := worker.ValidateResourceCheckVersion(containerInfo)

	if err != nil {
		return nil, false, err
	}

	if !valid {
		logger.Debug("check-container-version-outdated", lager.Data{
			"container-handle": containerInfo.Handle,
			"worker-name":      containerInfo.WorkerName,
		})

		return nil, false, nil
	}

	container, found, err := worker.LookupContainer(logger, containerInfo.Handle)
	if err != nil {
		return nil, false, err
	}

	if !found {
		logger.Info("reaping-container-not-found-on-worker", lager.Data{
			"container-handle": containerInfo.Handle,
			"worker-name":      containerInfo.WorkerName,
		})

		err := worker.provider.ReapContainer(containerInfo.Handle)
		if err != nil {
			return nil, false, err
		}

		return nil, false, err
	}

	return container, found, nil
}

func (worker *gardenWorker) LookupContainer(logger lager.Logger, handle string) (Container, bool, error) {
	gardenContainer, err := worker.gardenClient.Lookup(handle)
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
		worker.gardenClient,
		worker.baggageclaimClient,
		worker.db,
		worker.clock,
		worker.volumeFactory,
		worker.name,
	)

	if err != nil {
		logger.Error("failed-to-construct-container", err)
		return nil, false, err
	}

	return container, true, nil
}

func (worker *gardenWorker) ActiveContainers() int {
	return worker.activeContainers
}

func (worker *gardenWorker) Satisfying(spec WorkerSpec, resourceTypes atc.ResourceTypes) (Worker, error) {
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

func determineUnderlyingTypeName(typeName string, resourceTypes atc.ResourceTypes) string {
	resourceTypesMap := make(map[string]atc.ResourceType)
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

func (worker *gardenWorker) AllSatisfying(spec WorkerSpec, resourceTypes atc.ResourceTypes) ([]Worker, error) {
	return nil, errors.New("Not implemented")
}

func (worker *gardenWorker) Workers() ([]Worker, error) {
	return nil, errors.New("Not implemented")
}

func (worker *gardenWorker) GetWorker(name string) (Worker, error) {
	return nil, errors.New("Not implemented")
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
