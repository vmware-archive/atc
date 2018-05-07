package jobserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
)

func (s *Server) ListJobCombinationBuilds(pipeline db.Pipeline) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			builds []db.Build
			err    error
			until  int
			since  int
			limit  int
		)

		logger := s.logger.Session("list-job-combination-builds")

		jobName := r.FormValue(":job_name")
		teamName := r.FormValue(":team_name")

		urlUntil := r.FormValue(atc.PaginationQueryUntil)
		until, _ = strconv.Atoi(urlUntil)

		urlSince := r.FormValue(atc.PaginationQuerySince)
		since, _ = strconv.Atoi(urlSince)

		urlLimit := r.FormValue(atc.PaginationQueryLimit)
		limit, _ = strconv.Atoi(urlLimit)
		if limit == 0 {
			limit = atc.PaginationAPIDefaultLimit
		}

		job, found, err := pipeline.Job(jobName)
		if err != nil {
			logger.Error("failed-to-get-job", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobCombinationID := r.FormValue(":job_combination_id")
		jobCombination, err := job.FindJobCombinationByID(jobCombinationID)
		if err != nil {
			logger.Error("failed-to-get-job-combination", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		builds, pagination, err := jobCombination.Builds(db.Page{
			Since: since,
			Until: until,
			Limit: limit,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if pagination.Next != nil {
			s.addNextLink(w, teamName, pipeline.Name(), jobName, *pagination.Next)
		}

		if pagination.Previous != nil {
			s.addPreviousLink(w, teamName, pipeline.Name(), jobName, *pagination.Previous)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jobBuilds := make([]atc.Build, len(builds))
		for i := 0; i < len(builds); i++ {
			jobBuilds[i] = present.Build(builds[i])
		}

		err = json.NewEncoder(w).Encode(jobBuilds)
		if err != nil {
			logger.Error("failed-to-encode-job-builds", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
