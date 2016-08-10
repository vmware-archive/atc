package api_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
)

var _ = Describe("Versions API", func() {
	var pipelineDB *dbfakes.FakePipelineDB
	var expectedSavedPipeline db.SavedPipeline

	BeforeEach(func() {
		pipelineDB = new(dbfakes.FakePipelineDB)
		pipelineDBFactory.BuildReturns(pipelineDB)
		expectedSavedPipeline = db.SavedPipeline{}
		teamDB.GetPipelineByNameReturns(expectedSavedPipeline, true, nil)
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/versions", func() {
		var response *http.Response
		var queryParams string

		BeforeEach(func() {
			queryParams = ""
		})

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("GET", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/some-resource/versions"+queryParams, nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when not authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", 0, false, false)
			})

			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(false)
				})

				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			Context("and the pipeline is public", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(true)
					pipelineDB.GetResourceVersionsReturns([]db.SavedVersionedResource{}, db.Pagination{}, true, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", 1, true, true)
			})

			Context("when no params are passed", func() {
				It("does not set defaults for since and until", func() {
					Expect(pipelineDB.GetResourceVersionsCallCount()).To(Equal(1))

					resourceName, page := pipelineDB.GetResourceVersionsArgsForCall(0)
					Expect(resourceName).To(Equal("some-resource"))
					Expect(page).To(Equal(db.Page{
						Since: 0,
						Until: 0,
						Limit: 100,
					}))
				})
			})

			Context("when all the params are passed", func() {
				BeforeEach(func() {
					queryParams = "?since=2&until=3&limit=8"
				})

				It("passes them through", func() {
					Expect(pipelineDB.GetResourceVersionsCallCount()).To(Equal(1))

					resourceName, page := pipelineDB.GetResourceVersionsArgsForCall(0)
					Expect(resourceName).To(Equal("some-resource"))
					Expect(page).To(Equal(db.Page{
						Since: 2,
						Until: 3,
						Limit: 8,
					}))
				})
			})

			Context("when getting the versions succeeds", func() {
				var returnedVersions []db.SavedVersionedResource

				BeforeEach(func() {
					queryParams = "?since=5&limit=2"
					returnedVersions = []db.SavedVersionedResource{
						{
							ID:      4,
							Enabled: true,
							VersionedResource: db.VersionedResource{
								Resource: "some-resource",
								Type:     "some-type",
								Version: db.Version{
									"some": "version",
								},
								Metadata: []db.MetadataField{
									{
										Name:  "some",
										Value: "metadata",
									},
								},
								PipelineID: 42,
							},
						},
						{
							ID:      2,
							Enabled: false,
							VersionedResource: db.VersionedResource{
								Resource: "some-resource",
								Type:     "some-type",
								Version: db.Version{
									"some": "version",
								},
								Metadata: []db.MetadataField{
									{
										Name:  "some",
										Value: "metadata",
									},
								},
								PipelineID: 42,
							},
						},
					}

					pipelineDB.GetResourceVersionsReturns(returnedVersions, db.Pagination{}, true, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns content type application/json", func() {
					Expect(response.Header.Get("Content-type")).To(Equal("application/json"))
				})

				It("returns the json", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[
					{
						"id": 4,
						"enabled": true,
						"pipeline_id": 42,
						"resource": "some-resource",
						"type": "some-type",
						"version": {"some":"version"},
						"metadata": [
							{
								"name":"some",
								"value":"metadata"
							}
						]
					},
					{
						"id":2,
						"enabled": false,
						"pipeline_id": 42,
						"resource": "some-resource",
						"type": "some-type",
						"version": {"some":"version"},
						"metadata": [
							{
								"name":"some",
								"value":"metadata"
							}
						]
					}
				]`))
				})

				Context("when next/previous pages are available", func() {
					BeforeEach(func() {
						pipelineDB.GetPipelineNameReturns("some-pipeline")
						pipelineDB.GetResourceVersionsReturns(returnedVersions, db.Pagination{
							Previous: &db.Page{Until: 4, Limit: 2},
							Next:     &db.Page{Since: 2, Limit: 2},
						}, true, nil)
					})

					It("returns Link headers per rfc5988", func() {
						Expect(response.Header["Link"]).To(ConsistOf([]string{
							fmt.Sprintf(`<%s/api/v1/teams/a-team/pipelines/some-pipeline/resources/some-resource/versions?until=4&limit=2>; rel="previous"`, externalURL),
							fmt.Sprintf(`<%s/api/v1/teams/a-team/pipelines/some-pipeline/resources/some-resource/versions?since=2&limit=2>; rel="next"`, externalURL),
						}))
					})
				})
			})

			Context("when the versions can't be found", func() {
				BeforeEach(func() {
					pipelineDB.GetResourceVersionsReturns(nil, db.Pagination{}, false, nil)
				})

				It("returns 404 not found", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("when getting the versions fails", func() {
				BeforeEach(func() {
					pipelineDB.GetResourceVersionsReturns(nil, db.Pagination{}, false, errors.New("oh no!"))
				})

				It("returns 500 Internal Server Error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/versions/:resource_version_id/enable", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/versions/42/enable", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())

		})

		Context("when authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", 42, true, true)
			})

			It("injects the proper pipelineDB", func() {
				Expect(teamDB.GetPipelineByNameArgsForCall(0)).To(Equal("a-pipeline"))
				Expect(pipelineDBFactory.BuildCallCount()).To(Equal(1))
				actualSavedPipeline := pipelineDBFactory.BuildArgsForCall(0)
				Expect(actualSavedPipeline).To(Equal(expectedSavedPipeline))
			})

			Context("when enabling the resource succeeds", func() {
				BeforeEach(func() {
					pipelineDB.EnableVersionedResourceReturns(nil)
				})

				It("enabled the right versioned resource", func() {
					Expect(pipelineDB.EnableVersionedResourceArgsForCall(0)).To(Equal(42))
				})

				It("returns 200", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})

			Context("when enabling the resource fails", func() {
				BeforeEach(func() {
					pipelineDB.EnableVersionedResourceReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})

		Context("when not authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/versions/:resource_version_id/disable", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/resource-name/versions/42/disable", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", 42, true, true)
			})

			It("injects the proper pipelineDB", func() {
				Expect(teamDB.GetPipelineByNameCallCount()).To(Equal(1))
				Expect(teamDB.GetPipelineByNameArgsForCall(0)).To(Equal("a-pipeline"))
				Expect(pipelineDBFactory.BuildCallCount()).To(Equal(1))
				actualSavedPipeline := pipelineDBFactory.BuildArgsForCall(0)
				Expect(actualSavedPipeline).To(Equal(expectedSavedPipeline))
			})

			Context("when enabling the resource succeeds", func() {
				BeforeEach(func() {
					pipelineDB.DisableVersionedResourceReturns(nil)
				})

				It("disabled the right versioned resource", func() {
					Expect(pipelineDB.DisableVersionedResourceArgsForCall(0)).To(Equal(42))
				})

				It("returns 200", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})

			Context("when enabling the resource fails", func() {
				BeforeEach(func() {
					pipelineDB.DisableVersionedResourceReturns(errors.New("welp"))
				})

				It("returns 500", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})

		Context("when not authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/versions/:resource_version_id/input_to", func() {
		var response *http.Response
		var stringVersionID string

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("GET", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/some-resource/versions/"+stringVersionID+"/input_to", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			stringVersionID = "123"
		})

		Context("when not authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", 0, false, false)
			})

			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(false)
				})

				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			Context("and the pipeline is public", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(true)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", 1, true, true)
			})

			It("looks up the given version ID", func() {
				Expect(pipelineDB.GetBuildsWithVersionAsInputCallCount()).To(Equal(1))
				Expect(pipelineDB.GetBuildsWithVersionAsInputArgsForCall(0)).To(Equal(123))
			})

			Context("when getting the builds succeeds", func() {
				BeforeEach(func() {
					build1 := new(dbfakes.FakeBuild)
					build1.IDReturns(1024)
					build1.NameReturns("5")
					build1.JobNameReturns("some-job")
					build1.PipelineNameReturns("a-pipeline")
					build1.TeamNameReturns("a-team")
					build1.StatusReturns(db.StatusSucceeded)
					build1.StartTimeReturns(time.Unix(1, 0))
					build1.EndTimeReturns(time.Unix(100, 0))

					build2 := new(dbfakes.FakeBuild)
					build2.IDReturns(1025)
					build2.NameReturns("6")
					build2.JobNameReturns("some-job")
					build2.PipelineNameReturns("a-pipeline")
					build2.TeamNameReturns("a-team")
					build2.StatusReturns(db.StatusSucceeded)
					build2.StartTimeReturns(time.Unix(200, 0))
					build2.EndTimeReturns(time.Unix(300, 0))

					pipelineDB.GetBuildsWithVersionAsInputReturns([]db.Build{build1, build2}, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns content type application/json", func() {
					Expect(response.Header.Get("Content-type")).To(Equal("application/json"))
				})

				It("returns the json", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[
					{
						"id": 1024,
						"team_name": "a-team",
						"name": "5",
						"status": "succeeded",
						"job_name": "some-job",
						"url": "/teams/a-team/pipelines/a-pipeline/jobs/some-job/builds/5",
						"api_url": "/api/v1/builds/1024",
						"pipeline_name": "a-pipeline",
						"start_time": 1,
						"end_time": 100
					},
					{
						"id": 1025,
						"name": "6",
						"team_name": "a-team",
						"status": "succeeded",
						"job_name": "some-job",
						"url": "/teams/a-team/pipelines/a-pipeline/jobs/some-job/builds/6",
						"api_url": "/api/v1/builds/1025",
						"pipeline_name": "a-pipeline",
						"start_time": 200,
						"end_time": 300
					}
				]`))
				})
			})

			Context("when the version ID is invalid", func() {
				BeforeEach(func() {
					stringVersionID = "hello"
				})

				It("returns an empty list", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[]`))
				})
			})

			Context("when the call to get builds returns an error", func() {
				BeforeEach(func() {
					pipelineDB.GetBuildsWithVersionAsInputReturns(nil, errors.New("NOPE"))
				})

				It("returns a 500 internal server error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/resources/:resource_name/versions/:resource_version_id/output_of", func() {
		var response *http.Response
		var stringVersionID string

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("GET", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/resources/some-resource/versions/"+stringVersionID+"/output_of", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			stringVersionID = "123"
		})

		Context("when not authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", 0, false, false)
			})

			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(false)
				})

				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			Context("and the pipeline is public", func() {
				BeforeEach(func() {
					pipelineDB.IsPublicReturns(true)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})

		Context("when authorized", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", 1, true, true)
			})

			It("looks up the given version ID", func() {
				Expect(pipelineDB.GetBuildsWithVersionAsOutputCallCount()).To(Equal(1))
				Expect(pipelineDB.GetBuildsWithVersionAsOutputArgsForCall(0)).To(Equal(123))
			})

			Context("when getting the builds succeeds", func() {
				BeforeEach(func() {
					build1 := new(dbfakes.FakeBuild)
					build1.IDReturns(1024)
					build1.NameReturns("5")
					build1.JobNameReturns("some-job")
					build1.PipelineNameReturns("a-pipeline")
					build1.TeamNameReturns("a-team")
					build1.StatusReturns(db.StatusSucceeded)
					build1.StartTimeReturns(time.Unix(1, 0))
					build1.EndTimeReturns(time.Unix(100, 0))

					build2 := new(dbfakes.FakeBuild)
					build2.IDReturns(1025)
					build2.NameReturns("6")
					build2.JobNameReturns("some-job")
					build2.PipelineNameReturns("a-pipeline")
					build2.TeamNameReturns("a-team")
					build2.StatusReturns(db.StatusSucceeded)
					build2.StartTimeReturns(time.Unix(200, 0))
					build2.EndTimeReturns(time.Unix(300, 0))

					pipelineDB.GetBuildsWithVersionAsOutputReturns([]db.Build{build1, build2}, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns content type application/json", func() {
					Expect(response.Header.Get("Content-type")).To(Equal("application/json"))
				})

				It("returns the json", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[
					{
						"id": 1024,
						"name": "5",
						"status": "succeeded",
						"job_name": "some-job",
						"url": "/teams/a-team/pipelines/a-pipeline/jobs/some-job/builds/5",
						"api_url": "/api/v1/builds/1024",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"start_time": 1,
						"end_time": 100
					},
					{
						"id": 1025,
						"name": "6",
						"status": "succeeded",
						"job_name": "some-job",
						"url": "/teams/a-team/pipelines/a-pipeline/jobs/some-job/builds/6",
						"api_url": "/api/v1/builds/1025",
						"pipeline_name": "a-pipeline",
						"team_name": "a-team",
						"start_time": 200,
						"end_time": 300
					}
				]`))
				})
			})

			Context("when the version ID is invalid", func() {
				BeforeEach(func() {
					stringVersionID = "hello"
				})

				It("returns an empty list", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`[]`))
				})
			})

			Context("when the call to get builds returns an error", func() {
				BeforeEach(func() {
					pipelineDB.GetBuildsWithVersionAsOutputReturns(nil, errors.New("NOPE"))
				})

				It("returns a 500 internal server error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})
