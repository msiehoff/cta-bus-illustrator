package app

import (
	"context"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// HeadwayRepository persists observed headway gaps.
type HeadwayRepository interface {
	DeleteInRange(ctx context.Context, start, end time.Time) (int64, error)
	InsertBatch(ctx context.Context, headways []business.Headway) error
	List(ctx context.Context, filter HeadwayListFilter) ([]business.Headway, error)
	Count(ctx context.Context, filter HeadwayListFilter) (int64, error)
	// ListAll returns matching headways without pagination (capped) for summary computation.
	ListAll(ctx context.Context, filter HeadwayListFilter) ([]business.Headway, error)
}

// HeadwayListFilter scopes headway queries for admin / rider APIs.
type HeadwayListFilter struct {
	RouteID   string
	Direction string
	Stop      string // stop ID exact or name substring
	VehicleID string // matches from_vehicle_id or to_vehicle_id
	From      *time.Time
	To        *time.Time
	SortAsc   bool
	Limit     int
	Offset    int
}

// HeadwayJobRunRepository tracks headway rollup job metadata.
type HeadwayJobRunRepository interface {
	Create(ctx context.Context, run business.HeadwayJobRun) (business.HeadwayJobRun, error)
	Update(ctx context.Context, run business.HeadwayJobRun) error
	List(ctx context.Context, limit, offset int) ([]business.HeadwayJobRun, error)
	Get(ctx context.Context, id int64) (business.HeadwayJobRun, error)
}

// HeadwaySummaryRepository persists daily headway aggregates.
type HeadwaySummaryRepository interface {
	DeleteForServiceDate(ctx context.Context, serviceDate time.Time) (int64, error)
	InsertBatch(ctx context.Context, summaries []business.HeadwaySummary) error
	List(ctx context.Context, filter HeadwaySummaryFilter) ([]business.HeadwaySummary, error)
	Count(ctx context.Context, filter HeadwaySummaryFilter) (int64, error)
}

// HeadwaySummaryFilter scopes reads of persisted summaries.
// ServiceDate is optional; zero means all dates (or use From/To for a range).
type HeadwaySummaryFilter struct {
	ServiceDate time.Time
	From        *time.Time // inclusive service_date lower bound
	To          *time.Time // inclusive service_date upper bound
	Grain       string
	Method      string
	RouteID     string
	Direction   string
	StopID      string // exact stop id; empty = any
	Stop        string // id or name substring (for stop-grain list UX)
	SortAsc     bool
	Limit       int
	Offset      int
}
