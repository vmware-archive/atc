BEGIN;
  ALTER TABLE containers DROP COLUMN meta_resource_config_id;
COMMIT;
