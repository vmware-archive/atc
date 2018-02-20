package pipelineserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/auth"
	"github.com/concourse/atc/engine"
)

type Server struct {
	logger          lager.Logger
	rejector        auth.Rejector
	accessorFactory accessor.AccessorFactory
	engine          engine.Engine
}

func NewServer(
	logger lager.Logger,
	accessorFactory accessor.AccessorFactory,
	engine engine.Engine,
) *Server {
	return &Server{
		logger:          logger,
		rejector:        auth.UnauthorizedRejector{},
		accessorFactory: accessorFactory,
		engine:          engine,
	}
}
