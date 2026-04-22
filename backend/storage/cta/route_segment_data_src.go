package cta

import (
	"context"

	"log"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type RouteSegmentDataSource struct {
	client *Client
}

func NewRouteSegmentDataSource(client *Client) *RouteSegmentDataSource {
	return &RouteSegmentDataSource{client: client}
}

func (d *RouteSegmentDataSource) GetRouteSegments(ctx context.Context, routeID string) ([]business.RouteSegment, error) {
	routeSegments, err := d.client.GetRoutePattern(routeID)
	if err != nil {
		log.Printf("\nfailed to get route segments for route %s: %v", routeID, err)
		return nil, err
	}

	return SegmentsFromPatternResponse(routeSegments)
}
