BEGIN;

  CREATE TABLE job_combinations (
    id SERIAL PRIMARY KEY,
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    combination jsonb,
    inputs_determined boolean DEFAULT false
  );

  CREATE UNIQUE INDEX job_combinations_job_id_combination_key ON job_combinations (job_id, combination);

  ALTER TABLE builds RENAME job_id TO job_combination_id;
  ALTER TABLE builds DROP CONSTRAINT fkey_job_id;
  ALTER TABLE builds ADD CONSTRAINT fkey_job_combination_id FOREIGN KEY (job_combination_id) REFERENCES job_combinations(id) ON DELETE CASCADE;
  ALTER INDEX builds_job_id RENAME TO builds_job_combination_id;

COMMIT;
