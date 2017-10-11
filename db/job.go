package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/space"
)

//go:generate counterfeiter . Job

type Job interface {
	ID() int
	Name() string
	Paused() bool
	FirstLoggedBuildID() int
	PipelineID() int
	PipelineName() string
	TeamID() int
	TeamName() string
	Config() atc.JobConfig

	Reload() (bool, error)

	Pause() error
	Unpause() error

	Builds(page Page) ([]Build, Pagination, error)
	Build(name string) (Build, bool, error)
	FinishedAndNextBuild() (Build, Build, error)
	UpdateFirstLoggedBuildID(newFirstLoggedBuildID int) error

	SyncPermutations([]space.Permutation) ([]JobPermutation, error)

	// XXX: move to permutation
	GetPendingBuilds() ([]Build, error)
	GetRunningBuildsBySerialGroup(serialGroups []string) ([]Build, error)
	GetNextPendingBuildBySerialGroup(serialGroups []string) (Build, bool, error)

	SetMaxInFlightReached(bool) error
}

var jobsQuery = psql.Select("j.id", "j.name", "j.config", "j.paused", "j.first_logged_build_id", "p.id", "p.name", "p.team_id", "t.name", "j.nonce").
	From("jobs j").
	Join("pipelines p ON p.id = j.pipeline_id").
	Join("teams t ON p.team_id = t.id")

type FirstLoggedBuildIDDecreasedError struct {
	Job   string
	OldID int
	NewID int
}

func (e FirstLoggedBuildIDDecreasedError) Error() string {
	return fmt.Sprintf("first logged build id for job '%s' decreased from %d to %d", e.Job, e.OldID, e.NewID)
}

type job struct {
	id                 int
	name               string
	paused             bool
	firstLoggedBuildID int
	pipelineID         int
	pipelineName       string
	teamID             int
	teamName           string
	config             atc.JobConfig

	conn        Conn
	lockFactory lock.LockFactory
}

type Jobs []Job

func (jobs Jobs) Configs() atc.JobConfigs {
	var configs atc.JobConfigs

	for _, j := range jobs {
		configs = append(configs, j.Config())
	}

	return configs
}

func (j *job) ID() int                 { return j.id }
func (j *job) Name() string            { return j.name }
func (j *job) Paused() bool            { return j.paused }
func (j *job) FirstLoggedBuildID() int { return j.firstLoggedBuildID }
func (j *job) PipelineID() int         { return j.pipelineID }
func (j *job) PipelineName() string    { return j.pipelineName }
func (j *job) TeamID() int             { return j.teamID }
func (j *job) TeamName() string        { return j.teamName }
func (j *job) Config() atc.JobConfig   { return j.config }

func (j *job) Reload() (bool, error) {
	row := jobsQuery.Where(sq.Eq{"j.id": j.id}).
		RunWith(j.conn).
		QueryRow()

	err := scanJob(j, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (j *job) Pause() error {
	return j.updatePausedJob(true)
}

func (j *job) Unpause() error {
	return j.updatePausedJob(false)
}

func (j *job) FinishedAndNextBuild() (Build, Build, error) {
	row := buildsQuery.
		Where(sq.Eq{"j.id": j.id}).
		Where(sq.Expr("b.status NOT IN ('pending', 'started')")).
		OrderBy("b.id DESC").
		Limit(1).
		RunWith(j.conn).
		QueryRow()

	var finished, next Build

	finishedBuild := &build{conn: j.conn, lockFactory: j.lockFactory}
	err := scanBuild(finishedBuild, row, j.conn.EncryptionStrategy())
	if err == nil {
		finished = finishedBuild
	} else if err != sql.ErrNoRows {
		return nil, nil, err
	}

	row = buildsQuery.
		Where(sq.Eq{
			"j.id":     j.id,
			"b.status": []BuildStatus{BuildStatusPending, BuildStatusStarted},
		}).
		OrderBy("b.id ASC").
		Limit(1).
		RunWith(j.conn).
		QueryRow()

	nextBuild := &build{conn: j.conn, lockFactory: j.lockFactory}
	err = scanBuild(nextBuild, row, j.conn.EncryptionStrategy())
	if err == nil {
		next = nextBuild
	} else if err != sql.ErrNoRows {
		return nil, nil, err
	}

	return finished, next, nil
}

func (j *job) UpdateFirstLoggedBuildID(newFirstLoggedBuildID int) error {
	if j.firstLoggedBuildID > newFirstLoggedBuildID {
		return FirstLoggedBuildIDDecreasedError{
			Job:   j.name,
			OldID: j.firstLoggedBuildID,
			NewID: newFirstLoggedBuildID,
		}
	}

	result, err := psql.Update("jobs").
		Set("first_logged_build_id", newFirstLoggedBuildID).
		Where(sq.Eq{"id": j.id}).
		RunWith(j.conn).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return nonOneRowAffectedError{rowsAffected}
	}

	return nil
}

func (j *job) Builds(page Page) ([]Build, Pagination, error) {
	query := buildsQuery.Where(sq.Eq{"j.id": j.id})

	limit := uint64(page.Limit)

	var reverse bool
	if page.Since == 0 && page.Until == 0 {
		query = query.OrderBy("b.id DESC").Limit(limit)
	} else if page.Until != 0 {
		query = query.Where(sq.Gt{"b.id": page.Until}).OrderBy("b.id ASC").Limit(limit)
		reverse = true
	} else {
		query = query.Where(sq.Lt{"b.id": page.Since}).OrderBy("b.id DESC").Limit(limit)
	}

	rows, err := query.RunWith(j.conn).Query()
	if err != nil {
		return nil, Pagination{}, err
	}

	defer rows.Close()

	builds := []Build{}

	for rows.Next() {
		build := &build{conn: j.conn, lockFactory: j.lockFactory}
		err = scanBuild(build, rows, j.conn.EncryptionStrategy())
		if err != nil {
			return nil, Pagination{}, err
		}

		if reverse {
			builds = append([]Build{build}, builds...)
		} else {
			builds = append(builds, build)
		}
	}

	if len(builds) == 0 {
		return []Build{}, Pagination{}, nil
	}

	var maxID, minID int
	err = psql.Select("COALESCE(MAX(b.id), 0) as maxID", "COALESCE(MIN(b.id), 0) as minID").
		From("builds b").
		Join("job_permutations jp ON jp.id = b.job_permutation_id").
		Join("jobs j ON j.id = jp.job_id").
		Where(sq.Eq{"j.id": j.id}).
		RunWith(j.conn).
		QueryRow().
		Scan(&maxID, &minID)
	if err != nil {
		return nil, Pagination{}, err
	}

	firstBuild := builds[0]
	lastBuild := builds[len(builds)-1]

	var pagination Pagination

	if firstBuild.ID() < maxID {
		pagination.Previous = &Page{
			Until: firstBuild.ID(),
			Limit: page.Limit,
		}
	}

	if lastBuild.ID() > minID {
		pagination.Next = &Page{
			Since: lastBuild.ID(),
			Limit: page.Limit,
		}
	}

	return builds, pagination, nil
}

func (j *job) Build(name string) (Build, bool, error) {
	row := buildsQuery.Where(sq.Eq{
		"j.id":   j.id,
		"b.name": name,
	}).
		RunWith(j.conn).
		QueryRow()

	build := &build{conn: j.conn, lockFactory: j.lockFactory}

	err := scanBuild(build, row, j.conn.EncryptionStrategy())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return build, true, nil
}

// XXX: move to permutation
func (j *job) GetNextPendingBuildBySerialGroup(serialGroups []string) (Build, bool, error) {
	err := j.updateSerialGroups(serialGroups)
	if err != nil {
		return nil, false, err
	}

	row := buildsQuery.Options(`DISTINCT ON (b.id)`).
		Join(`jobs_serial_groups jsg ON j.id = jsg.job_id`).
		Where(sq.Eq{
			"jsg.serial_group":    serialGroups,
			"b.status":            BuildStatusPending,
			"j.paused":            false,
			"j.inputs_determined": true,
			"j.pipeline_id":       j.pipelineID}).
		OrderBy("b.id ASC").
		Limit(1).
		RunWith(j.conn).
		QueryRow()

	build := &build{conn: j.conn, lockFactory: j.lockFactory}
	err = scanBuild(build, row, j.conn.EncryptionStrategy())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return build, true, nil
}

func (j *job) GetRunningBuildsBySerialGroup(serialGroups []string) ([]Build, error) {
	err := j.updateSerialGroups(serialGroups)
	if err != nil {
		return nil, err
	}

	rows, err := buildsQuery.Options(`DISTINCT ON (b.id)`).
		Join(`jobs_serial_groups jsg ON j.id = jsg.job_id`).
		Where(sq.Eq{
			"jsg.serial_group": serialGroups,
			"j.pipeline_id":    j.pipelineID,
		}).
		Where(sq.Or{
			sq.Eq{"b.status": BuildStatusStarted},
			sq.Eq{"b.status": BuildStatusPending, "b.scheduled": true},
		}).
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bs := []Build{}

	for rows.Next() {
		build := &build{conn: j.conn, lockFactory: j.lockFactory}
		err = scanBuild(build, rows, j.conn.EncryptionStrategy())
		if err != nil {
			return nil, err
		}

		bs = append(bs, build)
	}

	return bs, nil
}

func (j *job) SetMaxInFlightReached(reached bool) error {
	result, err := psql.Update("jobs").
		Set("max_in_flight_reached", reached).
		Where(sq.Eq{
			"id": j.id,
		}).
		RunWith(j.conn).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return nonOneRowAffectedError{rowsAffected}
	}

	return nil
}

func (j *job) GetPendingBuilds() ([]Build, error) {
	builds := []Build{}

	row := jobsQuery.Where(sq.Eq{
		"j.id":     j.id,
		"j.active": true,
	}).RunWith(j.conn).QueryRow()

	job := &job{conn: j.conn, lockFactory: j.lockFactory}
	err := scanJob(job, row)
	if err != nil {
		return nil, err
	}

	rows, err := buildsQuery.
		Where(sq.Eq{
			"j.id":     j.id,
			"b.status": BuildStatusPending,
		}).
		OrderBy("b.id ASC").
		RunWith(j.conn).
		Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		build := &build{conn: j.conn, lockFactory: j.lockFactory}
		err = scanBuild(build, rows, j.conn.EncryptionStrategy())
		if err != nil {
			return nil, err
		}

		builds = append(builds, build)
	}

	return builds, nil
}

func (j *job) SyncPermutations(permutations []space.Permutation) ([]JobPermutation, error) {
	tx, err := j.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var keepIDs []int
	for _, p := range permutations {
		marshaled, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}

		var id int
		err = psql.Insert("job_permutations").
			Columns("job_id", "resource_spaces", "active").
			Values(j.id, marshaled, true).
			Suffix(`
				ON CONFLICT (job_id, resource_spaces) DO UPDATE SET active = true
				RETURNING id
			`).
			RunWith(tx).
			QueryRow().
			Scan(&id)
		if err != nil {
			return nil, err
		}

		keepIDs = append(keepIDs, id)
	}

	_, err = psql.Update("job_permutations").
		Set("active", false).
		Where(sq.Eq{"job_id": j.id}).
		Where(sq.NotEq{"id": keepIDs}).
		RunWith(tx).
		Exec()
	if err != nil {
		return nil, err
	}

	rows, err := jobPermutationsQuery.Where(sq.Eq{"jp.job_id": j.id, "jp.active": true}).RunWith(tx).Query()
	if err != nil {
		return nil, err
	}

	jobPermutations, err := scanJobPermutations(j.conn, j.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return jobPermutations, nil
}

func (j *job) updateSerialGroups(serialGroups []string) error {
	tx, err := j.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = psql.Delete("jobs_serial_groups").
		Where(sq.Eq{
			"job_id": j.id,
		}).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	for _, serialGroup := range serialGroups {
		_, err = psql.Insert("jobs_serial_groups (job_id, serial_group)").
			Values(j.id, serialGroup).
			RunWith(tx).
			Exec()
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (j *job) updatePausedJob(pause bool) error {
	result, err := psql.Update("jobs").
		Set("paused", pause).
		Where(sq.Eq{"id": j.id}).
		RunWith(j.conn).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return nonOneRowAffectedError{rowsAffected}
	}

	return nil
}

func scanJob(j *job, row scannable) error {
	var (
		configBlob []byte
		nonce      sql.NullString
	)

	err := row.Scan(&j.id, &j.name, &configBlob, &j.paused, &j.firstLoggedBuildID, &j.pipelineID, &j.pipelineName, &j.teamID, &j.teamName, &nonce)
	if err != nil {
		return err
	}

	es := j.conn.EncryptionStrategy()

	var noncense *string
	if nonce.Valid {
		noncense = &nonce.String
	}

	decryptedConfig, err := es.Decrypt(string(configBlob), noncense)
	if err != nil {
		return err
	}

	var config atc.JobConfig
	err = json.Unmarshal(decryptedConfig, &config)
	if err != nil {
		return err
	}

	j.config = config

	return nil
}

func scanJobs(conn Conn, lockFactory lock.LockFactory, rows *sql.Rows) (Jobs, error) {
	defer rows.Close()

	jobs := Jobs{}

	for rows.Next() {
		job := &job{conn: conn, lockFactory: lockFactory}

		err := scanJob(job, rows)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
