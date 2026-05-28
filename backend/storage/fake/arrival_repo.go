package fake

import (
	"context"
	"log"
	"sync"

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

// All returns a snapshot of all saved arrivals (safe for concurrent use).
func (r *ArrivalRepo) All() []business.Arrival {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]business.Arrival, len(r.Arrivals))
	copy(out, r.Arrivals)
	return out
}
