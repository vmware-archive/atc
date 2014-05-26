package grouper

import (
	"os"

	"github.com/tedsuo/ifrit"
)

type RunGroup map[string]ifrit.Runner

func (r RunGroup) Run(sig <-chan os.Signal, ready chan<- struct{}) error {
	p := EnvokeGroup(r)

	if ready != nil {
		close(ready)
	}

	for {
		select {
		case signal := <-sig:
			p.Signal(signal)
		case err := <-p.Wait():
			return err
		}
	}
}
