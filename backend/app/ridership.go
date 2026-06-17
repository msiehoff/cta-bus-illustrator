package app

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type RidershipRepository interface {
	GetLatestMonth() (time.Time, error)
	GetAvailableMonths() ([]time.Time, error)
	GetByMonth(month time.Time, ridershipType business.RidershipType) (map[string]*business.RidershipRecord, error)
	GetAllByRoute(routeExternalID string) ([]business.RidershipRecord, error)
	GetSystemTotals() ([]business.RidershipRecord, error)
	UpsertBatch(records []business.RidershipRecord) error
}

type RouteWithRidership struct {
	Route     business.Route
	Ridership *business.RidershipRecord
}

type RouteComparison struct {
	RouteID      string
	RouteName    string
	Current      float64
	YearAgo      *float64
	FiveYearsAgo *float64
	PreCovid2019 *float64
	RecoveryPct  *float64
	YearAgoPct   *float64
	FiveYearPct  *float64
}

type RoutesComparisonResult struct {
	CurrentMonth      time.Time
	BenchmarkMonth    time.Time
	YearAgoMonth      time.Time
	FiveYearsAgoMonth time.Time
	SystemCurrent     float64
	SystemPreCovid    *float64
	SystemRecovery    *float64
	Routes            []RouteComparison
}
