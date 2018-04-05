BEGIN;

  CREATE TABLE job_combinations (
    id SERIAL PRIMARY KEY,
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    combination jsonb NOT NULL,
    inputs_determined boolean DEFAULT false
  );

  CREATE UNIQUE INDEX job_combinations_job_id_combination_key ON job_combinations (job_id, combination);

COMMIT;
