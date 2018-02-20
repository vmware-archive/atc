package pipelineserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) GetPipeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger := s.logger.Session("get-pipeline")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipeline, err := acc.TeamPipeline(accessor.Read, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(present.Pipeline(pipeline))
	if err != nil {
		logger.Error("failed-to-encode-pipeline", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
