#!/usr/bin/env bash
# Download CTA bus ridership CSV from Chicago Open Data (SoDota) and import it
# into the deployed API, then refresh route segments.
#
# Required env vars:
# - RIDERSHIP_API_URL   Backend API origin base URL (not the public CDN/site host),
#                       e.g. https://your-api-host.example.com
#
# Optional env vars:
# - SODATA_DATASET_ID   SoDota view id (default: bynn-gwxy)
# - SODATA_DOMAIN       Open Data host (default: data.cityofchicago.org)
# - FORCE_IMPORT        If "true", import even when rowsUpdatedAt is unchanged
# - STAMP_DIR           Directory for last-import stamp (default: .ridership-dataset-stamp)
# - CURL_MAX_TIME       Max seconds per curl call (default: 600)
#
# Skip logic: fetches SoDota metadata and compares rowsUpdatedAt to STAMP_DIR.
# When unchanged (and FORCE_IMPORT is not true), exits 0 without calling the API.
# After a successful import, writes the new stamp so CI cache can persist it.
set -euo pipefail

RIDERSHIP_API_URL="${RIDERSHIP_API_URL:-}"
SODATA_DATASET_ID="${SODATA_DATASET_ID:-bynn-gwxy}"
SODATA_DOMAIN="${SODATA_DOMAIN:-data.cityofchicago.org}"
FORCE_IMPORT="${FORCE_IMPORT:-false}"
STAMP_DIR="${STAMP_DIR:-.ridership-dataset-stamp}"
CURL_MAX_TIME="${CURL_MAX_TIME:-600}"

if [[ -z "${RIDERSHIP_API_URL}" ]]; then
  echo "RIDERSHIP_API_URL is required (API origin base URL)" >&2
  exit 1
fi

# Trim trailing slash for consistent path joining.
RIDERSHIP_API_URL="${RIDERSHIP_API_URL%/}"

META_URL="https://${SODATA_DOMAIN}/api/views/${SODATA_DATASET_ID}.json"
CSV_URL="https://${SODATA_DOMAIN}/api/views/${SODATA_DATASET_ID}/rows.csv?accessType=DOWNLOAD"
STAMP_FILE="${STAMP_DIR}/rows_updated_at"

curl_api() {
  # Fail on HTTP errors; print response body on failure for Actions logs.
  local method="$1"
  local url="$2"
  shift 2
  local tmp
  tmp="$(mktemp)"
  local http_code
  set +e
  http_code="$(
    curl --silent --show-error --max-time "${CURL_MAX_TIME}" \
      -o "${tmp}" -w "%{http_code}" \
      -X "${method}" "${url}" "$@"
  )"
  local curl_status=$?
  set -e
  if [[ ${curl_status} -ne 0 ]]; then
    echo "curl failed (exit ${curl_status}) for ${method} ${url}" >&2
    cat "${tmp}" >&2 || true
    rm -f "${tmp}"
    exit 1
  fi
  if [[ "${http_code}" -lt 200 || "${http_code}" -ge 300 ]]; then
    echo "HTTP ${http_code} from ${method} ${url}" >&2
    cat "${tmp}" >&2 || true
    rm -f "${tmp}"
    exit 1
  fi
  echo "HTTP ${http_code} ${method} ${url}"
  cat "${tmp}"
  echo
  rm -f "${tmp}"
}

echo "Fetching SoDota metadata: ${META_URL}"
META_JSON="$(curl --fail --silent --show-error --max-time 60 "${META_URL}")"
ROWS_UPDATED_AT="$(
  python3 -c 'import json,sys; print(json.load(sys.stdin)["rowsUpdatedAt"])' <<<"${META_JSON}"
)"
if [[ -z "${ROWS_UPDATED_AT}" ]]; then
  echo "Could not parse rowsUpdatedAt from SoDota metadata" >&2
  exit 1
fi
echo "SoDota rowsUpdatedAt=${ROWS_UPDATED_AT}"

if [[ "${FORCE_IMPORT}" != "true" && -f "${STAMP_FILE}" ]]; then
  PREV="$(tr -d '[:space:]' <"${STAMP_FILE}")"
  if [[ "${PREV}" == "${ROWS_UPDATED_AT}" ]]; then
    echo "Dataset unchanged since last import; skipping (set FORCE_IMPORT=true to override)."
    exit 0
  fi
  echo "Dataset updated (previous stamp=${PREV}); importing."
else
  if [[ "${FORCE_IMPORT}" == "true" ]]; then
    echo "FORCE_IMPORT=true; importing regardless of stamp."
  else
    echo "No previous stamp found; importing."
  fi
fi

CSV_TMP="$(mktemp "${TMPDIR:-/tmp}/ridership.XXXXXX.csv")"
cleanup() {
  rm -f "${CSV_TMP}"
}
trap cleanup EXIT

echo "Downloading CSV: ${CSV_URL}"
curl --fail --silent --show-error --location --max-time "${CURL_MAX_TIME}" \
  -o "${CSV_TMP}" "${CSV_URL}"

if [[ ! -s "${CSV_TMP}" ]]; then
  echo "Downloaded CSV is empty" >&2
  exit 1
fi

# Sanity-check expected header columns from the SoDota export.
HEADER="$(head -n 1 "${CSV_TMP}")"
for col in route routename Month_Beginning Avg_Weekday_Rides Avg_Saturday_Rides MonthTotal; do
  if ! grep -q "${col}" <<<"${HEADER}"; then
    echo "CSV header missing expected column '${col}': ${HEADER}" >&2
    exit 1
  fi
done
echo "CSV OK ($(wc -l <"${CSV_TMP}" | tr -d ' ') lines including header)"

echo "POST ${RIDERSHIP_API_URL}/api/v1/ridership/import"
curl_api POST "${RIDERSHIP_API_URL}/api/v1/ridership/import" \
  -F "file=@${CSV_TMP};type=text/csv;filename=cta-ridership.csv"

echo "POST ${RIDERSHIP_API_URL}/api/v1/routes/import-segments"
curl_api POST "${RIDERSHIP_API_URL}/api/v1/routes/import-segments"

mkdir -p "${STAMP_DIR}"
printf '%s\n' "${ROWS_UPDATED_AT}" >"${STAMP_FILE}"
echo "Wrote stamp ${STAMP_FILE}=${ROWS_UPDATED_AT}"
echo "Ridership import complete."
