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
Initialize Pipeline:
**Stages**
*   **Get vehicles:** Each worker calls endpoint every X seconds for the configured routes, then passes then the vehicle pings to the Arrival Detector stage
*   **Arrival Detector:** Arrival detector sends messages to save to db stage when appropriate. If this is stateful & has multiple workers I'll need to be careful about sending vehicle pings. I'll probably only make this multi-worker if I have to in order to keep this stage simple
*   **Save to db:** (could be lumped in with Arrival Detector, not sure what's better)

**NOTE:** A log could be useful for replaying & debugging

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
# Pipeline Architecture
## Stage 1 — Vehicle Fetchers
Multiple workers may:
*   poll CTA `getvehicles`
*   per configured route
*   every 30–60 seconds
Emit:

```sql
vehicle_ping:
  vehicle_id
  route_id
  direction
  lat
  lon
  timestamp


```

Fetching can be parallelized safely.
* * *
## Stage 2 — Arrival Detector
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

This stage should initially be SINGLE-WORKER / SINGLE-THREADED for correctness and simplicity.
Reason:
*   arrival detection is stateful
*   duplicate arrivals are dangerous
*   ordering matters
Arrival detector emits:

```plain
arrival_event


```

* * *
## Stage 3 — DB Writer
Persists:
*   arrivals
*   optionally raw vehicle pings
May initially be combined with Arrival Detector for simplicity.
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
## Phase 1 — Basic Infrastructure
### Step 1
Create DB tables:
*   stops
*   arrivals
*   headways
### Step 2
Build route + stop importer from CTA APIs.
### Step 3
Create vehicle polling worker:
*   configurable routes
*   emits vehicle pings
Goal:
*   verify stable ingestion
* * *
## Phase 2 — Arrival Detection MVP
### Step 4
Create in-memory vehicle state store:

```css
map[vehicleID]VehicleState


```

### Step 5
Implement nearest-stop lookup.
Initially:
*   simple distance calculation
*   nearest stop on route
### Step 6
Implement arrival detection:
*   within threshold distance
*   cooldown protection
### Step 7
Persist arrivals.
Goal:
*   reliably produce clean arrival events
* * *
## Phase 3 — Headway Computation
### Step 8
Create background job:
*   process unprocessed arrivals
*   grouped by stop + direction
### Step 9
Compute:

```plain
headway = arrival[i] - arrival[i-1]


```

### Step 10
Persist headways.
Goal:
*   produce actual measured service frequency
* * *
## Phase 4 — Scheduled Headways
### Step 11
Import GTFS schedule data.
### Step 12
Compute scheduled headways per stop/time window.
### Step 13
Expose:
*   actual vs scheduled headway
Goal:
*   first rider-facing reliability metric
* * *
## Phase 5 — Reliability Metrics
### Step 14
Compute:
*   mean headway
*   std deviation
*   coefficient of variation
### Step 15
Detect:
*   large gaps
*   bunching events
*   ghost-bus-like behavior
* * *
# Architectural Constraints / Important Notes
*   Correctness is more important than parallelism initially.
*   Arrival detection is stateful and sensitive to ordering.
*   Vehicle state ownership must remain deterministic.
*   Do NOT parallelize arrival detection until necessary.
*   If scaling later:
    *   partition by vehicle\_id hash
    *   each vehicle always routed to same worker
* * *
# Success Criteria for Initial PR
The system can:
1. Poll vehicle positions reliably
2. Detect arrivals without duplicate spam
3. Compute actual headways per stop
4. Compare actual vs scheduled frequency
5. Produce stable metrics over time windows