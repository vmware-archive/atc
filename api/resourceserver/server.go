package resourceserver

import (
	"github.com/pivotal-golang/lager"

	"github.com/concourse/atc"
)

type server struct {
	logger lager.Logger

	configDB ConfigDB
}

type ConfigDB interface {
	GetConfig() (atc.Config, error)
}

func NewServer(
	logger lager.Logger,
	configDB ConfigDB,
) *server {
	return &server{
		logger:   logger,
		configDB: configDB,
	}
}
