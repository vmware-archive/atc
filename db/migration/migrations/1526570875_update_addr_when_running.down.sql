ALTER TABLE workers DROP CONSTRAINT addr_when_running;

ALTER TABLE workers ADD CONSTRAINT addr_when_running CHECK ((((state <> 'stalled'::worker_state) AND (state <> 'landed'::worker_state) AND ((addr IS NOT NULL) OR (baggageclaim_url IS NOT NULL))) OR (((state = 'stalled'::worker_state) OR (state = 'landed'::worker_state)) AND (addr IS NULL) AND (baggageclaim_url IS NULL))));

DROP TYPE worker_state;

CREATE TYPE worker_state AS ENUM (
    'running',
    'stalled',
    'landing',
    'landed',
    'retiring'
);
