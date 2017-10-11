package inputmapper

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/scheduler/inputmapper/inputconfig"
)

//go:generate counterfeiter . InputMapper

type InputMapper interface {
	SaveNextInputMapping(
		logger lager.Logger,
		versions *algorithm.VersionsDB,
		allJobPermutations map[db.Job][]db.JobPermutation,
		jobPermutation db.JobPermutation,
		inputConfigs []atc.JobInput,
	) (algorithm.InputMapping, error)
}

func NewInputMapper(pipeline db.Pipeline, transformer inputconfig.Transformer) InputMapper {
	return &inputMapper{pipeline: pipeline, transformer: transformer}
}

type inputMapper struct {
	pipeline    db.Pipeline
	transformer inputconfig.Transformer
}

func (i *inputMapper) SaveNextInputMapping(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	allJobPermutations map[db.Job][]db.JobPermutation,
	jobPermutation db.JobPermutation,
	inputConfigs []atc.JobInput,
) (algorithm.InputMapping, error) {
	logger = logger.Session("save-next-input-mapping")

	algorithmInputConfigs, err := i.transformer.TransformInputConfigs(versions, allJobPermutations, jobPermutation, inputConfigs)
	if err != nil {
		logger.Error("failed-to-get-algorithm-input-configs", err)
		return nil, err
	}

	independentMapping := algorithm.InputMapping{}
	for _, inputConfig := range algorithmInputConfigs {
		singletonMapping, ok := algorithm.InputConfigs{inputConfig}.Resolve(versions)
		if ok {
			independentMapping[inputConfig.Name] = singletonMapping[inputConfig.Name]
		}
	}

	err = jobPermutation.SaveIndependentInputMapping(independentMapping)
	if err != nil {
		logger.Error("failed-to-save-independent-input-mapping", err)
		return nil, err
	}

	if len(independentMapping) < len(inputConfigs) {
		// this is necessary to prevent builds from running with missing pinned versions
		err := jobPermutation.DeleteNextInputMapping()
		if err != nil {
			logger.Error("failed-to-delete-next-input-mapping-after-missing-pending", err)
		}

		return nil, err
	}

	resolvedMapping, ok := algorithmInputConfigs.Resolve(versions)
	if !ok {
		err := jobPermutation.DeleteNextInputMapping()
		if err != nil {
			logger.Error("failed-to-delete-next-input-mapping-after-failed-resolve", err)
		}

		return nil, err
	}

	err = jobPermutation.SaveNextInputMapping(resolvedMapping)
	if err != nil {
		logger.Error("failed-to-save-next-input-mapping", err)
		return nil, err
	}

	return resolvedMapping, nil
}
