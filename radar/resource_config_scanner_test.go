package radar_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/db/lock/lockfakes"
	"github.com/concourse/atc/radar/radarfakes"
	"github.com/concourse/atc/worker"

	. "github.com/concourse/atc/radar"
	"github.com/concourse/atc/resource"
	rfakes "github.com/concourse/atc/resource/resourcefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceConfigScanner", func() {
	var (
		epoch time.Time

		fakeResourceFactory                   *rfakes.FakeResourceFactory
		fakeResourceConfigCheckSessionFactory *dbfakes.FakeResourceConfigCheckSessionFactory
		fakeResourceConfigCheckSession        *dbfakes.FakeResourceConfigCheckSession
		fakeDBPipeline                        *dbfakes.FakePipeline
		fakeClock                             *fakeclock.FakeClock
		interval                              time.Duration
		variables                             creds.Variables

		fakeResourceType      *dbfakes.FakeResourceType
		versionedResourceType atc.VersionedResourceType

		scanner Scanner

		resourceConfig     atc.ResourceConfig
		fakeResourceConfig *dbfakes.FakeResourceConfig
		fakeScannable      *radarfakes.FakeScannable

		fakeLock *lockfakes.FakeLock
		teamID   = 123
	)

	BeforeEach(func() {
		epoch = time.Unix(123, 456).UTC()
		fakeLock = &lockfakes.FakeLock{}
		interval = 1 * time.Minute
		variables = template.StaticVariables{
			"source-params": "some-secret-sauce",
		}

		resourceConfig = atc.ResourceConfig{
			Name:   "some-resource",
			Type:   "git",
			Source: atc.Source{"uri": "some-secret-sauce"},
			Tags:   atc.Tags{"some-tag"},
		}

		versionedResourceType = atc.VersionedResourceType{
			ResourceType: atc.ResourceType{
				Name:   "some-custom-resource",
				Type:   "registry-image",
				Source: atc.Source{"custom": "((source-params))"},
			},
			Version: atc.Version{"custom": "version"},
		}

		fakeResourceFactory = new(rfakes.FakeResourceFactory)
		fakeResourceConfigCheckSession = new(dbfakes.FakeResourceConfigCheckSession)
		fakeResourceConfigCheckSessionFactory = new(dbfakes.FakeResourceConfigCheckSessionFactory)
		fakeResourceType = new(dbfakes.FakeResourceType)
		fakeDBPipeline = new(dbfakes.FakePipeline)
		fakeResourceConfig = new(dbfakes.FakeResourceConfig)
		fakeScannable = new(radarfakes.FakeScannable)

		fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionReturns(fakeResourceConfigCheckSession, nil)

		fakeResourceConfig.IDReturns(123)
		fakeResourceConfigCheckSession.ResourceConfigReturns(fakeResourceConfig)

		fakeDBPipeline.IDReturns(42)
		fakeDBPipeline.NameReturns("some-pipeline")
		fakeDBPipeline.TeamIDReturns(teamID)
		fakeClock = fakeclock.NewFakeClock(epoch)

		fakeDBPipeline.ReloadReturns(true, nil)
		fakeDBPipeline.ResourceTypesReturns([]db.ResourceType{fakeResourceType}, nil)

		fakeResourceType.IDReturns(1)
		fakeResourceType.NameReturns("some-custom-resource")
		fakeResourceType.TypeReturns("registry-image")
		fakeResourceType.SourceReturns(atc.Source{"custom": "((source-params))"})
		fakeResourceType.VersionReturns(atc.Version{"custom": "version"})

		fakeScannable.NameReturns("some-resource")
		fakeScannable.PausedReturns(false)
		fakeScannable.TypeReturns("git")
		fakeScannable.SourceReturns(atc.Source{"uri": "((source-params))"})
		fakeScannable.TagsReturns(atc.Tags{"some-tag"})
		fakeScannable.SetResourceConfigReturns(nil)

		scanner = NewResourceConfigScanner(
			fakeClock,
			fakeResourceFactory,
			fakeResourceConfigCheckSessionFactory,
			interval,
			fakeDBPipeline,
			"https://www.example.com",
			variables,
		)
	})

	Describe("Run", func() {
		var (
			fakeResource   *rfakes.FakeResource
			actualInterval time.Duration
			runErr         error
		)

		BeforeEach(func() {
			fakeResource = new(rfakes.FakeResource)
			fakeResourceFactory.NewResourceReturns(fakeResource, nil)
		})

		JustBeforeEach(func() {
			actualInterval, runErr = scanner.Run(lagertest.NewTestLogger("test"), fakeScannable)
		})

		Context("when the lock cannot be acquired", func() {
			BeforeEach(func() {
				fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckReturns(nil, false, nil)
			})

			It("does not check", func() {
				Expect(fakeResource.CheckCallCount()).To(Equal(0))
			})

			It("returns the configured interval", func() {
				Expect(runErr).To(Equal(ErrFailedToAcquireLock))
				Expect(actualInterval).To(Equal(interval))
			})
		})

		Context("when the lock can be acquired", func() {
			BeforeEach(func() {
				fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckReturns(fakeLock, true, nil)
			})

			It("checks immediately", func() {
				Expect(fakeResource.CheckCallCount()).To(Equal(1))
			})

			It("constructs the resource of the correct type", func() {
				Expect(fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionCallCount()).To(Equal(1))
				_, resourceType, resourceSource, resourceTypes, _ := fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionArgsForCall(0)
				Expect(resourceType).To(Equal("git"))
				Expect(resourceSource).To(Equal(atc.Source{"uri": "some-secret-sauce"}))
				Expect(resourceTypes).To(Equal(creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
					versionedResourceType,
				})))

				Expect(fakeScannable.SetResourceConfigCallCount()).To(Equal(1))
				resourceConfigID := fakeScannable.SetResourceConfigArgsForCall(0)
				Expect(resourceConfigID).To(Equal(123))

				_, _, owner, metadata, resourceSpec, resourceTypes, _ := fakeResourceFactory.NewResourceArgsForCall(0)
				Expect(owner).To(Equal(db.NewResourceConfigCheckSessionContainerOwner(fakeResourceConfigCheckSession, teamID)))
				Expect(metadata).To(Equal(db.ContainerMetadata{
					Type: db.ContainerTypeCheck,
				}))
				Expect(resourceSpec).To(Equal(worker.ContainerSpec{
					ImageSpec: worker.ImageSpec{
						ResourceType: "git",
					},
					Tags:   atc.Tags{"some-tag"},
					TeamID: 123,
				}))
				Expect(resourceTypes).To(Equal(creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
					versionedResourceType,
				})))
			})

			Context("when the scannable object has a specified check interval", func() {
				BeforeEach(func() {
					fakeScannable.CheckEveryReturns("10ms")
				})

				It("leases for the configured interval", func() {
					Expect(fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckCallCount()).To(Equal(1))

					_, leaseInterval, immediate := fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(0)
					Expect(leaseInterval).To(Equal(10 * time.Millisecond))
					Expect(immediate).To(BeFalse())

					Eventually(fakeLock.ReleaseCallCount).Should(Equal(1))
				})

				It("returns configured interval", func() {
					Expect(actualInterval).To(Equal(10 * time.Millisecond))
				})

				Context("when the interval cannot be parsed", func() {
					BeforeEach(func() {
						fakeScannable.CheckEveryReturns("bad-value")
					})

					It("sets the check error", func() {
						Expect(fakeScannable.SetCheckErrorCallCount()).To(Equal(1))

						resourceErr := fakeScannable.SetCheckErrorArgsForCall(0)
						Expect(resourceErr).To(MatchError("time: invalid duration bad-value"))
					})

					It("returns an error", func() {
						Expect(runErr).To(HaveOccurred())
					})
				})
			})

			It("grabs a periodic resource checking lock before checking, breaks lock after done", func() {
				Expect(fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckCallCount()).To(Equal(1))

				_, leaseInterval, immediate := fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(0)
				Expect(leaseInterval).To(Equal(interval))
				Expect(immediate).To(BeFalse())

				Eventually(fakeLock.ReleaseCallCount).Should(Equal(1))
			})

			Context("when the resource uses a custom type", func() {
				BeforeEach(func() {
					fakeScannable.TypeReturns("some-custom-resource")
				})

				Context("and the custom type has a version", func() {
					It("does not fail", func() {
						Expect(runErr).ToNot(HaveOccurred())
					})
				})

				Context("and the custom type does not have a version", func() {
					BeforeEach(func() {
						fakeResourceType.VersionReturns(nil)
					})

					It("returns silently", func() {
						Expect(runErr).ToNot(HaveOccurred())
					})

					// Context("when scanning for the custom type fails", func() {
					// 	var typeScanErr = errors.New("type scan failed")
					// 	BeforeEach(func() {
					// 		fakeScanner.ScanReturns(typeScanErr)
					// 	})
					// 	It("returns the error from scanning for the type", func() {
					// 		Expect(runErr).To(Equal(typeScanErr))
					// 	})
					// })

					// Context("when scanning for the custom type succeeds", func() {
					// 	BeforeEach(func() {
					// 		fakeScanner.ScanReturns(nil)
					// 	})
					// 	It("reloads the resource types", func() {
					// 		Expect(fakeDBPipeline.ResourceTypesCallCount()).To(Equal(2))
					// 	})
					// })
				})
			})

			Context("when there is no current version", func() {
				It("checks from nil", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(BeNil())
				})
			})

			Context("when there is a current version", func() {
				BeforeEach(func() {
					fakeResourceConfigVersion := new(dbfakes.FakeResourceConfigVersion)
					fakeResourceConfigVersion.IDReturns(1)
					fakeResourceConfigVersion.VersionReturns(db.Version{"version": "1"})

					fakeResourceConfig.GetLatestVersionReturns(fakeResourceConfigVersion, true, nil)
				})

				It("checks from it", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(Equal(atc.Version{"version": "1"}))
				})
			})

			Context("when the check returns versions", func() {
				var checkedFrom chan atc.Version

				var nextVersions []atc.Version

				BeforeEach(func() {
					checkedFrom = make(chan atc.Version, 100)

					nextVersions = []atc.Version{
						{"version": "1"},
						{"version": "2"},
						{"version": "3"},
					}

					checkResults := map[int][]atc.Version{
						0: nextVersions,
					}

					check := 0
					fakeResource.CheckStub = func(source atc.Source, from atc.Version) ([]atc.Version, error) {
						defer GinkgoRecover()

						Expect(source).To(Equal(resourceConfig.Source))

						checkedFrom <- from
						result := checkResults[check]
						check++

						return result, nil
					}
				})

				It("saves them all, in order", func() {
					Eventually(fakeResourceConfig.SaveVersionsCallCount).Should(Equal(1))

					versions := fakeResourceConfig.SaveVersionsArgsForCall(0)
					Expect(versions).To(Equal([]atc.Version{
						{"version": "1"},
						{"version": "2"},
						{"version": "3"},
					}))
				})

				Context("when saving versions fails", func() {
					BeforeEach(func() {
						fakeResourceConfig.SaveVersionsReturns(errors.New("failed"))
					})

					It("does not return an error", func() {
						Expect(runErr).NotTo(HaveOccurred())
					})
				})
			})

			Context("when checking fails internally", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeResource.CheckReturns(nil, disaster)
				})

				It("exits with the failure", func() {
					Expect(runErr).To(HaveOccurred())
					Expect(runErr).To(Equal(disaster))
				})
			})

			Context("when checking fails with ErrResourceScriptFailed", func() {
				scriptFail := resource.ErrResourceScriptFailed{}

				BeforeEach(func() {
					fakeResource.CheckReturns(nil, scriptFail)
				})

				It("returns no error", func() {
					Expect(runErr).NotTo(HaveOccurred())
				})
			})

			Context("when the pipeline is paused", func() {
				BeforeEach(func() {
					fakeDBPipeline.CheckPausedReturns(true, nil)
				})

				It("does not check", func() {
					Expect(fakeResource.CheckCallCount()).To(BeZero())
				})

				It("returns the default interval", func() {
					Expect(actualInterval).To(Equal(interval))
				})

				It("does not return an error", func() {
					Expect(runErr).NotTo(HaveOccurred())
				})
			})

			Context("when the scannable object is paused", func() {
				BeforeEach(func() {
					fakeScannable.PausedReturns(true)
				})

				It("does not check", func() {
					Expect(fakeResource.CheckCallCount()).To(BeZero())
				})

				It("returns the default interval", func() {
					Expect(actualInterval).To(Equal(interval))
				})

				It("does not return an error", func() {
					Expect(runErr).NotTo(HaveOccurred())
				})
			})

			Context("when checking if the pipeline is paused fails", func() {
				disaster := errors.New("disaster")

				BeforeEach(func() {
					fakeDBPipeline.CheckPausedReturns(false, disaster)
				})

				It("returns an error", func() {
					Expect(runErr).To(HaveOccurred())
					Expect(runErr).To(Equal(disaster))
				})
			})
		})
	})

	Describe("Scan", func() {
		var (
			fakeResource *rfakes.FakeResource

			scanErr error
		)

		BeforeEach(func() {
			fakeResource = new(rfakes.FakeResource)
			fakeResourceFactory.NewResourceReturns(fakeResource, nil)
		})

		JustBeforeEach(func() {
			scanErr = scanner.Scan(lagertest.NewTestLogger("test"), fakeScannable)
		})

		Context("if the lock can be acquired", func() {
			BeforeEach(func() {
				fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckReturns(fakeLock, true, nil)
			})

			// Context("Parent resource has no version and attempt to Scan fails", func() {
			// 	BeforeEach(func() {
			// 		var fakeGitResourceType *dbfakes.FakeResourceType
			// 		fakeGitResourceType = new(dbfakes.FakeResourceType)

			// 		fakeDBPipeline.ResourceTypesReturns([]db.ResourceType{fakeGitResourceType}, nil)

			// 		fakeGitResourceType.IDReturns(5)
			// 		fakeGitResourceType.NameReturns("git")
			// 		fakeGitResourceType.TypeReturns("registry-image")
			// 		fakeGitResourceType.SourceReturns(atc.Source{"custom": "((source-params))"})
			// 		fakeGitResourceType.VersionReturns(nil)

			// 		fakeScanner.ScanReturns(errors.New("some-resource-type-error"))
			// 	})

			// 	It("fails and returns error", func() {
			// 		Expect(fakeScanner.ScanCallCount()).To(Equal(1))
			// 		_, parentTypeName := fakeScanner.ScanArgsForCall(0)
			// 		Expect(parentTypeName.Name()).To(Equal("git"))

			// 		Expect(scanErr).To(HaveOccurred())
			// 		Expect(scanErr.Error()).To(Equal("some-resource-type-error"))
			// 	})

			// 	It("saves the error to check_error on resource row in db", func() {
			// 		Expect(fakeScannable.SetCheckErrorCallCount()).To(Equal(1))

			// 		err := fakeScannable.SetCheckErrorArgsForCall(0)
			// 		Expect(err).To(HaveOccurred())
			// 		Expect(err.Error()).To(Equal("some-resource-type-error"))
			// 	})
			// })

			It("succeeds", func() {
				Expect(scanErr).NotTo(HaveOccurred())
			})

			It("constructs the resource of the correct type", func() {
				Expect(fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionCallCount()).To(Equal(1))
				_, resourceType, resourceSource, resourceTypes, _ := fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionArgsForCall(0)
				Expect(resourceType).To(Equal("git"))
				Expect(resourceSource).To(Equal(atc.Source{"uri": "some-secret-sauce"}))
				Expect(resourceTypes).To(Equal(creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
					versionedResourceType,
				})))

				Expect(fakeScannable.SetResourceConfigCallCount()).To(Equal(1))
				resourceConfigID := fakeScannable.SetResourceConfigArgsForCall(0)
				Expect(resourceConfigID).To(Equal(123))

				_, _, owner, metadata, resourceSpec, resourceTypes, _ := fakeResourceFactory.NewResourceArgsForCall(0)
				Expect(owner).To(Equal(db.NewResourceConfigCheckSessionContainerOwner(fakeResourceConfigCheckSession, teamID)))
				Expect(metadata).To(Equal(db.ContainerMetadata{
					Type: db.ContainerTypeCheck,
				}))
				Expect(resourceSpec).To(Equal(worker.ContainerSpec{
					ImageSpec: worker.ImageSpec{
						ResourceType: "git",
					},
					Tags:   atc.Tags{"some-tag"},
					TeamID: 123,
				}))
				Expect(resourceTypes).To(Equal(creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
					versionedResourceType,
				})))
			})

			It("grabs an immediate resource checking lock before checking, breaks lock after done", func() {
				Expect(fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckCallCount()).To(Equal(1))

				_, leaseInterval, immediate := fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(0)
				Expect(leaseInterval).To(Equal(interval))
				Expect(immediate).To(BeTrue())

				Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
			})

			Context("when creating the resource config fails", func() {
				BeforeEach(func() {
					fakeResourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSessionReturns(nil, errors.New("catastrophe"))
				})

				It("sets the check error and returns the error", func() {
					Expect(scanErr).To(HaveOccurred())
					Expect(fakeScannable.SetCheckErrorCallCount()).To(Equal(1))

					resourceErr := fakeScannable.SetCheckErrorArgsForCall(0)
					Expect(resourceErr).To(MatchError("catastrophe"))
				})
			})

			Context("when updating the resource config id on the resource fails", func() {
				BeforeEach(func() {
					fakeScannable.SetResourceConfigReturns(errors.New("catastrophe"))
				})

				It("sets the check error and returns the error", func() {
					Expect(scanErr).To(HaveOccurred())
					Expect(fakeResourceConfig.SetCheckErrorCallCount()).To(Equal(1))

					resourceErr := fakeResourceConfig.SetCheckErrorArgsForCall(0)
					Expect(resourceErr).To(MatchError("catastrophe"))
				})
			})

			Context("when creating the resource checker fails", func() {
				BeforeEach(func() {
					fakeResourceFactory.NewResourceReturns(nil, errors.New("catastrophe"))
				})

				It("sets the check error and returns the error", func() {
					Expect(scanErr).To(HaveOccurred())
					Expect(fakeResourceConfig.SetCheckErrorCallCount()).To(Equal(1))

					resourceErr := fakeResourceConfig.SetCheckErrorArgsForCall(0)
					Expect(resourceErr).To(MatchError("catastrophe"))
				})
			})

			Context("when the resource config has a specified check interval", func() {
				BeforeEach(func() {
					fakeScannable.CheckEveryReturns("10ms")
				})

				It("leases for the configured interval", func() {
					Expect(fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckCallCount()).To(Equal(1))

					_, leaseInterval, immediate := fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(0)
					Expect(leaseInterval).To(Equal(10 * time.Millisecond))
					Expect(immediate).To(BeTrue())

					Eventually(fakeLock.ReleaseCallCount).Should(Equal(1))
				})

				Context("when the interval cannot be parsed", func() {
					BeforeEach(func() {
						fakeScannable.CheckEveryReturns("bad-value")
					})

					It("sets the check error and returns the error", func() {
						Expect(scanErr).To(HaveOccurred())
						Expect(fakeScannable.SetCheckErrorCallCount()).To(Equal(1))

						resourceErr := fakeScannable.SetCheckErrorArgsForCall(0)
						Expect(resourceErr).To(MatchError("time: invalid duration bad-value"))
					})
				})
			})

			Context("when the lock is not immediately available", func() {
				BeforeEach(func() {
					results := make(chan bool, 4)
					results <- false
					results <- false
					results <- true
					results <- true
					close(results)

					fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckStub = func(logger lager.Logger, interval time.Duration, immediate bool) (lock.Lock, bool, error) {
						if <-results {
							return fakeLock, true, nil
						} else {
							// allow the sleep to continue
							go fakeClock.WaitForWatcherAndIncrement(time.Second)
							return nil, false, nil
						}
					}
				})

				It("retries every second until it is", func() {
					Expect(fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckCallCount()).To(Equal(3))

					_, leaseInterval, immediate := fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(0)
					Expect(leaseInterval).To(Equal(interval))
					Expect(immediate).To(BeTrue())

					_, leaseInterval, immediate = fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(1)
					Expect(leaseInterval).To(Equal(interval))
					Expect(immediate).To(BeTrue())

					_, leaseInterval, immediate = fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckArgsForCall(2)
					Expect(leaseInterval).To(Equal(interval))
					Expect(immediate).To(BeTrue())

					Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
				})
			})

			It("clears the resource's check error", func() {
				Expect(fakeResourceConfig.SetCheckErrorCallCount()).To(Equal(1))

				err := fakeResourceConfig.SetCheckErrorArgsForCall(0)
				Expect(err).To(BeNil())
			})

			Context("when there is no current version", func() {
				BeforeEach(func() {
					fakeResourceConfig.GetLatestVersionReturns(nil, false, nil)
				})

				It("checks from nil", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(BeNil())
				})
			})

			Context("when getting the current version fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeResourceConfig.GetLatestVersionReturns(nil, false, disaster)
				})

				It("returns the error", func() {
					Expect(scanErr).To(Equal(disaster))
				})

				It("does not check", func() {
					Expect(fakeResource.CheckCallCount()).To(Equal(0))
				})
			})

			Context("when there is a current version", func() {
				var latestVersion db.Version
				BeforeEach(func() {
					latestVersion = db.Version{"version": "1"}
					fakeResourceConfigVersion := new(dbfakes.FakeResourceConfigVersion)
					fakeResourceConfigVersion.IDReturns(1)
					fakeResourceConfigVersion.VersionReturns(db.Version(latestVersion))

					fakeResourceConfig.GetLatestVersionReturns(fakeResourceConfigVersion, true, nil)
				})

				It("checks from it", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(Equal(atc.Version{"version": "1"}))
				})

				Context("when the check returns only the latest version", func() {
					BeforeEach(func() {
						fakeResource.CheckReturns([]atc.Version{atc.Version(latestVersion)}, nil)
					})

					It("does not save it", func() {
						Expect(fakeResourceConfig.SaveVersionsCallCount()).To(Equal(0))
					})
				})
			})

			Context("when the check returns versions", func() {
				var checkedFrom chan atc.Version

				var nextVersions []atc.Version

				BeforeEach(func() {
					checkedFrom = make(chan atc.Version, 100)

					nextVersions = []atc.Version{
						{"version": "1"},
						{"version": "2"},
						{"version": "3"},
					}

					checkResults := map[int][]atc.Version{
						0: nextVersions,
					}

					check := 0
					fakeResource.CheckStub = func(source atc.Source, from atc.Version) ([]atc.Version, error) {
						defer GinkgoRecover()

						Expect(source).To(Equal(resourceConfig.Source))

						checkedFrom <- from
						result := checkResults[check]
						check++

						return result, nil
					}
				})

				It("saves them all, in order", func() {
					Expect(fakeResourceConfig.SaveVersionsCallCount()).To(Equal(1))

					versions := fakeResourceConfig.SaveVersionsArgsForCall(0)
					Expect(versions).To(Equal([]atc.Version{
						{"version": "1"},
						{"version": "2"},
						{"version": "3"},
					}))
				})
			})

			Context("when checking fails internally", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeResource.CheckReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(scanErr).To(Equal(disaster))
				})

				It("sets the resource's check error", func() {
					Expect(fakeResourceConfig.SetCheckErrorCallCount()).To(Equal(1))

					err := fakeResourceConfig.SetCheckErrorArgsForCall(0)
					Expect(err).To(Equal(disaster))
				})
			})

			Context("when checking fails with ErrResourceScriptFailed", func() {
				scriptFail := resource.ErrResourceScriptFailed{}

				BeforeEach(func() {
					fakeResource.CheckReturns(nil, scriptFail)
				})

				It("returns no error", func() {
					Expect(scanErr).NotTo(HaveOccurred())
				})

				It("sets the resource's check error", func() {
					Expect(fakeResourceConfig.SetCheckErrorCallCount()).To(Equal(1))

					err := fakeResourceConfig.SetCheckErrorArgsForCall(0)
					Expect(err).To(Equal(scriptFail))
				})
			})
		})
	})

	Describe("ScanFromVersion", func() {
		var (
			fakeResource *rfakes.FakeResource
			fromVersion  atc.Version

			scanErr error
		)

		BeforeEach(func() {
			fakeResource = new(rfakes.FakeResource)
			fakeResourceFactory.NewResourceReturns(fakeResource, nil)
			fromVersion = nil
		})

		JustBeforeEach(func() {
			scanErr = scanner.ScanFromVersion(lagertest.NewTestLogger("test"), fakeScannable, fromVersion)
		})

		Context("if the lock can be acquired", func() {
			BeforeEach(func() {
				fakeResourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheckReturns(fakeLock, true, nil)
			})

			Context("when fromVersion is nil", func() {
				It("checks from nil", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(BeNil())
				})
			})

			Context("when fromVersion is specified", func() {
				BeforeEach(func() {
					fromVersion = atc.Version{
						"version": "1",
					}
				})

				It("checks from it", func() {
					_, version := fakeResource.CheckArgsForCall(0)
					Expect(version).To(Equal(atc.Version{"version": "1"}))
				})
			})

			Context("when checking fails with ErrResourceScriptFailed", func() {
				scriptFail := resource.ErrResourceScriptFailed{}

				BeforeEach(func() {
					fakeResource.CheckReturns(nil, scriptFail)
				})

				It("returns the error", func() {
					Expect(scanErr).To(Equal(scriptFail))
				})
			})
		})
	})
})
