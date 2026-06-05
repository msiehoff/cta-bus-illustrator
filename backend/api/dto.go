package api

import (
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// --- Routes / map response ---

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

// --- Ridership time-series responses ---

type RidershipDataPoint struct {
	Month    string  `json:"month"`    // "2024-01"
	Type     string  `json:"type"`     // "weekday" | "saturday" | "sunday"
	AvgRides float64 `json:"avgRides"`
}

type GetRidershipResponse struct {
	Records []RidershipDataPoint `json:"records"`
}

func toRidershipResponse(records []business.RidershipRecord) GetRidershipResponse {
	points := make([]RidershipDataPoint, len(records))
	for i, r := range records {
		points[i] = RidershipDataPoint{
			Month:    r.MonthBeginning.Format("2006-01"),
			Type:     string(r.Type),
			AvgRides: r.AvgRides,
		}
	}
	return GetRidershipResponse{Records: points}
}
