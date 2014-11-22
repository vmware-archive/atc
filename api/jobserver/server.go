package jobserver

import (
	"github.com/pivotal-golang/lager"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

type server struct {
	logger lager.Logger

	db       JobsDB
	configDB ConfigDB
}

type JobsDB interface {
	GetAllJobBuilds(job string) ([]db.Build, error)
	GetCurrentBuild(job string) (db.Build, error)
	GetJobBuild(job string, build string) (db.Build, error)
	GetJobFinishedAndNextBuild(job string) (*db.Build, *db.Build, error)
}

type ConfigDB interface {
	GetConfig() (atc.Config, error)
}

func NewServer(
	logger lager.Logger,
	db JobsDB,
	configDB ConfigDB,
) *server {
	return &server{
		logger:   logger,
		db:       db,
		configDB: configDB,
	}
}
