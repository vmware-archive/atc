package db

import (
	"database/sql"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc/db/algorithm"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/space"
)

//go:generate counterfeiter . Job

type JobPermutation interface {
	ID() int
	ResourceSpaces() space.Permutation

	CreateBuild() (Build, error)
	EnsurePendingBuildExists() error

	GetIndependentBuildInputs() (algorithm.InputMapping, error)
	GetNextBuildInputs() (algorithm.InputMapping, bool, error)
	SaveNextInputMapping(inputMapping algorithm.InputMapping) error
	SaveIndependentInputMapping(inputMapping algorithm.InputMapping) error
	DeleteNextInputMapping() error
}

var jobPermutationsQuery = psql.Select("jp.id, jp.resource_spaces, j.id, p.id, t.id").
	From("job_permutations jp").
	Join("jobs j ON j.id = jp.job_id").
	Join("pipelines p ON p.id = j.pipeline_id").
	Join("teams t ON t.id = p.team_id")

type jobPermutation struct {
	id             int
	resourceSpaces space.Permutation
	jobID          int
	pipelineID     int
	teamID         int

	conn        Conn
	lockFactory lock.LockFactory
}

func (j *jobPermutation) ID() int                           { return j.id }
func (j *jobPermutation) ResourceSpaces() space.Permutation { return j.resourceSpaces }

func (j *jobPermutation) CreateBuild() (Build, error) {
	tx, err := j.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	buildName, err := j.getNewBuildName(tx)
	if err != nil {
		return nil, err
	}

	build := &build{conn: j.conn, lockFactory: j.lockFactory}
	err = createBuild(tx, build, map[string]interface{}{
		"name":               buildName,
		"job_permutation_id": j.id,
		"pipeline_id":        j.pipelineID,
		"team_id":            j.teamID,
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

	_, err = j.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY next_builds_per_job_permutation`)
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (j *jobPermutation) EnsurePendingBuildExists() error {
	tx, err := j.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	buildName, err := j.getNewBuildName(tx)
	if err != nil {
		return err
	}

	rows, err := tx.Query(`
		INSERT INTO builds (name, job_permutation_id, pipeline_id, team_id, status)
		SELECT $1, $2, $3, $4, 'pending'
		WHERE NOT EXISTS
			(SELECT id FROM builds WHERE job_permutation_id = $2 AND status = 'pending')
		RETURNING id
	`, buildName, j.id, j.pipelineID, j.teamID)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		var buildID int
		err := rows.Scan(&buildID)
		if err != nil {
			return err
		}

		rows.Close()

		err = createBuildEventSeq(tx, buildID)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	return nil
}

func (j *jobPermutation) SaveIndependentInputMapping(inputMapping algorithm.InputMapping) error {
	return j.saveJobInputMapping("independent_build_inputs", inputMapping)
}

func (j *jobPermutation) SaveNextInputMapping(inputMapping algorithm.InputMapping) error {
	return j.saveJobInputMapping("next_build_inputs", inputMapping)
}

func (j *jobPermutation) GetIndependentBuildInputs() (algorithm.InputMapping, error) {
	return j.getBuildInputs("independent_build_inputs")
}

func (j *jobPermutation) GetNextBuildInputs() (algorithm.InputMapping, bool, error) {
	var found bool
	err := psql.Select("inputs_determined").
		From("job_permutations").
		Where(sq.Eq{"id": j.id}).
		RunWith(j.conn).
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
	buildInputs, err := j.getBuildInputs("next_build_inputs")
	return buildInputs, true, err
}

func (j *jobPermutation) DeleteNextInputMapping() error {
	tx, err := j.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = psql.Update("job_permutations").
		Set("inputs_determined", false).
		Where(sq.Eq{"id": j.id}).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	_, err = psql.Delete("next_build_inputs").
		Where(sq.Eq{"job_permutation_id": j.id}).
		RunWith(tx).Exec()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (j *jobPermutation) getBuildInputs(table string) (algorithm.InputMapping, error) {
	rows, err := psql.Select("input_name", "version_id", "first_occurrence").
		From(table + " i").
		Where(sq.Eq{"job_permutation_id": j.id}).
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	buildInputs := algorithm.InputMapping{}
	for rows.Next() {
		var (
			inputName       string
			versionID       int
			firstOccurrence bool
		)

		err := rows.Scan(&inputName, &versionID, &firstOccurrence)
		if err != nil {
			return nil, err
		}

		buildInputs[inputName] = algorithm.InputVersion{
			VersionID:       versionID,
			FirstOccurrence: firstOccurrence,
		}
	}

	return buildInputs, nil
}

func (j *jobPermutation) saveJobInputMapping(table string, inputMapping algorithm.InputMapping) error {
	tx, err := j.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if table == "next_build_inputs" {
		_, err = psql.Update("job_permutations").
			Set("inputs_determined", true).
			Where(sq.Eq{"id": j.id}).
			Where(sq.Expr("NOT inputs_determined")).
			RunWith(tx).
			Exec()
	}
	if err != nil {
		return err
	}

	rows, err := psql.Select("input_name, version_id, first_occurrence").
		From(table).
		Where(sq.Eq{"job_permutation_id": j.id}).
		RunWith(tx).
		Query()
	if err != nil {
		return err
	}

	oldInputMapping := algorithm.InputMapping{}
	for rows.Next() {
		var inputName string
		var inputVersion algorithm.InputVersion
		err := rows.Scan(&inputName, &inputVersion.VersionID, &inputVersion.FirstOccurrence)
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
					"job_permutation_id": j.id,
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
					"job_permutation_id": j.id,
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

func (j *jobPermutation) getNewBuildName(tx Tx) (string, error) {
	var buildName string
	err := psql.Update("jobs").
		Set("build_number_seq", sq.Expr("build_number_seq + 1")).
		Where(sq.Eq{"id": j.jobID}).
		Suffix("RETURNING build_number_seq").
		RunWith(tx).
		QueryRow().
		Scan(&buildName)

	return buildName, err
}

func scanJobPermutation(j *jobPermutation, row scannable) error {
	var (
		resourceSpacesBlob []byte
	)

	err := row.Scan(&j.id, &resourceSpacesBlob, &j.jobID, &j.pipelineID, &j.teamID)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resourceSpacesBlob, &j.resourceSpaces)
	if err != nil {
		return err
	}

	return nil
}

func scanJobPermutations(conn Conn, lockFactory lock.LockFactory, rows *sql.Rows) ([]JobPermutation, error) {
	defer rows.Close()

	jobPermutations := []JobPermutation{}

	for rows.Next() {
		jobPermutation := &jobPermutation{conn: conn, lockFactory: lockFactory}

		err := scanJobPermutation(jobPermutation, rows)
		if err != nil {
			return nil, err
		}

		jobPermutations = append(jobPermutations, jobPermutation)
	}

	return jobPermutations, nil
}
