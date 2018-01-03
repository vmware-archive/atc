BEGIN;

  CREATE TABLE job_combinations (
    id SERIAL PRIMARY KEY,
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    combination jsonb,
    inputs_determined boolean DEFAULT false
  );

  INSERT INTO job_combinations(id, job_id, combination)
  SELECT id, id, json_object_agg(resource_name, 'default')
  FROM jobs, LATERAL (SELECT DISTINCT(json_array_elements(config::json->'plan')->>'get') AS resource_name
                      UNION SELECT DISTINCT(json_array_elements(config::json->'plan')->>'put')) _
  WHERE resource_name is NOT NULL
  GROUP BY id;

  SELECT setval('job_combinations_id_seq', (SELECT max(id) FROM job_combinations));

  /* ALTER TABLE builds ADD COLUMN job_combination_id varchar(20); */
  /* ALTER TABLE builds ADD CONSTRAINT job_combination_id_fkey FOREIGN KEY (job_combination_id) REFERENCES job_resource_space_combinations (hash) ON DELETE CASCADE; */
  /* CREATE INDEX builds_job_id_ ON builds USING btree (interceptible, completed); */

  /* ALTER TABLE next_build_inputs ADD COLUMN job_combination_id varchar(20); */
  /* ALTER TABLE independent_build_inputs ADD COLUMN job_combination_id varchar(20) */
  /* ALTER TABLE worker_task_caches ADD COLUMN job_combination_id varchar(20) */

COMMIT;
