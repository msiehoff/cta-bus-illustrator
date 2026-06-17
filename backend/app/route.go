package app

import (
	"context"
	"log"
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

func (s *RouteService) GetRouteRidership(routeExternalID string) ([]business.RidershipRecord, error) {
	return s.ridershipRepo.GetAllByRoute(routeExternalID)
}

func (s *RouteService) GetSystemRidership() ([]business.RidershipRecord, error) {
	return s.ridershipRepo.GetSystemTotals()
}

func (s *RouteService) GetRoutesComparison(ridershipType business.RidershipType) (*RoutesComparisonResult, error) {
	latest, err := s.ridershipRepo.GetLatestMonth()
	if err != nil {
		return nil, err
	}

	yearAgo := offsetMonth(latest, 1)
	fiveYearsAgo := offsetMonth(latest, 5)
	preCovid := preCovidMonth(latest)

	currentByRoute, err := s.ridershipRepo.GetByMonth(latest, ridershipType)
	if err != nil {
		return nil, err
	}
	yearAgoByRoute, err := s.ridershipRepo.GetByMonth(yearAgo, ridershipType)
	if err != nil {
		return nil, err
	}
	fiveYearsAgoByRoute, err := s.ridershipRepo.GetByMonth(fiveYearsAgo, ridershipType)
	if err != nil {
		return nil, err
	}
	preCovidByRoute, err := s.ridershipRepo.GetByMonth(preCovid, ridershipType)
	if err != nil {
		return nil, err
	}

	routes, err := s.repo.GetRoutes()
	if err != nil {
		return nil, err
	}

	comparisons := make([]RouteComparison, 0, len(routes))
	var systemCurrent float64
	var systemPreCovid float64
	var hasPreCovid bool

	for _, route := range routes {
		currentRec := currentByRoute[route.ExternalID]
		if currentRec == nil {
			continue
		}

		cmp := RouteComparison{
			RouteID:   route.ExternalID,
			RouteName: route.Name,
			Current:   currentRec.AvgRides,
		}

		if rec := yearAgoByRoute[route.ExternalID]; rec != nil {
			cmp.YearAgo = &rec.AvgRides
			cmp.YearAgoPct = pctChange(currentRec.AvgRides, rec.AvgRides)
		}
		if rec := fiveYearsAgoByRoute[route.ExternalID]; rec != nil {
			cmp.FiveYearsAgo = &rec.AvgRides
			cmp.FiveYearPct = pctChange(currentRec.AvgRides, rec.AvgRides)
		}
		if rec := preCovidByRoute[route.ExternalID]; rec != nil {
			cmp.PreCovid2019 = &rec.AvgRides
			cmp.RecoveryPct = recoveryPct(currentRec.AvgRides, rec.AvgRides)
			systemPreCovid += rec.AvgRides
			hasPreCovid = true
		}

		systemCurrent += currentRec.AvgRides
		comparisons = append(comparisons, cmp)
	}

	result := &RoutesComparisonResult{
		CurrentMonth:      latest,
		BenchmarkMonth:    preCovid,
		YearAgoMonth:      yearAgo,
		FiveYearsAgoMonth: fiveYearsAgo,
		SystemCurrent:     systemCurrent,
		Routes:            comparisons,
	}
	if hasPreCovid && systemPreCovid > 0 {
		result.SystemPreCovid = &systemPreCovid
		result.SystemRecovery = recoveryPct(systemCurrent, systemPreCovid)
	}

	return result, nil
}

func offsetMonth(month time.Time, years int) time.Time {
	return time.Date(month.Year()-years, month.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func preCovidMonth(month time.Time) time.Time {
	return time.Date(2019, month.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func pctChange(current, baseline float64) *float64 {
	if baseline == 0 {
		return nil
	}
	v := ((current - baseline) / baseline) * 100
	return &v
}

func recoveryPct(current, baseline float64) *float64 {
	if baseline == 0 {
		return nil
	}
	v := (current / baseline) * 100
	return &v
}

func (s *RouteService) ImportRidership(records []business.RidershipRecord) error {
	return s.ridershipRepo.UpsertBatch(records)
}

func (s *RouteService) ImportRouteSegmentsFromSrc(ctx context.Context, dataSrc RouteSegmentDataSource) error {
	routes, err := s.repo.GetRoutes()
	if err != nil {
		return err
	}
	log.Printf("\nroutes retrieved: %d", len(routes))

	var errs *multierror.Error
	for _, route := range routes {
		segments, err := dataSrc.GetRouteSegments(ctx, route.ExternalID)
		if err != nil {
			log.Printf("\nfailed to get route segments for route %s: %v", route.ExternalID, err)
			errs = multierror.Append(errs, err)
			continue
		}

		log.Printf("\nsegments retrieved for route %s: %d", route.ExternalID, len(segments))
		if err := s.ImportRouteSegments(ctx, route.ExternalID, segments); err != nil {
			log.Printf("failed to import segments for route %s: %v", route.ExternalID, err)
			errs = multierror.Append(errs, err)
		}
	}

	log.Printf("\ntotal errors: %d", errs.Len())
	return errs.ErrorOrNil()
}

func (s *RouteService) ImportRouteSegments(ctx context.Context, routeID string, segments []business.RouteSegment) error {
	if err := s.repo.CreateSegments(routeID, segments); err != nil {
		log.Printf("failed to create segments for route %s: %v", routeID, err)
		return err
	}
	return nil
}
