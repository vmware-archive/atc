BEGIN;
  ALTER TABLE containers DROP COLUMN meta_resource_type_name;
COMMIT;
