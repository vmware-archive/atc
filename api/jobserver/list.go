package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) ListJobs(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-jobs")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

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

	var jobs []atc.Job

	include := r.FormValue("include")
	dashboard, groups, err := pipeline.Dashboard(include)

	if err != nil {
		logger.Error("failed-to-get-dashboard", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, job := range dashboard {
		jobs = append(
			jobs,
			present.Job(
				teamName,
				job.Job,
				groups,
				job.FinishedBuild,
				job.NextBuild,
				job.TransitionBuild,
			),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(jobs)
	if err != nil {
		logger.Error("failed-to-encode-jobs", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
