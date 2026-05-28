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

// ArrivalRepository is the port for persisting detected arrival events.
type ArrivalRepository interface {
	SaveArrival(ctx context.Context, arrival business.Arrival) error
}
