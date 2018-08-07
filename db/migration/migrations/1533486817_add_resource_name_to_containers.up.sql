BEGIN;
  ALTER TABLE containers ADD COLUMN meta_resource_name text NOT NULL DEFAULT '';
COMMIT;
