package business

import "time"

type RidershipType string

const (
	RidershipTypeWeekday  RidershipType = "weekday"
	RidershipTypeSaturday RidershipType = "saturday"
	RidershipTypeSunday   RidershipType = "sunday"
)

type RidershipRecord struct {
	RouteExternalID string
	MonthBeginning  time.Time
	Type            RidershipType
	AvgRides        float64
}
