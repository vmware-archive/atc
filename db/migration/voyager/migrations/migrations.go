package migrations

import (
	"database/sql"
	"reflect"

	"github.com/concourse/atc/db/encryption"
)

func NewGoMigrationsRunner(db *sql.DB, es encryption.Strategy) *GoMigrationsRunner {
	return &GoMigrationsRunner{db, es}
}

type GoMigrationsRunner struct {
	*sql.DB
	encryption.Strategy
}

func (self *GoMigrationsRunner) Run(name string) error {

	res := reflect.ValueOf(self).MethodByName(name).Call(nil)

	ret := res[0].Interface()

	if ret != nil {
		return ret.(error)
	}

	return nil
}
