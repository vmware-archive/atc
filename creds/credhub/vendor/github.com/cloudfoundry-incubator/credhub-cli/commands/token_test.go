package commands_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Token", func() {
	var (
		uaaServer *Server
	)

	BeforeEach(func() {
		uaaServer = NewServer()
	})

	AfterEach(func() {
		config.RemoveConfig()
	})

	Context("when the config file has a token", func() {

		BeforeEach(func() {
			cfg := config.ReadConfig()
			cfg.AccessToken = "2YotnFZFEjr1zCsicMWpAA"
			config.WriteConfig(cfg)

			uaaServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token/"),
					VerifyBody([]byte(`grant_type=refresh_token&refresh_token=revoked`)),
					RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"refresh_token":"erousflkajqwer",
						"token_type":"bearer",
						"expires_in":3600}`),
				),
			)

			setConfigAuthUrl(uaaServer.URL())
		})

		It("refreshes the token with --token", func() {
			session := runCommand("--token")

			Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			Expect(sout).To(ContainSubstring("Bearer 2YotnFZFEjr1zCsicMWpAA"))
		})
	})

	Context("when the config file does not have a token", func() {
		BeforeEach(func() {
			cfg := config.ReadConfig()
			cfg.AccessToken = ""
			config.WriteConfig(cfg)
		})

		It("displays nothing", func() {
			session := runCommand("--token")

			Eventually(session).Should(Exit(0))
			sout := string(session.Out.Contents())
			Expect(sout).To(Equal(""))
		})
	})
})
