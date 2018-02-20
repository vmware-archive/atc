package resourceserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) ListResources(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-resources")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	resources, err := acc.TeamPipelineResources(accessor.Read, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	var presentedResources []atc.Resource
	for _, resource := range resources {
		presentedResources = append(
			presentedResources,
			present.Resource(
				resource,
				teamName,
			),
		)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(presentedResources)
	if err != nil {
		logger.Error("failed-to-encode-resources", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
