package app

import (
	"context"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// CTAVehicleClient is the port for fetching live vehicle positions from the CTA.
// Implementations may call the real CTA Bus Tracker API or return fake data for local development.
type CTAVehicleClient interface {
	// GetVehicles returns the latest vehicle pings for up to 10 route IDs (CTA API limit).
	GetVehicles(ctx context.Context, routeIDs []string) ([]business.VehiclePing, error)

	// GetStops returns all stops for a given route and direction, ordered by sequence.
	GetStops(ctx context.Context, routeID, direction string) ([]business.Stop, error)
}

// ArrivalFilter controls paginated arrival queries for the admin UI.
type ArrivalFilter struct {
	RouteID   string
	Direction string
	Limit     int
	Offset    int
}

// ArrivalRepository is the port for persisting detected arrival events.
type ArrivalRepository interface {
	SaveArrival(ctx context.Context, arrival business.Arrival) error
	ListArrivals(ctx context.Context, filter ArrivalFilter) ([]business.Arrival, error)
	CountArrivals(ctx context.Context, filter ArrivalFilter) (int64, error)
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
