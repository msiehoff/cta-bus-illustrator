package app

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type RidershipRepository interface {
	GetLatestMonth() (time.Time, error)
	GetAvailableMonths() ([]time.Time, error)
	GetByMonth(month time.Time, ridershipType business.RidershipType) (map[string]*business.RidershipRecord, error)
	UpsertBatch(records []business.RidershipRecord) error
}

type RouteWithRidership struct {
	Route     business.Route
	Ridership *business.RidershipRecord
}
