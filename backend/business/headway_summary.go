package business

import "time"

// Summary grain / method constants persisted by the daily headway job.
const (
	HeadwayGrainStop           = "stop"
	HeadwayGrainRouteDirection = "route_direction"
	HeadwayGrainRoute          = "route" // both directions combined
	HeadwayGrainServiceDay     = "service_day"

	HeadwayMethodPooled    = "pooled"
	HeadwayMethodEqualStop = "equal_stop"

	// Temporal window keys (America/Chicago clock slices).
	HeadwayWindowOverall    = "overall"
	HeadwayWindowWeekdayAM  = "weekday_am"
	HeadwayWindowWeekdayMid = "weekday_mid"
	HeadwayWindowWeekdayPM  = "weekday_pm"
	HeadwayWindowWeekdayEve = "weekday_eve"
	HeadwayWindowSaturday   = "saturday"
	HeadwayWindowSunday     = "sunday"
)

// HeadwaySummaryStats is aggregate metrics over a set of observed headway gaps.
type HeadwaySummaryStats struct {
	Count          int
	MeanMinutes    float64
	MedianMinutes  float64
	StdDevMinutes  float64
	CV             float64 // stddev / mean; 0 when mean is 0 or count < 2
	AvgWaitMinutes float64 // mean/2 (regular-service approximation)
}

// HeadwayStopSummary is stats for one stop within a broader query.
type HeadwayStopSummary struct {
	StopID    string
	StopName  string
	RouteID   string
	RouteName string
	Direction string
	HeadwaySummaryStats
}

// HeadwaySummary is a persisted aggregate for a service date, window, grain, and method.
type HeadwaySummary struct {
	ServiceDate time.Time
	WindowStart time.Time
	WindowEnd   time.Time
	WindowKey   string // overall | weekday_am | … — temporal slice
	Grain       string
	Method      string
	StopID      string
	RouteID     string
	Direction   string
	StopName    string // join-populated on read
	RouteName   string // join-populated on read
	HeadwaySummaryStats
}
