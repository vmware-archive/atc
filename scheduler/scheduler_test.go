package scheduler_test

import (
	"errors"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/db/dbfakes"
	. "github.com/concourse/atc/scheduler"
	"github.com/concourse/atc/scheduler/inputmapper/inputmapperfakes"
	"github.com/concourse/atc/scheduler/schedulerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheduler", func() {
	var (
		fakePipeline     *dbfakes.FakePipeline
		fakeInputMapper  *inputmapperfakes.FakeInputMapper
		fakeBuildStarter *schedulerfakes.FakeBuildStarter
		fakeScanner      *schedulerfakes.FakeScanner

		scheduler *Scheduler

		disaster error
	)

	BeforeEach(func() {
		fakePipeline = new(dbfakes.FakePipeline)
		fakeInputMapper = new(inputmapperfakes.FakeInputMapper)
		fakeBuildStarter = new(schedulerfakes.FakeBuildStarter)
		fakeScanner = new(schedulerfakes.FakeScanner)

		scheduler = &Scheduler{
			Pipeline:     fakePipeline,
			InputMapper:  fakeInputMapper,
			BuildStarter: fakeBuildStarter,
			Scanner:      fakeScanner,
		}

		disaster = errors.New("bad thing")
	})

	Describe("Schedule", func() {
		var (
			versionsDB             *algorithm.VersionsDB
			fakeJobs               []db.Job
			fakeJob                *dbfakes.FakeJob
			fakeJob2               *dbfakes.FakeJob
			fakeJobCombinations    []db.JobCombination
			fakeJobCombination     *dbfakes.FakeJobCombination
			fakeJobCombination2    *dbfakes.FakeJobCombination
			fakeResource           *dbfakes.FakeResource
			nextPendingBuilds      []db.Build
			nextPendingBuildsJob1  []db.Build
			nextPendingBuildsJob2  []db.Build
			scheduleErr            error
			versionedResourceTypes atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			nextPendingBuilds = []db.Build{new(dbfakes.FakeBuild)}
			nextPendingBuildsJob1 = []db.Build{new(dbfakes.FakeBuild), new(dbfakes.FakeBuild)}
			nextPendingBuildsJob2 = []db.Build{new(dbfakes.FakeBuild)}
			fakePipeline.GetAllPendingBuildsReturns(map[string][]db.Build{
				"some-job":   nextPendingBuilds,
				"some-job-1": nextPendingBuildsJob1,
				"some-job-2": nextPendingBuildsJob2,
			}, nil)

			versionedResourceTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{Name: "some-resource-type"},
					Version:      atc.Version{"some": "version"},
				},
			}

			fakeResource = new(dbfakes.FakeResource)
			fakeResource.NameReturns("some-resource")
		})

		JustBeforeEach(func() {
			versionsDB = &algorithm.VersionsDB{JobCombinationIDs: map[string]int{"j1": 1}}

			var waiter Waiter
			_, scheduleErr = scheduler.Schedule(
				lagertest.NewTestLogger("test"),
				versionsDB,
				fakeJobs,
				db.Resources{fakeResource},
				versionedResourceTypes,
			)
			if waiter != nil {
				waiter.Wait()
			}
		})

		Context("when the job has no inputs", func() {
			BeforeEach(func() {
				fakeJobCombination = new(dbfakes.FakeJobCombination)
				fakeJobCombination2 = new(dbfakes.FakeJobCombination)

				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job-1")
				fakeJob.SyncResourceSpaceCombinationsReturns([]db.JobCombination{fakeJobCombination}, nil)

				fakeJob2 = new(dbfakes.FakeJob)
				fakeJob2.NameReturns("some-job-2")
				fakeJob2.SyncResourceSpaceCombinationsReturns([]db.JobCombination{fakeJobCombination2}, nil)

				fakeJobs = []db.Job{fakeJob, fakeJob2}
			})

			Context("when saving the next input mapping fails", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(scheduleErr).To(Equal(disaster))
				})
			})

			Context("when saving the next input mapping succeeds", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(algorithm.InputMapping{}, nil)
				})

				It("saved the next input mapping for the right job and versions", func() {
					Expect(fakeInputMapper.SaveNextInputMappingCallCount()).To(Equal(2))
					_, actualVersionsDB, actualJob, _ := fakeInputMapper.SaveNextInputMappingArgsForCall(0)
					Expect(actualVersionsDB).To(Equal(versionsDB))
					Expect(actualJob.Name()).To(Equal(fakeJob.Name()))

					_, actualVersionsDB, actualJob, _ = fakeInputMapper.SaveNextInputMappingArgsForCall(1)
					Expect(actualVersionsDB).To(Equal(versionsDB))
					Expect(actualJob.Name()).To(Equal(fakeJob2.Name()))
				})

				Context("when starting pending builds for job fails", func() {
					BeforeEach(func() {
						fakeBuildStarter.TryStartPendingBuildsForJobReturns(disaster)
					})

					It("returns the error", func() {
						Expect(scheduleErr).To(Equal(disaster))
					})

					It("started all pending builds for the right job", func() {
						Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(1))
						_, actualJob, actualJobCombination, actualResources, actualResourceTypes, actualPendingBuilds := fakeBuildStarter.TryStartPendingBuildsForJobArgsForCall(0)
						Expect(actualJob.Name()).To(Equal(fakeJob.Name()))
						Expect(actualJobCombination.ID()).To(Equal(fakeJobCombination.ID()))
						Expect(actualResources).To(Equal(db.Resources{fakeResource}))
						Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
						Expect(actualPendingBuilds).To(Equal(nextPendingBuildsJob1))
					})
				})

				Context("when starting all pending builds succeeds", func() {
					BeforeEach(func() {
						fakeBuildStarter.TryStartPendingBuildsForJobReturns(nil)
					})

					It("returns no error", func() {
						Expect(scheduleErr).NotTo(HaveOccurred())
					})

					It("didn't create a pending build", func() {
						//TODO: create a positive test case for this
						Expect(fakeJobCombination.EnsurePendingBuildExistsCallCount()).To(BeZero())
						Expect(fakeJobCombination2.EnsurePendingBuildExistsCallCount()).To(BeZero())
					})
				})
			})
		})

		Context("when the job has one trigger: true input", func() {
			BeforeEach(func() {
				fakeJob = new(dbfakes.FakeJob)
				fakeJob.NameReturns("some-job")
				fakeJob.ConfigReturns(atc.JobConfig{
					Plan: atc.PlanSequence{
						{Get: "a", Trigger: true},
						{Get: "b", Trigger: false},
					},
				})
				fakeJobs = []db.Job{fakeJob}

				fakeBuildStarter.TryStartPendingBuildsForJobReturns(nil)

				fakeJobCombination = new(dbfakes.FakeJobCombination)
				fakeJobCombination2 = new(dbfakes.FakeJobCombination)
				fakeJobCombinations = []db.JobCombination{fakeJobCombination, fakeJobCombination2}
				fakeJob.SyncResourceSpaceCombinationsReturns(fakeJobCombinations, nil)
			})

			Context("when no input mapping is found", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(algorithm.InputMapping{}, nil)
				})

				It("starts all pending builds and returns no error", func() {
					Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(2))
					Expect(scheduleErr).NotTo(HaveOccurred())

					_, actualJob, actualJobCombination, actualResources, actualResourceTypes, actualPendingBuilds := fakeBuildStarter.TryStartPendingBuildsForJobArgsForCall(0)
					Expect(actualJob.Name()).To(Equal(fakeJob.Name()))
					Expect(actualJobCombination.ID()).To(Equal(fakeJobCombination.ID()))
					Expect(actualResources).To(Equal(db.Resources{fakeResource}))
					Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
					Expect(actualPendingBuilds).To(Equal(nextPendingBuilds))

					_, actualJob, actualJobCombination, actualResources, actualResourceTypes, actualPendingBuilds = fakeBuildStarter.TryStartPendingBuildsForJobArgsForCall(1)
					Expect(actualJob.Name()).To(Equal(fakeJob.Name()))
					Expect(actualJobCombination.ID()).To(Equal(fakeJobCombination2.ID()))
					Expect(actualResources).To(Equal(db.Resources{fakeResource}))
					Expect(actualResourceTypes).To(Equal(versionedResourceTypes))
					Expect(actualPendingBuilds).To(Equal(nextPendingBuilds))
				})

				It("didn't create a pending build", func() {
					Expect(fakeJobCombination.EnsurePendingBuildExistsCallCount()).To(BeZero())
				})
			})

			Context("when no first occurrence input has trigger: true", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(algorithm.InputMapping{
						"a": algorithm.InputVersion{VersionID: 1, FirstOccurrence: false},
						"b": algorithm.InputVersion{VersionID: 2, FirstOccurrence: true},
					}, nil)
				})

				It("starts all pending builds and returns no error", func() {
					Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(2))
					Expect(scheduleErr).NotTo(HaveOccurred())
				})

				It("didn't create a pending build", func() {
					Expect(fakeJobCombination.EnsurePendingBuildExistsCallCount()).To(BeZero())
				})
			})

			Context("when a first occurrence input has trigger: true", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(algorithm.InputMapping{
						"a": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
						"b": algorithm.InputVersion{VersionID: 2, FirstOccurrence: false},
					}, nil)
				})

				Context("when creating a pending build fails", func() {
					BeforeEach(func() {
						fakeJobCombination.EnsurePendingBuildExistsReturns(disaster)
					})

					It("returns the error", func() {
						Expect(scheduleErr).To(Equal(disaster))
					})

					It("created a pending build for the right job", func() {
						Expect(fakeJobCombination.EnsurePendingBuildExistsCallCount()).To(Equal(1))
					})
				})

				Context("when creating a pending build succeeds", func() {
					BeforeEach(func() {
						fakeJobCombination.EnsurePendingBuildExistsReturns(nil)
					})

					It("starts all pending builds and returns no error", func() {
						Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(2))
						Expect(scheduleErr).NotTo(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("TriggerImmediately", func() {
		var (
			fakeJob            *dbfakes.FakeJob
			fakeJobCombination *dbfakes.FakeJobCombination
			fakeResource       *dbfakes.FakeResource
			triggeredBuild     db.Build
			triggerErr         error
			nextPendingBuilds  []db.Build
		)

		BeforeEach(func() {
			fakeJob = new(dbfakes.FakeJob)
			fakeJob.NameReturns("some-job")
			fakeJob.ConfigReturns(atc.JobConfig{Plan: atc.PlanSequence{{Get: "input-1"}, {Get: "input-2"}}})

			fakeJobCombination = new(dbfakes.FakeJobCombination)

			fakeResource = new(dbfakes.FakeResource)
			fakeResource.NameReturns("some-resource")
		})

		JustBeforeEach(func() {
			var waiter Waiter
			triggeredBuild, waiter, triggerErr = scheduler.TriggerImmediately(
				lagertest.NewTestLogger("test"),
				fakeJob,
				fakeJobCombination,
				db.Resources{fakeResource},
				atc.VersionedResourceTypes{
					{
						ResourceType: atc.ResourceType{Name: "some-resource-type"},
						Version:      atc.Version{"some": "version"},
					},
				},
			)
			if waiter != nil {
				waiter.Wait()
			}
		})

		Context("when creating the build fails", func() {
			BeforeEach(func() {
				fakeJobCombination.CreateBuildReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(triggerErr).To(Equal(disaster))
			})
		})

		Context("when creating the build succeeds", func() {
			var createdBuild *dbfakes.FakeBuild

			BeforeEach(func() {
				createdBuild = new(dbfakes.FakeBuild)
				createdBuild.IsManuallyTriggeredReturns(true)
				fakeJobCombination.CreateBuildReturns(createdBuild, nil)
			})

			It("tried to create a build for the right job", func() {
				Expect(fakeJobCombination.CreateBuildCallCount()).To(Equal(1))
			})

			Context("when get pending builds for job fails", func() {
				BeforeEach(func() {
					fakeJob.GetPendingBuildsReturns(nil, disaster)
				})

				It("does not try to start pending builds for job", func() {
					Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(0))
				})
			})

			Context("when get pending builds for job succeeds", func() {
				BeforeEach(func() {
					nextPendingBuilds = []db.Build{new(dbfakes.FakeBuild)}
					fakeJob.GetPendingBuildsReturns(nextPendingBuilds, nil)
				})

				It("tried to get pending builds for the right job", func() {
					Expect(fakeJob.GetPendingBuildsCallCount()).To(Equal(1))
				})

				Context("when trying to start pending builds succeeds", func() {
					BeforeEach(func() {
						fakeBuildStarter.TryStartPendingBuildsForJobReturns(nil)
					})

					It("tries to start builds for the right job", func() {
						Expect(fakeBuildStarter.TryStartPendingBuildsForJobCallCount()).To(Equal(1))
						_, _, _, _, _, b := fakeBuildStarter.TryStartPendingBuildsForJobArgsForCall(0)
						Expect(b).To(Equal(nextPendingBuilds))
					})
				})
			})
		})
	})

	Describe("SaveNextInputMapping", func() {
		var saveErr error
		var fakeJob *dbfakes.FakeJob

		JustBeforeEach(func() {
			fakeJob = new(dbfakes.FakeJob)
			fakeJob.NameReturns("some-job")

			saveErr = scheduler.SaveNextInputMapping(lagertest.NewTestLogger("test"), fakeJob, nil)
		})

		Context("when loading the versions DB fails", func() {
			BeforeEach(func() {
				fakePipeline.LoadVersionsDBReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(saveErr).To(Equal(disaster))
			})
		})

		Context("when loading the versions DB succeeds", func() {
			var versionsDB *algorithm.VersionsDB

			BeforeEach(func() {
				versionsDB = &algorithm.VersionsDB{JobCombinationIDs: map[string]int{"j1": 1}}
				fakePipeline.LoadVersionsDBReturns(versionsDB, nil)
			})

			Context("when saving the next input mapping fails", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(saveErr).To(Equal(disaster))
				})

				It("saved the next input mapping for the right job and versions", func() {
					Expect(fakeInputMapper.SaveNextInputMappingCallCount()).To(Equal(1))
					_, actualVersionsDB, actualJob, _ := fakeInputMapper.SaveNextInputMappingArgsForCall(0)
					Expect(actualVersionsDB).To(Equal(versionsDB))
					Expect(actualJob.Name()).To(Equal(fakeJob.Name()))
				})
			})

			Context("when saving the next input mapping succeeds", func() {
				BeforeEach(func() {
					fakeInputMapper.SaveNextInputMappingReturns(algorithm.InputMapping{
						"some-input": algorithm.InputVersion{VersionID: 1, FirstOccurrence: true},
					}, nil)
				})

				It("returns no error", func() {
					Expect(saveErr).NotTo(HaveOccurred())
				})
			})
		})
	})

	Describe("Combinations", func() {
		var (
			resourceSpaces map[string][]string
			combinations   []map[string]string
		)

		JustBeforeEach(func() {
			combinations = Combinations(resourceSpaces)
		})

		Context("when resource spaces are empty", func() {
			BeforeEach(func() {
				resourceSpaces = map[string][]string{}
			})

			It("returns an empty combination", func() {
				Expect(combinations).To(Equal([]map[string]string{map[string]string{}}))
			})
		})

		Context("when resource spaces contain exactly one resource", func() {
			BeforeEach(func() {
				resourceSpaces = map[string][]string{"some-resource": []string{"some-space", "another-space"}}
			})

			It("returns one combination", func() {
				Expect(len(combinations)).To(Equal(2))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "some-space"}))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "another-space"}))
			})
		})

		Context("when resource spaces contain multiple resources", func() {
			BeforeEach(func() {
				resourceSpaces = map[string][]string{"some-resource": []string{"some-space", "another-space"}, "foo": []string{"bar", "baz"}}
			})

			It("returns all combinations", func() {
				Expect(len(combinations)).To(Equal(4))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "some-space", "foo": "bar"}))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "another-space", "foo": "bar"}))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "some-space", "foo": "baz"}))
				Expect(combinations).To(ContainElement(map[string]string{"some-resource": "another-space", "foo": "baz"}))
			})
		})
	})
})
