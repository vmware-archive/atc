package accessor

import (
	"context"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/lock"
)

//go:generate counterfeiter . AccessorFactory
type AccessorFactory interface {
	CreateAccessor(context.Context) (Accessor, error)
}

type accessorFactory struct {
	conn        db.Conn
	lockFactory lock.LockFactory
}

func NewAccessorFactory(conn db.Conn, lockFactory lock.LockFactory) AccessorFactory {
	return &accessorFactory{
		conn:        conn,
		lockFactory: lockFactory,
	}
}

func (f *accessorFactory) CreateAccessor(ctx context.Context) (Accessor, error) {
	return nil, nil
}
