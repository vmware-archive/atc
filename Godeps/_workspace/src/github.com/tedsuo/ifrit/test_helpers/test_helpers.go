package test_helpers

import (
	"errors"
	"os"

	"github.com/tedsuo/ifrit"
)

type Ping struct{}

var PingerExitedFromPing = errors.New("pinger exited with a ping")
var PingerExitedFromSignal = errors.New("pinger exited with a signal")

type PingChan chan Ping

func (p PingChan) Run(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)
	select {
	case <-sigChan:
		return PingerExitedFromSignal
	case p <- Ping{}:
		return PingerExitedFromPing
	}
}

var NoReadyExitedNormally = errors.New("no ready exited normally")

var NoReadyRunner = ifrit.RunFunc(func(sigChan <-chan os.Signal, ready chan<- struct{}) error {
	return NoReadyExitedNormally
})
