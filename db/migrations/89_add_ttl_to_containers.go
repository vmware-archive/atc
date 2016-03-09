package migrations

import "github.com/BurntSushi/migration"

func AddTTLToContainers(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
		ALTER TABLE containers ADD COLUMN ttl text;
	`)
	if err != nil {
		return err
	}

	return nil
}
