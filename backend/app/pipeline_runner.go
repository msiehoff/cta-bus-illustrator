package app

import (
	"context"
	"log"
	"time"
)

// PipelineConfig holds the runtime knobs for the arrival detection pipeline.
type PipelineConfig struct {
	// RouteIDs is the set of CTA route IDs to monitor (e.g. ["8", "66", "77"]).
	RouteIDs []string

	// PollInterval controls how often vehicle positions are fetched from the CTA API.
	// The CTA updates vehicle positions roughly every 30 seconds.
	PollInterval time.Duration

	// Directions is the set of direction strings to load stops for (e.g. ["Northbound", "Southbound"]).
	// If empty, only "Northbound" and "Southbound" are used as defaults.
	Directions []string
}

// PipelineRunner wires together the three pipeline stages:
//
//	Stage 1 — Vehicle fetcher: polls CTAVehicleClient for vehicle positions.
//	Stage 2 — Arrival detector: stateful single-worker arrival detection.
//	Stage 3 — DB writer: embedded in the ArrivalDetector via ArrivalRepository.
type PipelineRunner struct {
	client   CTAVehicleClient
	detector *ArrivalDetector
	stops    StopRepository
	cfg      PipelineConfig
	status   *pipelineStatusTracker
}

// NewPipelineRunner creates a runner wired to the given client and arrival repo.
// stops may be nil; when set, loaded stops are upserted to the database.
func NewPipelineRunner(client CTAVehicleClient, repo ArrivalRepository, stops StopRepository, cfg PipelineConfig) *PipelineRunner {
	if len(cfg.Directions) == 0 {
		cfg.Directions = []string{"Northbound", "Southbound", "Eastbound", "Westbound"}
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 30 * time.Second
	}
	return &PipelineRunner{
		client:   client,
		detector: NewArrivalDetector(repo),
		stops:    stops,
		cfg:      cfg,
		status:   newPipelineStatusTracker(cfg),
	}
}

// Status returns a snapshot of pipeline health for observability.
func (r *PipelineRunner) Status() PipelineStatus {
	return r.status.snapshot()
}

// Run starts the pipeline and blocks until ctx is cancelled.
// Stop loading runs once at startup; vehicle polling then loops on PollInterval.
func (r *PipelineRunner) Run(ctx context.Context) error {
	log.Printf("pipeline: starting — routes=%v poll_interval=%v", r.cfg.RouteIDs, r.cfg.PollInterval)
	r.status.setRunning(true)
	defer r.status.setRunning(false)

	if err := r.loadAllStops(ctx); err != nil {
		// Non-fatal: log and continue. Pings for routes with missing stops are skipped.
		log.Printf("pipeline: stop loading incomplete: %v", err)
	}

	ticker := time.NewTicker(r.cfg.PollInterval)
	defer ticker.Stop()

	// Poll immediately on start, then on each tick.
	r.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("pipeline: shutting down")
			return ctx.Err()
		case <-ticker.C:
			r.poll(ctx)
		}
	}
}

// poll fetches vehicle pings for all configured routes and feeds them through the detector.
// The CTA API accepts up to 10 route IDs per call, so routes are batched automatically.
func (r *PipelineRunner) poll(ctx context.Context) {
	var totalPings int
	var pollErr error

	batches := chunkStrings(r.cfg.RouteIDs, 10)
	for _, batch := range batches {
		pings, err := r.client.GetVehicles(ctx, batch)
		if err != nil {
			log.Printf("pipeline: getvehicles error routes=%v: %v", batch, err)
			pollErr = err
			continue
		}
		log.Printf("pipeline: received %d pings for routes %v", len(pings), batch)
		totalPings += len(pings)
		if DebugEnabled() && len(pings) > 0 {
			sample := pings
			if len(sample) > 3 {
				sample = sample[:3]
			}
			for _, ping := range sample {
				Debugf("pipeline: sample ping vehicle=%s route=%s dir=%q lat=%.5f lon=%.5f ts=%s",
					ping.VehicleID, ping.RouteID, ping.Direction, ping.Lat, ping.Lon,
					ping.Timestamp.Format("15:04:05"))
			}
		}
		for _, ping := range pings {
			r.detector.ProcessPing(ctx, ping)
		}
	}

	r.status.recordPoll(totalPings, pollErr)
}

// loadAllStops fetches and caches stops for every configured route+direction combination.
func (r *PipelineRunner) loadAllStops(ctx context.Context) error {
	for _, routeID := range r.cfg.RouteIDs {
		for _, dir := range r.cfg.Directions {
			stops, err := r.client.GetStops(ctx, routeID, dir)
			if err != nil {
				log.Printf("pipeline: getstops route=%s dir=%s: %v", routeID, dir, err)
				continue
			}
			r.detector.LoadStops(routeID, dir, stops)
			if r.stops != nil {
				if err := r.stops.UpsertStops(ctx, stops); err != nil {
					log.Printf("pipeline: upsert stops route=%s dir=%s: %v", routeID, dir, err)
				}
			}
			log.Printf("pipeline: loaded %d stops route=%s dir=%s", len(stops), routeID, dir)
		}
	}
	return nil
}

// chunkStrings splits a slice into sub-slices of at most size n.
func chunkStrings(s []string, n int) [][]string {
	var chunks [][]string
	for len(s) > 0 {
		if len(s) < n {
			n = len(s)
		}
		chunks = append(chunks, s[:n])
		s = s[n:]
	}
	return chunks
}
