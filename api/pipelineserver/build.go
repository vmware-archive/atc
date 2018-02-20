package pipelineserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) CreateBuild(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("create-build")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

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

	var plan atc.Plan
	err = json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		logger.Info("malformed-request", lager.Data{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	build, err := pipeline.CreateOneOffBuild()
	if err != nil {
		logger.Error("failed-to-create-one-off-build", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	engineBuild, err := s.engine.CreateBuild(logger, build, plan)
	if err != nil {
		logger.Error("failed-to-start-build", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go engineBuild.Resume(logger)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(present.Build(build))
	if err != nil {
		logger.Error("failed-to-encode-build", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
