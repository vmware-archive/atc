package resourceserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) GetResource(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("get-resource")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	resourceName := r.FormValue(":resource_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	resource, err := acc.TeamPipelineResource(accessor.Read, teamName, pipelineName, resourceName)

	if err != nil {
		logger.Error("failed-to-get-resource", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(present.Resource(
		resource,
		teamName,
	))

	if err != nil {
		logger.Error("failed-to-encode-resource", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
