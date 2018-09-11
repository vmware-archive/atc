package migrations

import (
	"fmt"
	"time"
)

func (m *GoMigrationsRunner) Up_4000() error {
	type info struct {
		id     int
		tstamp time.Time
	}
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Commit()

	_, err = tx.Exec(`ALTER TABLE some_table ADD COLUMN name VARCHAR`)
	if err != nil {
		return err
	}

	rows, err := tx.Query("SELECT id FROM some_table")
	if err != nil {
		return err
	}

	infos := []info{}

	for rows.Next() {
		info := info{}
		err = rows.Scan(&info.id)
		if err != nil {
			return err
		}
		infos = append(infos, info)
	}

	for _, info := range infos {
		name := fmt.Sprintf("name_%v", info.id)
		tx.Exec(`UPDATE some_table SET name=$1 WHERE id=$2`, name, info.id)
	}

	return nil
}
