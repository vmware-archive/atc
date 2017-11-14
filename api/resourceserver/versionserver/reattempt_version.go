package versionserver

import (
	"net/http"
	"strconv"

	"github.com/concourse/atc/db"
	"github.com/tedsuo/rata"
)

func (s *Server) ReattemptResourceVersion(pipeline db.Pipeline) http.Handler {
	logger := s.logger.Session("reattempt-resource-version")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versionedResourceID, err := strconv.Atoi(rata.Param(r, "resource_version_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = pipeline.ReattemptVersionedResource(versionedResourceID)
		if err != nil {
			logger.Error("failed-to-reattempt-versioned-resource", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
