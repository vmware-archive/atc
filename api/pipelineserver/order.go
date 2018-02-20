package pipelineserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/api/accessor"
)

func (s *Server) OrderPipelines(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("order-pipelines")

	teamName := r.FormValue(":team_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	team, err := acc.Team(accessor.Write, teamName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	var pipelineNames []string
	if err := json.NewDecoder(r.Body).Decode(&pipelineNames); err != nil {
		logger.Error("invalid-json", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// team, found, err := s.teamFactory.FindTeam(teamName)
	// if err != nil {
	// 	logger.Error("failed-to-get-team", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// if !found {
	// 	logger.Info("team-not-found")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	err = team.OrderPipelines(pipelineNames)
	if err != nil {
		logger.Error("failed-to-order-pipelines", err, lager.Data{
			"pipeline-names": pipelineNames,
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
