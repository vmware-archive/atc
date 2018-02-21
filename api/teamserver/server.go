package teamserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/db"
)

type Server struct {
	logger          lager.Logger
	teamFactory     db.TeamFactory
	accessorFactory accessor.AccessorFactory
}

func NewServer(
	logger lager.Logger,
	teamFactory db.TeamFactory,
	accessorFactory accessor.AccessorFactory,
) *Server {
	return &Server{
		logger:          logger,
		teamFactory:     teamFactory,
		accessorFactory: accessorFactory,
	}
}
