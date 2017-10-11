package scheduler

import (
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/scheduler/inputmapper"
	"github.com/concourse/atc/space"
)

type Scheduler struct {
	Pipeline     db.Pipeline
	InputMapper  inputmapper.InputMapper
	BuildStarter BuildStarter
	Scanner      Scanner
}

//go:generate counterfeiter . Scanner

type Scanner interface {
	Scan(lager.Logger, string) error
}

func (s *Scheduler) Schedule(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	jobs []db.Job,
	resources db.Resources,
	resourceTypes atc.VersionedResourceTypes,
) (map[string]time.Duration, error) {
	jobSchedulingTime := map[string]time.Duration{}

	jobPermutations := map[db.Job][]db.JobPermutation{}
	for _, job := range jobs {
		permutations, err := job.SyncPermutations(space.FindCombinations(job.Config().Spaces))
		if err != nil {
			logger.Error("failed-to-sync-permutations", err)
			return nil, err
		}

		jobPermutations[job] = permutations
	}

	for job, permutations := range jobPermutations {
		jStart := time.Now()

		for _, permutation := range permutations {
			err := s.ensurePendingBuildExists(logger, versions, jobPermutations, permutation, job.Config().Inputs())
			if err != nil {
				logger.Error("failed-to-ensure-pending-build-exists", err)
				return nil, err
			}
		}

		jobSchedulingTime[job.Name()] = time.Since(jStart)
	}

	nextPendingBuilds, err := s.Pipeline.GetAllPendingBuilds()
	if err != nil {
		logger.Error("failed-to-get-all-next-pending-builds", err)
		return jobSchedulingTime, err
	}

	for job, permutations := range jobPermutations {
		jStart := time.Now()

		for _, permutation := range permutations {
			var nextPendingBuildsForPermutation []db.Build
			for _, build := range nextPendingBuilds {
				if build.JobPermutationID() == permutation.ID() {
					nextPendingBuildsForPermutation = append(nextPendingBuildsForPermutation, build)
				}
			}

			if len(nextPendingBuildsForPermutation) == 0 {
				continue
			}

			err := s.BuildStarter.TryStartPendingBuildsForJobPermutation(logger, job, jobPermutations, permutation, resources, resourceTypes, nextPendingBuildsForPermutation)
			if err != nil {
				return nil, err
			}
		}

		jobSchedulingTime[job.Name()] = jobSchedulingTime[job.Name()] + time.Since(jStart)
	}

	return jobSchedulingTime, nil
}

func (s *Scheduler) ensurePendingBuildExists(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	allJobPermutations map[db.Job][]db.JobPermutation,
	jobPermutation db.JobPermutation,
	inputConfigs []atc.JobInput,
) error {
	inputMapping, err := s.InputMapper.SaveNextInputMapping(logger, versions, allJobPermutations, jobPermutation, inputConfigs)
	if err != nil {
		return err
	}

	for _, inputConfig := range inputConfigs {
		inputVersion, ok := inputMapping[inputConfig.Name]

		//trigger: true, and the version has not been used
		if ok && inputVersion.FirstOccurrence && inputConfig.Trigger {
			err := jobPermutation.EnsurePendingBuildExists()
			if err != nil {
				logger.Error("failed-to-ensure-pending-build-exists", err)
				return err
			}

			break
		}
	}

	return nil
}

type Waiter interface {
	Wait()
}

func (s *Scheduler) TriggerImmediately(
	logger lager.Logger,
	job db.Job,
	resources db.Resources,
	resourceTypes atc.VersionedResourceTypes,
) (db.Build, Waiter, error) {
	logger = logger.Session("trigger-immediately", lager.Data{"job_name": job.Name()})

	return nil, nil, nil
	// build, err := job.CreateBuild()
	// if err != nil {
	// 	logger.Error("failed-to-create-job-build", err)
	// 	return nil, nil, err
	// }
	// wg := new(sync.WaitGroup)
	// wg.Add(1)

	// go func() {
	// 	defer wg.Done()

	// 	nextPendingBuilds, err := job.GetPendingBuilds()
	// 	if err != nil {
	// 		logger.Error("failed-to-get-next-pending-build-for-job", err)
	// 		return
	// 	}

	// 	err = s.BuildStarter.TryStartPendingBuildsForJob(logger, job, resources, resourceTypes, nextPendingBuilds)
	// 	if err != nil {
	// 		logger.Error("failed-to-start-next-pending-build-for-job", err, lager.Data{"job-name": job.Name()})
	// 		return
	// 	}
	// }()

	// return build, wg, nil
}

func (s *Scheduler) SaveNextInputMapping(logger lager.Logger, job db.Job) error {
	return nil
	// versions, err := s.Pipeline.LoadVersionsDB()
	// if err != nil {
	// 	logger.Error("failed-to-load-versions-db", err)
	// 	return err
	// }

	// _, err = s.InputMapper.SaveNextInputMapping(logger, versions, job)
	// return err
}
