package exec

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/worker"
)

const taskProcessID = "task"
const taskProcessPropertyName = "concourse:task-process"
const taskExitStatusPropertyName = "concourse:exit-status"

// MissingInputsError is returned when any of the task's required inputs are
// missing.
type MissingInputsError struct {
	Inputs []string
}

// Error prints a human-friendly message listing the inputs that were missing.
func (err MissingInputsError) Error() string {
	return fmt.Sprintf("missing inputs: %s", strings.Join(err.Inputs, ", "))
}

type MissingTaskImageSourceError struct {
	SourceName string
}

func (err MissingTaskImageSourceError) Error() string {
	return fmt.Sprintf(`missing image artifact source: %s

make sure there's a corresponding 'get' step, or a task that produces it as an output`, err.SourceName)
}

type TaskImageSourceParametersError struct {
	Err error
}

func (err TaskImageSourceParametersError) Error() string {
	return fmt.Sprintf("failed to evaluate image resource parameters: %s", err.Err)
}

//go:generate counterfeiter . TaskConfigSource

type TaskConfigSource interface {
	GetTaskConfig() (atc.TaskConfig, error)
}

type FetchConfigActionTaskConfigSource struct {
	Action FetchConfigResultAction
}

func (s *FetchConfigActionTaskConfigSource) GetTaskConfig() (atc.TaskConfig, error) {
	taskConfig := s.Action.Result()
	return taskConfig, nil
}

//go:generate counterfeiter . TaskBuildEventsDelegate

type TaskBuildEventsDelegate interface {
	Initializing(lager.Logger, atc.TaskConfig)
	Starting(lager.Logger, atc.TaskConfig)
}

// TaskAction executes a TaskConfig, whose inputs will be fetched from the
// worker.ArtifactRepository and outputs will be added to the worker.ArtifactRepository.
type TaskAction struct {
	privileged    Privileged
	configSource  TaskConfigSource
	tags          atc.Tags
	inputMapping  map[string]string
	outputMapping map[string]string

	// TODO: replace with RootFSSource
	artifactsRoot     string
	imageArtifactName string

	buildEventsDelegate TaskBuildEventsDelegate
	buildStepDelegate   BuildStepDelegate
	workerPool          worker.Client
	teamID              int
	buildID             int
	jobID               int
	stepName            string
	planID              atc.PlanID
	containerMetadata   db.ContainerMetadata

	resourceTypes creds.VersionedResourceTypes

	variables creds.Variables

	exitStatus ExitStatus
}

func NewTaskAction(
	privileged Privileged,
	configSource TaskConfigSource,
	tags atc.Tags,
	inputMapping map[string]string,
	outputMapping map[string]string,
	artifactsRoot string,
	imageArtifactName string,
	buildEventsDelegate TaskBuildEventsDelegate,
	buildStepDelegate BuildStepDelegate,
	workerPool worker.Client,
	teamID int,
	buildID int,
	jobID int,
	stepName string,
	planID atc.PlanID,
	containerMetadata db.ContainerMetadata,
	resourceTypes creds.VersionedResourceTypes,
	variables creds.Variables,
) *TaskAction {
	return &TaskAction{
		privileged:          privileged,
		configSource:        configSource,
		tags:                tags,
		inputMapping:        inputMapping,
		outputMapping:       outputMapping,
		artifactsRoot:       artifactsRoot,
		imageArtifactName:   imageArtifactName,
		buildEventsDelegate: buildEventsDelegate,
		buildStepDelegate:   buildStepDelegate,
		workerPool:          workerPool,
		teamID:              teamID,
		buildID:             buildID,
		jobID:               jobID,
		stepName:            stepName,
		planID:              planID,
		containerMetadata:   containerMetadata,
		resourceTypes:       resourceTypes,
		variables:           variables,
	}
}

// Run will first selects the worker based on the TaskConfig's platform, the
// TaskStep's tags, and prioritized by availability of volumes for the TaskConfig's
// inputs. Inputs that did not have volumes available on the worker will be streamed
// in to the container.
//
// If any inputs are not available in the worker.ArtifactRepository, MissingInputsError
// is returned.
//
// Once all the inputs are satisfies, the task's script will be executed, and
// the RunStep indicates that it's ready, and any signals will be forwarded to
// the script.
//
// If the script exits successfully, the outputs specified in the TaskConfig
// are registered with the worker.ArtifactRepository. If no outputs are specified, the
// task's entire working directory is registered as an ArtifactSource under the
// name of the task.
func (action *TaskAction) Run(
	logger lager.Logger,
	repository *worker.ArtifactRepository,

	// TODO: consider passing these as context
	signals <-chan os.Signal,
	ready chan<- struct{},
) error {
	config, err := action.configSource.GetTaskConfig()
	if err != nil {
		return err
	}

	action.buildEventsDelegate.Initializing(logger, config)

	containerSpec, err := action.containerSpec(logger, repository, config)
	if err != nil {
		return err
	}

	container, err := action.workerPool.FindOrCreateContainer(
		logger,
		signals,
		action.buildStepDelegate,
		db.NewBuildStepContainerOwner(action.buildID, action.planID),
		action.containerMetadata,
		containerSpec,
		action.resourceTypes,
	)
	if err != nil {
		return err
	}

	exitStatusProp, err := container.Property(taskExitStatusPropertyName)
	if err == nil {
		logger.Info("already-exited", lager.Data{"status": exitStatusProp})

		_, err = fmt.Sscanf(exitStatusProp, "%d", &action.exitStatus)
		if err != nil {
			return err
		}

		err = action.registerOutputs(logger, repository, config, container)
		if err != nil {
			return err
		}
		return nil
	}

	// for backwards compatibility with containers
	// that had their task process name set as property
	var processID string
	processID, err = container.Property(taskProcessPropertyName)
	if err != nil {
		processID = taskProcessID
	}

	processIO := garden.ProcessIO{
		Stdout: action.buildStepDelegate.Stdout(),
		Stderr: action.buildStepDelegate.Stderr(),
	}

	process, err := container.Attach(processID, processIO)
	if err == nil {
		logger.Info("already-running")
	} else {
		logger.Info("spawning")

		action.buildEventsDelegate.Starting(logger, config)

		process, err = container.Run(garden.ProcessSpec{
			ID: taskProcessID,

			Path: config.Run.Path,
			Args: config.Run.Args,

			Dir: path.Join(action.artifactsRoot, config.Run.Dir),
			TTY: &garden.TTYSpec{},
		}, processIO)
	}
	if err != nil {
		return err
	}

	logger.Info("attached")

	close(ready)

	exited := make(chan struct{})
	var processStatus int
	var processErr error

	go func() {
		processStatus, processErr = process.Wait()
		close(exited)
	}()

	select {
	case <-signals:
		err = action.registerOutputs(logger, repository, config, container)
		if err != nil {
			return err
		}

		err = container.Stop(false)
		if err != nil {
			logger.Error("stopping-container", err)
		}

		<-exited

		return ErrInterrupted

	case <-exited:
		if processErr != nil {
			return processErr
		}

		err = action.registerOutputs(logger, repository, config, container)
		if err != nil {
			return err
		}

		action.exitStatus = ExitStatus(processStatus)

		err = container.SetProperty(taskExitStatusPropertyName, fmt.Sprintf("%d", processStatus))
		if err != nil {
			return err
		}

		return nil
	}
}

// ExitStatus returns exit status of task script.
func (action *TaskAction) ExitStatus() ExitStatus {
	return action.exitStatus
}

func (action *TaskAction) containerSpec(logger lager.Logger, repository *worker.ArtifactRepository, config atc.TaskConfig) (worker.ContainerSpec, error) {
	imageSpec := worker.ImageSpec{
		Privileged: bool(action.privileged),
	}

	if action.imageArtifactName != "" {
		source, found := repository.SourceFor(worker.ArtifactName(action.imageArtifactName))
		if !found {
			return worker.ContainerSpec{}, MissingTaskImageSourceError{action.imageArtifactName}
		}

		imageSpec.ImageArtifactSource = source
		imageSpec.ImageArtifactName = worker.ArtifactName(action.imageArtifactName)
	} else if config.ImageResource != nil {
		imageSpec.ImageResource = &worker.ImageResource{
			Type:    config.ImageResource.Type,
			Source:  creds.NewSource(action.variables, config.ImageResource.Source),
			Params:  config.ImageResource.Params,
			Version: config.ImageResource.Version,
		}

	} else if config.RootfsURI != "" {
		imageSpec.ImageURL = config.RootfsURI
	}

	params, err := creds.NewTaskParams(action.variables, config.Params).Evaluate()
	if err != nil {
		return worker.ContainerSpec{}, err
	}

	containerSpec := worker.ContainerSpec{
		Platform:  config.Platform,
		Tags:      action.tags,
		TeamID:    action.teamID,
		ImageSpec: imageSpec,
		User:      config.Run.User,
		Dir:       action.artifactsRoot,
		Env:       action.envForParams(params),

		Inputs:  []worker.InputSource{},
		Outputs: worker.OutputPaths{},
	}

	var missingInputs []string
	for _, input := range config.Inputs {
		inputName := input.Name
		if sourceName, ok := action.inputMapping[inputName]; ok {
			inputName = sourceName
		}

		source, found := repository.SourceFor(worker.ArtifactName(inputName))
		if !found {
			missingInputs = append(missingInputs, inputName)
			continue
		}

		containerSpec.Inputs = append(containerSpec.Inputs, &taskInputSource{
			config:        input,
			source:        source,
			artifactsRoot: action.artifactsRoot,
		})
	}

	if len(missingInputs) > 0 {
		return worker.ContainerSpec{}, MissingInputsError{missingInputs}
	}

	for _, cacheConfig := range config.Caches {
		source := newTaskCacheSource(logger, action.teamID, action.jobID, action.stepName, cacheConfig.Path)
		containerSpec.Inputs = append(containerSpec.Inputs, &taskCacheInputSource{
			source:        source,
			artifactsRoot: action.artifactsRoot,
			cachePath:     cacheConfig.Path,
		})
	}

	for _, output := range config.Outputs {
		path := artifactsPath(output, action.artifactsRoot)
		containerSpec.Outputs[output.Name] = path
	}

	return containerSpec, nil
}

func (action *TaskAction) registerOutputs(logger lager.Logger, repository *worker.ArtifactRepository, config atc.TaskConfig, container worker.Container) error {
	volumeMounts := container.VolumeMounts()

	logger.Debug("registering-outputs", lager.Data{"outputs": config.Outputs})

	for _, output := range config.Outputs {
		outputName := output.Name
		if destinationName, ok := action.outputMapping[output.Name]; ok {
			outputName = destinationName
		}

		outputPath := artifactsPath(output, action.artifactsRoot)

		for _, mount := range volumeMounts {
			if mount.MountPath == outputPath {
				source := newTaskArtifactSource(logger, mount.Volume)
				repository.RegisterSource(worker.ArtifactName(outputName), source)
			}
		}
	}

	// Do not initialize caches for one-off builds
	if action.jobID != 0 {
		logger.Debug("initializing-caches", lager.Data{"caches": config.Caches})

		for _, cacheConfig := range config.Caches {
			for _, volumeMount := range volumeMounts {
				if volumeMount.MountPath == filepath.Join(action.artifactsRoot, cacheConfig.Path) {
					logger.Debug("initializing-cache", lager.Data{"path": volumeMount.MountPath})

					err := volumeMount.Volume.InitializeTaskCache(logger, action.jobID, action.stepName, cacheConfig.Path, bool(action.privileged))
					if err != nil {
						return err
					}

					continue
				}
			}
		}
	}

	return nil
}

func (TaskAction) envForParams(params map[string]string) []string {
	env := make([]string, 0, len(params))

	for k, v := range params {
		env = append(env, k+"="+v)
	}

	return env
}

type taskArtifactSource struct {
	logger lager.Logger
	volume worker.Volume
}

func newTaskArtifactSource(
	logger lager.Logger,
	volume worker.Volume,
) *taskArtifactSource {
	return &taskArtifactSource{
		logger: logger,
		volume: volume,
	}
}

func (src *taskArtifactSource) StreamTo(destination worker.ArtifactDestination) error {
	out, err := src.volume.StreamOut(".")
	if err != nil {
		return err
	}

	defer out.Close()

	return destination.StreamIn(".", out)
}

func (src *taskArtifactSource) StreamFile(filename string) (io.ReadCloser, error) {
	out, err := src.volume.StreamOut(filename)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(out)

	_, err = tarReader.Next()
	if err != nil {
		return nil, FileNotFoundError{Path: filename}
	}

	return fileReadCloser{
		Reader: tarReader,
		Closer: out,
	}, nil
}

func (src *taskArtifactSource) VolumeOn(w worker.Worker) (worker.Volume, bool, error) {
	return w.LookupVolume(src.logger, src.volume.Handle())
}

type taskInputSource struct {
	config        atc.TaskInputConfig
	source        worker.ArtifactSource
	artifactsRoot string
}

func (s *taskInputSource) Source() worker.ArtifactSource { return s.source }

func (s *taskInputSource) DestinationPath() string {
	subdir := s.config.Path
	if s.config.Path == "" {
		subdir = s.config.Name
	}

	return filepath.Join(s.artifactsRoot, subdir)
}

func artifactsPath(outputConfig atc.TaskOutputConfig, artifactsRoot string) string {
	outputSrc := outputConfig.Path
	if len(outputSrc) == 0 {
		outputSrc = outputConfig.Name
	}

	return path.Join(artifactsRoot, outputSrc) + "/"
}

type taskCacheInputSource struct {
	source        worker.ArtifactSource
	artifactsRoot string
	cachePath     string
}

func (s *taskCacheInputSource) Source() worker.ArtifactSource { return s.source }

func (s *taskCacheInputSource) DestinationPath() string {
	return filepath.Join(s.artifactsRoot, s.cachePath)
}

type taskCacheSource struct {
	logger   lager.Logger
	teamID   int
	jobID    int
	stepName string
	path     string
}

func newTaskCacheSource(
	logger lager.Logger,
	teamID int,
	jobID int,
	stepName string,
	path string,
) *taskCacheSource {
	return &taskCacheSource{
		logger:   logger,
		teamID:   teamID,
		jobID:    jobID,
		stepName: stepName,
		path:     path,
	}
}

func (src *taskCacheSource) StreamTo(destination worker.ArtifactDestination) error {
	// cache will be initialized every time on a new worker
	return nil
}

func (src *taskCacheSource) StreamFile(filename string) (io.ReadCloser, error) {
	return nil, errors.New("taskCacheSource.StreamFile not implemented")
}

func (src *taskCacheSource) VolumeOn(w worker.Worker) (worker.Volume, bool, error) {
	return w.FindVolumeForTaskCache(src.logger, src.teamID, src.jobID, src.stepName, src.path)
}
