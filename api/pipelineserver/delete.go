package pipelineserver

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
)

func (s *Server) DeletePipeline(w http.ResponseWriter, r *http.Request) {
	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

	logger := s.logger.Session("destroying-pipeline", lager.Data{
		"name": pipelineName,
	})

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipeline, err := acc.TeamPipeline(accessor.Write, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	logger.Info("start")

	err = pipeline.Destroy()
	if err != nil {
		logger.Error("failed-to-destroy-pipeline", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("done")

	w.WriteHeader(http.StatusNoContent)
}
