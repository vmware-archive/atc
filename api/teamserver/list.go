package teamserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) ListTeams(w http.ResponseWriter, r *http.Request) {
	hLog := s.logger.Session("list-teams")

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		hLog.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
	}

	teams, err := acc.Teams(accessor.Read)
	if err != nil {
		hLog.Error("failed-to-get-teams", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	presentedTeams := make([]atc.Team, len(teams))
	for i, team := range teams {
		presentedTeams[i] = present.Team(team)
	}

	err = json.NewEncoder(w).Encode(presentedTeams)
	if err != nil {
		hLog.Error("failed-to-encode-teams", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
