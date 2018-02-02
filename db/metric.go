package db

import (
	"github.com/concourse/atc/metric"
)

var slowQuery := dbConn.Query(
	`total_time / calls as avg_time,
	calls,
	total_time,
	rows,
	100.0 * shared_blks_hit / nullif(shared_blks_hit + shared_blks_read, 0) AS hit_percent,
  regexp_replace(query, '[\s\t\n]+', ' ', 'g') as sanitized_sql
	FROM pg_stat_statements
	WHERE query NOT LIKE '%EXPLAIN%'
	AND query NOT LIKE '%INDEX%'
	AND query NOT LIKE '%pg_stat_statements%'
	AND calls > 1
	ORDER BY avg_time DESC LIMIT 5
`)

func FetchSlowQueries() ([]SlowQuery error) {

	return nil, nil
}
