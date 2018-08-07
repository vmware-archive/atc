BEGIN;
  ALTER TABLE containers ADD COLUMN meta_resource_config_id integer NOT NULL DEFAULT 0;
COMMIT;
