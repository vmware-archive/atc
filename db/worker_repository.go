package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
	"github.com/concourse/atc/service"
)

// PostgresWorkerRepository implements service.WorkerRepository for a Postgres back-end
type PostgresWorkerRepository struct {
	conn Conn
}

// NewWorkerRepository returns a worker repository located at the given database connection
func NewWorkerRepository(conn Conn) service.WorkerRepository {
	return &PostgresWorkerRepository{
		conn: conn,
	}
}

var workersQuery = psql.Select(`
		w.name,
		w.version,
		w.addr,
		w.state,
		w.baggageclaim_url,
		w.certs_path,
		w.http_proxy_url,
		w.https_proxy_url,
		w.no_proxy,
		w.active_containers,
		w.resource_types,
		w.platform,
		w.tags,
		t.name,
		w.team_id,
		w.start_time,
		w.expires
	`).
	From("workers w").
	LeftJoin("teams t ON w.team_id = t.id")

// GetWorker will fetch a atc.Worker from our WorkerRepository
// Pass in name of the worker for the query
func (workerRepo *PostgresWorkerRepository) GetWorker(name string) (*atc.Worker, bool, error) {
	return getWorker(workerRepo.conn, workersQuery.Where(sq.Eq{"w.name": name}))
}

// VisibleWorkers returns a list of workers that are visible to the given team teamNames
// Returned workers are either specific to one of the given team names or they already
// workers accessible to all teams.
func (workerRepo *PostgresWorkerRepository) VisibleWorkers(teamNames []string) ([]atc.Worker, error) {
	workersQuery := workersQuery.
		Where(sq.Or{
			sq.Eq{"t.name": teamNames},
			sq.Eq{"w.team_id": nil},
		})

	workers, err := getWorkers(workerRepo.conn, workersQuery)
	if err != nil {
		return nil, err
	}

	return workers, nil
}

// Workers returns all available workers on in this Concourse instance
func (workerRepo *PostgresWorkerRepository) Workers() ([]atc.Worker, error) {
	return getWorkers(workerRepo.conn, workersQuery)
}

// Reload will reload the stored information for a worker from the repository
func (workerRepo *PostgresWorkerRepository) Reload(worker *atc.Worker) (bool, error) {
	row := workersQuery.Where(sq.Eq{"w.name": worker.Name}).
		RunWith(workerRepo.conn).
		QueryRow()

	err := scanWorker(worker, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getWorker(conn Conn, query sq.SelectBuilder) (*atc.Worker, bool, error) {
	row := query.
		RunWith(conn).
		QueryRow()

	worker := &atc.Worker{}
	err := scanWorker(worker, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return worker, true, nil
}

func getWorkers(conn Conn, query sq.SelectBuilder) ([]atc.Worker, error) {
	rows, err := query.RunWith(conn).Query()
	if err != nil {
		return nil, err
	}
	defer Close(rows)

	workers := []atc.Worker{}

	for rows.Next() {
		worker := &atc.Worker{}
		err := scanWorker(worker, rows)
		if err != nil {
			return nil, err
		}

		workers = append(workers, *worker)
	}

	return workers, nil
}

func scanWorker(worker *atc.Worker, row scannable) error {
	var (
		version  sql.NullString
		addStr   sql.NullString
		state    string
		bcURLStr sql.NullString
		//	reaperAddr    sql.NullString
		certsPathStr  sql.NullString
		httpProxyURL  sql.NullString
		httpsProxyURL sql.NullString
		noProxy       sql.NullString
		resourceTypes []byte
		platform      sql.NullString
		tags          []byte
		teamName      sql.NullString
		teamID        sql.NullInt64
		startTime     sql.NullInt64
		expiresAt     *time.Time
	)

	err := row.Scan(
		&worker.Name,
		&version,
		&addStr,
		&state,
		&bcURLStr,
		//	&reaperAddr,
		&certsPathStr,
		&httpProxyURL,
		&httpsProxyURL,
		&noProxy,
		&worker.ActiveContainers,
		&resourceTypes,
		&platform,
		&tags,
		&teamName,
		&teamID,
		&startTime,
		&expiresAt,
	)
	if err != nil {
		return err
	}

	if version.Valid {
		worker.Version = version.String
	}

	if addStr.Valid {
		worker.GardenAddr = addStr.String
	}

	if bcURLStr.Valid {
		worker.BaggageclaimURL = bcURLStr.String
	}

	// if reaperAddr.Valid {
	// 	worker.reaperAddr = &reaperAddr.String
	// }

	if certsPathStr.Valid {
		worker.CertsPath = &certsPathStr.String
	}

	worker.State = string(WorkerState(state))

	if startTime.Valid {
		worker.StartTime = startTime.Int64
	}

	if expiresAt != nil {
		worker.ExpiresAt = *expiresAt
	}

	if httpProxyURL.Valid {
		worker.HTTPProxyURL = httpProxyURL.String
	}

	if httpsProxyURL.Valid {
		worker.HTTPSProxyURL = httpsProxyURL.String
	}

	if noProxy.Valid {
		worker.NoProxy = noProxy.String
	}

	if teamName.Valid {
		worker.Team = teamName.String
	}

	if teamID.Valid {
		worker.TeamID = int(teamID.Int64)
	}

	if platform.Valid {
		worker.Platform = platform.String
	}

	err = json.Unmarshal(resourceTypes, &worker.ResourceTypes)
	if err != nil {
		return err
	}

	return json.Unmarshal(tags, &worker.Tags)
}

// HeartbeatWorker is used to update the TTL for the worker for the atc.Worker on the persistent store
// It allows Concourse to recognize when workers become unresponsive and should be put in 'stalled' state
func (workerRepo *PostgresWorkerRepository) HeartbeatWorker(atcWorker atc.Worker, ttl time.Duration) (*atc.Worker, error) {
	// In order to be able to calculate the ttl that we return to the caller
	// we must compare time.Now() to the worker.expires column
	// However, workers.expires column is a "timestamp (without timezone)"
	// So we format time.Now() without any timezone information and then
	// parse that using the same layout to strip the timezone information

	tx, err := workerRepo.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer Rollback(tx)

	expires := "NULL"
	if ttl != 0 {
		expires = fmt.Sprintf(`NOW() + '%d second'::INTERVAL`, int(ttl.Seconds()))
	}

	cSQL, _, err := sq.Case("state").
		When("'landing'::worker_state", "'landing'::worker_state").
		When("'landed'::worker_state", "'landed'::worker_state").
		When("'retiring'::worker_state", "'retiring'::worker_state").
		Else("'running'::worker_state").
		ToSql()

	if err != nil {
		return nil, err
	}

	addrSQL, _, err := sq.Case("state").
		When("'landed'::worker_state", "NULL").
		Else("'" + atcWorker.GardenAddr + "'").
		ToSql()
	if err != nil {
		return nil, err
	}

	bcSQL, _, err := sq.Case("state").
		When("'landed'::worker_state", "NULL").
		Else("'" + atcWorker.BaggageclaimURL + "'").
		ToSql()
	if err != nil {
		return nil, err
	}

	// reapSQL, _, err := sq.Case("state").
	// 	When("'landed'::worker_state", "NULL").
	// 	Else("'" + atcWorker.ReaperAddr + "'").
	// 	ToSql()
	// if err != nil {
	// 	return nil, err
	// }

	_, err = psql.Update("workers").
		Set("expires", sq.Expr(expires)).
		Set("addr", sq.Expr("("+addrSQL+")")).
		Set("baggageclaim_url", sq.Expr("("+bcSQL+")")).
		//	Set("reaper_addr", sq.Expr("("+reapSQL+")")).
		Set("active_containers", atcWorker.ActiveContainers).
		Set("state", sq.Expr("("+cSQL+")")).
		Where(sq.Eq{"name": atcWorker.Name}).
		RunWith(tx).
		Exec()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWorkerNotPresent
		}
		return nil, err
	}

	row := workersQuery.Where(sq.Eq{"w.name": atcWorker.Name}).
		RunWith(tx).
		QueryRow()

	worker := &atc.Worker{}
	err = scanWorker(worker, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWorkerNotPresent
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return worker, nil

}

// SaveWorker persists a atc.Worker to the data store
func (workerRepo *PostgresWorkerRepository) SaveWorker(atcWorker atc.Worker, ttl time.Duration) (*atc.Worker, error) {
	tx, err := workerRepo.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer Rollback(tx)

	savedWorker, err := saveWorker(tx, atcWorker, nil, ttl, workerRepo.conn)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return savedWorker, nil
}

// Delete a worker from the persistent store
func (workerRepo *PostgresWorkerRepository) Delete(worker *atc.Worker) error {
	_, err := sq.Delete("workers").
		Where(sq.Eq{
			"name": worker.Name,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(workerRepo.conn).
		Exec()

	return err
}

func saveWorker(tx Tx, atcWorker atc.Worker, teamID *int, ttl time.Duration, conn Conn) (*atc.Worker, error) {
	resourceTypes, err := json.Marshal(atcWorker.ResourceTypes)
	if err != nil {
		return nil, err
	}

	tags, err := json.Marshal(atcWorker.Tags)
	if err != nil {
		return nil, err
	}

	expires := "NULL"
	if ttl != 0 {
		expires = fmt.Sprintf(`NOW() + '%d second'::INTERVAL`, int(ttl.Seconds()))
	}

	var oldTeamID sql.NullInt64

	var workerState WorkerState
	if atcWorker.State != "" {
		workerState = WorkerState(atcWorker.State)
	} else {
		workerState = WorkerStateRunning
	}

	var workerVersion *string
	if atcWorker.Version != "" {
		workerVersion = &atcWorker.Version
	}

	err = psql.Select("team_id").From("workers").Where(sq.Eq{
		"name": atcWorker.Name,
	}).RunWith(tx).QueryRow().Scan(&oldTeamID)

	if err != nil {
		if err == sql.ErrNoRows {
			_, err = psql.Insert("workers").
				Columns(
					"addr",
					"expires",
					"active_containers",
					"resource_types",
					"tags",
					"platform",
					"baggageclaim_url",
					// "reaper_addr",
					"certs_path",
					"http_proxy_url",
					"https_proxy_url",
					"no_proxy",
					"name",
					"version",
					"start_time",
					"team_id",
					"state",
				).
				Values(
					atcWorker.GardenAddr,
					sq.Expr(expires),
					atcWorker.ActiveContainers,
					resourceTypes,
					tags,
					atcWorker.Platform,
					atcWorker.BaggageclaimURL,
					//			atcWorker.ReaperAddr,
					atcWorker.CertsPath,
					atcWorker.HTTPProxyURL,
					atcWorker.HTTPSProxyURL,
					atcWorker.NoProxy,
					atcWorker.Name,
					workerVersion,
					atcWorker.StartTime,
					teamID,
					string(workerState),
				).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if (oldTeamID.Valid == (teamID == nil)) ||
			(oldTeamID.Valid && (*teamID != int(oldTeamID.Int64))) {
			return nil, errors.New("update-of-other-teams-worker-not-allowed")
		}

		_, err = psql.Update("workers").
			Set("addr", atcWorker.GardenAddr).
			Set("expires", sq.Expr(expires)).
			Set("active_containers", atcWorker.ActiveContainers).
			Set("resource_types", resourceTypes).
			Set("tags", tags).
			Set("platform", atcWorker.Platform).
			Set("baggageclaim_url", atcWorker.BaggageclaimURL).
			//		Set("reaper_addr", atcWorker.ReaperAddr).
			Set("certs_path", atcWorker.CertsPath).
			Set("http_proxy_url", atcWorker.HTTPProxyURL).
			Set("https_proxy_url", atcWorker.HTTPSProxyURL).
			Set("no_proxy", atcWorker.NoProxy).
			Set("name", atcWorker.Name).
			Set("version", workerVersion).
			Set("start_time", atcWorker.StartTime).
			Set("state", string(workerState)).
			Where(sq.Eq{
				"name": atcWorker.Name,
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
	}

	var workerTeamID int
	if teamID != nil {
		workerTeamID = *teamID
	}

	savedWorker := &atc.Worker{
		Name:             atcWorker.Name,
		Version:          *workerVersion,
		State:            string(workerState),
		GardenAddr:       atcWorker.GardenAddr,
		BaggageclaimURL:  atcWorker.BaggageclaimURL,
		CertsPath:        atcWorker.CertsPath,
		HTTPProxyURL:     atcWorker.HTTPProxyURL,
		HTTPSProxyURL:    atcWorker.HTTPSProxyURL,
		NoProxy:          atcWorker.NoProxy,
		ActiveContainers: atcWorker.ActiveContainers,
		ResourceTypes:    atcWorker.ResourceTypes,
		Platform:         atcWorker.Platform,
		Tags:             atcWorker.Tags,
		Team:             atcWorker.Team,
		TeamID:           workerTeamID,
		StartTime:        atcWorker.StartTime,
	}

	workerBaseResourceTypeIDs := []int{}

	var (
		brt  BaseResourceType
		ubrt *UsedBaseResourceType
		uwrt *UsedWorkerResourceType
	)

	for _, resourceType := range atcWorker.ResourceTypes {
		workerResourceType := WorkerResourceType{
			Worker:  savedWorker,
			Image:   resourceType.Image,
			Version: resourceType.Version,
			BaseResourceType: &BaseResourceType{
				Name: resourceType.Type,
			},
		}

		brt = BaseResourceType{
			Name: resourceType.Type,
		}

		ubrt, err = brt.FindOrCreate(tx)
		if err != nil {
			return nil, err
		}

		_, err = psql.Delete("worker_base_resource_types").
			Where(sq.Eq{
				"worker_name":           atcWorker.Name,
				"base_resource_type_id": ubrt.ID,
			}).
			Where(sq.NotEq{
				"version": resourceType.Version,
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
		uwrt, err = workerResourceType.FindOrCreate(tx)
		if err != nil {
			return nil, err
		}

		workerBaseResourceTypeIDs = append(workerBaseResourceTypeIDs, uwrt.ID)
	}

	_, err = psql.Delete("worker_base_resource_types").
		Where(sq.Eq{
			"worker_name": atcWorker.Name,
		}).
		Where(sq.NotEq{
			"id": workerBaseResourceTypeIDs,
		}).
		RunWith(tx).
		Exec()
	if err != nil {
		return nil, err
	}

	if atcWorker.CertsPath != nil {
		_, err := WorkerResourceCerts{
			WorkerName: atcWorker.Name,
			CertsPath:  *atcWorker.CertsPath,
		}.FindOrCreate(tx)

		if err != nil {
			return nil, err
		}
	}

	return savedWorker, nil
}

// ResourceCerts returns the resource certs for the worker
func (workerRepo *PostgresWorkerRepository) ResourceCerts(worker *atc.Worker) (*UsedWorkerResourceCerts, bool, error) {
	if worker.CertsPath != nil {
		wrc := &WorkerResourceCerts{
			WorkerName: worker.Name,
			CertsPath:  *worker.CertsPath,
		}

		return wrc.Find(workerRepo.conn)
	}

	return nil, false, nil
}
