package business

import "time"

// VehiclePing is emitted by the vehicle fetcher stage for each polled vehicle position.
type VehiclePing struct {
	VehicleID string
	RouteID   string
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
