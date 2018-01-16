BEGIN;

  ALTER TABLE worker_task_caches RENAME job_combination_id TO job_id;

  ALTER TABLE worker_task_caches DROP CONSTRAINT worker_task_caches_job_combination_id_fkey;

  ALTER TABLE worker_task_caches ADD CONSTRAINT worker_task_caches_job_id_fkey FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE;

  ALTER INDEX worker_task_caches_job_combination_id RENAME TO worker_task_caches_job_id;

COMMIT;
