package accessor

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

//go:generate counterfeiter . Accessor
type Accessor interface {
	GetTeam(string) (db.Team, bool, error)
	PutTeam(string, TeamConfig) (db.Team, bool, error)
	PutTeamPipeline(string, string, PipelineConfig) (db.Pipeline, bool, error)
	//GetTeamPipelineConfig(string, string) (db.Pipeline, error)
	//	GetPipelines()
}

type TeamConfig struct {
	Team atc.Team
}

type PipelineConfig struct {
	Config  atc.Config
	Version db.ConfigVersion
	State   db.PipelinePausedState
}

type accessor struct {
	teamFactory db.TeamFactory

	teamNames []string
	isAdmin   bool
	logger    lager.Logger
}

func (a *accessor) PutTeam(teamName string, t TeamConfig) (db.Team, bool, error) {
	team, found, err := a.teamFactory.FindTeam(teamName)
	if err != nil {
		a.logger.Error("error-when-get-team", err)
		return nil, false, err
	}
	if !found {
		if a.isAdmin {
			team, err = a.teamFactory.CreateTeam(t.Team)
			if err != nil {
				a.logger.Error("error-when-creating-team-as-admin", err)
				return nil, false, err
			} else {
				a.logger.Debug("created-team", lager.Data{"name": teamName})
				return team, true, nil
			}
		} else {
			return nil, false, ErrForbidden
		}
	} else {
		if a.isTeamMember(teamName) || a.isAdmin {
			err = team.UpdateProviderAuth(t.Team.Auth)
			if err != nil {
				a.logger.Error("error-when-updating-provider-auth", err)
				return nil, false, err
			} else {
				a.logger.Debug("updated-team-auth", lager.Data{"name": teamName})
				return team, false, nil
			}
		} else {
			return nil, false, ErrForbidden
		}
	}
	// team, found, err := a.GetTeam(teamName)
	// a.logger.Debug("get-team", lager.Data{"name": teamName})
	// if err != nil {
	// 	a.logger.Error("error-when-get-team", err)
	// 	return nil, false, err
	// }

	// if found {
	// 	err = team.UpdateProviderAuth(t.Team.Auth)
	// 	if err != nil {
	// 		a.logger.Error("error-when-updating-provider-auth", err)
	// 		return nil, false, err
	// 	} else {
	// 		a.logger.Debug("updated-team-auth", lager.Data{"name": teamName})
	// 		return team, false, nil
	// 	}

	// } else if a.isAdmin {
	// 	team, err = a.teamFactory.CreateTeam(t.Team)
	// 	if err != nil {
	// 		a.logger.Error("error-when-creating-team", err)
	// 		return nil, false, err
	// 	} else {
	// 		a.logger.Debug("created-team", lager.Data{"name": teamName})
	// 		return team, true, nil
	// 	}

	// } else {
	// 	a.logger.Error("error-forbidden-team", err)
	// 	return nil, false, ErrForbidden
	// }
}

func (a *accessor) PutTeamPipeline(teamName string, pipelineName string, p PipelineConfig) (db.Pipeline, bool, error) {

	team, found, err := a.GetTeam(teamName)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, ErrNotFound
	}

	return team.SavePipeline(pipelineName, p.Config, p.Version, p.State)
}

func (a *accessor) GetTeam(teamName string) (db.Team, bool, error) {
	if !a.isTeamMember(teamName) && !a.isAdmin {
		return nil, false, ErrNotAuthorized
	}
	return a.teamFactory.FindTeam(teamName)
}

func (a *accessor) isTeamMember(teamName string) bool {
	for _, team := range a.teamNames {
		// if a.isAdmin {
		// 	return true
		// }
		if team == teamName {
			return true
		}
	}
	return false
}

// no token accessor --> no_teams --> resource (pub/priavte) --> resource
// invalid token --> no_teams -- resource (pub/pri) --> resource
// valid token no team access -- no_teams ()
// valid token w team access --> teams --> resource
// admin token --> isadmin true -> resource
// system token --> issystem true -> resource
