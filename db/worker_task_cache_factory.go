package db

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/lib/pq"
)

type UsedWorkerTaskCache struct {
	ID         int
	WorkerName string
}

//go:generate counterfeiter . WorkerTaskCacheFactory

type WorkerTaskCacheFactory interface {
	Find(jobCombinationID int, stepName string, path string, workerName string) (*UsedWorkerTaskCache, bool, error)
	FindOrCreate(jobCombinationID int, stepName string, path string, workerName string) (*UsedWorkerTaskCache, error)
}

type workerTaskCacheFactory struct {
	conn Conn
}

func NewWorkerTaskCacheFactory(conn Conn) WorkerTaskCacheFactory {
	return &workerTaskCacheFactory{
		conn: conn,
	}
}

func (f *workerTaskCacheFactory) Find(jobCombinationID int, stepName string, path string, workerName string) (*UsedWorkerTaskCache, bool, error) {
	var id int
	err := psql.Select("id").
		From("worker_task_caches").
		Where(sq.Eq{
			"job_combination_id": jobCombinationID,
			"step_name":          stepName,
			"worker_name":        workerName,
			"path":               path,
		}).
		RunWith(f.conn).
		QueryRow().
		Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &UsedWorkerTaskCache{
		ID:         id,
		WorkerName: workerName,
	}, true, nil
}

func (f *workerTaskCacheFactory) FindOrCreate(jobCombinationID int, stepName string, path string, workerName string) (*UsedWorkerTaskCache, error) {
	workerTaskCache := WorkerTaskCache{
		JobCombinationID: jobCombinationID,
		StepName:         stepName,
		WorkerName:       workerName,
		Path:             path,
	}

	var usedWorkerTaskCache *UsedWorkerTaskCache

	err := safeFindOrCreate(f.conn, func(tx Tx) error {
		var err error
		usedWorkerTaskCache, err = workerTaskCache.FindOrCreate(tx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return usedWorkerTaskCache, nil
}

type WorkerTaskCache struct {
	JobCombinationID int
	StepName         string
	WorkerName       string
	Path             string
}

func (wtc WorkerTaskCache) FindOrCreate(
	tx Tx,
) (*UsedWorkerTaskCache, error) {
	var id int
	err := psql.Select("id").
		From("worker_task_caches").
		Where(sq.Eq{
			"job_combination_id": wtc.JobCombinationID,
			"step_name":          wtc.StepName,
			"worker_name":        wtc.WorkerName,
			"path":               wtc.Path,
		}).
		RunWith(tx).
		QueryRow().
		Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = psql.Insert("worker_task_caches").
				Columns(
					"job_combination_id",
					"step_name",
					"worker_name",
					"path",
				).
				Values(
					wtc.JobCombinationID,
					wtc.StepName,
					wtc.WorkerName,
					wtc.Path,
				).
				Suffix("RETURNING id").
				RunWith(tx).
				QueryRow().
				Scan(&id)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == pqUniqueViolationErrCode {
					return nil, ErrSafeRetryFindOrCreate
				}

				return nil, err
			}

			return &UsedWorkerTaskCache{
				ID:         id,
				WorkerName: wtc.WorkerName,
			}, nil
		}

		return nil, err
	}

	return &UsedWorkerTaskCache{
		ID:         id,
		WorkerName: wtc.WorkerName,
	}, nil
}

func removeUnusedWorkerTaskCaches(tx Tx, pipelineID int, jobConfigs []atc.JobConfig) error {
	steps := make(map[string][]string)
	for _, jobConfig := range jobConfigs {
		for _, jobConfigPlan := range jobConfig.Plan {
			if jobConfigPlan.Task != "" {
				steps[jobConfig.Name] = append(steps[jobConfig.Name], jobConfigPlan.Task)
			}
		}
	}

	query := sq.Or{}
	for jobName, stepNames := range steps {
		query = append(query, sq.And{sq.Eq{"j.name": jobName}, sq.NotEq{"wtc.step_name": stepNames}})
	}

	_, err := psql.Delete("worker_task_caches wtc USING jobs j, job_combinations c").
		Where(sq.Or{
			query,
			sq.Eq{
				"j.pipeline_id": pipelineID,
				"j.active":      false,
			},
		}).
		Where(sq.Expr("c.id = wtc.job_combination_id")).
		Where(sq.Expr("c.job_id = j.id")).
		RunWith(tx).
		Exec()

	return err
}
