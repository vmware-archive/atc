package resourceserver

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/db"
)

// CheckResourceWebHook defines a handler for process a check resource request via an access token.
func (s *Server) CheckResourceWebHook(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("check-resource-webhook")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")
	resourceName := r.FormValue(":resource_name")
	webhookToken := r.URL.Query().Get("webhook_token")

	if webhookToken == "" {
		logger.Info("no-webhook-token", lager.Data{"error": "missing webhook_token"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		logger.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	pipeline, err := acc.TeamPipeline(accessor.Skip, teamName, pipelineName)
	if err != nil {
		logger.Error("failed-to-get-pipeline", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	resource, found, err := pipeline.Resource(resourceName)
	if err != nil {
		logger.Error("failed-to-get-resource", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !found {
		logger.Info("resource-not-found", lager.Data{"error": fmt.Sprintf("Resource not found %s", resourceName)})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token := resource.WebhookToken()
	if token != webhookToken {
		logger.Info("invalid-token", lager.Data{"error": fmt.Sprintf("invalid token for webhook %s", webhookToken)})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var fromVersion atc.Version
	latestVersion, found, err := pipeline.GetLatestVersionedResource(resourceName)
	if err != nil {
		logger.Error("failed-to-get-latest-versioned-resource", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if found {
		fromVersion = atc.Version(latestVersion.Version)
	}

	scanner := s.scannerFactory.NewResourceScanner(pipeline)
	err = scanner.ScanFromVersion(logger, resourceName, fromVersion)
	switch err.(type) {
	case db.ResourceNotFoundError:
		w.WriteHeader(http.StatusNotFound)
	case error:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
