package jobserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) CreateJobBuild(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logger := s.logger.Session("create-job-build")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	jobName := r.FormValue(":job_name")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipeline, err := acc.TeamPipeline(accessor.Write, teamName, pipelineName)
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

	build, _, err := scheduler.TriggerImmediately(logger, job, resources, versionedResourceTypes)
	if err != nil {
		logger.Error("failed-to-trigger", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to trigger: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode(present.Build(build))
	if err != nil {
		logger.Error("failed-to-encode-build", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
