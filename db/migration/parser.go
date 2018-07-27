package migration

import (
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(fileName string) ([]string, error) {

	var fileStatements []string
	migrationFileContents, err := Asset(fileName)
	if err != nil {
		return nil, err
	}
	var migrationStatements = []string{}
	fileStatements = append(fileStatements, strings.Split(string(migrationFileContents), ";")...)
	// last string is whitespace
	if strings.TrimSpace(fileStatements[len(fileStatements)-1]) == "" {
		fileStatements = fileStatements[:len(fileStatements)-1]
	}

	var isSqlStatement bool = false
	var sqlStatement string
	for _, statement := range fileStatements {
		statement = strings.TrimSpace(statement)

		if statement == "BEGIN" || statement == "COMMIT" {
			continue
		}
		if strings.Contains(statement, "BEGIN") {
			isSqlStatement = true
			sqlStatement = statement + ";"
		} else {

			if isSqlStatement {
				sqlStatement = strings.Join([]string{sqlStatement, statement, ";"}, "")
				if strings.HasPrefix(statement, "$$") {
					migrationStatements = append(migrationStatements, sqlStatement)
					isSqlStatement = false
				}
			} else {
				migrationStatements = append(migrationStatements, statement)
			}
		}
		// if strings.Contains(statement, "BEGIN") {
		// 	isSqlStatement = true
		// 	sqlStatement = statement + ";"
		// 	continue
		// }
		// if isSqlStatement {
		// 	sqlStatement = strings.Join([]string{sqlStatement, statement, ";"}, "")
		//
		// 	if strings.HasPrefix(statement, "$$") {
		// 		isSqlStatement = false
		// 		statement = sqlStatement
		// 		fmt.Println(statement)
		// 	} else {
		// 		continue
		// 	}
		// }
		// migrationStatements = append(migrationStatements, statement)
	}
	return migrationStatements, nil
}
