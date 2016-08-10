package db_test

import (
	"errors"
	"time"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/event"
	"github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PipelineDB", func() {
	var dbConn db.Conn
	var listener *pq.Listener

	var pipelineDBFactory db.PipelineDBFactory
	var sqlDB *db.SQLDB
	var teamDBFactory db.TeamDBFactory

	BeforeEach(func() {
		postgresRunner.Truncate()

		dbConn = db.Wrap(postgresRunner.Open())

		listener = pq.NewListener(postgresRunner.DataSourceName(), time.Second, time.Minute, nil)
		Eventually(listener.Ping, 5*time.Second).ShouldNot(HaveOccurred())
		bus := db.NewNotificationsBus(listener, dbConn)

		sqlDB = db.NewSQL(dbConn, bus)
		pipelineDBFactory = db.NewPipelineDBFactory(dbConn, bus)
		teamDBFactory = db.NewTeamDBFactory(dbConn, bus)
	})

	AfterEach(func() {
		err := dbConn.Close()
		Expect(err).NotTo(HaveOccurred())

		err = listener.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	pipelineConfig := atc.Config{
		Groups: atc.GroupConfigs{
			{
				Name:      "some-group",
				Jobs:      []string{"job-1", "job-2"},
				Resources: []string{"some-resource", "some-other-resource"},
			},
		},

		Resources: atc.ResourceConfigs{
			{
				Name: "some-resource",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
			{
				Name: "some-other-resource",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
			{
				Name: "some-really-other-resource",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
		},

		ResourceTypes: atc.ResourceTypes{
			{
				Name: "some-resource-type",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
		},

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
							Image: "some-image",
						},
					},
				},
			},
			{
				Name:   "some-other-job",
				Serial: true,
			},
			{
				Name: "a-job",
			},
			{
				Name: "shared-job",
			},
			{
				Name: "random-job",
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
	}

	otherPipelineConfig := atc.Config{
		Groups: atc.GroupConfigs{
			{
				Name:      "some-group",
				Jobs:      []string{"job-1", "job-2"},
				Resources: []string{"some-resource", "some-other-resource"},
			},
		},

		Resources: atc.ResourceConfigs{
			{
				Name: "some-resource",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
			{
				Name: "some-other-resource",
				Type: "some-type",
				Source: atc.Source{
					"source-config": "some-value",
				},
			},
		},

		Jobs: atc.JobConfigs{
			{
				Name: "some-job",
			},
			{
				Name: "some-other-job",
			},
			{
				Name: "a-job",
			},
			{
				Name: "shared-job",
			},
			{
				Name: "other-serial-group-job",
			},
		},
	}

	var (
		teamDB             db.TeamDB
		pipelineDB         db.PipelineDB
		otherPipelineDB    db.PipelineDB
		savedPipeline      db.SavedPipeline
		otherSavedPipeline db.SavedPipeline
	)

	BeforeEach(func() {
		var err error
		_, err = sqlDB.CreateTeam(db.Team{Name: "some-team"})
		Expect(err).NotTo(HaveOccurred())

		teamDB = teamDBFactory.GetTeamDB("some-team")

		savedPipeline, _, err = teamDB.SaveConfig("a-pipeline-name", pipelineConfig, 0, db.PipelineUnpaused)
		Expect(err).NotTo(HaveOccurred())

		otherSavedPipeline, _, err = teamDB.SaveConfig("other-pipeline-name", otherPipelineConfig, 0, db.PipelineUnpaused)
		Expect(err).NotTo(HaveOccurred())

		pipelineDB = pipelineDBFactory.Build(savedPipeline)
		otherPipelineDB = pipelineDBFactory.Build(otherSavedPipeline)
	})

	Describe("destroying a pipeline", func() {
		It("can be deleted", func() {
			// populate pipelines table
			pipelineThatWillBeDeleted, _, err := teamDB.SaveConfig("a-pipeline-that-will-be-deleted", pipelineConfig, 0, db.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			fetchedPipeline, found, err := teamDB.GetPipelineByName("a-pipeline-that-will-be-deleted")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			fetchedPipelineDB := pipelineDBFactory.Build(fetchedPipeline)

			// populate resources table and versioned_resources table

			savedResource, _, err := fetchedPipelineDB.GetResource("some-resource")
			Expect(err).NotTo(HaveOccurred())

			resourceConfig, found := pipelineConfig.Resources.Lookup("some-resource")
			Expect(found).To(BeTrue())

			fetchedPipelineDB.SaveResourceVersions(resourceConfig, []atc.Version{
				{
					"key": "value",
				},
			})

			// populate builds table
			build, err := fetchedPipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			oneOffBuild, err := teamDB.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())

			// populate jobs_serial_groups table
			_, err = fetchedPipelineDB.GetRunningBuildsBySerialGroup("some-job", []string{"serial-group"})
			Expect(err).NotTo(HaveOccurred())

			// populate build_inputs table
			_, err = fetchedPipelineDB.SaveInput(build.ID(), db.BuildInput{
				Name: "build-input",
				VersionedResource: db.VersionedResource{
					Resource:   "some-resource",
					PipelineID: savedPipeline.ID,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// In very old concourse deployments, build inputs and outputs seem to
			// have been created for one-off builds. This test makes sure they get
			// deleted. See story #109558152
			_, err = fetchedPipelineDB.SaveInput(oneOffBuild.ID(), db.BuildInput{
				Name: "one-off-build-input",
				VersionedResource: db.VersionedResource{
					Resource:   "some-resource",
					PipelineID: pipelineThatWillBeDeleted.ID,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			// populate build_outputs table

			_, err = fetchedPipelineDB.SaveOutput(build.ID(), db.VersionedResource{
				Resource:   "some-resource",
				PipelineID: pipelineThatWillBeDeleted.ID,
			}, false)
			Expect(err).NotTo(HaveOccurred())

			_, err = fetchedPipelineDB.SaveOutput(oneOffBuild.ID(), db.VersionedResource{
				Resource:   "some-resource",
				PipelineID: pipelineThatWillBeDeleted.ID,
			}, false)
			Expect(err).NotTo(HaveOccurred())

			// populate build_events table
			err = build.SaveEvent(event.StartTask{})
			Expect(err).NotTo(HaveOccurred())

			// populate image_resource_versions table
			err = build.SaveImageResourceVersion("some-plan-id", db.ResourceCacheIdentifier{
				ResourceVersion: atc.Version{"digest": "readers"},
				ResourceHash:    `docker{"some":"source"}`,
			})
			Expect(err).NotTo(HaveOccurred())

			err = fetchedPipelineDB.Destroy()
			Expect(err).NotTo(HaveOccurred())

			pipelines, err := sqlDB.GetAllPipelines()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipelines).NotTo(ContainElement(fetchedPipeline))

			_, _, found, err = fetchedPipelineDB.GetConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			resourceRows, err := dbConn.Query(`select id from resources where pipeline_id = $1`, fetchedPipeline.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resourceRows.Next()).To(BeFalse())

			resourceRows.Close()

			versionRows, err := dbConn.Query(`select id from versioned_resources where resource_id = $1`, savedResource.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(versionRows.Next()).To(BeFalse())

			versionRows.Close()

			buildRows, err := dbConn.Query(`select id from builds where id = $1`, build.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(buildRows.Next()).To(BeFalse())

			buildRows.Close()

			jobRows, err := dbConn.Query(`select id from jobs where pipeline_id = $1`, fetchedPipeline.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(jobRows.Next()).To(BeFalse())

			jobRows.Close()

			eventRows, err := dbConn.Query(`select build_id from build_events where build_id = $1`, build.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(eventRows.Next()).To(BeFalse())

			eventRows.Close()

			inputRows, err := dbConn.Query(`select build_id from build_inputs where build_id = $1`, build.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(inputRows.Next()).To(BeFalse())

			inputRows.Close()

			oneOffInputRows, err := dbConn.Query(`select build_id from build_inputs where build_id = $1`, oneOffBuild.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(oneOffInputRows.Next()).To(BeFalse())

			oneOffInputRows.Close()

			outputRows, err := dbConn.Query(`select build_id from build_outputs where build_id = $1`, build.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(outputRows.Next()).To(BeFalse())

			outputRows.Close()

			oneOffOutputRows, err := dbConn.Query(`select build_id from build_outputs where build_id = $1`, oneOffBuild.ID())
			Expect(err).NotTo(HaveOccurred())
			Expect(oneOffOutputRows.Next()).To(BeFalse())

			oneOffOutputRows.Close()

			foundImageVolumeIdentifiers, err := build.GetImageResourceCacheIdentifiers()
			Expect(err).NotTo(HaveOccurred())
			Expect(foundImageVolumeIdentifiers).To(BeEmpty())
		})
	})

	Describe("pausing and unpausing a pipeline", func() {
		It("starts out as unpaused", func() {
			Expect(savedPipeline.Paused).To(BeFalse())
		})

		It("can be paused", func() {
			err := pipelineDB.Pause()
			Expect(err).NotTo(HaveOccurred())

			pipelinePaused, err := pipelineDB.IsPaused()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipelinePaused).To(BeTrue())

			otherPipelinePaused, err := otherPipelineDB.IsPaused()
			Expect(err).NotTo(HaveOccurred())
			Expect(otherPipelinePaused).To(BeFalse())
		})

		It("can be unpaused", func() {
			err := pipelineDB.Pause()
			Expect(err).NotTo(HaveOccurred())

			err = otherPipelineDB.Pause()
			Expect(err).NotTo(HaveOccurred())

			err = pipelineDB.Unpause()
			Expect(err).NotTo(HaveOccurred())

			pipelinePaused, err := pipelineDB.IsPaused()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipelinePaused).To(BeFalse())

			otherPipelinePaused, err := otherPipelineDB.IsPaused()
			Expect(err).NotTo(HaveOccurred())
			Expect(otherPipelinePaused).To(BeTrue())
		})
	})

	Describe("UpdateName", func() {
		var teamDB db.TeamDB

		BeforeEach(func() {
			teamDB = teamDBFactory.GetTeamDB("some-team")
		})

		It("can update the name of a given pipeline", func() {
			err := pipelineDB.UpdateName("some-other-weird-name")
			Expect(err).NotTo(HaveOccurred())

			pipeline, found, err := teamDB.GetPipelineByName("some-other-weird-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(pipeline.Name).To(Equal("some-other-weird-name"))
		})

		Context("when there is a pipeline with the same name in another team", func() {
			var team2 db.SavedTeam
			var team2DB db.TeamDB

			BeforeEach(func() {
				var err error
				team2, err = sqlDB.CreateTeam(db.Team{Name: "some-other-team"})
				Expect(err).NotTo(HaveOccurred())

				team2DB = teamDBFactory.GetTeamDB(team2.Name)
				_, _, err = team2DB.SaveConfig("a-pipeline-name", pipelineConfig, 0, db.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())
			})

			It("doesn't rename the other pipeline", func() {
				err := pipelineDB.UpdateName("some-other-weird-name")
				Expect(err).NotTo(HaveOccurred())

				_, _, err = team2DB.GetPipelineByName("a-pipeline-name")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("ScopedName", func() {
		It("concatenates the pipeline name with the passed in name", func() {
			pipelineDB := pipelineDBFactory.Build(db.SavedPipeline{
				Pipeline: db.Pipeline{
					Name: "some-pipeline",
				},
			})
			Expect(pipelineDB.ScopedName("something-else")).To(Equal("some-pipeline:something-else"))
		})
	})

	Describe("getting the pipeline configuration", func() {
		It("can manage multiple pipeline configurations", func() {
			By("returning the saved config to later gets")
			returnedConfig, configVersion, found, err := pipelineDB.GetConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(returnedConfig).To(Equal(pipelineConfig))
			Expect(configVersion).NotTo(Equal(db.ConfigVersion(0)))

			otherReturnedConfig, otherConfigVersion, found, err := otherPipelineDB.GetConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(otherReturnedConfig).To(Equal(otherPipelineConfig))
			Expect(otherConfigVersion).NotTo(Equal(db.ConfigVersion(0)))

			updatedConfig := pipelineConfig

			updatedConfig.Groups = append(pipelineConfig.Groups, atc.GroupConfig{
				Name: "new-group",
				Jobs: []string{"new-job-1", "new-job-2"},
			})

			updatedConfig.Resources = append(pipelineConfig.Resources, atc.ResourceConfig{
				Name: "new-resource",
				Type: "new-type",
				Source: atc.Source{
					"new-source-config": "new-value",
				},
			})

			updatedConfig.Jobs = append(pipelineConfig.Jobs, atc.JobConfig{
				Name: "new-job",
				Plan: atc.PlanSequence{
					{
						Get:      "new-input",
						Resource: "new-resource",
						Params: atc.Params{
							"new-param": "new-value",
						},
					},
					{
						Task:           "some-task",
						TaskConfigPath: "new/config/path.yml",
					},
				},
			})

			By("being able to update the config with a valid config")
			_, _, err = teamDB.SaveConfig("a-pipeline-name", updatedConfig, configVersion, db.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())
			_, _, err = teamDB.SaveConfig("other-pipeline-name", updatedConfig, otherConfigVersion, db.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			By("returning the updated config")
			returnedConfig, newConfigVersion, found, err := pipelineDB.GetConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(returnedConfig).To(Equal(updatedConfig))
			Expect(newConfigVersion).NotTo(Equal(configVersion))

			otherReturnedConfig, newOtherConfigVersion, found, err := otherPipelineDB.GetConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(otherReturnedConfig).To(Equal(updatedConfig))
			Expect(newOtherConfigVersion).NotTo(Equal(otherConfigVersion))
		})
	})

	Context("Resources", func() {
		resourceName := "some-resource"
		otherResourceName := "some-other-resource"
		reallyOtherResourceName := "some-really-other-resource"

		var resource db.SavedResource
		var otherResource db.SavedResource
		var reallyOtherResource db.SavedResource

		BeforeEach(func() {
			var err error
			resource, _, err = pipelineDB.GetResource(resourceName)
			Expect(err).NotTo(HaveOccurred())

			otherResource, _, err = pipelineDB.GetResource(otherResourceName)
			Expect(err).NotTo(HaveOccurred())

			reallyOtherResource, _, err = pipelineDB.GetResource(reallyOtherResourceName)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns correct resource", func() {
			Expect(resource).To(Equal(db.SavedResource{
				ID: resource.ID,
				Resource: db.Resource{
					Name: "some-resource",
				},
				Paused:       false,
				PipelineName: "a-pipeline-name",
				CheckError:   nil,
			}))
		})

		Context("SaveResourceVersions", func() {
			var (
				originalVersionSlice []atc.Version
				resourceConfig       atc.ResourceConfig
			)

			BeforeEach(func() {
				resourceConfig = atc.ResourceConfig{
					Name:   resource.Name,
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}

				originalVersionSlice = []atc.Version{
					{"ref": "v1"},
					{"ref": "v3"},
				}
			})

			It("ensures versioned resources have the correct check_order", func() {
				err := pipelineDB.SaveResourceVersions(resourceConfig, originalVersionSlice)
				Expect(err).NotTo(HaveOccurred())

				latestVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(latestVR.Version).To(Equal(db.Version{"ref": "v3"}))
				Expect(latestVR.CheckOrder).To(Equal(2))

				pretendCheckResults := []atc.Version{
					{"ref": "v2"},
					{"ref": "v3"},
				}

				err = pipelineDB.SaveResourceVersions(resourceConfig, pretendCheckResults)
				Expect(err).NotTo(HaveOccurred())

				latestVR, found, err = pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(latestVR.Version).To(Equal(db.Version{"ref": "v3"}))
				Expect(latestVR.CheckOrder).To(Equal(4))
			})

			Context("resource not found in db", func() {
				BeforeEach(func() {
					resourceConfig = atc.ResourceConfig{
						Name:   "unknown-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					}

					originalVersionSlice = []atc.Version{{"ref": "v1"}}
				})

				It("returns an error", func() {
					err := pipelineDB.SaveResourceVersions(resourceConfig, originalVersionSlice)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("resource 'unknown-resource' not found"))
				})
			})
		})

		It("can load up versioned resource information relevant to scheduling", func() {
			job, err := pipelineDB.GetJob("some-job")
			Expect(err).NotTo(HaveOccurred())

			otherJob, err := pipelineDB.GetJob("some-other-job")
			Expect(err).NotTo(HaveOccurred())

			aJob, err := pipelineDB.GetJob("a-job")
			Expect(err).NotTo(HaveOccurred())

			sharedJob, err := pipelineDB.GetJob("shared-job")
			Expect(err).NotTo(HaveOccurred())

			randomJob, err := pipelineDB.GetJob("random-job")
			Expect(err).NotTo(HaveOccurred())

			otherSerialGroupJob, err := pipelineDB.GetJob("other-serial-group-job")
			Expect(err).NotTo(HaveOccurred())

			differentSerialGroupJob, err := pipelineDB.GetJob("different-serial-group-job")
			Expect(err).NotTo(HaveOccurred())

			versions, err := pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(BeEmpty())
			Expect(versions.BuildOutputs).To(BeEmpty())
			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"some-other-job":             otherJob.ID,
				"a-job":                      aJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("initially having no latest versioned resource")
			_, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			By("including saved versioned resources of the current pipeline")
			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR1, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(savedVR1.ModifiedTime).NotTo(BeNil())
			Expect(savedVR1.ModifiedTime).To(BeTemporally(">", time.Time{}))

			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "2"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR2, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(ConsistOf([]algorithm.ResourceVersion{
				{VersionID: savedVR1.ID, ResourceID: resource.ID, CheckOrder: savedVR1.CheckOrder},
				{VersionID: savedVR2.ID, ResourceID: resource.ID, CheckOrder: savedVR2.CheckOrder},
			}))

			Expect(versions.BuildOutputs).To(BeEmpty())
			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"some-other-job":             otherJob.ID,
				"a-job":                      aJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("not including saved versioned resources of other pipelines")
			otherPipelineResource, _, err := otherPipelineDB.GetResource("some-other-resource")
			Expect(err).NotTo(HaveOccurred())

			err = otherPipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   otherPipelineResource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}})
			Expect(err).NotTo(HaveOccurred())

			otherPipelineSavedVR, found, err := otherPipelineDB.GetLatestVersionedResource(otherPipelineResource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(ConsistOf([]algorithm.ResourceVersion{
				{VersionID: savedVR1.ID, ResourceID: resource.ID, CheckOrder: savedVR1.CheckOrder},
				{VersionID: savedVR2.ID, ResourceID: resource.ID, CheckOrder: savedVR2.CheckOrder},
			}))

			Expect(versions.BuildOutputs).To(BeEmpty())
			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"some-other-job":             otherJob.ID,
				"a-job":                      aJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("including outputs of successful builds")
			build1DB, err := pipelineDB.CreateJobBuild("a-job")
			Expect(err).NotTo(HaveOccurred())

			savedVR1, err = pipelineDB.SaveOutput(build1DB.ID(), savedVR1.VersionedResource, false)
			Expect(err).NotTo(HaveOccurred())

			err = build1DB.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(ConsistOf([]algorithm.ResourceVersion{
				{VersionID: savedVR1.ID, ResourceID: resource.ID, CheckOrder: savedVR1.CheckOrder},
				{VersionID: savedVR2.ID, ResourceID: resource.ID, CheckOrder: savedVR2.CheckOrder},
			}))

			Expect(versions.BuildOutputs).To(ConsistOf([]algorithm.BuildOutput{
				{
					ResourceVersion: algorithm.ResourceVersion{
						VersionID:  savedVR1.ID,
						ResourceID: resource.ID,
						CheckOrder: savedVR1.CheckOrder,
					},
					JobID:   aJob.ID,
					BuildID: build1DB.ID(),
				},
			}))

			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"a-job":                      aJob.ID,
				"some-other-job":             otherJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("not including outputs of failed builds")
			build2DB, err := pipelineDB.CreateJobBuild("a-job")
			Expect(err).NotTo(HaveOccurred())

			savedVR1, err = pipelineDB.SaveOutput(build2DB.ID(), savedVR1.VersionedResource, false)
			Expect(err).NotTo(HaveOccurred())

			err = build2DB.Finish(db.StatusFailed)
			Expect(err).NotTo(HaveOccurred())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(ConsistOf([]algorithm.ResourceVersion{
				{VersionID: savedVR1.ID, ResourceID: resource.ID, CheckOrder: savedVR1.CheckOrder},
				{VersionID: savedVR2.ID, ResourceID: resource.ID, CheckOrder: savedVR2.CheckOrder},
			}))

			Expect(versions.BuildOutputs).To(ConsistOf([]algorithm.BuildOutput{
				{
					ResourceVersion: algorithm.ResourceVersion{
						VersionID:  savedVR1.ID,
						ResourceID: resource.ID,
						CheckOrder: savedVR1.CheckOrder,
					},
					JobID:   aJob.ID,
					BuildID: build1DB.ID(),
				},
			}))

			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"a-job":                      aJob.ID,
				"some-other-job":             otherJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("not including outputs of builds in other pipelines")
			otherPipelineBuild, err := otherPipelineDB.CreateJobBuild("a-job")
			Expect(err).NotTo(HaveOccurred())

			_, err = otherPipelineDB.SaveOutput(otherPipelineBuild.ID(), otherPipelineSavedVR.VersionedResource, false)
			Expect(err).NotTo(HaveOccurred())

			err = otherPipelineBuild.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())
			Expect(versions.ResourceVersions).To(ConsistOf([]algorithm.ResourceVersion{
				{VersionID: savedVR1.ID, ResourceID: resource.ID, CheckOrder: savedVR1.CheckOrder},
				{VersionID: savedVR2.ID, ResourceID: resource.ID, CheckOrder: savedVR2.CheckOrder},
			}))

			Expect(versions.BuildOutputs).To(ConsistOf([]algorithm.BuildOutput{
				{
					ResourceVersion: algorithm.ResourceVersion{
						VersionID:  savedVR1.ID,
						ResourceID: resource.ID,
						CheckOrder: savedVR1.CheckOrder,
					},
					JobID:   aJob.ID,
					BuildID: build1DB.ID(),
				},
			}))

			Expect(versions.ResourceIDs).To(Equal(map[string]int{
				resource.Name:            resource.ID,
				otherResource.Name:       otherResource.ID,
				reallyOtherResource.Name: reallyOtherResource.ID,
			}))

			Expect(versions.JobIDs).To(Equal(map[string]int{
				"some-job":                   job.ID,
				"a-job":                      aJob.ID,
				"some-other-job":             otherJob.ID,
				"shared-job":                 sharedJob.ID,
				"random-job":                 randomJob.ID,
				"other-serial-group-job":     otherSerialGroupJob.ID,
				"different-serial-group-job": differentSerialGroupJob.ID,
			}))

			By("including build inputs")
			build1DB, err = pipelineDB.CreateJobBuild("a-job")
			Expect(err).NotTo(HaveOccurred())

			savedVR1, err = pipelineDB.SaveInput(build1DB.ID(), db.BuildInput{
				Name:              "some-input-name",
				VersionedResource: savedVR1.VersionedResource,
			})
			Expect(err).NotTo(HaveOccurred())

			err = build1DB.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			versions, err = pipelineDB.LoadVersionsDB()
			Expect(err).NotTo(HaveOccurred())

			Expect(versions.BuildInputs).To(ConsistOf([]algorithm.BuildInput{
				{
					ResourceVersion: algorithm.ResourceVersion{
						VersionID:  savedVR1.ID,
						ResourceID: resource.ID,
						CheckOrder: savedVR1.CheckOrder,
					},
					JobID:     aJob.ID,
					BuildID:   build1DB.ID(),
					InputName: "some-input-name",
				},
			}))
		})

		Context("when a version is disabled", func() {
			It("omits the version from the versions DB", func() {
				build1, err := pipelineDB.CreateJobBuild("a-job")
				Expect(err).NotTo(HaveOccurred())

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   resource.Name,
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "disabled"}})
				Expect(err).NotTo(HaveOccurred())

				disabledVersion, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				_, err = build1.SaveInput(db.BuildInput{
					Name:              "disabled-input",
					VersionedResource: disabledVersion.VersionedResource,
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = build1.SaveOutput(disabledVersion.VersionedResource, false)
				Expect(err).NotTo(HaveOccurred())

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   resource.Name,
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "enabled"}})
				Expect(err).NotTo(HaveOccurred())

				enabledVersion, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				_, err = build1.SaveInput(db.BuildInput{
					Name:              "enabled-input",
					VersionedResource: enabledVersion.VersionedResource,
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = build1.SaveOutput(enabledVersion.VersionedResource, false)
				Expect(err).NotTo(HaveOccurred())

				err = build1.Finish(db.StatusSucceeded)
				Expect(err).NotTo(HaveOccurred())

				pipelineDB.DisableVersionedResource(disabledVersion.ID)

				pipelineDB.DisableVersionedResource(enabledVersion.ID)
				pipelineDB.EnableVersionedResource(enabledVersion.ID)

				versions, err := pipelineDB.LoadVersionsDB()
				Expect(err).NotTo(HaveOccurred())

				aJob, err := pipelineDB.GetJob("a-job")
				Expect(err).NotTo(HaveOccurred())

				By("omitting it from the list of resource versions")
				Expect(versions.ResourceVersions).To(ConsistOf(
					algorithm.ResourceVersion{
						VersionID:  enabledVersion.ID,
						ResourceID: resource.ID,
						CheckOrder: enabledVersion.CheckOrder,
					},
				))

				By("omitting it from build outputs")
				Expect(versions.BuildOutputs).To(ConsistOf(
					algorithm.BuildOutput{
						ResourceVersion: algorithm.ResourceVersion{
							VersionID:  enabledVersion.ID,
							ResourceID: resource.ID,
							CheckOrder: enabledVersion.CheckOrder,
						},
						JobID:   aJob.ID,
						BuildID: build1.ID(),
					},
				))

				By("omitting it from build inputs")
				Expect(versions.BuildInputs).To(ConsistOf(
					algorithm.BuildInput{
						ResourceVersion: algorithm.ResourceVersion{
							VersionID:  enabledVersion.ID,
							ResourceID: resource.ID,
							CheckOrder: enabledVersion.CheckOrder,
						},
						JobID:     aJob.ID,
						BuildID:   build1.ID(),
						InputName: "enabled-input",
					},
				))
			})
		})

		Describe("GetVersionedResourceByVersion", func() {
			var savedVersion2 db.SavedVersionedResource
			BeforeEach(func() {
				err := pipelineDB.SaveResourceVersions(
					atc.ResourceConfig{
						Name: "some-resource",
						Type: "some-type",
						Source: atc.Source{
							"source-config": "some-value",
						},
					},
					[]atc.Version{
						{"version": "v1"},
						{"version": "v2"},
						{"version": "v3"}, // disabled
					},
				)
				Expect(err).NotTo(HaveOccurred())

				// save metadata for v2
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).ToNot(HaveOccurred())
				_, err = build.SaveInput(db.BuildInput{
					Name: "some-input",
					VersionedResource: db.VersionedResource{
						Resource:   "some-resource",
						Type:       "some-type",
						Version:    db.Version{"version": "v2"},
						Metadata:   []db.MetadataField{{Name: "name1", Value: "value1"}},
						PipelineID: pipelineDB.GetPipelineID(),
					},
					FirstOccurrence: true,
				})
				Expect(err).NotTo(HaveOccurred())

				savedVersions, _, found, err := pipelineDB.GetResourceVersions("some-resource", db.Page{Limit: 2})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(savedVersions).To(HaveLen(2))
				pipelineDB.DisableVersionedResource(savedVersions[0].ID)
				savedVersion2 = savedVersions[1]

				err = pipelineDB.SaveResourceVersions(
					atc.ResourceConfig{
						Name: "some-other-resource",
						Type: "some-type",
						Source: atc.Source{
							"source-config": "some-value",
						},
					},
					[]atc.Version{
						{"version": "v2"},
					},
				)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the SavedVersionedResource matching the given resource name and atc version", func() {
				By("returning versions that exist")
				actualSavedVersion, found, err := pipelineDB.GetVersionedResourceByVersion(
					atc.Version{"version": "v2"},
					"some-resource",
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(actualSavedVersion).To(Equal(savedVersion2))

				By("returning not found for versions that don't exist")
				_, found, err = pipelineDB.GetVersionedResourceByVersion(
					atc.Version{"versioni": "v2"},
					"some-resource",
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())

				By("returning not found for versions that only exist in another resource")
				_, found, err = pipelineDB.GetVersionedResourceByVersion(
					atc.Version{"version": "v1"},
					"some-other-resource",
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())

				By("returning not found for disabled versions")
				_, found, err = pipelineDB.GetVersionedResourceByVersion(
					atc.Version{"version": "v3"},
					"some-resource",
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})

		It("can load up the latest enabled versioned resource", func() {
			By("initially having no latest versioned resource")
			_, found, err := pipelineDB.GetLatestEnabledVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			By("including saved versioned resources of the current pipeline")
			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR1, found, err := pipelineDB.GetLatestEnabledVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "2"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR2, found, err := pipelineDB.GetLatestEnabledVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(savedVR1.Version).To(Equal(db.Version{"version": "1"}))
			Expect(savedVR1.PipelineID).To(Equal(pipelineDB.GetPipelineID()))
			Expect(savedVR2.Version).To(Equal(db.Version{"version": "2"}))
			Expect(savedVR2.PipelineID).To(Equal(pipelineDB.GetPipelineID()))

			By("not including saved versioned resources of other pipelines")
			_, _, err = otherPipelineDB.GetResource("some-other-resource")
			Expect(err).NotTo(HaveOccurred())

			err = otherPipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}, {"version": "2"}, {"version": "3"}})
			Expect(err).NotTo(HaveOccurred())

			otherPipelineSavedVR, found, err := otherPipelineDB.GetLatestEnabledVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(otherPipelineSavedVR.Version).To(Equal(db.Version{"version": "3"}))
			Expect(otherPipelineSavedVR.PipelineID).To(Equal(otherPipelineDB.GetPipelineID()))

			By("not including disabled versions")
			err = pipelineDB.DisableVersionedResource(savedVR2.ID)
			Expect(err).NotTo(HaveOccurred())

			savedVR3, found, err := pipelineDB.GetLatestEnabledVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(savedVR3.Version).To(Equal(db.Version{"version": "1"}))
		})

		It("can load up the latest versioned resource, enabled or not", func() {
			By("initially having no latest versioned resource")
			_, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())

			By("including saved versioned resources of the current pipeline")
			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR1, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "2"}})
			Expect(err).NotTo(HaveOccurred())

			savedVR2, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(savedVR1.Version).To(Equal(db.Version{"version": "1"}))
			Expect(savedVR1.PipelineID).To(Equal(pipelineDB.GetPipelineID()))
			Expect(savedVR2.Version).To(Equal(db.Version{"version": "2"}))
			Expect(savedVR2.PipelineID).To(Equal(pipelineDB.GetPipelineID()))

			By("not including saved versioned resources of other pipelines")
			_, _, err = otherPipelineDB.GetResource("some-other-resource")
			Expect(err).NotTo(HaveOccurred())

			err = otherPipelineDB.SaveResourceVersions(atc.ResourceConfig{
				Name:   resource.Name,
				Type:   "some-type",
				Source: atc.Source{"some": "source"},
			}, []atc.Version{{"version": "1"}, {"version": "2"}, {"version": "3"}})
			Expect(err).NotTo(HaveOccurred())

			otherPipelineSavedVR, found, err := otherPipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(otherPipelineSavedVR.Version).To(Equal(db.Version{"version": "3"}))
			Expect(otherPipelineSavedVR.PipelineID).To(Equal(otherPipelineDB.GetPipelineID()))

			By("including disabled versions")
			err = pipelineDB.DisableVersionedResource(savedVR2.ID)
			Expect(err).NotTo(HaveOccurred())

			latestVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			Expect(latestVR.Version).To(Equal(db.Version{"version": "2"}))
		})

		Describe("pausing and unpausing resources", func() {
			It("starts out as unpaused", func() {
				resource, _, err := pipelineDB.GetResource(resourceName)
				Expect(err).NotTo(HaveOccurred())

				Expect(resource.Paused).To(BeFalse())
			})

			It("can be paused", func() {
				err := pipelineDB.PauseResource(resourceName)
				Expect(err).NotTo(HaveOccurred())

				pausedResource, _, err := pipelineDB.GetResource(resourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(pausedResource.Paused).To(BeTrue())

				resource, _, err := otherPipelineDB.GetResource(resourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(resource.Paused).To(BeFalse())
			})

			It("can be unpaused", func() {
				err := pipelineDB.PauseResource(resourceName)
				Expect(err).NotTo(HaveOccurred())

				err = otherPipelineDB.PauseResource(resourceName)
				Expect(err).NotTo(HaveOccurred())

				err = pipelineDB.UnpauseResource(resourceName)
				Expect(err).NotTo(HaveOccurred())

				unpausedResource, _, err := pipelineDB.GetResource(resourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(unpausedResource.Paused).To(BeFalse())

				resource, _, err := otherPipelineDB.GetResource(resourceName)
				Expect(err).NotTo(HaveOccurred())
				Expect(resource.Paused).To(BeTrue())
			})
		})

		Describe("enabling and disabling versioned resources", func() {
			It("returns an error if the resource or version is bogus", func() {
				err := pipelineDB.EnableVersionedResource(42)
				Expect(err).To(HaveOccurred())

				err = pipelineDB.DisableVersionedResource(42)
				Expect(err).To(HaveOccurred())
			})

			It("does not affect explicitly fetching the latest version", func() {
				err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "1"}})
				Expect(err).NotTo(HaveOccurred())

				savedVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(savedVR.Resource).To(Equal("some-resource"))
				Expect(savedVR.Type).To(Equal("some-type"))
				Expect(savedVR.Version).To(Equal(db.Version{"version": "1"}))
				initialTime := savedVR.ModifiedTime

				err = pipelineDB.DisableVersionedResource(savedVR.ID)
				Expect(err).NotTo(HaveOccurred())

				disabledVR := savedVR
				disabledVR.Enabled = false

				latestVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(found).To(BeTrue())
				Expect(latestVR.Resource).To(Equal(disabledVR.Resource))
				Expect(latestVR.Type).To(Equal(disabledVR.Type))
				Expect(latestVR.Version).To(Equal(disabledVR.Version))
				Expect(latestVR.Enabled).To(BeFalse())
				Expect(latestVR.ModifiedTime).To(BeTemporally(">", initialTime))

				tmp_modified_time := latestVR.ModifiedTime

				err = pipelineDB.EnableVersionedResource(savedVR.ID)
				Expect(err).NotTo(HaveOccurred())

				enabledVR := savedVR
				enabledVR.Enabled = true

				latestVR, found, err = pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(found).To(BeTrue())
				Expect(latestVR.Resource).To(Equal(enabledVR.Resource))
				Expect(latestVR.Type).To(Equal(enabledVR.Type))
				Expect(latestVR.Version).To(Equal(enabledVR.Version))
				Expect(latestVR.Enabled).To(BeTrue())
				Expect(latestVR.ModifiedTime).To(BeTemporally(">", tmp_modified_time))
			})

			It("doesn't change the check_order when saving a new build input", func() {
				err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "1"},
					{"version": "2"},
					{"version": "3"},
				})
				Expect(err).NotTo(HaveOccurred())

				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				beforeVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "4"},
					{"version": "5"},
				})
				Expect(err).NotTo(HaveOccurred())

				input := db.BuildInput{
					Name:              "input-name",
					VersionedResource: beforeVR.VersionedResource,
				}

				afterVR, err := pipelineDB.SaveInput(build.ID(), input)
				Expect(err).NotTo(HaveOccurred())

				Expect(afterVR.CheckOrder).To(Equal(beforeVR.CheckOrder))
			})

			It("doesn't change the check_order when saving a new implicit build output", func() {
				err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "1"},
					{"version": "2"},
					{"version": "3"},
				})
				Expect(err).NotTo(HaveOccurred())

				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				beforeVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "4"},
					{"version": "5"},
				})
				Expect(err).NotTo(HaveOccurred())

				afterVR, err := pipelineDB.SaveOutput(build.ID(), beforeVR.VersionedResource, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(afterVR.CheckOrder).To(Equal(beforeVR.CheckOrder))
			})

			It("doesn't change the check_order when saving a new implicit build output", func() {
				err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "1"},
					{"version": "2"},
					{"version": "3"},
				})
				Expect(err).NotTo(HaveOccurred())

				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				beforeVR, found, err := pipelineDB.GetLatestVersionedResource(resource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{
					{"version": "4"},
					{"version": "5"},
				})
				Expect(err).NotTo(HaveOccurred())

				afterVR, err := pipelineDB.SaveOutput(build.ID(), beforeVR.VersionedResource, true)
				Expect(err).NotTo(HaveOccurred())

				Expect(afterVR.CheckOrder).To(Equal(beforeVR.CheckOrder))
			})
		})

		Describe("VersionsDB caching", func() {
			Context("when build outputs are added", func() {
				var build db.Build
				var savedVR db.SavedVersionedResource

				BeforeEach(func() {
					var err error
					build, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
						Name:   "some-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					}, []atc.Version{{"version": "1"}})
					Expect(err).NotTo(HaveOccurred())

					savedResource, _, err := pipelineDB.GetResource("some-resource")
					Expect(err).NotTo(HaveOccurred())

					var found bool
					savedVR, found, err = pipelineDB.GetLatestVersionedResource(savedResource.Name)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
				})

				It("will cache VersionsDB if no change has occured", func() {
					_, err := pipelineDB.SaveOutput(build.ID(), savedVR.VersionedResource, true)

					versionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())

					cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())
					Expect(versionsDB == cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be the same object")
				})

				It("will not cache VersionsDB if a change occured", func() {
					versionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())

					_, err = pipelineDB.SaveOutput(build.ID(), savedVR.VersionedResource, true)
					Expect(err).NotTo(HaveOccurred())

					cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())
					Expect(versionsDB != cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be different objects")
				})

				Context("when the build outputs are added for a different pipeline", func() {
					It("does not invalidate the cache for the original pipeline", func() {
						otherBuild, err := otherPipelineDB.CreateJobBuild("some-job")
						Expect(err).NotTo(HaveOccurred())

						err = otherPipelineDB.SaveResourceVersions(atc.ResourceConfig{
							Name:   "some-other-resource",
							Type:   "some-type",
							Source: atc.Source{"some": "source"},
						}, []atc.Version{{"version": "1"}})
						Expect(err).NotTo(HaveOccurred())

						otherSavedResource, _, err := otherPipelineDB.GetResource("some-other-resource")
						Expect(err).NotTo(HaveOccurred())

						otherSavedVR, found, err := otherPipelineDB.GetLatestVersionedResource(otherSavedResource.Name)
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeTrue())

						versionsDB, err := pipelineDB.LoadVersionsDB()
						Expect(err).NotTo(HaveOccurred())

						_, err = otherPipelineDB.SaveOutput(otherBuild.ID(), otherSavedVR.VersionedResource, true)
						Expect(err).NotTo(HaveOccurred())

						cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
						Expect(err).NotTo(HaveOccurred())
						Expect(versionsDB == cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be the same object")
					})
				})
			})

			Context("when versioned resources are added", func() {
				It("will cache VersionsDB if no change has occured", func() {
					err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
						Name:   "some-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					}, []atc.Version{{"version": "1"}})
					Expect(err).NotTo(HaveOccurred())

					versionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())

					cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())
					Expect(versionsDB == cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be the same object")
				})

				It("will not cache VersionsDB if a change occured", func() {
					err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
						Name:   "some-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					}, []atc.Version{{"version": "1"}})
					Expect(err).NotTo(HaveOccurred())

					versionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())

					err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
						Name:   "some-other-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					}, []atc.Version{{"version": "1"}})
					Expect(err).NotTo(HaveOccurred())

					cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
					Expect(err).NotTo(HaveOccurred())
					Expect(versionsDB != cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be different objects")
				})

				Context("when the versioned resources are added for a different pipeline", func() {
					It("does not invalidate the cache for the original pipeline", func() {
						err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
							Name:   "some-resource",
							Type:   "some-type",
							Source: atc.Source{"some": "source"},
						}, []atc.Version{{"version": "1"}})
						Expect(err).NotTo(HaveOccurred())

						versionsDB, err := pipelineDB.LoadVersionsDB()
						Expect(err).NotTo(HaveOccurred())

						err = otherPipelineDB.SaveResourceVersions(atc.ResourceConfig{
							Name:   "some-other-resource",
							Type:   "some-type",
							Source: atc.Source{"some": "source"},
						}, []atc.Version{{"version": "1"}})
						Expect(err).NotTo(HaveOccurred())

						cachedVersionsDB, err := pipelineDB.LoadVersionsDB()
						Expect(err).NotTo(HaveOccurred())
						Expect(versionsDB == cachedVersionsDB).To(BeTrue(), "Expected VersionsDB to be the same object")
					})
				})
			})
		})

		Describe("saving versioned resources", func() {
			It("updates the latest versioned resource", func() {
				err := pipelineDB.SaveResourceVersions(
					atc.ResourceConfig{
						Name:   "some-resource",
						Type:   "some-type",
						Source: atc.Source{"some": "source"},
					},
					[]atc.Version{{"version": "1"}},
				)
				Expect(err).NotTo(HaveOccurred())

				savedResource, _, err := pipelineDB.GetResource("some-resource")
				Expect(err).NotTo(HaveOccurred())

				savedVR, found, err := pipelineDB.GetLatestVersionedResource(savedResource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(savedVR.Resource).To(Equal("some-resource"))
				Expect(savedVR.Type).To(Equal("some-type"))
				Expect(savedVR.Version).To(Equal(db.Version{"version": "1"}))

				err = pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "2"}, {"version": "3"}})
				Expect(err).NotTo(HaveOccurred())

				savedVR, found, err = pipelineDB.GetLatestVersionedResource(savedResource.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(savedVR.Resource).To(Equal("some-resource"))
				Expect(savedVR.Type).To(Equal("some-type"))
				Expect(savedVR.Version).To(Equal(db.Version{"version": "3"}))
			})
		})

		It("initially reports zero builds for a job", func() {
			builds, err := pipelineDB.GetAllJobBuilds("some-job")
			Expect(err).NotTo(HaveOccurred())
			Expect(builds).To(BeEmpty())
		})

		It("initially has no pending build for a job", func() {
			_, found, err := pipelineDB.GetNextPendingBuild("some-job")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		Describe("marking resource checks as errored", func() {
			var resource db.SavedResource

			BeforeEach(func() {
				var err error
				resource, _, err = pipelineDB.GetResource("some-resource")
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the resource is first created", func() {
				It("is not errored", func() {
					Expect(resource.CheckError).To(BeNil())
				})
			})

			Context("when a resource check is marked as errored", func() {
				It("is then marked as errored", func() {
					originalCause := errors.New("on fire")

					err := pipelineDB.SetResourceCheckError(resource, originalCause)
					Expect(err).NotTo(HaveOccurred())

					returnedResource, _, err := pipelineDB.GetResource("some-resource")
					Expect(err).NotTo(HaveOccurred())

					Expect(returnedResource.CheckError).To(Equal(originalCause))
				})
			})

			Context("when a resource is cleared of check errors", func() {
				It("is not marked as errored again", func() {
					originalCause := errors.New("on fire")

					err := pipelineDB.SetResourceCheckError(resource, originalCause)
					Expect(err).NotTo(HaveOccurred())

					err = pipelineDB.SetResourceCheckError(resource, nil)
					Expect(err).NotTo(HaveOccurred())

					returnedResource, _, err := pipelineDB.GetResource("some-resource")
					Expect(err).NotTo(HaveOccurred())

					Expect(returnedResource.CheckError).To(BeNil())
				})
			})
		})
	})

	Describe("GetResourceType", func() {
		It("returns a SavedResourceType", func() {
			_, found, err := pipelineDB.GetResourceType("resource-type-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		Context("when the resource type has a version", func() {
			BeforeEach(func() {
				resourceType := atc.ResourceType{
					Name: "some-resource-type",
					Type: "some-type",
				}
				err := pipelineDB.SaveResourceTypeVersion(resourceType, atc.Version{"foo": "bar"})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns a SavedResourceType", func() {
				savedResourceType, found, err := pipelineDB.GetResourceType("some-resource-type")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(savedResourceType.Name).To(Equal("some-resource-type"))
				Expect(savedResourceType.Type).To(Equal("some-type"))
				Expect(savedResourceType.Version).To(Equal(db.Version{"foo": "bar"}))
			})
		})

		It("returns a SavedResourceType with no version when the resource type does not have a version", func() {
			savedResourceType, found, err := pipelineDB.GetResourceType("some-resource-type")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(savedResourceType.Name).To(Equal("some-resource-type"))
			Expect(savedResourceType.Type).To(Equal("some-type"))
			Expect(savedResourceType.Version).To(BeNil())
		})

		It("returns not found when the resource type cannot be found", func() {
			_, found, err := pipelineDB.GetResourceType("nonexistent-resource-type")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("GetResources", func() {
		var (
			dashboardResource1 db.DashboardResource
			dashboardResource2 db.DashboardResource
			dashboardResource3 db.DashboardResource
		)
		BeforeEach(func() {
			dashboardResource1 = db.DashboardResource{
				Resource: db.SavedResource{
					ID:           0,
					CheckError:   nil,
					Paused:       false,
					PipelineName: "a-pipeline-name",
					Resource:     db.Resource{Name: "some-resource"},
				},
				ResourceConfig: atc.ResourceConfig{
					Name: "some-resource",
					Type: "some-type",
					Source: atc.Source{
						"source-config": "some-value",
					},
				},
			}

			dashboardResource2 = db.DashboardResource{
				Resource: db.SavedResource{
					ID:           0,
					CheckError:   nil,
					Paused:       false,
					PipelineName: "a-pipeline-name",
					Resource:     db.Resource{Name: "some-other-resource"},
				},
				ResourceConfig: atc.ResourceConfig{
					Name: "some-other-resource",
					Type: "some-type",
					Source: atc.Source{
						"source-config": "some-value",
					},
				},
			}

			dashboardResource3 = db.DashboardResource{
				Resource: db.SavedResource{
					ID:           0,
					CheckError:   nil,
					Paused:       false,
					PipelineName: "a-pipeline-name",
					Resource:     db.Resource{Name: "some-really-other-resource"},
				},
				ResourceConfig: atc.ResourceConfig{
					Name: "some-really-other-resource",
					Type: "some-type",
					Source: atc.Source{
						"source-config": "some-value",
					},
				},
			}
		})
		Context("when there is a one-to-one correspondence between saved and configured resources", func() {
			It("returns all the saved resources", func() {
				dashboardResources, groupConfigs, found, err := pipelineDB.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				for i, _ := range dashboardResources {
					dashboardResources[i].Resource.ID = 0
				}

				Expect(dashboardResources).To(HaveLen(3))
				Expect(dashboardResources).To(ConsistOf(dashboardResource1, dashboardResource2, dashboardResource3))
				Expect(groupConfigs).To(Equal(atc.GroupConfigs{
					{
						Name:      "some-group",
						Jobs:      []string{"job-1", "job-2"},
						Resources: []string{"some-resource", "some-other-resource"},
					},
				}))
			})
		})

		Context("when there is a saved resource that is not configured", func() {
			BeforeEach(func() {
				pipelineConfigMinusResource := pipelineConfig
				pipelineConfigMinusResource.Resources = []atc.ResourceConfig{
					pipelineConfig.Resources[0],
					pipelineConfig.Resources[2],
				}

				_, _, err := teamDB.SaveConfig("a-pipeline-name", pipelineConfigMinusResource, 1, db.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the configured resources", func() {
				dashboardResources, groupConfigs, found, err := pipelineDB.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				for i, _ := range dashboardResources {
					dashboardResources[i].Resource.ID = 0
				}

				Expect(dashboardResources).To(HaveLen(2))
				Expect(dashboardResources).To(ConsistOf(dashboardResource1, dashboardResource3))
				Expect(groupConfigs).To(Equal(atc.GroupConfigs{
					{
						Name:      "some-group",
						Jobs:      []string{"job-1", "job-2"},
						Resources: []string{"some-resource", "some-other-resource"},
					},
				}))
			})
		})

		Context("when there is a configured resource that is not saved for some reason", func() {
			BeforeEach(func() {
				result, err := dbConn.Exec(`
					delete from resources
						where pipeline_id = $1
						and name = 'some-resource'
					`, pipelineDB.GetPipelineID())
				Expect(err).NotTo(HaveOccurred())
				Expect(result.RowsAffected()).To(Equal(int64(1)))
			})
			It("returns an error", func() {
				_, _, _, err := pipelineDB.GetResources()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(
					errors.New("found resource in pipeline configuration but not in database: some-resource"),
				))
			})
		})

		Context("when the pipeline has been deleted in the meantime", func() {
			BeforeEach(func() {
				err := pipelineDB.Destroy()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns not found", func() {
				_, _, found, err := pipelineDB.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("GetResource", func() {
		It("returns not found when the resource type cannot be found", func() {
			_, found, err := pipelineDB.GetResource("nonexistent-resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		Context("when resource exists", func() {
			BeforeEach(func() {
				err := pipelineDB.SaveResourceVersions(atc.ResourceConfig{
					Name:   "some-resource",
					Type:   "some-type",
					Source: atc.Source{"some": "source"},
				}, []atc.Version{{"version": "1"}})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns SavedResource when it exists", func() {
				savedResource, found, err := pipelineDB.GetResource("some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(savedResource.Resource.Name).To(Equal("some-resource"))
				Expect(savedResource.PipelineName).To(Equal("a-pipeline-name"))
			})
		})
	})

	Describe("SaveResourceTypeVersion", func() {
		Context("when resource type does not exist in database", func() {
			It("returns an error", func() {
				resourceType := atc.ResourceType{
					Name: "other-resource-type",
					Type: "resource-type-type",
				}
				err := pipelineDB.SaveResourceTypeVersion(resourceType, atc.Version{"foo": "bar"})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when resource type exists in database", func() {
			var resourceType atc.ResourceType
			BeforeEach(func() {
				resourceType = atc.ResourceType{
					Name: "some-resource-type",
					Type: "some-type",
				}

				err := pipelineDB.SaveResourceTypeVersion(resourceType, atc.Version{"foo": "bar"})
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates resource type", func() {
				err := pipelineDB.SaveResourceTypeVersion(resourceType, atc.Version{"baz": "qux"})
				Expect(err).NotTo(HaveOccurred())

				var savedResourceTypeName string
				var savedResourceTypeType string
				var versionJSON string
				err = dbConn.QueryRow(`
					SELECT name, type, version
					FROM resource_types
				`).Scan(&savedResourceTypeName, &savedResourceTypeType, &versionJSON)
				Expect(err).NotTo(HaveOccurred())
				Expect(savedResourceTypeName).To(Equal("some-resource-type"))
				Expect(savedResourceTypeType).To(Equal("some-type"))
				Expect(versionJSON).To(MatchJSON(`{"baz":"qux"}`))
			})
		})
	})

	Describe("Jobs", func() {
		Describe("GetDashboard", func() {
			It("returns a Dashboard object with a DashboardJob corresponding to each configured job", func() {
				pipelineDB.UpdateFirstLoggedBuildID("some-job", 57)

				job, err := pipelineDB.GetJob("some-job")
				Expect(err).NotTo(HaveOccurred())

				otherJob, err := pipelineDB.GetJob("some-other-job")
				Expect(err).NotTo(HaveOccurred())

				aJob, err := pipelineDB.GetJob("a-job")
				Expect(err).NotTo(HaveOccurred())

				sharedJob, err := pipelineDB.GetJob("shared-job")
				Expect(err).NotTo(HaveOccurred())

				randomJob, err := pipelineDB.GetJob("random-job")
				Expect(err).NotTo(HaveOccurred())

				otherSerialGroupJob, err := pipelineDB.GetJob("other-serial-group-job")
				Expect(err).NotTo(HaveOccurred())

				differentSerialGroupJob, err := pipelineDB.GetJob("different-serial-group-job")
				Expect(err).NotTo(HaveOccurred())

				lookupConfig := func(jobName string) atc.JobConfig {
					config, found := pipelineConfig.Jobs.Lookup(jobName)
					Expect(found).To(BeTrue())

					return config
				}

				By("returning jobs with no builds")
				expectedDashboard := db.Dashboard{
					{
						JobConfig:     lookupConfig(job.Name),
						Job:           job,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(otherJob.Name),
						Job:           otherJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(aJob.Name),
						Job:           aJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(sharedJob.Name),
						Job:           sharedJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(randomJob.Name),
						Job:           randomJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(otherSerialGroupJob.Name),
						Job:           otherSerialGroupJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
					{
						JobConfig:     lookupConfig(differentSerialGroupJob.Name),
						Job:           differentSerialGroupJob,
						NextBuild:     nil,
						FinishedBuild: nil,
					},
				}

				actualDashboard, groups, err := pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(groups).To(Equal(pipelineConfig.Groups))
				Expect(actualDashboard).To(ConsistOf(expectedDashboard))

				By("returning a job's most recent pending build if there are no running builds")
				jobBuildOldDB, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				expectedDashboard[0].NextBuild = jobBuildOldDB

				actualDashboard, _, err = pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(actualDashboard).To(ConsistOf(expectedDashboard))

				By("returning a job's most recent started build")
				jobBuildOldDB.Start("engine", "metadata")

				found, err := jobBuildOldDB.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				expectedDashboard[0].NextBuild = jobBuildOldDB

				actualDashboard, _, err = pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(actualDashboard).To(ConsistOf(expectedDashboard))

				By("returning a job's most recent started build even if there is a newer pending build")
				jobBuild, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				expectedDashboard[0].NextBuild = jobBuildOldDB

				actualDashboard, _, err = pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(actualDashboard).To(ConsistOf(expectedDashboard))

				By("returning a job's most recent finished build")
				err = jobBuild.Finish(db.StatusSucceeded)
				Expect(err).NotTo(HaveOccurred())

				found, err = jobBuild.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				expectedDashboard[0].FinishedBuild = jobBuild
				expectedDashboard[0].NextBuild = jobBuildOldDB

				actualDashboard, _, err = pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(actualDashboard).To(ConsistOf(expectedDashboard))

				By("returning a job's most recent finished build even when there is a newer unfinished build")
				jobBuildNewDB, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())
				jobBuildNewDB.Start("engine", "metadata")
				found, err = jobBuildNewDB.Reload()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				expectedDashboard[0].FinishedBuild = jobBuild
				expectedDashboard[0].NextBuild = jobBuildNewDB

				actualDashboard, _, err = pipelineDB.GetDashboard()
				Expect(err).NotTo(HaveOccurred())

				Expect(actualDashboard).To(ConsistOf(expectedDashboard))
			})
		})

		Describe("CreateJobBuild", func() {
			var build db.Build

			BeforeEach(func() {
				var err error
				build, err = pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets the properties of a build for a given job", func() {
				Expect(build.ID()).NotTo(BeZero())
				Expect(build.JobName()).To(Equal("some-job"))
				Expect(build.Name()).To(Equal("1"))
				Expect(build.Status()).To(Equal(db.StatusPending))
				Expect(build.IsScheduled()).To(BeFalse())
				Expect(build.TeamName()).To(Equal("some-team"))
			})
		})

		Describe("saving build inputs", func() {
			var (
				buildMetadata []db.MetadataField
				vr1           db.VersionedResource
				vr2           db.VersionedResource
			)

			BeforeEach(func() {
				buildMetadata = []db.MetadataField{
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
					PipelineID: savedPipeline.ID,
					Resource:   "some-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "1"},
					Metadata:   buildMetadata,
				}

				vr2 = db.VersionedResource{
					PipelineID: savedPipeline.ID,
					Resource:   "some-other-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "2"},
				}
			})

			It("fails to save build input if resource does not exist", func() {
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				vr := db.VersionedResource{
					PipelineID: savedPipeline.ID,
					Resource:   "unknown-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "2"},
				}

				input := db.BuildInput{
					Name:              "some-input",
					VersionedResource: vr,
				}

				_, err = build.SaveInput(input)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("resource 'unknown-resource' not found"))
			})

			It("updates metadata of existing versioned resources", func() {
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				_, err = build.SaveInput(db.BuildInput{
					Name:              "some-input",
					VersionedResource: vr2,
				})
				Expect(err).NotTo(HaveOccurred())

				inputs, _, err := build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: vr2, FirstOccurrence: true},
				}))

				withMetadata := vr2
				withMetadata.Metadata = buildMetadata

				_, err = build.SaveInput(db.BuildInput{
					Name:              "some-other-input",
					VersionedResource: withMetadata,
				})
				Expect(err).NotTo(HaveOccurred())

				inputs, _, err = build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))

				_, err = build.SaveInput(db.BuildInput{
					Name:              "some-input",
					VersionedResource: withMetadata,
				})
				Expect(err).NotTo(HaveOccurred())

				inputs, _, err = build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))

			})

			It("does not clobber metadata of existing versioned resources", func() {
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				withMetadata := vr2
				withMetadata.Metadata = buildMetadata

				withoutMetadata := vr2
				withoutMetadata.Metadata = nil

				savedVR, err := build.SaveInput(db.BuildInput{
					Name:              "some-input",
					VersionedResource: withMetadata,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(savedVR.Metadata).To(Equal(buildMetadata))

				inputs, _, err := build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))

				savedVR, err = build.SaveInput(db.BuildInput{
					Name:              "some-other-input",
					VersionedResource: withoutMetadata,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(savedVR.Metadata).To(Equal(buildMetadata))

				inputs, _, err = build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: withMetadata, FirstOccurrence: true},
					{Name: "some-other-input", VersionedResource: withMetadata, FirstOccurrence: true},
				}))

			})
		})

		Describe("saving inputs, implicit outputs, and explicit outputs", func() {
			var (
				vr1 db.VersionedResource
				vr2 db.VersionedResource
			)

			BeforeEach(func() {
				vr1 = db.VersionedResource{
					PipelineID: savedPipeline.ID,
					Resource:   "some-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "1"},
				}

				vr2 = db.VersionedResource{
					PipelineID: savedPipeline.ID,
					Resource:   "some-other-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "2"},
				}
			})

			It("correctly distinguishes them", func() {
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				// save a normal 'get'
				_, err = build.SaveInput(db.BuildInput{
					Name:              "some-input",
					VersionedResource: vr1,
				})
				Expect(err).NotTo(HaveOccurred())

				// save implicit output from 'get'
				_, err = build.SaveOutput(vr1, false)
				Expect(err).NotTo(HaveOccurred())

				// save explicit output from 'put'
				_, err = build.SaveOutput(vr2, true)
				Expect(err).NotTo(HaveOccurred())

				// save the dependent get
				_, err = build.SaveInput(db.BuildInput{
					Name:              "some-dependent-input",
					VersionedResource: vr2,
				})
				Expect(err).NotTo(HaveOccurred())

				// save the dependent 'get's implicit output
				_, err = build.SaveOutput(vr2, false)
				Expect(err).NotTo(HaveOccurred())

				inputs, outputs, err := build.GetResources()
				Expect(err).NotTo(HaveOccurred())
				Expect(inputs).To(ConsistOf([]db.BuildInput{
					{Name: "some-input", VersionedResource: vr1, FirstOccurrence: true},
				}))

				Expect(outputs).To(ConsistOf([]db.BuildOutput{
					{VersionedResource: vr2},
				}))

			})

			It("fails to save build output if resource does not exist", func() {
				build, err := pipelineDB.CreateJobBuild("some-job")
				Expect(err).NotTo(HaveOccurred())

				vr := db.VersionedResource{
					PipelineID: savedPipeline.ID,
					Resource:   "unknown-resource",
					Type:       "some-type",
					Version:    db.Version{"ver": "2"},
				}

				_, err = build.SaveOutput(vr, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("resource 'unknown-resource' not found"))
			})
		})

		Describe("pausing and unpausing jobs", func() {
			job := "some-job"

			It("starts out as unpaused", func() {
				job, err := pipelineDB.GetJob(job)
				Expect(err).NotTo(HaveOccurred())

				Expect(job.Paused).To(BeFalse())
			})

			It("can be paused", func() {
				err := pipelineDB.PauseJob(job)
				Expect(err).NotTo(HaveOccurred())

				err = otherPipelineDB.UnpauseJob(job)
				Expect(err).NotTo(HaveOccurred())

				pausedJob, err := pipelineDB.GetJob(job)
				Expect(err).NotTo(HaveOccurred())
				Expect(pausedJob.Paused).To(BeTrue())

				otherJob, err := otherPipelineDB.GetJob(job)
				Expect(err).NotTo(HaveOccurred())
				Expect(otherJob.Paused).To(BeFalse())
			})

			It("can be unpaused", func() {
				err := pipelineDB.PauseJob(job)
				Expect(err).NotTo(HaveOccurred())

				err = pipelineDB.UnpauseJob(job)
				Expect(err).NotTo(HaveOccurred())

				unpausedJob, err := pipelineDB.GetJob(job)
				Expect(err).NotTo(HaveOccurred())

				Expect(unpausedJob.Paused).To(BeFalse())
			})
		})

		Describe("UpdateFirstLoggedBuildID", func() {
			It("updates FirstLoggedBuildID on a job", func() {
				By("starting out as 0")
				job, err := pipelineDB.GetJob("some-job")
				Expect(err).NotTo(HaveOccurred())

				Expect(job.FirstLoggedBuildID).To(BeZero())

				By("increasing it to 57")

				err = pipelineDB.UpdateFirstLoggedBuildID("some-job", 57)
				Expect(err).NotTo(HaveOccurred())

				updatedJob, err := pipelineDB.GetJob("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedJob.FirstLoggedBuildID).To(Equal(57))

				By("not erroring when it's called with the same number")
				err = pipelineDB.UpdateFirstLoggedBuildID("some-job", 57)
				Expect(err).NotTo(HaveOccurred())

				By("erroring when the number decreases")
				err = pipelineDB.UpdateFirstLoggedBuildID("some-job", 56)
				Expect(err).To(Equal(db.FirstLoggedBuildIDDecreasedError{
					Job:   "some-job",
					OldID: 57,
					NewID: 56,
				}))
			})
		})

		Describe("GetJobBuild", func() {
			var firstBuild db.Build
			var job db.SavedJob

			BeforeEach(func() {
				var err error
				job, err = pipelineDB.GetJob("some-job")
				Expect(err).NotTo(HaveOccurred())

				firstBuild, err = pipelineDB.CreateJobBuild(job.Name)
				Expect(err).NotTo(HaveOccurred())
			})

			It("finds the build", func() {
				build, found, err := pipelineDB.GetJobBuild(job.Name, firstBuild.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(firstBuild.ID()))
				Expect(build.Status()).To(Equal(firstBuild.Status()))
			})
		})

		Describe("GetNextPendingBuildBySerialGroup", func() {
			var jobOneConfig atc.JobConfig
			var jobOneTwoConfig atc.JobConfig

			BeforeEach(func() {
				jobOneConfig = pipelineConfig.Jobs[0]    //serial-group
				jobOneTwoConfig = pipelineConfig.Jobs[5] //serial-group, really-different-group
			})

			Context("when some jobs have builds with inputs determined as false", func() {
				var actualBuild db.Build

				BeforeEach(func() {
					_, err := pipelineDB.CreateJobBuild(jobOneConfig.Name)
					Expect(err).NotTo(HaveOccurred())

					actualBuild, err = pipelineDB.CreateJobBuild(jobOneTwoConfig.Name)
					Expect(err).NotTo(HaveOccurred())

					err = pipelineDB.SaveNextInputMapping(nil, "other-serial-group-job")
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the next most pending build in a group of jobs", func() {
					build, found, err := pipelineDB.GetNextPendingBuildBySerialGroup(jobOneConfig.Name, []string{"serial-group"})
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(build.ID()).To(Equal(actualBuild.ID()))
				})
			})

			It("should return the next most pending build in a group of jobs", func() {
				buildOne, err := pipelineDB.CreateJobBuild(jobOneConfig.Name)
				Expect(err).NotTo(HaveOccurred())

				buildTwo, err := pipelineDB.CreateJobBuild(jobOneConfig.Name)
				Expect(err).NotTo(HaveOccurred())

				buildThree, err := pipelineDB.CreateJobBuild(jobOneTwoConfig.Name)
				Expect(err).NotTo(HaveOccurred())

				err = pipelineDB.SaveNextInputMapping(nil, "some-job")
				Expect(err).NotTo(HaveOccurred())
				err = pipelineDB.SaveNextInputMapping(nil, "other-serial-group-job")
				Expect(err).NotTo(HaveOccurred())

				build, found, err := pipelineDB.GetNextPendingBuildBySerialGroup(jobOneConfig.Name, []string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildOne.ID()))

				build, found, err = pipelineDB.GetNextPendingBuildBySerialGroup(jobOneTwoConfig.Name, []string{"serial-group", "really-different-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildOne.ID()))

				Expect(buildOne.Finish(db.StatusSucceeded)).To(Succeed())

				build, found, err = pipelineDB.GetNextPendingBuildBySerialGroup(jobOneConfig.Name, []string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildTwo.ID()))

				build, found, err = pipelineDB.GetNextPendingBuildBySerialGroup(jobOneTwoConfig.Name, []string{"serial-group", "really-different-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildTwo.ID()))

				scheduled, err := pipelineDB.UpdateBuildToScheduled(buildTwo.ID())
				Expect(err).NotTo(HaveOccurred())
				Expect(scheduled).To(BeTrue())
				Expect(buildTwo.Finish(db.StatusSucceeded)).To(Succeed())

				build, found, err = pipelineDB.GetNextPendingBuildBySerialGroup(jobOneConfig.Name, []string{"serial-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildThree.ID()))

				build, found, err = pipelineDB.GetNextPendingBuildBySerialGroup(jobOneTwoConfig.Name, []string{"serial-group", "really-different-group"})
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(build.ID()).To(Equal(buildThree.ID()))
			})
		})

		Describe("GetRunningBuildsBySerialGroup", func() {
			Describe("same job", func() {
				var startedBuild, scheduledBuild db.Build

				BeforeEach(func() {
					var err error
					_, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					startedBuild, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())
					_, err = startedBuild.Start("", "")
					Expect(err).NotTo(HaveOccurred())

					scheduledBuild, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					scheduled, err := pipelineDB.UpdateBuildToScheduled(scheduledBuild.ID())
					Expect(err).NotTo(HaveOccurred())
					Expect(scheduled).To(BeTrue())

					for _, s := range []db.Status{db.StatusSucceeded, db.StatusFailed, db.StatusErrored, db.StatusAborted} {
						finishedBuild, err := pipelineDB.CreateJobBuild("some-job")
						Expect(err).NotTo(HaveOccurred())

						scheduled, err = pipelineDB.UpdateBuildToScheduled(finishedBuild.ID())
						Expect(err).NotTo(HaveOccurred())
						Expect(scheduled).To(BeTrue())
						err = finishedBuild.Finish(s)
					}

					_, err = pipelineDB.CreateJobBuild("some-other-job")
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns a list of running or schedule builds for said job", func() {
					builds, err := pipelineDB.GetRunningBuildsBySerialGroup("some-job", []string{"serial-group"})
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
					_, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					serialGroupBuild, err = pipelineDB.CreateJobBuild("other-serial-group-job")
					Expect(err).NotTo(HaveOccurred())

					scheduled, err := pipelineDB.UpdateBuildToScheduled(serialGroupBuild.ID())
					Expect(err).NotTo(HaveOccurred())
					Expect(scheduled).To(BeTrue())

					differentSerialGroupBuild, err := pipelineDB.CreateJobBuild("different-serial-group-job")
					Expect(err).NotTo(HaveOccurred())

					scheduled, err = pipelineDB.UpdateBuildToScheduled(differentSerialGroupBuild.ID())
					Expect(err).NotTo(HaveOccurred())
					Expect(scheduled).To(BeTrue())
				})

				It("returns a list of builds in the same serial group", func() {
					builds, err := pipelineDB.GetRunningBuildsBySerialGroup("some-job", []string{"serial-group"})
					Expect(err).NotTo(HaveOccurred())

					Expect(len(builds)).To(Equal(1))
					Expect(builds[0].ID()).To(Equal(serialGroupBuild.ID()))
				})
			})
		})

		Context("when a build is created for a job", func() {
			var build1DB db.Build
			var jobConfig atc.JobConfig

			BeforeEach(func() {
				var err error
				build1DB, err = pipelineDB.CreateJobBuild("some-job")

				jobConfig = atc.JobConfig{
					Name:   "some-job",
					Serial: false,
				}
				Expect(err).NotTo(HaveOccurred())

				Expect(build1DB.ID()).NotTo(BeZero())
				Expect(build1DB.JobName()).To(Equal("some-job"))
				Expect(build1DB.Name()).To(Equal("1"))
				Expect(build1DB.Status()).To(Equal(db.StatusPending))
				Expect(build1DB.IsScheduled()).To(BeFalse())
			})

			It("becomes the next pending build", func() {
				nextPending, found, err := pipelineDB.GetNextPendingBuild("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(nextPending).To(Equal(build1DB))
			})

			It("is returned in the job's builds", func() {
				Expect(pipelineDB.GetAllJobBuilds("some-job")).To(ConsistOf([]db.Build{build1DB}))
			})

			Context("and another build for a different pipeline is created with the same job name", func() {
				BeforeEach(func() {
					otherBuild, err := otherPipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					Expect(otherBuild.ID()).NotTo(BeZero())
					Expect(otherBuild.JobName()).To(Equal("some-job"))
					Expect(otherBuild.Name()).To(Equal("1"))
					Expect(otherBuild.Status()).To(Equal(db.StatusPending))
					Expect(otherBuild.IsScheduled()).To(BeFalse())
				})

				It("does not change the next pending build", func() {
					nextPending, found, err := pipelineDB.GetNextPendingBuild("some-job")
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(nextPending).To(Equal(build1DB))
				})

				It("is not returned in the job's builds", func() {
					Expect(pipelineDB.GetAllJobBuilds("some-job")).To(ConsistOf([]db.Build{build1DB}))
				})
			})

			Context("when scheduled", func() {
				BeforeEach(func() {
					scheduled, err := pipelineDB.UpdateBuildToScheduled(build1DB.ID())
					Expect(err).NotTo(HaveOccurred())
					Expect(scheduled).To(BeTrue())
					build1DB.Reload()
				})

				It("remains the next pending build", func() {
					nextPending, found, err := pipelineDB.GetNextPendingBuild("some-job")
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(nextPending).To(Equal(build1DB))
				})
			})

			Context("when started", func() {
				BeforeEach(func() {
					started, err := build1DB.Start("some-engine", "some-metadata")
					Expect(err).NotTo(HaveOccurred())
					Expect(started).To(BeTrue())
				})

				It("saves the updated status, and the engine and engine metadata", func() {
					found, err := build1DB.Reload()
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(build1DB.Status()).To(Equal(db.StatusStarted))
					Expect(build1DB.Engine()).To(Equal("some-engine"))
					Expect(build1DB.EngineMetadata()).To(Equal("some-metadata"))
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
					err := build1DB.Finish(db.StatusSucceeded)
					Expect(err).NotTo(HaveOccurred())
				})

				It("sets the build's status and end time", func() {
					found, err := build1DB.Reload()
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(build1DB.Status()).To(Equal(db.StatusSucceeded))
					Expect(build1DB.EndTime().Unix()).To(BeNumerically("~", time.Now().Unix(), 3))
				})
			})

			Context("and another is created for the same job", func() {
				var build2DB db.Build

				BeforeEach(func() {
					var err error
					build2DB, err = pipelineDB.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())

					Expect(build2DB.ID()).NotTo(BeZero())
					Expect(build2DB.ID()).NotTo(Equal(build1DB.ID()))
					Expect(build2DB.Name()).To(Equal("2"))
					Expect(build2DB.Status()).To(Equal(db.StatusPending))
				})

				It("is returned in the job's builds, before the rest", func() {
					Expect(pipelineDB.GetAllJobBuilds("some-job")).To(Equal([]db.Build{build2DB, build1DB}))
				})

				Describe("the first build", func() {
					It("remains the next pending build", func() {
						nextPending, found, err := pipelineDB.GetNextPendingBuild("some-job")
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeTrue())
						Expect(nextPending.ID()).To(Equal(build1DB.ID()))
					})
				})
			})

			Context("and another is created for a different job", func() {
				var otherJobBuild db.Build

				BeforeEach(func() {
					var err error

					otherJobBuild, err = pipelineDB.CreateJobBuild("some-other-job")
					Expect(err).NotTo(HaveOccurred())

					Expect(otherJobBuild.ID()).NotTo(BeZero())
					Expect(otherJobBuild.Name()).To(Equal("1"))
					Expect(otherJobBuild.Status()).To(Equal(db.StatusPending))
				})

				It("shows up in its job's builds", func() {
					Expect(pipelineDB.GetAllJobBuilds("some-other-job")).To(Equal([]db.Build{otherJobBuild}))
				})

				It("does not show up in the first build's job's builds", func() {
					Expect(pipelineDB.GetAllJobBuilds("some-job")).To(Equal([]db.Build{build1DB}))
				})
			})
		})

		It("can report a job's latest running and finished builds", func() {
			finished, next, err := pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next).To(BeNil())
			Expect(finished).To(BeNil())

			finishedBuild, err := pipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			err = finishedBuild.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			otherFinishedBuild, err := otherPipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			err = otherFinishedBuild.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next).To(BeNil())
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			nextBuild, err := pipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			started, err := nextBuild.Start("some-engine", "meta")
			Expect(err).NotTo(HaveOccurred())
			Expect(started).To(BeTrue())

			otherNextBuild, err := otherPipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			otherStarted, err := otherNextBuild.Start("some-engine", "meta")
			Expect(err).NotTo(HaveOccurred())
			Expect(otherStarted).To(BeTrue())

			finished, next, err = pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID()))
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			anotherRunningBuild, err := pipelineDB.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID())) // not anotherRunningBuild
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			started, err = anotherRunningBuild.Start("some-engine", "meta")
			Expect(err).NotTo(HaveOccurred())
			Expect(started).To(BeTrue())

			finished, next, err = pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(nextBuild.ID())) // not anotherRunningBuild
			Expect(finished.ID()).To(Equal(finishedBuild.ID()))

			err = nextBuild.Finish(db.StatusSucceeded)
			Expect(err).NotTo(HaveOccurred())

			finished, next, err = pipelineDB.GetJobFinishedAndNextBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			Expect(next.ID()).To(Equal(anotherRunningBuild.ID()))
			Expect(finished.ID()).To(Equal(nextBuild.ID()))
		})
	})
})
