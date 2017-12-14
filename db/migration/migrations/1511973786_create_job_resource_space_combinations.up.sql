BEGIN;

  CREATE TABLE job_combinations (
    id SERIAL PRIMARY KEY,
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    combination jsonb NOT NULL,
    inputs_determined boolean DEFAULT false
  );

  CREATE TABLE job_combinations_resource_spaces (
    job_combination_id int REFERENCES job_combinations (id) ON DELETE CASCADE,
    resource_space_id int REFERENCES resource_spaces (id) ON DELETE CASCADE
  );

COMMIT;
