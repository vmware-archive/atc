package db

import (
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
)

var (
	ErrWorkerNotPresent         = errors.New("worker-not-present-in-db")
	ErrCannotPruneRunningWorker = errors.New("worker-not-stalled-for-pruning")
)

type WorkerState string

const (
	WorkerStateRunning  = WorkerState("running")
	WorkerStateStalled  = WorkerState("stalled")
	WorkerStateLanding  = WorkerState("landing")
	WorkerStateLanded   = WorkerState("landed")
	WorkerStateRetiring = WorkerState("retiring")
)

//go:generate counterfeiter . Worker

type Worker interface {
	Land() error
	Retire() error
	Prune() error
	Delete() error
}

type worker struct {
	conn Conn

	name             string
	version          *string
	state            WorkerState
	gardenAddr       *string
	baggageclaimURL  *string
	httpProxyURL     string
	httpsProxyURL    string
	noProxy          string
	activeContainers int
	resourceTypes    []atc.WorkerResourceType
	platform         string
	tags             []string
	teamID           int
	teamName         string
	startTime        int64
	expiresAt        time.Time
	certsPath        *string
}

func (worker *worker) Land() error {
	cSQL, _, err := sq.Case("state").
		When("'landed'::worker_state", "'landed'::worker_state").
		Else("'landing'::worker_state").
		ToSql()
	if err != nil {
		return err
	}

	result, err := psql.Update("workers").
		Set("state", sq.Expr("("+cSQL+")")).
		Where(sq.Eq{"name": worker.name}).
		RunWith(worker.conn).
		Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrWorkerNotPresent
	}

	return nil
}

func (worker *worker) Retire() error {
	result, err := psql.Update("workers").
		SetMap(map[string]interface{}{
			"state": string(WorkerStateRetiring),
		}).
		Where(sq.Eq{"name": worker.name}).
		RunWith(worker.conn).
		Exec()
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrWorkerNotPresent
	}

	return nil
}

func (worker *worker) Prune() error {
	rows, err := sq.Delete("workers").
		Where(sq.Eq{
			"name": worker.name,
		}).
		Where(sq.NotEq{
			"state": string(WorkerStateRunning),
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(worker.conn).
		Exec()

	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		//check whether the worker exists in the database at all
		var one int
		err := psql.Select("1").From("workers").Where(sq.Eq{"name": worker.name}).
			RunWith(worker.conn).
			QueryRow().
			Scan(&one)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrWorkerNotPresent
			}
			return err
		}

		return ErrCannotPruneRunningWorker
	}

	return nil
}

func (worker *worker) ResourceCerts() (*UsedWorkerResourceCerts, bool, error) {
	if worker.certsPath != nil {
		wrc := &WorkerResourceCerts{
			WorkerName: worker.name,
			CertsPath:  *worker.certsPath,
		}

		return wrc.Find(worker.conn)
	}

	return nil, false, nil
}
