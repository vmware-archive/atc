package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func jsonEncode(object interface{}) *bytes.Buffer {
	reqPayload, err := json.Marshal(object)
	Expect(err).NotTo(HaveOccurred())

	return bytes.NewBuffer(reqPayload)
}

func atcDBTeamEquality(atcTeam atc.Team, dbTeam db.Team) {
	Expect(dbTeam.Name).To(Equal(atcTeam.Name))
	if atcTeam.BasicAuth == nil {
		Expect(dbTeam.BasicAuth).To(BeNil())
	} else {
		Expect(dbTeam.BasicAuth).NotTo(BeNil())
		Expect(dbTeam.BasicAuth.BasicAuthUsername).To(Equal(atcTeam.BasicAuth.BasicAuthUsername))
		Expect(dbTeam.BasicAuth.BasicAuthPassword).To(Equal(atcTeam.BasicAuth.BasicAuthPassword))
	}
	if atcTeam.GitHubAuth == nil {
		Expect(dbTeam.GitHubAuth).To(BeNil())
	} else {
		Expect(dbTeam.GitHubAuth).NotTo(BeNil())
		Expect(dbTeam.GitHubAuth.ClientID).To(Equal(atcTeam.GitHubAuth.ClientID))
		Expect(dbTeam.GitHubAuth.ClientSecret).To(Equal(atcTeam.GitHubAuth.ClientSecret))
		Expect(dbTeam.GitHubAuth.Organizations).To(Equal(atcTeam.GitHubAuth.Organizations))
		Expect(dbTeam.GitHubAuth.Users).To(Equal(atcTeam.GitHubAuth.Users))
		Expect(len(dbTeam.GitHubAuth.Teams)).To(Equal(len(atcTeam.GitHubAuth.Teams)))
		for i, atcGitHubTeam := range atcTeam.GitHubAuth.Teams {
			dbGitHubTeam := dbTeam.GitHubAuth.Teams[i]
			Expect(dbGitHubTeam.OrganizationName).To(Equal(atcGitHubTeam.OrganizationName))
			Expect(dbGitHubTeam.TeamName).To(Equal(atcGitHubTeam.TeamName))
		}
	}
}

var _ = Describe("Teams API", func() {
	Describe("GET /api/v1/teams", func() {
		var response *http.Response

		JustBeforeEach(func() {
			path := fmt.Sprintf("%s/api/v1/teams", server.URL)

			request, err := http.NewRequest("GET", path, nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the database returns an error", func() {
			var disaster error

			BeforeEach(func() {
				disaster = errors.New("some error")
				teamServerDB.GetTeamsReturns(nil, disaster)
			})

			It("returns 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when the database returns teams", func() {
			BeforeEach(func() {
				teamServerDB.GetTeamsReturns([]db.SavedTeam{
					{
						ID: 5,
						Team: db.Team{
							Name: "avengers",
						},
					},
					{
						ID: 9,
						Team: db.Team{
							Name: "aliens",
							BasicAuth: &db.BasicAuth{
								BasicAuthUsername: "fake user",
								BasicAuthPassword: "no, bad",
							},
							GitHubAuth: &db.GitHubAuth{
								ClientID:      "fake id",
								ClientSecret:  "some secret",
								Organizations: []string{"a", "b", "c"},
								Teams: []db.GitHubTeam{
									{
										OrganizationName: "org1",
										TeamName:         "teama",
									},
									{
										OrganizationName: "org2",
										TeamName:         "teamb",
									},
								},
								Users: []string{"user1", "user2", "user3"},
							},
						},
					},
					{
						ID: 22,
						Team: db.Team{
							Name: "predators",
							UAAAuth: &db.UAAAuth{
								ClientID:     "fake id",
								ClientSecret: "some secret",
								CFSpaces:     []string{"myspace"},
								AuthURL:      "http://auth.url",
								TokenURL:     "http://token.url",
								CFURL:        "http://api.url",
							},
						},
					},
				}, nil)
			})

			It("returns 200 OK", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns the teams", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
					{
						"id": 5,
						"name": "avengers"
					},
					{
						"id": 9,
						"name": "aliens"
					},
					{
						"id": 22,
						"name": "predators"
					}
				]`))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name", func() {
		var request *http.Request
		var response *http.Response

		var team atc.Team
		var savedTeam db.SavedTeam
		var teamName string

		BeforeEach(func() {
			teamName = "team venture"

			team = atc.Team{}
			savedTeam = db.SavedTeam{
				ID: 2,
				Team: db.Team{
					Name: teamName,
				},
			}
		})

		Context("when the requester is authenticated for the right team (admin team)", func() {
			JustBeforeEach(func() {
				path := fmt.Sprintf("%s/api/v1/teams/%s", server.URL, teamName)

				var err error
				request, err = http.NewRequest("PUT", path, jsonEncode(team))
				Expect(err).NotTo(HaveOccurred())

				response, err = client.Do(request)
				Expect(err).NotTo(HaveOccurred())
			})

			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns(atc.DefaultTeamName, 1, true, true)
			})

			Describe("request body validation", func() {
				Describe("basic authenticaiton", func() {
					Context("BasicAuthUsername not filled in", func() {
						BeforeEach(func() {
							team = atc.Team{
								BasicAuth: &atc.BasicAuth{
									BasicAuthPassword: "Batman",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})

					Context("BasicAuthPassword not filled in", func() {
						BeforeEach(func() {
							team = atc.Team{
								BasicAuth: &atc.BasicAuth{
									BasicAuthUsername: "Hank Venture",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})
				})

				Describe("GitHub authenticaiton", func() {
					Context("ClientID not filled in", func() {
						BeforeEach(func() {
							team = atc.Team{
								GitHubAuth: &atc.GitHubAuth{
									ClientSecret: "09262-8765-001",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})

					Context("ClientSecret not filled in", func() {
						BeforeEach(func() {
							team = atc.Team{
								GitHubAuth: &atc.GitHubAuth{
									ClientID: "Brock Samson",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})

					Context("require at least one org, org/team, or username", func() {
						Context("when all are missing", func() {
							BeforeEach(func() {
								team = atc.Team{
									GitHubAuth: &atc.GitHubAuth{
										ClientID:     "Brock Samson",
										ClientSecret: "09262-8765-001",
									},
								}
							})

							It("returns a 400 Bad Request", func() {
								Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
							})
						})

						Context("when passed organizations", func() {
							BeforeEach(func() {
								team = atc.Team{
									GitHubAuth: &atc.GitHubAuth{
										ClientID:      "Brock Samson",
										ClientSecret:  "09262-8765-001",
										Organizations: []string{"United States Armed Forces", "Office of Secret Intelligence", "Team Venture", "S.P.H.I.N.X."},
									},
								}
							})

							It("does not error", func() {
								Expect(response.StatusCode).To(Equal(http.StatusCreated))
							})
						})

						Context("when passed a team", func() {
							BeforeEach(func() {
								team = atc.Team{
									GitHubAuth: &atc.GitHubAuth{
										ClientID:     "Brock Samson",
										ClientSecret: "09262-8765-001",
										Teams: []atc.GitHubTeam{
											{
												OrganizationName: "Office of Secret Intelligence",
												TeamName:         "Secret Agent",
											},
										},
									},
								}
							})

							It("does not error", func() {
								Expect(response.StatusCode).To(Equal(http.StatusCreated))
							})
						})

						Context("when passed users", func() {
							BeforeEach(func() {
								team = atc.Team{
									GitHubAuth: &atc.GitHubAuth{
										ClientID:     "S.P.H.I.N.X.",
										ClientSecret: "SPHINX Rising",
										Users: []string{
											"Col. Hunter Gathers",
											"Holy Diver/Shore Leave",
											"Mile High/Sky Pilot",
											"Brock Samson",
											"Unnamed German plastic surgeon",
										},
									},
								}
							})

							It("does not error", func() {
								Expect(response.StatusCode).To(Equal(http.StatusCreated))
							})
						})
					})
				})

				Describe("UAA authentication", func() {
					Context("when passed a valid team with UAA Auth", func() {
						BeforeEach(func() {
							team = atc.Team{
								UAAAuth: &atc.UAAAuth{
									ClientID:     "Brock Samson",
									ClientSecret: "09262-8765-001",
									CFSpaces:     []string{"myspace"},
									AuthURL:      "http://auth.url",
									TokenURL:     "http://token.url",
									CFURL:        "http://api.url",
								},
							}
						})

						It("responds with 201", func() {
							Expect(response.StatusCode).To(Equal(http.StatusCreated))
						})
					})

					Context("ClientID is not filled in", func() {
						BeforeEach(func() {
							team = atc.Team{
								UAAAuth: &atc.UAAAuth{
									ClientSecret: "09262-8765-001",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})

					Context("Spaces are not provided", func() {
						BeforeEach(func() {
							team = atc.Team{
								UAAAuth: &atc.UAAAuth{
									ClientID:     "S.P.H.I.N.X.",
									ClientSecret: "09262-8765-001",
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})

					Context("AuthURL is not provided", func() {
						BeforeEach(func() {
							team = atc.Team{
								UAAAuth: &atc.UAAAuth{
									ClientID:     "S.P.H.I.N.X.",
									ClientSecret: "09262-8765-001",
									CFSpaces:     []string{"myspace"},
								},
							}
						})

						It("returns a 400 Bad Request", func() {
							Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
						})
					})
				})
			})

			Context("when there's a problem finding teams", func() {
				BeforeEach(func() {
					teamDB.GetTeamReturns(db.SavedTeam{}, false, errors.New("a dingo ate my baby!"))
				})

				It("returns 500 Internal Server Error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when team exists", func() {
				BeforeEach(func() {
					teamDB.GetTeamReturns(savedTeam, true, nil)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("returns the updated team", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`{
					"id": 2,
					"name": "team venture"
				}`))

					Expect(teamServerDB.CreateTeamCallCount()).To(Equal(0))
				})

				Context("updating authentication", func() {
					var basicAuth *atc.BasicAuth
					var gitHubAuth *atc.GitHubAuth
					var uaaAuth *atc.UAAAuth
					var genericOAuth *atc.GenericOAuth

					BeforeEach(func() {
						basicAuth = &atc.BasicAuth{
							BasicAuthUsername: "Dean Venture",
							BasicAuthPassword: "Giant Boy Detective",
						}

						gitHubAuth = &atc.GitHubAuth{
							ClientID:     "Dean Venture",
							ClientSecret: "Giant Boy Detective",
							Users:        []string{"Dean Venture"},
						}

						uaaAuth = &atc.UAAAuth{
							ClientID:     "Dean Venture",
							ClientSecret: "Giant Boy Detective",
							CFSpaces:     []string{"CSI"},
							AuthURL:      "http://uaa.auth.url",
							TokenURL:     "http://uaa.token.url",
							CFURL:        "http://api.cf.url",
						}

						genericOAuth = &atc.GenericOAuth{
							ClientID:      "Dean Venture",
							ClientSecret:  "Giant Boy Detective",
							AuthURL:       "https://goa.auth.url",
							AuthURLParams: map[string]string{},
							TokenURL:      "https://goa.token.url",
							DisplayName:   "CSI",
						}
					})

					Context("when passed basic auth credentials", func() {
						BeforeEach(func() {
							teamDB.UpdateBasicAuthStub = func(basicAuth *db.BasicAuth) (db.SavedTeam, error) {
								team.Name = teamName
								Expect(basicAuth).NotTo(BeNil())
								Expect(team.BasicAuth).NotTo(BeNil())
								Expect(basicAuth.BasicAuthUsername).To(Equal(team.BasicAuth.BasicAuthUsername))
								Expect(basicAuth.BasicAuthPassword).To(Equal(team.BasicAuth.BasicAuthPassword))
								savedTeam.BasicAuth = basicAuth
								return savedTeam, nil
							}

							team.BasicAuth = basicAuth
						})

						It("updates the basic auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusOK))
							Expect(teamDB.UpdateBasicAuthCallCount()).To(Equal(1))
						})
					})

					Context("when passed GitHub auth credentials", func() {
						BeforeEach(func() {
							teamDB.UpdateGitHubAuthStub = func(gitHubAuth *db.GitHubAuth) (db.SavedTeam, error) {
								team.Name = teamName
								Expect(gitHubAuth.ClientID).To(Equal(team.GitHubAuth.ClientID))
								Expect(gitHubAuth.ClientSecret).To(Equal(team.GitHubAuth.ClientSecret))
								Expect(gitHubAuth.Organizations).To(Equal(team.GitHubAuth.Organizations))
								Expect(gitHubAuth.Teams).To(HaveLen(len(team.GitHubAuth.Teams)))
								for _, t := range gitHubAuth.Teams {
									Expect(team.GitHubAuth.Teams).To(ContainElement(db.GitHubTeam{
										OrganizationName: t.OrganizationName,
										TeamName:         t.TeamName,
									}))
								}
								Expect(gitHubAuth.Users).To(Equal(team.GitHubAuth.Users))
								savedTeam.GitHubAuth = gitHubAuth
								return savedTeam, nil
							}

							team.GitHubAuth = gitHubAuth
						})

						It("updates the GitHub auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusOK))
							Expect(teamDB.UpdateGitHubAuthCallCount()).To(Equal(1))
						})
					})

					Context("when passed UAA auth credentials", func() {
						BeforeEach(func() {
							teamDB.UpdateUAAAuthStub = func(uaaAuth *db.UAAAuth) (db.SavedTeam, error) {
								team.Name = teamName
								Expect(uaaAuth.ClientID).To(Equal(team.UAAAuth.ClientID))
								Expect(uaaAuth.ClientSecret).To(Equal(team.UAAAuth.ClientSecret))
								Expect(uaaAuth.AuthURL).To(Equal(team.UAAAuth.AuthURL))
								Expect(uaaAuth.TokenURL).To(Equal(team.UAAAuth.TokenURL))
								Expect(uaaAuth.CFSpaces).To(Equal(team.UAAAuth.CFSpaces))
								Expect(uaaAuth.CFURL).To(Equal(team.UAAAuth.CFURL))

								savedTeam.UAAAuth = uaaAuth
								return savedTeam, nil
							}

							team.UAAAuth = uaaAuth
						})

						It("updates the GitHub auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusOK))
							Expect(teamDB.UpdateUAAAuthCallCount()).To(Equal(1))
						})
					})

					Context("when passed generic OAuth auth credentials", func() {
						BeforeEach(func() {
							teamDB.UpdateGenericOAuthStub = func(genericOAuth *db.GenericOAuth) (db.SavedTeam, error) {
								team.Name = teamName
								Expect(genericOAuth.ClientID).To(Equal(team.GenericOAuth.ClientID))
								Expect(genericOAuth.ClientSecret).To(Equal(team.GenericOAuth.ClientSecret))
								Expect(genericOAuth.AuthURL).To(Equal(team.GenericOAuth.AuthURL))
								Expect(genericOAuth.TokenURL).To(Equal(team.GenericOAuth.TokenURL))
								Expect(genericOAuth.AuthURLParams).To(Equal(team.GenericOAuth.AuthURLParams))
								Expect(genericOAuth.DisplayName).To(Equal(team.GenericOAuth.DisplayName))

								savedTeam.GenericOAuth = genericOAuth
								return savedTeam, nil
							}

							team.GenericOAuth = genericOAuth
						})

						It("updates the GitHub auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusOK))
							Expect(teamDB.UpdateUAAAuthCallCount()).To(Equal(1))
						})
					})

				})
			})

			Context("when team does not exist", func() {
				BeforeEach(func() {
					teamDB.GetTeamReturns(db.SavedTeam{}, false, nil)

					teamServerDB.CreateTeamStub = func(submittedTeam db.Team) (db.SavedTeam, error) {
						team.Name = teamName
						atcDBTeamEquality(team, submittedTeam)
						return savedTeam, nil
					}
				})

				It("returns 201 Created", func() {
					Expect(response.StatusCode).To(Equal(http.StatusCreated))
				})

				It("returns the new team", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`{
					"id": 2,
					"name": "team venture"
				}`))

					Expect(teamServerDB.CreateTeamCallCount()).To(Equal(1))
				})

				Context("when there's a problem saving teams", func() {
					BeforeEach(func() {
						teamServerDB.CreateTeamReturns(db.SavedTeam{}, errors.New("Do not be too hasty in entering that room. I had Taco Bell for lunch!"))
					})

					It("returns 500 Internal Server Error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})

				Context("with authentication", func() {
					var basicAuth atc.BasicAuth
					var gitHubAuth atc.GitHubAuth

					BeforeEach(func() {
						basicAuth = atc.BasicAuth{
							BasicAuthUsername: "Dean Venture",
							BasicAuthPassword: "Giant Boy Detective",
						}

						gitHubAuth = atc.GitHubAuth{
							ClientID:     "Dean Venture",
							ClientSecret: "Giant Boy Detective",
							Users:        []string{"Dean Venture"},
						}
					})

					Context("when passed basic auth credentials", func() {
						BeforeEach(func() {
							team.BasicAuth = &basicAuth
						})

						It("updates the basic auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusCreated))
							Expect(teamServerDB.CreateTeamCallCount()).To(Equal(1))
						})
					})

					Context("when passed GitHub auth credentials", func() {
						BeforeEach(func() {
							team.GitHubAuth = &gitHubAuth
						})

						It("updates the GitHub auth for that team", func() {
							Expect(response.StatusCode).To(Equal(http.StatusCreated))
							Expect(teamServerDB.CreateTeamCallCount()).To(Equal(1))
						})
					})
				})
			})
		})

		Context("when the requester belongs to a non-admin team", func() {
			JustBeforeEach(func() {
				path := fmt.Sprintf("%s/api/v1/teams/%s", server.URL, "non-admin-team")

				var err error
				request, err = http.NewRequest("PUT", path, jsonEncode(team))
				Expect(err).NotTo(HaveOccurred())

				response, err = client.Do(request)
				Expect(err).NotTo(HaveOccurred())
			})

			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
			})

			Context("when updating their own team", func() {
				var savedTeam db.SavedTeam
				BeforeEach(func() {
					savedTeam = db.SavedTeam{
						ID: 5,
						Team: db.Team{
							Name: "non-admin-team",
						},
					}
					teamDB.GetTeamReturns(savedTeam, true, nil)
					userContextReader.GetTeamReturns("non-admin-team", 5, false, true)
				})

				It("returns 200 OK", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("always sets Admin property to false", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`{
						"id": 5,
						"name": "non-admin-team"
					}`))

					Expect(teamServerDB.CreateTeamCallCount()).To(Equal(0))
				})

				It("returns the updated team", func() {
					body, err := ioutil.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())

					Expect(body).To(MatchJSON(`{
						"id": 5,
						"name": "non-admin-team"
					}`))

					Expect(teamServerDB.CreateTeamCallCount()).To(Equal(0))
				})
			})

			Context("when updating another team", func() {
				BeforeEach(func() {
					userContextReader.GetTeamReturns("another-non-admin-team", 5, false, true)
				})

				It("returns 403 forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})

			Context("when team does not exist", func() {
				BeforeEach(func() {
					userContextReader.GetTeamReturns("non-admin-team", 5, false, true)
					teamDB.GetTeamReturns(db.SavedTeam{}, false, nil)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
				})
			})
		})

		Context("when the requester's team cannot be determined", func() {
			JustBeforeEach(func() {
				path := fmt.Sprintf("%s/api/v1/teams/%s", server.URL, teamName)

				var err error
				request, err = http.NewRequest("PUT", path, jsonEncode(team))
				Expect(err).NotTo(HaveOccurred())

				response, err = client.Do(request)
				Expect(err).NotTo(HaveOccurred())
			})

			BeforeEach(func() {
				authValidator.IsAuthenticatedReturns(true)
				userContextReader.GetTeamReturns("", 0, false, false)
			})

			It("returns 500 internal server error", func() {
				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
