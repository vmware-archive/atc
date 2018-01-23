package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/lock"
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

	ResourceSpaceCombinations(map[string][]string) []map[string]string
	JobCombination() (JobCombination, error)
	Builds(page Page) ([]Build, Pagination, error)
	Build(name string) (Build, bool, error)
	FinishedAndNextBuild() (Build, Build, error)
	UpdateFirstLoggedBuildID(newFirstLoggedBuildID int) error
	GetPendingBuilds() ([]Build, error)

	SetMaxInFlightReached(bool) error
	GetRunningBuildsBySerialGroup(serialGroups []string) ([]Build, error)
	GetNextPendingBuildBySerialGroup(serialGroups []string) (Build, bool, error)
	SyncResourceSpaceCombinations([]map[string]string) ([]JobCombination, error)
}

var jobsQuery = psql.Select("j.id", "j.name", "j.config", "j.paused", "j.first_logged_build_id", "j.pipeline_id", "p.name", "p.team_id", "t.name", "j.nonce").
	From("jobs j, pipelines p").
	LeftJoin("teams t ON p.team_id = t.id").
	Where(sq.Expr("j.pipeline_id = p.id"))

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
		Where(sq.Eq{
			"j.name":        j.name,
			"j.pipeline_id": j.pipelineID,
		}).
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
			"j.name":        j.name,
			"j.pipeline_id": j.pipelineID,
			"b.status":      []BuildStatus{BuildStatusPending, BuildStatusStarted},
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

func (j *job) ResourceSpaceCombinations(resourceSpaces map[string][]string) []map[string]string {
	var resource, space string
	var spaces []string

	combinations := []map[string]string{}

	if len(resourceSpaces) == 0 {
		return []map[string]string{map[string]string{}}
	}

	if len(resourceSpaces) == 1 {
		for resource, spaces = range resourceSpaces {
			for _, space = range spaces {
				combinations = append(combinations, map[string]string{resource: space})
			}
		}
		return combinations
	}

	for resource, spaces = range resourceSpaces {
		break
	}
	delete(resourceSpaces, resource)

	for _, combination := range j.ResourceSpaceCombinations(resourceSpaces) {
		for _, space = range spaces {
			copy := map[string]string{}
			for k, v := range combination {
				copy[k] = v
			}

			copy[resource] = space
			combinations = append(combinations, copy)
		}
	}

	return combinations
}

func (j *job) JobCombination() (JobCombination, error) {
	var jobCombinationID, jobID int
	var marshaled sql.NullString
	var combination map[string]string

	err := psql.Select("id, job_id, combination").
		From("job_combinations").
		Where(sq.Eq{"job_id": j.id}).
		OrderBy("id DESC").Limit(1).
		RunWith(j.conn).QueryRow().
		Scan(&jobCombinationID, &jobID, &marshaled)
	if err != nil {
		return nil, err
	}

	if marshaled.Valid {
		err = json.Unmarshal([]byte(marshaled.String), &combination)
		if err != nil {
			return nil, err
		}
	}

	jc := &jobCombination{conn: j.conn, lockFactory: j.lockFactory, id: jobCombinationID, jobID: jobID, pipelineID: j.pipelineID, teamID: j.teamID, combination: combination}
	return jc, nil
}

func (j *job) Builds(page Page) ([]Build, Pagination, error) {
	query := buildsQuery.Where(sq.Eq{"j.id": j.id})

	limit := uint64(page.Limit)

	var args []interface{}

	var reverse bool
	if page.Until != 0 {
		query = query.Where(sq.Gt{"b.id": page.Until}).OrderBy("b.id ASC").Limit(limit)
		reverse = true
	} else if page.From != 0 {
		query = query.Where(sq.LtOrEq{"b.id": page.From}).OrderBy("b.id DESC").Limit(limit)
	} else if page.To != 0 {
		query = query.Where(sq.GtOrEq{"b.id": page.To}).OrderBy("b.id ASC").Limit(limit)
		reverse = true
	} else if page.Since != 0 {
		query = query.Where(sq.Lt{"b.id": page.Since}).OrderBy("b.id DESC").Limit(limit)
	} else if page.Around != 0 {
		selectBuildIds := sq.Select("b.id from builds b").JoinClause("LEFT OUTER JOIN job_combinations c ON c.id = b.job_combination_id").Where(sq.Eq{"c.job_id": j.id})

		fromQuery, fromArgs, err := selectBuildIds.Where(sq.GtOrEq{"b.id": page.Around}).OrderBy("b.id ASC").Limit(limit).ToSql()
		if err != nil {
			return nil, Pagination{}, err
		}
		args = append(args, fromArgs...)

		sinceQuery, sinceArgs, err := selectBuildIds.Where(sq.Lt{"b.id": page.Around}).OrderBy("b.id DESC").Limit(limit).ToSql()
		if err != nil {
			return nil, Pagination{}, err
		}

		args = append(args, sinceArgs...)

		query = query.Where("b.id IN ( (" + fromQuery + ") UNION (" + sinceQuery + ") )").OrderBy("b.id DESC")
	} else {
		query = query.OrderBy("b.id DESC").Limit(limit)
	}

	queryString, qargs, err := query.ToSql()
	if err != nil {
		return nil, Pagination{}, err
	}

	args = append(qargs, args...)

	rows, err := j.conn.Query(queryString, args...)
	if err != nil {
		return nil, Pagination{}, err
	}

	defer Close(rows)

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

	if page.Around != 0 && len(builds) > int(limit) {
		buildsLeft := []Build{}
		buildsRight := []Build{}

		for i, b := range builds {
			if b.ID() != page.Around {
				buildsLeft = append(buildsLeft, b)
			} else {
				buildsLeft = append(buildsLeft, b)
				buildsRight = builds[i+1:]
				break
			}
		}

		if len(buildsLeft) == len(buildsRight) {
			half := float64(len(buildsRight)) / 2
			builds = append(buildsLeft[int(half):], buildsRight[:int(half)]...)
		} else if len(buildsLeft) < len(buildsRight) {
			remainder := len(buildsRight) - len(buildsLeft)
			builds = append(buildsLeft, buildsRight[:remainder]...)
		} else if len(buildsLeft) > len(buildsRight) {
			remainder := len(buildsLeft) - len(buildsRight)
			builds = append(buildsLeft[len(buildsLeft)-remainder:], buildsRight...)
		}

	}

	var maxID, minID int
	err = psql.Select("COALESCE(MAX(b.id), 0) as maxID", "COALESCE(MIN(b.id), 0) as minID").
		From("builds b").
		Join("job_combinations c ON c.id = b.job_combination_id").
		Join("jobs j ON c.job_id = j.id").
		Where(sq.Eq{
			"j.name":        j.name,
			"j.pipeline_id": j.pipelineID,
		}).
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
	var query sq.SelectBuilder

	if name == "latest" {
		query = buildsQuery.
			Where(sq.Eq{"j.id": j.id}).
			OrderBy("b.id DESC").
			Limit(1)
	} else {
		query = buildsQuery.Where(sq.Eq{
			"j.id":   j.id,
			"b.name": name,
		})
	}

	row := query.RunWith(j.conn).QueryRow()

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
			"c.inputs_determined": true,
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

	defer Close(rows)

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
		"j.name":        j.name,
		"j.active":      true,
		"j.pipeline_id": j.pipelineID,
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

	defer Close(rows)

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

func (j *job) SyncResourceSpaceCombinations(combinations []map[string]string) ([]JobCombination, error) {
	tx, err := j.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer Rollback(tx)

	jobCombinations, err := j.syncResourceSpaceCombinations(tx, combinations)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return jobCombinations, nil
}

func (j *job) updateSerialGroups(serialGroups []string) error {
	tx, err := j.conn.Begin()
	if err != nil {
		return err
	}

	defer Rollback(tx)

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

	return tx.Commit()
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
	defer Close(rows)

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

func (j *job) syncResourceSpaceCombinations(tx Tx, combinations []map[string]string) ([]JobCombination, error) {
	jobCombinations := []JobCombination{}

	for _, c := range combinations {
		var jobCombinationID int

		marshaled, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}

		needsMigration, err := findInvalidJobCombination(tx, j.ID())
		if err != nil {
			return nil, err
		}

		if needsMigration {
			err = psql.Update("job_combinations").
				Set("combination", string(marshaled)).
				Where(sq.Eq{"job_id": j.ID()}).
				Suffix("RETURNING id").
				RunWith(tx).QueryRow().Scan(&jobCombinationID)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = psql.Insert("job_combinations").
				Columns("job_id", "combination").
				Values(j.ID(), string(marshaled)).
				Suffix("ON CONFLICT (job_id, combination) DO NOTHING").
				RunWith(tx).Exec()
			if err != nil {
				return nil, err
			}

			err = psql.Select("id").
				From("job_combinations").
				Where(sq.Eq{
					"job_id":      j.ID(),
					"combination": string(marshaled),
				}).RunWith(tx).QueryRow().Scan(&jobCombinationID)
			if err != nil {
				return nil, err
			}
		}

		for resource, space := range c {
			var resourceSpaceID int

			err := psql.Select("rs.id").
				From("resource_spaces rs").
				Join("resources r ON r.id = rs.resource_id").
				Where(sq.Eq{
					"r.name":        resource,
					"r.pipeline_id": j.PipelineID(),
					"rs.name":       space,
				}).RunWith(tx).QueryRow().Scan(&resourceSpaceID)
			if err != nil {
				return nil, err
			}

			_, err = psql.Insert("job_combinations_resource_spaces").
				Columns("job_combination_id", "resource_space_id").
				Values(jobCombinationID, resourceSpaceID).
				Suffix("ON CONFLICT (job_combination_id, resource_space_id) DO NOTHING").
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, err
			}
		}

		jobCombination := &jobCombination{conn: j.conn, lockFactory: j.lockFactory, id: jobCombinationID, jobID: j.id, combination: c, pipelineID: j.PipelineID(), teamID: j.TeamID()}

		jobCombinations = append(jobCombinations, jobCombination)
	}

	return jobCombinations, nil
}

func findInvalidJobCombination(tx Tx, jobID int) (bool, error) {
	var combination sql.NullString
	combinationsPresent := false

	rows, err := psql.Select("combination").
		From("job_combinations").
		Where(sq.Eq{"job_id": jobID}).
		Limit(1).RunWith(tx).Query()
	if err != nil {
		return false, err
	}

	for rows.Next() {
		combinationsPresent = true
		err = rows.Scan(&combination)
		if err != nil {
			return false, err
		}
	}

	return combinationsPresent && !combination.Valid, nil
}
