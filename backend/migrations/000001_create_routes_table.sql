-- +goose Up
CREATE TABLE routes (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    external_id VARCHAR NOT NULL,
    name        VARCHAR NOT NULL
);

CREATE UNIQUE INDEX idx_routes_external_id ON routes (external_id);
CREATE INDEX idx_routes_deleted_at ON routes (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS routes;
