package grouper_test

import (
	"os"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/test_helpers"
)

var _ = Describe("ProcessGroup", func() {
	var rGroup grouper.RunGroup

	BeforeEach(func() {
		rGroup = grouper.RunGroup{
			"ping1": make(test_helpers.PingChan),
			"ping2": make(test_helpers.PingChan),
		}
	})

	Describe("EnvokeGroup()", func() {
		var pGroup grouper.ProcessGroup

		BeforeEach(func() {
			pGroup = grouper.EnvokeGroup(rGroup)
		})

		Context("when the linked process is signaled to stop", func() {
			BeforeEach(func() {
				pGroup.Signal(os.Kill)
			})

			It("exits with a combined error message", func() {
				err := <-pGroup.Wait()
				errMsg := err.Error()
				Ω(errMsg).Should(ContainSubstring("ping1"))
				Ω(errMsg).Should(ContainSubstring("ping2"))
				Ω(errMsg).Should(ContainSubstring(test_helpers.PingerExitedFromSignal.Error()))
			})

			It("emits an exit for each member", func(done Done) {
				members := []grouper.Member{}
				memChan := pGroup.Exits()
				for i := 0; i < len(rGroup); i++ {
					members = append(members, <-memChan)
				}
				Ω(members).Should(HaveLen(len(rGroup)))
				close(done)
			})
		})
	})

	Describe("ifrit.Envoke()", func() {
		var pGroup ifrit.Process

		BeforeEach(func() {
			pGroup = ifrit.Envoke(rGroup)
		})

		Context("when the linked process is signaled to stop", func() {
			BeforeEach(func() {
				pGroup.Signal(os.Kill)
			})

			It("exits with a combined error message", func() {
				err := <-pGroup.Wait()
				errMsg := err.Error()
				Ω(errMsg).Should(ContainSubstring("ping1"))
				Ω(errMsg).Should(ContainSubstring("ping2"))
				Ω(errMsg).Should(ContainSubstring(test_helpers.PingerExitedFromSignal.Error()))
			})
		})
	})
})
