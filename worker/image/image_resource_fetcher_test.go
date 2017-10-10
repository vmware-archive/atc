package image_test

import (
	"archive/tar"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/resource/resourcefakes"
	"github.com/concourse/atc/worker"
	"github.com/concourse/atc/worker/image"
	"github.com/concourse/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Image", func() {
	var fakeResourceFactory *resourcefakes.FakeResourceFactory
	var fakeImageResource *resourcefakes.FakeResource
	var fakeResourceFetcherFactory *resourcefakes.FakeFetcherFactory
	var fakeResourceFetcher *resourcefakes.FakeFetcher
	var fakeResourceFactoryFactory *resourcefakes.FakeResourceFactoryFactory
	var fakeResourceCacheFactory *dbfakes.FakeResourceCacheFactory
	var fakeResourceConfigFactory *dbfakes.FakeResourceConfigFactory
	var fakeCreatingContainer *dbfakes.FakeCreatingContainer

	var imageResourceFetcher image.ImageResourceFetcher

	var stderrBuf *gbytes.Buffer

	var logger lager.Logger
	var imageResource worker.ImageResource
	var version atc.Version
	var signals <-chan os.Signal
	var fakeImageFetchingDelegate *workerfakes.FakeImageFetchingDelegate
	var fakeWorker *workerfakes.FakeWorker
	var fakeClock *fakeclock.FakeClock
	var customTypes creds.VersionedResourceTypes
	var privileged bool

	var fetchedVolume worker.Volume
	var fetchedMetadataReader io.ReadCloser
	var fetchedVersion atc.Version
	var fetchErr error
	var teamID int

	BeforeEach(func() {
		fakeResourceFactory = new(resourcefakes.FakeResourceFactory)
		fakeImageResource = new(resourcefakes.FakeResource)
		fakeResourceFetcherFactory = new(resourcefakes.FakeFetcherFactory)
		fakeResourceFetcher = new(resourcefakes.FakeFetcher)
		fakeResourceFactoryFactory = new(resourcefakes.FakeResourceFactoryFactory)
		fakeResourceConfigFactory = new(dbfakes.FakeResourceConfigFactory)
		fakeResourceFetcherFactory.FetcherForReturns(fakeResourceFetcher)
		fakeResourceFactoryFactory.FactoryForReturns(fakeResourceFactory)
		fakeCreatingContainer = new(dbfakes.FakeCreatingContainer)
		fakeClock = fakeclock.NewFakeClock(time.Now())
		stderrBuf = gbytes.NewBuffer()

		variables := template.StaticVariables{
			"source-param":   "super-secret-sauce",
			"a-source-param": "super-secret-a-source",
			"b-source-param": "super-secret-b-source",
		}

		logger = lagertest.NewTestLogger("test")
		imageResource = worker.ImageResource{
			Type:   "docker",
			Source: creds.NewSource(variables, atc.Source{"some": "((source-param))"}),
			Params: &atc.Params{"some": "params"},
		}
		version = nil
		signals = make(chan os.Signal)
		fakeImageFetchingDelegate = new(workerfakes.FakeImageFetchingDelegate)
		fakeImageFetchingDelegate.StderrReturns(stderrBuf)
		fakeWorker = new(workerfakes.FakeWorker)
		fakeWorker.TagsReturns(atc.Tags{"worker", "tags"})
		teamID = 123

		customTypes = creds.NewVersionedResourceTypes(variables, atc.VersionedResourceTypes{
			{
				ResourceType: atc.ResourceType{
					Name:   "custom-type-a",
					Type:   "base-type",
					Source: atc.Source{"some": "((a-source-param))"},
				},
				Version: atc.Version{"some": "a-version"},
			},
			{
				ResourceType: atc.ResourceType{
					Name:   "custom-type-b",
					Type:   "custom-type-a",
					Source: atc.Source{"some": "((b-source-param))"},
				},
				Version: atc.Version{"some": "b-version"},
			},
		})

		fakeResourceCacheFactory = new(dbfakes.FakeResourceCacheFactory)
	})

	JustBeforeEach(func() {
		imageResourceFetcher = image.NewImageResourceFetcherFactory(
			fakeResourceFetcherFactory,
			fakeResourceFactoryFactory,
			fakeResourceCacheFactory,
			fakeResourceConfigFactory,
			fakeClock,
		).NewImageResourceFetcher(
			fakeWorker,
			imageResource,
			version,
			teamID,
			customTypes,
			fakeImageFetchingDelegate,
		)

		fetchedVolume, fetchedMetadataReader, fetchedVersion, fetchErr = imageResourceFetcher.Fetch(
			logger,
			signals,
			fakeCreatingContainer,
			privileged,
		)
	})

	Context("when no version is specified", func() {
		BeforeEach(func() {
			version = nil
		})

		Context("when initializing the Check resource works", func() {
			var (
				fakeCheckResource *resourcefakes.FakeResource
				fakeBuildResource *resourcefakes.FakeResource
			)

			BeforeEach(func() {
				fakeCheckResource = new(resourcefakes.FakeResource)
				fakeBuildResource = new(resourcefakes.FakeResource)
				fakeResourceFactory.NewResourceReturns(fakeCheckResource, nil)
			})

			Context("when check returns a version", func() {
				BeforeEach(func() {
					fakeCheckResource.CheckReturns([]atc.Version{{"v": "1"}}, nil)
				})

				Context("when saving the version in the database succeeds", func() {
					BeforeEach(func() {
						fakeImageFetchingDelegate.ImageVersionDeterminedReturns(nil)
					})

					Context("when fetching resource fails", func() {
						BeforeEach(func() {
							fakeResourceFetcher.FetchReturns(nil, resource.ErrInterrupted)
						})

						It("returns error", func() {
							Expect(fetchErr).To(Equal(resource.ErrInterrupted))
						})
					})

					Context("when fetching resource succeeds", func() {
						var (
							fakeVersionedSource   *resourcefakes.FakeVersionedSource
							fakeUsedResourceCache *db.UsedResourceCache
						)

						BeforeEach(func() {
							fakeVersionedSource = new(resourcefakes.FakeVersionedSource)
							fakeResourceFetcher.FetchReturns(fakeVersionedSource, nil)

							fakeVersionedSource.StreamOutReturns(tarStreamWith("some-tar-contents"), nil)
							fakeVolume := new(workerfakes.FakeVolume)
							fakeVersionedSource.VolumeReturns(fakeVolume)

							fakeUsedResourceCache = new(db.UsedResourceCache)
							fakeResourceCacheFactory.FindOrCreateResourceCacheReturns(fakeUsedResourceCache, nil)
						})

						Context("when the resource has a volume", func() {
							var (
								fakeVolume *workerfakes.FakeVolume
								volumePath string
							)

							BeforeEach(func() {
								fakeVolume = new(workerfakes.FakeVolume)
								volumePath = "C:/Documents and Settings/Evan/My Documents"

								fakeVolume.PathReturns(volumePath)
								fakeVersionedSource.VolumeReturns(fakeVolume)

								privileged = true
							})

							Context("calling NewResource", func() {
								BeforeEach(func() {
									fakeResourceFactory.NewResourceReturns(fakeCheckResource, nil)
								})

								It("created the 'check' resource with the correct session, with the currently fetching type removed from the set", func() {
									Expect(fakeResourceFactory.NewResourceCallCount()).To(Equal(1))
									_, csig, owner, metadata, resourceSpec, actualCustomTypes, delegate := fakeResourceFactory.NewResourceArgsForCall(0)
									Expect(csig).To(Equal(signals))
									Expect(owner).To(Equal(db.NewImageCheckContainerOwner(fakeCreatingContainer)))
									Expect(metadata).To(Equal(db.ContainerMetadata{
										Type: db.ContainerTypeCheck,
									}))
									Expect(resourceSpec).To(Equal(worker.ContainerSpec{
										ImageSpec: worker.ImageSpec{
											ResourceType: "docker",
										},
										Tags:   []string{"worker", "tags"},
										TeamID: 123,
									}))
									Expect(actualCustomTypes).To(Equal(customTypes))
									Expect(delegate).To(Equal(fakeImageFetchingDelegate))
								})
							})

							It("succeeds", func() {
								Expect(fetchErr).To(BeNil())
							})

							It("returns the image volume", func() {
								Expect(fetchedVolume).To(Equal(fakeVolume))
							})

							It("calls StreamOut on the versioned source with the right metadata path", func() {
								Expect(fakeVersionedSource.StreamOutCallCount()).To(Equal(1))
								Expect(fakeVersionedSource.StreamOutArgsForCall(0)).To(Equal("metadata.json"))
							})

							It("returns a tar stream containing the contents of metadata.json", func() {
								Expect(ioutil.ReadAll(fetchedMetadataReader)).To(Equal([]byte("some-tar-contents")))
							})

							It("has the version on the image", func() {
								Expect(fetchedVersion).To(Equal(atc.Version{"v": "1"}))
							})

							It("created the 'check' resource with the correct session, with the currently fetching type removed from the set", func() {
								Expect(fakeResourceFactory.NewResourceCallCount()).To(Equal(1))
								_, csig, owner, metadata, resourceSpec, actualCustomTypes, delegate := fakeResourceFactory.NewResourceArgsForCall(0)
								Expect(csig).To(Equal(signals))
								Expect(owner).To(Equal(db.NewImageCheckContainerOwner(fakeCreatingContainer)))
								Expect(metadata).To(Equal(db.ContainerMetadata{
									Type: db.ContainerTypeCheck,
								}))
								Expect(resourceSpec).To(Equal(worker.ContainerSpec{
									ImageSpec: worker.ImageSpec{
										ResourceType: "docker",
									},
									Tags:   []string{"worker", "tags"},
									TeamID: 123,
								}))
								Expect(actualCustomTypes).To(Equal(customTypes))
								Expect(delegate).To(Equal(fakeImageFetchingDelegate))
							})

							It("ran 'check' with the right config", func() {
								Expect(fakeCheckResource.CheckCallCount()).To(Equal(1))
								checkSource, checkVersion := fakeCheckResource.CheckArgsForCall(0)
								Expect(checkVersion).To(BeNil())
								Expect(checkSource).To(Equal(atc.Source{"some": "super-secret-sauce"}))
							})

							It("saved the image resource version in the database", func() {
								Expect(fakeImageFetchingDelegate.ImageVersionDeterminedCallCount()).To(Equal(1))
								Expect(fakeImageFetchingDelegate.ImageVersionDeterminedArgsForCall(0)).To(Equal(fakeUsedResourceCache))
							})

							It("fetches resource with correct session", func() {
								Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(1))
								_, session, tags, actualTeamID, actualCustomTypes, resourceInstance, metadata, delegate, _, _ := fakeResourceFetcher.FetchArgsForCall(0)
								Expect(metadata).To(Equal(resource.EmptyMetadata{}))
								Expect(session).To(Equal(resource.Session{
									Metadata: db.ContainerMetadata{
										Type: db.ContainerTypeGet,
									},
								}))
								Expect(tags).To(Equal(atc.Tags{"worker", "tags"}))
								Expect(actualTeamID).To(Equal(teamID))
								Expect(resourceInstance).To(Equal(resource.NewResourceInstance(
									"docker",
									atc.Version{"v": "1"},
									atc.Source{"some": "super-secret-sauce"},
									atc.Params{"some": "params"},
									customTypes,
									&db.UsedResourceCache{},
									db.NewImageGetContainerOwner(fakeCreatingContainer),
								)))
								Expect(actualCustomTypes).To(Equal(customTypes))
								Expect(delegate).To(Equal(fakeImageFetchingDelegate))
								expectedLockName := fmt.Sprintf("%x",
									sha256.Sum256([]byte(
										`{"type":"docker","version":{"v":"1"},"source":{"some":"super-secret-sauce"},"params":{"some":"params"},"worker_name":"fake-worker-name"}`,
									)),
								)
								Expect(resourceInstance.LockName("fake-worker-name")).To(Equal(expectedLockName))
							})

							It("gets the volume", func() {
								Expect(fakeVersionedSource.VolumeCallCount()).To(Equal(1))
							})

							Context("when streaming the metadata out fails", func() {
								disaster := errors.New("nope")

								BeforeEach(func() {
									fakeVersionedSource.StreamOutReturns(nil, disaster)
								})

								It("returns the error", func() {
									Expect(fetchErr).To(Equal(disaster))
								})
							})

							Context("when the resource still does not have a volume for some reason", func() {
								BeforeEach(func() {
									fakeVersionedSource.VolumeReturns(nil)
								})

								It("returns an appropriate error", func() {
									Expect(fetchErr).To(Equal(image.ErrImageGetDidNotProduceVolume))
								})
							})
						})
					})
				})

				Context("when saving the version in the database fails", func() {
					var imageVersionSavingCalamity error
					BeforeEach(func() {
						imageVersionSavingCalamity = errors.New("hang in there bud")
						fakeImageFetchingDelegate.ImageVersionDeterminedReturns(imageVersionSavingCalamity)
					})

					It("returns the error", func() {
						Expect(fetchErr).To(Equal(imageVersionSavingCalamity))
					})

					It("does not construct the 'get' resource", func() {
						Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(0))
					})
				})
			})

			Context("when check returns no versions", func() {
				BeforeEach(func() {
					fakeCheckResource.CheckReturns([]atc.Version{}, nil)
				})

				It("exits with ErrImageUnavailable", func() {
					Expect(fetchErr).To(Equal(image.ErrImageUnavailable))
				})

				It("does not attempt to save any versions in the database", func() {
					Expect(fakeImageFetchingDelegate.ImageVersionDeterminedCallCount()).To(Equal(0))
				})
			})

			Context("when check returns an error", func() {
				var (
					disaster error
				)

				BeforeEach(func() {
					disaster = errors.New("wah")
					fakeCheckResource.CheckReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(fetchErr).To(Equal(disaster))
				})

				It("does not construct the 'get' resource", func() {
					Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(0))
				})
			})
		})

		Context("when initializing the Check resource fails", func() {
			var (
				disaster error
			)

			BeforeEach(func() {
				disaster = errors.New("wah")
				fakeResourceFactory.NewResourceReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(fetchErr).To(Equal(disaster))
			})

			It("does not construct the 'get' resource", func() {
				Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(0))
			})
		})
	})

	Context("when a version is specified", func() {
		BeforeEach(func() {
			version = atc.Version{"some": "version"}
		})

		Context("when saving the version in the database succeeds", func() {
			BeforeEach(func() {
				fakeImageFetchingDelegate.ImageVersionDeterminedReturns(nil)
			})

			Context("when fetching resource fails", func() {
				BeforeEach(func() {
					fakeResourceFetcher.FetchReturns(nil, resource.ErrInterrupted)
				})

				It("returns error", func() {
					Expect(fetchErr).To(Equal(resource.ErrInterrupted))
				})
			})

			Context("when fetching resource succeeds", func() {
				var (
					fakeVersionedSource   *resourcefakes.FakeVersionedSource
					fakeUsedResourceCache *db.UsedResourceCache
				)

				BeforeEach(func() {
					fakeVersionedSource = new(resourcefakes.FakeVersionedSource)
					fakeResourceFetcher.FetchReturns(fakeVersionedSource, nil)

					fakeVersionedSource.StreamOutReturns(tarStreamWith("some-tar-contents"), nil)
					fakeVolume := new(workerfakes.FakeVolume)
					fakeVersionedSource.VolumeReturns(fakeVolume)

					fakeUsedResourceCache = new(db.UsedResourceCache)
					fakeResourceCacheFactory.FindOrCreateResourceCacheReturns(fakeUsedResourceCache, nil)
				})

				Context("when the resource has a volume", func() {
					var (
						fakeVolume *workerfakes.FakeVolume
						volumePath string
					)

					BeforeEach(func() {
						fakeVolume = new(workerfakes.FakeVolume)
						volumePath = "C:/Documents and Settings/Evan/My Documents"

						fakeVolume.PathReturns(volumePath)
						fakeVersionedSource.VolumeReturns(fakeVolume)

						privileged = true
					})

					It("does not construct a new resource for checking", func() {
						Expect(fakeResourceFactory.NewResourceCallCount()).To(BeZero())
					})

					It("succeeds", func() {
						Expect(fetchErr).To(BeNil())
					})

					It("returns the image volume", func() {
						Expect(fetchedVolume).To(Equal(fakeVolume))
					})

					It("calls StreamOut on the versioned source with the right metadata path", func() {
						Expect(fakeVersionedSource.StreamOutCallCount()).To(Equal(1))
						Expect(fakeVersionedSource.StreamOutArgsForCall(0)).To(Equal("metadata.json"))
					})

					It("returns a tar stream containing the contents of metadata.json", func() {
						Expect(ioutil.ReadAll(fetchedMetadataReader)).To(Equal([]byte("some-tar-contents")))
					})

					It("has the version on the image", func() {
						Expect(fetchedVersion).To(Equal(atc.Version{"some": "version"}))
					})

					It("saved the image resource version in the database", func() {
						Expect(fakeImageFetchingDelegate.ImageVersionDeterminedCallCount()).To(Equal(1))
						Expect(fakeImageFetchingDelegate.ImageVersionDeterminedArgsForCall(0)).To(Equal(fakeUsedResourceCache))
					})

					It("fetches resource with correct session", func() {
						Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(1))
						_, session, tags, actualTeamID, actualCustomTypes, resourceInstance, metadata, delegate, _, _ := fakeResourceFetcher.FetchArgsForCall(0)
						Expect(metadata).To(Equal(resource.EmptyMetadata{}))
						Expect(session).To(Equal(resource.Session{
							Metadata: db.ContainerMetadata{
								Type: db.ContainerTypeGet,
							},
						}))
						Expect(tags).To(Equal(atc.Tags{"worker", "tags"}))
						Expect(actualTeamID).To(Equal(teamID))
						Expect(resourceInstance).To(Equal(resource.NewResourceInstance(
							"docker",
							atc.Version{"some": "version"},
							atc.Source{"some": "super-secret-sauce"},
							atc.Params{"some": "params"},
							customTypes,
							&db.UsedResourceCache{},
							db.NewImageGetContainerOwner(fakeCreatingContainer),
						)))
						Expect(actualCustomTypes).To(Equal(customTypes))
						Expect(delegate).To(Equal(fakeImageFetchingDelegate))
					})

					It("gets the volume", func() {
						Expect(fakeVersionedSource.VolumeCallCount()).To(Equal(1))
					})

					Context("when streaming the metadata out fails", func() {
						disaster := errors.New("nope")

						BeforeEach(func() {
							fakeVersionedSource.StreamOutReturns(nil, disaster)
						})

						It("returns the error", func() {
							Expect(fetchErr).To(Equal(disaster))
						})
					})

					Context("when the resource still does not have a volume for some reason", func() {
						BeforeEach(func() {
							fakeVersionedSource.VolumeReturns(nil)
						})

						It("returns an appropriate error", func() {
							Expect(fetchErr).To(Equal(image.ErrImageGetDidNotProduceVolume))
						})
					})
				})
			})
		})

		Context("when saving the version in the database fails", func() {
			var imageVersionSavingCalamity error
			BeforeEach(func() {
				imageVersionSavingCalamity = errors.New("hang in there bud")
				fakeImageFetchingDelegate.ImageVersionDeterminedReturns(imageVersionSavingCalamity)
			})

			It("returns the error", func() {
				Expect(fetchErr).To(Equal(imageVersionSavingCalamity))
			})

			It("does not construct the 'get' resource", func() {
				Expect(fakeResourceFetcher.FetchCallCount()).To(Equal(0))
			})
		})
	})
})

func tarStreamWith(metadata string) io.ReadCloser {
	buffer := gbytes.NewBuffer()

	tarWriter := tar.NewWriter(buffer)
	err := tarWriter.WriteHeader(&tar.Header{
		Name: "metadata.json",
		Mode: 0600,
		Size: int64(len(metadata)),
	})
	Expect(err).NotTo(HaveOccurred())

	_, err = io.WriteString(tarWriter, metadata)
	Expect(err).NotTo(HaveOccurred())

	err = tarWriter.Close()
	Expect(err).NotTo(HaveOccurred())

	return buffer
}
