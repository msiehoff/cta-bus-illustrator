-- +goose Up

ALTER TABLE headway_summaries
    ADD COLUMN IF NOT EXISTS window_key VARCHAR NOT NULL DEFAULT 'overall';

DROP INDEX IF EXISTS idx_headway_summaries_unique;
CREATE UNIQUE INDEX idx_headway_summaries_unique
    ON headway_summaries (service_date, window_key, grain, method, route_id, direction, stop_id);

DROP INDEX IF EXISTS idx_headway_summaries_lookup;
CREATE INDEX idx_headway_summaries_lookup
    ON headway_summaries (service_date, window_key, route_id, direction, stop_id);

-- +goose Down

DROP INDEX IF EXISTS idx_headway_summaries_unique;
DROP INDEX IF EXISTS idx_headway_summaries_lookup;

ALTER TABLE headway_summaries DROP COLUMN IF EXISTS window_key;

CREATE UNIQUE INDEX idx_headway_summaries_unique
    ON headway_summaries (service_date, grain, method, route_id, direction, stop_id);

CREATE INDEX idx_headway_summaries_lookup
    ON headway_summaries (service_date, route_id, direction, stop_id);
