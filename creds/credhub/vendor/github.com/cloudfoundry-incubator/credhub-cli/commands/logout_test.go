package commands_test

import (
	"net/http"

	"runtime"

	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Logout", func() {
	AfterEach(func() {
		config.RemoveConfig()
	})

	It("marks the access token and refresh token as revoked if no config exists", func() {
		config.RemoveConfig()
		runLogoutCommand()
	})

	It("leaves the access token and refresh token as revoked if config exists and they were already revoked", func() {
		cfg := config.Config{RefreshToken: "revoked", AccessToken: "revoked"}
		config.WriteConfig(cfg)
		runLogoutCommand()
	})

	It("asks UAA to revoke the refresh token (and UAA succeeds)", func() {
		doRevoke(http.StatusOK)
	})

	It("asks UAA to revoke the refresh token (and reports no error when UAA fails)", func() {
		doRevoke(http.StatusUnauthorized)
	})

	ItBehavesLikeHelp("logout", "o", func(session *Session) {
		Expect(session.Err).To(Say("Usage:"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] logout"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] logout"))
		}
	})
})

func doRevoke(uaaResponseStatus int) {
	cfg := config.Config{
		RefreshToken: "5b9c9fd51ba14838ac2e6b222d487106-r",
		AccessToken:  "myAccessToken",
		AuthURL:      authServer.URL(),
	}
	config.WriteConfig(cfg)

	authServer.AppendHandlers(
		CombineHandlers(
			VerifyRequest("DELETE", "/oauth/token/revoke/5b9c9fd51ba14838ac2e6b222d487106-r"),
			RespondWith(uaaResponseStatus, ""),
		),
	)
	runLogoutCommand()
}

func runLogoutCommand() {
	session := runCommand("logout")
	Eventually(session).Should(Exit(0))
	Eventually(session).Should(Say("Logout Successful"))
	cfg := config.ReadConfig()
	Expect(cfg.AccessToken).To(Equal("revoked"))
	Expect(cfg.RefreshToken).To(Equal("revoked"))
}
