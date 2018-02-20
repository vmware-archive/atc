package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) GetJobBuild(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("get-job-build")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	jobName := r.FormValue(":job_name")
	buildName := r.FormValue(":build_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	build, err := acc.TeamPipelineJobBuild(accessor.Read, teamName, pipelineName, jobName, buildName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(present.Build(build))
	if err != nil {
		logger.Error("failed-to-encode-build", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
