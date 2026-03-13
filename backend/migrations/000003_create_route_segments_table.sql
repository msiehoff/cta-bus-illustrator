-- +goose Up
CREATE TABLE route_segments (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    route_id   BIGINT NOT NULL REFERENCES routes (id),
    sequence   INT NOT NULL,
    lat        DOUBLE PRECISION NOT NULL,
    lng        DOUBLE PRECISION NOT NULL
);

CREATE INDEX idx_route_segments_route_id ON route_segments (route_id);
CREATE INDEX idx_route_segments_deleted_at ON route_segments (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS route_segments;
