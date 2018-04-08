BEGIN;
  ALTER TABLE "workers" DROP COLUMN "type";
  DROP TYPE worker_type;
COMMIT;
