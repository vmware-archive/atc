package pipelineserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

// show all public pipelines and team private pipelines if authorized
func (s *Server) ListAllPipelines(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-all-pipelines")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipelines, err := acc.Pipelines(accessor.Read)
	if err != nil {
		logger.Error("failed-to-get-pipelines", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(present.Pipelines(pipelines))
	if err != nil {
		logger.Error("failed-to-encode-pipelines", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
