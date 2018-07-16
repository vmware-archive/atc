-- +goose Up
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers ADD COLUMN reaper_addr text;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

  ALTER TABLE workers DROP COLUMN reaper_addr;

COMMIT;
-- +goose StatementEnd
