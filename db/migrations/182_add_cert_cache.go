package migrations

import "github.com/concourse/atc/db/migration"

func AddCertCache(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
		CREATE TABLE cert_cache (
			key text PRIMARY KEY,
			data text NOT NULL,
			nonce text
		)
	`)
	return err
}
