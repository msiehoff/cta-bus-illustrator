#!/usr/bin/env bash
set -euo pipefail

# Imports route segments into the deployed API by:
# 1) fetching route IDs from /api/v1/routes
# 2) fetching CTA getpatterns JSON per route
# 3) POSTing the JSON to /api/v1/routes/:externalId/segments
#
# Required env vars:
# - API_BASE_URL   e.g. https://d1opmbjrozxley.cloudfront.net
# - CTA_API_KEY    CTA Bus Tracker API key
#
# Optional env vars:
# - ROUTE_IDS      comma-separated list; when set, skips /routes fetch
# - SLEEP_SECONDS  delay between routes (default: 0.2)

if [[ -z "${API_BASE_URL:-}" ]]; then
  echo "API_BASE_URL is required (e.g. https://your-api-host)" >&2
  exit 1
fi

if [[ -z "${CTA_API_KEY:-}" ]]; then
  echo "CTA_API_KEY is required" >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required for JSON parsing and URL encoding" >&2
  exit 1
fi

SLEEP_SECONDS="${SLEEP_SECONDS:-0.2}"

urlencode() {
  python3 -c 'import sys, urllib.parse; print(urllib.parse.quote(sys.argv[1], safe=""))' "$1"
}

if [[ -n "${ROUTE_IDS:-}" ]]; then
  IFS=',' read -r -a ROUTE_IDS_ARR <<< "${ROUTE_IDS}"
else
  echo "Fetching route IDs from ${API_BASE_URL}/api/v1/routes ..."
  ROUTES_JSON="$(curl --fail --silent --show-error "${API_BASE_URL}/api/v1/routes")"
  mapfile -t ROUTE_IDS_ARR < <(
    python3 -c '
import json, sys
doc = json.loads(sys.stdin.read())
for f in doc.get("features", []):
    rid = (f.get("properties") or {}).get("routeId")
    if rid:
        print(rid)
' <<< "${ROUTES_JSON}"
  )
fi

if [[ ${#ROUTE_IDS_ARR[@]} -eq 0 ]]; then
  echo "No route IDs found to import." >&2
  exit 1
fi

echo "Importing segments for ${#ROUTE_IDS_ARR[@]} routes ..."

ok=0
fail=0

for route_id in "${ROUTE_IDS_ARR[@]}"; do
  trimmed_route_id="$(echo "${route_id}" | xargs)"
  if [[ -z "${trimmed_route_id}" ]]; then
    continue
  fi

  cta_json_tmp="$(mktemp)"
  encoded_route_id="$(urlencode "${trimmed_route_id}")"

  echo "[$((ok + fail + 1))/${#ROUTE_IDS_ARR[@]}] route ${trimmed_route_id}: fetching CTA pattern ..."
  if ! curl --fail --silent --show-error \
    --get "https://www.ctabustracker.com/bustime/api/v3/getpatterns" \
    --data-urlencode "key=${CTA_API_KEY}" \
    --data-urlencode "format=json" \
    --data-urlencode "rt=${trimmed_route_id}" \
    > "${cta_json_tmp}"; then
    echo "  -> CTA fetch failed for route ${trimmed_route_id}" >&2
    rm -f "${cta_json_tmp}"
    fail=$((fail + 1))
    sleep "${SLEEP_SECONDS}"
    continue
  fi

  echo "  -> posting to ${API_BASE_URL}/api/v1/routes/${encoded_route_id}/segments"
  if curl --fail --silent --show-error \
    -X POST "${API_BASE_URL}/api/v1/routes/${encoded_route_id}/segments" \
    -H "Content-Type: application/json" \
    --data-binary "@${cta_json_tmp}" \
    > /dev/null; then
    ok=$((ok + 1))
  else
    echo "  -> segment import failed for route ${trimmed_route_id}" >&2
    fail=$((fail + 1))
  fi

  rm -f "${cta_json_tmp}"
  sleep "${SLEEP_SECONDS}"
done

echo "Done. Imported: ${ok}, failed: ${fail}"
if [[ ${fail} -gt 0 ]]; then
  exit 1
fi
