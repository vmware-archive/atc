package db

import (
	"database/sql"
	"encoding/json"

	"github.com/concourse/atc/db/lock"
)

type JobCombination interface {
	ID() string
	JobID() int
	ResourceSpaceID() int
	ResourceSpaces() map[string]string
}

type jobCombination struct {
	id              string
	jobID           int
	resourceSpaceID int
	resourceSpaces  map[string]string
}

func (c *jobCombination) ID() string {
	return c.id
}

func (c *jobCombination) JobID() int {
	return c.jobID
}

func (c *jobCombination) ResourceSpaceID() int {
	return c.resourceSpaceID
}

func (c *jobCombination) ResourceSpaces() map[string]string {
	return c.resourceSpaces
}

func scanJobCombination(c *jobCombination, row scannable) error {
	var resourceSpacesBlob []byte

	err := row.Scan(&c.id, &c.jobID, &resourceSpacesBlob)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resourceSpacesBlob, &c.resourceSpaces)
	return err
}

func scanJobCombinations(conn Conn, lockFactory lock.LockFactory, rows *sql.Rows) ([]JobCombination, error) {
	defer Close(rows)

	jobCombinations := []JobCombination{}

	for rows.Next() {
		jobCombination := &jobCombination{}

		err := scanJobCombination(jobCombination, rows)
		if err != nil {
			return nil, err
		}

		jobCombinations = append(jobCombinations, jobCombination)
	}

	return jobCombinations, nil
}
