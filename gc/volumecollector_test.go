package gc_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/gc"
	"github.com/concourse/atc/gc/gcfakes"
	"github.com/concourse/atc/worker/workerfakes"
	"github.com/concourse/baggageclaim/baggageclaimfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeCollector", func() {
	var (
		volumeCollector gc.Collector
		fakeJobRunner   *gcfakes.FakeWorkerJobRunner

		volumeFactory      db.VolumeFactory
		workerFactory      db.WorkerFactory
		fakeBCVolume       *baggageclaimfakes.FakeVolume
		createdVolume      db.CreatedVolume
		creatingContainer1 db.CreatingContainer
		creatingContainer2 db.CreatingContainer
		team               db.Team
		worker             db.Worker

		fakeWorker             *workerfakes.FakeWorker
		fakeBaggageclaimClient *baggageclaimfakes.FakeClient
	)

	BeforeEach(func() {
		postgresRunner.Truncate()

		volumeFactory = db.NewVolumeFactory(dbConn)
		workerFactory = db.NewWorkerFactory(dbConn)

		fakeBaggageclaimClient = new(baggageclaimfakes.FakeClient)

		fakeBCVolume = new(baggageclaimfakes.FakeVolume)

		fakeWorker = new(workerfakes.FakeWorker)
		fakeBaggageclaimClient.LookupVolumeReturns(fakeBCVolume, true, nil)
		fakeWorker.BaggageclaimClientReturns(fakeBaggageclaimClient)

		fakeJobRunner = new(gcfakes.FakeWorkerJobRunner)
		fakeJobRunner.TryStub = func(logger lager.Logger, workerName string, job gc.Job) {
			job.Run(fakeWorker)
		}

		logger := lagertest.NewTestLogger("volume-collector")
		volumeCollector = gc.NewVolumeCollector(
			logger,
			volumeFactory,
			fakeJobRunner,
		)
	})

	Describe("Run", func() {

		Context("when there are failed volumes", func() {
			var failedVolume1 db.FailedVolume

			BeforeEach(func() {
				var err error
				team, err = teamFactory.CreateTeam(atc.Team{Name: "some-team"})
				Expect(err).ToNot(HaveOccurred())

				build, err := team.CreateOneOffBuild()
				Expect(err).ToNot(HaveOccurred())

				worker, err = workerFactory.SaveWorker(atc.Worker{
					Name:            "some-worker",
					GardenAddr:      "1.2.3.4:7777",
					BaggageclaimURL: "1.2.3.4:7788",
				}, 5*time.Minute)
				Expect(err).ToNot(HaveOccurred())

				creatingContainer1, err = team.CreateContainer(worker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
					Type:     "task",
					StepName: "some-task",
				})
				Expect(err).ToNot(HaveOccurred())

				creatingVolume1, err := volumeFactory.CreateContainerVolume(team.ID(), worker.Name(), creatingContainer1, "some-path-1")
				Expect(err).NotTo(HaveOccurred())

				failedVolume1, err = creatingVolume1.Failed()
				Expect(err).NotTo(HaveOccurred())
			})

			It("deletes all the failed volumes from the database", func() {
				failedVolumes, err := volumeFactory.GetFailedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(failedVolumes).To(HaveLen(1))

				err = volumeCollector.Run()
				Expect(err).NotTo(HaveOccurred())

				failedVolumes, err = volumeFactory.GetFailedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(failedVolumes).To(HaveLen(0))
			})
		})

		Context("when there are orphaned volumes", func() {
			BeforeEach(func() {
				var err error
				team, err = teamFactory.CreateTeam(atc.Team{Name: "some-team"})
				Expect(err).ToNot(HaveOccurred())

				build, err := team.CreateOneOffBuild()
				Expect(err).ToNot(HaveOccurred())

				worker, err = workerFactory.SaveWorker(atc.Worker{
					Name:            "some-worker",
					GardenAddr:      "1.2.3.4:7777",
					BaggageclaimURL: "1.2.3.4:7788",
				}, 5*time.Minute)
				Expect(err).ToNot(HaveOccurred())

				creatingContainer1, err = team.CreateContainer(worker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
					Type:     "task",
					StepName: "some-task",
				})
				Expect(err).ToNot(HaveOccurred())

				creatingContainer2, err = team.CreateContainer(worker.Name(), db.NewBuildStepContainerOwner(build.ID(), "some-plan"), db.ContainerMetadata{
					Type:     "task",
					StepName: "some-task",
				})
				Expect(err).ToNot(HaveOccurred())

				creatingVolume1, err := volumeFactory.CreateContainerVolume(team.ID(), worker.Name(), creatingContainer1, "some-path-1")
				Expect(err).NotTo(HaveOccurred())
				createdVolume, err = creatingVolume1.Created()
				Expect(err).NotTo(HaveOccurred())

				_, err = volumeFactory.CreateContainerVolume(team.ID(), worker.Name(), creatingContainer2, "some-path-2")
				Expect(err).NotTo(HaveOccurred())

				creatingVolume3, err := volumeFactory.CreateContainerVolume(team.ID(), worker.Name(), creatingContainer1, "some-path-3")
				Expect(err).NotTo(HaveOccurred())
				createdVolume3, err := creatingVolume3.Created()
				Expect(err).NotTo(HaveOccurred())
				_, err = createdVolume3.Destroying()
				Expect(err).NotTo(HaveOccurred())

				createdContainer1, err := creatingContainer1.Created()
				Expect(err).NotTo(HaveOccurred())
				destroyingContainer1, err := createdContainer1.Destroying()
				Expect(err).NotTo(HaveOccurred())
				destroyed, err := destroyingContainer1.Destroy()
				Expect(err).NotTo(HaveOccurred())
				Expect(destroyed).To(BeTrue())
			})

			It("deletes created and destroying orphaned volumes", func() {
				createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(createdVolumes).To(HaveLen(1))
				Expect(destoryingVolumes).To(HaveLen(1))

				err = volumeCollector.Run()
				Expect(err).NotTo(HaveOccurred())

				createdVolumes, destoryingVolumes, err = volumeFactory.GetOrphanedVolumes()
				Expect(err).NotTo(HaveOccurred())
				Expect(createdVolumes).To(HaveLen(0))
				Expect(destoryingVolumes).To(HaveLen(0))

				Expect(fakeBCVolume.DestroyCallCount()).To(Equal(2))
			})

			Context("when destroying the volume in db fails because volume has children", func() {
				BeforeEach(func() {
					_, err := createdVolume.CreateChildForContainer(creatingContainer2, "some-path-1")
					Expect(err).NotTo(HaveOccurred())
				})

				It("leaves the volume in the db", func() {
					createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
					Expect(err).NotTo(HaveOccurred())
					Expect(createdVolumes).To(HaveLen(1))
					createdVolumeHandle := createdVolumes[0].Handle()
					Expect(destoryingVolumes).To(HaveLen(1))

					err = volumeCollector.Run()
					Expect(err).NotTo(HaveOccurred())

					createdVolumes, destoryingVolumes, err = volumeFactory.GetOrphanedVolumes()
					Expect(err).NotTo(HaveOccurred())
					Expect(createdVolumes).To(HaveLen(1))
					Expect(destoryingVolumes).To(HaveLen(0))
					Expect(createdVolumes[0].Handle()).To(Equal(createdVolumeHandle))
				})
			})

			Context("when destroying the volume in baggageclaim fails", func() {
				BeforeEach(func() {
					fakeBCVolume.DestroyReturns(errors.New("oh no!"))
				})

				It("leaves the volume in the db", func() {
					createdVolumes, destoryingVolumes, err := volumeFactory.GetOrphanedVolumes()
					Expect(err).NotTo(HaveOccurred())
					Expect(createdVolumes).To(HaveLen(1))
					Expect(destoryingVolumes).To(HaveLen(1))

					err = volumeCollector.Run()
					Expect(err).NotTo(HaveOccurred())

					createdVolumes, destoryingVolumes, err = volumeFactory.GetOrphanedVolumes()
					Expect(err).NotTo(HaveOccurred())
					Expect(createdVolumes).To(HaveLen(0))
					Expect(destoryingVolumes).To(HaveLen(2))
				})
			})
		})
	})
})
