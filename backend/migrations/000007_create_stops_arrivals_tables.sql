-- +goose Up

CREATE TABLE stops (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    stop_id     VARCHAR NOT NULL,
    route_id    VARCHAR NOT NULL,
    direction   VARCHAR NOT NULL,
    name        VARCHAR NOT NULL,
    lat         DOUBLE PRECISION NOT NULL,
    lon         DOUBLE PRECISION NOT NULL,
    sequence    INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_stops_stop_route_dir ON stops (stop_id, route_id, direction);
CREATE INDEX idx_stops_route_dir ON stops (route_id, direction);

CREATE TABLE arrivals (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    stop_id     VARCHAR NOT NULL,
    route_id    VARCHAR NOT NULL,
    direction   VARCHAR NOT NULL,
    vehicle_id  VARCHAR NOT NULL,
    timestamp   TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Primary query pattern: fetch arrivals at a stop in time order for headway calculation.
CREATE INDEX idx_arrivals_stop_route_dir_time ON arrivals (stop_id, route_id, direction, timestamp);
CREATE INDEX idx_arrivals_vehicle_time ON arrivals (vehicle_id, timestamp);

-- +goose Down
DROP TABLE IF EXISTS arrivals;
DROP TABLE IF EXISTS stops;
