package app

import (
	"context"
	"sort"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// RidershipRouteProvider resolves pipeline routes from the most recent ridership import.
type RidershipRouteProvider struct {
	repo RidershipRepository
}

func NewRidershipRouteProvider(repo RidershipRepository) *RidershipRouteProvider {
	return &RidershipRouteProvider{repo: repo}
}

func (p *RidershipRouteProvider) GetRouteIDs(ctx context.Context) ([]string, error) {
	_ = ctx

	month, err := p.repo.GetLatestMonth()
	if err != nil {
		return nil, err
	}

	ids := make(map[string]struct{})
	for _, ridershipType := range []business.RidershipType{
		business.RidershipTypeWeekday,
		business.RidershipTypeSaturday,
		business.RidershipTypeSunday,
	} {
		byRoute, err := p.repo.GetByMonth(month, ridershipType)
		if err != nil {
			return nil, err
		}
		for routeID := range byRoute {
			ids[routeID] = struct{}{}
		}
	}

	routeIDs := make([]string, 0, len(ids))
	for routeID := range ids {
		routeIDs = append(routeIDs, routeID)
	}
	sort.Strings(routeIDs)
	return routeIDs, nil
}
