package migrator

import (
	"strings"

	"github.com/concourse/atc/db/migration"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(fileNames []string) ([]string, error) {

	var migrationStatements []string
	for _, migrationFileName := range fileNames {
		migrationFileContents, err := migration.Asset(migrationFileName)
		if err != nil {
			return nil, err
		}

		migrationStatements = append(migrationStatements, strings.Split(string(migrationFileContents), ";")...)
		// last string is whitespace
		if strings.TrimSpace(migrationStatements[len(migrationStatements)-1]) == "" {
			migrationStatements = migrationStatements[:len(migrationStatements)-1]
		}
	}
	return migrationStatements, nil
}
