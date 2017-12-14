BEGIN;

  /* ALTER TABLE builds ADD COLUMN job_combination_id varchar(20); */
  /* ALTER TABLE builds ADD CONSTRAINT job_combination_id_fkey FOREIGN KEY (job_combination_id) REFERENCES job_resource_space_combinations (hash) ON DELETE CASCADE; */
  /* CREATE INDEX builds_job_id_ ON builds USING btree (interceptible, completed); */

  /* ALTER TABLE next_build_inputs ADD COLUMN job_combination_id varchar(20); */
  /* ALTER TABLE independent_build_inputs ADD COLUMN job_combination_id varchar(20) */
  /* ALTER TABLE worker_task_caches ADD COLUMN job_combination_id varchar(20) */

COMMIT;
