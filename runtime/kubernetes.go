package runtime

import (
	"context"
	"errors"

	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/worker"

	"k8s.io/client-go/kubernetes"
)

type K8sOrchestrator struct {
	Client        *kubernetes.Clientset
	Namespace     string
	TeamFactory   db.TeamFactory
	VolumeFactory db.VolumeFactory
	LockFactory   lock.LockFactory
	Worker        db.Worker
}

func (s *K8sOrchestrator) RunTask(
	ctx context.Context,
	delegate TaskExecutionDelegate,
	owner db.ContainerOwner,
	metadata db.ContainerMetadata,
	containerSpec worker.ContainerSpec,
	resourceTypes creds.VersionedResourceTypes,
	ioConfig IOConfig,
	config atc.TaskConfig,
) (chan TaskResult, []worker.VolumeMount, error) {

	// job, err := s.Client.BatchV1().Jobs(s.Namespace).Create(
	// 	&v1.Job{
	// 		Spec: v1.JobSpec{
	// 			Parallelism:  1,
	// 			Completions:  1,
	// 			BackoffLimit: 1,
	// 			Template: corev1.PodTemplateSpec{
	// 				Spec: corev1.PodSpec{
	// 					Containers: []corev1.Container{
	// 						corev1.Container{
	// 							Name:       "task",
	// 							Image:      config.RootfsURI,
	// 							Command:    config.Run.Path,
	// 							Args:       config.Run.Args,
	// 							WorkingDir: config.Run.Dir,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// )
	return nil, nil, errors.New("not implemented")
}
