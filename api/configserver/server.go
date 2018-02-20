package configserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
)

type Server struct {
	logger          lager.Logger
	accessorFactory accessor.AccessorFactory
}

func NewServer(
	logger lager.Logger,
	accessorFactory accessor.AccessorFactory,
) *Server {
	return &Server{
		logger:          logger,
		accessorFactory: accessorFactory,
	}
}
