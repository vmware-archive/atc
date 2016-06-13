package pipelineserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListPipelines(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-pipelines")
	teamName := r.FormValue(":team_name")

	pipelines, err := s.getPipelinesForTeam(teamName)
	if err != nil {
		logger.Error("failed-to-get-all-active-pipelines", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(pipelines)
}
