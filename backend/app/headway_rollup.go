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

// HeadwayRollup computes observed headways and daily summaries for a service date.
type HeadwayRollup struct {
	arrivals  ArrivalRepository
	headways  HeadwayRepository
	summaries HeadwaySummaryRepository
	runs      HeadwayJobRunRepository
}

func NewHeadwayRollup(
	arrivals ArrivalRepository,
	headways HeadwayRepository,
	summaries HeadwaySummaryRepository,
	runs HeadwayJobRunRepository,
) *HeadwayRollup {
	return &HeadwayRollup{arrivals: arrivals, headways: headways, summaries: summaries, runs: runs}
}

// RunResult is returned after a rollup completes (success or failure still updates the run row).
type RunResult struct {
	Run business.HeadwayJobRun
}

// Run recomputes observed headways and summaries for serviceDate. Idempotent for that day.
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

	result, runErr := r.execute(ctx, serviceDateOnly, start, end)
	finished := time.Now().UTC()
	run.FinishedAt = &finished
	run.ArrivalsProcessed = result.arrivalsProcessed
	run.HeadwaysWritten = result.headwaysWritten
	run.SummariesWritten = result.summariesWritten

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

	log.Printf("headway rollup: service_date=%s arrivals=%d headways=%d summaries=%d trigger=%s",
		start.Format("2006-01-02"), result.arrivalsProcessed, result.headwaysWritten, result.summariesWritten, triggeredBy)
	return RunResult{Run: run}, nil
}

type executeResult struct {
	arrivalsProcessed int
	headwaysWritten   int
	summariesWritten  int
}

func (r *HeadwayRollup) execute(ctx context.Context, serviceDate, start, end time.Time) (executeResult, error) {
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

	summaries := BuildPersistedSummaries(observed, serviceDate, start, end)
	if _, err := r.summaries.DeleteForServiceDate(ctx, serviceDate); err != nil {
		return executeResult{}, fmt.Errorf("delete existing summaries: %w", err)
	}
	if err := r.summaries.InsertBatch(ctx, summaries); err != nil {
		return executeResult{}, fmt.Errorf("insert summaries: %w", err)
	}

	return executeResult{
		arrivalsProcessed: len(arrivals),
		headwaysWritten:   len(observed),
		summariesWritten:  len(summaries),
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
			// Same vehicle twice is usually a re-detection, not a real bus gap.
			// The day's first arrival has no predecessor; the last has no successor —
			// both are naturally excluded because we only emit consecutive pairs.
			if prev.VehicleID == curr.VehicleID {
				continue
			}
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

// LoadStoredSummary builds the admin summary response from persisted rows for a service date.
// Falls back to false when no stop-grain rows exist (job not run yet).
func (r *HeadwayRollup) LoadStoredSummary(ctx context.Context, filter HeadwaySummaryFilter) (business.HeadwaySummaryStats, business.HeadwaySummaryStats, []business.HeadwayStopSummary, bool, error) {
	if r.summaries == nil {
		return business.HeadwaySummaryStats{}, business.HeadwaySummaryStats{}, nil, false, nil
	}

	stopFilter := filter
	stopFilter.Grain = business.HeadwayGrainStop
	stopFilter.Method = business.HeadwayMethodPooled
	stopFilter.Limit = 50_000
	stops, err := r.summaries.List(ctx, stopFilter)
	if err != nil {
		return business.HeadwaySummaryStats{}, business.HeadwaySummaryStats{}, nil, false, err
	}
	if len(stops) == 0 {
		return business.HeadwaySummaryStats{}, business.HeadwaySummaryStats{}, nil, false, nil
	}

	byStop := make([]business.HeadwayStopSummary, len(stops))
	for i, s := range stops {
		byStop[i] = business.HeadwayStopSummary{
			StopID:              s.StopID,
			StopName:            s.StopName,
			RouteID:             s.RouteID,
			RouteName:           s.RouteName,
			Direction:           s.Direction,
			HeadwaySummaryStats: s.HeadwaySummaryStats,
		}
	}

	pooled, equal := resolveStoredOverall(ctx, r.summaries, filter, byStop)
	return pooled, equal, byStop, true, nil
}

func resolveStoredOverall(
	ctx context.Context,
	repo HeadwaySummaryRepository,
	filter HeadwaySummaryFilter,
	byStop []business.HeadwayStopSummary,
) (pooled, equal business.HeadwaySummaryStats) {
	equal = MeanOfStopMeans(byStop)

	if filter.StopID != "" || filter.Stop != "" {
		if len(byStop) == 1 {
			return byStop[0].HeadwaySummaryStats, byStop[0].HeadwaySummaryStats
		}
		return SummarizeHeadwaysFromStopStats(byStop), equal
	}

	if filter.RouteID != "" && filter.Direction != "" {
		rows, err := repo.List(ctx, HeadwaySummaryFilter{
			ServiceDate: filter.ServiceDate,
			Grain:       business.HeadwayGrainRouteDirection,
			RouteID:     filter.RouteID,
			Direction:   filter.Direction,
			Limit:       50,
		})
		if err == nil {
			for _, row := range rows {
				switch row.Method {
				case business.HeadwayMethodPooled:
					pooled = row.HeadwaySummaryStats
				case business.HeadwayMethodEqualStop:
					equal = row.HeadwaySummaryStats
				}
			}
			if pooled.Count > 0 {
				return pooled, equal
			}
		}
	}

	if filter.RouteID == "" && filter.Direction == "" && filter.Stop == "" && filter.StopID == "" {
		rows, err := repo.List(ctx, HeadwaySummaryFilter{
			ServiceDate: filter.ServiceDate,
			Grain:       business.HeadwayGrainServiceDay,
			Limit:       50,
		})
		if err == nil {
			for _, row := range rows {
				switch row.Method {
				case business.HeadwayMethodPooled:
					pooled = row.HeadwaySummaryStats
				case business.HeadwayMethodEqualStop:
					equal = row.HeadwaySummaryStats
				}
			}
			if pooled.Count > 0 {
				return pooled, equal
			}
		}
	}

	return SummarizeHeadwaysFromStopStats(byStop), equal
}

// SummarizeHeadwaysFromStopStats approximates pooled stats when only stop aggregates exist.
// Uses count-weighted mean of stop means (not a true recompute over raw gaps).
func SummarizeHeadwaysFromStopStats(stops []business.HeadwayStopSummary) business.HeadwaySummaryStats {
	if len(stops) == 0 {
		return business.HeadwaySummaryStats{}
	}
	var totalCount int
	var weightedMean float64
	for _, s := range stops {
		totalCount += s.Count
		weightedMean += s.MeanMinutes * float64(s.Count)
	}
	if totalCount == 0 {
		return business.HeadwaySummaryStats{}
	}
	mean := weightedMean / float64(totalCount)
	return business.HeadwaySummaryStats{
		Count:          totalCount,
		MeanMinutes:    mean,
		MedianMinutes:  MeanOfStopMeans(stops).MedianMinutes,
		StdDevMinutes:  MeanOfStopMeans(stops).StdDevMinutes,
		CV:             MeanOfStopMeans(stops).CV,
		AvgWaitMinutes: mean / 2,
	}
}
