package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func SpaceJob(job db.SpaceJob) atc.SpaceJob {
	sanitizedInputs := []atc.JobInput{}
	for _, input := range job.Job.Config().Inputs() {
		sanitizedInputs = append(sanitizedInputs, atc.JobInput{
			Name:     input.Name,
			Resource: input.Resource,
			Passed:   input.Passed,
			Trigger:  input.Trigger,
		})
	}

	sanitizedOutputs := []atc.JobOutput{}
	for _, output := range job.Job.Config().Outputs() {
		sanitizedOutputs = append(sanitizedOutputs, atc.JobOutput{
			Name:     output.Name,
			Resource: output.Resource,
		})
	}

	combinations := []atc.SpaceJobCombination{}
	for _, combination := range job.SpaceJobCombinations {
		combinations = append(combinations, SpaceJobCombination(combination))
	}

	return atc.SpaceJob{
		ID: job.Job.ID(),

		Name:                 job.Job.Name(),
		PipelineName:         job.Job.PipelineName(),
		TeamName:             job.Job.TeamName(),
		DisableManualTrigger: job.Job.Config().DisableManualTrigger,
		Paused:               job.Job.Paused(),
		FirstLoggedBuildID:   job.Job.FirstLoggedBuildID(),

		Inputs:  sanitizedInputs,
		Outputs: sanitizedOutputs,

		Groups: job.Job.Tags(),

		Combinations: combinations,
	}
}
