package jobserver

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
	"github.com/google/jsonapi"
)

func (s *Server) GetJobBuild(pipeline db.Pipeline) http.Handler {
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
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Job Not Found Error",
				Detail: fmt.Sprintf("Job with name '%s' not found.", jobName),
				Status: "404",
			}})
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
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Build Not Found Error",
				Detail: fmt.Sprintf("Build with name '%s' not found.", buildName),
				Status: "404",
			}})
			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(present.Build(build))
	})
}
