BEGIN;

  ALTER TABLE jobs ADD COLUMN inputs_determined boolean DEFAULT false NOT NULL;

COMMIT;
