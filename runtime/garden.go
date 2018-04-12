package runtime

import (
	"context"
	"fmt"
	"path"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

const taskProcessID = "task"
const taskProcessPropertyName = "concourse:task-process"
const taskExitStatusPropertyName = "concourse:exit-status"

type GardenOrchestrator struct {
	WorkerPool worker.Client
}

func (s *GardenOrchestrator) SetWorkerPool(pool worker.Client) {
	s.WorkerPool = pool
}

func (s *GardenOrchestrator) RunTask(
	ctx context.Context,
	delegate TaskExecutionDelegate,
	owner db.ContainerOwner,
	metadata db.ContainerMetadata,
	containerSpec worker.ContainerSpec,
	resourceTypes creds.VersionedResourceTypes,
	ioConfig IOConfig,
	config atc.TaskConfig,
) (chan TaskResult, []worker.VolumeMount, error) {

	logger := lagerctx.FromContext(ctx)

	container, err := s.WorkerPool.FindOrCreateContainer(
		ctx,
		logger,
		delegate,
		owner,
		metadata,
		containerSpec,
		resourceTypes,
	)
	if err != nil {
		return nil, nil, err
	}

	var exitStatus int

	exitStatusProp, err := container.Property(taskExitStatusPropertyName)
	if err == nil {
		logger.Info("already-exited", lager.Data{"status": exitStatusProp})

		_, err = fmt.Sscanf(exitStatusProp, "%d", &exitStatus)
		if err != nil {
			return nil, nil, err
		}

		return nil, container.VolumeMounts(), nil
	}

	task := &gardenTask{
		container:     container,
		containerSpec: containerSpec,
		delegate:      delegate,
	}

	return task.run(ctx, ioConfig, config)
}

type gardenTask struct {
	container     worker.Container
	containerSpec worker.ContainerSpec
	delegate      TaskExecutionDelegate
}

type TaskResult struct {
	ExitStatus ExitStatus
	Err        error
}

func (t *gardenTask) run(
	ctx context.Context,
	ioConfig IOConfig,
	config atc.TaskConfig,
) (chan TaskResult, []worker.VolumeMount, error) {
	logger := lagerctx.FromContext(ctx)
	// for backwards compatibility with containers
	// that had their task process name set as property
	var processID string
	processID, err := t.container.Property(taskProcessPropertyName)
	if err != nil {
		processID = taskProcessID
	}

	processIO := garden.ProcessIO{
		Stdout: ioConfig.Stdout,
		Stderr: ioConfig.Stderr,
	}

	process, err := t.container.Attach(processID, processIO)
	if err == nil {
		logger.Info("already-running")
	} else {
		logger.Info("spawning")

		t.delegate.Starting(logger, config)

		process, err = t.container.Run(garden.ProcessSpec{
			ID: taskProcessID,

			Path: config.Run.Path,
			Args: config.Run.Args,

			Dir: path.Join(t.containerSpec.Dir, config.Run.Dir),

			// Guardian sets the default TTY window size to width: 80, height: 24,
			// which creates ANSI control sequences that do not work with other window sizes
			TTY: &garden.TTYSpec{WindowSize: &garden.WindowSize{Columns: 500, Rows: 500}},
		}, processIO)
	}
	if err != nil {
		return nil, t.container.VolumeMounts(), err
	}

	result := make(chan TaskResult)
	exited := make(chan struct{})
	var processStatus int
	var processErr error

	go func() {
		processStatus, processErr = process.Wait()
		close(exited)
	}()

	go func() {
		select {
		case <-ctx.Done():
			err := t.container.Stop(false)
			if err != nil {
				logger.Error("failed-to-stop-container", err, lager.Data{"handle": t.container.Handle()})
			}
		case <-exited:
			if processErr != nil {
				result <- TaskResult{Err: processErr}
			}

			err = t.container.SetProperty(taskExitStatusPropertyName, fmt.Sprintf("%d", processStatus))
			if err != nil {
				result <- TaskResult{Err: err}
			}

			result <- TaskResult{
				ExitStatus: ExitStatus(processStatus),
			}
		}
	}()

	return result, t.container.VolumeMounts(), nil
}

func (t *gardenTask) VolumeMounts() []worker.VolumeMount {
	return t.container.VolumeMounts()
}

func (s *GardenOrchestrator) GetResource(
	logger lager.Logger,
	ctx context.Context,
	targetWorker worker.Worker,
	containerSpec worker.ContainerSpec,
	resourceInstance resource.ResourceInstance,
	session resource.Session,
	resourceTypes creds.VersionedResourceTypes,
	delegate worker.ImageFetchingDelegate,
) (resource.VersionedSource, worker.Volume, error) {

	resourceFactory := resource.NewResourceFactory(targetWorker)
	res, err := resourceFactory.NewResource(
		ctx,
		logger,
		resourceInstance.ContainerOwner(),
		session.Metadata,
		containerSpec,
		resourceTypes,
		delegate,
	)
	if err != nil {
		logger.Error("failed-to-construct-resource", err)
		return nil, nil, err
	}

	var volume worker.Volume
	for _, mount := range res.Container().VolumeMounts() {
		if mount.MountPath == resource.ResourcesDir("get") {
			volume = mount.Volume
			break
		}
	}

	versionedSource, err := res.Get(
		ctx,
		volume,
		resource.IOConfig{
			Stdout: delegate.Stdout(),
			Stderr: delegate.Stderr(),
		},
		resourceInstance.Source(),
		resourceInstance.Params(),
		resourceInstance.Version(),
	)
	if err != nil {
		logger.Error("failed-to-fetch-resource", err)
		return nil, nil, err
	}

	return versionedSource, volume, nil
}
