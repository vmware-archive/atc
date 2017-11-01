package migrations

import "github.com/concourse/atc/db/migration"

func AddAuthFieldsToTeams(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
    ALTER TABLE teams
    ADD COLUMN basic_auth json null;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    ALTER TABLE teams
    ADD COLUMN github_auth json null;
  `)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    ALTER TABLE teams
    ADD COLUMN ldap_basic_auth text null;
  `)

	_, err = tx.Exec(`
		ALTER TABLE teams
		ADD COLUMN ldap_nonce text null;
	`)

	return err
}
