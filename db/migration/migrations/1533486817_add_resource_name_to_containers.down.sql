BEGIN;
  ALTER TABLE containers DROP COLUMN meta_resource_name;
COMMIT;
