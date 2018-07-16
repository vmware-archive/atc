-- +goose Up
-- +goose StatementBegin
BEGIN;
  DELETE FROM "build_outputs" WHERE NOT "explicit";
  ALTER TABLE "build_outputs" DROP COLUMN "explicit";
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
  ALTER TABLE "build_outputs" ADD COLUMN "explicit" boolean;
COMMIT;
-- +goose StatementEnd
