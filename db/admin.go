package db

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/lock"
)

//go:generate counterfeiter . Admin
type Admin interface {
	CreateTeam(atc.Team) (Team, error)
}

type admin struct {
	conn        Conn
	lockFactory lock.LockFactory
}

func NewAdmin(conn Conn, lockFactory lock.LockFactory) Admin {
	return &admin{
		conn:        conn,
		lockFactory: lockFactory,
	}
}

func (a *admin) CreateTeam(team atc.Team) (Team, error) {
	return nil, nil
}
