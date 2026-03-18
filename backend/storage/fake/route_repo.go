package fake

import "github.com/msiehoff/cta-bus-illustrator/backend/business"

type RouteRepo struct{}

func (r *RouteRepo) GetRoutes() ([]business.Route, error) {
	return []business.Route{
		{
			ExternalID: "66",
			Name:       "Chicago",
			Segments: []business.RouteSegment{
				{Lng: -87.8065, Lat: 41.8957}, // Austin & Chicago
				{Lng: -87.7754, Lat: 41.8957}, // Central & Chicago
				{Lng: -87.7523, Lat: 41.8957}, // Laramie & Chicago
				{Lng: -87.7341, Lat: 41.8957}, // Cicero & Chicago
				{Lng: -87.7137, Lat: 41.8957}, // Pulaski & Chicago
				{Lng: -87.6948, Lat: 41.8957}, // Kedzie & Chicago
				{Lng: -87.6726, Lat: 41.8957}, // Western & Chicago
				{Lng: -87.6558, Lat: 41.8957}, // Damen & Chicago
				{Lng: -87.6418, Lat: 41.8957}, // Ashland & Chicago
				{Lng: -87.6264, Lat: 41.8957}, // Halsted & Chicago
				{Lng: -87.6134, Lat: 41.8957}, // Morgan & Chicago
				{Lng: -87.6062, Lat: 41.8957}, // Peoria & Chicago
				{Lng: -87.5961, Lat: 41.8966}, // State & Chicago
				{Lng: -87.5854, Lat: 41.8978}, // Michigan & Chicago
				{Lng: -87.5769, Lat: 41.8983}, // McClurg & Chicago (Navy Pier turn)
			},
		},
	}, nil
}

func (r *RouteRepo) GetRoute(id string) (business.Route, error) {
	return business.Route{
		ExternalID: id,
		Name:       "Chicago",
	}, nil
}

func (r *RouteRepo) CreateSegments(routeID string, segments []business.RouteSegment) error {
	return nil
}
