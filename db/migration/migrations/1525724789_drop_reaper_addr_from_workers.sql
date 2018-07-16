-- +goose Up
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers DROP COLUMN reaper_addr;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers ADD reaper_addr text DEFAULT '';

COMMIT;
-- +goose StatementEnd
