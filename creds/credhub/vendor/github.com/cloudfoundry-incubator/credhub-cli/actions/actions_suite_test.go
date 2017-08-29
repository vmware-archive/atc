package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	test_util "github.com/cloudfoundry-incubator/credhub-cli/test"
)

func TestActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Actions Suite")
}

var _ = BeforeSuite(func() {
	test_util.CleanEnv()
})
