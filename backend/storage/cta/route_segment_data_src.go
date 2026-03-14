package cta

import (
	"context"

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
		return nil, err
	}

	segments := make([]business.RouteSegment, len(routeSegments.BustimeResponse.Ptr))

	direction := routeSegments.BustimeResponse.Ptr[0]
	for i, segment := range direction.Pt {
		segments[i] = business.RouteSegment{
			Sequence: segment.Seq,
			Lat:      segment.Lat,
			Lng:      segment.Lon,
		}
	}

	return segments, nil
}
