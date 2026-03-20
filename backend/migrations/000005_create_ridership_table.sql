-- +goose Up
CREATE TABLE ridership (
    id                 BIGSERIAL PRIMARY KEY,
    created_at         TIMESTAMP WITH TIME ZONE,
    updated_at         TIMESTAMP WITH TIME ZONE,
    deleted_at         TIMESTAMP WITH TIME ZONE,
    route_id           BIGINT NOT NULL REFERENCES routes (id),
    month_beginning    DATE NOT NULL,
    avg_weekday_rides  DECIMAL(10, 2),
    avg_saturday_rides DECIMAL(10, 2),
    avg_sunday_rides   DECIMAL(10, 2),
    month_total        INTEGER
);

CREATE UNIQUE INDEX idx_ridership_route_month ON ridership (route_id, month_beginning);
CREATE INDEX idx_ridership_deleted_at ON ridership (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS ridership;
