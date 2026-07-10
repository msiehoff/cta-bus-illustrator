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
## headway
Represents the gap between two consecutive arrivals at a stop.

```sql
headway:
  stop_id
  route_id
  direction
  timestamp
  headway_minutes


```

Timestamp should represent the later arrival.
* * *
# Pipeline Execution Model

How the pipeline process runs in production and locally.

| Option | Description | Pros | Cons | When to use |
|--------|-------------|------|------|-------------|
| **A. Goroutine in `main.go`** | Background goroutine alongside Gin HTTP server | One deployable, shared DB, simple local dev | API restarts kill ingestion; can't scale independently | **v1 default** — single-container deploy |
| **B. Separate `cmd/pipeline` binary** | Second entry point, same repo | Clean separation; restart API without stopping ingestion | Two processes locally; shared config | When ingestion should survive API deploys |
| **C. External cron** | Batch job runs on schedule, exits | Simple ops for batch work | **Bad for live polling** — loses in-memory vehicle state between runs | GTFS import, nightly rollups only |
| **D. Message queue** | Pings published to Redis/SQS, workers consume | Horizontal scale, natural fan-out | Heavy infra; ordering guarantees are hard | Later, if polling all 127 routes |

**v1 approach:** Option A, gated by `PIPELINE_ENABLED=true`. Extract to Option B when deploying API + worker as separate services.

### Environment variables (pipeline)

| Variable | Default | Description |
|----------|---------|-------------|
| `PIPELINE_ENABLED` | `false` | Start the ingestion pipeline on app boot |
| `PIPELINE_ROUTES` | *(unset)* | Comma-separated route IDs. When unset, routes are resolved from the latest ridership import |
| `PIPELINE_POLL_INTERVAL` | `30s` | How often to poll `getvehicles` |
| `PIPELINE_USE_FAKE_CTA` | `false` | Use simulated vehicle data (no API key needed) |
| `CTA_API_KEY` | — | Required when pipeline enabled and not using fake CTA |

* * *
# Admin UI

An admin UI makes the pipeline observable in production and easy to test locally.

### Goals
*   See whether the pipeline is running, when it last polled, and how many pings/arrivals were processed
*   Inspect recent arrivals and errors without querying the DB directly
*   Manually trigger stop reload or view configured routes

### Authentication
*   Login page protecting all `/admin/*` routes
*   **v1:** hardcoded credentials from environment variables (no user DB)
    *   `ADMIN_USERNAME` — required to access admin UI
    *   `ADMIN_PASSWORD` — required to access admin UI
*   Session cookie after successful login (HTTP-only, signed with `ADMIN_SESSION_SECRET` or derived from password)
*   Future: OAuth or proper user management if needed

### Planned pages (PR 2 — Admin UI)
1. **Login** — username/password form
2. **Dashboard** — pipeline status (running, last poll time, routes, ping count, arrival count, last error)
3. **Arrivals** — paginated table of recent arrivals, filterable by route/direction
4. **Stops** — loaded stops per route (verify importer)

### Admin API endpoints (backend, PR 2)
*   `POST /api/v1/admin/login` — authenticate, set session cookie
*   `POST /api/v1/admin/logout`
*   `GET /api/v1/admin/session` — check auth state (public)
*   `GET /api/v1/admin/pipeline/status` — pipeline health + stats
*   `GET /api/v1/admin/arrivals` — recent arrivals (paginated)

Admin UI routes (`/admin/login`, `/admin`, `/admin/arrivals`) are React pages. API calls go under `/api/v1/admin/*` so the Vite dev proxy does not intercept page navigation.

PR 1 wires the pipeline and persists data; PR 2 adds the admin API + React admin pages.

* * *
# Pipeline Architecture (detail)
## Stage 1 — Vehicle Fetchers
Poll CTA `getvehicles` per configured routes (batched, 10 per call), every 30–60 seconds.
Emit:

```yaml
vehicle_ping:
  vehicle_id
  route_id
  direction
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

Vehicle state:

```yaml
vehicle_state:
  last_stop_id
  last_arrival_time
  last_position
  last_timestamp
```

**SINGLE-WORKER / SINGLE-THREADED** for correctness. Emits `arrival_event` and persists to `arrivals`.
* * *
## Stage 3 — DB Writer
Initially embedded in ArrivalAnalyzer. Each future analyzer owns its own persistence.
* * *
# Arrival Detection Strategy
## Initial implementation
A bus is considered to have arrived if:
*   vehicle is within X meters of stop
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
* * *
# Headway Calculation
For each:

```css
(stop_id, route_id, direction)


```

Sort arrivals by timestamp:

```erlang
t1, t2, t3, ...


```

Compute:

```plain
headway_i = t_i - t_(i-1)


```

Store as `headway`.
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
# Important Metrics (Future)
## Reliability / Bunching
Possible metrics:

```erlang
stddev(headway) / mean(headway)


```

or:

```erlang
% of headways > 2x scheduled


```

* * *
## Rider Wait Time
For irregular service:

```plain
avg_wait ≈ mean(headway^2) / (2 * mean(headway))


```

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

## PR 1 — Wire Up Ingestion ✅ (in progress)

**Goal:** Pipeline runs locally and in deploy; persists arrivals (and stops) to Postgres.

| Step | Task | Status |
|------|------|--------|
| 1 | `stops` + `arrivals` tables | ✅ `000007` |
| 2 | `headways` table migration | ✅ `000008` |
| 3 | Postgres `ArrivalRepository` + `StopRepository` | ✅ |
| 4 | Real `CTAVehicleClient` adapter over `storage/cta` | ✅ |
| 5 | Wire pipeline into `main.go` (`PIPELINE_ENABLED`, env config) | ✅ |
| 6 | Verify with fake CTA client → rows in `arrivals` | manual |

**Success criteria:** Stable ingestion for 2 routes over a few hours, no duplicate spam.

* * *
## PR 2 — Admin UI + Pipeline Observability ✅

**Goal:** Login-protected admin UI to monitor pipeline health and inspect arrivals locally and in production.

| Step | Task | Status |
|------|------|--------|
| 1 | `PipelineStatus` tracker | ✅ |
| 2 | Admin auth middleware (`ADMIN_USERNAME`, `ADMIN_PASSWORD`, session cookie) | ✅ |
| 3 | Admin API: login, logout, pipeline status, recent arrivals | ✅ |
| 4 | React admin pages: login, dashboard, arrivals table | ✅ |

* * *
## PR 3 — Headway Computation + Rider API

**Goal:** Actual headways per stop/direction.

| Step | Task |
|------|------|
| 1 | Headway rollup job (ticker or post-poll trigger) |
| 2 | Query: sort arrivals by `(stop_id, route_id, direction)`, compute `t[i] - t[i-1]` |
| 3 | `GET /routes/:id/headways` endpoint |
| 4 | Frontend chart on RoutePage (replace "Coming soon") |

* * *
## PR 4 — GTFS Scheduled Headways

**Goal:** Actual vs scheduled comparison (gap rate, avg wait time).

| Step | Task |
|------|------|
| 1 | GTFS import (`trips.txt`, `stop_times.txt`) |
| 2 | Compute scheduled headways per stop |
| 3 | API: actual vs scheduled side-by-side |
| 4 | Frontend: gap rate, avg wait (`headway / 2`) |

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

**Goal:** CV, bunching (`headway > 2x scheduled`), time-window rollups.

* * *
## Refactor: Modular Analyzer Fan-Out

After PR 1 stabilizes ingestion, refactor `PipelineRunner` to dispatch pings through a `[]Analyzer` slice. Arrival detection becomes `ArrivalAnalyzer`; future analyzers plug in without changing ingestion.

* * *
# Architectural Constraints / Important Notes
*   Correctness is more important than parallelism initially.
*   Arrival detection is stateful and sensitive to ordering.
*   Vehicle state ownership must remain deterministic.
*   Do NOT parallelize arrival detection until necessary.
*   Each analyzer owns its persistence — keeps analysis modular and testable.
*   If scaling later:
    *   partition by `vehicle_id` hash
    *   each vehicle always routed to same worker
* * *
# Success Criteria

### PR 1 (ingestion)
1. Poll vehicle positions reliably for configured routes
2. Detect arrivals without duplicate spam
3. Persist arrivals and stops to Postgres

### Full v1 (through PR 4)
4. Compute actual headways per stop
5. Compare actual vs scheduled frequency
6. Produce stable metrics over time windows

### Admin UI (PR 2)
7. Observe pipeline health and recent arrivals via login-protected UI