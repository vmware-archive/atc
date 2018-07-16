package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db/encryption"
	"github.com/concourse/atc/db/lock"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source"
	"github.com/mattes/migrate/source/go-bindata"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func NewOpenHelper(driver, name string, lockFactory lock.LockFactory, strategy encryption.Strategy) *OpenHelper {
	return &OpenHelper{
		driver,
		name,
		lockFactory,
		strategy,
	}
}

type OpenHelper struct {
	driver         string
	dataSourceName string
	lockFactory    lock.LockFactory
	strategy       encryption.Strategy
}

func (self *OpenHelper) CurrentVersion() (int, error) {
	db, err := sql.Open(self.driver, self.dataSourceName)
	if err != nil {
		return -1, err
	}

	defer db.Close()

	return NewMigrator(db, self.lockFactory, self.strategy).CurrentVersion()
}

func (self *OpenHelper) SupportedVersion() (int, error) {
	db, err := sql.Open(self.driver, self.dataSourceName)
	if err != nil {
		return -1, err
	}

	defer db.Close()

	return NewMigrator(db, self.lockFactory, self.strategy).SupportedVersion()
}

func (self *OpenHelper) Open() (*sql.DB, error) {
	db, err := sql.Open(self.driver, self.dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := NewMigrator(db, self.lockFactory, self.strategy).Up(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func (self *OpenHelper) OpenAtVersion(version int) (*sql.DB, error) {
	db, err := sql.Open(self.driver, self.dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := NewMigrator(db, self.lockFactory, self.strategy).Migrate(version); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func (self *OpenHelper) MigrateToVersion(version int) error {
	db, err := sql.Open(self.driver, self.dataSourceName)
	if err != nil {
		return err
	}

	defer db.Close()

	if err := NewMigrator(db, self.lockFactory, self.strategy).Migrate(version); err != nil {
		return err
	}

	return nil
}

type Migrator interface {
	CurrentVersion() (int, error)
	SupportedVersion() (int, error)
	Migrate(version int) error
	Up() error
}

var extractOnce = &sync.Once{}

func NewMigrator(db *sql.DB, lockFactory lock.LockFactory, strategy encryption.Strategy) Migrator {
	return NewMigratorForMigrations(db, lockFactory, strategy, AssetNames())
}

func NewMigratorForMigrations(db *sql.DB, lockFactory lock.LockFactory, strategy encryption.Strategy, migrations []string) Migrator {

	var tempdir string
	extractOnce.Do(func() {
		tempdir, _ = ioutil.TempDir("", "migrations")
		//defer os.RemoveAll(tempdir) // clean up
		fmt.Println("=======================================", tempdir)
		var content []byte
		for _, filename := range AssetNames() {
			content, _ = Asset(filename)
			tmpfn := filepath.Join(tempdir, filename)
			_ = ioutil.WriteFile(tmpfn, content, 0644)
		}
	})

	return &migrator{
		db,
		lockFactory,
		strategy,
		lager.NewLogger("migrations"),
		migrations,
		tempdir,
	}
}

type migrator struct {
	db            *sql.DB
	lockFactory   lock.LockFactory
	strategy      encryption.Strategy
	logger        lager.Logger
	migrations    filenames
	migrationsdir string
}

func (self *migrator) Stategy() encryption.Strategy {
	return self.Stategy()
}

func (self *migrator) SupportedVersion() (int, error) {

	v, err := goose.GetDBVersion(self.db)
	if err != nil {
		return -1, err
	}
	return int(v), nil

	// latest := self.migrations.Latest()
	//
	// m, err := source.Parse(latest)
	// if err != nil {
	// 	return -1, err
	// }
	//
	//return int(m.Version), nil
}

func (self *migrator) CurrentVersion() (int, error) {
	_, lock, err := self.openWithLock()
	if err != nil {
		return -1, err
	}

	if lock != nil {
		defer lock.Release()
	}

	v, err := goose.GetDBVersion(self.db)
	if err != nil {
		return -1, err
	}
	return int(v), nil

	// version, _, err := m.Version()
	// if err != nil {
	// 	return -1, err
	// }
	//
	// return int(version), nil
}

func (self *migrator) Migrate(version int) error {

	_, lock, err := self.openWithLock()
	if err != nil {
		return err
	}

	if lock != nil {
		defer lock.Release()
	}

	if err = goose.UpTo(self.db, self.migrationsdir, int64(version)); err != nil {
		if err.Error() != "no change" {
			return err
		}
	}
	// if err = m.Migrate(uint(version)); err != nil {
	// 	if err.Error() != "no change" {
	// 		return err
	// 	}
	// }

	return nil
}

func (self *migrator) Up() error {

	_, lock, err := self.openWithLock()
	if err != nil {
		return err
	}

	if lock != nil {
		defer lock.Release()
	}

	if err = goose.Up(self.db, self.migrationsdir); err != nil {
		if err.Error() != "no change" {
			return err
		}
	}
	// if err = m.Up(); err != nil {
	// 	if err.Error() != "no change" {
	// 		return err
	// 	}
	// }

	return nil
}

func (self *migrator) open() (*migrate.Migrate, error) {

	//func (self *migrator) open() (*migrate.Migrate, error) {

	forceVersion, err := self.checkLegacyVersion()
	if err != nil {
		return nil, err
	}

	s, err := bindata.WithInstance(bindata.Resource(
		self.migrations,
		func(name string) ([]byte, error) {
			return Asset(name)
		}),
	)

	d, err := postgres.WithInstance(self.db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	driver := NewDriver(d, self.db, self.strategy)

	m, err := migrate.NewWithInstance("go-bindata", s, "postgres", driver)
	if err != nil {
		return nil, err
	}

	if forceVersion > 0 {
		if err = m.Force(forceVersion); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (self *migrator) openWithLock() (*migrate.Migrate, lock.Lock, error) {

	var err error
	var acquired bool
	var newLock lock.Lock

	if self.lockFactory != nil {
		for {
			newLock, acquired, err = self.lockFactory.Acquire(self.logger, lock.NewDatabaseMigrationLockID())

			if err != nil {
				return nil, nil, err
			}

			if acquired {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}

	m, err := self.open()

	if err != nil && newLock != nil {
		newLock.Release()
		return nil, nil, err
	}

	return m, newLock, err
}

func (self *migrator) existLegacyVersion() bool {
	var exists bool
	err := self.db.QueryRow("SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = 'migration_version')").Scan(&exists)
	return err != nil || exists
}

func (self *migrator) checkLegacyVersion() (int, error) {
	oldMigrationLastVersion := 189
	newMigrationStartVersion := 1510262030

	var err error
	var dbVersion int

	exists := self.existLegacyVersion()
	if !exists {
		return -1, nil
	}

	if err = self.db.QueryRow("SELECT version FROM migration_version").Scan(&dbVersion); err != nil {
		return -1, nil
	}

	if dbVersion != oldMigrationLastVersion {
		return -1, fmt.Errorf("Must upgrade from db version %d (concourse 3.6.0), current db version: %d", oldMigrationLastVersion, dbVersion)
	}

	if _, err = self.db.Exec("DROP TABLE IF EXISTS migration_version"); err != nil {
		return -1, err
	}

	return newMigrationStartVersion, nil
}

type filenames []string

func (m filenames) Len() int {
	return len(m)
}

func (m filenames) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m filenames) Less(i, j int) bool {
	m1, _ := source.Parse(m[i])
	m2, _ := source.Parse(m[j])
	return m1.Version < m2.Version
}

func (m filenames) Latest() string {
	matches := []string{}

	for _, match := range m {
		if _, err := source.Parse(match); err == nil {
			matches = append(matches, match)
		}
	}

	sort.Sort(filenames(matches))

	return matches[len(matches)-1]
}
