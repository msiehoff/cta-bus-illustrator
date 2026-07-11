package business

import "time"

type Arrival struct {
	StopID    string
	StopName  string // populated on list queries via join with stops
	RouteID   string
	RouteName string // populated on list queries via join with routes
	Direction string
	VehicleID string
	Timestamp time.Time
}
