package maxinflight_test

import (
	"errors"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/dbng/dbngfakes"
	"github.com/concourse/atc/scheduler/maxinflight"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updater", func() {
	var (
		fakePipeline *dbngfakes.FakePipeline
		fakeJob      *dbngfakes.FakeJob
		updater      maxinflight.Updater
		disaster     error
	)

	BeforeEach(func() {
		fakePipeline = new(dbngfakes.FakePipeline)
		fakeJob = new(dbngfakes.FakeJob)
		updater = maxinflight.NewUpdater(fakePipeline)
		disaster = errors.New("bad thing")
	})

	Describe("UpdateMaxInFlightReached", func() {
		var rawMaxInFlight int
		var serialGroups []string
		var updateErr error
		var reached bool

		JustBeforeEach(func() {
			reached, updateErr = updater.UpdateMaxInFlightReached(
				lagertest.NewTestLogger("test"),
				atc.JobConfig{
					Name:           "some-job",
					SerialGroups:   serialGroups,
					RawMaxInFlight: rawMaxInFlight,
				},
				57,
			)
		})

		itReturnsFalseAndNoError := func() {
			It("returns false and no error", func() {
				Expect(updateErr).NotTo(HaveOccurred())
				Expect(reached).To(BeFalse())
				Expect(fakeJob.SetMaxInFlightReachedCallCount()).To(Equal(1))
				actualReached := fakeJob.SetMaxInFlightReachedArgsForCall(0)
				Expect(actualReached).To(BeFalse())
			})
		}

		itReturnsTrueAndNoError := func() {
			It("returns true and no error", func() {
				Expect(updateErr).NotTo(HaveOccurred())
				Expect(reached).To(BeTrue())
				Expect(fakeJob.SetMaxInFlightReachedCallCount()).To(Equal(1))
				actualReached := fakeJob.SetMaxInFlightReachedArgsForCall(0)
				Expect(actualReached).To(BeTrue())
			})
		}

		itReturnsTheError := func() {
			It("returns the error", func() {
				Expect(updateErr).To(Equal(disaster))
				Expect(fakeJob.SetMaxInFlightReachedCallCount()).To(Equal(0))
			})
		}

		Context("when the job is found", func() {
			BeforeEach(func() {
				fakeJob.NameReturns("some-job")
				fakePipeline.JobReturns(fakeJob, true, nil)
			})

			Context("when the job config doesn't specify max in flight", func() {
				BeforeEach(func() {
					rawMaxInFlight = 0
					serialGroups = []string{}
				})

				itReturnsFalseAndNoError()

				It("doesn't look at the database", func() {
					Expect(fakeJob.GetRunningBuildsBySerialGroupCallCount()).To(BeZero())
					Expect(fakeJob.GetNextPendingBuildBySerialGroupCallCount()).To(BeZero())
				})

				Context("when setting max in flight reached fails", func() {
					BeforeEach(func() {
						fakeJob.SetMaxInFlightReachedReturns(disaster)
					})

					It("returns the error", func() {
						Expect(updateErr).To(Equal(disaster))
					})
				})
			})

			itReturnsFalseIfOurBuildIsNext := func() {
				Context("when the build we are trying to run is no longer pending", func() {
					BeforeEach(func() {
						fakeJob.GetNextPendingBuildBySerialGroupReturns(nil, false, nil)
					})

					itReturnsTrueAndNoError()
				})

				Context("when there is another build ahead of us in line", func() {
					var fakeBuild *dbngfakes.FakeBuild

					BeforeEach(func() {
						fakeBuild = new(dbngfakes.FakeBuild)
						fakeBuild.IDReturns(101)
						fakeJob.GetNextPendingBuildBySerialGroupReturns(fakeBuild, true, nil)
					})

					itReturnsTrueAndNoError()
				})

				Context("when the build we are trying to run is first in line", func() {
					var fakeBuild *dbngfakes.FakeBuild

					BeforeEach(func() {
						fakeBuild = new(dbngfakes.FakeBuild)
						fakeBuild.IDReturns(57)
						fakeJob.GetNextPendingBuildBySerialGroupReturns(fakeBuild, true, nil)
					})

					itReturnsFalseAndNoError()
				})
			}

			Context("when the job config specifies max in flight = 3", func() {
				BeforeEach(func() {
					rawMaxInFlight = 3
					serialGroups = []string{}
				})

				Context("when looking up the running builds fails", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns(nil, disaster)
					})

					itReturnsTheError()

					It("looked up the running builds with the right job name and serial group", func() {
						Expect(fakeJob.GetRunningBuildsBySerialGroupCallCount()).To(Equal(1))
						actualSerialGroups := fakeJob.GetRunningBuildsBySerialGroupArgsForCall(0)
						Expect(actualSerialGroups).To(ConsistOf("some-job"))
					})
				})

				Context("when there are 3 builds of the job running", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns([]dbng.Build{nil, nil, nil}, nil)
					})

					itReturnsTrueAndNoError()

					It("doesn't look up the next pending build", func() {
						Expect(fakeJob.GetNextPendingBuildBySerialGroupCallCount()).To(BeZero())
					})
				})

				Context("when there are 2 builds of the job running", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns([]dbng.Build{nil, nil}, nil)
					})

					Context("when looking up the next pending build returns an error", func() {
						BeforeEach(func() {
							fakeJob.GetNextPendingBuildBySerialGroupReturns(nil, false, disaster)
						})

						itReturnsTheError()

						It("looked up the next pending build with the right job name and serial group", func() {
							Expect(fakeJob.GetNextPendingBuildBySerialGroupCallCount()).To(Equal(1))
							actualSerialGroups := fakeJob.GetNextPendingBuildBySerialGroupArgsForCall(0)
							Expect(actualSerialGroups).To(ConsistOf("some-job"))
						})
					})

					itReturnsFalseIfOurBuildIsNext()
				})
			})

			Context("when the job is in serial groups", func() {
				BeforeEach(func() {
					rawMaxInFlight = 0
					serialGroups = []string{"serial-group-1", "serial-group-2"}
				})

				Context("when looking up the running builds fails", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns(nil, disaster)
					})

					itReturnsTheError()

					It("looked up the running builds with the right job name and serial group", func() {
						Expect(fakeJob.GetRunningBuildsBySerialGroupCallCount()).To(Equal(1))
						actualSerialGroups := fakeJob.GetRunningBuildsBySerialGroupArgsForCall(0)
						Expect(actualSerialGroups).To(ConsistOf("serial-group-1", "serial-group-2"))
					})
				})

				Context("when a job in the serial group is running", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns([]dbng.Build{nil}, nil)
					})

					itReturnsTrueAndNoError()

					It("doesn't look up the next pending build", func() {
						Expect(fakeJob.GetNextPendingBuildBySerialGroupCallCount()).To(BeZero())
					})
				})

				Context("when no job in the serial group is running", func() {
					BeforeEach(func() {
						fakeJob.GetRunningBuildsBySerialGroupReturns([]dbng.Build{}, nil)
					})

					Context("when looking up the next pending build returns an error", func() {
						BeforeEach(func() {
							fakeJob.GetNextPendingBuildBySerialGroupReturns(nil, false, disaster)
						})

						itReturnsTheError()

						It("looked up the next pending build with the right job name and serial group", func() {
							Expect(fakeJob.GetNextPendingBuildBySerialGroupCallCount()).To(Equal(1))
							actualSerialGroups := fakeJob.GetNextPendingBuildBySerialGroupArgsForCall(0)
							Expect(actualSerialGroups).To(ConsistOf("serial-group-1", "serial-group-2"))
						})
					})

					itReturnsFalseIfOurBuildIsNext()
				})
			})
		})

		Context("when the job is not found", func() {
			BeforeEach(func() {
				fakePipeline.JobReturns(nil, false, nil)
			})

			It("returns true and no error", func() {
				Expect(updateErr).NotTo(HaveOccurred())
				Expect(reached).To(BeTrue())
			})
		})

		Context("when finding the job fails", func() {
			BeforeEach(func() {
				fakePipeline.JobReturns(nil, false, errors.New("AH"))
			})

			It("returns true and no error", func() {
				Expect(updateErr).To(HaveOccurred())
				Expect(reached).To(BeFalse())
			})
		})
	})
})
