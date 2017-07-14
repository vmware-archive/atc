// +build !windows

package config_test

import (
	"github.com/cloudfoundry-incubator/credhub-cli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var cfg config.Config

	BeforeEach(func() {
		cfg = config.Config{
			ApiURL:  "http://api.example.com",
			AuthURL: "http://auth.example.com",
		}
	})

	It("places the config file in .cm in the home directory", func() {
		Expect(config.ConfigPath()).To(HaveSuffix(`/.credhub/config.json`))
	})
})
