package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) GetJobBuild(pipeline dbng.Pipeline) http.Handler {
	logger := s.logger.Session("get-job-build")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobName := r.FormValue(":job_name")
		buildName := r.FormValue(":build_name")

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

		build, found, err := job.Build(buildName)
		if err != nil {
			logger.Error("failed-to-get-job-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(present.Build(build))
	})
}
