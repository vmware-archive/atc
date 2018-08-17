package teamserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
)

func (s *Server) ListTeams(w http.ResponseWriter, r *http.Request) {
	hLog := s.logger.Session("list-teams")

	teams, err := s.teamFactory.GetTeams()
	if err != nil {
		hLog.Error("failed-to-get-teams", errors.New("sorry"))
		w.WriteHeader(http.StatusInternalServerError)
	}

	acc := accessor.GetAccessor(r)
	presentedTeams := make([]atc.Team, 0)
	for _, team := range teams {
		if acc.IsAdmin() || acc.IsAuthorized(team.Name()) {
			presentedTeams = append(presentedTeams, present.Team(team))
		}
	}

	sortTeams := make([]interface{}, 0)
	for _, team := range presentedTeams {
		sortTeams = append(sortTeams, team)
	}

	alphabeticalSort := func(team1, team2 *interface{}) bool {
		t1 := (*team1).(atc.Team)
		t2 := (*team2).(atc.Team)
		return t1.Name < t2.Name
	}

	sortedTeams := Sorter{items: sortTeams}.GenericSort(alphabeticalSort)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sortedTeams)
	if err != nil {
		hLog.Error("failed-to-encode-teams", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
