ALTER TABLE workers DROP CONSTRAINT addr_when_running;

ALTER TABLE workers ADD CONSTRAINT addr_when_running CHECK ((
((state = 'stalled'::worker_state) AND (addr IS NULL) OR (baggageclaim_url IS NULL))
OR
((state = 'running'::worker_state) AND (addr IS NOT NULL) AND (baggageclaim_url IS NOT NULL))
OR
((state = 'retiring'::worker_state) AND ((addr IS NOT NULL) OR (addr IS NULL)) AND ((baggageclaim_url IS NOT NULL) OR (baggageclaim_url IS NULL)))
));

UPDATE workers SET state='retiring' WHERE state='landing' OR state='landed';

DROP TYPE worker_state;

CREATE TYPE worker_state AS ENUM (
      'running',
      'stalled',
      'retiring'
);
