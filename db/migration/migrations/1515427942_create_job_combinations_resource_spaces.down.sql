BEGIN;

  DROP INDEX job_combinations_resource_spaces_job_combination_id_resource_space_id_key;

  DROP TABLE job_combinations_resource_spaces;

COMMIT;
