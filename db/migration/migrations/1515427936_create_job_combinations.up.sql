BEGIN;

  CREATE TABLE job_combinations (
    id SERIAL PRIMARY KEY,
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    combination jsonb,
    inputs_determined boolean DEFAULT false
  );
  CREATE UNIQUE INDEX job_combinations_job_id_combination_key ON job_combinations (job_id, combination);

  INSERT INTO job_combinations(id, job_id)
  SELECT id, id
  FROM jobs;

  ALTER TABLE builds RENAME job_id TO job_combination_id;
  ALTER TABLE builds DROP CONSTRAINT fkey_job_id;
  ALTER TABLE builds ADD CONSTRAINT fkey_job_combination_id FOREIGN KEY (job_combination_id) REFERENCES job_combinations(id) ON DELETE CASCADE;
  ALTER INDEX builds_job_id RENAME TO builds_job_combination_id;

  DROP INDEX next_builds_per_job_id;
  DROP INDEX latest_completed_builds_per_job_id;
  DROP INDEX transition_builds_per_job_id;
  DROP MATERIALIZED VIEW next_builds_per_job;
  DROP MATERIALIZED VIEW transition_builds_per_job;
  DROP MATERIALIZED VIEW latest_completed_builds_per_job;

  CREATE MATERIALIZED VIEW next_builds_per_job_combination AS
   WITH latest_build_ids_per_job_combination AS (
           SELECT min(b_1.id) AS build_id
             FROM (builds b_1
               JOIN job_combinations c ON ((c.id = b_1.job_combination_id)))
            WHERE (b_1.status = ANY (ARRAY['pending'::build_status, 'started'::build_status]))
            GROUP BY b_1.job_combination_id
          )
   SELECT b.id,
      b.name,
      b.status,
      b.scheduled,
      b.start_time,
      b.end_time,
      b.engine,
      b.engine_metadata,
      b.completed,
      b.job_combination_id,
      b.reap_time,
      b.team_id,
      b.manually_triggered,
      b.interceptible,
      b.nonce,
      b.public_plan,
      b.pipeline_id
     FROM (builds b
       JOIN latest_build_ids_per_job_combination l ON ((l.build_id = b.id)))
    WITH NO DATA;
  REFRESH MATERIALIZED VIEW next_builds_per_job_combination;

  CREATE MATERIALIZED VIEW latest_completed_builds_per_job_combination AS
   WITH latest_build_ids_per_job_combination AS (
           SELECT max(b_1.id) AS build_id
             FROM (builds b_1
               JOIN job_combinations c ON ((c.id = b_1.job_combination_id)))
            WHERE (b_1.status <> ALL (ARRAY['pending'::build_status, 'started'::build_status]))
            GROUP BY b_1.job_combination_id
          )
   SELECT b.id,
      b.name,
      b.status,
      b.scheduled,
      b.start_time,
      b.end_time,
      b.engine,
      b.engine_metadata,
      b.completed,
      b.job_combination_id,
      b.reap_time,
      b.team_id,
      b.manually_triggered,
      b.interceptible,
      b.nonce,
      b.public_plan,
      b.pipeline_id
     FROM (builds b
       JOIN latest_build_ids_per_job_combination l ON ((l.build_id = b.id)))
    WITH NO DATA;
  REFRESH MATERIALIZED VIEW latest_completed_builds_per_job_combination;

  CREATE MATERIALIZED VIEW transition_builds_per_job_combination AS
   WITH builds_before_transition AS (
           SELECT b_1.job_combination_id,
              max(b_1.id) AS max
             FROM ((builds b_1
               LEFT JOIN job_combinations c ON ((c.id = b_1.job_combination_id)))
               LEFT JOIN latest_completed_builds_per_job_combination s ON ((b_1.job_combination_id = s.job_combination_id)))
            WHERE ((b_1.status <> s.status) AND (b_1.status <> ALL (ARRAY['pending'::build_status, 'started'::build_status])))
            GROUP BY b_1.job_combination_id
          )
   SELECT DISTINCT ON (b.job_combination_id) b.id,
      b.name,
      b.status,
      b.scheduled,
      b.start_time,
      b.end_time,
      b.engine,
      b.engine_metadata,
      b.completed,
      b.job_combination_id,
      b.reap_time,
      b.team_id,
      b.manually_triggered,
      b.interceptible,
      b.nonce,
      b.public_plan,
      b.pipeline_id
     FROM (builds b
       LEFT JOIN builds_before_transition ON ((b.job_combination_id = builds_before_transition.job_combination_id)))
    WHERE (((builds_before_transition.max IS NULL) AND (b.status <> ALL (ARRAY['pending'::build_status, 'started'::build_status]))) OR (b.id > builds_before_transition.max))
    ORDER BY b.job_combination_id, b.id
    WITH NO DATA;
  REFRESH MATERIALIZED VIEW transition_builds_per_job_combination;

  CREATE UNIQUE INDEX next_builds_per_job_combination_id ON next_builds_per_job_combination USING btree (id);
  CREATE UNIQUE INDEX latest_completed_builds_per_job_combination_id ON latest_completed_builds_per_job_combination USING btree (id);
  CREATE UNIQUE INDEX transition_builds_per_job_combination_id ON transition_builds_per_job_combination USING btree (id);

COMMIT;
