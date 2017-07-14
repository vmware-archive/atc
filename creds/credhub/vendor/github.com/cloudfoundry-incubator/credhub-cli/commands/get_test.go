package commands_test

import (
	"net/http"

	"runtime"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Get", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get", "-n", "test-credential")
	ItAutomaticallyLogsIn("GET", "get", "-n", "test-credential")

	ItBehavesLikeHelp("get", "g", func(session *Session) {
		Expect(session.Err).To(Say("Usage"))
		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("credhub-cli.exe \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		} else {
			Expect(session.Err).To(Say("credhub-cli \\[OPTIONS\\] get \\[get-OPTIONS\\]"))
		}
	})

	It("displays missing required parameter", func() {
		session := runCommand("get")

		Eventually(session).Should(Exit(1))

		if runtime.GOOS == "windows" {
			Expect(session.Err).To(Say("A name or ID must be provided. Please update and retry your request."))
		} else {
			Expect(session.Err).To(Say("A name or ID must be provided. Please update and retry your request."))
		}
	})

	It("gets a value secret", func() {
		responseJson := fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "value", "my-value", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-value&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-value")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyValuePotatoesYaml))
	})

	It("gets a password secret", func() {
		responseJson := fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "password", "my-password", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-password&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-password")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPasswordPotatoesYaml))
	})

	It("gets a json secret", func() {
		serverResponse := fmt.Sprintf(JSON_CREDENTIAL_ARRAY_RESPONSE_JSON, "json-secret", `{"foo":"bar","nested":{"a":1},"an":["array"]}`)

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=json-secret&current=true"),
				RespondWith(http.StatusOK, serverResponse),
			),
		)

		session := runCommand("get", "-n", "json-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyJsonFormatYaml))
	})

	It("gets a certificate secret", func() {
		responseJson := fmt.Sprintf(CERTIFICATE_CREDENTIAL_ARRAY_RESPONSE_JSON, "my-secret", "my-ca", "my-cert", "my-priv")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-secret&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-secret")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyCertificateYaml))
	})

	It("gets an rsa secret", func() {
		responseJson := fmt.Sprintf(RSA_SSH_CREDENTIAL_ARRAY_RESPONSE_JSON, "rsa", "foo-rsa-key", "some-public-key", "some-private-key")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=foo-rsa-key&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "foo-rsa-key")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyRSAFooYaml))
	})

	It("can output json", func() {
		responseJson := fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "password", "my-password", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-password&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-password", "--output-json")

		Eventually(session).Should(Exit(0))
		Eventually(string(session.Out.Contents())).Should(MatchJSON(`{
			"id": "` + UUID + `",
			"type": "password",
			"name": "my-password",
			"version_created_at": "` + TIMESTAMP + `",
			"value": "potatoes"
		}`))
	})

	It("gets a user secret", func() {
		responseJson := fmt.Sprintf(USER_CREDENTIAL_ARRAY_RESPONSE_JSON, "my-username-credential", "my-username", "test-password", "passw0rd-H4$h")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=my-username-credential&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "my-username-credential")

		Eventually(session).Should(Exit(0))
		Expect(session.Out.Contents()).To(ContainSubstring(responseMyUsernameYaml))
	})

	It("gets a secret by ID", func() {
		responseJson := fmt.Sprintf(STRING_CREDENTIAL_ARRAY_RESPONSE_JSON, "password", "my-password", "potatoes")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data/"+UUID),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "--id", UUID)

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say(responseMyPasswordPotatoesYaml))
	})

	It("does not use Printf on user-supplied data", func() {
		responseJson := fmt.Sprintf(STRING_CREDENTIAL_RESPONSE_JSON, "password", "injected", "et''%/7(V&`|?m|Ckih$")

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/api/v1/data", "name=injected&current=true"),
				RespondWith(http.StatusOK, responseJson),
			),
		)

		session := runCommand("get", "-n", "injected")

		Eventually(session).Should(Exit(0))
		Eventually(session.Out).Should(Say("et''%/7\\(V&`|\\?m\\|Ckih\\$" + TIMESTAMP))
	})
})
