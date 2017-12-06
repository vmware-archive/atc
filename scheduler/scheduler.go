package scheduler

import (
	"sort"
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/scheduler/inputmapper"
	"github.com/mitchellh/hashstructure"
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

type Foo []map[string]string

func (s Foo) Len() int {
	return len(s)
}
func (s Foo) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Foo) Less(i, j int) bool {
	v1, _ := hashstructure.Hash(s[i], nil)
	v2, _ := hashstructure.Hash(s[j], nil)
	return v1 < v2
}

func Combinations(resourceSpaces map[string][]string) []map[string]string {
	var resource, space string
	var spaces []string

	combinations := []map[string]string{}

	if len(resourceSpaces) == 0 {
		return combinations
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

	sort.Sort(Foo(combinations))
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
		_, err := job.SyncResourceSpaceCombinations(Combinations(map[string][]string{}))
		jStart := time.Now()
		err = s.ensurePendingBuildExists(logger, versions, job)
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

		err := s.BuildStarter.TryStartPendingBuildsForJob(logger, job, resources, resourceTypes, nextPendingBuildsForJob)
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
) error {
	inputMapping, err := s.InputMapper.SaveNextInputMapping(logger, versions, job)
	if err != nil {
		return err
	}

	for _, inputConfig := range job.Config().Inputs() {
		inputVersion, ok := inputMapping[inputConfig.Name]

		//trigger: true, and the version has not been used
		if ok && inputVersion.FirstOccurrence && inputConfig.Trigger {
			err := job.EnsurePendingBuildExists()
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

	build, err := job.CreateBuild()
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

		err = s.BuildStarter.TryStartPendingBuildsForJob(logger, job, resources, resourceTypes, nextPendingBuilds)
		if err != nil {
			logger.Error("failed-to-start-next-pending-build-for-job", err, lager.Data{"job-name": job.Name()})
			return
		}
	}()

	return build, wg, nil
}

func (s *Scheduler) SaveNextInputMapping(logger lager.Logger, job db.Job) error {
	versions, err := s.Pipeline.LoadVersionsDB()
	if err != nil {
		logger.Error("failed-to-load-versions-db", err)
		return err
	}

	_, err = s.InputMapper.SaveNextInputMapping(logger, versions, job)
	return err
}
