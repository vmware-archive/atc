package migrations

func (m *GoMigrationsRunner) Down_4000() error {

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`ALTER TABLE some_table DROP COLUMN name`)
	if err != nil {
		return err
	}

	return nil
}
