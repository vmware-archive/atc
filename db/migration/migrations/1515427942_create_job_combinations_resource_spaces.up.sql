BEGIN;

  CREATE TABLE job_combinations_resource_spaces (
    job_combination_id int REFERENCES job_combinations (id) ON DELETE CASCADE,
    resource_space_id int REFERENCES resource_spaces (id) ON DELETE CASCADE
  );

  CREATE UNIQUE INDEX job_combinations_resource_spaces_job_combination_id_resource_space_id_key ON job_combinations_resource_spaces (job_combination_id, resource_space_id);

COMMIT;
