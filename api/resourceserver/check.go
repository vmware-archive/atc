package resourceserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/resource"
	"github.com/tedsuo/rata"
)

func (s *Server) CheckResource(dbPipeline dbng.Pipeline) http.Handler {
	logger := s.logger.Session("check-resource")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceName := rata.Param(r, "resource_name")

		var reqBody atc.CheckRequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			logger.Info("malformed-request", lager.Data{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fromVersion := reqBody.From
		if fromVersion == nil {
			latestVersion, found, err := dbPipeline.GetLatestVersionedResource(resourceName)
			if err != nil {
				logger.Info("failed-to-get-latest-versioned-resource", lager.Data{"error": err.Error()})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if found {
				fromVersion = atc.Version(latestVersion.Version)
			}
		}

		scanner := s.scannerFactory.NewResourceScanner(dbPipeline)

		err = scanner.ScanFromVersion(logger, resourceName, fromVersion)
		switch scanErr := err.(type) {
		case resource.ErrResourceScriptFailed:
			checkResponseBody := atc.CheckResponseBody{
				ExitStatus: scanErr.ExitStatus,
				Stderr:     scanErr.Stderr,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(checkResponseBody)
		case dbng.ResourceNotFoundError:
			w.WriteHeader(http.StatusNotFound)
		case error:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})
}
