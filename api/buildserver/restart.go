package buildserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
)

type rebuildJob struct {
	db.Job
	oldBuild db.Build
}

func (r rebuildJob) CreateBuild() (db.Build, error) {
	return r.CreateRebuild(r.oldBuild)
}

func (s *Server) RebuildBuild(build db.Build) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		logger := s.logger.Session("restart-build")

		pipeline, found, err := build.Pipeline()
		if err != nil {
			logger.Error("failed-resolve-pipeline", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		job, found, err := pipeline.Job(build.JobName())
		if err != nil {
			logger.Error("failed-resolve-job", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if job.Config().DisableManualTrigger {
			w.WriteHeader(http.StatusConflict)
			return
		}

		scheduler := s.schedulerFactory.BuildScheduler(pipeline, s.externalURL, s.variablesFactory.NewVariables(pipeline.TeamName(), pipeline.Name()))

		resourceTypes, err := pipeline.ResourceTypes()
		if err != nil {
			logger.Error("failed-to-get-resource-types", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		versionedResourceTypes := resourceTypes.Deserialize()

		resources, err := pipeline.Resources()
		if err != nil {
			logger.Error("failed-to-get-resources", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rebuild := rebuildJob{
			Job:      job,
			oldBuild: build,
		}

		dbBuild, _, err := scheduler.TriggerImmediately(logger, rebuild, resources, versionedResourceTypes)
		if err != nil {
			logger.Error("failed-to-trigger", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to trigger: %s", err)
			return
		}

		err = json.NewEncoder(w).Encode(present.Build(dbBuild))
		if err != nil {
			logger.Error("failed-to-encode-build", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
