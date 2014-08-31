package ifrit

import "os"

type Process interface {
	Wait() <-chan error
	Signal(os.Signal)
}

func Envoke(r Runner) Process {
	return envokeProcess(r)
}

func envokeProcess(r Runner) Process {
	p := &process{
		runner: r,
		sig:    make(chan os.Signal),
		ready:  make(chan struct{}),
		exited: make(chan struct{}),
	}

	go p.run()

	select {
	case <-p.ready:
	case <-p.Wait():
	}

	return p
}

type process struct {
	runner     Runner
	sig        chan os.Signal
	ready      chan struct{}
	exited     chan struct{}
	exitStatus error
}

func (p *process) run() {
	p.exitStatus = p.runner.Run(p.sig, p.ready)
	close(p.exited)
}

func (p *process) Wait() <-chan error {
	exitChan := make(chan error, 1)

	go func() {
		<-p.exited
		exitChan <- p.exitStatus
	}()

	return exitChan
}

func (p *process) Signal(signal os.Signal) {
	go func() {
		select {
		case p.sig <- signal:
		case <-p.exited:
		}
	}()
}
