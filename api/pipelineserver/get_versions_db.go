package pipelineserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
)

func (s *Server) GetVersionsDB(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("get-version-db")

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

	versionsDB, _ := pipeline.LoadVersionsDB()
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(versionsDB)
	if err != nil {
		logger.Error("failed-to-encode-version-db", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
