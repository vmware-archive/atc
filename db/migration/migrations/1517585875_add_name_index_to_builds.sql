-- +goose Up
-- +goose StatementBegin
BEGIN;
  CREATE INDEX builds_name ON builds USING btree (name);
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
  DROP INDEX builds_name;
COMMIT;
-- +goose StatementEnd
