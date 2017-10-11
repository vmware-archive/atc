package inputconfig

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
)

//go:generate counterfeiter . Transformer

type Transformer interface {
	TransformInputConfigs(db *algorithm.VersionsDB, allJobPermutations map[db.Job][]db.JobPermutation, jobPermutation db.JobPermutation, inputs []atc.JobInput) (algorithm.InputConfigs, error)
}

func NewTransformer(pipeline db.Pipeline) Transformer {
	return &transformer{pipeline: pipeline}
}

type transformer struct {
	pipeline db.Pipeline
}

func (i *transformer) TransformInputConfigs(versionsDB *algorithm.VersionsDB, allJobPermutations map[db.Job][]db.JobPermutation, jobPermutation db.JobPermutation, inputs []atc.JobInput) (algorithm.InputConfigs, error) {
	inputConfigs := algorithm.InputConfigs{}

	for _, input := range inputs {
		if input.Version == nil {
			input.Version = &atc.VersionConfig{Latest: true}
		}

		pinnedVersionID := 0
		if input.Version.Pinned != nil {
			savedVersion, found, err := i.pipeline.GetVersionedResourceByVersion(input.Version.Pinned, input.Resource)
			if err != nil {
				return nil, err
			}

			if !found {
				continue
			}

			pinnedVersionID = savedVersion.ID
		}

		ourSpaces := jobPermutation.ResourceSpaces()

		jobs := algorithm.JobPermutationSet{}
		for _, passedJobName := range input.Passed {
			var passedJob db.Job
			for job, _ := range allJobPermutations {
				if job.Name() == passedJobName {
					passedJob = job
					break
				}
			}

			if passedJob == nil {
				panic("ruh roh")
			}

			for _, permutation := range allJobPermutations[passedJob] {
				otherSpaces := permutation.ResourceSpaces()

				var mismatch bool
				for resource, space := range ourSpaces {
					otherSpace, found := otherSpaces[resource]
					if found && otherSpace != space {
						mismatch = true
						break
					}
				}

				if !mismatch {
					jobs.Add(permutation.ID())
				}
			}
		}

		inputConfigs = append(inputConfigs, algorithm.InputConfig{
			Name:             input.Name,
			UseEveryVersion:  input.Version.Every,
			PinnedVersionID:  pinnedVersionID,
			ResourceSpaceID:  versionsDB.ResourceSpaceIDs[input.Resource+"{"+ourSpaces[input.Resource]+"}"], // XXX: please no
			Passed:           jobs,
			JobPermutationID: jobPermutation.ID(),
		})
	}

	return inputConfigs, nil
}
