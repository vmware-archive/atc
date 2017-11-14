package migrations

import "github.com/concourse/atc/db/migration"

func AddAttemptsToVersionedResources(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
    ALTER TABLE versioned_resources ADD COLUMN attempts int DEFAULT(0) NOT NULL
  `)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    CREATE UNIQUE INDEX versioned_resources_resource_id_type_version_attempts
      ON versioned_resources (resource_id, type, version, attempts)
  `)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    DROP INDEX versioned_resources_resource_id_type_version
  `)

	return err
}
