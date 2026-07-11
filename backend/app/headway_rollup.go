package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// Chicago is the service-day timezone for CTA operations.
func ChicagoLocation() *time.Location {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return time.FixedZone("CST", -6*60*60)
	}
	return loc
}

// ServiceDateBounds returns [start, end) for a calendar date in America/Chicago.
func ServiceDateBounds(serviceDate time.Time) (start, end time.Time) {
	loc := ChicagoLocation()
	y, m, d := serviceDate.In(loc).Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, loc)
	end = start.Add(24 * time.Hour)
	return start, end
}

// YesterdayServiceDate returns yesterday's calendar date in Chicago (midnight local).
func YesterdayServiceDate() time.Time {
	now := time.Now().In(ChicagoLocation())
	y, m, d := now.AddDate(0, 0, -1).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, ChicagoLocation())
}

// ParseServiceDate parses YYYY-MM-DD as a Chicago calendar date.
func ParseServiceDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", s, ChicagoLocation())
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid service_date %q (want YYYY-MM-DD): %w", s, err)
	}
	return t, nil
}

// HeadwayRollup computes observed headways for a service date.
type HeadwayRollup struct {
	arrivals ArrivalRepository
	headways HeadwayRepository
	runs     HeadwayJobRunRepository
}

func NewHeadwayRollup(arrivals ArrivalRepository, headways HeadwayRepository, runs HeadwayJobRunRepository) *HeadwayRollup {
	return &HeadwayRollup{arrivals: arrivals, headways: headways, runs: runs}
}

// RunResult is returned after a rollup completes (success or failure still updates the run row).
type RunResult struct {
	Run business.HeadwayJobRun
}

// Run recomputes observed headways for serviceDate. Idempotent for that day.
func (r *HeadwayRollup) Run(ctx context.Context, serviceDate time.Time, triggeredBy business.HeadwayJobTrigger) (RunResult, error) {
	start, end := ServiceDateBounds(serviceDate)
	serviceDateOnly := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	run, err := r.runs.Create(ctx, business.HeadwayJobRun{
		ServiceDate: serviceDateOnly,
		Status:      business.HeadwayJobRunning,
		TriggeredBy: triggeredBy,
		StartedAt:   time.Now().UTC(),
	})
	if err != nil {
		return RunResult{}, fmt.Errorf("create job run: %w", err)
	}

	result, runErr := r.execute(ctx, start, end)
	finished := time.Now().UTC()
	run.FinishedAt = &finished
	run.ArrivalsProcessed = result.arrivalsProcessed
	run.HeadwaysWritten = result.headwaysWritten

	if runErr != nil {
		run.Status = business.HeadwayJobFailed
		run.ErrorMessage = runErr.Error()
		_ = r.runs.Update(ctx, run)
		return RunResult{Run: run}, runErr
	}

	run.Status = business.HeadwayJobSuccess
	if err := r.runs.Update(ctx, run); err != nil {
		return RunResult{Run: run}, fmt.Errorf("update job run: %w", err)
	}

	log.Printf("headway rollup: service_date=%s arrivals=%d headways=%d trigger=%s",
		start.Format("2006-01-02"), result.arrivalsProcessed, result.headwaysWritten, triggeredBy)
	return RunResult{Run: run}, nil
}

type executeResult struct {
	arrivalsProcessed int
	headwaysWritten   int
}

func (r *HeadwayRollup) execute(ctx context.Context, start, end time.Time) (executeResult, error) {
	arrivals, err := r.arrivals.ListArrivalsInRange(ctx, start, end)
	if err != nil {
		return executeResult{}, fmt.Errorf("list arrivals: %w", err)
	}

	observed := ComputeObservedHeadways(arrivals)

	if _, err := r.headways.DeleteInRange(ctx, start, end); err != nil {
		return executeResult{}, fmt.Errorf("delete existing headways: %w", err)
	}
	if err := r.headways.InsertBatch(ctx, observed); err != nil {
		return executeResult{}, fmt.Errorf("insert headways: %w", err)
	}

	return executeResult{
		arrivalsProcessed: len(arrivals),
		headwaysWritten:   len(observed),
	}, nil
}

// ComputeObservedHeadways groups arrivals by stop/route/direction and emits consecutive gaps.
func ComputeObservedHeadways(arrivals []business.Arrival) []business.Headway {
	type key struct {
		StopID, RouteID, Direction string
	}

	groups := make(map[key][]business.Arrival)
	order := make([]key, 0)
	for _, a := range arrivals {
		k := key{StopID: a.StopID, RouteID: a.RouteID, Direction: a.Direction}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], a)
	}

	out := make([]business.Headway, 0)
	for _, k := range order {
		list := groups[k]
		// ListArrivalsInRange is ordered globally; within a group timestamps should already be ascending.
		for i := 1; i < len(list); i++ {
			prev, curr := list[i-1], list[i]
			mins := curr.Timestamp.Sub(prev.Timestamp).Minutes()
			if mins <= 0 {
				continue
			}
			out = append(out, business.Headway{
				StopID:         k.StopID,
				RouteID:        k.RouteID,
				Direction:      k.Direction,
				Timestamp:      curr.Timestamp,
				HeadwayMinutes: mins,
				FromVehicleID:  prev.VehicleID,
				ToVehicleID:    curr.VehicleID,
			})
		}
	}
	return out
}

// ListJobRuns returns recent headway job runs for the admin UI.
func (r *HeadwayRollup) ListJobRuns(ctx context.Context, limit, offset int) ([]business.HeadwayJobRun, error) {
	if limit <= 0 {
		limit = 50
	}
	return r.runs.List(ctx, limit, offset)
}
