package worker_test

import (
	"context"
	"errors"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	. "github.com/concourse/atc/worker"
	"github.com/concourse/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pool", func() {
	var (
		logger       *lagertest.TestLogger
		fakeProvider *workerfakes.FakeWorkerProvider
		fakeStrategy *workerfakes.FakeContainerPlacementStrategy
		pool         Client
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeProvider = new(workerfakes.FakeWorkerProvider)
		fakeStrategy = new(workerfakes.FakeContainerPlacementStrategy)

		pool = NewPool(fakeProvider, fakeStrategy)
	})

	Describe("Satisfying", func() {
		var (
			spec WorkerSpec

			satisfyingErr    error
			satisfyingWorker Worker
			resourceTypes    creds.VersionedResourceTypes
		)

		BeforeEach(func() {
			spec = WorkerSpec{
				Platform: "some-platform",
				Tags:     []string{"step", "tags"},
			}

			variables := template.StaticVariables{
				"secret-source": "super-secret-source",
			}

			resourceTypes = creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "some-resource-type",
						Type:   "some-underlying-type",
						Source: atc.Source{"some": "((secret-source))"},
					},
					Version: atc.Version{"some": "version"},
				},
			})
		})

		JustBeforeEach(func() {
			satisfyingWorker, satisfyingErr = pool.Satisfying(logger, spec, resourceTypes)
		})

		Context("with multiple satisfactory workers", func() {
			var (
				workerA *workerfakes.FakeWorker
				workerB *workerfakes.FakeWorker
				workerC *workerfakes.FakeWorker
			)

			BeforeEach(func() {
				workerA = new(workerfakes.FakeWorker)
				workerB = new(workerfakes.FakeWorker)
				workerC = new(workerfakes.FakeWorker)
				fakeProvider.AllSatisfyingReturns([]Worker{workerA, workerB}, nil)
			})

			It("succeeds", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})

			It("returns a random worker satisfying the spec", func() {
				chosenCount := map[Worker]int{workerA: 0, workerB: 0, workerC: 0}
				for i := 0; i < 100; i++ {
					satisfyingWorker, satisfyingErr = pool.Satisfying(logger, spec, resourceTypes)
					Expect(satisfyingErr).NotTo(HaveOccurred())
					chosenCount[satisfyingWorker]++
				}
				Expect(chosenCount[workerA]).To(BeNumerically("~", chosenCount[workerB], 50))
				Expect(chosenCount[workerC]).To(BeZero())
			})

		})

		Context("when getting the satisfactory workers fails (no workers, no satisfactory workers)", func() {
			disaster := ErrNoWorkers

			BeforeEach(func() {
				fakeProvider.AllSatisfyingReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(satisfyingErr).To(Equal(disaster))
			})
		})
	})

	Describe("FindContainerByHandle", func() {
		var (
			foundContainer Container
			found          bool
			findErr        error
		)

		JustBeforeEach(func() {
			foundContainer, found, findErr = pool.FindContainerByHandle(
				logger,
				4567,
				"some-handle",
			)
		})

		Context("when a worker is found with the container", func() {
			var fakeWorker *workerfakes.FakeWorker
			var fakeContainer *workerfakes.FakeContainer

			BeforeEach(func() {
				fakeWorker = new(workerfakes.FakeWorker)
				fakeProvider.FindWorkerForContainerReturns(fakeWorker, true, nil)

				fakeContainer = new(workerfakes.FakeContainer)
				fakeWorker.FindContainerByHandleReturns(fakeContainer, true, nil)
			})

			It("succeeds", func() {
				Expect(found).To(BeTrue())
				Expect(findErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(foundContainer).To(Equal(fakeContainer))
			})

			It("finds on the particular worker", func() {
				Expect(fakeWorker.FindContainerByHandleCallCount()).To(Equal(1))

				_, actualTeamID, actualHandle := fakeProvider.FindWorkerForContainerArgsForCall(0)
				Expect(actualTeamID).To(Equal(4567))
				Expect(actualHandle).To(Equal("some-handle"))
			})
		})

		Context("when no worker is found with the container", func() {
			BeforeEach(func() {
				fakeProvider.FindWorkerForContainerReturns(nil, false, nil)
			})

			It("returns no container, false, and no error", func() {
				Expect(foundContainer).To(BeNil())
				Expect(found).To(BeFalse())
				Expect(findErr).ToNot(HaveOccurred())
			})
		})
	})

	Describe("FindOrCreateContainer", func() {
		var (
			ctx                       context.Context
			fakeImageFetchingDelegate *workerfakes.FakeImageFetchingDelegate
			metadata                  db.ContainerMetadata
			spec                      ContainerSpec
			resourceTypes             creds.VersionedResourceTypes
			fakeOwner                 *dbfakes.FakeContainerOwner

			fakeContainer *workerfakes.FakeContainer

			createdContainer Container
			createErr        error

			incompatibleWorker *workerfakes.FakeWorker
			compatibleWorker   *workerfakes.FakeWorker
		)

		BeforeEach(func() {
			ctx = context.Background()

			fakeImageFetchingDelegate = new(workerfakes.FakeImageFetchingDelegate)

			fakeOwner = new(dbfakes.FakeContainerOwner)

			fakeInput1 := new(workerfakes.FakeInputSource)
			fakeInput1AS := new(workerfakes.FakeArtifactSource)
			fakeInput1AS.VolumeOnStub = func(worker Worker) (Volume, bool, error) {
				switch worker {
				case compatibleWorkerOneCache1, compatibleWorkerOneCache2, compatibleWorkerTwoCaches:
					return new(workerfakes.FakeVolume), true, nil
				default:
					return nil, false, nil
				}
			}
			fakeInput1.SourceReturns(fakeInput1AS)

			fakeInput2 := new(workerfakes.FakeInputSource)
			fakeInput2AS := new(workerfakes.FakeArtifactSource)
			fakeInput2AS.VolumeOnStub = func(worker Worker) (Volume, bool, error) {
				switch worker {
				case compatibleWorkerTwoCaches:
					return new(workerfakes.FakeVolume), true, nil
				default:
					return nil, false, nil
				}
			}
			fakeInput2.SourceReturns(fakeInput2AS)

			spec = ContainerSpec{
				ImageSpec: ImageSpec{ResourceType: "some-type"},

				TeamID: 4567,

				Inputs: []InputSource{
					fakeInput1,
					fakeInput2,
				},
			}

			variables := template.StaticVariables{
				"secret-source": "super-secret-source",
			}

			resourceTypes = creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-b",
						Type:   "custom-type-a",
						Source: atc.Source{"some": "((secret-source))"},
					},
					Version: atc.Version{"some": "version"},
				},
			})
			fakeContainer = new(workerfakes.FakeContainer)

			incompatibleWorker = new(workerfakes.FakeWorker)
			compatibleWorker = new(workerfakes.FakeWorker)
			compatibleWorker.FindOrCreateContainerReturns(fakeContainer, nil)

			fakeProvider.AllSatisfyingReturns([]Worker{compatibleWorker}, nil)
		})

		JustBeforeEach(func() {
			createdContainer, createErr = pool.FindOrCreateContainer(
				ctx,
				logger,
				fakeImageFetchingDelegate,
				fakeOwner,
				metadata,
				spec,
				resourceTypes,
			)
		})

		Context("when a worker is found with the container", func() {
			var fakeWorker *workerfakes.FakeWorker

			BeforeEach(func() {
				fakeWorker = new(workerfakes.FakeWorker)
				fakeProvider.FindWorkerForContainerByOwnerReturns(fakeWorker, true, nil)
				fakeWorker.FindOrCreateContainerReturns(fakeContainer, nil)
			})

			It("succeeds", func() {
				Expect(createErr).NotTo(HaveOccurred())
			})

			It("returns the created container", func() {
				Expect(createdContainer).To(Equal(fakeContainer))
			})

			It("'find-or-create's on the particular worker", func() {
				Expect(fakeWorker.FindOrCreateContainerCallCount()).To(Equal(1))

				_, actualTeamID, actualOwner := fakeProvider.FindWorkerForContainerByOwnerArgsForCall(0)
				Expect(actualTeamID).To(Equal(4567))
				Expect(actualOwner).To(Equal(fakeOwner))
			})
		})

		Context("when no worker is found with the container", func() {
			BeforeEach(func() {
				fakeProvider.FindWorkerForContainerByOwnerReturns(nil, false, nil)
			})

			Context("with no workers available", func() {
				BeforeEach(func() {
					fakeProvider.AllSatisfyingReturns(nil, ErrNoWorkers)
				})

				It("returns ErrNoWorkers", func() {
					Expect(createErr).To(Equal(ErrNoWorkers))
				})
			})

			Context("with no compatible workers available", func() {
				var noCompatibleWorkersErr error
				BeforeEach(func() {
					noCompatibleWorkersErr = NoCompatibleWorkersError{
						Spec:    spec.WorkerSpec(),
						Workers: []Worker{incompatibleWorker},
					}
					fakeProvider.AllSatisfyingReturns([]Worker{}, noCompatibleWorkersErr)
				})

				It("returns NoCompatibleWorkersError", func() {
					Expect(createErr).To(Equal(noCompatibleWorkersErr))
				})
			})

			Context("with compatible workers available", func() {
				BeforeEach(func() {
					fakeProvider.AllSatisfyingReturns([]Worker{
						compatibleWorker,
					}, nil)
				})

				Context("when strategy returns a worker", func() {
					BeforeEach(func() {
						fakeStrategy.ChooseReturns(compatibleWorker, nil)
					})

					It("chooses a worker", func() {
						Expect(createErr).ToNot(HaveOccurred())
						Expect(fakeStrategy.ChooseCallCount()).To(Equal(1))
						Expect(compatibleWorker.FindOrCreateContainerCallCount()).To(Equal(1))
						Expect(createdContainer).To(Equal(fakeContainer))
					})
				})

				Context("when strategy errors", func() {
					var (
						strategyError error
					)

					BeforeEach(func() {
						strategyError = errors.New("strategical explosion")
						fakeStrategy.ChooseReturns(nil, strategyError)
					})

					It("returns an error", func() {
						Expect(createErr).To(Equal(strategyError))
					})
				})
			})
		})
	})
})
