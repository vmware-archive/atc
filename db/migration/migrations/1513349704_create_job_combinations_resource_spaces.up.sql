BEGIN;

  CREATE TABLE job_combinations_resource_spaces (
    job_combination_id int REFERENCES job_combinations (id) ON DELETE CASCADE,
    resource_space_id int REFERENCES resource_spaces (id) ON DELETE CASCADE
  );

  INSERT INTO job_combinations_resource_spaces(job_combination_id, resource_space_id)
  SELECT job_resources.id as job_combination_id, resources.id as resource_space_id FROM (
    SELECT id, (json_each_text(json_array_elements((config::json->'plan')))).value FROM jobs
  ) job_resources JOIN resources ON resources.name = job_resources.value;

COMMIT;
