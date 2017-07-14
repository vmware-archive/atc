package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Import", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get", "-n", "test-credential")

	Describe("importing a file with mixed credentials", func() {
		It("sets the all credentials", func() {
			setUpImportRequests()

			session := runCommand("import", "-f", "../test/test_import_file.yml")

			Eventually(session).Should(Exit(0))

			Eventually(session.Out).Should(Say(`name: /test/password
type: password
value: test-password-value`))
			Eventually(session.Out).Should(Say(`name: /test/value
type: value
value: test-value`))
			Eventually(session.Out).Should(Say(`name: /test/certificate
type: certificate
value:
  ca: ca-certificate
  certificate: certificate
  private_key: private-key`))
			Eventually(session.Out).Should(Say(`name: /test/rsa
type: rsa
value:
  private_key: private-key
  public_key: public-key`))
			Eventually(session.Out).Should(Say(`name: /test/ssh
type: ssh
value:
  private_key: private-key
  public_key: ssh-public-key`))
			Eventually(session.Out).Should(Say(`name: /test/user
type: user
value:
  password: test-user-password
  password_hash: P455W0rd-H45H
  username: covfefe`))
			Eventually(session.Out).Should(Say(`name: /test/json
type: json
value:
  "1": key is not a string
  "3.14": pi
  arbitrary_object:
    nested_array:
    - array_val1
    - array_object_subvalue: covfefe
  "true": key is a bool
`))
		})
	})

	Describe("when importing file with no name specified", func() {
		It("passes through the server error", func() {
			jsonBody := `{"type":"password","value":"test-password","overwrite":true}`
			SetupPutBadRequestServer(jsonBody)

			session := runCommand("import", "-f", "../test/test_import_missing_name.yml")

			Eventually(session.Err).Should(Say(`test error`))
		})
	})

	Describe("when importing file with incorrect YAML", func() {
		It("returns an error message", func() {
			errorMessage := `The referenced file does not contain valid yaml structure. Please update and retry your request.`

			session := runCommand("import", "-f", "../test/test_import_incorrect_yaml.yml")

			Eventually(session.Err).Should(Say(errorMessage))
		})
	})
})

func setUpImportRequests() {
	SetupOverwritePutValueServer("/test/password", "password", "test-password-value", true)
	SetupOverwritePutValueServer("/test/value", "value", "test-value", true)
	SetupPutCertificateServer("/test/certificate",
		`ca-certificate`,
		`certificate`,
		`private-key`)
	SetupPutRsaSshServer("/test/rsa", "rsa", "public-key", "private-key", true)
	SetupPutRsaSshServer("/test/ssh", "ssh", "ssh-public-key", "private-key", true)
	SetupPutUserServer("/test/user", `{"username": "covfefe", "password": "test-user-password"}`, "covfefe", "test-user-password", "P455W0rd-H45H", true)
	setupPutJsonServer("/test/json", `{"1":"key is not a string","3.14":"pi","true":"key is a bool","arbitrary_object":{"nested_array":["array_val1",{"array_object_subvalue":"covfefe"}]}}`)
}
