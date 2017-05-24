package dbng_test

import (
	"encoding/json"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Team", func() {
	var (
		team      dbng.Team
		otherTeam dbng.Team
	)

	BeforeEach(func() {
		var err error
		team, err = teamFactory.CreateTeam(atc.Team{Name: "some-team"})
		Expect(err).ToNot(HaveOccurred())
		otherTeam, err = teamFactory.CreateTeam(atc.Team{Name: "some-other-team"})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			team, found, err := teamFactory.FindTeam("some-other-team")
			Expect(team.Name()).To(Equal("some-other-team"))
			Expect(found).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			err = otherTeam.Delete()
			Expect(err).ToNot(HaveOccurred())
		})

		It("deletes the team", func() {
			team, found, err := teamFactory.FindTeam("some-other-team")
			Expect(team).To(BeNil())
			Expect(found).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("SaveWorker", func() {
		var (
			team      dbng.Team
			otherTeam dbng.Team
			atcWorker atc.Worker
			err       error
		)

		BeforeEach(func() {
			postgresRunner.Truncate()
			team, err = teamFactory.CreateTeam(atc.Team{Name: "team"})
			Expect(err).NotTo(HaveOccurred())

			otherTeam, err = teamFactory.CreateTeam(atc.Team{Name: "some-other-team"})
			Expect(err).NotTo(HaveOccurred())
			atcWorker = atc.Worker{
				GardenAddr:       "some-garden-addr",
				BaggageclaimURL:  "some-bc-url",
				HTTPProxyURL:     "some-http-proxy-url",
				HTTPSProxyURL:    "some-https-proxy-url",
				NoProxy:          "some-no-proxy",
				ActiveContainers: 140,
				ResourceTypes: []atc.WorkerResourceType{
					{
						Type:    "some-resource-type",
						Image:   "some-image",
						Version: "some-version",
					},
					{
						Type:    "other-resource-type",
						Image:   "other-image",
						Version: "other-version",
					},
				},
				Platform:  "some-platform",
				Tags:      atc.Tags{"some", "tags"},
				Name:      "some-name",
				StartTime: 55,
			}
		})

		Context("the worker already exists", func() {
			Context("the worker is not in stalled state", func() {
				Context("the team_id of the new worker is the same", func() {
					BeforeEach(func() {
						_, err := team.SaveWorker(atcWorker, 5*time.Minute)
						Expect(err).NotTo(HaveOccurred())
					})
					It("overwrites all the data", func() {
						atcWorker.GardenAddr = "new-garden-addr"
						savedWorker, err := team.SaveWorker(atcWorker, 5*time.Minute)
						Expect(err).NotTo(HaveOccurred())
						Expect(savedWorker.Name()).To(Equal("some-name"))
						Expect(*savedWorker.GardenAddr()).To(Equal("new-garden-addr"))
						Expect(savedWorker.State()).To(Equal(dbng.WorkerStateRunning))
					})
				})
				Context("the team_id of the new worker is different", func() {
					BeforeEach(func() {
						_, err = otherTeam.SaveWorker(atcWorker, 5*time.Minute)
						Expect(err).NotTo(HaveOccurred())
					})
					It("errors", func() {
						_, err = team.SaveWorker(atcWorker, 5*time.Minute)
						Expect(err).To(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("Workers", func() {
		var (
			team      dbng.Team
			otherTeam dbng.Team
			atcWorker atc.Worker
			err       error
		)

		BeforeEach(func() {
			postgresRunner.Truncate()
			team, err = teamFactory.CreateTeam(atc.Team{Name: "team"})
			Expect(err).NotTo(HaveOccurred())

			otherTeam, err = teamFactory.CreateTeam(atc.Team{Name: "some-other-team"})
			Expect(err).NotTo(HaveOccurred())
			atcWorker = atc.Worker{
				GardenAddr:       "some-garden-addr",
				BaggageclaimURL:  "some-bc-url",
				HTTPProxyURL:     "some-http-proxy-url",
				HTTPSProxyURL:    "some-https-proxy-url",
				NoProxy:          "some-no-proxy",
				ActiveContainers: 140,
				ResourceTypes: []atc.WorkerResourceType{
					{
						Type:    "some-resource-type",
						Image:   "some-image",
						Version: "some-version",
					},
					{
						Type:    "other-resource-type",
						Image:   "other-image",
						Version: "other-version",
					},
				},
				Platform:  "some-platform",
				Tags:      atc.Tags{"some", "tags"},
				Name:      "some-name",
				StartTime: 55,
			}
		})

		Context("when there are global workers and workers for the team", func() {
			BeforeEach(func() {
				_, err = team.SaveWorker(atcWorker, 0)
				Expect(err).NotTo(HaveOccurred())

				atcWorker.Name = "some-new-worker"
				atcWorker.GardenAddr = "some-other-garden-addr"
				atcWorker.BaggageclaimURL = "some-other-bc-url"
				_, err = workerFactory.SaveWorker(atcWorker, 0)
				Expect(err).NotTo(HaveOccurred())
			})

			It("finds them without error", func() {
				workers, err := team.Workers()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(workers)).To(Equal(2))

				Expect(workers[0].Name()).To(Equal("some-name"))
				Expect(*workers[0].GardenAddr()).To(Equal("some-garden-addr"))
				Expect(*workers[0].BaggageclaimURL()).To(Equal("some-bc-url"))

				Expect(workers[1].Name()).To(Equal("some-new-worker"))
				Expect(*workers[1].GardenAddr()).To(Equal("some-other-garden-addr"))
				Expect(*workers[1].BaggageclaimURL()).To(Equal("some-other-bc-url"))
			})
		})

		Context("when there are workers for another team", func() {
			BeforeEach(func() {
				atcWorker.Name = "some-other-team-worker"
				atcWorker.GardenAddr = "some-other-garden-addr"
				atcWorker.BaggageclaimURL = "some-other-bc-url"
				_, err = otherTeam.SaveWorker(atcWorker, 5*time.Minute)
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not find the other team workers", func() {
				workers, err := team.Workers()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(workers)).To(Equal(0))
			})
		})

		Context("when there are no workers", func() {
			It("returns an error", func() {
				workers, err := workerFactory.Workers()
				Expect(err).NotTo(HaveOccurred())
				Expect(workers).To(BeEmpty())
			})
		})
	})

	Describe("FindContainersByMetadata", func() {
		var sampleMetadata []dbng.ContainerMetadata
		var metaContainers map[dbng.ContainerMetadata][]dbng.Container

		BeforeEach(func() {
			baseMetadata := fullMetadata

			diffType := fullMetadata
			diffType.Type = dbng.ContainerTypeCheck

			diffStepName := fullMetadata
			diffStepName.StepName = fullMetadata.StepName + "-other"

			diffAttempt := fullMetadata
			diffAttempt.Attempt = fullMetadata.Attempt + ",2"

			diffPipelineID := fullMetadata
			diffPipelineID.PipelineID = fullMetadata.PipelineID + 1

			diffJobID := fullMetadata
			diffJobID.JobID = fullMetadata.JobID + 1

			diffBuildID := fullMetadata
			diffBuildID.BuildID = fullMetadata.BuildID + 1

			diffWorkingDirectory := fullMetadata
			diffWorkingDirectory.WorkingDirectory = fullMetadata.WorkingDirectory + "/other"

			diffUser := fullMetadata
			diffUser.User = fullMetadata.User + "-other"

			sampleMetadata = []dbng.ContainerMetadata{
				baseMetadata,
				diffType,
				diffStepName,
				diffAttempt,
				diffPipelineID,
				diffJobID,
				diffBuildID,
				diffWorkingDirectory,
				diffUser,
			}

			build, err := defaultPipeline.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			metaContainers = make(map[dbng.ContainerMetadata][]dbng.Container)
			for _, meta := range sampleMetadata {
				firstContainerCreating, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), build.ID(), atc.PlanID("some-job"), meta)
				Expect(err).NotTo(HaveOccurred())

				metaContainers[meta] = append(metaContainers[meta], firstContainerCreating)

				secondContainerCreating, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), build.ID(), atc.PlanID("some-job"), meta)
				Expect(err).NotTo(HaveOccurred())

				secondContainerCreated, err := secondContainerCreating.Created()
				Expect(err).NotTo(HaveOccurred())

				metaContainers[meta] = append(metaContainers[meta], secondContainerCreated)

				thirdContainerCreating, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), build.ID(), atc.PlanID("some-job"), meta)
				Expect(err).NotTo(HaveOccurred())

				thirdContainerCreated, err := thirdContainerCreating.Created()
				Expect(err).NotTo(HaveOccurred())

				// third container is not appended; we don't want Destroying containers
				thirdContainerDestroying, err := thirdContainerCreated.Destroying()
				Expect(err).NotTo(HaveOccurred())

				metaContainers[meta] = append(metaContainers[meta], thirdContainerDestroying)
			}
		})

		It("finds creating, created, and destroying containers for the team, matching the metadata in full", func() {
			for _, meta := range sampleMetadata {
				expectedHandles := []string{}
				for _, c := range metaContainers[meta] {
					expectedHandles = append(expectedHandles, c.Handle())
				}

				containers, err := defaultTeam.FindContainersByMetadata(meta)
				Expect(err).ToNot(HaveOccurred())

				foundHandles := []string{}
				for _, c := range containers {
					foundHandles = append(foundHandles, c.Handle())
				}

				// should always find a Creating container and a Created container
				Expect(foundHandles).To(HaveLen(3))
				Expect(foundHandles).To(ConsistOf(expectedHandles))
			}
		})

		It("finds containers for the team, matching partial metadata", func() {
			containers, err := defaultTeam.FindContainersByMetadata(dbng.ContainerMetadata{
				Type: dbng.ContainerTypeTask,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(containers).ToNot(BeEmpty())

			foundHandles := []string{}
			for _, c := range containers {
				foundHandles = append(foundHandles, c.Handle())
			}

			var notFound int
			for meta, cs := range metaContainers {
				if meta.Type == dbng.ContainerTypeTask {
					for _, c := range cs {
						Expect(foundHandles).To(ContainElement(c.Handle()))
					}
				} else {
					for _, c := range cs {
						Expect(foundHandles).ToNot(ContainElement(c.Handle()))
						notFound++
					}
				}
			}

			// just to assert test setup is valid
			Expect(notFound).ToNot(BeZero())
		})

		It("finds all containers for the team when given empty metadata", func() {
			containers, err := defaultTeam.FindContainersByMetadata(dbng.ContainerMetadata{})
			Expect(err).ToNot(HaveOccurred())
			Expect(containers).ToNot(BeEmpty())

			foundHandles := []string{}
			for _, c := range containers {
				foundHandles = append(foundHandles, c.Handle())
			}

			for _, cs := range metaContainers {
				for _, c := range cs {
					Expect(foundHandles).To(ContainElement(c.Handle()))
				}
			}
		})

		It("does not find containers for other teams", func() {
			for _, meta := range sampleMetadata {
				containers, err := otherTeam.FindContainersByMetadata(meta)
				Expect(err).ToNot(HaveOccurred())
				Expect(containers).To(BeEmpty())
			}
		})
	})

	Describe("FindCheckContainers", func() {
		Context("when pipeline exists", func() {
			Context("when resource exists", func() {
				Context("when check container for resource exists", func() {
					var resourceContainer dbng.CreatingContainer
					var usedResourceConfig *dbng.UsedResourceConfig

					BeforeEach(func() {
						pipelineResourceTypes, err := defaultPipeline.ResourceTypes()
						Expect(err).NotTo(HaveOccurred())

						usedResourceConfig, err = resourceConfigFactory.FindOrCreateResourceConfig(
							logger,
							dbng.ForResource(defaultResource.ID()),
							defaultResource.Type(),
							defaultResource.Source(),
							pipelineResourceTypes.Deserialize(),
						)
						Expect(err).NotTo(HaveOccurred())

						resourceContainer, err = defaultTeam.CreateResourceCheckContainer(
							"default-worker",
							usedResourceConfig,
							dbng.ContainerMetadata{},
						)
						Expect(err).NotTo(HaveOccurred())
					})

					It("returns check container for resource", func() {
						containers, err := defaultTeam.FindCheckContainers(logger, "default-pipeline", "some-resource")
						Expect(err).NotTo(HaveOccurred())
						Expect(containers).To(ContainElement(resourceContainer))
					})

					Context("when another team has a container with the same resource config", func() {
						BeforeEach(func() {
							_, err := otherTeam.CreateResourceCheckContainer(
								"default-worker",
								usedResourceConfig,
								dbng.ContainerMetadata{},
							)
							Expect(err).NotTo(HaveOccurred())
						})

						It("only returns container for current team", func() {
							containers, err := defaultTeam.FindCheckContainers(logger, "default-pipeline", "some-resource")
							Expect(err).NotTo(HaveOccurred())
							Expect(containers).To(HaveLen(1))
							Expect(containers).To(ContainElement(resourceContainer))
						})
					})
				})

				Context("when check container does not exist", func() {
					It("returns empty list", func() {
						containers, err := defaultTeam.FindCheckContainers(logger, "default-pipeline", "some-resource")
						Expect(err).NotTo(HaveOccurred())
						Expect(containers).To(BeEmpty())
					})
				})
			})

			Context("when resource does not exist", func() {
				It("returns empty list", func() {
					containers, err := defaultTeam.FindCheckContainers(logger, "default-pipeline", "non-existent-resource")
					Expect(err).NotTo(HaveOccurred())
					Expect(containers).To(BeEmpty())
				})
			})
		})

		Context("when pipeline does not exist", func() {
			It("returns empty list", func() {
				containers, err := defaultTeam.FindCheckContainers(logger, "non-existent-pipeline", "some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(containers).To(BeEmpty())
			})
		})
	})

	Describe("FindContainerByHandle", func() {
		var createdContainer dbng.CreatedContainer

		BeforeEach(func() {
			build, err := defaultPipeline.CreateJobBuild("some-job")
			Expect(err).NotTo(HaveOccurred())

			creatingContainer, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), build.ID(), atc.PlanID("some-job"), dbng.ContainerMetadata{Type: "task", StepName: "some-task"})
			Expect(err).NotTo(HaveOccurred())

			createdContainer, err = creatingContainer.Created()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when worker is no longer in database", func() {
			BeforeEach(func() {
				err := defaultWorker.Delete()
				Expect(err).NotTo(HaveOccurred())
			})

			It("the container goes away from the db", func() {
				_, found, err := defaultTeam.FindContainerByHandle(createdContainer.Handle())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})
		})

		It("finds a container for the team", func() {
			container, found, err := defaultTeam.FindContainerByHandle(createdContainer.Handle())
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(container).ToNot(BeNil())
			Expect(container.Handle()).To(Equal(createdContainer.Handle()))
		})

		It("does not find container for another team", func() {
			_, found, err := otherTeam.FindContainerByHandle(createdContainer.Handle())
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
	})

	Describe("FindWorkerForResourceCheckContainer", func() {
		var resourceConfig *dbng.UsedResourceConfig

		BeforeEach(func() {
			var err error
			resourceConfig, err = resourceConfigFactory.FindOrCreateResourceConfig(
				logger,
				dbng.ForResource(defaultResource.ID()),
				"some-base-resource-type",
				atc.Source{"some": "source"},
				atc.VersionedResourceTypes{},
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is a creating container", func() {
			BeforeEach(func() {
				_, err := defaultTeam.CreateResourceCheckContainer(defaultWorker.Name(), resourceConfig, dbng.ContainerMetadata{Type: "check"})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForResourceCheckContainer(resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForResourceCheckContainer(resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})

		Context("when there is a created container", func() {
			var originalCreatedContainer dbng.CreatedContainer

			BeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateResourceCheckContainer(defaultWorker.Name(), resourceConfig, dbng.ContainerMetadata{Type: "check"})
				Expect(err).NotTo(HaveOccurred())
				originalCreatedContainer, err = creatingContainer.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForResourceCheckContainer(resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForResourceCheckContainer(resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})

			Context("when container is expired", func() {
				BeforeEach(func() {
					_, err := psql.Update("containers").
						Set("best_if_used_by", sq.Expr("NOW() - '1 second'::INTERVAL")).
						Where(sq.Eq{"id": originalCreatedContainer.ID()}).
						RunWith(dbConn).Exec()
					Expect(err).NotTo(HaveOccurred())
				})

				It("does not find it", func() {
					worker, found, err := defaultTeam.FindWorkerForResourceCheckContainer(resourceConfig)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeFalse())
					Expect(worker).To(BeNil())
				})
			})
		})

		Context("when there is no container", func() {
			It("returns nil", func() {
				worker, found, err := defaultTeam.FindWorkerForResourceCheckContainer(resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})
	})

	Describe("FindResourceCheckContainerOnWorker", func() {
		var resourceConfig *dbng.UsedResourceConfig

		BeforeEach(func() {
			var err error
			resourceConfig, err = resourceConfigFactory.FindOrCreateResourceConfig(
				logger,
				dbng.ForResource(defaultResource.ID()),
				"some-base-resource-type",
				atc.Source{"some": "source"},
				atc.VersionedResourceTypes{},
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is a creating container", func() {
			BeforeEach(func() {
				_, err := defaultTeam.CreateResourceCheckContainer(defaultWorker.Name(), resourceConfig, dbng.ContainerMetadata{Type: "check"})
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).To(BeNil())
				Expect(creatingContainer).NotTo(BeNil())
			})

			It("does not find container for another team", func() {
				creatingContainer, createdContainer, err := otherTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingContainer).To(BeNil())
				Expect(createdContainer).To(BeNil())
			})
		})

		Context("when there is a created container", func() {
			var originalCreatedContainer dbng.CreatedContainer

			BeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateResourceCheckContainer(defaultWorker.Name(), resourceConfig, dbng.ContainerMetadata{Type: "check"})
				Expect(err).NotTo(HaveOccurred())
				originalCreatedContainer, err = creatingContainer.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).NotTo(BeNil())
				Expect(creatingContainer).To(BeNil())
			})

			It("does not find container for another team", func() {
				creatingContainer, createdContainer, err := otherTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingContainer).To(BeNil())
				Expect(createdContainer).To(BeNil())
			})

			Context("when container is expired", func() {
				BeforeEach(func() {
					_, err := psql.Update("containers").
						Set("best_if_used_by", sq.Expr("NOW() - '1 second'::INTERVAL")).
						Where(sq.Eq{"id": originalCreatedContainer.ID()}).
						RunWith(dbConn).Exec()
					Expect(err).NotTo(HaveOccurred())
				})

				It("does not find it", func() {
					creatingContainer, createdContainer, err := defaultTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
					Expect(err).NotTo(HaveOccurred())
					Expect(creatingContainer).To(BeNil())
					Expect(createdContainer).To(BeNil())
				})
			})
		})

		Context("when there is no container", func() {
			It("returns nil", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindResourceCheckContainerOnWorker(defaultWorker.Name(), resourceConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).To(BeNil())
				Expect(creatingContainer).To(BeNil())
			})
		})
	})

	Describe("FindWorkerForContainer", func() {
		var containerMetadata dbng.ContainerMetadata
		var defaultBuild dbng.Build

		BeforeEach(func() {
			var err error
			containerMetadata = dbng.ContainerMetadata{
				Type:     "task",
				StepName: "some-task",
			}
			defaultBuild, err = defaultTeam.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is a creating container", func() {
			var container dbng.CreatingContainer

			BeforeEach(func() {
				var err error
				container, err = defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForContainer(container.Handle())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForContainer(container.Handle())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})

		Context("when there is a created container", func() {
			var container dbng.CreatedContainer

			BeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).NotTo(HaveOccurred())

				container, err = creatingContainer.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForContainer(container.Handle())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForContainer(container.Handle())
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})

		Context("when there is no container", func() {
			It("returns nil", func() {
				worker, found, err := defaultTeam.FindWorkerForContainer("bogus-handle")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})
	})

	Describe("FindWorkerForBuildContainer", func() {
		var containerMetadata dbng.ContainerMetadata
		var defaultBuild dbng.Build

		BeforeEach(func() {
			containerMetadata = dbng.ContainerMetadata{
				Type:     "task",
				StepName: "some-task",
			}
			var err error
			defaultBuild, err = defaultTeam.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is a creating container", func() {
			BeforeEach(func() {
				_, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForBuildContainer(defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForBuildContainer(defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})

		Context("when there is a created container", func() {
			BeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).NotTo(HaveOccurred())
				_, err = creatingContainer.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				worker, found, err := defaultTeam.FindWorkerForBuildContainer(defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(worker).NotTo(BeNil())
				Expect(worker.Name()).To(Equal(defaultWorker.Name()))
			})

			It("does not find container for another team", func() {
				worker, found, err := otherTeam.FindWorkerForBuildContainer(defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})

		Context("when there is no container", func() {
			It("returns nil", func() {
				worker, found, err := defaultTeam.FindWorkerForBuildContainer(defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(worker).To(BeNil())
			})
		})
	})

	Describe("FindBuildContainerOnWorker", func() {
		var containerMetadata dbng.ContainerMetadata
		var defaultBuild dbng.Build

		BeforeEach(func() {
			containerMetadata = dbng.ContainerMetadata{
				Type:     "task",
				StepName: "some-task",
			}
			var err error
			defaultBuild, err = defaultTeam.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is a creating container", func() {
			BeforeEach(func() {
				_, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns it", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindBuildContainerOnWorker(defaultWorker.Name(), defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).To(BeNil())
				Expect(creatingContainer).NotTo(BeNil())
			})

			It("does not find container for another team", func() {
				creatingContainer, createdContainer, err := otherTeam.FindBuildContainerOnWorker(defaultWorker.Name(), defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingContainer).To(BeNil())
				Expect(createdContainer).To(BeNil())
			})
		})

		Context("when there is a created container", func() {
			BeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateBuildContainer(defaultWorker.Name(), defaultBuild.ID(), "some-plan", containerMetadata)
				Expect(err).NotTo(HaveOccurred())
				_, err = creatingContainer.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns it", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindBuildContainerOnWorker(defaultWorker.Name(), defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).NotTo(BeNil())
				Expect(creatingContainer).To(BeNil())
			})

			It("does not find container for another team", func() {
				creatingContainer, createdContainer, err := otherTeam.FindBuildContainerOnWorker(defaultWorker.Name(), defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingContainer).To(BeNil())
				Expect(createdContainer).To(BeNil())
			})
		})

		Context("when there is no container", func() {
			It("returns nil", func() {
				creatingContainer, createdContainer, err := defaultTeam.FindBuildContainerOnWorker(defaultWorker.Name(), defaultBuild.ID(), "some-plan")
				Expect(err).NotTo(HaveOccurred())
				Expect(createdContainer).To(BeNil())
				Expect(creatingContainer).To(BeNil())
			})
		})
	})

	Describe("Updating Auth", func() {
		var (
			basicAuth    *atc.BasicAuth
			authProvider map[string]*json.RawMessage
		)

		BeforeEach(func() {
			basicAuth = &atc.BasicAuth{
				BasicAuthUsername: "fake user",
				BasicAuthPassword: "no, bad",
			}

			data := []byte(`{"credit_card":"please"}`)
			authProvider = map[string]*json.RawMessage{
				"fake-provider": (*json.RawMessage)(&data),
			}
		})

		Describe("UpdateBasicAuth", func() {
			It("saves basic auth team info without overwriting the provider auth", func() {
				err := team.UpdateProviderAuth(authProvider)
				Expect(err).NotTo(HaveOccurred())

				err = team.UpdateBasicAuth(basicAuth)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.Auth()).To(Equal(authProvider))
			})

			It("saves basic auth team info to the existing team", func() {
				err := team.UpdateBasicAuth(basicAuth)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.BasicAuth().BasicAuthUsername).To(Equal(basicAuth.BasicAuthUsername))
				Expect(bcrypt.CompareHashAndPassword([]byte(team.BasicAuth().BasicAuthPassword),
					[]byte(basicAuth.BasicAuthPassword))).To(BeNil())
			})

			It("nulls basic auth when has a blank username", func() {
				basicAuth.BasicAuthUsername = ""
				err := team.UpdateBasicAuth(basicAuth)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.BasicAuth()).To(BeNil())
			})

			It("nulls basic auth when has a blank password", func() {
				basicAuth.BasicAuthPassword = ""
				err := team.UpdateBasicAuth(basicAuth)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.BasicAuth()).To(BeNil())
			})
		})

		Describe("UpdateProviderAuth", func() {
			It("saves auth team info to the existing team", func() {
				err := team.UpdateProviderAuth(authProvider)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.Auth()).To(Equal(authProvider))
			})

			It("saves github auth team info without over writing the basic auth", func() {
				err := team.UpdateBasicAuth(basicAuth)
				Expect(err).NotTo(HaveOccurred())

				err = team.UpdateProviderAuth(authProvider)
				Expect(err).NotTo(HaveOccurred())

				Expect(team.BasicAuth().BasicAuthUsername).To(Equal(basicAuth.BasicAuthUsername))
				Expect(bcrypt.CompareHashAndPassword([]byte(team.BasicAuth().BasicAuthPassword),
					[]byte(basicAuth.BasicAuthPassword))).To(BeNil())
			})
		})
	})

	Describe("Pipelines", func() {
		var (
			pipelines []dbng.Pipeline
			pipeline1 dbng.Pipeline
			pipeline2 dbng.Pipeline
		)

		JustBeforeEach(func() {
			var err error
			pipelines, err = team.Pipelines()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the team has configured pipelines", func() {
			BeforeEach(func() {
				var err error
				pipeline1, _, err = team.SavePipeline("fake-pipeline", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-name"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())

				pipeline2, _, err = team.SavePipeline("fake-pipeline-two", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-fake"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{pipeline1, pipeline2}))
			})
		})
		Context("when the team has no configured pipelines", func() {
			It("returns no pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{}))
			})
		})
	})

	Describe("PublicPipelines", func() {
		var (
			pipelines []dbng.Pipeline
			pipeline1 dbng.Pipeline
			pipeline2 dbng.Pipeline
		)

		JustBeforeEach(func() {
			var err error
			pipelines, err = team.PublicPipelines()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the team has configured pipelines", func() {
			BeforeEach(func() {
				var err error
				pipeline1, _, err = team.SavePipeline("fake-pipeline", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-name"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())

				pipeline2, _, err = team.SavePipeline("fake-pipeline-two", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-fake"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())

				err = pipeline2.Expose()
				Expect(err).ToNot(HaveOccurred())

				found, err := pipeline2.Reload()
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())
			})

			It("returns the pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{pipeline2}))
			})
		})
		Context("when the team has no configured pipelines", func() {
			It("returns no pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{}))
			})
		})
	})

	Describe("VisiblePipelines", func() {
		var (
			pipelines []dbng.Pipeline
			pipeline1 dbng.Pipeline
			pipeline2 dbng.Pipeline
		)

		JustBeforeEach(func() {
			var err error
			pipelines, err = team.VisiblePipelines()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the team has configured pipelines", func() {
			BeforeEach(func() {
				var err error
				pipeline1, _, err = team.SavePipeline("fake-pipeline", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-name"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())

				pipeline2, _, err = otherTeam.SavePipeline("fake-pipeline-two", atc.Config{
					Jobs: atc.JobConfigs{
						{Name: "job-fake"},
					},
				}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).ToNot(HaveOccurred())

				Expect(pipeline2.Expose()).To(Succeed())
				Expect(pipeline2.Reload()).To(BeTrue())
			})

			It("returns the pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{pipeline1, pipeline2}))
			})

			Context("when the other team has a private pipeline", func() {
				var pipeline3 dbng.Pipeline
				BeforeEach(func() {
					var err error
					pipeline3, _, err = otherTeam.SavePipeline("fake-pipeline-three", atc.Config{
						Jobs: atc.JobConfigs{
							{Name: "job-fake-again"},
						},
					}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
					Expect(err).ToNot(HaveOccurred())
				})

				It("does not return the other team private pipeline", func() {
					Expect(pipelines).To(Equal([]dbng.Pipeline{pipeline1, pipeline2}))
				})
			})
		})

		Context("when the team has no configured pipelines", func() {
			It("returns no pipelines", func() {
				Expect(pipelines).To(Equal([]dbng.Pipeline{}))
			})
		})
	})

	Describe("OrderPipelines", func() {
		var pipeline1 dbng.Pipeline
		var pipeline2 dbng.Pipeline
		var otherPipeline1 dbng.Pipeline
		var otherPipeline2 dbng.Pipeline

		BeforeEach(func() {
			var err error
			pipeline1, _, err = team.SavePipeline("pipeline-name-a", atc.Config{}, 0, dbng.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())
			pipeline2, _, err = team.SavePipeline("pipeline-name-b", atc.Config{}, 0, dbng.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())

			otherPipeline1, _, err = otherTeam.SavePipeline("pipeline-name-a", atc.Config{}, 0, dbng.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())
			otherPipeline2, _, err = otherTeam.SavePipeline("pipeline-name-b", atc.Config{}, 0, dbng.PipelineUnpaused)
			Expect(err).ToNot(HaveOccurred())
		})

		It("orders pipelines that belong to team (case insensitive)", func() {
			err := team.OrderPipelines([]string{"pipeline-name-b", "pipeline-name-a"})
			Expect(err).NotTo(HaveOccurred())

			err = otherTeam.OrderPipelines([]string{"pipeline-name-a", "pipeline-name-b"})
			Expect(err).NotTo(HaveOccurred())

			orderedPipelines, err := team.Pipelines()

			Expect(err).NotTo(HaveOccurred())
			Expect(orderedPipelines).To(HaveLen(2))
			Expect(orderedPipelines[0].ID()).To(Equal(pipeline2.ID()))
			Expect(orderedPipelines[1].ID()).To(Equal(pipeline1.ID()))

			otherTeamOrderedPipelines, err := otherTeam.Pipelines()
			Expect(err).NotTo(HaveOccurred())
			Expect(otherTeamOrderedPipelines).To(HaveLen(2))
			Expect(otherTeamOrderedPipelines[0].ID()).To(Equal(otherPipeline1.ID()))
			Expect(otherTeamOrderedPipelines[1].ID()).To(Equal(otherPipeline2.ID()))
		})
	})

	Describe("CreateOneOffBuild", func() {
		var (
			oneOffBuild dbng.Build
			err         error
		)

		BeforeEach(func() {
			oneOffBuild, err = team.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())
		})

		It("can create one-off builds", func() {
			Expect(oneOffBuild.ID()).NotTo(BeZero())
			Expect(oneOffBuild.JobName()).To(BeZero())
			Expect(oneOffBuild.PipelineName()).To(BeZero())
			Expect(oneOffBuild.Name()).To(Equal(strconv.Itoa(oneOffBuild.ID())))
			Expect(oneOffBuild.TeamName()).To(Equal(team.Name()))
			Expect(oneOffBuild.Status()).To(Equal(dbng.BuildStatusPending))
		})
	})

	Describe("PrivateAndPublicBuilds", func() {
		Context("when there are no builds", func() {
			It("returns an empty list of builds", func() {
				builds, pagination, err := team.PrivateAndPublicBuilds(dbng.Page{Limit: 2})
				Expect(err).NotTo(HaveOccurred())

				Expect(pagination.Next).To(BeNil())
				Expect(pagination.Previous).To(BeNil())
				Expect(builds).To(BeEmpty())
			})
		})

		Context("when there are builds", func() {
			var allBuilds [5]dbng.Build
			var pipeline dbng.Pipeline
			var pipelineBuilds [2]dbng.Build

			BeforeEach(func() {
				for i := 0; i < 3; i++ {
					build, err := team.CreateOneOffBuild()
					Expect(err).NotTo(HaveOccurred())
					allBuilds[i] = build
				}

				config := atc.Config{
					Jobs: atc.JobConfigs{
						{
							Name: "some-job",
						},
					},
				}
				var err error
				pipeline, _, err = team.SavePipeline("some-pipeline", config, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				for i := 3; i < 5; i++ {
					build, err := pipeline.CreateJobBuild("some-job")
					Expect(err).NotTo(HaveOccurred())
					allBuilds[i] = build
					pipelineBuilds[i-3] = build
				}
			})

			It("returns all team builds with correct pagination", func() {
				builds, pagination, err := team.PrivateAndPublicBuilds(dbng.Page{Limit: 2})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(2))
				Expect(builds[0]).To(Equal(allBuilds[4]))
				Expect(builds[1]).To(Equal(allBuilds[3]))

				Expect(pagination.Previous).To(BeNil())
				Expect(pagination.Next).To(Equal(&dbng.Page{Since: allBuilds[3].ID(), Limit: 2}))

				builds, pagination, err = team.PrivateAndPublicBuilds(*pagination.Next)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(2))

				Expect(builds[0]).To(Equal(allBuilds[2]))
				Expect(builds[1]).To(Equal(allBuilds[1]))

				Expect(pagination.Previous).To(Equal(&dbng.Page{Until: allBuilds[2].ID(), Limit: 2}))
				Expect(pagination.Next).To(Equal(&dbng.Page{Since: allBuilds[1].ID(), Limit: 2}))

				builds, pagination, err = team.PrivateAndPublicBuilds(*pagination.Next)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(1))
				Expect(builds[0]).To(Equal(allBuilds[0]))

				Expect(pagination.Previous).To(Equal(&dbng.Page{Until: allBuilds[0].ID(), Limit: 2}))
				Expect(pagination.Next).To(BeNil())

				builds, pagination, err = team.PrivateAndPublicBuilds(*pagination.Previous)
				Expect(err).NotTo(HaveOccurred())

				Expect(len(builds)).To(Equal(2))
				Expect(builds[0]).To(Equal(allBuilds[2]))
				Expect(builds[1]).To(Equal(allBuilds[1]))

				Expect(pagination.Previous).To(Equal(&dbng.Page{Until: allBuilds[2].ID(), Limit: 2}))
				Expect(pagination.Next).To(Equal(&dbng.Page{Since: allBuilds[1].ID(), Limit: 2}))
			})

			Context("when there are builds that belong to different teams", func() {
				var teamABuilds [3]dbng.Build
				var teamBBuilds [3]dbng.Build

				var caseInsensitiveTeamA dbng.Team
				var caseInsensitiveTeamB dbng.Team

				BeforeEach(func() {
					_, err := teamFactory.CreateTeam(atc.Team{Name: "team-a"})
					Expect(err).NotTo(HaveOccurred())

					_, err = teamFactory.CreateTeam(atc.Team{Name: "team-b"})
					Expect(err).NotTo(HaveOccurred())

					var found bool
					caseInsensitiveTeamA, found, err = teamFactory.FindTeam("team-A")
					Expect(found).To(BeTrue())
					Expect(err).ToNot(HaveOccurred())

					caseInsensitiveTeamB, found, err = teamFactory.FindTeam("team-B")
					Expect(found).To(BeTrue())
					Expect(err).ToNot(HaveOccurred())

					for i := 0; i < 3; i++ {
						teamABuilds[i], err = caseInsensitiveTeamA.CreateOneOffBuild()
						Expect(err).NotTo(HaveOccurred())

						teamBBuilds[i], err = caseInsensitiveTeamB.CreateOneOffBuild()
						Expect(err).NotTo(HaveOccurred())
					}
				})

				Context("when other team builds are private", func() {
					It("returns only builds for requested team", func() {
						builds, _, err := caseInsensitiveTeamA.PrivateAndPublicBuilds(dbng.Page{Limit: 10})
						Expect(err).NotTo(HaveOccurred())

						Expect(len(builds)).To(Equal(3))
						Expect(builds).To(ConsistOf(teamABuilds))

						builds, _, err = caseInsensitiveTeamB.PrivateAndPublicBuilds(dbng.Page{Limit: 10})
						Expect(err).NotTo(HaveOccurred())

						Expect(len(builds)).To(Equal(3))
						Expect(builds).To(ConsistOf(teamBBuilds))
					})
				})

				Context("when other team builds are public", func() {
					BeforeEach(func() {
						pipeline.Expose()
					})

					It("returns builds for requested team and public builds", func() {
						builds, _, err := caseInsensitiveTeamA.PrivateAndPublicBuilds(dbng.Page{Limit: 10})
						Expect(err).NotTo(HaveOccurred())

						Expect(builds).To(HaveLen(5))
						expectedBuilds := []dbng.Build{}
						for _, b := range teamABuilds {
							expectedBuilds = append(expectedBuilds, b)
						}
						for _, b := range pipelineBuilds {
							expectedBuilds = append(expectedBuilds, b)
						}
						Expect(builds).To(ConsistOf(expectedBuilds))
					})
				})
			})
		})
	})

	Describe("SavePipeline", func() {
		type SerialGroup struct {
			JobID int
			Name  string
		}
		var config atc.Config
		var otherConfig atc.Config

		BeforeEach(func() {
			config = atc.Config{
				Groups: atc.GroupConfigs{
					{
						Name:      "some-group",
						Jobs:      []string{"job-1", "job-2"},
						Resources: []string{"resource-1", "resource-2"},
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

						Serial:       true,
						SerialGroups: []string{"serial-group-1", "serial-group-2"},

						Plan: atc.PlanSequence{
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
									RootFsUri: "some-image",
								},
							},
							{
								Put: "some-resource",
								Params: atc.Params{
									"some-param": "some-value",
								},
							},
						},
					},
				},
			}

			otherConfig = atc.Config{
				Groups: atc.GroupConfigs{
					{
						Name:      "some-group",
						Jobs:      []string{"job-1", "job-2"},
						Resources: []string{"resource-1", "resource-2"},
					},
				},

				Resources: atc.ResourceConfigs{
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
						Name: "some-other-job",
					},
				},
			}
		})

		Context("on initial create", func() {
			var pipelineName string
			BeforeEach(func() {
				pipelineName = "some-pipeline"
			})

			It("returns true for created", func() {
				_, created, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())
				Expect(created).To(BeTrue())
			})

			It("caches the team id", func() {
				_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline.TeamID()).To(Equal(team.ID()))
			})

			It("can be saved as paused", func() {
				_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelinePaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(pipeline.Paused()).To(BeTrue())
			})

			It("can be saved as unpaused", func() {
				_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(pipeline.Paused()).To(BeFalse())
			})

			It("defaults to paused", func() {
				_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(pipeline.Paused()).To(BeTrue())
			})

			It("creates all of the resources from the pipeline in the database", func() {
				savedPipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				resource, found, err := savedPipeline.Resource("some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(resource.Type()).To(Equal("some-type"))
				Expect(resource.Source()).To(Equal(atc.Source{
					"source-config": "some-value",
				}))
			})

			It("updates resource config", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.Resources[0].Source = atc.Source{
					"source-other-config": "some-other-value",
				}

				savedPipeline, _, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				resource, found, err := savedPipeline.Resource("some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(resource.Type()).To(Equal("some-type"))
				Expect(resource.Source()).To(Equal(atc.Source{
					"source-other-config": "some-other-value",
				}))
			})

			It("marks resource as inactive if it is no longer in config", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.Resources = []atc.ResourceConfig{}

				savedPipeline, _, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				_, found, err := savedPipeline.Resource("some-resource")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})

			It("creates all of the resource types from the pipeline in the database", func() {
				savedPipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				resourceType, found, err := savedPipeline.ResourceType("some-resource-type")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(resourceType.Type()).To(Equal("some-type"))
				Expect(resourceType.Source()).To(Equal(atc.Source{
					"source-config": "some-value",
				}))
			})

			It("updates resource type config from the pipeline in the database", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.ResourceTypes[0].Source = atc.Source{
					"source-other-config": "some-other-value",
				}

				savedPipeline, _, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				resourceType, found, err := savedPipeline.ResourceType("some-resource-type")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(resourceType.Type()).To(Equal("some-type"))
				Expect(resourceType.Source()).To(Equal(atc.Source{
					"source-other-config": "some-other-value",
				}))
			})

			It("marks resource type as inactive if it is no longer in config", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.ResourceTypes = []atc.ResourceType{}

				savedPipeline, _, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				_, found, err := savedPipeline.ResourceType("some-resource-type")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})

			It("creates all of the jobs from the pipeline in the database", func() {
				savedPipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				job, found, err := savedPipeline.Job("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(job.Config()).To(Equal(config.Jobs[0]))
			})

			It("updates job config", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.Jobs[0].Public = false

				_, _, err = team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				job, found, err := pipeline.Job("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(job.Config().Public).To(BeFalse())
			})

			It("marks job inactive", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				config.Jobs = []atc.JobConfig{}

				savedPipeline, _, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				_, found, err := savedPipeline.Job("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			})

			It("creates all of the serial groups from the jobs in the database", func() {
				savedPipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				serialGroups := []SerialGroup{}
				rows, err := dbConn.Query("SELECT job_id, serial_group FROM jobs_serial_groups")
				Expect(err).NotTo(HaveOccurred())

				for rows.Next() {
					var serialGroup SerialGroup
					err = rows.Scan(&serialGroup.JobID, &serialGroup.Name)
					Expect(err).NotTo(HaveOccurred())
					serialGroups = append(serialGroups, serialGroup)
				}

				job, found, err := savedPipeline.Job("some-job")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(serialGroups).To(ConsistOf([]SerialGroup{
					{
						JobID: job.ID(),
						Name:  "serial-group-1",
					},
					{
						JobID: job.ID(),
						Name:  "serial-group-2",
					},
				}))
			})
		})

		Context("on updates", func() {
			var pipelineName string

			BeforeEach(func() {
				pipelineName = "a-pipeline-name"
			})

			It("it returns created as false", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				_, created, err := team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())
				Expect(created).To(BeFalse())
			})

			It("updating from paused to unpaused", func() {
				pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelinePaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline.Paused()).To(BeTrue())

				_, _, configVersion, err := pipeline.Config()
				Expect(err).NotTo(HaveOccurred())

				_, _, err = team.SavePipeline(pipelineName, config, configVersion, dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err = team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline.Paused()).To(BeFalse())
			})

			It("updating from unpaused to paused", func() {
				_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err := team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline.Paused()).To(BeFalse())

				_, _, err = team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelinePaused)
				Expect(err).NotTo(HaveOccurred())

				pipeline, found, err = team.Pipeline(pipelineName)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(pipeline.Paused()).To(BeTrue())
			})

			Context("updating with no change", func() {
				It("maintains paused if the pipeline is paused", func() {
					_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelinePaused)
					Expect(err).NotTo(HaveOccurred())

					pipeline, found, err := team.Pipeline(pipelineName)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(pipeline.Paused()).To(BeTrue())

					_, _, err = team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
					Expect(err).NotTo(HaveOccurred())

					pipeline, found, err = team.Pipeline(pipelineName)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(pipeline.Paused()).To(BeTrue())
				})

				It("maintains unpaused if the pipeline is unpaused", func() {
					_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineUnpaused)
					Expect(err).NotTo(HaveOccurred())

					pipeline, found, err := team.Pipeline(pipelineName)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(pipeline.Paused()).To(BeFalse())

					_, _, err = team.SavePipeline(pipelineName, config, pipeline.ConfigVersion(), dbng.PipelineNoChange)
					Expect(err).NotTo(HaveOccurred())

					pipeline, found, err = team.Pipeline(pipelineName)
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())
					Expect(pipeline.Paused()).To(BeFalse())
				})
			})
		})

		It("can lookup a pipeline by name", func() {
			pipelineName := "a-pipeline-name"
			otherPipelineName := "an-other-pipeline-name"

			_, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())
			_, _, err = team.SavePipeline(otherPipelineName, otherConfig, 0, dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			pipeline, found, err := team.Pipeline(pipelineName)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(pipeline.Name()).To(Equal(pipelineName))
			Expect(pipeline.ID()).NotTo(Equal(0))
			configPipeline, _, _, err := pipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(configPipeline).To(Equal(config))

			otherPipeline, found, err := team.Pipeline(otherPipelineName)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(otherPipeline.Name()).To(Equal(otherPipelineName))
			Expect(otherPipeline.ID()).NotTo(Equal(0))
			configPipeline, _, _, err = otherPipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(configPipeline).To(Equal(otherConfig))
		})

		It("can manage multiple pipeline configurations", func() {
			pipelineName := "a-pipeline-name"
			otherPipelineName := "an-other-pipeline-name"

			By("being able to save the config")
			pipeline, _, err := team.SavePipeline(pipelineName, config, 0, dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			otherPipeline, _, err := team.SavePipeline(otherPipelineName, otherConfig, 0, dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			By("returning the saved config to later gets")
			returnedConfig, returnedRawConfig, configVersion, err := pipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedConfig).To(Equal(config))
			jsonBytes, err := json.Marshal(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedRawConfig).To(MatchJSON(jsonBytes))
			Expect(configVersion).NotTo(Equal(dbng.ConfigVersion(0)))

			otherReturnedConfig, otherReturnedRawConfig, otherConfigVersion, err := otherPipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(otherReturnedConfig).To(Equal(otherConfig))
			jsonBytes, err = json.Marshal(otherConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(otherReturnedRawConfig).To(MatchJSON(jsonBytes))
			Expect(otherConfigVersion).NotTo(Equal(dbng.ConfigVersion(0)))

			updatedConfig := config

			updatedConfig.Groups = append(config.Groups, atc.GroupConfig{
				Name: "new-group",
				Jobs: []string{"new-job-1", "new-job-2"},
			})

			updatedConfig.Resources = append(config.Resources, atc.ResourceConfig{
				Name: "new-resource",
				Type: "new-type",
				Source: atc.Source{
					"new-source-config": "new-value",
				},
			})

			updatedConfig.Jobs = append(config.Jobs, atc.JobConfig{
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

			By("not allowing non-sequential updates")
			_, _, err = team.SavePipeline(pipelineName, updatedConfig, pipeline.ConfigVersion()-1, dbng.PipelineUnpaused)
			Expect(err).To(Equal(dbng.ErrConfigComparisonFailed))

			_, _, err = team.SavePipeline(pipelineName, updatedConfig, pipeline.ConfigVersion()+10, dbng.PipelineUnpaused)
			Expect(err).To(Equal(dbng.ErrConfigComparisonFailed))

			_, _, err = team.SavePipeline(otherPipelineName, updatedConfig, otherPipeline.ConfigVersion()-1, dbng.PipelineUnpaused)
			Expect(err).To(Equal(dbng.ErrConfigComparisonFailed))

			_, _, err = team.SavePipeline(otherPipelineName, updatedConfig, otherPipeline.ConfigVersion()+10, dbng.PipelineUnpaused)
			Expect(err).To(Equal(dbng.ErrConfigComparisonFailed))

			By("being able to update the config with a valid con")
			pipeline, _, err = team.SavePipeline(pipelineName, updatedConfig, pipeline.ConfigVersion(), dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())
			otherPipeline, _, err = team.SavePipeline(otherPipelineName, updatedConfig, otherPipeline.ConfigVersion(), dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			By("returning the updated config")
			returnedConfig, returnedRawConfig, newConfigVersion, err := pipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedConfig).To(Equal(updatedConfig))
			rawConfigJSONBytes, err := json.Marshal(updatedConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedRawConfig).To(MatchJSON(rawConfigJSONBytes))
			Expect(newConfigVersion).NotTo(Equal(configVersion))

			otherReturnedConfig, _, newOtherConfigVersion, err := otherPipeline.Config()
			Expect(err).NotTo(HaveOccurred())
			Expect(otherReturnedConfig).To(Equal(updatedConfig))
			Expect(returnedRawConfig).To(MatchJSON(rawConfigJSONBytes))
			Expect(newOtherConfigVersion).NotTo(Equal(otherConfigVersion))

			By("being able to retrieve invalid config")
			invalidPipelineName := "invalid-config"
			invalidPipeline, _, err := team.SavePipeline(invalidPipelineName, config, 1, dbng.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			plaintext := []byte("bad-json")
			invalidConfig, invalidNonce, err := key.Encrypt(plaintext)
			Expect(err).NotTo(HaveOccurred())

			dbConn.Exec(`
		UPDATE pipelines
		SET config = $1, nonce = $2
		WHERE name = 'invalid-config'
		`, invalidConfig, invalidNonce)

			_, _, invalidConfigVersion, err := invalidPipeline.Config()
			Expect(err).To(BeAssignableToTypeOf(atc.MalformedConfigError{}))
			Expect(err.Error()).To(ContainSubstring("malformed config:"))
			Expect(invalidConfigVersion).NotTo(Equal(dbng.ConfigVersion(1)))
		})

		Context("when there are multiple teams", func() {
			It("can allow pipelines with the same name across teams", func() {
				teamPipeline, _, err := team.SavePipeline("steve", config, 0, dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				By("allowing you to save a pipeline with the same name in another team")
				otherTeamPipeline, _, err := otherTeam.SavePipeline("steve", otherConfig, 0, dbng.PipelineUnpaused)
				Expect(err).NotTo(HaveOccurred())

				By("updating the pipeline config for the correct team's pipeline")
				teamPipeline, _, err = team.SavePipeline("steve", otherConfig, teamPipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				_, _, err = otherTeam.SavePipeline("steve", config, otherTeamPipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).NotTo(HaveOccurred())

				By("pausing the correct team's pipeline")
				_, _, err = team.SavePipeline("steve", otherConfig, teamPipeline.ConfigVersion(), dbng.PipelinePaused)
				Expect(err).NotTo(HaveOccurred())

				pausedPipeline, found, err := team.Pipeline("steve")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				unpausedPipeline, found, err := otherTeam.Pipeline("steve")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(pausedPipeline.Paused()).To(BeTrue())
				Expect(unpausedPipeline.Paused()).To(BeFalse())

				By("cannot cross update configs")
				_, _, err = team.SavePipeline("steve", otherConfig, otherTeamPipeline.ConfigVersion(), dbng.PipelineNoChange)
				Expect(err).To(HaveOccurred())

				_, _, err = team.SavePipeline("steve", otherConfig, otherTeamPipeline.ConfigVersion(), dbng.PipelinePaused)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
