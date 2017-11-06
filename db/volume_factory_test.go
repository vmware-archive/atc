package db_test

import (
	"time"

	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeFactory", func() {
	var (
		team2             db.Team
		usedResourceCache *db.UsedResourceCache
		build             db.Build
	)

	BeforeEach(func() {
		var err error
		build, err = defaultTeam.CreateOneOffBuild()
		Expect(err).ToNot(HaveOccurred())

		usedResourceCache, err = resourceCacheFactory.FindOrCreateResourceCache(
			logger,
			db.ForBuild(build.ID()),
			"some-type",
			atc.Version{"some": "version"},
			atc.Source{
				"some": "source",
			},
			atc.Params{"some": "params"},
			creds.NewVersionedResourceTypes(
				template.StaticVariables{},
				atc.VersionedResourceTypes{
					atc.VersionedResourceType{
						ResourceType: atc.ResourceType{
							Name: "some-type",
							Type: "some-base-resource-type",
							Source: atc.Source{
								"some-type": "source",
							},
						},
						Version: atc.Version{"some-type": "version"},
					},
				},
			),
		)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("GetTeamVolumes", func() {
		var (
			team1handles []string
			team2handles []string
		)

		It("returns task cache volumes", func() {
			taskCache, err := workerTaskCacheFactory.FindOrCreate(defaultJob.ID(), "some-step", "some-path", defaultWorker.Name())
			Expect(err).NotTo(HaveOccurred())

			creatingVolume, err := volumeFactory.CreateTaskCacheVolume(defaultTeam.ID(), taskCache)
			Expect(err).NotTo(HaveOccurred())

			createdVolume, err := creatingVolume.Created()
			Expect(err).NotTo(HaveOccurred())

			volumes, err := volumeFactory.GetTeamVolumes(defaultTeam.ID())
			Expect(err).NotTo(HaveOccurred())

			Expect(volumes).To(HaveLen(1))
			Expect(volumes[0].Handle()).To(Equal(createdVolume.Handle()))
			Expect(volumes[0].Type()).To(Equal(db.VolumeTypeTaskCache))
		})

		Context("with container volumes", func() {
			JustBeforeEach(func() {
				creatingContainer, err := defaultTeam.CreateContainer(defaultWorker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
					Type:     "task",
					StepName: "some-task",
				})
				Expect(err).ToNot(HaveOccurred())

				team1handles = []string{}
				team2handles = []string{}

				team2, err = teamFactory.CreateTeam(atc.Team{Name: "some-other-defaultTeam"})
				Expect(err).ToNot(HaveOccurred())

				creatingVolume1, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-1")
				Expect(err).NotTo(HaveOccurred())
				createdVolume1, err := creatingVolume1.Created()
				Expect(err).NotTo(HaveOccurred())
				team1handles = append(team1handles, createdVolume1.Handle())

				creatingVolume2, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-2")
				Expect(err).NotTo(HaveOccurred())
				createdVolume2, err := creatingVolume2.Created()
				Expect(err).NotTo(HaveOccurred())
				team1handles = append(team1handles, createdVolume2.Handle())

				creatingVolume3, err := volumeFactory.CreateContainerVolume(team2.ID(), defaultWorker.Name(), creatingContainer, "some-path-3")
				Expect(err).NotTo(HaveOccurred())
				createdVolume3, err := creatingVolume3.Created()
				Expect(err).NotTo(HaveOccurred())
				team2handles = append(team2handles, createdVolume3.Handle())
			})

			It("returns only the matching defaultTeam's volumes", func() {
				createdVolumes, err := volumeFactory.GetTeamVolumes(defaultTeam.ID())
				Expect(err).NotTo(HaveOccurred())
				createdHandles := []string{}
				for _, vol := range createdVolumes {
					createdHandles = append(createdHandles, vol.Handle())
				}
				Expect(createdHandles).To(Equal(team1handles))

				createdVolumes2, err := volumeFactory.GetTeamVolumes(team2.ID())
				Expect(err).NotTo(HaveOccurred())
				createdHandles2 := []string{}
				for _, vol := range createdVolumes2 {
					createdHandles2 = append(createdHandles2, vol.Handle())
				}
				Expect(createdHandles2).To(Equal(team2handles))
			})

			Context("when worker is stalled", func() {
				BeforeEach(func() {
					var err error
					defaultWorker, err = workerFactory.SaveWorker(defaultWorkerPayload, -10*time.Minute)
					Expect(err).NotTo(HaveOccurred())
					stalledWorkers, err := workerLifecycle.StallUnresponsiveWorkers()
					Expect(err).NotTo(HaveOccurred())
					Expect(stalledWorkers).To(ContainElement(defaultWorker.Name()))
				})

				It("returns volumes", func() {
					createdVolumes, err := volumeFactory.GetTeamVolumes(defaultTeam.ID())
					Expect(err).NotTo(HaveOccurred())
					createdHandles := []string{}
					for _, vol := range createdVolumes {
						createdHandles = append(createdHandles, vol.Handle())
					}
					Expect(createdHandles).To(Equal(team1handles))

					createdVolumes2, err := volumeFactory.GetTeamVolumes(team2.ID())
					Expect(err).NotTo(HaveOccurred())
					createdHandles2 := []string{}
					for _, vol := range createdVolumes2 {
						createdHandles2 = append(createdHandles2, vol.Handle())
					}
					Expect(createdHandles2).To(Equal(team2handles))
				})
			})
		})
	})

	Describe("GetOrphanedVolumes", func() {
		var (
			expectedCreatedHandles    []string
			expectedDestroyingHandles []string
		)

		BeforeEach(func() {
			creatingContainer, err := defaultTeam.CreateContainer(defaultWorker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
				Type:     "task",
				StepName: "some-task",
			})
			Expect(err).ToNot(HaveOccurred())

			expectedCreatedHandles = []string{}
			expectedDestroyingHandles = []string{}

			creatingVolume1, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-1")
			Expect(err).NotTo(HaveOccurred())
			createdVolume1, err := creatingVolume1.Created()
			Expect(err).NotTo(HaveOccurred())
			expectedCreatedHandles = append(expectedCreatedHandles, createdVolume1.Handle())

			creatingVolume2, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-2")
			Expect(err).NotTo(HaveOccurred())
			createdVolume2, err := creatingVolume2.Created()
			Expect(err).NotTo(HaveOccurred())
			expectedCreatedHandles = append(expectedCreatedHandles, createdVolume2.Handle())

			creatingVolume3, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-3")
			Expect(err).NotTo(HaveOccurred())
			createdVolume3, err := creatingVolume3.Created()
			Expect(err).NotTo(HaveOccurred())
			destroyingVolume3, err := createdVolume3.Destroying()
			Expect(err).NotTo(HaveOccurred())
			expectedDestroyingHandles = append(expectedDestroyingHandles, destroyingVolume3.Handle())

			resourceCacheVolume, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-4")
			Expect(err).NotTo(HaveOccurred())
			expectedCreatedHandles = append(expectedCreatedHandles, resourceCacheVolume.Handle())

			resourceCacheVolumeCreated, err := resourceCacheVolume.Created()
			Expect(err).NotTo(HaveOccurred())

			err = resourceCacheVolumeCreated.InitializeResourceCache(usedResourceCache)
			Expect(err).NotTo(HaveOccurred())

			deleted, err := build.Delete()
			Expect(err).NotTo(HaveOccurred())
			Expect(deleted).To(BeTrue())

			deleteTx, err := dbConn.Begin()
			Expect(err).ToNot(HaveOccurred())
			deleted, err = usedResourceCache.Destroy(deleteTx)
			Expect(err).NotTo(HaveOccurred())
			Expect(deleted).To(BeTrue())
			Expect(deleteTx.Commit()).To(Succeed())

			createdContainer, err := creatingContainer.Created()
			Expect(err).NotTo(HaveOccurred())
			destroyingContainer, err := createdContainer.Destroying()
			Expect(err).NotTo(HaveOccurred())
			destroyed, err := destroyingContainer.Destroy()
			Expect(err).NotTo(HaveOccurred())
			Expect(destroyed).To(BeTrue())
		})

		It("returns orphaned volumes", func() {
			createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
			Expect(err).NotTo(HaveOccurred())
			createdHandles := []string{}

			for _, vol := range createdVolumes {
				createdHandles = append(createdHandles, vol.Handle())
				Expect(vol.WorkerName()).To(Equal("default-worker"))
			}
			Expect(createdHandles).To(Equal(expectedCreatedHandles))

			destroyingHandles := []string{}
			for _, vol := range destoryingVolumes {
				destroyingHandles = append(destroyingHandles, vol.Handle())
				Expect(vol.WorkerName()).To(Equal("default-worker"))
			}

			Expect(destroyingHandles).To(Equal(expectedDestroyingHandles))
		})

		Context("when worker is stalled", func() {
			BeforeEach(func() {
				var err error
				defaultWorker, err = workerFactory.SaveWorker(defaultWorkerPayload, -10*time.Minute)
				Expect(err).NotTo(HaveOccurred())
				stalledWorkers, err := workerLifecycle.StallUnresponsiveWorkers()
				Expect(err).NotTo(HaveOccurred())
				Expect(stalledWorkers).To(ContainElement(defaultWorker.Name()))
			})

			It("does not return volumes", func() {
				createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(createdVolumes).To(HaveLen(0))
				Expect(destoryingVolumes).To(HaveLen(0))
			})
		})

		Context("when worker is landed", func() {
			BeforeEach(func() {
				err := defaultWorker.Land()
				Expect(err).NotTo(HaveOccurred())
				landedWorkers, err := workerLifecycle.LandFinishedLandingWorkers()
				Expect(err).NotTo(HaveOccurred())
				Expect(landedWorkers).To(ContainElement(defaultWorker.Name()))
			})

			It("does not return volumes", func() {
				createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(createdVolumes).To(HaveLen(0))
				Expect(destoryingVolumes).To(HaveLen(0))
			})
		})
	})

	Describe("GetFailedVolumes", func() {

		var expectedFailedHandles []string

		BeforeEach(func() {
			creatingContainer, err := defaultTeam.CreateContainer(defaultWorker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
				Type:     "task",
				StepName: "some-task",
			})
			Expect(err).ToNot(HaveOccurred())

			expectedFailedHandles = []string{}

			creatingVolume1, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-1")
			Expect(err).NotTo(HaveOccurred())
			failedVolume1, err := creatingVolume1.Failed()
			Expect(err).NotTo(HaveOccurred())

			expectedFailedHandles = append(expectedFailedHandles, failedVolume1.Handle())
		})

		It("returns failed volumes", func() {
			failedVolumes, err := volumeFactory.GetFailedVolumes()
			Expect(err).NotTo(HaveOccurred())
			failedHandles := []string{}

			for _, vol := range failedVolumes {
				failedHandles = append(failedHandles, vol.Handle())
				Expect(vol.WorkerName()).To(Equal("default-worker"))
			}
			Expect(failedHandles).To(Equal(expectedFailedHandles))
		})

	})

	Describe("FindBaseResourceTypeVolume", func() {
		var usedWorkerBaseResourceType *db.UsedWorkerBaseResourceType
		BeforeEach(func() {
			workerBaseResourceTypeFactory := db.NewWorkerBaseResourceTypeFactory(dbConn)
			var err error
			var found bool
			usedWorkerBaseResourceType, found, err = workerBaseResourceTypeFactory.Find("some-base-resource-type", defaultWorker)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
		})

		Context("when there is a created volume for base resource type", func() {
			var existingVolume db.CreatedVolume

			BeforeEach(func() {
				var err error
				volume, err := volumeFactory.CreateBaseResourceTypeVolume(defaultTeam.ID(), usedWorkerBaseResourceType)
				Expect(err).NotTo(HaveOccurred())
				existingVolume, err = volume.Created()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns created volume", func() {
				creatingVolume, createdVolume, err := volumeFactory.FindBaseResourceTypeVolume(defaultTeam.ID(), usedWorkerBaseResourceType)
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingVolume).To(BeNil())
				Expect(createdVolume).ToNot(BeNil())
				Expect(createdVolume.Handle()).To(Equal(existingVolume.Handle()))
			})
		})

		Context("when there is a creating volume for base resource type", func() {
			var existingVolume db.CreatingVolume

			BeforeEach(func() {
				var err error
				existingVolume, err = volumeFactory.CreateBaseResourceTypeVolume(defaultTeam.ID(), usedWorkerBaseResourceType)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns creating volume", func() {
				creatingVolume, createdVolume, err := volumeFactory.FindBaseResourceTypeVolume(defaultTeam.ID(), usedWorkerBaseResourceType)
				Expect(err).NotTo(HaveOccurred())
				Expect(creatingVolume).ToNot(BeNil())
				Expect(creatingVolume.Handle()).To(Equal(existingVolume.Handle()))
				Expect(createdVolume).To(BeNil())
			})
		})
	})

	Describe("FindResourceCacheVolume", func() {
		var usedResourceCache *db.UsedResourceCache

		BeforeEach(func() {
			build, err := defaultPipeline.CreateOneOffBuild()
			Expect(err).NotTo(HaveOccurred())

			usedResourceCache, err = resourceCacheFactory.FindOrCreateResourceCache(
				logger,
				db.ForBuild(build.ID()),
				"some-type",
				atc.Version{"some": "version"},
				atc.Source{
					"some": "source",
				},
				atc.Params{"some": "params"},
				creds.NewVersionedResourceTypes(
					template.StaticVariables{"source-param": "some-secret-sauce"},
					atc.VersionedResourceTypes{
						atc.VersionedResourceType{
							ResourceType: atc.ResourceType{
								Name: "some-type",
								Type: "some-base-resource-type",
								Source: atc.Source{
									"some-type": "source",
								},
							},
							Version: atc.Version{"some-type": "version"},
						},
					},
				),
			)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when there is a created volume for resource cache", func() {
			var existingVolume db.CreatedVolume

			BeforeEach(func() {
				var err error
				creatingContainer, err := defaultTeam.CreateContainer(defaultWorker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
					Type:     "get",
					StepName: "some-resource",
				})
				Expect(err).ToNot(HaveOccurred())

				resourceCacheVolume, err := volumeFactory.CreateContainerVolume(defaultTeam.ID(), defaultWorker.Name(), creatingContainer, "some-path-4")
				Expect(err).NotTo(HaveOccurred())

				existingVolume, err = resourceCacheVolume.Created()
				Expect(err).NotTo(HaveOccurred())

				err = existingVolume.InitializeResourceCache(usedResourceCache)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns created volume", func() {
				createdVolume, found, err := volumeFactory.FindResourceCacheVolume(defaultWorker.Name(), usedResourceCache)
				Expect(err).NotTo(HaveOccurred())
				Expect(createdVolume.Handle()).To(Equal(existingVolume.Handle()))
				Expect(found).To(BeTrue())
			})
		})
	})
})
