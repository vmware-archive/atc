BEGIN;

  ALTER TABLE jobs DROP COLUMN inputs_determined;

COMMIT;
