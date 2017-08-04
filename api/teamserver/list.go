package teamserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/db"
)

func (s *Server) ListTeams(w http.ResponseWriter, r *http.Request) {
	hLog := s.logger.Session("list-teams")

	teams, err := s.teamFactory.GetTeams()
	if err != nil {
		hLog.Error("failed-to-get-teams", errors.New("sorry"))
		w.WriteHeader(http.StatusInternalServerError)
	}

	var presenter func(db.Team) atc.Team
	authTeam, authTeamFound := auth.GetTeam(r)
	if authTeamFound && authTeam.IsAdmin() {
		presenter = present.TeamWithAdmin
	} else {
		presenter = present.Team
	}

	presentedTeams := make([]atc.Team, len(teams))
	for i, team := range teams {
		presentedTeams[i] = presenter(team)
	}

	json.NewEncoder(w).Encode(presentedTeams)
}
