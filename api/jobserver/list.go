package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) ListJobs(pipeline dbng.Pipeline) http.Handler {
	logger := s.logger.Session("list-jobs")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jobs []atc.Job

		dashboard, groups, err := pipeline.Dashboard()
		if err != nil {
			logger.Error("failed-to-get-dashboard", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		teamName := r.FormValue(":team_name")

		for _, job := range dashboard {
			jobs = append(
				jobs,
				present.Job(
					teamName,
					job.Job,
					groups,
					job.FinishedBuild,
					job.NextBuild,
				),
			)
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(jobs)
	})
}
