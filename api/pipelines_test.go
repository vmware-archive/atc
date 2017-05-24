package api_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/dbng/dbngfakes"

	"github.com/concourse/atc/db/algorithm"
)

var _ = Describe("Pipelines API", func() {
	var (
		dbPipeline *dbngfakes.FakePipeline
		fakeTeam   *dbngfakes.FakeTeam

		publicPipeline        *dbngfakes.FakePipeline
		anotherPublicPipeline *dbngfakes.FakePipeline
		privatePipeline       *dbngfakes.FakePipeline
	)
	BeforeEach(func() {
		dbPipeline = new(dbngfakes.FakePipeline)
		fakeTeam = new(dbngfakes.FakeTeam)

		publicPipeline = new(dbngfakes.FakePipeline)
		publicPipeline.IDReturns(1)
		publicPipeline.PausedReturns(true)
		publicPipeline.PublicReturns(true)
		publicPipeline.TeamNameReturns("main")
		publicPipeline.NameReturns("public-pipeline")
		publicPipeline.ConfigReturns(atc.Config{
			Groups: atc.GroupConfigs{
				{
					Name:      "group2",
					Jobs:      []string{"job3", "job4"},
					Resources: []string{"resource3", "resource4"},
				},
			},
		}, "", 0, nil)

		anotherPublicPipeline = new(dbngfakes.FakePipeline)
		anotherPublicPipeline.IDReturns(2)
		anotherPublicPipeline.PausedReturns(true)
		anotherPublicPipeline.PublicReturns(true)
		anotherPublicPipeline.TeamNameReturns("another")
		anotherPublicPipeline.NameReturns("another-pipeline")

		privatePipeline = new(dbngfakes.FakePipeline)
		privatePipeline.IDReturns(3)
		privatePipeline.PausedReturns(false)
		privatePipeline.PublicReturns(false)
		privatePipeline.TeamNameReturns("main")
		privatePipeline.NameReturns("private-pipeline")
		privatePipeline.ConfigReturns(atc.Config{
			Groups: atc.GroupConfigs{
				{
					Name:      "group1",
					Jobs:      []string{"job1", "job2"},
					Resources: []string{"resource1", "resource2"},
				},
			},
		}, "", 0, nil)

		fakeTeam.PipelinesReturns([]dbng.Pipeline{
			privatePipeline,
			publicPipeline,
		}, nil)

		fakeTeam.VisiblePipelinesReturns([]dbng.Pipeline{
			privatePipeline,
			publicPipeline,
			anotherPublicPipeline,
		}, nil)
		fakeTeam.PublicPipelinesReturns([]dbng.Pipeline{publicPipeline}, nil)

		dbPipelineFactory.PublicPipelinesReturns([]dbng.Pipeline{publicPipeline, anotherPublicPipeline}, nil)
	})

	Describe("GET /api/v1/pipelines", func() {
		var response *http.Response

		JustBeforeEach(func() {
			req, err := http.NewRequest("GET", server.URL+"/api/v1/pipelines", nil)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 OK", func() {
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns application/json", func() {
			Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
		})

		Context("when team is set in user context", func() {
			BeforeEach(func() {
				userContextReader.GetTeamReturns("some-team", false, true)
			})

			It("constructs teamDB with provided team name", func() {
				Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
				Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("some-team"))
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				userContextReader.GetTeamReturns("", false, false)
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns only public pipelines", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
				{
					"id": 1,
					"name": "public-pipeline",
					"url": "/teams/main/pipelines/public-pipeline",
					"paused": true,
					"public": true,
					"team_name": "main",
					"groups": [
						{
							"name": "group2",
							"jobs": ["job3", "job4"],
							"resources": ["resource3", "resource4"]
						}
					]
				},
				{
					"id": 2,
					"name": "another-pipeline",
					"url": "/teams/another/pipelines/another-pipeline",
					"paused": true,
					"public": true,
					"team_name": "another"
				}]`))
			})
		})

		Context("when authenticated", func() {
			BeforeEach(func() {
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				userContextReader.GetTeamReturns("main", false, true)
				authValidator.IsAuthenticatedReturns(true)
			})

			It("returns all pipelines of the team + all public pipelines", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
				{
					"id": 3,
					"name": "private-pipeline",
					"url": "/teams/main/pipelines/private-pipeline",
					"paused": false,
					"public": false,
					"team_name": "main",
					"groups": [
						{
							"name": "group1",
							"jobs": ["job1", "job2"],
							"resources": ["resource1", "resource2"]
						}
					]
				},
				{
					"id": 1,
					"name": "public-pipeline",
					"url": "/teams/main/pipelines/public-pipeline",
					"paused": true,
					"public": true,
					"team_name": "main",
					"groups": [
						{
							"name": "group2",
							"jobs": ["job3", "job4"],
							"resources": ["resource3", "resource4"]
						}
					]
				},
				{
					"id": 2,
					"name": "another-pipeline",
					"url": "/teams/another/pipelines/another-pipeline",
					"paused": true,
					"public": true,
					"team_name": "another"
				}]`))
			})

			Context("when the call to get active pipelines fails", func() {
				BeforeEach(func() {
					fakeTeam.VisiblePipelinesReturns(nil, errors.New("disaster"))
				})

				It("returns 500 internal server error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines", func() {
		var response *http.Response

		JustBeforeEach(func() {
			req, err := http.NewRequest("GET", server.URL+"/api/v1/teams/main/pipelines", nil)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated as requested team", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("main", false, true)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
			})

			It("returns 200 OK", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns application/json", func() {
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("constructs team with provided team name", func() {
				Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
				Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("main"))
			})

			It("returns all team's pipelines", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
					{
						"id": 3,
						"name": "private-pipeline",
						"url": "/teams/main/pipelines/private-pipeline",
						"paused": false,
						"public": false,
						"team_name": "main",
						"groups": [
							{
								"name": "group1",
								"jobs": ["job1", "job2"],
								"resources": ["resource1", "resource2"]
							}
						]
					},
					{
						"id": 1,
						"name": "public-pipeline",
						"url": "/teams/main/pipelines/public-pipeline",
						"paused": true,
						"public": true,
						"team_name": "main",
						"groups": [
							{
								"name": "group2",
								"jobs": ["job3", "job4"],
								"resources": ["resource3", "resource4"]
							}
						]
					}]`))
			})

			Context("when the call to get active pipelines fails", func() {
				BeforeEach(func() {
					fakeTeam.PipelinesReturns(nil, errors.New("disaster"))
				})

				It("returns 500 internal server error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})

		Context("when authenticated as another team", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("another-team", false, true)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
			})

			It("returns only team's public pipelines", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
					{
						"id": 1,
						"name": "public-pipeline",
						"url": "/teams/main/pipelines/public-pipeline",
						"paused": true,
						"public": true,
						"team_name": "main",
						"groups": [
							{
								"name": "group2",
								"jobs": ["job3", "job4"],
								"resources": ["resource3", "resource4"]
							}
						]
					}]`))
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", false, false)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
			})

			It("returns only team's public pipelines", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
					{
						"id": 1,
						"name": "public-pipeline",
						"url": "/teams/main/pipelines/public-pipeline",
						"paused": true,
						"public": true,
						"team_name": "main",
						"groups": [
							{
								"name": "group2",
								"jobs": ["job3", "job4"],
								"resources": ["resource3", "resource4"]
							}
						]
					}]`))
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name", func() {
		var response *http.Response
		var fakePipeline *dbngfakes.FakePipeline

		BeforeEach(func() {
			fakePipeline = new(dbngfakes.FakePipeline)
			fakePipeline.IDReturns(4)
			fakePipeline.NameReturns("some-specific-pipeline")
			fakePipeline.PausedReturns(false)
			fakePipeline.PublicReturns(true)
			fakePipeline.TeamNameReturns("a-team")
			fakePipeline.ConfigReturns(atc.Config{
				Groups: atc.GroupConfigs{
					{
						Name:      "group1",
						Jobs:      []string{"job1", "job2"},
						Resources: []string{"resource1", "resource2"},
					},
					{
						Name:      "group2",
						Jobs:      []string{"job3", "job4"},
						Resources: []string{"resource3", "resource4"},
					},
				},
			}, "", 0, nil)
		})

		JustBeforeEach(func() {
			req, err := http.NewRequest("GET", server.URL+"/api/v1/teams/a-team/pipelines/some-specific-pipeline", nil)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", false, false)
			})

			It("returns 401", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when authenticated as requested team", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", true, true)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				fakeTeam.PipelineReturns(fakePipeline, true, nil)
			})

			It("returns 200 ok", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns application/json", func() {
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("returns a pipeline JSON", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`
					{
						"id": 4,
						"name": "some-specific-pipeline",
						"url": "/teams/a-team/pipelines/some-specific-pipeline",
						"paused": false,
						"public": true,
						"team_name": "a-team",
						"groups": [
							{
								"name": "group1",
								"jobs": ["job1", "job2"],
								"resources": ["resource1", "resource2"]
							},
							{
								"name": "group2",
								"jobs": ["job3", "job4"],
								"resources": ["resource3", "resource4"]
							}
						]
					}`))
			})
		})

		Context("when authenticated as another team", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("another-team", true, true)
				fakeTeam.PipelineReturns(fakePipeline, true, nil)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
			})

			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					fakePipeline.PublicReturns(false)
				})

				It("returns 403", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})

			Context("and the pipeline is public", func() {
				BeforeEach(func() {
					fakeTeam.PipelineReturns(fakePipeline, true, nil)
					fakePipeline.PublicReturns(true)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})

		Context("when not authenticated at all", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
				userContextReader.GetTeamReturns("", true, false)
				dbTeam.PipelineReturns(fakePipeline, true, nil)
			})

			Context("and the pipeline is private", func() {
				BeforeEach(func() {
					fakePipeline.PublicReturns(false)
				})

				It("returns 401", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})
			})

			Context("and the pipeline is public", func() {
				BeforeEach(func() {
					fakePipeline.PublicReturns(true)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})
	})

	Describe("DELETE /api/v1/teams/:team_name/pipelines/:pipeline_name", func() {
		var response *http.Response

		JustBeforeEach(func() {
			pipelineName := "a-pipeline-name"
			req, err := http.NewRequest("DELETE", server.URL+"/api/v1/teams/a-team/pipelines/"+pipelineName, nil)
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")

			response, err = client.Do(req)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
					dbPipeline.NameReturns("a-pipeline-name")
					fakeTeam.PipelineReturns(dbPipeline, true, nil)
				})

				It("returns 204 No Content", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipelineDB", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline-name"))
				})

				It("deletes the named pipeline from the database", func() {
					Expect(dbPipeline.DestroyCallCount()).To(Equal(1))
				})

				Context("when an error occurs destroying the pipeline", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						err := errors.New("disaster!")
						dbPipeline.DestroyReturns(err)
					})

					It("returns a 500 Internal Server Error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when the user is not logged in", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/pause", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/pause", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipelineDB", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline"))
				})

				Context("when pausing the pipeline succeeds", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.PauseReturns(nil)
					})

					It("returns 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})
				})

				Context("when pausing the pipeline fails", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.PauseReturns(errors.New("welp"))
					})

					It("returns 500", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/unpause", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/unpause", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
					fakeTeam.PipelineReturns(dbPipeline, true, nil)
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipelineDB", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline"))
				})

				Context("when unpausing the pipeline succeeds", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.UnpauseReturns(nil)
					})

					It("returns 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})
				})

				Context("when unpausing the pipeline fails", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.UnpauseReturns(errors.New("welp"))
					})

					It("returns 500", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/expose", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/expose", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipelineDB", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline"))
				})

				Context("when exposing the pipeline succeeds", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.ExposeReturns(nil)
					})

					It("returns 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})
				})

				Context("when exposing the pipeline fails", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.ExposeReturns(errors.New("welp"))
					})

					It("returns 500", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/hide", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/hide", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipeline", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline"))
				})

				Context("when hiding the pipeline succeeds", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.HideReturns(nil)
					})

					It("returns 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})
				})

				Context("when hiding the pipeline fails", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.HideReturns(errors.New("welp"))
					})

					It("returns 500", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/ordering", func() {
		var response *http.Response
		var body io.Reader

		BeforeEach(func() {
			body = bytes.NewBufferString(`
				[
					"a-pipeline",
					"another-pipeline",
					"yet-another-pipeline",
					"one-final-pipeline",
					"just-kidding"
				]
			`)
		})

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/ordering", body)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belonbgs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				Context("with invalid json", func() {
					BeforeEach(func() {
						body = bytes.NewBufferString(`{}`)
					})

					It("returns 400", func() {
						Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
					})
				})

				It("constructs team with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				Context("when ordering the pipelines succeeds", func() {
					BeforeEach(func() {
						fakeTeam.OrderPipelinesReturns(nil)
					})

					It("orders the pipelines", func() {
						Expect(fakeTeam.OrderPipelinesCallCount()).To(Equal(1))
						pipelineNames := fakeTeam.OrderPipelinesArgsForCall(0)
						Expect(pipelineNames).To(Equal(
							[]string{
								"a-pipeline",
								"another-pipeline",
								"yet-another-pipeline",
								"one-final-pipeline",
								"just-kidding",
							},
						))

					})

					It("returns 200", func() {
						Expect(response.StatusCode).To(Equal(http.StatusOK))
					})
				})

				Context("when ordering the pipelines fails", func() {
					BeforeEach(func() {
						fakeTeam.OrderPipelinesReturns(errors.New("welp"))
					})

					It("returns 500", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})

			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("GET /api/v1/teams/:team_name/pipelines/:pipeline_name/versions-db", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("GET", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/versions-db", nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("a-team", true, true)
				dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				fakeTeam.PipelineReturns(dbPipeline, true, nil)
				//construct Version db

				dbPipeline.LoadVersionsDBReturns(
					&algorithm.VersionsDB{
						ResourceVersions: []algorithm.ResourceVersion{
							{
								VersionID:  73,
								ResourceID: 127,
								CheckOrder: 123,
							},
						},
						BuildOutputs: []algorithm.BuildOutput{
							{
								ResourceVersion: algorithm.ResourceVersion{
									VersionID:  73,
									ResourceID: 127,
									CheckOrder: 123,
								},
								BuildID: 66,
								JobID:   13,
							},
						},
						BuildInputs: []algorithm.BuildInput{
							{
								ResourceVersion: algorithm.ResourceVersion{
									VersionID:  66,
									ResourceID: 77,
									CheckOrder: 88,
								},
								BuildID:   66,
								JobID:     13,
								InputName: "some-input-name",
							},
						},
						JobIDs: map[string]int{
							"bad-luck-job": 13,
						},
						ResourceIDs: map[string]int{
							"resource-127": 127,
						},
					},
					nil,
				)
			})

			It("constructs teamDB with provided team name", func() {
				Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
				Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
			})

			It("returns 200", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns application/json", func() {
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json"))
			})

			It("returns a json representation of all the versions in the pipeline", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`{
				"ResourceVersions": [
					{
						"VersionID": 73,
						"ResourceID": 127,
						"CheckOrder": 123
			    }
				],
				"BuildOutputs": [
					{
						"VersionID": 73,
						"ResourceID": 127,
						"BuildID": 66,
						"JobID": 13,
						"CheckOrder": 123
					}
				],
				"BuildInputs": [
					{
						"VersionID": 66,
						"ResourceID": 77,
						"BuildID": 66,
						"JobID": 13,
						"CheckOrder": 88,
						"InputName": "some-input-name"
					}
				],
				"JobIDs": {
						"bad-luck-job": 13
				},
				"ResourceIDs": {
					"resource-127": 127
				}
				}`))
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/pipelines/:pipeline_name/rename", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			request, err := http.NewRequest("PUT", server.URL+"/api/v1/teams/a-team/pipelines/a-pipeline/rename", bytes.NewBufferString(`{"name":"some-new-name"}`))
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when authenticated", func() {
			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("a-team", true, true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
					fakeTeam.PipelineReturns(dbPipeline, true, nil)
				})

				It("constructs teamDB with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("injects the proper pipeline", func() {
					pipelineName := fakeTeam.PipelineArgsForCall(0)
					Expect(pipelineName).To(Equal("a-pipeline"))
				})

				It("returns 204", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("renames the pipeline to the name provided", func() {
					Expect(dbPipeline.RenameCallCount()).To(Equal(1))
					Expect(dbPipeline.RenameArgsForCall(0)).To(Equal("some-new-name"))
				})

				Context("when an error occurs on update", func() {
					BeforeEach(func() {
						fakeTeam.PipelineReturns(dbPipeline, true, nil)
						dbPipeline.RenameReturns(errors.New("whoops"))
					})

					It("returns a 500 internal server error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					authValidator.IsAuthenticatedReturns(true)
					userContextReader.GetTeamReturns("another-team", true, true)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
