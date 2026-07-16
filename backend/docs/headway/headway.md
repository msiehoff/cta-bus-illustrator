# CTA Transit Lab | Headway

## Headway
### Metrics to Calculate
*   **Scheduled Headway:** scheduled average duration between bus arrival times
*   **Actual Headway:** average duration between bus arrival times
*   **Gap Rate:** Headway difference (between actual + scheduled)
*   **Average Wait Time:** actual headway/2
*   Overall & by time window (rush hour, evening, weekend...etc.)

### Calculating Headway
*   Detect stop arrivals (bus within radius of stop or past it) & save timestamp
*   Calculate Average headway at a stop
*   Calculate by time window (e.g. rush hour, weekend...etc.)

### API Request Limit
*   100,000 requests per day limit
*   127 routes
*   getvehicles → 10 routes max
*   updated every 30 seconds
*   13 x 2 x 60 x 18 → 28,080 requests to get all route data 18 hours of the day (37,440 for 24 hours)
*   Halsted bus runs 4:00 am to midnight (may be typical)

## Data
*   **stop:** name, lat/lon, route\_id, direction
*   **arrival:** stop\_id, route\_id, direction, timestamp, vehicle\_id
*   **headway:** route\_id, stop\_id, direction, timestamp (of the later bus), headway\_minutes
## Arrival Detector
*   **Problem:** It's easy to add duplicate arrivals. Buses may sit at a stop for ~20-40 seconds and GPS jitter can move it in/out of the radious
*   **Solution 1:** Stateful arrival detector
    *   Add a lightweight vehicle state

```plain
vehicle_state:
  vehicle_id
  last_stop_id
  last_arrival_time
```

*       *   Arrival Logic: for each vehicle ping
        *   Find nearest downstream stop along the route shap
        *   Check
            *   distance < threshold (30-50m)
            *   AND vehicle has not already recorded arrival at this stop recently
            *   Only create if `(vehicle_id, stop_id) not seen in last N minutes)`
*   **Solution 2:** detect passing, not just proximity (Even Better?)
    *   Track vehicle movement across time
    *   Detect when it crosses the stop position along the route

```cs
if previous_position is before stop
AND current_position is after stop
-> arrival event
```

(Could incorporate both)
**NOTE:** Need to account for terminal stops (buses linger there). Could exclude them.

## Pipeline

The pipeline has two layers: **ingestion** (continuous polling) and **analysis** (modular fan-out).

```
┌─────────────────┐
│  Data Ingestion │  polls CTA getvehicles, emits VehiclePing
└────────┬────────┘
         │ pings
         ▼
┌──────────────────────────────────────────────┐
│  Analyzer fan-out (each owns its persistence) │
├──────────────┬───────────────┬────────────────┤
│  Arrival     │  Segment      │  Ghost Bus     │  (future PR)
│  Analyzer    │  Speed        │  Analyzer      │
│  → arrivals  │  → speeds     │  → ghosts      │
└──────────────┴───────────────┴────────────────┘
         │
         ▼ (batch / on schedule)
┌─────────────────┐
│ Headway Rollup  │  reads arrivals → writes headways
└─────────────────┘
```

**Ingestion stage**
*   Poll CTA `getvehicles` every 30–60 seconds for configured routes (batched, 10 routes per call)
*   Load stops from CTA `getstops` at startup and persist to DB
*   Emit `VehiclePing` to analyzers

**Analyzer interface (design goal)**
Each analyzer is isolated, owns its DB writes, and consumes the same ping stream:

```go
type Analyzer interface {
    Name() string
    ProcessPing(ctx context.Context, ping business.VehiclePing) error
}
```

*   **ArrivalAnalyzer** — proximity + cooldown detection → `arrivals` table (single-worker, order-sensitive)
*   **SpeedAnalyzer** — consecutive stop arrivals per vehicle → `segment_speeds` table (v1 add-on, PR 4)
*   **GhostBusAnalyzer** — scheduled vs actual correlation → own PR (needs GTFS)
*   **HeadwayRollup** — batch job on `arrivals`, not per-ping → `headways` table

Arrival detection must remain **single-worker** for correctness. Fan-out happens after ingestion; v1 dispatches pings to analyzers sequentially in one goroutine.

**NOTE:** A raw vehicle ping log could be useful for replaying & debugging.

* * *
# Technical Specification
# CTA Bus Headway Tracking — Design Context + Incremental Implementation Plan
## Project Context
I’m building a CTA bus analytics site focused on:
*   bus reliability
*   headways
*   bunching
*   corridor performance
*   eventually segment speed visualization and ghost bus metrics
Current data sources:
*   CTA Bus Tracker API
    *   routes
    *   stops
    *   vehicle locations (`getvehicles`)
Vehicle data is polled every ~30–60 seconds.
Tech stack:
*   Golang backend
*   Likely a pipeline/channel-based architecture
*   Relational DB for persistence
Design Principles
*   Onion architecture
*   CTA API as an adapter
*   Pipeline & concurrency logic should be separate from domain logic (including detecting headways) and should be unit testable
Goal for this PR:
1. Scheduled route headway
2. Actual route headway
Future goals:
3\. Segment speed analysis
4\. Reliability / bunching metrics
5\. Ghost bus detection
* * *
# Core Conceptual Model
## Important distinction
Headway must initially be calculated:
*   per stop
*   per direction
*   over a time window
Do NOT average blindly across all stops.
Everything downstream depends on correctly detecting:

```css
(bus arrives at stop X at time T)


```

This “arrival event” becomes the foundational primitive for:
*   actual headways
*   wait time
*   bunching
*   reliability
*   speed between stops
* * *
# Proposed Data Models
## stop

```yaml
stop:
  stop_id
  route_id
  direction
  name
  lat
  lon
  sequence


```

* * *
## arrival
Represents a detected real-world bus arrival at a stop.

```sql
arrival:
  stop_id
  route_id
  direction
  vehicle_id
  timestamp


```

* * *
## headway (observed)
Represents the gap between two consecutive arrivals at a **single stop**.
This is the atomic fact — one row per consecutive bus pair.

```yaml
headway:
  stop_id
  route_id
  direction
  timestamp              # later of the two arrivals
  headway_minutes
  from_vehicle_id        # optional — earlier bus
  to_vehicle_id          # optional — later bus
```

**Idempotency:** unique on `(stop_id, route_id, direction, timestamp)`.
Re-running a day upserts/replaces gaps for that day without duplicates.

Timestamp should represent the later arrival.

* * *
## headway_job_runs
Tracks each headway computation run (daily cron or admin-triggered).

```yaml
headway_job_runs:
  id
  service_date           # calendar day being computed (America/Chicago)
  status                 # pending | running | success | failed
  triggered_by           # cron | admin | api
  started_at
  finished_at
  arrivals_processed
  headways_written
  summaries_written
  error_message          # nullable
```

* * *
## headway_summaries
Persisted aggregates written by the same daily job as observed headways.

```yaml
headway_summaries:
  service_date
  window_start / window_end   # America/Chicago day bounds
  grain                      # stop | route_direction | service_day
  method                     # pooled | equal_stop
  stop_id / route_id / direction  # empty string when unused for grain
  observation_count
  mean_minutes / median_minutes / stddev_minutes / cv / avg_wait_minutes
```

* * *
# Headway Computation Model

## Conceptual model (transit practice)

**Atomic unit = observed headway at a stop.**
At `(route, direction, stop)`, sort arrivals `t1, t2, t3, …` and compute:

```plain
headway_i = t_i − t_(i−1)
```

That is “how long until the next bus showed up.” Time-window averages
(rush hour, Saturday, overall, …) are derived from those gaps.

**Do not** blindly pool gaps from every stop into one “route headway”
without documenting the aggregation. Prefer:

| Grain | Meaning |
|-------|---------|
| **Stop + route + direction** | Primary — most correct |
| **Route + direction** | e.g. mean of per-stop means (equal weight per stop) |
| **Route** | Both directions or worse direction — least specific; document choice |

## Metrics for a time period

| Metric | Formula / notes | When |
|--------|-----------------|------|
| **Mean headway** | average of observed gaps | PR 3 |
| **Median headway** | less skewed by huge gaps | PR 3 |
| **Count** | number of observed gaps | PR 3 |
| **CV (coefficient of variation)** | `stddev(headway) / mean(headway)` | PR 3 |
| **Avg wait** | ≈ mean/2 if regular; else `E[h²]/(2E[h])` | PR 3 (simple), refine later |
| **% gapped / bunched** | vs scheduled (e.g. >2× or <0.5×) | PR 4 (needs GTFS) |

### What is CV at a stop?

**CV = coefficient of variation** = standard deviation of observed headways
divided by their mean, for one `(stop, route, direction)` in a time window.

- **CV ≈ 0** — buses arrive like clockwork (gaps nearly identical)
- **CV ≈ 0.3–0.5** — typical moderate irregularity
- **High CV** (e.g. > 0.6–0.8) — unreliable: mixes bunches and long gaps

Mean alone can look fine while CV is terrible (e.g. average 10 min made of
2-min and 18-min gaps). CV is the main **reliability** signal before we have
scheduled headways for bunching/gap % metrics.

## Observed vs summaries

| Layer | Storage | Role |
|-------|---------|------|
| **Observed headways** | `headways` table | Facts — idempotent daily rollup output |
| **Summaries** | `headway_summaries` table | Mean/median/CV by grain × method for that service date |

The daily job writes **both** layers for `service_date = D` (delete + insert).
Admin/rider APIs prefer stored summaries when a service date is set; they fall back
to compute-on-read only for ad-hoc filters (e.g. vehicle) or when no job has run yet.

Do **not** mix `"observed"` and `"weekday_am_avg"` as interchangeable row types
in one undifferentiated table. Facts are one grain; summaries are another
(`grain` + `method` on `headway_summaries`).

### headway_summaries grains

| Grain | Method | Meaning |
|-------|--------|---------|
| `stop` | `pooled` | Stats for one stop+route+direction |
| `route_direction` | `pooled` | All gaps on that route+direction |
| `route_direction` | `equal_stop` | Mean of per-stop means |
| `route` | `pooled` / `equal_stop` | Both directions on a route combined |
| `service_day` | `pooled` / `equal_stop` | Whole service date |

* * *
# Headway Daily Job

Headway rollup is a **separate batch job**, not part of the live vehicle
pipeline. Arrivals accumulate all day; headways are computed for a completed
service date (typically “yesterday” in America/Chicago).

## Why separate from the pipeline?

| Live pipeline | Headway job |
|---------------|-------------|
| Continuous, stateful | Batch, stateless over a date range |
| Latency-sensitive polling | Can take seconds–minutes |
| Restart loses in-memory vehicle state | Re-run any day safely (idempotent) |

## Triggering

One authenticated endpoint powers both admin UI and automation:

```http
POST /api/v1/admin/headways/run
{ "service_date": "2026-07-10" }   # optional; default = yesterday (Chicago)
```

**Auth (two mechanisms, same handler):**
1. **Admin session cookie** — admin UI “Run for date…” button
2. **Job token header** — cron / CI: `Authorization: Bearer $HEADWAY_JOB_TOKEN`
   (or `X-Job-Token`). Separate from user password so rotation is easy.

Not publicly callable without one of the above.

**Admin UI:**
- List `headway_job_runs` (date, status, counts, errors, duration)
- Trigger run for a specific date
- Link to sample headways for that day after success

## Idempotent day recompute

For `service_date = D`:
1. Insert `headway_job_runs` row (`status=running`)
2. Load arrivals with timestamps in `[D 00:00, D+1 00:00)` Chicago
3. Group by `(stop_id, route_id, direction)`, sort, compute gaps
4. Upsert into `headways` on unique `(stop_id, route_id, direction, timestamp)`
   — or delete existing observed headways whose timestamp falls on `D`, then insert
5. Mark run `success` / `failed` with metadata

Re-running the same day is safe and replaces prior results for that day.

## Daily cron — keep it simple (no AWS required)

Recommended order of simplicity:

| Option | How | Pros | Cons |
|--------|-----|------|------|
| **1. GitHub Actions `schedule`** | Workflow cron → `curl` with `HEADWAY_JOB_TOKEN` | Free, no extra infra, secrets in repo settings | Depends on GitHub uptime; ~15min schedule drift |
| **2. Platform cron** | Railway / Render / Fly cron job hitting the same URL | Same region as app | Only if already on that host |
| **3. Host crontab** | `0 6 * * * curl -H "Authorization: Bearer …" …` on a VPS | Classic, transparent | Need a host |

**v1 recommendation: GitHub Actions scheduled workflow** calling
`POST /api/v1/admin/headways/run` with the job token. Example sketch:

```yaml
# .github/workflows/headway-daily.yml
on:
  schedule:
    - cron: "0 11 * * *"   # 06:00 America/Chicago ≈ 11:00 UTC (adjust for DST)
  workflow_dispatch:       # manual re-run from GitHub UI
jobs:
  run:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger headway rollup
        run: |
          curl -sf -X POST "$APP_URL/api/v1/admin/headways/run" \
            -H "Authorization: Bearer $HEADWAY_JOB_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{}'
```

Avoid EventBridge, Lambda, SQS, etc. until there’s a clear need.

**Local / Makefile:** `make headway-run DATE=2026-07-10` → curl the same endpoint
(or `go run` a small `cmd/headway` that shares the service code).

* * *
# Scheduled Headways
CTA Bus Tracker does NOT provide schedule data.
Need CTA GTFS feed.
Relevant GTFS tables:
*   trips.txt
*   stop\_times.txt
Compute scheduled headways exactly the same way:
*   differences between scheduled arrivals at stops
* * *
# Important Metrics (detail)
## Reliability / Bunching
CV at a stop (see above):

```erlang
stddev(headway) / mean(headway)
```

With GTFS (PR 4+):

```erlang
% of headways > 2x scheduled   # gaps
% of headways < 0.5x scheduled # bunches
```

* * *
## Rider Wait Time
For irregular service:

```plain
avg_wait ≈ mean(headway^2) / (2 * mean(headway))
```

For roughly regular service, `mean(headway) / 2` is a fine first approximation.
* * *
## Segment Speed
Measure:

```css
time(stop B) - time(stop A)
```

per vehicle.
* * *
# Incremental Implementation Plan

Organized as small, reviewable PRs. Status reflects current progress.

## PR 1 — Wire Up Ingestion ✅

**Goal:** Pipeline runs locally and in deploy; persists arrivals (and stops) to Postgres.

| Step | Task | Status |
|------|------|--------|
| 1 | `stops` + `arrivals` tables | ✅ `000007` |
| 2 | `headways` table migration | ✅ `000008` (extend in PR 3) |
| 3 | Postgres `ArrivalRepository` + `StopRepository` | ✅ |
| 4 | Real `CTAVehicleClient` adapter over `storage/cta` | ✅ |
| 5 | Wire pipeline into `main.go` (`PIPELINE_ENABLED`, env config) | ✅ |
| 6 | Verify with fake CTA client → rows in `arrivals` | ✅ |

* * *
## PR 2 — Admin UI + Pipeline Observability ✅

**Goal:** Login-protected admin UI to monitor pipeline health and inspect arrivals.

| Step | Task | Status |
|------|------|--------|
| 1–4 | Auth, status, arrivals API + React pages | ✅ |
| Follow-up | Pattern-ID direction fix, arrivals UX, debug logs | ✅ #14 |

* * *
## PR 3 — Headway Rollup Job + Rider API

**Goal:** Idempotent daily (and on-demand) computation of **observed** headways;
admin visibility into job runs; rider-facing actual headway stats.

| Step | Task |
|------|------|
| 1 | Extend `headways` (vehicle IDs, unique constraint) + `headway_job_runs` migration |
| 2 | `HeadwayRollup` service: arrivals for `service_date` → observed gaps + summaries (idempotent) |
| 3 | `POST /api/v1/admin/headways/run` — session **or** `HEADWAY_JOB_TOKEN` |
| 4 | `GET /api/v1/admin/headways/runs` — list job metadata for admin UI |
| 5 | Admin page: run history + “Run for date” trigger |
| 6 | `headway_summaries` table + stats written by daily job; admin summary API reads stored rows |
| 7 | `GET /api/v1/routes/:id/headways?direction=&stop=&window=&from=&to=` |
| 8 | GitHub Actions daily workflow (simplest cron) |
| 9 | Frontend: replace “Coming soon” with headway summary + simple chart |

**Success criteria:** Re-running the same `service_date` twice does not duplicate rows;
admin can trigger and inspect runs; route page shows actual mean headway for a window.

* * *
## PR 4 — GTFS Scheduled Headways

**Goal:** Actual vs scheduled comparison (gap rate, avg wait time).

| Step | Task |
|------|------|
| 1 | GTFS import (`trips.txt`, `stop_times.txt`) |
| 2 | Compute scheduled headways per stop |
| 3 | API: actual vs scheduled side-by-side |
| 4 | Frontend: gap rate, avg wait |

* * *
## PR 5 — Segment Speed Analyzer

**Goal:** Avg bus speed by segment (same arrival primitives + stop sequence).

| Step | Task |
|------|------|
| 1 | `segment_speeds` table |
| 2 | `SpeedAnalyzer` in fan-out — pairs consecutive arrivals per vehicle |
| 3 | API + optional map overlay on `route_segments` |

* * *
## PR 6 — Ghost Bus Detection (separate PR)

**Goal:** Detect buses in schedule/predictions that never materialize.

Needs GTFS + prediction correlation. Larger scope — own PR.

* * *
## PR 7 — Reliability Metrics

**Goal:** Richer CV presentation, bunching (`headway > 2x scheduled`), time-window rollups UI.

* * *
## Refactor: Modular Analyzer Fan-Out

After ingestion is stable, refactor `PipelineRunner` to dispatch pings through a
`[]Analyzer` slice. Arrival detection becomes `ArrivalAnalyzer`; future analyzers
plug in without changing ingestion. Headway remains a **batch job**, not a ping analyzer.

* * *
# Architectural Constraints / Important Notes
*   Correctness is more important than parallelism initially.
*   Arrival detection is stateful and sensitive to ordering.
*   Vehicle state ownership must remain deterministic.
*   Do NOT parallelize arrival detection until necessary.
*   Each analyzer owns its persistence — keeps analysis modular and testable.
*   Headway rollup is batch/idempotent and separate from the live pipeline.
*   If scaling later:
    *   partition by `vehicle_id` hash
    *   each vehicle always routed to same worker
* * *
# Success Criteria

### PR 1 (ingestion)
1. Poll vehicle positions reliably for configured routes
2. Detect arrivals without duplicate spam
3. Persist arrivals and stops to Postgres

### PR 3 (headways)
4. Compute observed headways per stop/route/direction for a service date
5. Idempotent re-runs; job runs visible in admin
6. Summaries (mean, median, CV) over time windows via API

### Full v1 (through PR 4)
7. Compare actual vs scheduled frequency
8. Produce stable metrics over time windows

### Admin UI
9. Observe pipeline health, arrivals, and headway job runs via login-protected UI

* * *
# Pipeline Execution Model

How the pipeline process runs in production and locally.

| Option | Description | Pros | Cons | When to use |
|--------|-------------|------|------|-------------|
| **A. Goroutine in `main.go`** | Background goroutine alongside Gin HTTP server | One deployable, shared DB, simple local dev | API restarts kill ingestion; can't scale independently | **v1 default** — live vehicle polling |
| **B. Separate `cmd/pipeline` binary** | Second entry point, same repo | Clean separation; restart API without stopping ingestion | Two processes locally; shared config | When ingestion should survive API deploys |
| **C. External cron / HTTP trigger** | Batch job on a schedule (or admin-triggered) | Simple ops for batch work; idempotent re-runs | Not for live stateful polling | **Headway rollup, GTFS import** |
| **D. Message queue** | Pings published to Redis/SQS, workers consume | Horizontal scale, natural fan-out | Heavy infra; ordering guarantees are hard | Later, if polling all 127 routes |

**Live ingestion:** Option A (`PIPELINE_ENABLED=true`).
**Headway rollup:** Option C — daily GitHub Actions (or platform cron) → authenticated `POST /api/v1/admin/headways/run`.

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PIPELINE_ENABLED` | `false` | Start the ingestion pipeline on app boot |
| `PIPELINE_ROUTES` | *(unset)* | Comma-separated route IDs. When unset, routes are resolved from the latest ridership import |
| `PIPELINE_POLL_INTERVAL` | `30s` | How often to poll `getvehicles` |
| `PIPELINE_USE_FAKE_CTA` | `false` | Use simulated vehicle data (no API key needed) |
| `CTA_API_KEY` | — | Required when pipeline enabled and not using fake CTA |
| `LOG_LEVEL` | `info` | Set `debug` for arrival-detection diagnostics |
| `ADMIN_USERNAME` / `ADMIN_PASSWORD` | — | Admin UI login |
| `HEADWAY_JOB_TOKEN` | — | Bearer token for cron / automated headway runs |

* * *
# Admin UI

An admin UI makes the pipeline observable in production and easy to test locally.

### Goals
*   See whether the pipeline is running, when it last polled, and how many pings/arrivals were processed
*   Inspect recent arrivals and errors without querying the DB directly
*   Inspect headway job runs (status, counts, errors) and trigger a run for a service date

### Authentication
*   Login page protecting admin UI routes
*   **v1:** hardcoded credentials from env (`ADMIN_USERNAME`, `ADMIN_PASSWORD`)
*   Session cookie after login (HTTP-only)
*   **Job token:** `HEADWAY_JOB_TOKEN` accepted on headway run endpoint for cron (no session needed)

### Pages
1. **Login** ✅
2. **Dashboard** — pipeline status ✅
3. **Arrivals** — filters, stop names, time sort ✅
4. **Headway jobs** (PR 3) — run history + "Run for date"
5. **Stops** — optional later

### Admin API
*   `POST /api/v1/admin/login` / `logout` / `GET session` ✅
*   `GET /api/v1/admin/pipeline/status` ✅
*   `GET /api/v1/admin/arrivals` ✅
*   `POST /api/v1/admin/headways/run` — session **or** job token (PR 3)
*   `GET /api/v1/admin/headways/runs` — job run history (PR 3)

* * *
# Pipeline Architecture (detail)
## Stage 1 — Vehicle Fetchers
Poll CTA `getvehicles` per configured routes (batched, 10 per call), every 30–60 seconds.
Emit:

```yaml
vehicle_ping:
  vehicle_id
  route_id
  pattern_id   # pid — resolve direction via getpatterns
  direction    # usually empty from CTA; filled from pattern map
  lat
  lon
  timestamp
```

Fetching can be parallelized safely in a future version. v1 uses a synchronous poll loop.
* * *
## Stage 2 — Arrival Detector (ArrivalAnalyzer)
Consumes vehicle pings and maintains:

```css
map[vehicle_id]vehicle_state
```

**Direction:** resolve `pattern_id → rtdir` at startup from `getpatterns`; match stops only for that route+direction.

**SINGLE-WORKER / SINGLE-THREADED** for correctness. Emits `arrival_event` and persists to `arrivals`.
* * *
## Stage 3 — DB Writer
Initially embedded in ArrivalAnalyzer. Each future analyzer owns its own persistence.
* * *
## Stage 4 — Headway rollup (batch, not in poll loop)
Daily / on-demand job over a `service_date` → observed `headways` + `headway_job_runs`.
See **Headway Daily Job** above.
* * *
# Arrival Detection Strategy
## Initial implementation
A bus is considered to have arrived if:
*   vehicle is within X meters of stop (same route+direction as pattern)
*   AND vehicle has not already been recorded at this stop recently
Need cooldown logic:

```scheme
(vehicle_id, stop_id)
```

to prevent duplicate arrivals from:
*   repeated polling
*   GPS jitter
*   buses dwelling at stops
Possible cooldown:
*   2–5 minutes
* * *
## Future improvement
Instead of proximity-only:
*   detect when bus crosses/passes stop along route progression
This is more accurate and avoids false positives.
