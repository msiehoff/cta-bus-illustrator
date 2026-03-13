-- +goose Up
INSERT INTO route_segments (route_id, sequence, lat, lng, created_at, updated_at)
SELECT r.id, s.sequence, s.lat, s.lng, NOW(), NOW()
FROM routes r
JOIN (VALUES
    (1,  41.8957, -87.8065),
    (2,  41.8957, -87.7754),
    (3,  41.8957, -87.7523),
    (4,  41.8957, -87.7341),
    (5,  41.8957, -87.7137),
    (6,  41.8957, -87.6948),
    (7,  41.8957, -87.6726),
    (8,  41.8957, -87.6558),
    (9,  41.8957, -87.6418),
    (10, 41.8957, -87.6264),
    (11, 41.8957, -87.6134),
    (12, 41.8957, -87.6062),
    (13, 41.8966, -87.5961),
    (14, 41.8978, -87.5854),
    (15, 41.8983, -87.5769)
) AS s(sequence, lat, lng) ON TRUE
WHERE r.external_id = '66';

-- +goose Down
DELETE FROM route_segments
WHERE route_id = (SELECT id FROM routes WHERE external_id = '66');
