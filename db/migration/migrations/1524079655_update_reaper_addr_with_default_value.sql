-- +goose Up
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers ALTER COLUMN reaper_addr SET DEFAULT ''::text;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers ALTER COLUMN reaper_addr DROP DEFAULT;

COMMIT;
-- +goose StatementEnd
