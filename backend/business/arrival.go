package business

import "time"

type Arrival struct {
	StopID    string
	StopName  string // populated on list queries via join with stops
	RouteID   string
	Direction string
	VehicleID string
	Timestamp time.Time
}
