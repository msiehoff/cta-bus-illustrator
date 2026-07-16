#!/usr/bin/env bash
# Local-only: insert sample arrivals and run headway rollups so /headways UI has data.
# Does NOT run in production migrations — invoke via `make seed-headways`.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Load backend/.env if present (DATABASE_URL, HEADWAY_JOB_TOKEN, …).
if [[ -f "$ROOT/backend/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT/backend/.env"
  set +a
fi

APP_URL="${APP_URL:-http://localhost:8080}"
TOKEN="${TOKEN:-${HEADWAY_JOB_TOKEN:-}}"

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is required (set in backend/.env or the environment)." >&2
  echo "In-memory fake mode has no durable arrivals to seed." >&2
  exit 1
fi

if ! command -v psql >/dev/null 2>&1; then
  echo "psql is required to seed arrivals." >&2
  exit 1
fi

echo "Seeding sample arrivals into local database…"

psql "$DATABASE_URL" -v ON_ERROR_STOP=1 <<'SQL'
-- Idempotent: remove prior seed rows only.
DELETE FROM arrivals WHERE stop_id LIKE 'seed-%';

-- Three Chicago service days, two routes, both directions, ~6–12 min gaps.
INSERT INTO arrivals (stop_id, route_id, direction, vehicle_id, timestamp)
SELECT
  'seed-' || route_id || '-' || dir_code || '-' || stop_n,
  route_id,
  direction,
  'seed-v-' || route_id || '-' || dir_code || '-' || d || '-' || bus_n,
  (
    (DATE '2026-07-10' + d)
    + TIME '07:00'
    + (stop_n * INTERVAL '90 seconds')
    + (bus_n * gap)
  ) AT TIME ZONE 'America/Chicago'
FROM (
  VALUES
    ('8',  'Northbound', 'nb', INTERVAL '8 minutes'),
    ('8',  'Southbound', 'sb', INTERVAL '10 minutes'),
    ('22', 'Northbound', 'nb', INTERVAL '7 minutes'),
    ('22', 'Southbound', 'sb', INTERVAL '9 minutes')
) AS routes(route_id, direction, dir_code, gap),
  generate_series(0, 2) AS d,          -- Jul 10–12
  generate_series(1, 3) AS stop_n,     -- 3 stops per direction
  generate_series(0, 15) AS bus_n;     -- 16 arrivals → 15 gaps per stop

SELECT COUNT(*) AS seed_arrivals FROM arrivals WHERE stop_id LIKE 'seed-%';
SQL

DATES=(2026-07-10 2026-07-11 2026-07-12)

if [[ -z "$TOKEN" ]]; then
  echo
  echo "Arrivals seeded. To build headway summaries, start the API and run:"
  for d in "${DATES[@]}"; do
    echo "  make headway-run TOKEN=<HEADWAY_JOB_TOKEN> DATE=$d"
  done
  echo "Or open Admin → Headway Jobs and run those dates."
  exit 0
fi

echo "Running headway rollups for ${DATES[*]} against $APP_URL …"

for d in "${DATES[@]}"; do
  echo "  → $d"
  curl -sf -X POST "$APP_URL/api/v1/admin/headways/run" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"service_date\":\"$d\"}" \
    || {
      echo >&2
      echo "Rollup failed for $d. Is the API running at $APP_URL with HEADWAY_JOB_TOKEN set?" >&2
      echo "Arrivals are seeded; re-run rollups once the server is up:" >&2
      echo "  make headway-run TOKEN=\$HEADWAY_JOB_TOKEN DATE=$d" >&2
      exit 1
    }
  echo
done

echo "Done. Open /headways/routes and /headways/system (routes 8 and 22)."
