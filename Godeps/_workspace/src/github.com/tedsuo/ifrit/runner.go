package ifrit

import "os"

type Runner interface {
	Run(signals <-chan os.Signal, ready chan<- struct{}) error
}

type RunFunc func(signals <-chan os.Signal, ready chan<- struct{}) error

func (r RunFunc) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	return r(signals, ready)
}
