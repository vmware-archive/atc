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

type Config struct {
	Command           *exec.Cmd
	Name              string
	AnsiColorCode     string
	StartCheck        string
	StartCheckTimeout time.Duration
	Cleanup           func()
}

type Runner struct {
	Command           *exec.Cmd
	Name              string
	AnsiColorCode     string
	StartCheck        string
	StartCheckTimeout time.Duration
	Cleanup           func()
	session           *gexec.Session
	sessionReady      chan struct{}
}

func New(config Config) *Runner {
	return &Runner{
		Name:              config.Name,
		Command:           config.Command,
		AnsiColorCode:     config.AnsiColorCode,
		StartCheck:        config.StartCheck,
		StartCheckTimeout: config.StartCheckTimeout,
		Cleanup:           config.Cleanup,
		sessionReady:      make(chan struct{}),
	}
}

func (r *Runner) Buffer() *gbytes.Buffer {
	if r.sessionReady == nil {
		panic("ginkgomon improperly created without using New")
	}
	<-r.sessionReady
	return r.session.Buffer()
}

func (r *Runner) Run(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	defer ginkgo.GinkgoRecover()

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
			timeout = 5 * time.Second
		}

		Eventually(allOutput, timeout).Should(gbytes.Say(r.StartCheck))
	}

	r.session = session
	if r.sessionReady != nil {
		close(r.sessionReady)
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
