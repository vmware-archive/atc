package worker_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	gfakes "github.com/cloudfoundry-incubator/garden/gardenfakes"
	"github.com/cloudfoundry-incubator/garden/server"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	. "github.com/concourse/atc/worker"
	"github.com/concourse/atc/worker/workerfakes"
	"github.com/concourse/baggageclaim"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type immediateRetryPolicy struct{}

func (immediateRetryPolicy) DelayFor(uint) (time.Duration, bool) {
	return 0, true
}

var _ = Describe("DBProvider", func() {
	var (
		fakeDB *workerfakes.FakeWorkerDB

		logger *lagertest.TestLogger

		fakeGardenBackend  *gfakes.FakeBackend
		gardenAddr         string
		baggageclaimServer *ghttp.Server
		gardenServer       *server.GardenServer
		provider           WorkerProvider

		fakeImageFetcher          *workerfakes.FakeImageFetcher
		fakeImageFetchingDelegate *workerfakes.FakeImageFetchingDelegate

		workers    []Worker
		workersErr error
	)

	BeforeEach(func() {
		baggageclaimServer = ghttp.NewServer()

		baggageclaimServer.RouteToHandler("POST", "/volumes", ghttp.RespondWithJSONEncoded(
			http.StatusCreated,
			baggageclaim.VolumeResponse{Handle: "vol-handle"},
		))
		baggageclaimServer.RouteToHandler("PUT", "/volumes/vol-handle/ttl", ghttp.RespondWith(
			http.StatusNoContent,
			nil,
		))
		baggageclaimServer.RouteToHandler("GET", "/volumes/vol-handle", ghttp.RespondWithJSONEncoded(
			http.StatusOK,
			baggageclaim.VolumeResponse{Handle: "vol-handle"},
		))
		baggageclaimServer.RouteToHandler("GET", "/volumes/vol-handle/stats", ghttp.RespondWithJSONEncoded(
			http.StatusOK,
			baggageclaim.VolumeStatsResponse{SizeInBytes: 1024},
		))

		fakeDB = new(workerfakes.FakeWorkerDB)
		fakeDB.GetVolumeTTLReturns(1*time.Millisecond, true, nil)
		fakeDB.GetContainerReturns(db.SavedContainer{}, true, nil)
		fakeDB.FindWorkerCheckResourceTypeVersionReturns("", false, nil)

		gardenAddr = fmt.Sprintf("0.0.0.0:%d", 8888+GinkgoParallelNode())
		fakeGardenBackend = new(gfakes.FakeBackend)
		logger = lagertest.NewTestLogger("test")
		gardenServer = server.New("tcp", gardenAddr, 0, fakeGardenBackend, logger)
		err := gardenServer.Start()
		Expect(err).NotTo(HaveOccurred())

		fakeImageFetcher = new(workerfakes.FakeImageFetcher)
		fakeImageFetchingDelegate = new(workerfakes.FakeImageFetchingDelegate)

		provider = NewDBWorkerProvider(logger, fakeDB, nil, immediateRetryPolicy{}, fakeImageFetcher)
	})

	AfterEach(func() {
		gardenServer.Stop()

		Eventually(func() error {
			conn, err := net.Dial("tcp", gardenAddr)
			if err == nil {
				conn.Close()
			}

			return err
		}).Should(HaveOccurred())

		baggageclaimServer.Close()
	})

	Context("when we call to get multiple workers", func() {
		JustBeforeEach(func() {
			workers, workersErr = provider.Workers()
		})

		Context("when the database yields workers", func() {
			BeforeEach(func() {
				fakeDB.WorkersReturns([]db.SavedWorker{
					{
						WorkerInfo: db.WorkerInfo{
							Name:             "some-worker",
							GardenAddr:       gardenAddr,
							BaggageclaimURL:  baggageclaimServer.URL(),
							ActiveContainers: 2,
							ResourceTypes: []atc.WorkerResourceType{
								{Type: "some-resource-a", Image: "some-image-a"},
							},
						},
					},
					{
						WorkerInfo: db.WorkerInfo{
							Name:             "some-other-worker",
							GardenAddr:       gardenAddr,
							ActiveContainers: 2,
							ResourceTypes: []atc.WorkerResourceType{
								{Type: "some-resource-b", Image: "some-image-b"},
							},
						},
					},
				}, nil)
			})

			It("succeeds", func() {
				Expect(workersErr).NotTo(HaveOccurred())
			})

			It("returns a worker for each one", func() {
				Expect(workers).To(HaveLen(2))
			})

			Context("creating the connection to garden", func() {
				var id Identifier
				var spec ContainerSpec

				JustBeforeEach(func() {
					id = Identifier{
						ResourceID: 1234,
					}

					spec = ContainerSpec{
						ImageSpec: ImageSpec{
							ResourceType: "some-resource-a",
						},
					}

					fakeContainer := new(gfakes.FakeContainer)
					fakeContainer.HandleReturns("created-handle")

					fakeGardenBackend.CreateReturns(fakeContainer, nil)
					fakeGardenBackend.LookupReturns(fakeContainer, nil)

					By("connecting to the worker")
					fakeDB.GetWorkerReturns(db.SavedWorker{WorkerInfo: db.WorkerInfo{GardenAddr: gardenAddr}}, true, nil)
					container, err := workers[0].CreateContainer(logger, nil, fakeImageFetchingDelegate, id, Metadata{}, spec, nil)
					Expect(err).NotTo(HaveOccurred())

					err = container.Destroy()
					Expect(err).NotTo(HaveOccurred())

					By("restarting the worker with a new address")
					gardenServer.Stop()

					Eventually(func() error {
						conn, err := net.Dial("tcp", gardenAddr)
						if err == nil {
							conn.Close()
						}

						return err
					}).Should(HaveOccurred())

					gardenAddr = fmt.Sprintf("0.0.0.0:%d", 7777+GinkgoParallelNode())

					gardenServer = server.New("tcp", gardenAddr, 0, fakeGardenBackend, logger)
					err = gardenServer.Start()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Describe("a created container", func() {
				BeforeEach(func() {
					fakeDB.GetWorkerReturns(db.SavedWorker{WorkerInfo: db.WorkerInfo{GardenAddr: gardenAddr}}, true, nil)
				})

				It("calls through to garden", func() {
					id := Identifier{
						ResourceID: 1234,
					}

					spec := ContainerSpec{
						ImageSpec: ImageSpec{
							ResourceType: "some-resource-a",
						},
					}

					fakeContainer := new(gfakes.FakeContainer)
					fakeContainer.HandleReturns("created-handle")

					fakeGardenBackend.CreateReturns(fakeContainer, nil)
					fakeGardenBackend.LookupReturns(fakeContainer, nil)

					container, err := workers[0].CreateContainer(logger, nil, fakeImageFetchingDelegate, id, Metadata{}, spec, nil)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeDB.CreateContainerCallCount()).To(Equal(1))
					createdInfo, _, _, _ := fakeDB.CreateContainerArgsForCall(0)
					Expect(createdInfo.WorkerName).To(Equal("some-worker"))

					Expect(container.Handle()).To(Equal("created-handle"))

					Expect(fakeGardenBackend.CreateCallCount()).To(Equal(1))

					err = container.Destroy()
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeGardenBackend.DestroyCallCount()).To(Equal(1))
					Expect(fakeGardenBackend.DestroyArgsForCall(0)).To(Equal("created-handle"))
				})
			})

			Describe("a looked-up container", func() {
				BeforeEach(func() {
					fakeDB.GetWorkerReturns(db.SavedWorker{WorkerInfo: db.WorkerInfo{GardenAddr: gardenAddr}}, true, nil)
				})

				It("calls through to garden", func() {
					fakeContainer := new(gfakes.FakeContainer)
					fakeContainer.HandleReturns("some-handle")

					fakeGardenBackend.ContainersReturns([]garden.Container{fakeContainer}, nil)
					fakeGardenBackend.LookupReturns(fakeContainer, nil)

					returnContainer := db.SavedContainer{
						Container: db.Container{
							ContainerMetadata: db.ContainerMetadata{
								Handle: "some-handle",
							},
						},
					}
					fakeDB.FindContainerByIdentifierReturns(returnContainer, true, nil)

					container, found, err := workers[0].FindContainerForIdentifier(logger, Identifier{
						ResourceID: 1234,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())

					Expect(container.Handle()).To(Equal("some-handle"))

					err = container.Destroy()
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeGardenBackend.DestroyCallCount()).To(Equal(1))
					Expect(fakeGardenBackend.DestroyArgsForCall(0)).To(Equal("some-handle"))
				})
			})
		})

		Context("when the database fails to return workers", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeDB.WorkersReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(workersErr).To(Equal(disaster))
			})
		})
	})

	Context("when we call to get a single worker", func() {
		var found bool
		var worker Worker

		Context("when looking up workers returns an error", func() {
			It("returns an error", func() {
				fakeDB.GetWorkerReturns(db.SavedWorker{}, true, errors.New("disaster"))

				worker, found, workersErr = provider.GetWorker("a-worker")
				Expect(workersErr).To(HaveOccurred())
				Expect(worker).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when we find no workers", func() {
			It("returns found as false", func() {
				fakeDB.GetWorkerReturns(db.SavedWorker{}, false, nil)

				worker, found, workersErr = provider.GetWorker("no-worker")
				Expect(workersErr).NotTo(HaveOccurred())
				Expect(worker).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})
	})

	Context("when we call to get a container info by identifier", func() {
		It("calls through to the db object", func() {
			provider.FindContainerForIdentifier(Identifier{
				BuildID: 1234,
				PlanID:  atc.PlanID("planid"),
			})

			Expect(fakeDB.FindContainerByIdentifierCallCount()).To(Equal(1))

			Expect(fakeDB.FindContainerByIdentifierArgsForCall(0)).To(Equal(db.ContainerIdentifier{
				BuildID: 1234,
				PlanID:  atc.PlanID("planid"),
			}))
		})
	})
})
