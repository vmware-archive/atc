package resourceserver

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
)

func (s *Server) UnpauseResource(w http.ResponseWriter, r *http.Request) {

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	resourceName := r.FormValue(":resource_name")

	logger := s.logger.Session("unpause-resource", lager.Data{
		"resource": resourceName,
	})

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	resource, err := acc.TeamPipelineResource(accessor.Write, teamName, pipelineName, resourceName)
	if err != nil {
		logger.Error("failed-to-get-resource", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	err = resource.Unpause()
	if err != nil {
		logger.Error("failed-to-unpause", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
