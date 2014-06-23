package grouper

import (
	"os"
	"time"
)

const (
	NoTimeout           = time.Duration(0)
	DefaultStopTimeout  = 10 * time.Second
	DefaultGraceTimeout = 5 * time.Minute
)

type Restart struct {
	AttemptRestart bool
	Signal         os.Signal
	Timeout        time.Duration
}

type RestartPolicy func() Restart

var (
	RestartMePolicy = RestartPolicy(func() Restart {
		return Restart{true, Continue, NoTimeout}
	})

	StopMePolicy = RestartPolicy(func() Restart {
		return Restart{false, Continue, NoTimeout}
	})

	RestartGroupPolicy = RestartPolicy(func() Restart {
		return Restart{true, os.Interrupt, DefaultGraceTimeout}
	})

	StopGroupPolicy = RestartPolicy(func() Restart {
		return Restart{false, os.Interrupt, DefaultGraceTimeout}
	})
)
