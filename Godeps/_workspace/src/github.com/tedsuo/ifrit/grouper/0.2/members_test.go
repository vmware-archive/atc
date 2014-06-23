package grouper_test

import (
	"os"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper/0.2"
	"github.com/tedsuo/ifrit/test_helpers"
)

var _ = Describe("Members", func() {
	Describe("ifrit.Envoke()", func() {
		var pGroup ifrit.Process

		BeforeEach(func() {
			pinger1 := make(test_helpers.PingChan)
			pinger2 := make(test_helpers.PingChan)
			pinger3 := make(test_helpers.PingChan)
			pinger4 := make(test_helpers.PingChan)
			pGroup = ifrit.Envoke(grouper.Members{
				{"ping1", pinger1, grouper.RestartMePolicy()},
				{"ping2", pinger2, grouper.StopMePolicy()},
				{"ping3", pinger3, grouper.RestartGroupPolicy()},
				{"ping4", pinger4, grouper.StopGroupPolicy()},
			})
		})

		Context("when the linked process is signaled to stop", func() {
			BeforeEach(func() {
				pGroup.Signal(os.Kill)
			})

			It("exits with nil", func() {
				Eventually(pGroup.Wait()).Should(Receive(nil))
			})
		})
	})

	Describe("Signal propogation", func() {
		var recorder1 *test_helpers.SignalRecoder
		var recorder2 *test_helpers.SignalRecoder
		var pGroup ifrit.Process

		BeforeEach(func() {
			recorder1 = test_helpers.NewSignalRecorder(syscall.SIGVTALRM)
			recorder2 = test_helpers.NewSignalRecorder()
			pGroup = ifrit.Envoke(grouper.Members{
				{"recorder1", recorder1, grouper.StopGroupPolicy()},
				{"recorder2", recorder2, grouper.StopGroupPolicy()},
			})

			pGroup.Signal(syscall.SIGVTALRM)
			Eventually(pGroup.Wait()).Should(Receive())
		})

		It("should propogate the initial signal, and the first processes exit signal", func() {
			Ω(recorder2.ReceivedSignals()).Should(Equal([]os.Signal{syscall.SIGVTALRM, os.Interrupt}))
		})
	})
})
