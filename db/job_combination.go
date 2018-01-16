package db

import (
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/db/lock"
)

type JobCombination interface {
	ID() int
	JobID() int
	Combination() map[string]string

	CreateBuild() (Build, error)
	EnsurePendingBuildExists() error

	GetNextBuildInputs() ([]BuildInput, bool, error)
	GetIndependentBuildInputs() ([]BuildInput, error)

	SaveNextInputMapping(inputMapping algorithm.InputMapping) error
	SaveIndependentInputMapping(inputMapping algorithm.InputMapping) error
	DeleteNextInputMapping() error
}

type jobCombination struct {
	id          int
	jobID       int
	combination map[string]string

	pipelineID int
	teamID     int

	conn        Conn
	lockFactory lock.LockFactory
}

func (c *jobCombination) ID() int {
	return c.id
}

func (c *jobCombination) JobID() int {
	return c.jobID
}

func (c *jobCombination) Combination() map[string]string {
	return c.combination
}

func (c *jobCombination) CreateBuild() (Build, error) {
	tx, err := c.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer Rollback(tx)

	buildName, err := c.getNewBuildName(tx)
	if err != nil {
		return nil, err
	}

	build := &build{conn: c.conn, lockFactory: c.lockFactory}
	err = createBuild(tx, build, map[string]interface{}{
		"name":               buildName,
		"job_combination_id": c.id,
		"pipeline_id":        c.pipelineID,
		"team_id":            c.teamID,
		"status":             BuildStatusPending,
		"manually_triggered": true,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY next_builds_per_job_combination`)
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (c *jobCombination) EnsurePendingBuildExists() error {
	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}

	defer Rollback(tx)

	buildName, err := c.getNewBuildName(tx)
	if err != nil {
		return err
	}

	rows, err := tx.Query(`
		INSERT INTO builds (name, job_combination_id, pipeline_id, team_id, status)
		SELECT $1, $2, $3, $4, 'pending'
		WHERE NOT EXISTS
			(SELECT id FROM builds WHERE job_combination_id = $2 AND status = 'pending')
		RETURNING id
	`, buildName, c.ID(), c.pipelineID, c.teamID)
	if err != nil {
		return err
	}

	defer Close(rows)

	if rows.Next() {
		var buildID int
		err := rows.Scan(&buildID)
		if err != nil {
			return err
		}

		err = rows.Close()
		if err != nil {
			return err
		}

		err = createBuildEventSeq(tx, buildID)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	return nil
}

func (c *jobCombination) GetNextBuildInputs() ([]BuildInput, bool, error) {
	var found bool
	err := psql.Select("inputs_determined").
		From("job_combinations").
		Where(sq.Eq{"id": c.id}).
		RunWith(c.conn).
		QueryRow().
		Scan(&found)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	// there is a possible race condition where found is true at first but the
	// inputs are deleted by the time we get here
	buildInputs, err := c.getBuildInputs("next_build_inputs")
	return buildInputs, true, err
}

func (c *jobCombination) GetIndependentBuildInputs() ([]BuildInput, error) {
	return c.getBuildInputs("independent_build_inputs")
}

func (c *jobCombination) SaveNextInputMapping(inputMapping algorithm.InputMapping) error {
	return c.saveJobInputMapping("next_build_inputs", inputMapping)
}

func (c *jobCombination) SaveIndependentInputMapping(inputMapping algorithm.InputMapping) error {
	return c.saveJobInputMapping("independent_build_inputs", inputMapping)
}

func (c *jobCombination) DeleteNextInputMapping() error {
	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}

	defer Rollback(tx)

	_, err = psql.Update("job_combinations").
		Set("inputs_determined", false).
		Where(sq.Eq{"id": c.id}).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	_, err = psql.Delete("next_build_inputs").
		Where(sq.Eq{"job_combination_id": c.id}).
		RunWith(tx).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (c *jobCombination) getNewBuildName(tx Tx) (string, error) {
	var buildName string
	err := psql.Update("jobs").
		Set("build_number_seq", sq.Expr("build_number_seq + 1")).
		Where(sq.Eq{"id": c.jobID}).
		Suffix("RETURNING build_number_seq").
		RunWith(tx).
		QueryRow().
		Scan(&buildName)

	return buildName, err
}

func (c *jobCombination) getBuildInputs(table string) ([]BuildInput, error) {
	rows, err := psql.Select("i.input_name, i.first_occurrence, r.name, vr.type, vr.version, vr.metadata").
		From(table + " i").
		Join("job_combinations c ON c.id = i.job_combination_id").
		Join("versioned_resources vr ON vr.id = i.version_id").
		Join("resource_spaces rs ON rs.id = vr.resource_space_id").
		Join("resources r ON r.id = rs.resource_id").
		Where(sq.Eq{"c.id": c.id}).
		RunWith(c.conn).
		Query()
	if err != nil {
		return nil, err
	}

	buildInputs := []BuildInput{}
	for rows.Next() {
		var (
			inputName       string
			firstOccurrence bool
			resourceName    string
			resourceType    string
			versionBlob     string
			metadataBlob    string
			version         ResourceVersion
			metadata        []ResourceMetadataField
		)

		err := rows.Scan(&inputName, &firstOccurrence, &resourceName, &resourceType, &versionBlob, &metadataBlob)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(versionBlob), &version)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(metadataBlob), &metadata)
		if err != nil {
			return nil, err
		}

		buildInputs = append(buildInputs, BuildInput{
			Name: inputName,
			VersionedResource: VersionedResource{
				Resource: resourceName,
				Type:     resourceType,
				Version:  version,
				Metadata: metadata,
			},
			FirstOccurrence: firstOccurrence,
		})
	}
	return buildInputs, nil
}

func (c *jobCombination) saveJobInputMapping(table string, inputMapping algorithm.InputMapping) error {
	tx, err := c.conn.Begin()
	if err != nil {
		return err
	}

	defer Rollback(tx)

	if table == "next_build_inputs" {
		_, err = psql.Update("job_combinations").
			Set("inputs_determined", true).
			Where(sq.Eq{"id": c.id}).
			Where(sq.Expr("NOT inputs_determined")).
			RunWith(tx).
			Exec()
	}
	if err != nil {
		return err
	}

	rows, err := psql.Select("input_name, version_id, first_occurrence").
		From(table).
		Where(sq.Eq{"job_combination_id": c.ID()}).
		RunWith(tx).
		Query()
	if err != nil {
		return err
	}

	oldInputMapping := algorithm.InputMapping{}
	for rows.Next() {
		var inputName string
		var inputVersion algorithm.InputVersion
		err = rows.Scan(&inputName, &inputVersion.VersionID, &inputVersion.FirstOccurrence)
		if err != nil {
			return err
		}

		oldInputMapping[inputName] = inputVersion
	}

	for inputName, oldInputVersion := range oldInputMapping {
		inputVersion, found := inputMapping[inputName]
		if !found || inputVersion != oldInputVersion {
			_, err = psql.Delete(table).
				Where(sq.Eq{
					"job_combination_id": c.id,
					"input_name":         inputName,
				}).
				RunWith(tx).
				Exec()
			if err != nil {
				return err
			}
		}
	}

	for inputName, inputVersion := range inputMapping {
		oldInputVersion, found := oldInputMapping[inputName]
		if !found || inputVersion != oldInputVersion {
			_, err := psql.Insert(table).
				SetMap(map[string]interface{}{
					"job_combination_id": c.id,
					"input_name":         inputName,
					"version_id":         inputVersion.VersionID,
					"first_occurrence":   inputVersion.FirstOccurrence,
				}).
				RunWith(tx).
				Exec()
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
