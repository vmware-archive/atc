-- +goose Up
-- +goose StatementBegin
BEGIN;
  ALTER TABLE pipelines ADD COLUMN cache_index integer NOT NULL DEFAULT 1;
  ALTER TABLE versioned_resources DROP COLUMN modified_time;
  ALTER TABLE build_inputs DROP COLUMN modified_time;
  ALTER TABLE build_outputs DROP COLUMN modified_time;
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
  ALTER TABLE pipelines DROP COLUMN cache_index;
  ALTER TABLE versioned_resources ADD COLUMN modified_time timestamp without time zone DEFAULT now() NOT NULL;
  ALTER TABLE build_inputs ADD COLUMN modified_time timestamp without time zone DEFAULT now() NOT NULL;
  ALTER TABLE build_outputs ADD COLUMN modified_time timestamp without time zone DEFAULT now() NOT NULL;
COMMIT;
-- +goose StatementEnd
