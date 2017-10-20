package resource

import (
	"os"

	"github.com/concourse/atc"
)

type putRequest struct {
	Source atc.Source `json:"source"`
	Params atc.Params `json:"params,omitempty"`
	Space  string     `json:"space,omitempty"`
}

func (resource *resource) Put(
	ioConfig IOConfig,
	source atc.Source,
	params atc.Params,
	space string,
	signals <-chan os.Signal,
	ready chan<- struct{},
) (VersionedSource, error) {
	resourceDir := ResourcesDir("put")

	vs := &putVersionedSource{
		container:   resource.container,
		resourceDir: resourceDir,
	}

	runner := resource.runScript(
		"/opt/resource/out",
		[]string{resourceDir},
		putRequest{
			Params: params,
			Source: source,
			Space:  space,
		},
		&vs.versionResult,
		ioConfig.Stderr,
		true,
	)

	err := runner.Run(signals, ready)
	if err != nil {
		return nil, err
	}

	return vs, nil
}
