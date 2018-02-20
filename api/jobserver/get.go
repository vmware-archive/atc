package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) GetJob(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("get-job")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	jobName := r.FormValue(":job_name")

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

	job, found, err := pipeline.Job(jobName)
	if err != nil {
		logger.Error("failed-to-get-resource-types", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	finished, next, err := job.FinishedAndNextBuild()
	if err != nil {
		logger.Error("could-not-get-job-finished-and-next-build", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(present.Job(
		teamName,
		job,
		pipeline.Groups(),
		finished,
		next,
		nil,
	))
	if err != nil {
		logger.Error("failed-to-encode-job", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
