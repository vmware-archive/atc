package voyager

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db/encryption"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/db/migration/migrations"
	multierror "github.com/hashicorp/go-multierror"
	_ "github.com/lib/pq"
)

type Migrator interface {
	CurrentVersion() (int, error)
	SupportedVersion() (int, error)
	Migrate(version int) error
	Up() error
	Migrations() ([]migration, error)
}

func NewMigrator(db *sql.DB, lockFactory lock.LockFactory, strategy encryption.Strategy, source Source) Migrator {
	return NewMigratorForMigrations(db, lockFactory, strategy, source)
}

func NewMigratorForMigrations(db *sql.DB, lockFactory lock.LockFactory, strategy encryption.Strategy, source Source) Migrator {
	return &migrator{
		db,
		lockFactory,
		strategy,
		lager.NewLogger("migrations"),
		source,
	}
}

type migrator struct {
	db          *sql.DB
	lockFactory lock.LockFactory
	strategy    encryption.Strategy
	logger      lager.Logger
	source      Source
}

func (m *migrator) SupportedVersion() (int, error) {
	matches := []migration{}

	assets := m.source.AssetNames()

	var parser = NewParser(m.source)
	for _, match := range assets {
		if migration, err := parser.ParseMigrationFilename(match); err == nil {
			matches = append(matches, migration)
		}
	}
	sortMigrations(matches)
	return matches[len(matches)-1].Version, nil
}

func (self *migrator) CurrentVersion() (int, error) {
	var currentVersion int
	var direction string
	err := self.db.QueryRow("SELECT version, direction FROM migrations_history WHERE status!='failed' ORDER BY tstamp DESC LIMIT 1").Scan(&currentVersion, &direction)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return -1, err
	}
	migrations, err := self.Migrations()
	if err != nil {
		return -1, err
	}
	versions := []int{migrations[0].Version}
	for _, m := range migrations {
		if m.Version > versions[len(versions)-1] {
			versions = append(versions, m.Version)
		}
	}
	for i, version := range versions {
		if currentVersion == version && direction == "down" {
			currentVersion = versions[i-1]
			break
		}
	}
	return currentVersion, nil
}

func (self *migrator) Migrate(toVersion int) error {

	lock, err := self.acquireLock()
	if err != nil {
		return err
	}

	if lock != nil {
		defer lock.Release()
	}

	existingDBVersion, err := self.migrateFromSchemaMigrations()
	if err != nil {
		return err
	}

	_, err = self.db.Exec("CREATE TABLE IF NOT EXISTS migrations_history (version bigint, tstamp timestamp with time zone, direction varchar, status varchar, dirty boolean)")
	if err != nil {
		return err
	}

	if existingDBVersion > 0 {
		var containsOldMigrationInfo bool
		err = self.db.QueryRow("SELECT EXISTS (SELECT 1 FROM migrations_history where version=$1)", existingDBVersion).Scan(&containsOldMigrationInfo)

		if !containsOldMigrationInfo {
			_, err = self.db.Exec("INSERT INTO migrations_history (version, tstamp, direction, status, dirty) VALUES ($1, current_timestamp, 'up', 'passed', false)", existingDBVersion)
			if err != nil {
				return err
			}
		}
	}

	currentVersion, err := self.CurrentVersion()
	if err != nil {
		return err
	}

	migrations, err := self.Migrations()
	if err != nil {
		return err
	}

	if currentVersion <= toVersion {
		for _, m := range migrations {
			if currentVersion < m.Version && m.Version <= toVersion && m.Direction == "up" {
				err = self.runMigration(m)
				if err != nil {
					return err
				}
			}
		}
	} else {
		for i := len(migrations) - 1; i >= 0; i-- {
			if currentVersion >= migrations[i].Version && migrations[i].Version > toVersion && migrations[i].Direction == "down" {
				err = self.runMigration(migrations[i])
				if err != nil {
					return err
				}

			}
		}

		err = self.migrateToSchemaMigrations(toVersion)
		if err != nil {
			return err
		}
	}
	return nil
}

type Strategy int

const (
	GoMigration Strategy = iota
	SQLTransaction
	SQLNoTransaction
)

type migration struct {
	Name       string
	Version    int
	Direction  string
	Statements []string
	Strategy   Strategy
}

func (m *migrator) recordMigrationFailure(migration migration, err error, dirty bool) error {
	_, dbErr := m.db.Exec("INSERT INTO migrations_history (version, tstamp, direction, status, dirty) VALUES ($1, current_timestamp, $2, 'failed', $3)", migration.Version, migration.Direction, dirty)
	return multierror.Append(fmt.Errorf("Migration '%s' failed: %v", migration.Name, err), dbErr)
}

func (m *migrator) runMigration(migration migration) error {
	var err error

	switch migration.Strategy {
	case GoMigration:
		err = migrations.NewMigrations(m.db, m.strategy).Run(migration.Name)
		if err != nil {
			return m.recordMigrationFailure(migration, err, false)
		}
	case SQLTransaction:
		tx, err := m.db.Begin()
		for _, statement := range migration.Statements {
			_, err = tx.Exec(statement)
			if err != nil {
				tx.Rollback()
				err = multierror.Append(fmt.Errorf("Transaction %v failed, rolled back the migration", statement), err)
				if err != nil {
					return m.recordMigrationFailure(migration, err, false)
				}
			}
		}
		err = tx.Commit()
	case SQLNoTransaction:
		_, err = m.db.Exec(migration.Statements[0])
		if err != nil {
			return m.recordMigrationFailure(migration, err, true)
		}
	}

	_, err = m.db.Exec("INSERT INTO migrations_history (version, tstamp, direction, status, dirty) VALUES ($1, current_timestamp, $2, 'passed', false)", migration.Version, migration.Direction)
	return err
}

func (self *migrator) Migrations() ([]migration, error) {
	migrationList := []migration{}
	assets := self.source.AssetNames()
	var parser = NewParser(self.source)
	for _, assetName := range assets {
		parsedMigration, err := parser.ParseFileToMigration(assetName)
		if err != nil {
			return nil, err
		}
		migrationList = append(migrationList, parsedMigration)
	}

	sortMigrations(migrationList)

	return migrationList, nil
}

func (self *migrator) Up() error {
	migrations, err := self.Migrations()
	if err != nil {
		return err
	}
	return self.Migrate(migrations[len(migrations)-1].Version)
}

func (self *migrator) acquireLock() (lock.Lock, error) {

	var err error
	var acquired bool
	var newLock lock.Lock

	if self.lockFactory != nil {
		for {
			newLock, acquired, err = self.lockFactory.Acquire(self.logger, lock.NewDatabaseMigrationLockID())

			if err != nil {
				return nil, err
			}

			if acquired {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}

	return newLock, err
}

func CheckTableExist(db *sql.DB, tableName string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name=$1)", tableName).Scan(&exists)
	return err != nil || exists
}

func (self *migrator) migrateFromSchemaMigrations() (int, error) {
	if !CheckTableExist(self.db, "schema_migrations") || CheckTableExist(self.db, "migrations_history") {
		return 0, nil
	}

	var isDirty = false
	var existingVersion int
	err := self.db.QueryRow("SELECT dirty, version FROM schema_migrations LIMIT 1").Scan(&isDirty, &existingVersion)
	if err != nil {
		return 0, err
	}

	if isDirty {
		return 0, errors.New("cannot begin migration. Database is in a dirty state")
	}

	return existingVersion, nil
}

type filenames []string

func sortMigrations(migrationList []migration) {
	sort.Slice(migrationList, func(i, j int) bool {
		return migrationList[i].Version < migrationList[j].Version
	})
}

func (self *migrator) migrateToSchemaMigrations(toVersion int) error {
	newMigrationsHistoryFirstVersion := 1532706545

	if toVersion >= newMigrationsHistoryFirstVersion {
		return nil
	}

	if !CheckTableExist(self.db, "schema_migrations") {
		_, err := self.db.Exec("CREATE TABLE schema_migrations (version bigint, dirty boolean)")
		if err != nil {
			return err
		}

		_, err = self.db.Exec("INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)", toVersion)
		if err != nil {
			return err
		}
	} else {
		_, err := self.db.Exec("UPDATE schema_migrations SET version=$1, dirty=false", toVersion)
		if err != nil {
			return err
		}
	}

	return nil
}
