-- +goose Up

ALTER TABLE headways
    ADD COLUMN IF NOT EXISTS from_vehicle_id VARCHAR,
    ADD COLUMN IF NOT EXISTS to_vehicle_id VARCHAR;

-- Idempotent upserts for a service day: one observed gap per later-arrival timestamp.
CREATE UNIQUE INDEX IF NOT EXISTS idx_headways_unique_observed
    ON headways (stop_id, route_id, direction, timestamp);

CREATE TABLE headway_job_runs (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    service_date        DATE NOT NULL,
    status              VARCHAR NOT NULL,
    triggered_by        VARCHAR NOT NULL,
    started_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    finished_at         TIMESTAMP WITH TIME ZONE,
    arrivals_processed  INTEGER NOT NULL DEFAULT 0,
    headways_written    INTEGER NOT NULL DEFAULT 0,
    error_message       TEXT
);

CREATE INDEX idx_headway_job_runs_service_date ON headway_job_runs (service_date DESC);
CREATE INDEX idx_headway_job_runs_started_at ON headway_job_runs (started_at DESC);

-- +goose Down

DROP TABLE IF EXISTS headway_job_runs;
DROP INDEX IF EXISTS idx_headways_unique_observed;
ALTER TABLE headways
    DROP COLUMN IF EXISTS from_vehicle_id,
    DROP COLUMN IF EXISTS to_vehicle_id;
