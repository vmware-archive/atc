package scheduler

import (
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/scheduler/inputmapper"
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

func Combinations(resourceSpaces map[string][]string) []map[string]string {
	var resource, space string
	var spaces []string

	combinations := []map[string]string{}

	if len(resourceSpaces) == 0 {
		return []map[string]string{map[string]string{}}
	}

	if len(resourceSpaces) == 1 {
		for resource, spaces = range resourceSpaces {
			for _, space = range spaces {
				combinations = append(combinations, map[string]string{resource: space})
			}
		}
		return combinations
	}

	for resource, spaces = range resourceSpaces {
		break
	}
	delete(resourceSpaces, resource)

	for _, combination := range Combinations(resourceSpaces) {
		for _, space = range spaces {
			copy := map[string]string{}
			for k, v := range combination {
				copy[k] = v
			}

			copy[resource] = space
			combinations = append(combinations, copy)
		}
	}

	return combinations
}

func (s *Scheduler) Schedule(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	jobs []db.Job,
	resources db.Resources,
	resourceTypes atc.VersionedResourceTypes,
) (map[string]time.Duration, error) {
	jobSchedulingTime := map[string]time.Duration{}

	for _, job := range jobs {
		jobCombinations, err := job.SyncResourceSpaceCombinations(Combinations(job.Config().Spaces()))
		if err != nil {
			logger.Error("failed-to-sync-resource-space-combinations", err)
		}

		jStart := time.Now()

		for _, jobCombination := range jobCombinations {
			err = s.ensurePendingBuildExists(logger, versions, job, jobCombination)
			if err != nil {
				break
			}
		}

		jobSchedulingTime[job.Name()] = time.Since(jStart)

		if err != nil {
			return jobSchedulingTime, err
		}
	}

	nextPendingBuilds, err := s.Pipeline.GetAllPendingBuilds()
	if err != nil {
		logger.Error("failed-to-get-all-next-pending-builds", err)
		return jobSchedulingTime, err
	}

	for _, job := range jobs {
		jStart := time.Now()
		nextPendingBuildsForJob, ok := nextPendingBuilds[job.Name()]
		if !ok {
			continue
		}

		err := s.BuildStarter.TryStartPendingBuildsForJob(logger, job, nil, resources, resourceTypes, nextPendingBuildsForJob)
		jobSchedulingTime[job.Name()] = jobSchedulingTime[job.Name()] + time.Since(jStart)

		if err != nil {
			return jobSchedulingTime, err
		}
	}

	return jobSchedulingTime, nil
}

func (s *Scheduler) ensurePendingBuildExists(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	job db.Job,
	jobCombination db.JobCombination,
) error {
	inputMapping, err := s.InputMapper.SaveNextInputMapping(logger, versions, job, jobCombination)
	if err != nil {
		return err
	}

	for _, inputConfig := range job.Config().Inputs() {
		inputVersion, ok := inputMapping[inputConfig.Name]

		//trigger: true, and the version has not been used
		if ok && inputVersion.FirstOccurrence && inputConfig.Trigger {
			err := jobCombination.EnsurePendingBuildExists()
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
	jobCombination db.JobCombination,
	resources db.Resources,
	resourceTypes atc.VersionedResourceTypes,
) (db.Build, Waiter, error) {
	logger = logger.Session("trigger-immediately", lager.Data{"job_name": job.Name()})

	build, err := jobCombination.CreateBuild()
	if err != nil {
		logger.Error("failed-to-create-job-build", err)
		return nil, nil, err
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()

		nextPendingBuilds, err := job.GetPendingBuilds()
		if err != nil {
			logger.Error("failed-to-get-next-pending-build-for-job", err)
			return
		}

		err = s.BuildStarter.TryStartPendingBuildsForJob(logger, job, jobCombination, resources, resourceTypes, nextPendingBuilds)
		if err != nil {
			logger.Error("failed-to-start-next-pending-build-for-job", err, lager.Data{"job-name": job.Name()})
			return
		}
	}()

	return build, wg, nil
}

func (s *Scheduler) SaveNextInputMapping(logger lager.Logger, job db.Job, jobCombination db.JobCombination) error {
	versions, err := s.Pipeline.LoadVersionsDB()
	if err != nil {
		logger.Error("failed-to-load-versions-db", err)
		return err
	}

	_, err = s.InputMapper.SaveNextInputMapping(logger, versions, job, jobCombination)
	return err
}
