package business

import "time"

// Headway represents the gap between two consecutive arrivals at a stop.
// Timestamp reflects the later of the two arrivals.
type Headway struct {
	StopID         string
	RouteID        string
	Direction      string
	Timestamp      time.Time
	HeadwayMinutes float64
}
