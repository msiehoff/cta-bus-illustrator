package api

import "github.com/msiehoff/cta-bus-illustrator/backend/business"

type routeProperties struct {
	RouteID   string `json:"routeId"`
	RouteName string `json:"routeName"`
	Color     string `json:"color"`
	Ridership int    `json:"ridership"`
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

func toGetRoutesResponse(routes []business.Route) GetRoutesResponse {
	features := make([]routeFeature, len(routes))
	for i, r := range routes {
		features[i] = routeFeature{
			Type: "Feature",
			Properties: routeProperties{
				RouteID:   r.ExternalID,
				RouteName: r.Name,
				Color:     r.Color,
				Ridership: r.Ridership,
			},
			Geometry: routeGeometry{
				Type:        "LineString",
				Coordinates: r.Coordinates,
			},
		}
	}
	return GetRoutesResponse{Type: "FeatureCollection", Features: features}
}
