package commands_test

import (
	"net/http"

	"runtime"

	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Generate", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get", "-n", "test-credential")
	ItAutomaticallyLogsIn("POST", "generate", "-n", "test-credential", "-t", "password")

	It("requires a type", func() {
		session := runCommand("generate", "-n", "my-credential")
		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(Say("A type must be specified when generating a credential. Valid types include 'password', 'user', 'certificate', 'ssh' and 'rsa'."))
	})

	Describe("Without password parameters", func() {
		It("uses default parameters", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson("my-password", `{}`, true))

			session := runCommand("generate", "-n", "my-password", "-t", "password")
			Eventually(session).Should(Exit(0))
		})

		It("prints the generated password secret", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson("my-password", `{}`, true))

			session := runCommand("generate", "-n", "my-password", "-t", "password")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyPasswordPotatoesYaml))
		})

		It("can print the generated password secret as JSON", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson("my-password", `{}`, true))

			session := runCommand("generate", "-n", "my-password", "-t", "password", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(session.Out.Contents()).To(MatchJSON(`{
				"id" :"` + UUID + `",
				"type": "password",
				"name": "my-password",
				"version_created_at": "` + TIMESTAMP + `",
				"value": "potatoes"
			}`))
		})
	})

	Describe("with a variety of password parameters", func() {
		It("prints the secret", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson("my-password", `{}`, true))

			session := runCommand("generate", "-n", "my-password", "-t", "password")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyPasswordPotatoesYaml))
		})

		It("can print the secret as JSON", func() {
			setupPasswordPostServer("my-password", "potatoes", generateDefaultTypeRequestJson("my-password", `{}`, true))

			session := runCommand(
				"generate",
				"-n", "my-password",
				"-t", "password",
				"--output-json",
			)

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"id" :"` + UUID + `",
				"type": "password",
				"name": "my-password",
				"version_created_at": "` + TIMESTAMP + `",
				"value": "potatoes"
			}`))
		})

		It("with with no-overwrite", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{}`, false))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{"length":42}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{"exclude_upper":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{"exclude_lower":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("including special characters", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{"include_special":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--include-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupPasswordPostServer("my-password", "potatoes", generateRequestJson("password", "my-password", `{"exclude_number":true}`, true))
			session := runCommand("generate", "-n", "my-password", "-t", "password", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of SSH parameters", func() {
		It("prints the SSH key", func() {
			setupRsaSshPostServer("foo-ssh-key", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", "foo-ssh-key", `{}`, true))

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMySSHFooYaml))
		})

		It("can print the SSH key as JSON", func() {
			setupRsaSshPostServer("foo-ssh-key", "ssh", "some-public-key", "fake-private-key", generateRequestJson("ssh", "foo-ssh-key", `{}`, true))

			session := runCommand("generate", "-n", "foo-ssh-key", "-t", "ssh", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"id" :"` + UUID + `",
				"type": "ssh",
				"name": "foo-ssh-key",
				"version_created_at": "` + TIMESTAMP + `",
				"value": {
					"public_key": "some-public-key",
					"private_key": "fake-private-key"
				}
			}`))
		})

		It("with with no-overwrite", func() {
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", "my-ssh", `{}`, false))
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", generateRequestJson("ssh", "my-ssh", `{"key_length":3072}`, true))
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})

		It("including comment", func() {
			expectedRequestJson := generateRequestJson("ssh", "my-ssh", `{"ssh_comment":"i am an ssh comment"}`, true)
			setupRsaSshPostServer("my-ssh", "ssh", "some-public-key", "some-private-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-ssh", "-t", "ssh", "-m", "i am an ssh comment")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of RSA parameters", func() {
		It("prints the RSA key", func() {
			setupRsaSshPostServer("foo-rsa-key", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", "foo-rsa-key", `{}`, true))

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyRSAFooYaml))
		})

		It("can print the RSA key as JSON", func() {
			setupRsaSshPostServer("foo-rsa-key", "rsa", "some-public-key", "fake-private-key", generateRequestJson("rsa", "foo-rsa-key", `{}`, true))

			session := runCommand("generate", "-n", "foo-rsa-key", "-t", "rsa", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"id" :"` + UUID + `",
				"type": "rsa",
				"name": "foo-rsa-key",
				"version_created_at": "` + TIMESTAMP + `",
				"value": {
					"public_key": "some-public-key",
					"private_key": "fake-private-key"
				}
			}`))
		})

		It("with with no-overwrite", func() {
			setupRsaSshPostServer("my-rsa", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", "my-rsa", `{}`, false))
			session := runCommand("generate", "-n", "my-rsa", "-t", "rsa", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupRsaSshPostServer("my-rsa", "rsa", "some-public-key", "some-private-key", generateRequestJson("rsa", "my-rsa", `{"key_length":3072}`, true))
			session := runCommand("generate", "-n", "my-rsa", "-t", "rsa", "-k", "3072")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of certificate parameters", func() {
		It("prints the certificate", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "my-ca", "my-cert", "my-priv", expectedRequestJson)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")

			Eventually(session).Should(Exit(0))
			Expect(session.Out).To(Say(responseMyCertificateYaml))
		})

		It("can print the certificate as JSON", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "my-ca", "my-cert", "my-priv", expectedRequestJson)

			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--output-json")

			Eventually(session).Should(Exit(0))
			Expect(string(session.Out.Contents())).To(MatchJSON(`{
				"id" :"` + UUID + `",
				"type": "certificate",
				"name": "my-secret",
				"version_created_at": "` + TIMESTAMP + `",
				"value": {
					"ca": "my-ca",
					"certificate": "my-cert",
					"private_key": "my-priv"
				}
			}`))
		})

		It("including common name", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"common_name":"common.name.io"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io")
			Eventually(session).Should(Exit(0))
		})

		It("including common name with no-overwrite", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"common_name":"common.name.io"}`, false)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--common-name", "common.name.io", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including organization", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"organization":"organization.io"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization", "organization.io")
			Eventually(session).Should(Exit(0))
		})

		It("including organization unit", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"organization_unit":"My Unit"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--organization-unit", "My Unit")
			Eventually(session).Should(Exit(0))
		})

		It("including locality", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"locality":"My Locality"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--locality", "My Locality")
			Eventually(session).Should(Exit(0))
		})

		It("including state", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"state":"My State"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--state", "My State")
			Eventually(session).Should(Exit(0))
		})

		It("including country", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"country":"My Country"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--country", "My Country")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple alternative names", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"alternative_names": [ "Alt1", "Alt2" ]}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--alternative-name", "Alt1", "--alternative-name", "Alt2")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple extended key usage settings", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"extended_key_usage": [ "server_auth", "client_auth" ]}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-e", "server_auth", "--ext-key-usage=client_auth")
			Eventually(session).Should(Exit(0))
		})

		It("including multiple key usage settings", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"key_usage": ["digital_signature", "non_repudiation"]}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-g", "digital_signature", "--key-usage=non_repudiation")
			Eventually(session).Should(Exit(0))
		})

		It("including key length", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"key_length":2048}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--key-length", "2048")
			Eventually(session).Should(Exit(0))
		})

		It("including duration", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"duration":1000}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--duration", "1000")
			Eventually(session).Should(Exit(0))
		})

		It("including certificate authority", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"ca":"my_ca"}`, true)
			setupCertificatePostServer("my-secret", "potatoes-ca", "potatoes-cert", "potatoes-priv-key", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "--ca", "my_ca")
			Eventually(session).Should(Exit(0))
		})

		It("including self-signed flag", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"self_sign": true, "common_name": "my.name.io"}`, true)
			setupCertificatePostServer("my-secret", "", "", "", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-c", "my.name.io", "--self-sign")
			Eventually(session).Should(Exit(0))
		})

		It("including is-ca flag", func() {
			expectedRequestJson := generateRequestJson("certificate", "my-secret", `{"is_ca": true, "common_name": "my.name.io"}`, true)
			setupCertificatePostServer("my-secret", "", "", "", expectedRequestJson)
			session := runCommand("generate", "-n", "my-secret", "-t", "certificate", "-c", "my.name.io", "--is-ca")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("with a variety of user parameters", func() {
		name := "my-username-credential"
		It("prints the secret", func() {
			expectedRequestJson := generateRequestJson("user", name, `{}`, true)
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				expectedRequestJson)

			session := runCommand("generate", "-n", name, "-t", "user")

			Eventually(session).Should(Exit(0))
			Expect(session.Out.Contents()).To(ContainSubstring(responseMyUsernameYaml))
		})

		It("should accept a statically provided username", func() {
			expectedRequestJson := generateUserRequestJson(name, `{}`, `{"username": "my-username"}`, true)
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				expectedRequestJson)

			session := runCommand("generate", "-n", name, "-t", "user", "-z", "my-username")

			Eventually(session).Should(Exit(0))
			Expect(session.Out.Contents()).To(ContainSubstring(responseMyUsernameYaml))
		})

		It("with with no-overwrite", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{}`, false))
			session := runCommand("generate", "-n", name, "-t", "user", "--no-overwrite")
			Eventually(session).Should(Exit(0))
		})

		It("including length", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{"length": 42}`, true))
			session := runCommand("generate", "-n", name, "-t", "user", "-l", "42")
			Eventually(session).Should(Exit(0))
		})

		It("excluding upper case", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{"exclude_upper": true}`, true))
			session := runCommand("generate", "-n", name, "-t", "user", "--exclude-upper")
			Eventually(session).Should(Exit(0))
		})

		It("excluding lower case", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{"exclude_lower": true}`, true))
			session := runCommand("generate", "-n", name, "-t", "user", "--exclude-lower")
			Eventually(session).Should(Exit(0))
		})

		It("including special characters", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{"include_special": true}`, true))
			session := runCommand("generate", "-n", name, "-t", "user", "--include-special")
			Eventually(session).Should(Exit(0))
		})

		It("excluding numbers", func() {
			setupUserPostServer(
				name,
				"my-username",
				"test-password",
				"passw0rd-H4$h",
				generateRequestJson("user", name, `{"exclude_number": true}`, true))
			session := runCommand("generate", "-n", name, "-t", "user", "--exclude-number")
			Eventually(session).Should(Exit(0))
		})
	})

	Describe("Help", func() {
		ItBehavesLikeHelp("generate", "n", func(session *Session) {
			Expect(session.Err).To(Say("generate"))
			Expect(session.Err).To(Say("name"))
			Expect(session.Err).To(Say("length"))
		})

		It("short flags", func() {
			Expect(commands.GenerateCommand{}).To(SatisfyAll(
				commands.HaveFlag("name", "n"),
				commands.HaveFlag("type", "t"),
				commands.HaveFlag("no-overwrite", "O"),
				commands.HaveFlag("length", "l"),
				commands.HaveFlag("include-special", "S"),
				commands.HaveFlag("exclude-number", "N"),
				commands.HaveFlag("exclude-upper", "U"),
				commands.HaveFlag("exclude-lower", "L"),
				commands.HaveFlag("common-name", "c"),
				commands.HaveFlag("organization", "o"),
				commands.HaveFlag("organization-unit", "u"),
				commands.HaveFlag("locality", "i"),
				commands.HaveFlag("state", "s"),
				commands.HaveFlag("country", "y"),
				commands.HaveFlag("alternative-name", "a"),
				commands.HaveFlag("key-length", "k"),
				commands.HaveFlag("duration", "d"),
			))
		})

		It("displays missing 'n' option as required parameters", func() {
			session := runCommand("generate")

			Eventually(session).Should(Exit(1))

			if runtime.GOOS == "windows" {
				Expect(session.Err).To(Say("the required flag `/n, /name' was not specified"))
			} else {
				Expect(session.Err).To(Say("the required flag `-n, --name' was not specified"))
			}
		})

		It("displays the server provided error when an error is received", func() {
			server.AppendHandlers(
				RespondWith(http.StatusBadRequest, `{"error": "test error"}`),
			)

			session := runCommand("generate", "-n", "my-value", "-t", "value")

			Eventually(session).Should(Exit(1))

			Expect(session.Err).To(Say("test error"))
		})
	})
})

func setupUserPostServer(name, username, password, passwordHash, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(USER_CREDENTIAL_RESPONSE_JSON, name, username, password, passwordHash)),
		),
	)
}

func setupPasswordPostServer(name, value, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(STRING_CREDENTIAL_RESPONSE_JSON, "password", name, value)),
		),
	)
}

func setupRsaSshPostServer(name, credentialType, publicKey, privateKey, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(RSA_SSH_CREDENTIAL_RESPONSE_JSON, credentialType, name, publicKey, privateKey)),
		),
	)
}

func setupCertificatePostServer(name, ca, certificate, privateKey, requestJson string) {
	server.AppendHandlers(
		CombineHandlers(
			VerifyRequest("POST", "/api/v1/data"),
			VerifyJSON(requestJson),
			RespondWith(http.StatusOK, fmt.Sprintf(CERTIFICATE_CREDENTIAL_RESPONSE_JSON, name, ca, certificate, privateKey)),
		),
	)
}

func generateRequestJson(credentialType, name, params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_CREDENTIAL_REQUEST_JSON, name, credentialType, overwrite, params)
}

func generateUserRequestJson(name, params, value string, overwrite bool) string {
	return fmt.Sprintf(USER_GENERATE_CREDENTIAL_REQUEST_JSON, name, overwrite, params, value)
}

func generateDefaultTypeRequestJson(name, params string, overwrite bool) string {
	return fmt.Sprintf(GENERATE_DEFAULT_TYPE_REQUEST_JSON, name, overwrite, params)
}
