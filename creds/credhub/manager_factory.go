package credhub

import (
	"github.com/concourse/atc/creds"
	flags "github.com/jessevdk/go-flags"
)

type CredhubManagerFactory struct{}

func init() {
	creds.Register("credhub", NewCredhubManagerFactory())
}

func NewCredhubManagerFactory() creds.ManagerFactory {
	return &CredhubManagerFactory{}
}

func (factory *CredhubManagerFactory) AddConfig(group *flags.Group) creds.Manager {
	manager := &CredhubManager{}

	subGroup, err := group.AddGroup("Credhub Credential Management", "", manager)
	if err != nil {
		panic(err)
	}

	subGroup.Namespace = "credhub"

	return manager
}
