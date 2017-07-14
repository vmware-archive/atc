package commands_test

import (
	"net/http"

	"fmt"

	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Login", func() {
	var (
		uaaServer *Server
	)

	BeforeEach(func() {
		uaaServer = NewServer()
	})

	AfterEach(func() {
		config.RemoveConfig()
	})

	Describe("with mixed password and client parameters", func() {
		Context("with a client name and username", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-name", "test_client", "--username", "test-username")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Client and password credentials may not be combined. Please update and retry your request with a single login method."))
			})
		})

		Context("with a client secret and username", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-secret", "test_secret", "--username", "test-username")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Client and password credentials may not be combined. Please update and retry your request with a single login method."))
			})
		})

		Context("with a client name and password", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-name", "test_client", "--password", "test-password")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Client and password credentials may not be combined. Please update and retry your request with a single login method."))
			})
		})

		Context("with a client secret and password", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-secret", "test_secret", "--password", "test-password")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Client and password credentials may not be combined. Please update and retry your request with a single login method."))
			})
		})

		Context("with all parameters from both", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-name", "test_client", "--client-secret", "test_secret", "--username", "test-username", "--password", "test-password")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Client and password credentials may not be combined. Please update and retry your request with a single login method."))
			})
		})
	})

	Describe("password flow", func() {
		BeforeEach(func() {
			uaaServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token/"),
					VerifyBody([]byte(`grant_type=password&password=pass&response_type=token&username=user`)),
					RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"refresh_token":"erousflkajqwer",
						"token_type":"bearer",
						"expires_in":3600}`),
				),
			)

			setConfigAuthUrl(uaaServer.URL())
		})

		Context("with a username and a password", func() {
			It("authenticates with the UAA server and saves a token", func() {
				session := runCommand("login", "-u", "user", "-p", "pass")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Login Successful"))
				Eventually(session.Out.Contents()).ShouldNot(ContainSubstring("Setting the target url:"))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with a username and no password", func() {
			It("prompts for a password", func() {
				session := runCommandWithStdin(strings.NewReader("pass\n"), "login", "-u", "user")
				Eventually(session.Out).Should(Say("password:"))
				Eventually(session.Wait("10s").Out).Should(Say("Login Successful"))
				Eventually(session).Should(Exit(0))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with a password and no username", func() {
			It("fails authentication with an error message", func() {
				session := runCommand("login", "-p", "pass")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The combination of parameters in the request is not allowed. Please validate your input and retry your request."))
			})
		})

		Context("with neither a username nor a password", func() {
			It("prompts for a username and password", func() {
				// TODO:  devise an input which echoes the input characters for the user name, much as gopass.GetPasswdMasked()
				// echoes '*', for that we may regression-test the echoing of the username
				setConfigAuthUrl(uaaServer.URL())
				session := runCommandWithStdin(strings.NewReader("user\npass\n"), "login")
				Eventually(session.Out).Should(Say("username:"))
				Eventually(session.Out).Should(Say("password:"))
				Eventually(session.Wait("10s").Out).Should(Say("Login Successful"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	Describe("client flow", func() {
		BeforeEach(func() {
			uaaServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token/"),
					VerifyBody([]byte(`client_id=test_client&client_secret=test_secret&grant_type=client_credentials&response_type=token`)),
					RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"refresh_token":"erousflkajqwer",
						"token_type":"bearer",
						"expires_in":3600}`),
				),
			)

			setConfigAuthUrl(uaaServer.URL())
		})

		Context("with the client name and secret from the CLI", func() {
			It("authenticates with the UAA server and saves a token", func() {
				session := runCommand("login", "--client-name", "test_client", "--client-secret", "test_secret")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Login Successful"))
				Eventually(session.Out.Contents()).ShouldNot(ContainSubstring("Setting the target url:"))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with the client name and secret from the env", func() {
			It("authenticates with the UAA server and saves a token", func() {
				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client", "CREDHUB_SECRET=test_secret"}, "login")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Login Successful"))
				Eventually(session.Out.Contents()).ShouldNot(ContainSubstring("Setting the target url:"))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with the client name from the env and client secret from the CLI", func() {
			It("authenticates with the UAA server and saves a token", func() {
				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client"}, "login", "--client-secret", "test_secret")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Login Successful"))
				Eventually(session.Out.Contents()).ShouldNot(ContainSubstring("Setting the target url:"))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with the client name from the CLI and secret from the env", func() {
			It("authenticates with the UAA server and saves a token", func() {
				session := runCommandWithEnv([]string{"CREDHUB_SECRET=test_secret"}, "login", "--client-name", "test_client")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Login Successful"))
				Eventually(session.Out.Contents()).ShouldNot(ContainSubstring("Setting the target url:"))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			})
		})

		Context("with the client name from the CLI and no client secret", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-name", "test_client")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Both client name and client secret must be provided to authenticate. Please update and retry your request."))
			})
		})

		Context("with the client name from the environment and no client secret", func() {
			It("fails with an error message", func() {
				session := runCommandWithEnv([]string{"CREDHUB_CLIENT=test_client"}, "login")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Both client name and client secret must be provided to authenticate. Please update and retry your request."))
			})
		})

		Context("with the client secret from the CLI and no client name", func() {
			It("fails with an error message", func() {
				session := runCommand("login", "--client-secret", "test_secret")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Both client name and client secret must be provided to authenticate. Please update and retry your request."))
			})
		})

		Context("with the client secret from the environment and no client name", func() {
			It("fails with an error message", func() {
				session := runCommandWithEnv([]string{"CREDHUB_SECRET=test_secret"}, "login")

				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Both client name and client secret must be provided to authenticate. Please update and retry your request."))
			})
		})
	})

	Context("when logging in with server api target", func() {
		var (
			uaaServer *Server
			apiServer *Server
		)

		BeforeEach(func() {
			uaaServer = NewServer()
			uaaServer.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token/"),
					VerifyBody([]byte(`grant_type=password&password=pass&response_type=token&username=user`)),
					RespondWith(http.StatusOK, `{
						"access_token":"2YotnFZFEjr1zCsicMWpAA",
						"refresh_token":"erousflkajqwer",
						"token_type":"bearer",
						"expires_in":3600}`),
				),
			)

			apiServer = NewServer()
			setupServer(apiServer, uaaServer.URL())
		})

		AfterEach(func() {
			apiServer.Close()
			uaaServer.Close()
		})

		It("sets the target to the server's url and auth server url", func() {
			session := runCommand("login", "-u", "user", "-p", "pass", "-s", apiServer.URL())

			Expect(apiServer.ReceivedRequests()).Should(HaveLen(1))
			Expect(uaaServer.ReceivedRequests()).Should(HaveLen(1))
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(Say("Login Successful"))
			cfg := config.ReadConfig()
			Expect(cfg.ApiURL).To(Equal(apiServer.URL()))
			Expect(cfg.AuthURL).To(Equal(uaaServer.URL()))
		})

		It("saves caCert to config when it is provided", func() {
			session := runCommand("login", "-u", "user", "-p", "pass", "-s", apiServer.URL(), "--ca-cert", "../test/test-ca.pem")

			Expect(session).Should(Exit(0))
			cfg := config.ReadConfig()

			Expect(cfg.CaCert).Should(Equal([]string{"../test/test-ca.pem"}))
		})

		Context("when the user skips TLS validation", func() {

			It("prints warning and deprecation notice when --skip-tls-validation flag is present", func() {
				apiServer.Close()
				apiServer = NewTLSServer()
				setupServer(apiServer, uaaServer.URL())
				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass", "--skip-tls-validation")

				Eventually(session).Should(Exit(0))
				Eventually(session.Out).Should(Say("Warning: The targeted TLS certificate has not been verified for this connection."))
				Eventually(session.Out).Should(Say("Warning: The --skip-tls-validation flag is deprecated. Please use --ca-cert instead."))
			})

			It("sets skip-tls flag in the config file", func() {
				apiServer.Close()
				apiServer = NewTLSServer()
				setupServer(apiServer, uaaServer.URL())
				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass", "--skip-tls-validation")

				Eventually(session).Should(Exit(0))
				cfg := config.ReadConfig()
				Expect(cfg.InsecureSkipVerify).To(Equal(true))
			})

			It("resets skip-tls flag in the config file", func() {
				cfg := config.ReadConfig()
				cfg.InsecureSkipVerify = true
				err := config.WriteConfig(cfg)
				Expect(err).NotTo(HaveOccurred())

				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass")

				Eventually(session).Should(Exit(0))
				cfg = config.ReadConfig()
				Expect(cfg.InsecureSkipVerify).To(Equal(false))
			})

			It("using a TLS server without the skip-tls flag set will fail on certificate verification", func() {
				apiServer.Close()
				apiServer = NewTLSServer()
				setupServer(apiServer, uaaServer.URL())
				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("Error connecting to the targeted API"))
			})

			It("using a TLS server with the skip-tls flag set will succeed", func() {
				apiServer.Close()
				apiServer = NewTLSServer()
				setupServer(apiServer, uaaServer.URL())
				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass", "--skip-tls-validation")

				Eventually(session).Should(Exit(0))
			})

			It("records skip-tls into config file even with http URLs (will do nothing with that value)", func() {
				session := runCommand("login", "-s", apiServer.URL(), "-u", "user", "-p", "pass", "--skip-tls-validation")
				cfg := config.ReadConfig()

				Eventually(session).Should(Exit(0))
				Expect(cfg.InsecureSkipVerify).To(Equal(true))
			})
		})

		It("saves the oauth tokens", func() {
			runCommand("login", "-u", "user", "-p", "pass", "-s", apiServer.URL())

			cfg := config.ReadConfig()
			Expect(cfg.AccessToken).To(Equal("2YotnFZFEjr1zCsicMWpAA"))
			Expect(cfg.RefreshToken).To(Equal("erousflkajqwer"))
		})

		Context("when api server is unavailable", func() {
			var (
				badServer *Server
			)

			BeforeEach(func() {
				badServer = NewServer()
				badServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("GET", "/info"),
						RespondWith(http.StatusBadGateway, nil),
					),
				)
			})

			It("should not login", func() {
				session := runCommand("login", "-u", "user", "-p", "pass", "-s", badServer.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The targeted API does not appear to be valid. Please validate the API address and retry your request."))
				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
			})

			It("should not override config's existing API URL value", func() {
				cfg := config.ReadConfig()
				cfg.ApiURL = "foo"
				config.WriteConfig(cfg)

				session := runCommand("login", "-u", "user", "-p", "pass", "-s", badServer.URL())

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The targeted API does not appear to be valid. Please validate the API address and retry your request."))
				Expect(uaaServer.ReceivedRequests()).Should(HaveLen(0))
				cfg2 := config.ReadConfig()
				Expect(cfg2.ApiURL).To(Equal("foo"))
			})
		})

		Context("when credentials are invalid", func() {
			var (
				apiServer    *Server
				badUaaServer *Server
				session      *Session
			)

			BeforeEach(func() {
				badUaaServer = NewServer()
				badUaaServer.AppendHandlers(
					CombineHandlers(
						VerifyRequest("POST", "/oauth/token/"),
						VerifyBody([]byte(`grant_type=password&password=pass&response_type=token&username=user`)),
						RespondWith(http.StatusUnauthorized, `{
						"error":"unauthorized",
						"error_description":"An Authentication object was not found in the SecurityContext"
						}`),
					),
					CombineHandlers(
						VerifyRequest("DELETE", "/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r"),
						RespondWith(http.StatusOK, ""),
					),
				)

				apiServer = NewServer()
				setupServer(apiServer, badUaaServer.URL())

				cfg := config.ReadConfig()
				cfg.AuthURL = badUaaServer.URL()
				cfg.AccessToken = "fake_token"
				cfg.RefreshToken = "5b9c9fd51ba14838ac2e6b222d487106-r"
				config.WriteConfig(cfg)
			})

			It("fails to login", func() {
				session = runCommand("login", "-u", "user", "-p", "pass")
				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("The provided username and password combination are incorrect. Please validate your input and retry your request."))
				Expect(badUaaServer.ReceivedRequests()).Should(HaveLen(2))
			})

			It("revokes any existing tokens", func() {
				session = runCommand("login", "-u", "user", "-p", "pass")
				Eventually(session).Should(Exit(1))
				cfg := config.ReadConfig()
				Expect(cfg.AccessToken).To(Equal("revoked"))
				Expect(cfg.RefreshToken).To(Equal("revoked"))
				Expect(badUaaServer.ReceivedRequests()).Should(HaveLen(2))
			})

			It("doesn't print 'Setting the target url' message with -s flag", func() {
				session = runCommand("login", "-u", "user", "-p", "pass", "-s", apiServer.URL())
				Eventually(session).Should(Exit(1))
				Expect(session.Out).NotTo(Say("Setting the target url: " + apiServer.URL()))
			})
		})
	})

	Describe("when logging in without server api target", func() {
		var (
			apiUrl string
		)

		BeforeEach(func() {
			cfg := config.ReadConfig()
			apiUrl = cfg.ApiURL
			cfg.ApiURL = ""
			config.WriteConfig(cfg)
		})

		AfterEach(func() {
			cfg := config.ReadConfig()
			cfg.ApiURL = apiUrl
			config.WriteConfig(cfg)
		})

		Context("with no user or password flags", func() {
			It("returns an error message", func() {
				session := runCommand("login")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue."))
			})
		})

		Context("with user and password flags", func() {
			It("returns an error message", func() {
				session := runCommand("login", "-u", "user", "-p", "pass")

				Eventually(session).Should(Exit(1))
				Eventually(session.Err).Should(Say("An API target is not set. Please target the location of your server with `credhub api --server api.example.com` to continue."))
			})
		})
	})

	Describe("Help", func() {
		ItBehavesLikeHelp("login", "l", func(session *Session) {
			Expect(session.Err).To(Say("login"))
			Expect(session.Err).To(Say("username"))
			Expect(session.Err).To(Say("password"))
			Expect(session.Err).To(Say("client-name"))
			Expect(session.Err).To(Say("client-secret"))
		})
	})
})

func setConfigAuthUrl(authUrl string) {
	cfg := config.ReadConfig()
	cfg.AuthURL = authUrl
	config.WriteConfig(cfg)
}

func setupServer(theServer *Server, uaaUrl string) {
	theServer.AppendHandlers(
		CombineHandlers(
			VerifyRequest("GET", "/info"),
			RespondWith(http.StatusOK, fmt.Sprintf(`{
					"app":{"version":"0.1.0 build DEV","name":"CredHub"},
					"auth-server":{"url":"%s"}
					}`, uaaUrl)),
		),
	)
}
