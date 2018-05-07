package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
)

func (s *Server) ListSpaceJobs(pipeline db.Pipeline) http.Handler {
	logger := s.logger.Session("list-space-jobs")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dashboard, err := s.spaceJobFactory.PipelineJobs(pipeline.TeamName(), pipeline.Name())
		if err != nil {
			logger.Error("failed-to-get-pipeline-jobs", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var spaceJobs []atc.SpaceJob

		for _, spaceJob := range dashboard {
			spaceJobs = append(spaceJobs, present.SpaceJob(spaceJob))
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(spaceJobs)
		if err != nil {
			logger.Error("failed-to-encode-jobs", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
