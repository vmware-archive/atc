package teamserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/skymarshal/provider"
)

func (s *Server) SetTeam(w http.ResponseWriter, r *http.Request) {
	hLog := s.logger.Session("set-team")

	hLog.Debug("setting-team")

	teamName := r.FormValue(":team_name")

	atcTeam, err := s.decodeTeam(r)
	if err != nil {
		hLog.Error("failed-to-decode-team", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := s.accessorFactory.CreateAccessor(r.Context())
	if err != nil {
		hLog.Error("failed-to-get-user", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	team, err := acc.Team(accessor.Write, teamName)

	if err == nil {
		err = team.UpdateProviderAuth(atcTeam.Auth)
		if err != nil {
			hLog.Error("failed-to-update-provider-auth-for-team", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else if err == accessor.ErrNotFound {
		admin, err := acc.Admin(accessor.Write)
		if err != nil {
			hLog.Error("failed-to-get-admin-user", err)
			w.WriteHeader(accessor.HttpStatus(err))
			return
		}

		team, err = admin.CreateTeam(atcTeam)
		if err != nil {
			hLog.Error("failed-to-create-team-for-admin", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	} else {
		hLog.Error("failed-to-get-team", err)
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(present.Team(team))
	if err != nil {
		hLog.Error("failed-to-encode-team", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) decodeTeam(r *http.Request) (atc.Team, error) {

	defer r.Body.Close()

	var team atc.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		return team, err
	}

	team.Name = r.FormValue(":team_name")

	providers := provider.GetProviders()

	for providerName, config := range team.Auth {
		p, found := providers[providerName]
		if !found {
			return team, errors.New("provider-not-found")
		}

		authConfig, err := p.UnmarshalConfig(config)
		if err != nil {
			return team, err
		}

		err = authConfig.Validate()
		if err != nil {
			return team, err
		}

		err = authConfig.Finalize()
		if err != nil {
			return team, err
		}

		jsonConfig, err := p.MarshalConfig(authConfig)
		if err != nil {
			return team, err
		}

		team.Auth[providerName] = jsonConfig
	}

	return team, nil
}
