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

COMMIT;
