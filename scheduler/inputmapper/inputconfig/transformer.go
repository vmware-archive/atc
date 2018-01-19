package inputconfig

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
)

//go:generate counterfeiter . Transformer

type Transformer interface {
	TransformInputConfigs(versionsDB *algorithm.VersionsDB, jobCombination db.JobCombination, inputs []atc.JobInput) (algorithm.InputConfigs, error)
}

func NewTransformer(pipeline db.Pipeline) Transformer {
	return &transformer{pipeline: pipeline}
}

type transformer struct {
	pipeline db.Pipeline
}

func (i *transformer) TransformInputConfigs(versionsDB *algorithm.VersionsDB, jobCombination db.JobCombination, inputs []atc.JobInput) (algorithm.InputConfigs, error) {
	inputConfigs := algorithm.InputConfigs{}

	var jobCombinationID int
	if jobCombination == nil {
		jobCombinationID = 0
	} else {
		jobCombinationID = jobCombination.ID()
	}

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

		jobs := algorithm.JobSet{}
		for _, passedJobName := range input.Passed {
			jobs[versionsDB.JobCombinationIDs[passedJobName]] = struct{}{}
		}

		inputConfigs = append(inputConfigs, algorithm.InputConfig{
			Name:             input.Name,
			UseEveryVersion:  input.Version.Every,
			PinnedVersionID:  pinnedVersionID,
			ResourceSpaceID:  versionsDB.ResourceSpaceIDs[input.Resource],
			Passed:           jobs,
			JobCombinationID: jobCombinationID,
		})
	}

	return inputConfigs, nil
}
