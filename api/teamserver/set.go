package teamserver

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor"
	"github.com/concourse/atc/api/present"
	"github.com/concourse/skymarshal/provider"
)

func (s *Server) SetTeam(w http.ResponseWriter, r *http.Request) {
	hLog := s.logger.Session("set-team")

	hLog.Debug("setting-team")

	teamName := r.FormValue(":team_name")

	acc := s.accessorFactory.Create(r)

	// authTeam, authTeamFound := auth.GetTeam(r)
	// if !authTeamFound {
	// 	hLog.Error("failed-to-get-team-from-auth", errors.New("failed-to-get-team-from-auth"))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	var atcTeam atc.Team
	err := json.NewDecoder(r.Body).Decode(&atcTeam)
	if err != nil {
		hLog.Error("malformed-request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atcTeam.Name = teamName

	// team, found, err := acc.GetTeam(teamName)
	// if err != nil {
	// 	hLog.Error("failed-to-lookup-team", err, lager.Data{"teamName": teamName})
	// 	w.WriteHeader(accessor.HttpStatus(err))
	// 	return
	// }
	// if !authTeam.IsAdmin() && !authTeam.IsAuthorized(teamName) {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	providers := provider.GetProviders()

	for providerName, config := range atcTeam.Auth {
		p, found := providers[providerName]
		if !found {
			hLog.Error("failed-to-find-provider", err, lager.Data{"provider": providerName})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authConfig, err := p.UnmarshalConfig(config)
		if err != nil {
			hLog.Error("failed-to-unmarshal-auth", err, lager.Data{"provider": providerName})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = authConfig.Validate()
		if err != nil {
			hLog.Error("request-body-validation-error", err, lager.Data{"provider": providerName})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = authConfig.Finalize()
		if err != nil {
			hLog.Error("cannot-finalize-auth-config", err, lager.Data{"provider": providerName})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonConfig, err := p.MarshalConfig(authConfig)
		if err != nil {
			hLog.Error("cannot-marshal-auth-config", err, lager.Data{"provider": providerName})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		atcTeam.Auth[providerName] = jsonConfig
	}

	// team, found, err := s.teamFactory.FindTeam(teamName)
	// if err != nil {
	// 	hLog.Error("failed-to-lookup-team", err, lager.Data{"teamName": teamName})
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	team, created, err := acc.PutTeam(teamName, accessor.TeamConfig{atcTeam})
	if err != nil {
		hLog.Error("failed-to-update-team", err, lager.Data{"teamName": teamName})
		w.WriteHeader(accessor.HttpStatus(err))
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// if found {
	// 	hLog.Debug("updating-credentials")
	// 	err = team.UpdateProviderAuth(atcTeam.Auth)
	// 	if err != nil {
	// 		hLog.Error("failed-to-update-team", err, lager.Data{"teamName": teamName})
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// } else if authTeam.IsAdmin() {
	// 	hLog.Debug("creating team")

	// 	team, err = s.teamFactory.CreateTeam(atcTeam)
	// 	if err != nil {
	// 		hLog.Error("failed-to-save-team", err, lager.Data{"teamName": teamName})
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.WriteHeader(http.StatusCreated)
	// } else {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	err = json.NewEncoder(w).Encode(present.Team(team))
	if err != nil {
		hLog.Error("failed-to-encode-team", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
