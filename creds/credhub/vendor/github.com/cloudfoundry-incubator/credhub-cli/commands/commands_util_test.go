package commands_test

import (
	"os"

	"github.com/cloudfoundry-incubator/credhub-cli/commands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	credhub_errors "github.com/cloudfoundry-incubator/credhub-cli/errors"
)

var _ = Describe("Util", func() {
	Describe("#ReadFile", func() {
		It("reads a file into memory", func() {
			tempDir := createTempDir("filesForTesting")
			fileContents := "My Test String"
			filename := createCredentialFile(tempDir, "file.txt", fileContents)
			readContents, err := commands.ReadFile(filename)
			Expect(readContents).To(Equal(fileContents))
			Expect(err).To(BeNil())
			os.RemoveAll(tempDir)
		})

		It("returns an error message if a file cannot be read", func() {
			readContents, err := commands.ReadFile("Foo")
			Expect(readContents).To(Equal(""))
			Expect(err).To(MatchError(credhub_errors.NewFileLoadError()))
		})
	})

	Describe("#AddDefaultSchemeIfNecessary", func() {
		It("adds the default scheme (https://) to a server which has none", func() {
			transformedUrl := commands.AddDefaultSchemeIfNecessary("foo.com:8080")
			Expect(transformedUrl).To(Equal("https://foo.com:8080"))
		})

		It("does not add the default scheme if one is already there", func() {
			transformedUrl := commands.AddDefaultSchemeIfNecessary("ftp://foo.com:8080")
			Expect(transformedUrl).To(Equal("ftp://foo.com:8080"))
		})
	})
})
