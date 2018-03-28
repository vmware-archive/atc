package runtime

import (
	"context"
	"io"
	"path"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/worker"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func NewK8sOrchestrator(
	client *kubernetes.Clientset,
	namespace string,
	dbVolumeFactory db.VolumeFactory,
	dbTeamFactory db.TeamFactory,
	lockFactory lock.LockFactory,
	clock clock.Clock,
) Orchestrator {
	return &k8sOrchestrator{
		client:    client,
		namespace: namespace,
		containerProvider: &worker.DbContainerProvider{
			DbVolumeFactory: dbVolumeFactory,
			DbTeamFactory:   dbTeamFactory,
			WorkerName:      "kubernetes",
			LockFactory:     lockFactory,
			Clock:           clock,
		},
	}
}

type k8sOrchestrator struct {
	client            *kubernetes.Clientset
	namespace         string
	containerProvider *worker.DbContainerProvider
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

	// Create job
	job, err := o.client.BatchV1().Jobs(o.namespace).Create(o.jobForTask(config, containerSpec, creatingContainer))
	if err != nil {
		logger.Error("failed-to-create-job", err)
		return nil, nil, err
	}

	//watch pod events relating to job
	err = o.trackJobProgress(ctx, job, ioConfig, result)

	if err != nil {
		logger.Error("failed-to-track-job-in-k8s", err)
		return nil, nil, err
	}

	return result, []worker.VolumeMount{}, nil
}

func (o *k8sOrchestrator) trackJobProgress(
	ctx context.Context,
	job *batchv1.Job,
	ioConfig IOConfig,
	result chan TaskResult,
) error {
	logger := lagerctx.FromContext(ctx)
	//watch pod events relating to job
	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return err
	}

	podEvents, err := o.client.Core().Pods(o.namespace).Watch(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return err
	}

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
					logs, result, err := o.handlePodUpdate(logger, pod, ioConfig)
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
		result <- taskResult
	}()

	return nil
}

func (o *k8sOrchestrator) jobForTask(
	config atc.TaskConfig,
	spec worker.ContainerSpec,
	creatingContainer db.CreatingContainer,
) *batchv1.Job {
	var one int32 = 1
	job := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Parallelism:  &one,
			Completions:  &one,
			BackoffLimit: &one,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{
						v1.Container{
							Name:       "task",
							Image:      config.RootfsURI,
							Command:    []string{config.Run.Path},
							Args:       config.Run.Args,
							WorkingDir: path.Join(spec.Dir, config.Run.Dir),
						},
					},
				},
			},
		},
	}
	job.Name = creatingContainer.Handle()
	return job
}

func (o k8sOrchestrator) streamLogs(podName string, opts *v1.PodLogOptions, stdout io.Writer) (io.ReadCloser, error) {
	// TODO: deal with case when there are no logs to stream :
	// failed to open log file "/var/log/pods/3abdb0a3-32bd-11e8-a0f6-080027a5a1ce/task_0.log":
	// open /var/log/pods/3abdb0a3-32bd-11e8-a0f6-080027a5a1ce/task_0.log: no such file or directory
	logs, err := o.client.Core().Pods(o.namespace).GetLogs(podName, opts).Stream()
	if err != nil {
		return nil, err
	}
	io.Copy(stdout, logs)
	return logs, nil
}

func (o *k8sOrchestrator) handlePodUpdate(logger lager.Logger, pod *v1.Pod, ioConfig IOConfig) (io.ReadCloser, *TaskResult, error) {
	var (
		logStream  io.ReadCloser
		taskResult *TaskResult
		err        error
	)

	if len(pod.Status.ContainerStatuses) > 0 {
		state := pod.Status.ContainerStatuses[0].State
		if pod.Status.ContainerStatuses[0].Ready {
			opts := &v1.PodLogOptions{
				Follow: true,
			}

			logStream, err = o.streamLogs(pod.Name, opts, ioConfig.Stdout)
			if err != nil {
				logger.Error("failed-to-stream-logs-to-stdout", err)
				return nil, nil, err
			}
		}

		if state.Terminated != nil {
			termination := *state.Terminated
			opts := &v1.PodLogOptions{}

			logStream, err = o.streamLogs(pod.Name, opts, ioConfig.Stdout)
			if err != nil {
				logger.Error("failed-to-stream-logs-to-stdout", err)
				return nil, nil, err
			}

			taskResult = &TaskResult{
				ExitStatus: ExitStatus(termination.ExitCode),
			}
		}
	}

	return logStream, taskResult, nil
}
