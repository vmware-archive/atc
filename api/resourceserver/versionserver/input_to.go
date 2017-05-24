package versionserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/dbng"
)

func (s *Server) ListBuildsWithVersionAsInput(pipeline dbng.Pipeline) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versionIDString := r.FormValue(":resource_version_id")
		versionID, _ := strconv.Atoi(versionIDString)

		builds, err := pipeline.GetBuildsWithVersionAsInput(versionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		presentedBuilds := []atc.Build{}
		for _, build := range builds {
			presentedBuilds = append(presentedBuilds, present.Build(build))
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(presentedBuilds)
	})
}
