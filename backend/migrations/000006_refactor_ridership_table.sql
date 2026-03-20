-- +goose Up
DROP TABLE IF EXISTS ridership;

CREATE TABLE ridership (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE,
    route_id        BIGINT NOT NULL REFERENCES routes (id),
    month_beginning DATE NOT NULL,
    type            VARCHAR NOT NULL,
    avg_rides       DECIMAL(10, 2) NOT NULL
);

CREATE UNIQUE INDEX idx_ridership_route_month_type ON ridership (route_id, month_beginning, type);
CREATE INDEX idx_ridership_deleted_at ON ridership (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS ridership;
