package migrations

import "github.com/concourse/atc/db/migration"

func AddVersionedResourcesCheckOrderIndex(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
		CREATE INDEX versioned_resources_check_order ON versioned_resources (check_order DESC);
	`)
	if err != nil {
		return err
	}

	return nil
}
