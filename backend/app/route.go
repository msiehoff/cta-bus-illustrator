package app

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type RouteRepository interface {
	GetRoutes() ([]business.Route, error)
	GetRoute(id string) (business.Route, error)
	CreateSegments(routeID string, segments []business.RouteSegment) error
}

type RouteSegmentDataSource interface {
	GetRouteSegments(ctx context.Context, routeID string) ([]business.RouteSegment, error)
}

type RouteService struct {
	repo          RouteRepository
	ridershipRepo RidershipRepository
}

func NewRouteService(repo RouteRepository, ridershipRepo RidershipRepository) *RouteService {
	return &RouteService{repo: repo, ridershipRepo: ridershipRepo}
}

func (s *RouteService) GetRoutesForMonth(month time.Time, ridershipType business.RidershipType) ([]RouteWithRidership, error) {
	routes, err := s.repo.GetRoutes()
	if err != nil {
		return nil, err
	}

	ridershipByRoute, err := s.ridershipRepo.GetByMonth(month, ridershipType)
	if err != nil {
		return nil, err
	}

	result := make([]RouteWithRidership, len(routes))
	for i, route := range routes {
		result[i] = RouteWithRidership{
			Route:     route,
			Ridership: ridershipByRoute[route.ExternalID],
		}
	}
	return result, nil
}

func (s *RouteService) GetLatestRidershipMonth() (time.Time, error) {
	return s.ridershipRepo.GetLatestMonth()
}

func (s *RouteService) GetAvailableRidershipMonths() ([]time.Time, error) {
	return s.ridershipRepo.GetAvailableMonths()
}

func (s *RouteService) ImportRidership(records []business.RidershipRecord) error {
	return s.ridershipRepo.UpsertBatch(records)
}

func (s *RouteService) ImportRouteSegments(ctx context.Context, dataSrc RouteSegmentDataSource) error {
	routes, err := s.repo.GetRoutes()
	if err != nil {
		return err
	}

	var errs *multierror.Error
	for _, route := range routes {
		segments, err := dataSrc.GetRouteSegments(ctx, route.ExternalID)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if err := s.repo.CreateSegments(route.ExternalID, segments); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}
