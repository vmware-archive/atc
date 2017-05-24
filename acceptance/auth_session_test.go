package acceptance_test

import (
	"net/http"

	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"

	"github.com/concourse/atc/auth"

	"encoding/json"
	"io/ioutil"

	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Session", func() {
	var atcCommand *ATCCommand
	var page *agouti.Page
	var pipeline dbng.Pipeline

	BeforeEach(func() {
		atcCommand = NewATCCommand(atcBin, 1, postgresRunner.DataSourceName(), []string{}, NO_AUTH)
		err := atcCommand.Start()
		Expect(err).NotTo(HaveOccurred())

		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())

		teamFactory := dbng.NewTeamFactory(dbngConn, lockFactory, dbng.NewNoEncryption())
		defaultTeam, found, err := teamFactory.FindTeam(atc.DefaultTeamName)
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue()) // created by postgresRunner

		pipeline, _, err = defaultTeam.SavePipeline("main", atc.Config{
			Jobs: atc.JobConfigs{
				{
					Name: "job-name",
				},
			},
		}, dbng.ConfigVersion(1), dbng.PipelineUnpaused)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())

		atcCommand.Stop()
	})

	It("generates auth token cookie", func() {
		LoginWithNoAuth(page, atcCommand.URL(""))
		cookies, err := page.GetCookies()
		Expect(err).NotTo(HaveOccurred())
		var authCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == auth.AuthCookieName {
				authCookie = cookie
			}
		}
		Expect(authCookie).NotTo(BeNil())
		Expect(authCookie.HttpOnly).To(BeTrue())
		Expect(authCookie.Secure).To(BeFalse())
	})

	Context("when request does not contain CSRF token", func() {
		It("returns 400 Bad Request", func() {
			LoginWithNoAuth(page, atcCommand.URL(""))
			Expect(page.RunScript("return localStorage.removeItem('csrf_token')", nil, nil)).To(Succeed())

			Expect(page.Navigate(atcCommand.URL("/teams/main/pipelines/main/jobs/job-name"))).To(Succeed())
			Eventually(page.Find(".build-action")).Should(BeFound())
			Expect(page.Find(".build-action").Click()).To(Succeed())

			// API request will return bad request
			// no changes in UI, no builds will be scheduled
			job, found, err := pipeline.Job("job-name")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())

			builds, _, err := job.Builds(dbng.Page{Limit: 1})
			Expect(err).ToNot(HaveOccurred())

			Expect(builds).To(HaveLen(0))
		})
	})

	Context("when request contains invalid CSRF token", func() {
		It("returns 401 Not Authorized and redirects to login page", func() {
			LoginWithNoAuth(page, atcCommand.URL(""))
			Expect(page.RunScript("return localStorage.setItem('csrf_token', 'invalid')", nil, nil)).To(Succeed())

			Expect(page.Navigate(atcCommand.URL("/teams/main/pipelines/main/jobs/job-name"))).To(Succeed())
			Eventually(page.Find(".build-action")).Should(BeFound())
			Expect(page.Find(".build-action").Click()).To(Succeed())

			Eventually(page.FindByButton("login")).Should(BeFound())
		})
	})

	Context("when CSRF token and session token are not associated", func() {
		It("returns 401 Not Authorized and redirects to login page", func() {
			LoginWithNoAuth(page, atcCommand.URL(""))

			var firstCSRFToken string
			Expect(page.RunScript("return localStorage.getItem('csrf_token')", nil, &firstCSRFToken)).To(Succeed())
			Expect(firstCSRFToken).NotTo(BeNil())

			LoginWithNoAuth(page, atcCommand.URL(""))

			Expect(page.RunScript("return localStorage.setItem('csrf_token', '"+firstCSRFToken+"')", nil, nil)).To(Succeed())

			Expect(page.Navigate(atcCommand.URL("/teams/main/pipelines/main/jobs/job-name"))).To(Succeed())
			Eventually(page.Find(".build-action")).Should(BeFound())
			Expect(page.Find(".build-action").Click()).To(Succeed())

			Eventually(page.FindByButton("login")).Should(BeFound())
		})
	})

	Context("when request contains valid CSRF with associated session token", func() {
		It("returns 200 OK", func() {
			LoginWithNoAuth(page, atcCommand.URL(""))
			Expect(page.Navigate(atcCommand.URL("/teams/main/pipelines/main/jobs/job-name"))).To(Succeed())
			Eventually(page.Find(".build-action")).Should(BeFound())
			Expect(page.Find(".build-action").Click()).To(Succeed())
			Eventually(page.All("#builds li").Count).Should(Equal(1))
		})
	})

	Context("when request has authorization token in header", func() {
		var atcToken atc.AuthToken
		var client *http.Client

		BeforeEach(func() {
			request, err := http.NewRequest("GET", atcCommand.URL("/api/v1/teams/main/auth/token"), nil)
			Expect(err).NotTo(HaveOccurred())

			client = &http.Client{
				Transport: &http.Transport{},
			}
			resp, err := client.Do(request)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			err = json.Unmarshal(body, &atcToken)
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not require CSRF token", func() {
			request, err := http.NewRequest("GET", atcCommand.URL("/api/v1/teams/main/pipelines/main"), nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("Authorization", atcToken.Type+" "+atcToken.Value)

			response, err := client.Do(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
