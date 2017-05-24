package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
)

func BuildInput(input dbng.BuildInput, config atc.JobInput, source atc.Source) atc.BuildInput {
	return atc.BuildInput{
		Name:     input.Name,
		Resource: input.Resource,
		Type:     input.Type,
		Source:   source,
		Params:   config.Params,
		Version:  atc.Version(input.Version),
		Tags:     config.Tags,
	}
}
