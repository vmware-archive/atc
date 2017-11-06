package gc_test

import (
	"errors"
	"time"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/gc"
	"github.com/concourse/atc/gc/gcfakes"
	"github.com/concourse/atc/worker/workerfakes"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc/db/dbfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContainerCollector", func() {
	var (
		fakeContainerRepository *dbfakes.FakeContainerRepository
		fakeJobRunner           *gcfakes.FakeWorkerJobRunner

		logger *lagertest.TestLogger

		fakeWorker       *workerfakes.FakeWorker
		fakeGardenClient *gardenfakes.FakeClient

		creatingContainer *dbfakes.FakeCreatingContainer

		collector gc.Collector
	)

	BeforeEach(func() {
		fakeContainerRepository = new(dbfakes.FakeContainerRepository)

		fakeWorker = new(workerfakes.FakeWorker)
		fakeGardenClient = new(gardenfakes.FakeClient)
		fakeWorker.GardenClientReturns(fakeGardenClient)
		fakeJobRunner = new(gcfakes.FakeWorkerJobRunner)
		fakeJobRunner.TryStub = func(logger lager.Logger, workerName string, job gc.Job) {
			job.Run(fakeWorker)
		}

		logger = lagertest.NewTestLogger("test")

		collector = gc.NewContainerCollector(
			logger,
			fakeContainerRepository,
			fakeJobRunner,
		)
	})

	Describe("Run", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = collector.Run()
		})

		It("succeeds", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("Failed Containers", func() {
			var (
				failedContainer1 *dbfakes.FakeFailedContainer
				failedContainer2 *dbfakes.FakeFailedContainer
			)

			BeforeEach(func() {
				failedContainer1 = new(dbfakes.FakeFailedContainer)
				failedContainer1.HandleReturns("some-handle-1")
				failedContainer1.WorkerNameReturns("bar")

				failedContainer2 = new(dbfakes.FakeFailedContainer)
				failedContainer2.HandleReturns("some-handle-2")
				failedContainer2.WorkerNameReturns("bar")

				fakeContainerRepository.FindFailedContainersReturns(
					[]db.FailedContainer{
						failedContainer1,
						failedContainer2,
					},
					nil,
				)
			})

			Context("when there are failed containers", func() {
				It("deletes them from the database", func() {
					Expect(failedContainer1.DestroyCallCount()).To(Equal(1))
					Expect(failedContainer2.DestroyCallCount()).To(Equal(1))
				})

				Context("when deleting one of the containers fails", func() {
					BeforeEach(func() {
						failedContainer1.DestroyReturns(false, errors.New("There is no failure except in no longer trying"))
					})

					It("still deletes the other failed containers", func() {
						Expect(failedContainer2.DestroyCallCount()).To(Equal(1))
					})
				})

				Context("when finding failed containers fails", func() {
					BeforeEach(func() {
						fakeContainerRepository.FindFailedContainersReturns(
							[]db.FailedContainer{},
							errors.New("You have to be able to accept failure to get better"),
						)
					})

					It("still tries to remove the orphaned containers", func() {
						Expect(fakeContainerRepository.FindOrphanedContainersCallCount()).To(Equal(1))
					})

				})
			})
		})

		Describe("Orphaned Containers", func() {

			var (
				createdContainer               *dbfakes.FakeCreatedContainer
				destroyingContainerFromCreated *dbfakes.FakeDestroyingContainer

				destroyingContainer *dbfakes.FakeDestroyingContainer
			)

			BeforeEach(func() {
				creatingContainer = new(dbfakes.FakeCreatingContainer)
				creatingContainer.HandleReturns("some-handle-1")

				createdContainer = new(dbfakes.FakeCreatedContainer)
				createdContainer.HandleReturns("some-handle-2")
				createdContainer.WorkerNameReturns("foo")

				destroyingContainerFromCreated = new(dbfakes.FakeDestroyingContainer)
				createdContainer.DestroyingReturns(destroyingContainerFromCreated, nil)
				destroyingContainerFromCreated.HandleReturns("some-handle-2")
				destroyingContainerFromCreated.WorkerNameReturns("foo")

				destroyingContainer = new(dbfakes.FakeDestroyingContainer)
				destroyingContainer.HandleReturns("some-handle-3")
				destroyingContainer.WorkerNameReturns("bar")

				fakeContainerRepository.FindOrphanedContainersReturns(
					[]db.CreatingContainer{
						creatingContainer,
					},
					[]db.CreatedContainer{
						createdContainer,
					},
					[]db.DestroyingContainer{
						destroyingContainer,
					},
					nil,
				)

				destroyingContainerFromCreated.DestroyReturns(true, nil)
				destroyingContainer.DestroyReturns(true, nil)
			})

			Context("when there are created containers in hijacked state", func() {
				var (
					fakeGardenContainer *gardenfakes.FakeContainer
				)

				BeforeEach(func() {
					createdContainer.IsHijackedReturns(true)
					fakeGardenContainer = new(gardenfakes.FakeContainer)
				})

				Context("when container still exists in garden", func() {
					BeforeEach(func() {
						fakeGardenClient.LookupReturns(fakeGardenContainer, nil)
					})

					It("tells garden to set the TTL to 5 Min", func() {
						Expect(fakeGardenClient.LookupCallCount()).To(Equal(1))
						lookupHandle := fakeGardenClient.LookupArgsForCall(0)
						Expect(lookupHandle).To(Equal("some-handle-2"))

						Expect(fakeGardenContainer.SetGraceTimeCallCount()).To(Equal(1))
						graceTime := fakeGardenContainer.SetGraceTimeArgsForCall(0)
						Expect(graceTime).To(Equal(5 * time.Minute))
					})

					It("marks container as discontinued in database", func() {
						Expect(createdContainer.DiscontinueCallCount()).To(Equal(1))
					})
				})

				Context("when container does not exist in garden", func() {
					BeforeEach(func() {
						fakeGardenClient.LookupReturns(nil, garden.ContainerNotFoundError{Handle: "im-fake-and-still-hijacked"})
					})

					It("marks container as destroying", func() {
						Expect(createdContainer.DestroyingCallCount()).To(Equal(1))
					})
				})
			})

			It("marks all found containers as destroying, tells garden to destroy it, and then removes it from the DB", func() {
				Expect(fakeContainerRepository.FindOrphanedContainersCallCount()).To(Equal(1))

				Expect(createdContainer.DestroyingCallCount()).To(Equal(1))
				Expect(destroyingContainerFromCreated.DestroyCallCount()).To(Equal(1))

				Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))

				Expect(fakeJobRunner.TryCallCount()).To(Equal(2))
				_, try1Worker, _ := fakeJobRunner.TryArgsForCall(0)
				Expect(try1Worker).To(Equal("foo"))
				_, try3Worker, _ := fakeJobRunner.TryArgsForCall(1)
				Expect(try3Worker).To(Equal("bar"))

				Expect(fakeGardenClient.DestroyCallCount()).To(Equal(2))
				Expect(fakeGardenClient.DestroyArgsForCall(0)).To(Equal("some-handle-2"))
				Expect(fakeGardenClient.DestroyArgsForCall(1)).To(Equal("some-handle-3"))
			})

			Context("when there are destroying containers that are discontinued", func() {
				BeforeEach(func() {
					destroyingContainer.IsDiscontinuedReturns(true)
				})

				Context("when container exists in garden", func() {
					BeforeEach(func() {
						fakeGardenClient.LookupReturns(new(gardenfakes.FakeContainer), nil)
					})

					It("does not delete container and lets it expire in garden first", func() {
						Expect(fakeGardenClient.DestroyCallCount()).To(Equal(1))
						Expect(fakeGardenClient.DestroyArgsForCall(0)).To(Equal("some-handle-2"))

						Expect(destroyingContainer.DestroyCallCount()).To(Equal(0))
					})
				})

				Context("when container does not exist in garden", func() {
					BeforeEach(func() {
						fakeGardenClient.LookupReturns(nil, garden.ContainerNotFoundError{})
					})

					It("deletes container in database", func() {
						Expect(fakeGardenClient.DestroyCallCount()).To(Equal(1))
						Expect(fakeGardenClient.DestroyArgsForCall(0)).To(Equal("some-handle-2"))

						Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))
					})
				})
			})

			Context("when finding containers for deletion fails", func() {
				BeforeEach(func() {
					fakeContainerRepository.FindOrphanedContainersReturns(nil, nil, nil, errors.New("some-error"))
				})

				It("returns and logs the error", func() {
					Expect(err).To(MatchError("container collector failed"))
					Expect(fakeContainerRepository.FindOrphanedContainersCallCount()).To(Equal(1))
					Expect(fakeJobRunner.TryCallCount()).To(Equal(0))
				})
			})

			Context("when destroying a garden container errors", func() {
				BeforeEach(func() {
					fakeGardenClient.DestroyStub = func(handle string) error {
						switch handle {
						case "some-handle-1":
							return errors.New("some-error")
						case "some-handle-2":
							return nil
						case "some-handle-3":
							return nil
						default:
							return nil
						}
					}
				})

				It("continues destroying the rest of the containers", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeJobRunner.TryCallCount()).To(Equal(2))
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(2))

					Expect(destroyingContainerFromCreated.DestroyCallCount()).To(Equal(1))
					Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))
				})
			})

			Context("when destroying a garden container errors because container is not found", func() {
				BeforeEach(func() {
					fakeGardenClient.DestroyStub = func(handle string) error {
						switch handle {
						case "some-handle-1":
							return garden.ContainerNotFoundError{Handle: "some-handle"}
						case "some-handle-2":
							return nil
						case "some-handle-3":
							return nil
						default:
							return nil
						}
					}
				})

				It("deletes container from database", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeJobRunner.TryCallCount()).To(Equal(2))
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(2))

					Expect(destroyingContainerFromCreated.DestroyCallCount()).To(Equal(1))
					Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))
				})
			})

			Context("when destroying a container in the DB errors", func() {
				BeforeEach(func() {
					destroyingContainerFromCreated.DestroyReturns(false, errors.New("some-error"))
				})

				It("continues destroying the rest of the containers", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeJobRunner.TryCallCount()).To(Equal(2))
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(2))
					Expect(destroyingContainerFromCreated.DestroyCallCount()).To(Equal(1))
					Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))
				})
			})

			Context("when it can't find a container to destroy", func() {
				BeforeEach(func() {
					destroyingContainerFromCreated.DestroyReturns(false, nil)
				})

				It("continues destroying the rest of the containers", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeJobRunner.TryCallCount()).To(Equal(2))
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(2))
					Expect(destroyingContainerFromCreated.DestroyCallCount()).To(Equal(1))
					Expect(destroyingContainer.DestroyCallCount()).To(Equal(1))
				})
			})
		})

	})
})
