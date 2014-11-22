package configserver

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/pivotal-golang/lager"
)

type server struct {
	logger lager.Logger

	db       ConfigDB
	validate ConfigValidator
}

type ConfigDB interface {
	db.ConfigDB
}

type ConfigValidator func(atc.Config) error

func NewServer(
	logger lager.Logger,
	db ConfigDB,
	validator ConfigValidator,
) *server {
	return &server{
		logger:   logger,
		db:       db,
		validate: validator,
	}
}
