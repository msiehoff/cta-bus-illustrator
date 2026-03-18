package api

import "github.com/msiehoff/cta-bus-illustrator/backend/app"

type routeProperties struct {
	RouteID   string   `json:"routeId"`
	RouteName string   `json:"routeName"`
	AvgRides  *float64 `json:"avgRides,omitempty"`
}

type routeGeometry struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

type routeFeature struct {
	Type       string          `json:"type"`
	Properties routeProperties `json:"properties"`
	Geometry   routeGeometry   `json:"geometry"`
}

type GetRoutesResponse struct {
	Type     string         `json:"type"`
	Features []routeFeature `json:"features"`
}

func toGetRoutesResponse(routes []app.RouteWithRidership) GetRoutesResponse {
	features := make([]routeFeature, len(routes))
	for i, rwr := range routes {
		r := rwr.Route
		coordinates := make([][2]float64, len(r.Segments))
		for j, s := range r.Segments {
			coordinates[j] = [2]float64{s.Lng, s.Lat}
		}

		props := routeProperties{
			RouteID:   r.ExternalID,
			RouteName: r.Name,
		}
		if rwr.Ridership != nil {
			props.AvgRides = &rwr.Ridership.AvgRides
		}

		features[i] = routeFeature{
			Type:       "Feature",
			Properties: props,
			Geometry: routeGeometry{
				Type:        "LineString",
				Coordinates: coordinates,
			},
		}
	}
	return GetRoutesResponse{Type: "FeatureCollection", Features: features}
}
