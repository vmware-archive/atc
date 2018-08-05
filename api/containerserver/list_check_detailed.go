package containerserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/db"
)

func (s *Server) ListCheckContainersDetailed(team db.Team) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hLog := s.logger.Session("list-containers")

		w.Header().Set("Content-Type", "application/json")

		hLog.Debug("listing-check-containers-detailed")

		containers, err := team.FindCheckContainerDetails(true)
		if err != nil {
			hLog.Error("failed-to-find-containers", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hLog.Debug("listed", lager.Data{"container-count": len(containers)})

		presentedContainers := make([]atc.Container, len(containers))
		for i := 0; i < len(containers); i++ {
			container := containers[i]
			presentedContainers[i] = present.Container(container)
		}

		err = json.NewEncoder(w).Encode(presentedContainers)
		if err != nil {
			hLog.Error("failed-to-encode-containers", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
