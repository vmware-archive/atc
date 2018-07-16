-- +goose Up
-- +goose StatementBegin
BEGIN;
  ALTER TABLE resource_caches DROP CONSTRAINT resource_caches_resource_config_id_version_params_hash_key;

  CREATE UNIQUE INDEX resource_caches_resource_config_id_version_params_hash_key ON resource_caches (resource_config_id, md5(version), params_hash);
COMMIT;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
BEGIN;
  DROP INDEX resource_caches_resource_config_id_version_params_hash_key;

  ALTER TABLE ONLY resource_caches ADD CONSTRAINT resource_caches_resource_config_id_version_params_hash_key UNIQUE (resource_config_id, version, params_hash);
COMMIT;
-- +goose StatementEnd
