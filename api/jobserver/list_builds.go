package jobserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) ListJobBuilds(pipeline dbng.Pipeline) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			builds []dbng.Build
			err    error
			until  int
			since  int
			limit  int
		)

		logger := s.logger.Session("list-job-builds")

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

		builds, pagination, err := job.Builds(dbng.Page{
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

		w.WriteHeader(http.StatusOK)

		jobBuilds := make([]atc.Build, len(builds))
		for i := 0; i < len(builds); i++ {
			jobBuilds[i] = present.Build(builds[i])
		}
		json.NewEncoder(w).Encode(jobBuilds)
	})
}

func (s *Server) addNextLink(w http.ResponseWriter, teamName, pipelineName, jobName string, page dbng.Page) {
	w.Header().Add("Link", fmt.Sprintf(
		`<%s/api/v1/teams/%s/pipelines/%s/jobs/%s/builds?%s=%d&%s=%d>; rel="%s"`,
		s.externalURL,
		teamName,
		pipelineName,
		jobName,
		atc.PaginationQuerySince,
		page.Since,
		atc.PaginationQueryLimit,
		page.Limit,
		atc.LinkRelNext,
	))
}

func (s *Server) addPreviousLink(w http.ResponseWriter, teamName, pipelineName, jobName string, page dbng.Page) {
	w.Header().Add("Link", fmt.Sprintf(
		`<%s/api/v1/teams/%s/pipelines/%s/jobs/%s/builds?%s=%d&%s=%d>; rel="%s"`,
		s.externalURL,
		teamName,
		pipelineName,
		jobName,
		atc.PaginationQueryUntil,
		page.Until,
		atc.PaginationQueryLimit,
		page.Limit,
		atc.LinkRelPrevious,
	))
}
