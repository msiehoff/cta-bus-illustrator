-- +goose Up

CREATE TABLE headway_summaries (
    id                  BIGSERIAL PRIMARY KEY,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    service_date        DATE NOT NULL,
    window_start        TIMESTAMP WITH TIME ZONE NOT NULL,
    window_end          TIMESTAMP WITH TIME ZONE NOT NULL,
    -- stop | route_direction | service_day
    grain               VARCHAR NOT NULL,
    -- pooled | equal_stop
    method              VARCHAR NOT NULL,
    stop_id             VARCHAR NOT NULL DEFAULT '',
    route_id            VARCHAR NOT NULL DEFAULT '',
    direction           VARCHAR NOT NULL DEFAULT '',
    observation_count   INTEGER NOT NULL DEFAULT 0,
    mean_minutes        DOUBLE PRECISION NOT NULL DEFAULT 0,
    median_minutes      DOUBLE PRECISION NOT NULL DEFAULT 0,
    stddev_minutes      DOUBLE PRECISION NOT NULL DEFAULT 0,
    cv                  DOUBLE PRECISION NOT NULL DEFAULT 0,
    avg_wait_minutes    DOUBLE PRECISION NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_headway_summaries_unique
    ON headway_summaries (service_date, grain, method, route_id, direction, stop_id);

CREATE INDEX idx_headway_summaries_lookup
    ON headway_summaries (service_date, route_id, direction, stop_id);

ALTER TABLE headway_job_runs
    ADD COLUMN IF NOT EXISTS summaries_written INTEGER NOT NULL DEFAULT 0;

-- +goose Down

ALTER TABLE headway_job_runs
    DROP COLUMN IF EXISTS summaries_written;

DROP TABLE IF EXISTS headway_summaries;
