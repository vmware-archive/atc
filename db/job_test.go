package db_test

import (
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	var (
		jobCombination db.JobCombination
		job            db.Job
		pipeline       db.Pipeline
		team           db.Team
	)

	BeforeEach(func() {
		var err error
		team, err = teamFactory.CreateTeam(atc.Team{Name: "some-team"})
		Expect(err).ToNot(HaveOccurred())

		var created bool
		pipeline, created, err = team.SavePipeline("fake-pipeline", atc.Config{
			Jobs: atc.JobConfigs{
				{
					Name: "some-job",

					Public: true,

					Serial: true,

					SerialGroups: []string{"serial-group"},

					Plan: atc.PlanSequence{
						{
							Put: "some-resource",
							Params: atc.Params{
								"some-param": "some-value",
							},
						},
						{
							Get:      "some-input",
							Resource: "some-resource",
							Params: atc.Params{
								"some-param": "some-value",
							},
							Passed:  []string{"job-1", "job-2"},
							Trigger: true,
						},
						{
							Task:           "some-task",
							Privileged:     true,
							TaskConfigPath: "some/config/path.yml",
							TaskConfig: &atc.TaskConfig{
								RootfsURI: "some-image",
							},
						},
					},
				},
				{
					Name: "some-other-job",
				},
				{
					Name:         "other-serial-group-job",
					SerialGroups: []string{"serial-group", "really-different-group"},
				},
				{
					Name:         "different-serial-group-job",
					SerialGroups: []string{"different-serial-group"},
				},
			},
			Resources: atc.ResourceConfigs{
				{
					Name: "some-resource",
					Type: "some-type",
				},
				{
					Name: "some-other-resource",
					Type: "some-type",
				},
			},
		}, db.ConfigVersion(0), db.PipelineUnpaused)
		Expect(err).ToNot(HaveOccurred())
		Expect(created).To(BeTrue())

		var found bool
		job, found, err = pipeline.Job("some-job")
		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(BeTrue())

		jobCombination, err = job.JobCombination()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Pause and Unpause", func() {
		It("starts out as unpaused", func() {
			Expect(job.Paused()).To(BeFalse())
		})

		It("can be paused", func() {
			err := job.Pause()
			Expect(err).NotTo(HaveOccurred())

			found, err := job.Reload()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(job.Paused()).To(BeTrue())
		})

		It("can be unpaused", func() {
			err := job.Unpause()
			Expect(err).NotTo(HaveOccurred())

			found, err := job.Reload()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(job.Paused()).To(BeFalse())
		})
	})

	Describe("FinishedAndNextBuild", func() {
		var otherPipeline db.Pipeline
		var otherJob db.Job
		var otherJobCombination db.JobCombination

		BeforeEach(func() {
			var created bool
			var err error
			otherPipeline, created, err = team.SavePipeline("other-pipeline", atc.Config{
				Jobs: atc.JobConfigs{
					{Name: "some-job"},
				},
			}, db.ConfigVersion(0), db.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeTrue())

			var found bool
			otherJob, found, err = otherPipeline.Job("some-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())

			otherJobCombination, err = otherJob.JobCombination()
			Expect(err).ToNot(HaveOccurred())
		})

		It("can report a job's latest running and finished builds", func() {
			finished, next, err := job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next).To(BeNil())
			Expect(finished).To(BeNil())

			finishedBuild, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			err = finishedBuild.Finish(db.BuildStatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			otherFinishedBuild, err := otherJobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			err = otherFinishedBuild.Finish(db.BuildStatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next).To(BeNil())
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			nextBuild, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			started, err := nextBuild.Start("some-engine", `{"id":"1"}`, atc.Plan{})
			Expect(err).NotTo(HaveOccurred())
			Expect(started).To(BeTrue())

			otherNextBuild, err := otherJobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			otherStarted, err := otherNextBuild.Start("some-engine", `{"id":"1"}`, atc.Plan{})
			Expect(err).NotTo(HaveOccurred())
			Expect(otherStarted).To(BeTrue())

			finished, next, err = job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID()))
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			anotherRunningBuild, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID())) // not anotherRunningBuild
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			started, err = anotherRunningBuild.Start("some-engine", `{"meta":"data"}`, atc.Plan{})
			Expect(err).NotTo(HaveOccurred())
			Expect(started).To(BeTrue())

			finished, next, err = job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID())) // not anotherRunningBuild
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			err = nextBuild.Finish(db.BuildStatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = job.FinishedAndNextBuild()
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(anotherRunningBuild.ID()))
			Expect(finished.ID()).To(Equal(nextBuild.ID()))
		})
	})

	Describe("UpdateFirstLoggedBuildID", func() {
		It("updates FirstLoggedBuildID on a job", func() {
			By("starting out as 0")
			Expect(job.FirstLoggedBuildID()).To(BeZero())

			By("increasing it to 57")
			err := job.UpdateFirstLoggedBuildID(57)
			Expect(err).NotTo(HaveOccurred())

			found, err := job.Reload()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(job.FirstLoggedBuildID()).To(Equal(57))

			By("not erroring when it's called with the same number")
			err = job.UpdateFirstLoggedBuildID(57)
			Expect(err).NotTo(HaveOccurred())

			By("erroring when the number decreases")
			err = job.UpdateFirstLoggedBuildID(56)
			Expect(err).To(Equal(db.FirstLoggedBuildIDDecreasedError{
				Job:   "some-job",
				OldID: 57,
				NewID: 56,
			}))
		})
	})

	Describe("JobCombination", func() {
		It("returns the job combination given its id", func() {
			combination, err := job.JobCombination()
			Expect(err).NotTo(HaveOccurred())
			Expect(combination.ID()).To(Equal(jobCombination.ID()))
			Expect(combination.JobID()).To(Equal(jobCombination.JobID()))
			Expect(combination.Combination()).To(Equal(jobCombination.Combination()))
		})
	})

	Describe("JobCombinations", func() {
		It("returns all the job combinations", func() {
			combinations, err := job.JobCombinations()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(combinations)).To(Equal(1))
			Expect(combinations[0].ID()).To(Equal(jobCombination.ID()))
			Expect(combinations[0].JobID()).To(Equal(jobCombination.JobID()))
			Expect(combinations[0].Combination()).To(Equal(jobCombination.Combination()))
		})
	})

	Describe("Builds", func() {
		var (
			builds       [10]db.Build
			someJob      db.Job
			someOtherJob db.Job
		)

		BeforeEach(func() {
			var found bool
			var err error
			someJob, found, err = pipeline.Job("some-job")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			someOtherJob, found, err = pipeline.Job("some-other-job")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			someJobCombination := jobCombination

			someOtherJobCombination, err := someOtherJob.JobCombination()
			Expect(err).ToNot(HaveOccurred())

			for i := 0; i < 10; i++ {
				build, err := someJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				_, err = someOtherJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				builds[i] = build
			}
		})

		Context("when there are no builds to be found", func() {
			It("returns the builds, with previous/next pages", func() {
				buildsPage, pagination, err := someOtherJob.Builds(db.Page{})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{}))
				Expect(pagination).To(Equal(db.Pagination{}))
			})
		})

		Context("with no since/until", func() {
			It("returns the first page, with the given limit, and a next page", func() {
				buildsPage, pagination, err := someJob.Builds(db.Page{Limit: 2})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{builds[9], builds[8]}))
				Expect(pagination.Previous).To(BeNil())
				Expect(pagination.Next).To(Equal(&db.Page{Since: builds[8].ID(), Limit: 2}))
			})
		})

		Context("with a since that places it in the middle of the builds", func() {
			It("returns the builds, with previous/next pages", func() {
				buildsPage, pagination, err := someJob.Builds(db.Page{Since: builds[6].ID(), Limit: 2})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{builds[5], builds[4]}))
				Expect(pagination.Previous).To(Equal(&db.Page{Until: builds[5].ID(), Limit: 2}))
				Expect(pagination.Next).To(Equal(&db.Page{Since: builds[4].ID(), Limit: 2}))
			})
		})

		Context("with a since that places it at the end of the builds", func() {
			It("returns the builds, with previous/next pages", func() {
				buildsPage, pagination, err := someJob.Builds(db.Page{Since: builds[2].ID(), Limit: 2})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{builds[1], builds[0]}))
				Expect(pagination.Previous).To(Equal(&db.Page{Until: builds[1].ID(), Limit: 2}))
				Expect(pagination.Next).To(BeNil())
			})
		})

		Context("with an until that places it in the middle of the builds", func() {
			It("returns the builds, with previous/next pages", func() {
				buildsPage, pagination, err := someJob.Builds(db.Page{Until: builds[6].ID(), Limit: 2})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{builds[8], builds[7]}))
				Expect(pagination.Previous).To(Equal(&db.Page{Until: builds[8].ID(), Limit: 2}))
				Expect(pagination.Next).To(Equal(&db.Page{Since: builds[7].ID(), Limit: 2}))
			})
		})

		Context("with a until that places it at the beginning of the builds", func() {
			It("returns the builds, with previous/next pages", func() {
				buildsPage, pagination, err := someJob.Builds(db.Page{Until: builds[7].ID(), Limit: 2})
				Expect(err).ToNot(HaveOccurred())
				Expect(buildsPage).To(Equal([]db.Build{builds[9], builds[8]}))
				Expect(pagination.Previous).To(BeNil())
				Expect(pagination.Next).To(Equal(&db.Page{Since: builds[8].ID(), Limit: 2}))
			})
		})
	})

	Describe("Build", func() {
		var firstBuild db.Build

		Context("when a build exists", func() {
			BeforeEach(func() {
				var err error
				firstBuild, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())
			})

			It("finds the latest build", func() {
				secondBuild, err := jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				build, found, err := job.Build("latest")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(secondBuild.ID()))
				Expect(build.Status()).To(Equal(secondBuild.Status()))
			})

			It("finds the build", func() {
				build, found, err := job.Build(firstBuild.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(firstBuild.ID()))
				Expect(build.Status()).To(Equal(firstBuild.Status()))
			})
		})

		Context("when the build does not exist", func() {
			It("does not error", func() {
				build, found, err := job.Build("bogus-build")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(build).To(BeNil())
			})

			It("does not error finding the latest", func() {
				build, found, err := job.Build("latest")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(build).To(BeNil())
			})
		})
	})

	Describe("GetRunningBuildsBySerialGroup", func() {
		Describe("same job", func() {
			var startedBuild, scheduledBuild db.Build

			BeforeEach(func() {
				var err error
				_, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				startedBuild, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())
				_, err = startedBuild.Start("", "{}", atc.Plan{})
				Expect(err).NotTo(HaveOccurred())

				scheduledBuild, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				scheduled, err := scheduledBuild.Schedule()
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduled).To(BeTrue())

				for _, s := range []db.BuildStatus{db.BuildStatusSucceeded, db.BuildStatusFailed, db.BuildStatusErrored, db.BuildStatusAborted} {
					finishedBuild, err := jobCombination.CreateBuild()
					Expect(err).NotTo(HaveOccurred())

					scheduled, err = finishedBuild.Schedule()
					Expect(err).NotTo(HaveOccurred())
					Expect(scheduled).To(BeTrue())

					err = finishedBuild.Finish(s)
					Expect(err).NotTo(HaveOccurred())
				}

				otherJob, found, err := pipeline.Job("some-other-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				otherJobCombination, err := otherJob.JobCombination()
				Expect(err).ToNot(HaveOccurred())

				_, err = otherJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a list of running or schedule builds for said job", func() {
				builds, err := job.GetRunningBuildsBySerialGroup([]string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(2))
				ids := []int{}
				for _, build := range builds {
					ids = append(ids, build.ID())
				}
				Expect(ids).To(ConsistOf([]int{startedBuild.ID(), scheduledBuild.ID()}))
			})
		})

		Describe("multiple jobs with same serial group", func() {
			var serialGroupBuild db.Build

			BeforeEach(func() {
				var err error
				_, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				otherSerialJob, found, err := pipeline.Job("other-serial-group-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				otherSerialJobCombination, err := otherSerialJob.JobCombination()
				Expect(err).ToNot(HaveOccurred())

				serialGroupBuild, err = otherSerialJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				scheduled, err := serialGroupBuild.Schedule()
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduled).To(BeTrue())

				differentSerialJob, found, err := pipeline.Job("different-serial-group-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				differentSerialJobCombination, err := differentSerialJob.JobCombination()
				Expect(err).ToNot(HaveOccurred())

				differentSerialGroupBuild, err := differentSerialJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				scheduled, err = differentSerialGroupBuild.Schedule()
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduled).To(BeTrue())
			})

			It("returns a list of builds in the same serial group", func() {
				builds, err := job.GetRunningBuildsBySerialGroup([]string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(1))
				Expect(builds[0].ID()).To(Equal(serialGroupBuild.ID()))
			})
		})
	})

	Describe("GetNextPendingBuildBySerialGroup", func() {
		var job1, job2 db.Job
		var job2Combination db.JobCombination

		BeforeEach(func() {
			var found bool
			var err error

			job1 = job

			job2, found, err = pipeline.Job("other-serial-group-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())

			job2Combination, err = job2.JobCombination()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when some jobs have builds with inputs determined as false", func() {
			var actualBuild db.Build

			BeforeEach(func() {
				_, err := jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				actualBuild, err = job2Combination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				err = job2Combination.SaveNextInputMapping(nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the next most pending build in a group of jobs", func() {
				build, found, err := job1.GetNextPendingBuildBySerialGroup([]string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(actualBuild.ID()))
			})
		})

		It("should return the next most pending build in a group of jobs", func() {
			buildOne, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			buildTwo, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			buildThree, err := job2Combination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			err = jobCombination.SaveNextInputMapping(nil)
			Expect(err).NotTo(HaveOccurred())
			err = job2Combination.SaveNextInputMapping(nil)
			Expect(err).NotTo(HaveOccurred())

			build, found, err := job1.GetNextPendingBuildBySerialGroup([]string{"serial-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildOne.ID()))

			err = job1.Pause()
			Expect(err).NotTo(HaveOccurred())

			build, found, err = job1.GetNextPendingBuildBySerialGroup([]string{"serial-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildThree.ID()))

			err = job1.Unpause()
			Expect(err).NotTo(HaveOccurred())

			build, found, err = job2.GetNextPendingBuildBySerialGroup([]string{"serial-group", "really-different-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildOne.ID()))

			Expect(buildOne.Finish(db.BuildStatusSucceeded)).To(Succeed())

			build, found, err = job1.GetNextPendingBuildBySerialGroup([]string{"serial-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildTwo.ID()))

			build, found, err = job2.GetNextPendingBuildBySerialGroup([]string{"serial-group", "really-different-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildTwo.ID()))

			scheduled, err := buildTwo.Schedule()
			Expect(err).NotTo(HaveOccurred())
			Expect(scheduled).To(BeTrue())
			Expect(buildTwo.Finish(db.BuildStatusSucceeded)).To(Succeed())

			build, found, err = job1.GetNextPendingBuildBySerialGroup([]string{"serial-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildThree.ID()))

			build, found, err = job2.GetNextPendingBuildBySerialGroup([]string{"serial-group", "really-different-group"})
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(build.ID()).To(Equal(buildThree.ID()))
		})
	})

	Describe("saving build inputs", func() {
		var (
			buildMetadata []db.ResourceMetadataField
			vr1           db.VersionedResource
		)

		BeforeEach(func() {
			buildMetadata = []db.ResourceMetadataField{
				{
					Name:  "meta1",
					Value: "value1",
				},
				{
					Name:  "meta2",
					Value: "value2",
				},
			}

			vr1 = db.VersionedResource{
				Resource: "some-other-resource",
				Type:     "some-type",
				Version:  db.ResourceVersion{"ver": "2"},
			}
		})

		It("fails to save build input if resource does not exist", func() {
			build, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			vr := db.VersionedResource{
				Resource: "unknown-resource",
				Type:     "some-type",
				Version:  db.ResourceVersion{"ver": "2"},
			}

			input := db.BuildInput{
				Name:              "some-input",
				VersionedResource: vr,
			}

			err = build.SaveInput(input)
			Expect(err).To(HaveOccurred())
		})

		It("updates metadata of existing versioned resources", func() {
			build, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			err = build.SaveInput(db.BuildInput{
				Name:              "some-input",
				VersionedResource: vr1,
			})
			Expect(err).NotTo(HaveOccurred())

			inputs, _, err := build.Resources()
			Expect(err).NotTo(HaveOccurred())
			Expect(inputs).To(ConsistOf([]db.BuildInput{
				{Name: "some-input", VersionedResource: vr1, FirstOccurrence: true},
			}))

			withMetadata := vr1
			withMetadata.Metadata = buildMetadata

			err = build.SaveInput(db.BuildInput{
				Name:              "some-other-input",
				VersionedResource: withMetadata,
			})
			Expect(err).NotTo(HaveOccurred())

			inputs, _, err = build.Resources()
			Expect(err).NotTo(HaveOccurred())
			Expect(inputs).To(ConsistOf([]db.BuildInput{
				{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
				{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
			}))

			err = build.SaveInput(db.BuildInput{
				Name:              "some-input",
				VersionedResource: withMetadata,
			})
			Expect(err).NotTo(HaveOccurred())

			inputs, _, err = build.Resources()
			Expect(err).NotTo(HaveOccurred())
			Expect(inputs).To(ConsistOf([]db.BuildInput{
				{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
				{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
			}))

		})

		It("does not clobber metadata of existing versioned resources", func() {
			build, err := jobCombination.CreateBuild()
			Expect(err).NotTo(HaveOccurred())

			withMetadata := vr1
			withMetadata.Metadata = buildMetadata

			withoutMetadata := vr1
			withoutMetadata.Metadata = nil

			err = build.SaveInput(db.BuildInput{
				Name:              "some-input",
				VersionedResource: withMetadata,
			})
			Expect(err).NotTo(HaveOccurred())

			inputs, _, err := build.Resources()
			Expect(err).NotTo(HaveOccurred())
			Expect(inputs).To(ConsistOf([]db.BuildInput{
				{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
			}))

			err = build.SaveInput(db.BuildInput{
				Name:              "some-other-input",
				VersionedResource: withoutMetadata,
			})
			Expect(err).NotTo(HaveOccurred())

			inputs, _, err = build.Resources()
			Expect(err).NotTo(HaveOccurred())
			Expect(inputs).To(ConsistOf([]db.BuildInput{
				{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
				{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
			}))
		})
	})

	Describe("a build is created for a job", func() {
		var (
			build1DB      db.Build
			otherPipeline db.Pipeline
			otherJob      db.Job
		)

		BeforeEach(func() {
			pipelineConfig := atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",
					},
				},
				Resources: atc.ResourceConfigs{
					{
						Name: "some-other-resource",
						Type: "some-type",
					},
				},
			}
			var err error
			otherPipeline, _, err = team.SavePipeline("some-other-pipeline", pipelineConfig, db.ConfigVersion(1), db.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())

			build1DB, err = jobCombination.CreateBuild()
			Expect(err).ToNot(HaveOccurred())

			Expect(build1DB.ID()).NotTo(BeZero())
			Expect(build1DB.JobName()).To(Equal("some-job"))
			Expect(build1DB.Name()).To(Equal("1"))
			Expect(build1DB.Status()).To(Equal(db.BuildStatusPending))
			Expect(build1DB.IsScheduled()).To(BeFalse())

			var found bool
			otherJob, found, err = otherPipeline.Job("some-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
		})

		It("becomes the next pending build for job", func() {
			nextPendings, err := job.GetPendingBuilds()
			Expect(err).NotTo(HaveOccurred())
			Expect(nextPendings).NotTo(BeEmpty())
			Expect(nextPendings[0].ID()).To(Equal(build1DB.ID()))
		})

		It("is in the list of pending builds", func() {
			nextPendingBuilds, err := pipeline.GetAllPendingBuilds()
			Expect(err).NotTo(HaveOccurred())
			Expect(nextPendingBuilds["some-job"]).To(HaveLen(1))
			Expect(nextPendingBuilds["some-job"]).To(Equal([]db.Build{build1DB}))
		})

		Context("and another build for a different pipeline is created with the same job name", func() {
			BeforeEach(func() {
				otherJobCombination, err := otherJob.JobCombination()
				Expect(err).ToNot(HaveOccurred())

				otherBuild, err := otherJobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				Expect(otherBuild.ID()).NotTo(BeZero())
				Expect(otherBuild.JobName()).To(Equal("some-job"))
				Expect(otherBuild.Name()).To(Equal("1"))
				Expect(otherBuild.Status()).To(Equal(db.BuildStatusPending))
				Expect(otherBuild.IsScheduled()).To(BeFalse())
			})

			It("does not change the next pending build for job", func() {
				nextPendingBuilds, err := job.GetPendingBuilds()
				Expect(err).NotTo(HaveOccurred())
				Expect(nextPendingBuilds).To(Equal([]db.Build{build1DB}))
			})

			It("does not change pending builds", func() {
				nextPendingBuilds, err := pipeline.GetAllPendingBuilds()
				Expect(err).NotTo(HaveOccurred())
				Expect(nextPendingBuilds["some-job"]).To(HaveLen(1))
				Expect(nextPendingBuilds["some-job"]).To(Equal([]db.Build{build1DB}))
			})
		})

		Context("when scheduled", func() {
			BeforeEach(func() {
				var err error
				var found bool
				found, err = build1DB.Schedule()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
			})

			It("remains the next pending build for job", func() {
				nextPendingBuilds, err := job.GetPendingBuilds()
				Expect(err).NotTo(HaveOccurred())
				Expect(nextPendingBuilds).NotTo(BeEmpty())
				Expect(nextPendingBuilds[0].ID()).To(Equal(build1DB.ID()))
			})

			It("remains in the list of pending builds", func() {
				nextPendingBuilds, err := pipeline.GetAllPendingBuilds()
				Expect(err).NotTo(HaveOccurred())
				Expect(nextPendingBuilds["some-job"]).To(HaveLen(1))
				Expect(nextPendingBuilds["some-job"][0].ID()).To(Equal(build1DB.ID()))
			})
		})

		Context("when started", func() {
			BeforeEach(func() {
				started, err := build1DB.Start("some-engine", `{"some":"metadata"}`, atc.Plan{})
				Expect(err).NotTo(HaveOccurred())
				Expect(started).To(BeTrue())
			})

			It("saves the updated status, and the engine and engine metadata", func() {
				found, err := build1DB.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build1DB.Status()).To(Equal(db.BuildStatusStarted))
				Expect(build1DB.Engine()).To(Equal("some-engine"))
				Expect(build1DB.EngineMetadata()).To(Equal(`{"some":"metadata"}`))
			})

			It("saves the build's start time", func() {
				found, err := build1DB.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build1DB.StartTime().Unix()).To(BeNumerically("~", time.Now().Unix(), 3))
			})
		})

		Context("when the build finishes", func() {
			BeforeEach(func() {
				err := build1DB.Finish(db.BuildStatusSucceeded)
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets the build's status and end time", func() {
				found, err := build1DB.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build1DB.Status()).To(Equal(db.BuildStatusSucceeded))
				Expect(build1DB.EndTime().Unix()).To(BeNumerically("~", time.Now().Unix(), 3))
			})
		})

		Context("and another is created for the same job", func() {
			var build2DB db.Build

			BeforeEach(func() {
				var err error
				build2DB, err = jobCombination.CreateBuild()
				Expect(err).NotTo(HaveOccurred())

				Expect(build2DB.ID()).NotTo(BeZero())
				Expect(build2DB.ID()).NotTo(Equal(build1DB.ID()))
				Expect(build2DB.Name()).To(Equal("2"))
				Expect(build2DB.Status()).To(Equal(db.BuildStatusPending))
			})

			Describe("the first build", func() {
				It("remains the next pending build", func() {
					nextPendingBuilds, err := job.GetPendingBuilds()
					Expect(err).NotTo(HaveOccurred())
					Expect(nextPendingBuilds).To(HaveLen(2))
					Expect(nextPendingBuilds[0].ID()).To(Equal(build1DB.ID()))
					Expect(nextPendingBuilds[1].ID()).To(Equal(build2DB.ID()))
				})

				It("remains in the list of pending builds", func() {
					nextPendingBuilds, err := pipeline.GetAllPendingBuilds()
					Expect(err).NotTo(HaveOccurred())
					Expect(nextPendingBuilds["some-job"]).To(HaveLen(2))
					Expect(nextPendingBuilds["some-job"]).To(ConsistOf(build1DB, build2DB))
				})
			})
		})
	})

	Describe("SyncResourceSpaceCombinations", func() {
		BeforeEach(func() {
			otherPipeline, created, err := team.SavePipeline("other-fake-pipeline", atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "some-job",

						Public: true,

						Serial: true,

						SerialGroups: []string{"serial-group"},

						Plan: atc.PlanSequence{
							{
								Put: "some-resource",
								Params: atc.Params{
									"some-param": "some-value",
								},
							},
							{
								Get:      "some-input",
								Resource: "some-resource",
								Params: atc.Params{
									"some-param": "some-value",
								},
								Passed:  []string{"job-1", "job-2"},
								Trigger: true,
							},
							{
								Task:           "some-task",
								Privileged:     true,
								TaskConfigPath: "some/config/path.yml",
								TaskConfig: &atc.TaskConfig{
									RootfsURI: "some-image",
								},
							},
						},
					},
					{
						Name: "some-other-job",
					},
					{
						Name:         "other-serial-group-job",
						SerialGroups: []string{"serial-group", "really-different-group"},
					},
					{
						Name:         "different-serial-group-job",
						SerialGroups: []string{"different-serial-group"},
					},
				},
				Resources: atc.ResourceConfigs{
					{
						Name: "some-resource",
						Type: "some-type",
					},
					{
						Name: "some-other-resource",
						Type: "some-type",
					},
				},
			}, db.ConfigVersion(0), db.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())
			Expect(created).To(BeTrue())

			var found bool
			job, found, err = otherPipeline.Job("some-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())

			someResource, found, err := otherPipeline.Resource("some-resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			otherResource, found, err := otherPipeline.Resource("some-other-resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			err = otherPipeline.SaveSpaces(someResource, []string{"some-space"})
			Expect(err).NotTo(HaveOccurred())

			err = otherPipeline.SaveSpaces(otherResource, []string{"some-other-space", "some-another-space"})
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a job_resource_space_combination for each space in every combination", func() {
			combination1 := map[string]string{"some-resource": "some-space", "some-other-resource": "some-other-space"}
			combination2 := map[string]string{"some-resource": "some-space", "some-other-resource": "some-another-space"}

			jobCombinations, err := job.SyncResourceSpaceCombinations([]map[string]string{combination1, combination2})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobCombinations).To(HaveLen(2))
			Expect(jobCombinations[0].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[0].Combination()).To(Equal(combination1))
			Expect(jobCombinations[1].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[1].Combination()).To(Equal(combination2))
		})

		It("Creates an empty job_combination when no combinations are provided", func() {
			combination := map[string]string{}

			jobCombinations, err := job.SyncResourceSpaceCombinations([]map[string]string{combination})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobCombinations).To(HaveLen(1))
			Expect(jobCombinations[0].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[0].Combination()).To(Equal(combination))
		})

		It("Updates a job_combination when its combination is null", func() {
			tx, err := dbConn.Begin()
			Expect(err).NotTo(HaveOccurred())

			_, err = psql.Insert("job_combinations").
				Columns("job_id").
				Values(job.ID()).
				RunWith(tx).
				Exec()
			Expect(err).NotTo(HaveOccurred())

			err = tx.Commit()
			Expect(err).NotTo(HaveOccurred())

			db.Rollback(tx)

			combination := map[string]string{"some-resource": "some-space", "some-other-resource": "some-other-space"}

			jobCombinations, err := job.SyncResourceSpaceCombinations([]map[string]string{combination})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobCombinations).To(HaveLen(1))
			Expect(jobCombinations[0].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[0].Combination()).To(Equal(combination))
		})

		It("Create new job_combinations when combinations have changed", func() {
			combination1 := map[string]string{"some-resource": "some-space", "some-other-resource": "some-other-space"}
			combination2 := map[string]string{"some-resource": "some-space", "some-other-resource": "some-another-space"}

			_, err := job.SyncResourceSpaceCombinations([]map[string]string{combination1})
			Expect(err).NotTo(HaveOccurred())

			jobCombinations, err := job.JobCombinations()
			Expect(err).NotTo(HaveOccurred())
			oldCombinationCount := len(jobCombinations)

			jobCombinations, err = job.SyncResourceSpaceCombinations([]map[string]string{combination2})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobCombinations).To(HaveLen(1))
			Expect(jobCombinations[0].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[0].Combination()).To(Equal(combination2))

			jobCombinations, err = job.JobCombinations()
			Expect(err).NotTo(HaveOccurred())
			newCombinationCount := len(jobCombinations)
			Expect(newCombinationCount).To(Equal(oldCombinationCount + 1))
		})

		It("Does not create new job_combinations when the same combinations already exist", func() {
			combination := map[string]string{"some-resource": "some-space", "some-other-resource": "some-other-space"}

			_, err := job.SyncResourceSpaceCombinations([]map[string]string{combination})
			Expect(err).NotTo(HaveOccurred())

			jobCombinations, err := job.JobCombinations()
			Expect(err).NotTo(HaveOccurred())
			oldCombinationCount := len(jobCombinations)

			jobCombinations, err = job.SyncResourceSpaceCombinations([]map[string]string{combination})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobCombinations).To(HaveLen(1))
			Expect(jobCombinations[0].JobID()).To(Equal(job.ID()))
			Expect(jobCombinations[0].Combination()).To(Equal(combination))

			jobCombinations, err = job.JobCombinations()
			Expect(err).NotTo(HaveOccurred())
			newCombinationCount := len(jobCombinations)
			Expect(newCombinationCount).To(Equal(oldCombinationCount))
		})
	})

	Describe("ResourceSpaceCombinations", func() {
		var (
			resourceSpaces map[string][]string
			combinations   []map[string]string
		)

		JustBeforeEach(func() {
			combinations = job.ResourceSpaceCombinations(resourceSpaces)
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
