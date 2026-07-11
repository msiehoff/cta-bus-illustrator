package business

import "time"

// VehiclePing is emitted by the vehicle fetcher stage for each polled vehicle position.
type VehiclePing struct {
	VehicleID string
	RouteID   string
	// PatternID is the CTA pattern ID (pid) for the trip the vehicle is running.
	// Used to resolve Direction via getpatterns.
	PatternID int
	// Direction is the route direction (e.g. "Northbound"). Prefer resolving from
	// PatternID; may be empty on real CTA pings until patterns are loaded.
	Direction string
	Lat       float64
	Lon       float64
	Timestamp time.Time
}

// VehicleState is maintained by the arrival detector for each tracked vehicle.
// It holds just enough context to detect arrivals and prevent duplicate recordings.
type VehicleState struct {
	LastStopID      string
	LastArrivalTime time.Time
	LastLat         float64
	LastLon         float64
	LastTimestamp   time.Time
}
