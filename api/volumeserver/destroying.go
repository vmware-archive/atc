package volumeserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
)

func (s *Server) ListDestroyingVolumes(w http.ResponseWriter, r *http.Request) {
	workerName := r.URL.Query().Get("worker_name")
	w.Header().Set("Content-Type", "application/json")

	hLog := s.logger.Session("marked-volumes-for-worker", lager.Data{"worker_name": workerName})

	if workerName != "" {
		destroyingVolumesHandles, err := s.destroyer.FindOrphanedVolumesasDestroying(workerName)
		if err != nil {
			hLog.Error("failed-to-find-destroying-volumes", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hLog.Debug("list", lager.Data{"destroying-volume-count": len(destroyingVolumesHandles)})

		err = json.NewEncoder(w).Encode(destroyingVolumesHandles)
		if err != nil {
			hLog.Error("failed-to-encode-volumes", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		hLog.Info("failed-to-find-worker")
		w.WriteHeader(http.StatusNotFound)
	}
}
