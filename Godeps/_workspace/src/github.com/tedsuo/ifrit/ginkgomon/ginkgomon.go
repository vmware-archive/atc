package ginkgomon

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type Runner struct {
	Command           *exec.Cmd
	Name              string
	AnsiColorCode     string
	StartCheck        string
	StartCheckTimeout time.Duration
	Cleanup           func()
}

func (r *Runner) Run(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	allOutput := gbytes.NewBuffer()

	session, err := gexec.Start(
		r.Command,
		gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name),
			io.MultiWriter(allOutput, ginkgo.GinkgoWriter),
		),
		gexec.NewPrefixedWriter(
			fmt.Sprintf("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name),
			io.MultiWriter(allOutput, ginkgo.GinkgoWriter),
		),
	)

	Ω(err).ShouldNot(HaveOccurred())

	if r.StartCheck != "" {
		timeout := r.StartCheckTimeout
		if timeout == 0 {
			timeout = time.Second
		}

		Eventually(allOutput, timeout).Should(gbytes.Say(r.StartCheck))
	}

	close(ready)

	var signal os.Signal

	for {
		select {
		case signal = <-sigChan:
			session.Signal(signal)

		case <-session.Exited:
			if r.Cleanup != nil {
				r.Cleanup()
			}

			if session.ExitCode() == 0 {
				return nil
			}

			return fmt.Errorf("exit status %d", session.ExitCode())
		}
	}
}
