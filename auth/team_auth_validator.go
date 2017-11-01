package auth

import (
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/auth/provider"
	"github.com/concourse/atc/db"
)

type teamAuthValidator struct {
	teamFactory  db.TeamFactory
	jwtValidator Validator
}

func NewTeamAuthValidator(
	teamFactory db.TeamFactory,
	jwtValidator Validator,
) Validator {
	return &teamAuthValidator{
		teamFactory:  teamFactory,
		jwtValidator: jwtValidator,
	}
}

func (v teamAuthValidator) IsAuthenticated(logger lager.Logger, r *http.Request) bool {
	teamName := r.FormValue(":team_name")
	team, found, err := v.teamFactory.FindTeam(teamName)
	if err != nil || !found {
		return false
	}

	if !isAuthConfigured(team) {
		return true
	}

	if team.BasicAuth != nil && NewBasicAuthValidator(team).IsAuthenticated(logger, r) {
		return true
	}
	// try ldap if configured
	if team.LdapBasicAuth != nil && NewLdapBasicAuthValidator(team).IsAuthenticated(logger, r) {
		return true
	}

	return v.jwtValidator.IsAuthenticated(logger, r)

}

func isAuthConfigured(t db.Team) bool {
	if t.BasicAuth() != nil {
		return true
	}

	for name := range provider.GetProviders() {
		_, configured := t.Auth()[name]
		if configured {
			return true
		}
	}

	return false
}
