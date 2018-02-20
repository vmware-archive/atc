package jobserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/auth"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/scheduler"
)

//go:generate counterfeiter . SchedulerFactory

type SchedulerFactory interface {
	BuildScheduler(db.Pipeline, string, creds.Variables) scheduler.BuildScheduler
}

type Server struct {
	logger lager.Logger

	schedulerFactory SchedulerFactory
	externalURL      string
	rejector         auth.Rejector
	variablesFactory creds.VariablesFactory
	accessorFactory  accessor.AccessorFactory
}

func NewServer(
	logger lager.Logger,
	schedulerFactory SchedulerFactory,
	externalURL string,
	variablesFactory creds.VariablesFactory,
	accessorFactory accessor.AccessorFactory,
) *Server {
	return &Server{
		logger:           logger,
		schedulerFactory: schedulerFactory,
		externalURL:      externalURL,
		rejector:         auth.UnauthorizedRejector{},
		variablesFactory: variablesFactory,
		accessorFactory:  accessorFactory,
	}
}
