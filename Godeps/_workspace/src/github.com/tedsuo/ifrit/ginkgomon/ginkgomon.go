package runner

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type Config struct {
	Name              string
	BinPath           string
	AnsiColorCode     string
	StartCheck        string
	StartCheckTimeout time.Duration
	Args              []string
}

type Runner struct {
	config Config
}

func New(config Config) *GexecRunner {
	return &Runner{
		config: Config,
	}
}

func (r *Runner) Run(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	session, err := gexec.Start(
		exec.Command(
			r.BinPath,
			r.Args...,
		),
		gexec.NewPrefixedWriter("\x1b[32m[o]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name, ginkgo.GinkgoWriter),
		gexec.NewPrefixedWriter("\x1b[91m[e]\x1b[%s[%s]\x1b[0m ", r.AnsiColorCode, r.Name, ginkgo.GinkgoWriter),
	)

	Ω(err).ShouldNot(HaveOccurred())

	if r.StartCheck != "" {
		Eventually(r.Session, r.StartCheckTimeout).Should(gbytes.Say(r.StartCheck))
	}

	close(ready)

	var signal os.Signal

	for {
		select {

		case signal = <-sigChan:
			session.Signal(signal)

		case <-session.Exited:
			if session.ExitCode() == 0 {
				return nil
			}
			return fmt.Errorf("exit status %d", session.ExitCode())
		}
	}
}
