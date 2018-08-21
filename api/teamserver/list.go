package teamserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

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
	authorizedTeams := make([]atc.Team, 0)
	publicTeams := make([]atc.Team, 0)
	// presentedTeams := make([]atc.Team, 0)

	emptyTeam := atc.Team{}
	for _, team := range teams {
		if acc.IsAdmin() || acc.IsAuthorized(team.Name()) {
			if reflect.DeepEqual(team.Auth, emptyTeam.Auth) {
				publicTeams = append(publicTeams, present.Team(team))
			} else {
				authorizedTeams = append(authorizedTeams, present.Team(team))
			}
		}
	}

	sortAuthorizedTeams := make([]interface{}, 0)
	for _, team := range authorizedTeams {
		sortAuthorizedTeams = append(sortAuthorizedTeams, team)
	}

	sortPublicTeams := make([]interface{}, 0)
	for _, team := range publicTeams {
		sortPublicTeams = append(sortPublicTeams, team)
	}

	alphabeticalSort := func(team1, team2 *interface{}) bool {
		t1 := (*team1).(atc.Team)
		t2 := (*team2).(atc.Team)
		return t1.Name < t2.Name
	}

	sortAuthorizedTeams = Sorter{items: sortAuthorizedTeams}.GenericSort(alphabeticalSort)
	sortPublicTeams = Sorter{items: sortPublicTeams}.GenericSort(alphabeticalSort)

	// fmt.Printf("authrized & public %v", append(sortAuthorizedTeams, sortPublicTeams))
	// sortAuthorizedTeams = append(sortAuthorizedTeams, sortPublicTeams)
	// sortedTeams := Sorter{items: presentedTeams}.GenericSort(alphabeticalSort)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(append(sortAuthorizedTeams, sortPublicTeams))
	if err != nil {
		hLog.Error("failed-to-encode-teams", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
