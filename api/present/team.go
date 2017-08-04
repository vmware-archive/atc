package present

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
)

func Team(team db.Team) atc.Team {
	return atc.Team{
		ID:   team.ID(),
		Name: team.Name(),
	}
}

func TeamWithAdmin(team db.Team) atc.Team {
	return atc.Team{
		ID:        team.ID(),
		Name:      team.Name(),
		BasicAuth: team.BasicAuth(),
		Auth:      team.Auth(),
	}
}
