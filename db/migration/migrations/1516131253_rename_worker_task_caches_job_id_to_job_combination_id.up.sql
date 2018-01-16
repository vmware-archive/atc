BEGIN;

  ALTER TABLE worker_task_caches RENAME job_id TO job_combination_id;

  ALTER TABLE worker_task_caches DROP CONSTRAINT worker_task_caches_job_id_fkey;

  ALTER TABLE worker_task_caches ADD CONSTRAINT worker_task_caches_job_combination_id_fkey FOREIGN KEY (job_combination_id) REFERENCES job_combinations(id) ON DELETE CASCADE;

  ALTER INDEX worker_task_caches_job_id RENAME TO worker_task_caches_job_combination_id;

COMMIT;
