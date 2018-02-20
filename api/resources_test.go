package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	"github.com/concourse/atc/radar/radarfakes"
	"github.com/concourse/atc/resource"
)

var _ = Describe("Resources API", func() {
	var (
		fakePipeline *dbfakes.FakePipeline
		resource1    *dbfakes.FakeResource
		resource2    *dbfakes.FakeResource
		resource3    *dbfakes.FakeResource
	)

	BeforeEach(func() {
		fakePipeline = new(dbfakes.FakePipeline)
		dbTeamFactory.FindTeamReturns(dbTeam, true, nil)
		dbTeam.PipelineReturns(fakePipeline, true, nil)
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/resources", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			response, err = client.Get(server.URL + "/api/v1/teams/a-team/pipelines/a-pipeline/resources")
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when getting the dashboard resources succeeds", func() {
			BeforeEach(func() {
				resource1 = new(dbfakes.FakeResource)
				resource1.IDReturns(1)
				resource1.CheckErrorReturns(nil)
				resource1.PausedReturns(true)
				resource1.PipelineNameReturns("a-pipeline")
				resource1.NameReturns("resource-1")
				resource1.TypeReturns("type-1")
				resource1.LastCheckedReturns(time.Unix(1513364881, 0))

				resource2 = new(dbfakes.FakeResource)
				resource2.IDReturns(2)
				resource2.CheckErrorReturns(errors.New("sup"))
				resource2.FailingToCheckReturns(true)
				resource2.PausedReturns(false)
				resource2.PipelineNameReturns("a-pipeline")
				resource2.NameReturns("resource-2")
				resource2.TypeReturns("type-2")

				resource3 = new(dbfakes.FakeResource)
				resource3.IDReturns(3)
				resource3.CheckErrorReturns(nil)
				resource3.PausedReturns(true)
				resource3.PipelineNameReturns("a-pipeline")
				resource3.NameReturns("resource-3")
				resource3.TypeReturns("type-3")
			})

			Context("when not authorized", func() {
				Context("and the pipeline is private", func() {
					BeforeEach(func() {
						fakeAccessor.TeamPipelineResourcesReturns(nil, accessor.ErrNotAuthorized)
					})

					It("returns 401", func() {
						Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
					})
				})

				Context("and the pipeline is public", func() {
					BeforeEach(func() {
						resource1.CheckErrorReturns(nil)
						resource2.CheckErrorReturns(nil)
						resource3.CheckErrorReturns(nil)

						fakeAccessor.TeamPipelineResourcesReturns([]db.Resource{resource1, resource2, resource3}, nil)
					})

					It("returns 200 OK", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})

					It("returns each resource, excluding their check failure", func() {
						body, err := ioutil.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						Expect(body).To(MatchJSON(`[
					{
						"name": "resource-1",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"type": "type-1",
						"paused": true,
						"last_checked": 1513364881
					},
					{
						"name": "resource-2",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"type": "type-2",
						"failing_to_check": true
					},
					{
						"name": "resource-3",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"type": "type-3",
						"paused": true
					}
				]`))
					})
				})
			})

			Context("when authorized", func() {
				BeforeEach(func() {
					resource2.CheckErrorReturns(errors.New("sup"))
					fakeAccessor.TeamPipelineResourcesReturns([]db.Resource{resource1, resource2, resource3}, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns each resource, including their check failure", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[
						{
							"name": "resource-1",
							"pipeline_name": "a-pipeline",
							"team_name": "a-team",
							"type": "type-1",
							"paused": true,
							"last_checked": 1513364881
						},
						{
							"name": "resource-2",
							"pipeline_name": "a-pipeline",
							"team_name": "a-team",
							"type": "type-2",
							"failing_to_check": true,
							"check_error": "sup"
						},
						{
							"name": "resource-3",
							"pipeline_name": "a-pipeline",
							"team_name": "a-team",
							"type": "type-3",
							"paused": true
						}
					]`))
				})

				Context("when getting the resource config fails", func() {
					Context("when the resources are not found", func() {
						BeforeEach(func() {
							fakeAccessor.TeamPipelineResourcesReturns(nil, accessor.ErrNotFound)
						})

						It("returns 404", func() {
							Expect(response.StatusCode).To(Equal(http.StatusNotFound))
						})
					})

					Context("with an unknown error", func() {
						BeforeEach(func() {
							fakeAccessor.TeamPipelineResourcesReturns(nil, errors.New("error"))
						})

						It("returns 500", func() {
							Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
						})
					})
				})
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			response, err = client.Get(server.URL + "/api/v1/teams/a-team/pipelines/a-pipeline/resources/some-resource")
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when not authorized", func() {
			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotAuthorized)
				})

				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})
		})

		Context("when the pipeline is public", func() {
			BeforeEach(func() {
				resource1 := new(dbfakes.FakeResource)
				resource1.CheckErrorReturns(nil)
				resource1.PipelineNameReturns("a-pipeline")
				resource1.NameReturns("resource-1")
				resource1.FailingToCheckReturns(true)
				resource1.TypeReturns("type-1")
				resource1.LastCheckedReturns(time.Unix(1513364881, 0))

				fakeAccessor.TeamPipelineResourceReturns(resource1, nil)
			})

			It("returns 200 OK", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns the resource json without the check error", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`
					{
						"name": "resource-1",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"type": "type-1",
						"last_checked": 1513364881,
						"failing_to_check": true
					}`))
			})
		})

		Context("when authorized", func() {

			BeforeEach(func() {
				resource1 := new(dbfakes.FakeResource)
				resource1.NameReturns("resource-1")

				fakeAccessor.TeamPipelineResourceReturns(resource1, nil)
			})

			It("looks it up in the database", func() {
				Expect(fakeAccessor.TeamPipelineResourceCallCount()).To(Equal(1))

				access, teamName, pipelineName, resourceName := fakeAccessor.TeamPipelineResourceArgsForCall(0)
				Expect(teamName).To(Equal("a-team"))
				Expect(pipelineName).To(Equal("a-pipeline"))
				Expect(resourceName).To(Equal("some-resource"))
				Expect(access).To(Equal(accessor.Read))
			})

			Context("when the resource cannot be found in the database", func() {
				BeforeEach(func() {
					fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotFound)
				})

				It("returns a 404", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when the call to the db returns an error", func() {
				BeforeEach(func() {
					fakeAccessor.TeamPipelineResourceReturns(nil, errors.New("error"))
				})

				It("returns a 500 error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when the call to get a resource succeeds", func() {
				BeforeEach(func() {
					resource1 := new(dbfakes.FakeResource)
					resource1.CheckErrorReturns(errors.New("sup"))
					resource1.PausedReturns(true)
					resource1.PipelineNameReturns("a-pipeline")
					resource1.NameReturns("resource-1")
					resource1.FailingToCheckReturns(true)
					resource1.TypeReturns("type-1")
					resource1.LastCheckedReturns(time.Unix(1513364881, 0))

					fakeAccessor.TeamPipelineResourceReturns(resource1, nil)
				})

				It("returns 200 ok", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns the resource json with the check error", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`
							{
								"name": "resource-1",
								"pipeline_name": "a-pipeline",
								"team_name": "a-team",
								"type": "type-1",
								"last_checked": 1513364881,
								"paused": true,
								"failing_to_check": true,
								"check_error": "sup"
							}`))
				})
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/pause", func() {
		var (
			response     *http.Response
			fakeResource *dbfakes.FakeResource
		)

		BeforeEach(func() {
			fakeResource = new(dbfakes.FakeResource)
			fakeResource.NameReturns("resource-name")
		})

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/pause", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineResourceReturns(fakeResource, nil)
			})

			Context("when pausing the resource succeeds", func() {
				BeforeEach(func() {
					fakeResource.PauseReturns(nil)
				})

				It("returns 200", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})

			Context("when resource can not be found", func() {
				BeforeEach(func() {
					fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotFound)
				})

				It("returns 404", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when pausing the resource fails", func() {
				BeforeEach(func() {
					fakeResource.PauseReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotAuthorized)
			})

			It("returns Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/unpause", func() {
		var (
			response     *http.Response
			fakeResource *dbfakes.FakeResource
		)

		BeforeEach(func() {
			fakeResource = new(dbfakes.FakeResource)
			fakeResource.NameReturns("resource-name")
		})

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/unpause", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineResourceReturns(fakeResource, nil)
			})

			Context("when unpausing the resource succeeds", func() {
				BeforeEach(func() {
					fakeResource.UnpauseReturns(nil)
				})

				It("returns 200", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})

			Context("when resource can not be found", func() {
				BeforeEach(func() {
					fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotFound)
				})

				It("returns 404", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when unpausing the resource fails", func() {
				BeforeEach(func() {
					fakeResource.UnpauseReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineResourceReturns(nil, accessor.ErrNotAuthorized)
			})

			It("returns Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("POST /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/check", func() {
		var fakeScanner *radarfakes.FakeScanner
		var checkRequestBody atc.CheckRequestBody
		var response *http.Response

		BeforeEach(func() {
			fakeScanner = new(radarfakes.FakeScanner)
			fakeScannerFactory.NewResourceScannerReturns(fakeScanner)

			checkRequestBody = atc.CheckRequestBody{}
		})

		JustBeforeEach(func() {
			reqPayload, err := json.Marshal(checkRequestBody)
			Expect(err).NotTo(HaveOccurred())

			request, err := http.NewRequest("POST", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/check", bytes.NewBuffer(reqPayload))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineReturns(fakePipeline, nil)
			})

			It("requests the correct pipeline", func() {
				Expect(fakeAccessor.TeamPipelineCallCount()).To(Equal(1))
				access, teamName, pipelineName := fakeAccessor.TeamPipelineArgsForCall(0)
				Expect(teamName).To(Equal("a-team"))
				Expect(pipelineName).To(Equal("a-pipeline"))
				Expect(access).To(Equal(accessor.Write))
			})

			It("tries to scan with no version specified", func() {
				Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
				_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
				Expect(actualResourceName).To(Equal("resource-name"))
				Expect(actualFromVersion).To(BeNil())
			})

			It("returns 200", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			Context("when checking with a version specified", func() {
				BeforeEach(func() {
					checkRequestBody = atc.CheckRequestBody{
						From: atc.Version{
							"some-version-key": "some-version-value",
						},
					}
				})

				It("tries to scan with the version specified", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
					_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
					Expect(actualResourceName).To(Equal("resource-name"))
					Expect(actualFromVersion).To(Equal(checkRequestBody.From))
				})
			})

			Context("when the resource already has versions", func() {
				BeforeEach(func() {
					returnedVersion := db.SavedVersionedResource{
						ID:      4,
						Enabled: true,
						VersionedResource: db.VersionedResource{
							Resource: "some-resource",
							Type:     "some-type",
							Version: db.ResourceVersion{
								"some": "version",
							},
							Metadata: []db.ResourceMetadataField{
								{
									Name:  "some",
									Value: "metadata",
								},
							},
						},
					}
					fakePipeline.GetLatestVersionedResourceReturns(returnedVersion, true, nil)
				})

				It("tries to scan with the latest version when no version is passed", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
					_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
					Expect(actualResourceName).To(Equal("resource-name"))
					Expect(actualFromVersion).To(Equal(atc.Version{"some": "version"}))
				})
			})

			Context("when failing to get latest version for resource", func() {
				BeforeEach(func() {
					fakePipeline.GetLatestVersionedResourceReturns(db.SavedVersionedResource{}, false, errors.New("disaster"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("does not scan from version", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(0))
				})
			})

			Context("when checking fails with ResourceNotFoundError", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(db.ResourceNotFoundError{})
				})

				It("returns 404", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when checking the resource fails internally", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when checking the resource fails with ErrResourceScriptFailed", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(
						resource.ErrResourceScriptFailed{
							ExitStatus: 42,
							Stderr:     "my tooth",
						},
					)
				})

				It("returns 400", func() {
					Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				})

				It("returns the script's exit status and stderr", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`{
						"exit_status": 42,
						"stderr": "my tooth"
					}`))
				})

				It("returns application/json", func() {
					Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineReturns(nil, accessor.ErrNotAuthorized)
			})

			It("returns 401", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("POST /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/check/webhook", func() {
		var (
			fakeScanner      *radarfakes.FakeScanner
			checkRequestBody atc.CheckRequestBody
			response         *http.Response
			versionMap       atc.Version
			fakeResource     *dbfakes.FakeResource
		)

		BeforeEach(func() {
			fakeScanner = new(radarfakes.FakeScanner)
			fakeScannerFactory.NewResourceScannerReturns(fakeScanner)
			checkRequestBody = atc.CheckRequestBody{}

			fakeResource = new(dbfakes.FakeResource)
			fakeResource.NameReturns("resource-name")
		})

		JustBeforeEach(func() {
			reqPayload, err := json.Marshal(checkRequestBody)
			Expect(err).NotTo(HaveOccurred())

			request, err := http.NewRequest("POST", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/check/webhook?webhook_token=fake-token", bytes.NewBuffer(reqPayload))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Set("Content-Type", "application/json")

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				fakeAccessor.TeamPipelineReturns(fakePipeline, nil)
				fakeResource.WebhookTokenReturns("fake-token")
				fakePipeline.ResourceReturns(fakeResource, true, nil)
			})

			It("requests pipeline with propers args", func() {
				Expect(fakeAccessor.TeamPipelineCallCount()).To(Equal(1))
				access, teamName, pipelineName := fakeAccessor.TeamPipelineArgsForCall(0)
				Expect(access).To(Equal(accessor.Skip))
				Expect(teamName).To(Equal("a-team"))
				Expect(pipelineName).To(Equal("a-pipeline"))
			})

			It("requests resources with propers args", func() {
				Expect(fakePipeline.ResourceCallCount()).To(Equal(1))
				resourceName := fakePipeline.ResourceArgsForCall(0)
				Expect(resourceName).To(Equal("resource-name"))
			})

			It("tries to scan with no version specified", func() {
				Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
				_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
				Expect(actualResourceName).To(Equal("resource-name"))
				Expect(actualFromVersion).To(BeNil())
			})

			It("returns 200", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			Context("when checking with a version specified", func() {
				BeforeEach(func() {
					resourceVersion := map[string]string{"some-version": "some-key"}
					versionMap = atc.Version(resourceVersion)
					fakePipeline.GetLatestVersionedResourceReturns(
						db.SavedVersionedResource{VersionedResource: db.VersionedResource{Version: resourceVersion}}, true, nil)
				})

				It("tries to scan with the version specified", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
					_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
					Expect(actualResourceName).To(Equal("resource-name"))
					Expect(actualFromVersion).To(Equal(versionMap))
				})
			})

			Context("when the resource already has versions", func() {
				BeforeEach(func() {
					returnedVersion := db.SavedVersionedResource{
						ID:      4,
						Enabled: true,
						VersionedResource: db.VersionedResource{
							Resource: "some-resource",
							Type:     "some-type",
							Version: db.ResourceVersion{
								"some": "version",
							},
							Metadata: []db.ResourceMetadataField{
								{
									Name:  "some",
									Value: "metadata",
								},
							},
						},
					}
					fakePipeline.GetLatestVersionedResourceReturns(returnedVersion, true, nil)
				})

				It("tries to scan with the latest version when no version is passed", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(1))
					_, actualResourceName, actualFromVersion := fakeScanner.ScanFromVersionArgsForCall(0)
					Expect(actualResourceName).To(Equal("resource-name"))
					Expect(actualFromVersion).To(Equal(atc.Version{"some": "version"}))
				})
			})

			Context("when failing to get latest version for resource", func() {
				BeforeEach(func() {
					fakePipeline.GetLatestVersionedResourceReturns(db.SavedVersionedResource{}, false, errors.New("disaster"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("does not scan from version", func() {
					Expect(fakeScanner.ScanFromVersionCallCount()).To(Equal(0))
				})
			})

			Context("when checking fails with ResourceNotFoundError", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(db.ResourceNotFoundError{})
				})

				It("returns 404", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when checking the resource fails internally", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when checking the resource fails with err", func() {
				BeforeEach(func() {
					fakeScanner.ScanFromVersionReturns(errors.New("error"))
				})

				It("returns 400", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
			Context("when unauthorized", func() {
				BeforeEach(func() {
					fakeResource.WebhookTokenReturns("wrong-token")
					fakePipeline.ResourceReturns(fakeResource, true, nil)
				})
				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})
		})
	})
})
