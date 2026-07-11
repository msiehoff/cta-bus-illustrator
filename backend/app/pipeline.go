package app

import (
	"context"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// CTAVehicleClient is the port for fetching live vehicle positions from the CTA.
// Implementations may call the real CTA Bus Tracker API or return fake data for local development.
type CTAVehicleClient interface {
	// GetVehicles returns the latest vehicle pings for up to 10 route IDs (CTA API limit).
	GetVehicles(ctx context.Context, routeIDs []string) ([]business.VehiclePing, error)

	// GetStops returns all stops for a given route and direction, ordered by sequence.
	GetStops(ctx context.Context, routeID, direction string) ([]business.Stop, error)

	// GetPatterns returns a map of pattern ID → normalized direction (e.g. "Eastbound")
	// for all patterns on a route.
	GetPatterns(ctx context.Context, routeID string) (map[int]string, error)
}

// ArrivalFilter controls paginated arrival queries for the admin UI.
type ArrivalFilter struct {
	RouteID   string
	Direction string
	// Stop matches stop_id exactly or stop name (case-insensitive substring).
	Stop string
	// VehicleID filters by vehicle_id (exact match).
	VehicleID string
	From      *time.Time
	To        *time.Time
	// SortAsc sorts by timestamp ascending when true; default is newest first.
	SortAsc bool
	Limit   int
	Offset  int
}

// ArrivalRepository is the port for persisting detected arrival events.
type ArrivalRepository interface {
	SaveArrival(ctx context.Context, arrival business.Arrival) error
	ListArrivals(ctx context.Context, filter ArrivalFilter) ([]business.Arrival, error)
	CountArrivals(ctx context.Context, filter ArrivalFilter) (int64, error)
	// ListArrivalsInRange returns all arrivals with timestamp in [start, end), ordered for headway computation.
	ListArrivalsInRange(ctx context.Context, start, end time.Time) ([]business.Arrival, error)
}

// StopRepository is the port for persisting stop metadata loaded from the CTA API.
type StopRepository interface {
	UpsertStops(ctx context.Context, stops []business.Stop) error
}

// PipelineRouteProvider supplies route IDs for the pipeline when PIPELINE_ROUTES is unset.
// Production implementations typically derive routes from the latest ridership import.
type PipelineRouteProvider interface {
	GetRouteIDs(ctx context.Context) ([]string, error)
}
