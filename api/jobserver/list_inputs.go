package jobserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) ListJobInputs(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-job-inputs")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	jobName := r.FormValue(":job_name")

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

	variables := s.variablesFactory.NewVariables(teamName, pipelineName)
	scheduler := s.schedulerFactory.BuildScheduler(pipeline, s.externalURL, variables)

	err = scheduler.SaveNextInputMapping(logger, job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buildInputs, found, err := job.GetNextBuildInputs()
	if err != nil {
		logger.Error("failed-to-get-next-build-inputs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resources, err := pipeline.Resources()
	if err != nil {
		logger.Error("failed-to-get-resources", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jobInputs := job.Config().Inputs()
	presentedBuildInputs := make([]atc.BuildInput, len(buildInputs))
	for i, input := range buildInputs {
		resource, _ := resources.Lookup(input.Resource)

		var config atc.JobInput
		for _, jobInput := range jobInputs {
			if jobInput.Name == input.Name {
				config = jobInput
				break
			}
		}

		presentedBuildInputs[i] = present.BuildInput(input, config, resource.Source())
	}

	err = json.NewEncoder(w).Encode(presentedBuildInputs)
	if err != nil {
		logger.Error("failed-to-encode-build-inputs", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
