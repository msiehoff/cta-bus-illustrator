-- +goose Up

CREATE TABLE headways (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    stop_id         VARCHAR NOT NULL,
    route_id        VARCHAR NOT NULL,
    direction       VARCHAR NOT NULL,
    timestamp       TIMESTAMP WITH TIME ZONE NOT NULL,
    headway_minutes DOUBLE PRECISION NOT NULL
);

CREATE INDEX idx_headways_stop_route_dir_time ON headways (stop_id, route_id, direction, timestamp);

-- +goose Down
DROP TABLE IF EXISTS headways;
