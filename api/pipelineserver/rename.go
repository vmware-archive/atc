package pipelineserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
)

func (s *Server) RenamePipeline(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("rename-pipeline")

	teamName := r.FormValue(":team_name")
	pipelineName := r.FormValue(":pipeline_name")

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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rename atc.RenameRequest
	err = json.Unmarshal(data, &rename)
	if err != nil {
		logger.Error("failed-to-unmarshal-body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = pipeline.Rename(rename.NewName)
	if err != nil {
		logger.Error("failed-to-update-name", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
