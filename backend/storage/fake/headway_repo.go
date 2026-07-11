package fake

import (
	"context"
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

func (r *HeadwayRepo) ListInRange(_ context.Context, start, end time.Time, filter app.HeadwayListFilter) ([]business.Headway, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]business.Headway, 0)
	for _, h := range r.Headways {
		if h.Timestamp.Before(start) || !h.Timestamp.Before(end) {
			continue
		}
		if filter.RouteID != "" && h.RouteID != filter.RouteID {
			continue
		}
		if filter.Direction != "" && h.Direction != filter.Direction {
			continue
		}
		if filter.StopID != "" && h.StopID != filter.StopID {
			continue
		}
		out = append(out, h)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	if filter.Offset >= len(out) {
		return []business.Headway{}, nil
	}
	endIdx := filter.Offset + limit
	if endIdx > len(out) {
		endIdx = len(out)
	}
	return out[filter.Offset:endIdx], nil
}

func (r *HeadwayRepo) CountInRange(_ context.Context, start, end time.Time, filter app.HeadwayListFilter) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int64
	for _, h := range r.Headways {
		if h.Timestamp.Before(start) || !h.Timestamp.Before(end) {
			continue
		}
		if filter.RouteID != "" && h.RouteID != filter.RouteID {
			continue
		}
		if filter.Direction != "" && h.Direction != filter.Direction {
			continue
		}
		if filter.StopID != "" && h.StopID != filter.StopID {
			continue
		}
		count++
	}
	return count, nil
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
	// Newest first (append order is chronological).
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
