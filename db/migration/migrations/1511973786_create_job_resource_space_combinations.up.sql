BEGIN;

  CREATE TABLE job_resource_space_combinations (
    job_id int REFERENCES jobs (id) ON DELETE CASCADE,
    resource_space_id int REFERENCES resource_spaces (id) ON DELETE CASCADE,
    combination jsonb NOT NULL,
    hash varchar(20) NOT NULL,
    inputs_determined boolean DEFAULT false,
    UNIQUE (job_id, resource_space_id, hash)
  );

COMMIT;
