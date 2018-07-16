-- +goose Up
-- +goose StatementBegin
BEGIN;

  ALTER TABLE jobs ADD COLUMN tags text[];

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

  ALTER TABLE jobs DROP COLUMN tags;

COMMIT;
-- +goose StatementEnd
