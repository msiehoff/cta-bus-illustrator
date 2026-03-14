package app

import (
	"context"

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
	repo RouteRepository
}

func NewRouteService(repo RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) GetRoutes() ([]business.Route, error) {
	return s.repo.GetRoutes()
}

func (s *RouteService) ImportRouteSegments(ctx context.Context, dataSrc RouteSegmentDataSource, routeID string) error {
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
