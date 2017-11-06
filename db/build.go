package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db/lock"
	"github.com/concourse/atc/event"
	"github.com/lib/pq"
)

type BuildStatus string

const (
	BuildStatusPending   BuildStatus = "pending"
	BuildStatusStarted   BuildStatus = "started"
	BuildStatusAborted   BuildStatus = "aborted"
	BuildStatusSucceeded BuildStatus = "succeeded"
	BuildStatusFailed    BuildStatus = "failed"
	BuildStatusErrored   BuildStatus = "errored"
)

var buildsQuery = psql.Select("b.id, b.name, b.job_id, b.team_id, b.status, b.manually_triggered, b.scheduled, b.engine, b.engine_metadata, b.public_plan, b.start_time, b.end_time, b.reap_time, j.name, b.pipeline_id, p.name, t.name, b.nonce").
	From("builds b").
	JoinClause("LEFT OUTER JOIN jobs j ON b.job_id = j.id").
	JoinClause("LEFT OUTER JOIN pipelines p ON b.pipeline_id = p.id").
	JoinClause("LEFT OUTER JOIN teams t ON b.team_id = t.id")

//go:generate counterfeiter . Build

type Build interface {
	ID() int
	Name() string
	JobID() int
	JobName() string
	PipelineID() int
	PipelineName() string
	TeamID() int
	TeamName() string
	Engine() string
	EngineMetadata() string
	PublicPlan() *json.RawMessage
	Status() BuildStatus
	StartTime() time.Time
	EndTime() time.Time
	ReapTime() time.Time
	IsManuallyTriggered() bool
	IsScheduled() bool

	IsRunning() bool

	Reload() (bool, error)

	Interceptible() (bool, error)
	AcquireTrackingLock(logger lager.Logger, interval time.Duration) (lock.Lock, bool, error)
	Preparation() (BuildPreparation, bool, error)

	Start(string, string, atc.Plan) (bool, error)
	FinishWithError(cause error) error
	Finish(BuildStatus) error

	SetInterceptible(bool) error

	Events(uint) (EventSource, error)
	SaveEvent(event atc.Event) error

	SaveInput(input BuildInput) error
	SaveOutput(vr VersionedResource, explicit bool) error
	UseInputs(inputs []BuildInput) error

	Resources() ([]BuildInput, []BuildOutput, error)
	GetVersionedResources() (SavedVersionedResources, error)
	SaveImageResourceVersion(*UsedResourceCache) error

	Pipeline() (Pipeline, bool, error)

	Delete() (bool, error)
	MarkAsAborted() error
	AbortNotifier() (Notifier, error)
	Schedule() (bool, error)
}

type build struct {
	id        int
	name      string
	status    BuildStatus
	scheduled bool

	teamID   int
	teamName string

	pipelineID   int
	pipelineName string
	jobID        int
	jobName      string

	isManuallyTriggered bool

	engine         string
	engineMetadata string
	publicPlan     *json.RawMessage

	startTime time.Time
	endTime   time.Time
	reapTime  time.Time

	conn        Conn
	lockFactory lock.LockFactory
}

var ErrBuildDisappeared = errors.New("build-disappeared-from-db")

func (b *build) ID() int                      { return b.id }
func (b *build) Name() string                 { return b.name }
func (b *build) JobID() int                   { return b.jobID }
func (b *build) JobName() string              { return b.jobName }
func (b *build) PipelineID() int              { return b.pipelineID }
func (b *build) PipelineName() string         { return b.pipelineName }
func (b *build) TeamID() int                  { return b.teamID }
func (b *build) TeamName() string             { return b.teamName }
func (b *build) IsManuallyTriggered() bool    { return b.isManuallyTriggered }
func (b *build) Engine() string               { return b.engine }
func (b *build) EngineMetadata() string       { return b.engineMetadata }
func (b *build) PublicPlan() *json.RawMessage { return b.publicPlan }
func (b *build) StartTime() time.Time         { return b.startTime }
func (b *build) EndTime() time.Time           { return b.endTime }
func (b *build) ReapTime() time.Time          { return b.reapTime }
func (b *build) Status() BuildStatus          { return b.status }
func (b *build) IsScheduled() bool            { return b.scheduled }

func (b *build) IsRunning() bool {
	switch b.status {
	case BuildStatusPending, BuildStatusStarted:
		return true
	default:
		return false
	}
}

func (b *build) Reload() (bool, error) {
	row := buildsQuery.Where(sq.Eq{"b.id": b.id}).
		RunWith(b.conn).
		QueryRow()

	err := scanBuild(b, row, b.conn.EncryptionStrategy())
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (b *build) Interceptible() (bool, error) {
	var interceptible bool

	err := psql.Select("interceptible").
		From("builds").
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		QueryRow().Scan(&interceptible)

	if err != nil {
		return true, err
	}

	return interceptible, nil
}

func (b *build) SetInterceptible(i bool) error {
	rows, err := psql.Update("builds").
		Set("interceptible", i).
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrBuildDisappeared
	}

	return nil

}

func (b *build) Start(engine, metadata string, plan atc.Plan) (bool, error) {
	tx, err := b.conn.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	encryptedMetadata, nonce, err := b.conn.EncryptionStrategy().Encrypt([]byte(metadata))
	if err != nil {
		return false, err
	}

	var startTime time.Time

	err = psql.Update("builds").
		Set("status", "started").
		Set("start_time", sq.Expr("now()")).
		Set("engine", engine).
		Set("engine_metadata", encryptedMetadata).
		Set("public_plan", plan.Public()).
		Set("nonce", nonce).
		Where(sq.Eq{
			"id":     b.id,
			"status": "pending",
		}).
		Suffix("RETURNING start_time").
		RunWith(tx).
		QueryRow().
		Scan(&startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = b.saveEvent(tx, event.Status{
		Status: atc.StatusStarted,
		Time:   startTime.Unix(),
	})
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	err = b.conn.Bus().Notify(buildEventsChannel(b.id))
	if err != nil {
		return false, err
	}

	_, err = b.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY next_builds_per_job`)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *build) FinishWithError(cause error) error {
	err := b.SaveEvent(event.Error{
		Message: cause.Error(),
	})
	if err != nil {
		return err
	}

	return b.Finish(BuildStatusErrored)
}

func (b *build) Finish(status BuildStatus) error {
	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var endTime time.Time

	err = psql.Update("builds").
		Set("status", status).
		Set("end_time", sq.Expr("now()")).
		Set("completed", true).
		Set("engine_metadata", nil).
		Set("nonce", nil).
		Where(sq.Eq{"id": b.id}).
		Suffix("RETURNING end_time").
		RunWith(tx).
		QueryRow().
		Scan(&endTime)
	if err != nil {
		return err
	}

	err = b.saveEvent(tx, event.Status{
		Status: atc.BuildStatus(status),
		Time:   endTime.Unix(),
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
		DROP SEQUENCE %s
	`, buildEventSeq(b.id)))
	if err != nil {
		return err
	}

	if b.jobID != 0 && status == BuildStatusSucceeded {
		_, err = psql.Delete("build_image_resource_caches birc USING builds b").
			Where(sq.Expr("birc.build_id = b.id")).
			Where(sq.Lt{"build_id": b.id}).
			Where(sq.Eq{"b.job_id": b.jobID}).
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

	err = b.conn.Bus().Notify(buildEventsChannel(b.id))
	if err != nil {
		return err
	}

	_, err = b.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY latest_completed_builds_per_job`)
	if err != nil {
		return err
	}

	_, err = b.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY next_builds_per_job`)
	if err != nil {
		return err
	}

	_, err = b.conn.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY transition_builds_per_job`)
	if err != nil {
		return err
	}

	return nil
}

func (b *build) Delete() (bool, error) {
	rows, err := psql.Delete("builds").
		Where(sq.Eq{
			"id": b.id,
		}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, ErrBuildDisappeared
	}

	return true, nil
}

// MarkAsAborted will send the abort notification to all build abort
// channel listeners. It will set the status to aborted that will make
// AbortNotifier send notification in case if tracking ATC misses the first
// notification on abort channel.
// Setting status as aborted will also make Start() return false in case where
// build was aborted before it was started.
func (b *build) MarkAsAborted() error {
	_, err := psql.Update("builds").
		Set("status", string(BuildStatusAborted)).
		Where(sq.Eq{"id": b.id}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return err
	}

	err = b.conn.Bus().Notify(buildAbortChannel(b.id))
	if err != nil {
		return err
	}

	return nil
}

// AbortNotifier returns a Notifier that can be watched for when the build
// is marked as aborted. Once the build is marked as aborted it will send a
// notification to finish the build to ATC that is tracking this build.
func (b *build) AbortNotifier() (Notifier, error) {
	return newConditionNotifier(b.conn.Bus(), buildAbortChannel(b.id), func() (bool, error) {
		var aborted bool
		err := psql.Select("status = 'aborted'").
			From("builds").
			Where(sq.Eq{"id": b.id}).
			RunWith(b.conn).
			QueryRow().
			Scan(&aborted)

		return aborted, err
	})
}

func (b *build) Schedule() (bool, error) {
	result, err := psql.Update("builds").
		Set("scheduled", true).
		Where(sq.Eq{"id": b.id}).
		RunWith(b.conn).
		Exec()
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
}

func (b *build) Pipeline() (Pipeline, bool, error) {
	if b.pipelineID == 0 {
		return nil, false, nil
	}

	row := pipelinesQuery.
		Where(sq.Eq{"p.id": b.pipelineID}).
		RunWith(b.conn).
		QueryRow()

	pipeline := newPipeline(b.conn, b.lockFactory)
	err := scanPipeline(pipeline, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return pipeline, true, nil
}

func (b *build) SaveImageResourceVersion(rc *UsedResourceCache) error {
	_, err := psql.Insert("build_image_resource_caches").
		Columns("resource_cache_id", "build_id").
		Values(rc.ID, b.id).
		RunWith(b.conn).
		Exec()
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil
		}

		return err
	}

	return nil
}

func (b *build) AcquireTrackingLock(logger lager.Logger, interval time.Duration) (lock.Lock, bool, error) {
	lock, acquired, err := b.lockFactory.Acquire(
		logger.Session("lock", lager.Data{
			"build_id": b.id,
		}),
		lock.NewBuildTrackingLockID(b.id),
	)
	if err != nil {
		return nil, false, err
	}

	if !acquired {
		return nil, false, nil
	}

	return lock, true, nil
}

func (b *build) Preparation() (BuildPreparation, bool, error) {
	if b.jobID == 0 || b.status != BuildStatusPending {
		return BuildPreparation{
			BuildID:             b.id,
			PausedPipeline:      BuildPreparationStatusNotBlocking,
			PausedJob:           BuildPreparationStatusNotBlocking,
			MaxRunningBuilds:    BuildPreparationStatusNotBlocking,
			Inputs:              map[string]BuildPreparationStatus{},
			InputsSatisfied:     BuildPreparationStatusNotBlocking,
			MissingInputReasons: MissingInputReasons{},
		}, true, nil
	}

	var (
		pausedPipeline     bool
		pausedJob          bool
		maxInFlightReached bool
		pipelineID         int
		jobName            string
	)
	err := psql.Select("p.paused, j.paused, j.max_in_flight_reached, j.pipeline_id, j.name").
		From("builds b").
		Join("jobs j ON b.job_id = j.id").
		Join("pipelines p ON j.pipeline_id = p.id").
		Where(sq.Eq{"b.id": b.id}).
		RunWith(b.conn).
		QueryRow().
		Scan(&pausedPipeline, &pausedJob, &maxInFlightReached, &pipelineID, &jobName)
	if err != nil {
		if err == sql.ErrNoRows {
			return BuildPreparation{}, false, nil
		}
		return BuildPreparation{}, false, err
	}

	pausedPipelineStatus := BuildPreparationStatusNotBlocking
	if pausedPipeline {
		pausedPipelineStatus = BuildPreparationStatusBlocking
	}

	pausedJobStatus := BuildPreparationStatusNotBlocking
	if pausedJob {
		pausedJobStatus = BuildPreparationStatusBlocking
	}

	maxInFlightReachedStatus := BuildPreparationStatusNotBlocking
	if maxInFlightReached {
		maxInFlightReachedStatus = BuildPreparationStatusBlocking
	}

	tf := NewTeamFactory(b.conn, b.lockFactory)
	t, found, err := tf.FindTeam(b.teamName)
	if err != nil {
		return BuildPreparation{}, false, err
	}

	if !found {
		return BuildPreparation{}, false, nil
	}

	pipeline, found, err := t.Pipeline(b.pipelineName)
	if err != nil {
		return BuildPreparation{}, false, err
	}

	if !found {
		return BuildPreparation{}, false, nil
	}

	job, found, err := pipeline.Job(jobName)
	if err != nil {
		return BuildPreparation{}, false, err
	}

	if !found {
		return BuildPreparation{}, false, nil
	}

	configInputs := job.Config().Inputs()

	nextBuildInputs, found, err := job.GetNextBuildInputs()
	if err != nil {
		return BuildPreparation{}, false, err
	}

	inputsSatisfiedStatus := BuildPreparationStatusBlocking
	inputs := map[string]BuildPreparationStatus{}
	missingInputReasons := MissingInputReasons{}

	if found {
		inputsSatisfiedStatus = BuildPreparationStatusNotBlocking
		for _, buildInput := range nextBuildInputs {
			inputs[buildInput.Name] = BuildPreparationStatusNotBlocking
		}
	} else {
		buildInputs, err := job.GetIndependentBuildInputs()
		if err != nil {
			return BuildPreparation{}, false, err
		}

		for _, configInput := range configInputs {
			found := false
			for _, buildInput := range buildInputs {
				if buildInput.Name == configInput.Name {
					found = true
					break
				}
			}
			if found {
				inputs[configInput.Name] = BuildPreparationStatusNotBlocking
			} else {
				inputs[configInput.Name] = BuildPreparationStatusBlocking
				if len(configInput.Passed) > 0 {
					if configInput.Version != nil && configInput.Version.Pinned != nil {
						_, found, err := pipeline.GetVersionedResourceByVersion(configInput.Version.Pinned, configInput.Resource)
						if err != nil {
							return BuildPreparation{}, false, err
						}

						if found {
							missingInputReasons.RegisterPassedConstraint(configInput.Name)
						} else {
							versionJSON, err := json.Marshal(configInput.Version.Pinned)
							if err != nil {
								return BuildPreparation{}, false, err
							}

							missingInputReasons.RegisterPinnedVersionUnavailable(configInput.Name, string(versionJSON))
						}
					} else {
						missingInputReasons.RegisterPassedConstraint(configInput.Name)
					}
				} else {
					if configInput.Version != nil && configInput.Version.Pinned != nil {
						versionJSON, err := json.Marshal(configInput.Version.Pinned)
						if err != nil {
							return BuildPreparation{}, false, err
						}

						missingInputReasons.RegisterPinnedVersionUnavailable(configInput.Name, string(versionJSON))
					} else {
						missingInputReasons.RegisterNoVersions(configInput.Name)
					}
				}
			}
		}
	}

	buildPreparation := BuildPreparation{
		BuildID:             b.id,
		PausedPipeline:      pausedPipelineStatus,
		PausedJob:           pausedJobStatus,
		MaxRunningBuilds:    maxInFlightReachedStatus,
		Inputs:              inputs,
		InputsSatisfied:     inputsSatisfiedStatus,
		MissingInputReasons: missingInputReasons,
	}

	return buildPreparation, true, nil
}

func (b *build) Events(from uint) (EventSource, error) {
	notifier, err := newConditionNotifier(b.conn.Bus(), buildEventsChannel(b.id), func() (bool, error) {
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	table := fmt.Sprintf("team_build_events_%d", b.teamID)
	if b.pipelineID != 0 {
		table = fmt.Sprintf("pipeline_build_events_%d", b.pipelineID)
	}

	return newBuildEventSource(
		b.id,
		table,
		b.conn,
		notifier,
		from,
	), nil
}

func (b *build) SaveEvent(event atc.Event) error {
	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = b.saveEvent(tx, event)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = b.conn.Bus().Notify(buildEventsChannel(b.id))
	if err != nil {
		return err
	}

	return nil
}

func (b *build) SaveInput(input BuildInput) error {
	if b.pipelineID == 0 {
		return nil
	}

	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	row := pipelinesQuery.
		Where(sq.Eq{"p.id": b.pipelineID}).
		RunWith(tx).
		QueryRow()

	pipeline := &pipeline{conn: b.conn, lockFactory: b.lockFactory}
	err = scanPipeline(pipeline, row)
	if err != nil {
		return err
	}

	err = pipeline.saveInputTx(tx, b.id, input)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (b *build) SaveOutput(vr VersionedResource, explicit bool) error {
	if b.pipelineID == 0 {
		return nil
	}

	row := pipelinesQuery.
		Where(sq.Eq{"p.id": b.pipelineID}).
		RunWith(b.conn).
		QueryRow()
	pipeline := &pipeline{conn: b.conn, lockFactory: b.lockFactory}
	err := scanPipeline(pipeline, row)
	if err != nil {
		return err
	}

	return pipeline.saveOutput(b.id, vr, explicit)
}

func (b *build) UseInputs(inputs []BuildInput) error {
	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = psql.Delete("build_inputs").
		Where(sq.Eq{"build_id": b.id}).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	row := pipelinesQuery.
		Where(sq.Eq{"p.id": b.pipelineID}).
		RunWith(tx).
		QueryRow()

	pipeline := &pipeline{conn: b.conn, lockFactory: b.lockFactory}
	err = scanPipeline(pipeline, row)
	if err != nil {
		return err
	}

	for _, input := range inputs {
		err = pipeline.saveInputTx(tx, b.id, input)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (b *build) Resources() ([]BuildInput, []BuildOutput, error) {
	inputs := []BuildInput{}
	outputs := []BuildOutput{}

	rows, err := b.conn.Query(`
		SELECT i.name, r.name, vr.type, vr.version, vr.metadata,
		NOT EXISTS (
			SELECT 1
			FROM build_inputs ci, builds cb
			WHERE versioned_resource_id = vr.id
			AND cb.job_id = b.job_id
			AND ci.build_id = cb.id
			AND ci.build_id < b.id
		)
		FROM versioned_resources vr, build_inputs i, builds b, resource_spaces rs, resources r
		WHERE b.id = $1
		AND i.build_id = b.id
		AND i.versioned_resource_id = vr.id
		AND rs.id = vr.resource_space_id
    AND r.id = rs.resource_id
		AND NOT EXISTS (
			SELECT 1
			FROM build_outputs o
			WHERE o.versioned_resource_id = vr.id
			AND o.build_id = i.build_id
			AND o.explicit
		)
	`, b.id)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var inputName string
		var vr VersionedResource
		var firstOccurrence bool

		var version, metadata string
		err := rows.Scan(&inputName, &vr.Resource, &vr.Type, &version, &metadata, &firstOccurrence)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal([]byte(version), &vr.Version)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal([]byte(metadata), &vr.Metadata)
		if err != nil {
			return nil, nil, err
		}

		inputs = append(inputs, BuildInput{
			Name:              inputName,
			VersionedResource: vr,
			FirstOccurrence:   firstOccurrence,
		})
	}

	rows, err = b.conn.Query(`
		SELECT r.name, vr.type, vr.version, vr.metadata
		FROM versioned_resources vr, build_outputs o, builds b, resource_spaces rs, resources r
		WHERE b.id = $1
		AND o.build_id = b.id
		AND o.versioned_resource_id = vr.id
		AND rs.id = vr.resource_space_id
    AND r.id = rs.resource_id
		AND o.explicit
	`, b.id)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vr VersionedResource

		var version, metadata string
		err := rows.Scan(&vr.Resource, &vr.Type, &version, &metadata)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal([]byte(version), &vr.Version)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal([]byte(metadata), &vr.Metadata)
		if err != nil {
			return nil, nil, err
		}

		outputs = append(outputs, BuildOutput{
			VersionedResource: vr,
		})
	}

	return inputs, outputs, nil
}

func (b *build) GetVersionedResources() (SavedVersionedResources, error) {
	return b.getVersionedResources(`
		SELECT vr.id,
			vr.enabled,
			vr.version,
			vr.metadata,
			vr.type,
			r.name,
			vr.modified_time
		FROM builds b
		INNER JOIN jobs j ON j.id = b.job_id
		INNER JOIN build_inputs bi ON bi.build_id = b.id
		INNER JOIN versioned_resources vr ON vr.id = bi.versioned_resource_id
		INNER JOIN resource_spaces rs ON rs.id = vr.resource_space_id
		INNER JOIN resources r ON r.id = rs.resource_id
		WHERE b.id = $1

		UNION ALL

		SELECT vr.id,
			vr.enabled,
			vr.version,
			vr.metadata,
			vr.type,
			r.name,
			vr.modified_time
		FROM builds b
		INNER JOIN jobs j ON j.id = b.job_id
		INNER JOIN build_outputs bo ON bo.build_id = b.id
		INNER JOIN versioned_resources vr ON vr.id = bo.versioned_resource_id
		INNER JOIN resource_spaces rs ON rs.id = vr.resource_space_id
		INNER JOIN resources r ON r.id = rs.resource_id
		WHERE b.id = $1 AND bo.explicit`)
}

func (b *build) getVersionedResources(resourceRequest string) (SavedVersionedResources, error) {
	rows, err := b.conn.Query(resourceRequest, b.id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	savedVersionedResources := SavedVersionedResources{}

	for rows.Next() {
		var versionedResource SavedVersionedResource
		var versionJSON []byte
		var metadataJSON []byte
		err = rows.Scan(&versionedResource.ID, &versionedResource.Enabled, &versionJSON, &metadataJSON, &versionedResource.Type, &versionedResource.Resource, &versionedResource.ModifiedTime)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(versionJSON, &versionedResource.Version)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(metadataJSON, &versionedResource.Metadata)
		if err != nil {
			return nil, err
		}

		savedVersionedResources = append(savedVersionedResources, versionedResource)
	}

	return savedVersionedResources, nil
}

func createBuildEventSeq(tx Tx, buildid int) error {
	_, err := tx.Exec(fmt.Sprintf(`
		CREATE SEQUENCE %s MINVALUE 0
	`, buildEventSeq(buildid)))
	return err
}

func buildEventSeq(buildid int) string {
	return fmt.Sprintf("build_event_id_seq_%d", buildid)
}

func scanBuild(b *build, row scannable, encryptionStrategy EncryptionStrategy) error {
	var (
		jobID, pipelineID                                         sql.NullInt64
		engine, engineMetadata, jobName, pipelineName, publicPlan sql.NullString
		startTime, endTime, reapTime                              pq.NullTime
		nonce                                                     sql.NullString

		status string
	)

	err := row.Scan(&b.id, &b.name, &jobID, &b.teamID, &status, &b.isManuallyTriggered, &b.scheduled, &engine, &engineMetadata, &publicPlan, &startTime, &endTime, &reapTime, &jobName, &pipelineID, &pipelineName, &b.teamName, &nonce)
	if err != nil {
		return err
	}

	b.status = BuildStatus(status)
	b.jobName = jobName.String
	b.jobID = int(jobID.Int64)
	b.pipelineName = pipelineName.String
	b.pipelineID = int(pipelineID.Int64)
	b.engine = engine.String
	b.startTime = startTime.Time
	b.endTime = endTime.Time
	b.reapTime = reapTime.Time

	var noncense *string
	if nonce.Valid {
		noncense = &nonce.String
		decryptedEngineMetadata, err := encryptionStrategy.Decrypt(string(engineMetadata.String), noncense)
		if err != nil {
			return err
		}
		b.engineMetadata = string(decryptedEngineMetadata)
	} else {
		b.engineMetadata = engineMetadata.String
	}

	if publicPlan.Valid {
		err = json.Unmarshal([]byte(publicPlan.String), &b.publicPlan)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *build) saveEvent(tx Tx, event atc.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	table := fmt.Sprintf("team_build_events_%d", b.teamID)
	if b.pipelineID != 0 {
		table = fmt.Sprintf("pipeline_build_events_%d", b.pipelineID)
	}
	_, err = psql.Insert(table).
		Columns("event_id", "build_id", "type", "version", "payload").
		Values(sq.Expr("nextval('"+buildEventSeq(b.id)+"')"), b.id, string(event.EventType()), string(event.Version()), payload).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func createBuild(tx Tx, build *build, vals map[string]interface{}) error {
	var buildID int
	err := psql.Insert("builds").
		SetMap(vals).
		Suffix("RETURNING id").
		RunWith(tx).
		QueryRow().
		Scan(&buildID)
	if err != nil {
		return err
	}

	err = scanBuild(build, buildsQuery.
		Where(sq.Eq{"b.id": buildID}).
		RunWith(tx).
		QueryRow(),
		build.conn.EncryptionStrategy(),
	)
	if err != nil {
		return err
	}

	err = createBuildEventSeq(tx, buildID)
	if err != nil {
		return err
	}

	return nil
}

func buildEventsChannel(buildID int) string {
	return fmt.Sprintf("build_events_%d", buildID)
}

func buildAbortChannel(buildID int) string {
	return fmt.Sprintf("build_abort_%d", buildID)
}
