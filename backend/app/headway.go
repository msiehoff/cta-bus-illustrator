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
	ListInRange(ctx context.Context, start, end time.Time, filter HeadwayListFilter) ([]business.Headway, error)
	CountInRange(ctx context.Context, start, end time.Time, filter HeadwayListFilter) (int64, error)
}

// HeadwayListFilter scopes headway queries for admin / rider APIs.
type HeadwayListFilter struct {
	RouteID   string
	Direction string
	StopID    string
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
