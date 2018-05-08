BEGIN;

  ALTER TABLE jobs ADD COLUMN inputs_determined boolean DEFAULT false NOT NULL;
  ALTER TABLE jobs ADD COLUMN build_number_seq integer DEFAULT 0 NOT NULL;

  UPDATE jobs SET build_number_seq = (
    SELECT COUNT(*) FROM builds b
    LEFT JOIN job_combinations c ON b.job_combination_id = c.id
    LEFT JOIN jobs j ON c.job_id = j.id
    WHERE j.id = jobs.id
  );

COMMIT;
