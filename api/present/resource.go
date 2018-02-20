package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func Resource(resource db.Resource, teamName string) atc.Resource {

	var checkErrString string
	if resource.CheckError() != nil {
		checkErrString = resource.CheckError().Error()
	}

	atcResource := atc.Resource{
		Name:         resource.Name(),
		PipelineName: resource.PipelineName(),
		TeamName:     teamName,
		Type:         resource.Type(),

		Paused: resource.Paused(),

		FailingToCheck: resource.FailingToCheck(),
		CheckError:     checkErrString,
	}

	if !resource.LastChecked().IsZero() {
		atcResource.LastChecked = resource.LastChecked().Unix()
	}

	return atcResource
}
