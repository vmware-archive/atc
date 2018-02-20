package resourceserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
)

func (s *Server) CheckResource(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("check-resource")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	resourceName := r.FormValue(":resource_name")

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

	var reqBody atc.CheckRequestBody
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		logger.Info("malformed-request", lager.Data{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromVersion := reqBody.From
	if fromVersion == nil {
		latestVersion, found, err := pipeline.GetLatestVersionedResource(resourceName)
		if err != nil {
			logger.Info("failed-to-get-latest-versioned-resource", lager.Data{"error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if found {
			fromVersion = atc.Version(latestVersion.Version)
		}
	}

	scanner := s.scannerFactory.NewResourceScanner(pipeline)
	err = scanner.ScanFromVersion(logger, resourceName, fromVersion)
	switch scanErr := err.(type) {
	case resource.ErrResourceScriptFailed:

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(atc.CheckResponseBody{
			ExitStatus: scanErr.ExitStatus,
			Stderr:     scanErr.Stderr,
		})
		if err != nil {
			logger.Error("failed-to-encode-check-response-body", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case db.ResourceNotFoundError:
		w.WriteHeader(http.StatusNotFound)
	case error:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
