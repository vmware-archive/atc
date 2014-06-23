package grouper_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper/0.2"
	"github.com/tedsuo/ifrit/test_helpers"
)

var _ = Describe("RestartPolicies", func() {
	var pinger1 test_helpers.PingChan
	var pinger2 test_helpers.PingChan
	var pGroup ifrit.Process

	BeforeEach(func() {
		pinger1 = make(test_helpers.PingChan)
		pinger2 = make(test_helpers.PingChan)
		pGroup = ifrit.Envoke(grouper.Members{
			{"ping1", pinger1, grouper.RestartMePolicy()},
			{"ping2", pinger2, grouper.StopGroupPolicy()},
		})
	})

	Describe("RestartMe", func() {
		Context("when the member exits", func() {
			BeforeEach(func() {
				<-pinger1
			})

			It("does not exit the group", func() {
				Eventually(pGroup.Wait()).ShouldNot(Receive())
			})

			It("restarts the process", func() {
				Eventually(pinger1).Should(Receive())
			})

			Context("and then the group is sent a Signal", func() {
				BeforeEach(func() {
					Eventually(pinger1).Should(Receive())
					pGroup.Signal(os.Kill)
				})

				It("exits the group", func() {
					Eventually(pGroup.Wait()).Should(Receive())
				})
			})
		})
	})

	Describe("StopGroup", func() {
		BeforeEach(func() {
			<-pinger2
		})

		It("exits the group", func() {
			Eventually(pGroup.Wait()).Should(Receive())
		})

		It("does not restart itself", func() {
			Consistently(pinger2).ShouldNot(Receive())
		})

		It("stops all other processes", func() {
			Eventually(pinger1).Should(Receive())
			Consistently(pinger1).ShouldNot(Receive())
		})
	})
})
