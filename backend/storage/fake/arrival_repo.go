package fake

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// ArrivalRepo is an in-memory ArrivalRepository for local development and testing.
// It logs each arrival and holds them in a slice so they can be inspected in tests.
type ArrivalRepo struct {
	mu       sync.Mutex
	Arrivals []business.Arrival
}

func (r *ArrivalRepo) SaveArrival(_ context.Context, arrival business.Arrival) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Arrivals = append(r.Arrivals, arrival)
	log.Printf("fake arrival repo: saved arrival vehicle=%s route=%s dir=%s stop=%s at=%s",
		arrival.VehicleID, arrival.RouteID, arrival.Direction, arrival.StopID,
		arrival.Timestamp.Format("15:04:05"))
	return nil
}

func (r *ArrivalRepo) ListArrivals(_ context.Context, filter app.ArrivalFilter) ([]business.Arrival, error) {
	all := r.filteredArrivals(filter)
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	start := filter.Offset
	if start >= len(all) {
		return []business.Arrival{}, nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	out := make([]business.Arrival, end-start)
	copy(out, all[start:end])
	return out, nil
}

func (r *ArrivalRepo) CountArrivals(_ context.Context, filter app.ArrivalFilter) (int64, error) {
	return int64(len(r.filteredArrivals(filter))), nil
}

func (r *ArrivalRepo) ListArrivalsInRange(_ context.Context, start, end time.Time) ([]business.Arrival, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]business.Arrival, 0)
	for _, arrival := range r.Arrivals {
		if !arrival.Timestamp.Before(start) && arrival.Timestamp.Before(end) {
			out = append(out, arrival)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.RouteID != b.RouteID {
			return a.RouteID < b.RouteID
		}
		if a.Direction != b.Direction {
			return a.Direction < b.Direction
		}
		if a.StopID != b.StopID {
			return a.StopID < b.StopID
		}
		return a.Timestamp.Before(b.Timestamp)
	})
	return out, nil
}

func (r *ArrivalRepo) filteredArrivals(filter app.ArrivalFilter) []business.Arrival {
	r.mu.Lock()
	defer r.mu.Unlock()

	matches := make([]business.Arrival, 0, len(r.Arrivals))
	for _, arrival := range r.Arrivals {
		if filter.RouteID != "" && arrival.RouteID != filter.RouteID {
			continue
		}
		if filter.Direction != "" && arrival.Direction != filter.Direction {
			continue
		}
		if filter.VehicleID != "" && arrival.VehicleID != filter.VehicleID {
			continue
		}
		if filter.Stop != "" {
			stopQ := strings.ToLower(filter.Stop)
			idMatch := arrival.StopID == filter.Stop
			nameMatch := strings.Contains(strings.ToLower(arrival.StopName), stopQ)
			if !idMatch && !nameMatch {
				continue
			}
		}
		matches = append(matches, arrival)
	}

	sort.Slice(matches, func(i, j int) bool {
		if filter.SortAsc {
			return matches[i].Timestamp.Before(matches[j].Timestamp)
		}
		return matches[i].Timestamp.After(matches[j].Timestamp)
	})
	return matches
}

// All returns a snapshot of all saved arrivals (safe for concurrent use).
func (r *ArrivalRepo) All() []business.Arrival {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]business.Arrival, len(r.Arrivals))
	copy(out, r.Arrivals)
	return out
}
