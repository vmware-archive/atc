package pipelineserver

import (
	"fmt"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/jobserver"
	"github.com/concourse/atc/db"
)

func (s *Server) PipelineBadge(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("pipeline-badge")
	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	jobs, err := acc.TeamPipelineJobs(accessor.Write, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-jobs", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	var build db.Build

	jobStatusPrecedence := map[db.BuildStatus]int{
		db.BuildStatusFailed:    1,
		db.BuildStatusErrored:   2,
		db.BuildStatusAborted:   3,
		db.BuildStatusSucceeded: 4,
	}

	for _, job := range jobs {
		b, _, err := job.FinishedAndNextBuild()
		if err != nil {
			logger.Error("could-not-get-finished-and-next-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if b == nil {
			continue
		}

		if build == nil || jobStatusPrecedence[b.Status()] < jobStatusPrecedence[build.Status()] {
			build = b
		}
	}

	w.Header().Set("Content-type", "image/svg+xml")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, jobserver.BadgeForBuild(build))
}
