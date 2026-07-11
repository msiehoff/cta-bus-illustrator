package fake

import (
	"context"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// HeadwayRepo is an in-memory HeadwayRepository for local development and testing.
type HeadwayRepo struct {
	mu       sync.Mutex
	Headways []business.Headway
}

func (r *HeadwayRepo) DeleteInRange(_ context.Context, start, end time.Time) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	kept := make([]business.Headway, 0, len(r.Headways))
	var deleted int64
	for _, h := range r.Headways {
		if !h.Timestamp.Before(start) && h.Timestamp.Before(end) {
			deleted++
			continue
		}
		kept = append(kept, h)
	}
	r.Headways = kept
	return deleted, nil
}

func (r *HeadwayRepo) InsertBatch(_ context.Context, headways []business.Headway) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Headways = append(r.Headways, headways...)
	return nil
}

func (r *HeadwayRepo) List(_ context.Context, filter app.HeadwayListFilter) ([]business.Headway, error) {
	matched := r.filtered(filter)
	sort.SliceStable(matched, func(i, j int) bool {
		if filter.SortAsc {
			return matched[i].Timestamp.Before(matched[j].Timestamp)
		}
		return matched[j].Timestamp.Before(matched[i].Timestamp)
	})

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if filter.Offset >= len(matched) {
		return []business.Headway{}, nil
	}
	end := filter.Offset + limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[filter.Offset:end], nil
}

func (r *HeadwayRepo) ListAll(_ context.Context, filter app.HeadwayListFilter) ([]business.Headway, error) {
	matched := r.filtered(filter)
	sort.SliceStable(matched, func(i, j int) bool {
		if filter.SortAsc {
			return matched[i].Timestamp.Before(matched[j].Timestamp)
		}
		return matched[j].Timestamp.Before(matched[i].Timestamp)
	})
	if len(matched) > 50_000 {
		matched = matched[:50_000]
	}
	return matched, nil
}

func (r *HeadwayRepo) Count(_ context.Context, filter app.HeadwayListFilter) (int64, error) {
	return int64(len(r.filtered(filter))), nil
}

func (r *HeadwayRepo) filtered(filter app.HeadwayListFilter) []business.Headway {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]business.Headway, 0)
	stopQ := strings.ToLower(filter.Stop)
	for _, h := range r.Headways {
		if filter.RouteID != "" && h.RouteID != filter.RouteID {
			continue
		}
		if filter.Direction != "" && h.Direction != filter.Direction {
			continue
		}
		if filter.VehicleID != "" && h.FromVehicleID != filter.VehicleID && h.ToVehicleID != filter.VehicleID {
			continue
		}
		if filter.Stop != "" {
			idMatch := h.StopID == filter.Stop
			nameMatch := strings.Contains(strings.ToLower(h.StopName), stopQ)
			if !idMatch && !nameMatch {
				continue
			}
		}
		if filter.From != nil && h.Timestamp.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !h.Timestamp.Before(*filter.To) {
			continue
		}
		out = append(out, h)
	}
	return out
}

// HeadwayJobRunRepo is an in-memory HeadwayJobRunRepository.
type HeadwayJobRunRepo struct {
	mu   sync.Mutex
	seq  atomic.Int64
	Runs []business.HeadwayJobRun
}

func (r *HeadwayJobRunRepo) Create(_ context.Context, run business.HeadwayJobRun) (business.HeadwayJobRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	run.ID = r.seq.Add(1)
	r.Runs = append(r.Runs, run)
	return run, nil
}

func (r *HeadwayJobRunRepo) Update(_ context.Context, run business.HeadwayJobRun) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.Runs {
		if existing.ID == run.ID {
			r.Runs[i] = run
			return nil
		}
	}
	return nil
}

func (r *HeadwayJobRunRepo) List(_ context.Context, limit, offset int) ([]business.HeadwayJobRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limit <= 0 {
		limit = 50
	}
	n := len(r.Runs)
	if offset >= n {
		return []business.HeadwayJobRun{}, nil
	}
	out := make([]business.HeadwayJobRun, 0, limit)
	for i := n - 1 - offset; i >= 0 && len(out) < limit; i-- {
		out = append(out, r.Runs[i])
	}
	return out, nil
}

func (r *HeadwayJobRunRepo) Get(_ context.Context, id int64) (business.HeadwayJobRun, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, run := range r.Runs {
		if run.ID == id {
			return run, nil
		}
	}
	return business.HeadwayJobRun{}, nil
}

// HeadwaySummaryRepo is an in-memory HeadwaySummaryRepository.
type HeadwaySummaryRepo struct {
	mu        sync.Mutex
	Summaries []business.HeadwaySummary
}

func (r *HeadwaySummaryRepo) DeleteForServiceDate(_ context.Context, serviceDate time.Time) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	day := serviceDate.Format("2006-01-02")
	kept := make([]business.HeadwaySummary, 0, len(r.Summaries))
	var deleted int64
	for _, s := range r.Summaries {
		if s.ServiceDate.Format("2006-01-02") == day {
			deleted++
			continue
		}
		kept = append(kept, s)
	}
	r.Summaries = kept
	return deleted, nil
}

func (r *HeadwaySummaryRepo) InsertBatch(_ context.Context, summaries []business.HeadwaySummary) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Summaries = append(r.Summaries, summaries...)
	return nil
}

func (r *HeadwaySummaryRepo) List(_ context.Context, filter app.HeadwaySummaryFilter) ([]business.HeadwaySummary, error) {
	matched := r.filtered(filter)
	sort.SliceStable(matched, func(i, j int) bool {
		if matched[i].ServiceDate.Equal(matched[j].ServiceDate) {
			if filter.SortAsc {
				return matched[i].MeanMinutes < matched[j].MeanMinutes
			}
			return matched[i].MeanMinutes > matched[j].MeanMinutes
		}
		if filter.SortAsc {
			return matched[i].ServiceDate.Before(matched[j].ServiceDate)
		}
		return matched[i].ServiceDate.After(matched[j].ServiceDate)
	})

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if filter.Offset >= len(matched) {
		return []business.HeadwaySummary{}, nil
	}
	end := filter.Offset + limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[filter.Offset:end], nil
}

func (r *HeadwaySummaryRepo) Count(_ context.Context, filter app.HeadwaySummaryFilter) (int64, error) {
	return int64(len(r.filtered(filter))), nil
}

func (r *HeadwaySummaryRepo) filtered(filter app.HeadwaySummaryFilter) []business.HeadwaySummary {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]business.HeadwaySummary, 0)
	for _, s := range r.Summaries {
		if !filter.ServiceDate.IsZero() && s.ServiceDate.Format("2006-01-02") != filter.ServiceDate.Format("2006-01-02") {
			continue
		}
		if filter.Grain != "" && s.Grain != filter.Grain {
			continue
		}
		if filter.Method != "" && s.Method != filter.Method {
			continue
		}
		if filter.RouteID != "" && s.RouteID != filter.RouteID {
			continue
		}
		if filter.Direction != "" && s.Direction != filter.Direction {
			continue
		}
		if filter.StopID != "" && s.StopID != filter.StopID {
			continue
		}
		if filter.Stop != "" {
			stopQ := strings.ToLower(filter.Stop)
			if s.StopID != filter.Stop && !strings.Contains(strings.ToLower(s.StopName), stopQ) {
				continue
			}
		}
		out = append(out, s)
	}
	return out
}
