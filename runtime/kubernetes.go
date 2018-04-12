package runtime

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
	"github.com/concourse/baggageclaim"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

const TaskContainerName = "task"

func NewK8sOrchestrator(
	client *kubernetes.Clientset,
	config *restclient.Config,
	namespace string,
	dbVolumeFactory db.VolumeFactory,
	dbTeamFactory db.TeamFactory,
	lockFactory lock.LockFactory,
	clock clock.Clock,
) Orchestrator {
	return &k8sOrchestrator{
		client:    client,
		config:    config,
		namespace: namespace,
		containerProvider: &worker.DbContainerProvider{
			DbVolumeFactory: dbVolumeFactory,
			DbTeamFactory:   dbTeamFactory,
			LockFactory:     lockFactory,
			Clock:           clock,
		},
		volumeFactory: dbVolumeFactory,
	}
}

type k8sOrchestrator struct {
	client            *kubernetes.Clientset
	config            *restclient.Config
	namespace         string
	workerPool        worker.Client
	containerProvider *worker.DbContainerProvider
	volumeFactory     db.VolumeFactory
}

func (s *k8sOrchestrator) SetWorkerPool(pool worker.Client) {
	s.workerPool = pool
}

func (s *k8sOrchestrator) GetResource(
	logger lager.Logger,
	ctx context.Context,
	worker worker.Worker,
	containerSpec worker.ContainerSpec,
	resourceInstance resource.ResourceInstance,
	session resource.Session,
	resourceTypes creds.VersionedResourceTypes,
	delegate worker.ImageFetchingDelegate,
) (resource.VersionedSource, worker.Volume, error) {
	return nil, nil, errors.New("nope")
}

func (o *k8sOrchestrator) k8sWorker(logger lager.Logger) (worker.Worker, error) {
	workers, err := o.workerPool.RunningWorkers(logger)
	if err != nil {
		return nil, err
	}

	for _, w := range workers {
		if w.Type() == "garden" {
			o.containerProvider.WorkerName = w.Name()
			return w, nil
		}
	}

	return nil, errors.New("no worker found for k8s orchestrator")
}

func (o *k8sOrchestrator) RunTask(
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
	result := make(chan TaskResult)

	_, err := o.k8sWorker(logger)
	if err != nil {
		logger.Error("failed-to-find-worker", err)
		return nil, nil, err
	}

	creatingContainer, createdContainer, err := o.containerProvider.Find(containerSpec, owner)

	if err != nil {
		logger.Error("failed-to-find-container-in-db", err)
		return nil, nil, err
	}

	if createdContainer != nil {
		logger = logger.WithData(lager.Data{"container": createdContainer.Handle()})
		logger.Debug("found-created-container-in-db")
		// look up job in k8s
		// get logs
	}

	if creatingContainer != nil {
		// look up job in k8s
		// check status
		// get logs
	} else {
		creatingContainer, err = o.containerProvider.Create(containerSpec, owner, metadata)
		if err != nil {
			logger.Error("failed-to-create-container-in-db", err)
			return nil, nil, err
		}
	}

	task, err := o.createTask(logger, config, containerSpec, ioConfig, creatingContainer)
	if err != nil {
		logger.Error("failed-to-create-task-step", err)
		return nil, nil, err
	}

	//watch pod events relating to job
	err = o.trackProgress(ctx, task, result)
	if err != nil {
		logger.Error("failed-to-track-job-in-k8s", err)
		return nil, nil, err
	}

	return result, []worker.VolumeMount{}, nil
}

type k8sTask struct {
	job             batchv1.Job
	pod             *v1.Pod
	ioConfig        IOConfig
	hasStreamedLogs bool
	outputArgs      []string
}

func (o *k8sOrchestrator) createTask(
	logger lager.Logger,
	config atc.TaskConfig,
	containerSpec worker.ContainerSpec,
	ioConfig IOConfig,
	creatingContainer db.CreatingContainer,
) (*k8sTask, error) {
	chosenWorker, err := o.k8sWorker(logger.Session("find worker"))
	if err != nil {
		return nil, err
	}

	jobSpec, outputArgs, err := o.jobForTask(logger, config, containerSpec, creatingContainer, chosenWorker)
	if err != nil {
		return nil, err
	}
	// Create job
	job, err := o.client.BatchV1().Jobs(o.namespace).Create(jobSpec)
	if err != nil {
		logger.Error("failed-to-create-job", err)
		return nil, err
	}

	return &k8sTask{
		job:        *job,
		ioConfig:   ioConfig,
		outputArgs: outputArgs,
	}, nil

}

func (o *k8sOrchestrator) trackProgress(
	ctx context.Context,
	task *k8sTask,
	result chan TaskResult,
) error {
	logger := lagerctx.FromContext(ctx)

	//watch pod events relating to job
	selector, err := metav1.LabelSelectorAsSelector(task.job.Spec.Selector)
	if err != nil {
		return err
	}

	podEvents, err := o.client.Core().Pods(o.namespace).Watch(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return err
	}
	// filter for modification events, which includes events about the underlying container
	modifyEvent := watch.Filter(podEvents, func(e watch.Event) (watch.Event, bool) {
		return e, e.Type == watch.Modified
	})

	go func() {
		var logStream io.ReadCloser
		var taskResult = TaskResult{}
	L:
		for {
			select {
			case event := <-modifyEvent.ResultChan():
				pod, ok := event.Object.(*v1.Pod)
				if ok {
					task.pod = pod
					logs, result, err := o.handleUpdate(logger, task)
					logStream = logs
					if err != nil {
						taskResult = TaskResult{
							Err: err,
						}
						break L
					}
					if result != nil {
						taskResult = *result
						break L
					}
				}
				continue
			case <-ctx.Done():
				taskResult = TaskResult{
					Err: ctx.Err(),
				}
				break L
			}
		}

		if logStream != nil {
			logStream.Close()
		}
		podEvents.Stop()

		if task.pod != nil {
			outputCmd := append([]string{"conveyor", "outputs"}, task.outputArgs...)
			err := o.execInContainer("prep-outputs", outputCmd, *task.pod, task.ioConfig)
			if err != nil {
				logger.Error("failed-to-prep-outputs", err)
			}
		}

		result <- taskResult
	}()

	return nil
}

func (o *k8sOrchestrator) prepareOutputs(
	logger lager.Logger,
	spec worker.ContainerSpec,
	creatingContainer db.CreatingContainer,
	chosenWorker worker.Worker,
	volumes *[]v1.Volume,
	mounts *[]v1.VolumeMount,
) ([]string, error) {

	v := *volumes
	m := *mounts

	baggageclaimURL := chosenWorker.BaggageclaimURL()
	outputArgs := []string{"-b", *baggageclaimURL}

	for name, outputPath := range spec.Outputs {
		volume := v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		}

		v = append(v, volume)
		mount := v1.VolumeMount{
			Name:      name,
			MountPath: outputPath,
		}
		m = append(m, mount)

		bcVolume, volumeErr := chosenWorker.VolumeClient().FindOrCreateVolumeForContainer(
			logger,
			worker.VolumeSpec{
				Strategy:   baggageclaim.EmptyStrategy{},
				Privileged: false,
			},
			creatingContainer,
			spec.TeamID,
			outputPath,
		)
		if volumeErr != nil {
			return []string{}, volumeErr
		}

		outputArgs = append(outputArgs, "-v", fmt.Sprintf("%s=%s", outputPath, bcVolume.Handle()))
	}

	*volumes = v
	*mounts = m

	return outputArgs, nil

}

func (o *k8sOrchestrator) prepareInputs(
	spec worker.ContainerSpec,
	chosenWorker worker.Worker,
	volumes *[]v1.Volume,
	mounts *[]v1.VolumeMount,
) []string {
	baggageclaimURL := chosenWorker.BaggageclaimURL()
	inputArgs := []string{"-b", *baggageclaimURL}

	v := *volumes
	m := *mounts

	for _, input := range spec.Inputs {
		volume, found, _ := input.Source().VolumeOn(chosenWorker)
		if found {
			k8sVolume := v1.Volume{
				Name: volume.Handle(),
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			}

			v = append(v, k8sVolume)
			mount := v1.VolumeMount{
				Name:      volume.Handle(),
				MountPath: input.DestinationPath(),
			}

			m = append(m, mount)
			inputArgs = append(inputArgs, "-v", fmt.Sprintf("%s=%s", volume.Handle(), input.DestinationPath()))
		}
	}

	*volumes = v
	*mounts = m

	return inputArgs
}

func (o *k8sOrchestrator) jobForTask(
	logger lager.Logger,
	config atc.TaskConfig,
	spec worker.ContainerSpec,
	creatingContainer db.CreatingContainer,
	chosenWorker worker.Worker,
) (*batchv1.Job, []string, error) {
	workDir := path.Join(spec.Dir, config.Run.Dir)
	volumes := []v1.Volume{}
	mounts := []v1.VolumeMount{}

	outputArgs, err := o.prepareOutputs(logger, spec, creatingContainer, chosenWorker, &volumes, &mounts)
	if err != nil {
		return nil, []string{}, err
	}

	inputArgs := o.prepareInputs(spec, chosenWorker, &volumes, &mounts)

	var envVars []v1.EnvVar
	for key, val := range config.Params {
		envVars = append(envVars, v1.EnvVar{Name: key, Value: val})
	}

	var one int32 = 1
	job := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &one,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes:       volumes,
					RestartPolicy: v1.RestartPolicyNever,
					InitContainers: []v1.Container{
						v1.Container{
							Name:            "prep-inputs",
							Image:           "topherbullock/conveyor",
							ImagePullPolicy: v1.PullNever,
							Command:         []string{"conveyor", "inputs"},
							Args:            inputArgs,
							WorkingDir:      workDir,
							VolumeMounts:    mounts,
						},
					},
					Containers: []v1.Container{
						v1.Container{
							Name:         TaskContainerName,
							Image:        config.RootfsURI,
							Command:      []string{config.Run.Path},
							Args:         config.Run.Args,
							WorkingDir:   workDir,
							VolumeMounts: mounts,
							Env:          envVars,
						},
						v1.Container{
							Name:            "prep-outputs",
							Image:           "topherbullock/conveyor",
							ImagePullPolicy: v1.PullNever,
							// Command:      []string{"sleep"},
							// Args:         []string{"10000"},
							VolumeMounts: mounts,
							TTY:          true,
							Stdin:        true,
						},
					},
				},
			},
		},
	}

	job.Name = creatingContainer.Handle()
	return job, outputArgs, nil
}

func (o k8sOrchestrator) streamLogs(podName string, opts *v1.PodLogOptions, stdout io.Writer) (io.ReadCloser, error) {
	logs, err := o.client.Core().Pods(o.namespace).GetLogs(podName, opts).Stream()
	if err != nil {
		errMessage := err.Error()
		if strings.HasPrefix(errMessage, "failed to open log file") {
			// TODO: find a better way to deal with case when there are no logs to stream
			return nil, nil
		}
		return nil, err
	}
	io.Copy(stdout, logs)
	return logs, nil
}

func containerStatus(name string, podStatus v1.PodStatus) *v1.ContainerStatus {
	for _, containerStatus := range podStatus.ContainerStatuses {
		if containerStatus.Name == name {
			return &containerStatus
		}
	}
	return nil
}

func (o *k8sOrchestrator) execInContainer(
	containerName string,
	command []string,
	pod v1.Pod,
	ioConfig IOConfig,
) error {

	req := o.client.RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		Param("container", containerName)

	req.VersionedParams(
		&v1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(o.config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{})
	if err != nil {
		return errors.New("failed to stream")
	}

	return nil
}

func (o *k8sOrchestrator) handleUpdate(logger lager.Logger, task *k8sTask) (io.ReadCloser, *TaskResult, error) {
	var (
		logStream  io.ReadCloser
		taskResult *TaskResult
		err        error
	)

	pod := task.pod
	taskContainerStatus := containerStatus(TaskContainerName, pod.Status)

	if taskContainerStatus != nil {
		state := taskContainerStatus.State
		if taskContainerStatus.Ready {
			opts := &v1.PodLogOptions{
				Follow:    true,
				Container: TaskContainerName,
			}

			if !task.hasStreamedLogs {
				logStream, err = o.streamLogs(pod.Name, opts, task.ioConfig.Stdout)
				if err != nil {
					logger.Error("failed-to-stream-logs-to-stdout", err)
					return nil, nil, err
				}
				task.hasStreamedLogs = true
			}
		}

		if state.Terminated != nil {
			termination := *state.Terminated
			opts := &v1.PodLogOptions{
				Container: TaskContainerName,
			}

			switch termination.Reason {
			case "ContainerCannotRun":
			case "Error":
				if termination.Message != "" {
					task.ioConfig.Stdout.Write([]byte(termination.Message))
				}
			default:
				if !task.hasStreamedLogs {
					logStream, err = o.streamLogs(pod.Name, opts, task.ioConfig.Stdout)
					if err != nil {
						logger.Error("failed-to-stream-logs-to-stdout", err)
						return nil, nil, err
					}
					task.hasStreamedLogs = true
				}
			}

			taskResult = &TaskResult{
				ExitStatus: ExitStatus(termination.ExitCode),
			}
		}
	}

	return logStream, taskResult, nil
}
