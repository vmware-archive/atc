// Code generated by counterfeiter. DO NOT EDIT.
package schedulerfakes

import (
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/scheduler"
)

type FakeBuildScheduler struct {
	ScheduleStub        func(logger lager.Logger, versions *algorithm.VersionsDB, jobs []db.Job, resources db.Resources, resourceTypes atc.VersionedResourceTypes) (map[string]time.Duration, error)
	scheduleMutex       sync.RWMutex
	scheduleArgsForCall []struct {
		logger        lager.Logger
		versions      *algorithm.VersionsDB
		jobs          []db.Job
		resources     db.Resources
		resourceTypes atc.VersionedResourceTypes
	}
	scheduleReturns struct {
		result1 map[string]time.Duration
		result2 error
	}
	scheduleReturnsOnCall map[int]struct {
		result1 map[string]time.Duration
		result2 error
	}
	TriggerImmediatelyStub        func(logger lager.Logger, job db.Job, resources db.Resources, resourceTypes atc.VersionedResourceTypes) (db.Build, scheduler.Waiter, error)
	triggerImmediatelyMutex       sync.RWMutex
	triggerImmediatelyArgsForCall []struct {
		logger        lager.Logger
		job           db.Job
		resources     db.Resources
		resourceTypes atc.VersionedResourceTypes
	}
	triggerImmediatelyReturns struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}
	triggerImmediatelyReturnsOnCall map[int]struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}
	SaveNextInputMappingStub        func(logger lager.Logger, job db.Job, jobCombination db.JobCombination) error
	saveNextInputMappingMutex       sync.RWMutex
	saveNextInputMappingArgsForCall []struct {
		logger         lager.Logger
		job            db.Job
		jobCombination db.JobCombination
	}
	saveNextInputMappingReturns struct {
		result1 error
	}
	saveNextInputMappingReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildScheduler) Schedule(logger lager.Logger, versions *algorithm.VersionsDB, jobs []db.Job, resources db.Resources, resourceTypes atc.VersionedResourceTypes) (map[string]time.Duration, error) {
	var jobsCopy []db.Job
	if jobs != nil {
		jobsCopy = make([]db.Job, len(jobs))
		copy(jobsCopy, jobs)
	}
	fake.scheduleMutex.Lock()
	ret, specificReturn := fake.scheduleReturnsOnCall[len(fake.scheduleArgsForCall)]
	fake.scheduleArgsForCall = append(fake.scheduleArgsForCall, struct {
		logger        lager.Logger
		versions      *algorithm.VersionsDB
		jobs          []db.Job
		resources     db.Resources
		resourceTypes atc.VersionedResourceTypes
	}{logger, versions, jobsCopy, resources, resourceTypes})
	fake.recordInvocation("Schedule", []interface{}{logger, versions, jobsCopy, resources, resourceTypes})
	fake.scheduleMutex.Unlock()
	if fake.ScheduleStub != nil {
		return fake.ScheduleStub(logger, versions, jobs, resources, resourceTypes)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.scheduleReturns.result1, fake.scheduleReturns.result2
}

func (fake *FakeBuildScheduler) ScheduleCallCount() int {
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	return len(fake.scheduleArgsForCall)
}

func (fake *FakeBuildScheduler) ScheduleArgsForCall(i int) (lager.Logger, *algorithm.VersionsDB, []db.Job, db.Resources, atc.VersionedResourceTypes) {
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	return fake.scheduleArgsForCall[i].logger, fake.scheduleArgsForCall[i].versions, fake.scheduleArgsForCall[i].jobs, fake.scheduleArgsForCall[i].resources, fake.scheduleArgsForCall[i].resourceTypes
}

func (fake *FakeBuildScheduler) ScheduleReturns(result1 map[string]time.Duration, result2 error) {
	fake.ScheduleStub = nil
	fake.scheduleReturns = struct {
		result1 map[string]time.Duration
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildScheduler) ScheduleReturnsOnCall(i int, result1 map[string]time.Duration, result2 error) {
	fake.ScheduleStub = nil
	if fake.scheduleReturnsOnCall == nil {
		fake.scheduleReturnsOnCall = make(map[int]struct {
			result1 map[string]time.Duration
			result2 error
		})
	}
	fake.scheduleReturnsOnCall[i] = struct {
		result1 map[string]time.Duration
		result2 error
	}{result1, result2}
}

func (fake *FakeBuildScheduler) TriggerImmediately(logger lager.Logger, job db.Job, resources db.Resources, resourceTypes atc.VersionedResourceTypes) (db.Build, scheduler.Waiter, error) {
	fake.triggerImmediatelyMutex.Lock()
	ret, specificReturn := fake.triggerImmediatelyReturnsOnCall[len(fake.triggerImmediatelyArgsForCall)]
	fake.triggerImmediatelyArgsForCall = append(fake.triggerImmediatelyArgsForCall, struct {
		logger        lager.Logger
		job           db.Job
		resources     db.Resources
		resourceTypes atc.VersionedResourceTypes
	}{logger, job, resources, resourceTypes})
	fake.recordInvocation("TriggerImmediately", []interface{}{logger, job, resources, resourceTypes})
	fake.triggerImmediatelyMutex.Unlock()
	if fake.TriggerImmediatelyStub != nil {
		return fake.TriggerImmediatelyStub(logger, job, resources, resourceTypes)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.triggerImmediatelyReturns.result1, fake.triggerImmediatelyReturns.result2, fake.triggerImmediatelyReturns.result3
}

func (fake *FakeBuildScheduler) TriggerImmediatelyCallCount() int {
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	return len(fake.triggerImmediatelyArgsForCall)
}

func (fake *FakeBuildScheduler) TriggerImmediatelyArgsForCall(i int) (lager.Logger, db.Job, db.Resources, atc.VersionedResourceTypes) {
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	return fake.triggerImmediatelyArgsForCall[i].logger, fake.triggerImmediatelyArgsForCall[i].job, fake.triggerImmediatelyArgsForCall[i].resources, fake.triggerImmediatelyArgsForCall[i].resourceTypes
}

func (fake *FakeBuildScheduler) TriggerImmediatelyReturns(result1 db.Build, result2 scheduler.Waiter, result3 error) {
	fake.TriggerImmediatelyStub = nil
	fake.triggerImmediatelyReturns = struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildScheduler) TriggerImmediatelyReturnsOnCall(i int, result1 db.Build, result2 scheduler.Waiter, result3 error) {
	fake.TriggerImmediatelyStub = nil
	if fake.triggerImmediatelyReturnsOnCall == nil {
		fake.triggerImmediatelyReturnsOnCall = make(map[int]struct {
			result1 db.Build
			result2 scheduler.Waiter
			result3 error
		})
	}
	fake.triggerImmediatelyReturnsOnCall[i] = struct {
		result1 db.Build
		result2 scheduler.Waiter
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeBuildScheduler) SaveNextInputMapping(logger lager.Logger, job db.Job, jobCombination db.JobCombination) error {
	fake.saveNextInputMappingMutex.Lock()
	ret, specificReturn := fake.saveNextInputMappingReturnsOnCall[len(fake.saveNextInputMappingArgsForCall)]
	fake.saveNextInputMappingArgsForCall = append(fake.saveNextInputMappingArgsForCall, struct {
		logger         lager.Logger
		job            db.Job
		jobCombination db.JobCombination
	}{logger, job, jobCombination})
	fake.recordInvocation("SaveNextInputMapping", []interface{}{logger, job, jobCombination})
	fake.saveNextInputMappingMutex.Unlock()
	if fake.SaveNextInputMappingStub != nil {
		return fake.SaveNextInputMappingStub(logger, job, jobCombination)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.saveNextInputMappingReturns.result1
}

func (fake *FakeBuildScheduler) SaveNextInputMappingCallCount() int {
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	return len(fake.saveNextInputMappingArgsForCall)
}

func (fake *FakeBuildScheduler) SaveNextInputMappingArgsForCall(i int) (lager.Logger, db.Job, db.JobCombination) {
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	return fake.saveNextInputMappingArgsForCall[i].logger, fake.saveNextInputMappingArgsForCall[i].job, fake.saveNextInputMappingArgsForCall[i].jobCombination
}

func (fake *FakeBuildScheduler) SaveNextInputMappingReturns(result1 error) {
	fake.SaveNextInputMappingStub = nil
	fake.saveNextInputMappingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildScheduler) SaveNextInputMappingReturnsOnCall(i int, result1 error) {
	fake.SaveNextInputMappingStub = nil
	if fake.saveNextInputMappingReturnsOnCall == nil {
		fake.saveNextInputMappingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveNextInputMappingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildScheduler) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.scheduleMutex.RLock()
	defer fake.scheduleMutex.RUnlock()
	fake.triggerImmediatelyMutex.RLock()
	defer fake.triggerImmediatelyMutex.RUnlock()
	fake.saveNextInputMappingMutex.RLock()
	defer fake.saveNextInputMappingMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildScheduler) recordInvocation(key string, args []interface{}) {
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

var _ scheduler.BuildScheduler = new(FakeBuildScheduler)
