package runtime

import (
	"context"

	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

type NoOpScheduler struct{}

func (*NoOpScheduler) Task(context.Context, TaskExecutionDelegate, db.ContainerOwner, db.ContainerMetadata, worker.ContainerSpec, creds.VersionedResourceTypes) Task {
	return &noopTask{}
}

func (*NoOpScheduler) Resource(atc.Source, atc.Params) Resource {
	return &noopResource{}
}

type noopResource struct{}

func (*noopResource) Get(context.Context, worker.Volume, IOConfig, atc.Source, atc.Params, atc.Version) (resource.VersionedSource, error) {
	return nil, nil
}
func (*noopResource) Put(context.Context, IOConfig, atc.Source, atc.Params) (resource.VersionedSource, error) {
	return nil, nil
}
func (*noopResource) Check(atc.Source, atc.Version) ([]atc.Version, error) {
	return []atc.Version{}, nil
}

type noopTask struct{}

func (*noopTask) Run(context.Context, IOConfig, atc.TaskConfig) (chan TaskResult, []worker.VolumeMount, error) {
	return make(chan TaskResult), []worker.VolumeMount{}, nil
}
