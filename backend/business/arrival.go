package business

import "time"

type Arrival struct {
	StopID    string
	RouteID   string
	Direction string
	VehicleID string
	Timestamp time.Time
}
