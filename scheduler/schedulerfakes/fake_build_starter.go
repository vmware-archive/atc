// Code generated by counterfeiter. DO NOT EDIT.
package schedulerfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/scheduler"
)

type FakeBuildStarter struct {
	TryStartPendingBuildsForJobCombinationStub        func(logger lager.Logger, job db.Job, jobCombination db.JobCombination, resources db.Resources, resourceTypes atc.VersionedResourceTypes, nextPendingBuilds []db.Build) error
	tryStartPendingBuildsForJobCombinationMutex       sync.RWMutex
	tryStartPendingBuildsForJobCombinationArgsForCall []struct {
		logger            lager.Logger
		job               db.Job
		jobCombination    db.JobCombination
		resources         db.Resources
		resourceTypes     atc.VersionedResourceTypes
		nextPendingBuilds []db.Build
	}
	tryStartPendingBuildsForJobCombinationReturns struct {
		result1 error
	}
	tryStartPendingBuildsForJobCombinationReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCombination(logger lager.Logger, job db.Job, jobCombination db.JobCombination, resources db.Resources, resourceTypes atc.VersionedResourceTypes, nextPendingBuilds []db.Build) error {
	var nextPendingBuildsCopy []db.Build
	if nextPendingBuilds != nil {
		nextPendingBuildsCopy = make([]db.Build, len(nextPendingBuilds))
		copy(nextPendingBuildsCopy, nextPendingBuilds)
	}
	fake.tryStartPendingBuildsForJobCombinationMutex.Lock()
	ret, specificReturn := fake.tryStartPendingBuildsForJobCombinationReturnsOnCall[len(fake.tryStartPendingBuildsForJobCombinationArgsForCall)]
	fake.tryStartPendingBuildsForJobCombinationArgsForCall = append(fake.tryStartPendingBuildsForJobCombinationArgsForCall, struct {
		logger            lager.Logger
		job               db.Job
		jobCombination    db.JobCombination
		resources         db.Resources
		resourceTypes     atc.VersionedResourceTypes
		nextPendingBuilds []db.Build
	}{logger, job, jobCombination, resources, resourceTypes, nextPendingBuildsCopy})
	fake.recordInvocation("TryStartPendingBuildsForJobCombination", []interface{}{logger, job, jobCombination, resources, resourceTypes, nextPendingBuildsCopy})
	fake.tryStartPendingBuildsForJobCombinationMutex.Unlock()
	if fake.TryStartPendingBuildsForJobCombinationStub != nil {
		return fake.TryStartPendingBuildsForJobCombinationStub(logger, job, jobCombination, resources, resourceTypes, nextPendingBuilds)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.tryStartPendingBuildsForJobCombinationReturns.result1
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCombinationCallCount() int {
	fake.tryStartPendingBuildsForJobCombinationMutex.RLock()
	defer fake.tryStartPendingBuildsForJobCombinationMutex.RUnlock()
	return len(fake.tryStartPendingBuildsForJobCombinationArgsForCall)
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCombinationArgsForCall(i int) (lager.Logger, db.Job, db.JobCombination, db.Resources, atc.VersionedResourceTypes, []db.Build) {
	fake.tryStartPendingBuildsForJobCombinationMutex.RLock()
	defer fake.tryStartPendingBuildsForJobCombinationMutex.RUnlock()
	return fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].logger, fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].job, fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].jobCombination, fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].resources, fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].resourceTypes, fake.tryStartPendingBuildsForJobCombinationArgsForCall[i].nextPendingBuilds
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCombinationReturns(result1 error) {
	fake.TryStartPendingBuildsForJobCombinationStub = nil
	fake.tryStartPendingBuildsForJobCombinationReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildStarter) TryStartPendingBuildsForJobCombinationReturnsOnCall(i int, result1 error) {
	fake.TryStartPendingBuildsForJobCombinationStub = nil
	if fake.tryStartPendingBuildsForJobCombinationReturnsOnCall == nil {
		fake.tryStartPendingBuildsForJobCombinationReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.tryStartPendingBuildsForJobCombinationReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildStarter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.tryStartPendingBuildsForJobCombinationMutex.RLock()
	defer fake.tryStartPendingBuildsForJobCombinationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildStarter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ scheduler.BuildStarter = new(FakeBuildStarter)
