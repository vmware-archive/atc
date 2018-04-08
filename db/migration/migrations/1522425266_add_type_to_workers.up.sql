BEGIN;
  CREATE TYPE worker_type AS ENUM ( 'garden', 'kubernetes' );
  ALTER TABLE "workers" ADD COLUMN "type" worker_type DEFAULT 'garden'::worker_type;
COMMIT;
