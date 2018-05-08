BEGIN;

  ALTER TABLE jobs DROP COLUMN inputs_determined;
  ALTER TABLE jobs DROP COLUMN build_number_seq;

COMMIT;
