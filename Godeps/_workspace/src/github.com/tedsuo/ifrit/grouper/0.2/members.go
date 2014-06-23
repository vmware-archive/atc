package grouper

import (
	"fmt"
	"os"

	"github.com/tedsuo/ifrit"
)

type Loader interface {
	Load(error) (ifrit.Runner, bool)
}

type Members []Member

type Member struct {
	Name    string
	Runner  ifrit.Runner
	Restart Restart
}

func (members Members) Load(err error) (ifrit.Runner, bool) {
	return members, true
}

func (members Members) Run(sig <-chan os.Signal, ready chan<- struct{}) error {

	group := make(pGroup, len(members))

	startedChan := make(pMemberChan, len(members))
	startedChan.envokeGroup(members, group)

	exitedChan := make(exitedChannel, len(group))
	exitedChan.waitForGroup(group)

	var errToReturn error
	desiredCount := len(group)
	signaledToStop := false

	if ready != nil {
		close(ready)
	}

	for {
		if desiredCount == 0 {
			return errToReturn
		}

		select {

		case signal := <-sig:
			signaledToStop = true
			group.Signal(signal)

		case pm := <-startedChan:
			group[pm.Process] = pm
			go exitedChan.waitForProcess(pm.Process)

		case e := <-exitedChan:
			member, ok := group[e.Process]
			if !ok {
				panic(fmt.Errorf("Exit for missing process! \nExit: \nErr:%s \n Process: %#v", e.Process))
			}

			delete(group, e.Process)
			desiredCount--
			restart := member.Restart

			if restart.Signal != Continue {
				group.Signal(restart.Signal)
			}

			if signaledToStop {
				continue
			}

			if e.error != nil {
				errToReturn = fmt.Errorf("%s exited with error: %s", member.Name, e.error)
			}

			if !restart.AttemptRestart {
				if restart.Signal != Continue {
					signaledToStop = true
				}
				continue
			}

			loader, ok := member.Runner.(Loader)
			if !ok {
				continue
			}

			nextRunner, ok := loader.Load(e.error)
			if !ok {
				continue
			}

			desiredCount++
			newMember := Member{member.Name, nextRunner, member.Restart}
			go startedChan.envokeMember(newMember)
		}
	}
}

type exit struct {
	ifrit.Process
	error
}

type exitedChannel chan exit

func (exitedChan exitedChannel) waitForGroup(group pGroup) {
	for p, _ := range group {
		go exitedChan.waitForProcess(p)
	}
}

func (exitedChan exitedChannel) waitForProcess(p ifrit.Process) {
	err := <-p.Wait()
	exitedChan <- exit{p, err}
}

type pGroup map[ifrit.Process]pMember

func (group pGroup) Signal(signal os.Signal) {
	for p, _ := range group {
		p.Signal(signal)
	}
}

type pMember struct {
	ifrit.Process
	Member
}

type pMemberChan chan pMember

func (pmChan pMemberChan) envokeMember(member Member) {
	process := ifrit.Envoke(member.Runner)
	pmChan <- pMember{process, member}
}

func (pmChan pMemberChan) envokeGroup(group Members, p pGroup) {
	for _, member := range group {
		go pmChan.envokeMember(member)
	}

	for _ = range group {
		pm := <-pmChan
		p[pm.Process] = pm
	}
}
