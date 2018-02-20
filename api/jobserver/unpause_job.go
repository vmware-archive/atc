package jobserver

import (
	"net/http"

	"github.com/concourse/atc/api/accessor"
)

func (s *Server) UnpauseJob(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("unpause-job")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	jobName := r.FormValue(":job_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	job, err := acc.TeamPipelineJob(accessor.Write, teamName, pipelineName, jobName)
	if err != nil {
		logger.Error("failed-to-get-job", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	err = job.Unpause()
	if err != nil {
		logger.Error("failed-to-unpause-job", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
