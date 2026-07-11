package fake

import (
	"context"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

// CTAClient returns realistic-looking Chicago bus data without hitting the real CTA API.
// Vehicles are simulated moving along the #8 Halsted and #66 Chicago corridors.
type CTAClient struct {
	tick int // advances each call to simulate vehicle movement
}

func NewCTAClient() *CTAClient {
	return &CTAClient{}
}

// GetVehicles returns a fixed set of simulated vehicle pings.
// Positions shift slightly on each call to mimic real bus movement.
func (c *CTAClient) GetVehicles(_ context.Context, routeIDs []string) ([]business.VehiclePing, error) {
	c.tick++
	now := time.Now()

	// Index routeIDs for quick lookup so we only return pings for requested routes.
	requested := make(map[string]bool, len(routeIDs))
	for _, id := range routeIDs {
		requested[id] = true
	}

	all := c.allVehicles(now)
	out := make([]business.VehiclePing, 0, len(all))
	for _, p := range all {
		if requested[p.RouteID] {
			out = append(out, p)
		}
	}
	return out, nil
}

// GetStops returns a fixed set of stops for the supported fake routes.
func (c *CTAClient) GetStops(_ context.Context, routeID, direction string) ([]business.Stop, error) {
	key := routeID + ":" + direction
	if stops, ok := fakeStops[key]; ok {
		return stops, nil
	}
	return []business.Stop{}, nil
}

// GetPatterns returns pattern ID → direction for supported fake routes.
func (c *CTAClient) GetPatterns(_ context.Context, routeID string) (map[int]string, error) {
	switch routeID {
	case "8":
		return map[int]string{801: "Northbound", 802: "Southbound"}, nil
	case "66":
		return map[int]string{661: "Eastbound", 662: "Westbound"}, nil
	default:
		return map[int]string{}, nil
	}
}

// allVehicles builds simulated pings. The tick offset nudges vehicles along their routes.
// Direction is left empty — resolution uses PatternID via GetPatterns, matching real CTA.
func (c *CTAClient) allVehicles(now time.Time) []business.VehiclePing {
	offset := float64(c.tick) * 0.0003 // ~33m northward per tick

	return []business.VehiclePing{
		// Route 8 Halsted — Northbound
		{VehicleID: "8001", RouteID: "8", PatternID: 801,
			Lat: 41.8480 + offset, Lon: -87.6441, Timestamp: now},
		{VehicleID: "8002", RouteID: "8", PatternID: 801,
			Lat: 41.8780 + offset, Lon: -87.6441, Timestamp: now},
		{VehicleID: "8003", RouteID: "8", PatternID: 801,
			Lat: 41.9050 + offset, Lon: -87.6441, Timestamp: now},

		// Route 8 Halsted — Southbound
		{VehicleID: "8004", RouteID: "8", PatternID: 802,
			Lat: 41.9600 - offset, Lon: -87.6441, Timestamp: now},
		{VehicleID: "8005", RouteID: "8", PatternID: 802,
			Lat: 41.9200 - offset, Lon: -87.6441, Timestamp: now},

		// Route 66 Chicago — Eastbound
		{VehicleID: "6601", RouteID: "66", PatternID: 661,
			Lat: 41.8957, Lon: -87.7600 + offset, Timestamp: now},
		{VehicleID: "6602", RouteID: "66", PatternID: 661,
			Lat: 41.8957, Lon: -87.7100 + offset, Timestamp: now},

		// Route 66 Chicago — Westbound
		{VehicleID: "6603", RouteID: "66", PatternID: 662,
			Lat: 41.8957, Lon: -87.6200 - offset, Timestamp: now},
	}
}

// fakeStops is a keyed map of realistic stops for supported routes.
// Coordinates are real CTA stop positions.
var fakeStops = map[string][]business.Stop{
	"8:Northbound": {
		{StopID: "8N-1", RouteID: "8", Direction: "Northbound", Name: "Halsted & 79th", Lat: 41.7508, Lon: -87.6441, Sequence: 1},
		{StopID: "8N-2", RouteID: "8", Direction: "Northbound", Name: "Halsted & 63rd", Lat: 41.7798, Lon: -87.6441, Sequence: 2},
		{StopID: "8N-3", RouteID: "8", Direction: "Northbound", Name: "Halsted & 47th", Lat: 41.8096, Lon: -87.6441, Sequence: 3},
		{StopID: "8N-4", RouteID: "8", Direction: "Northbound", Name: "Halsted & 35th", Lat: 41.8309, Lon: -87.6441, Sequence: 4},
		{StopID: "8N-5", RouteID: "8", Direction: "Northbound", Name: "Halsted & Cermak", Lat: 41.8529, Lon: -87.6471, Sequence: 5},
		{StopID: "8N-6", RouteID: "8", Direction: "Northbound", Name: "Halsted & Roosevelt", Lat: 41.8671, Lon: -87.6474, Sequence: 6},
		{StopID: "8N-7", RouteID: "8", Direction: "Northbound", Name: "Halsted & 16th", Lat: 41.8584, Lon: -87.6463, Sequence: 7},
		{StopID: "8N-8", RouteID: "8", Direction: "Northbound", Name: "Halsted & Madison", Lat: 41.8820, Lon: -87.6474, Sequence: 8},
		{StopID: "8N-9", RouteID: "8", Direction: "Northbound", Name: "Halsted & Chicago", Lat: 41.8963, Lon: -87.6474, Sequence: 9},
		{StopID: "8N-10", RouteID: "8", Direction: "Northbound", Name: "Halsted & North", Lat: 41.9100, Lon: -87.6474, Sequence: 10},
		{StopID: "8N-11", RouteID: "8", Direction: "Northbound", Name: "Halsted & Fullerton", Lat: 41.9250, Lon: -87.6486, Sequence: 11},
		{StopID: "8N-12", RouteID: "8", Direction: "Northbound", Name: "Halsted & Diversey", Lat: 41.9325, Lon: -87.6486, Sequence: 12},
	},
	"8:Southbound": {
		{StopID: "8S-1", RouteID: "8", Direction: "Southbound", Name: "Halsted & Diversey", Lat: 41.9325, Lon: -87.6486, Sequence: 1},
		{StopID: "8S-2", RouteID: "8", Direction: "Southbound", Name: "Halsted & Fullerton", Lat: 41.9250, Lon: -87.6486, Sequence: 2},
		{StopID: "8S-3", RouteID: "8", Direction: "Southbound", Name: "Halsted & North", Lat: 41.9100, Lon: -87.6474, Sequence: 3},
		{StopID: "8S-4", RouteID: "8", Direction: "Southbound", Name: "Halsted & Chicago", Lat: 41.8963, Lon: -87.6474, Sequence: 4},
		{StopID: "8S-5", RouteID: "8", Direction: "Southbound", Name: "Halsted & Madison", Lat: 41.8820, Lon: -87.6474, Sequence: 5},
		{StopID: "8S-6", RouteID: "8", Direction: "Southbound", Name: "Halsted & Roosevelt", Lat: 41.8671, Lon: -87.6474, Sequence: 6},
		{StopID: "8S-7", RouteID: "8", Direction: "Southbound", Name: "Halsted & Cermak", Lat: 41.8529, Lon: -87.6471, Sequence: 7},
		{StopID: "8S-8", RouteID: "8", Direction: "Southbound", Name: "Halsted & 35th", Lat: 41.8309, Lon: -87.6441, Sequence: 8},
	},
	"66:Eastbound": {
		{StopID: "66E-1", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Austin", Lat: 41.8957, Lon: -87.8065, Sequence: 1},
		{StopID: "66E-2", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Cicero", Lat: 41.8957, Lon: -87.7455, Sequence: 2},
		{StopID: "66E-3", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Pulaski", Lat: 41.8957, Lon: -87.7241, Sequence: 3},
		{StopID: "66E-4", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Kedzie", Lat: 41.8957, Lon: -87.7054, Sequence: 4},
		{StopID: "66E-5", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Western", Lat: 41.8957, Lon: -87.6876, Sequence: 5},
		{StopID: "66E-6", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Ashland", Lat: 41.8957, Lon: -87.6637, Sequence: 6},
		{StopID: "66E-7", RouteID: "66", Direction: "Eastbound", Name: "Chicago & Halsted", Lat: 41.8963, Lon: -87.6474, Sequence: 7},
		{StopID: "66E-8", RouteID: "66", Direction: "Eastbound", Name: "Chicago & State", Lat: 41.8966, Lon: -87.6279, Sequence: 8},
	},
	"66:Westbound": {
		{StopID: "66W-1", RouteID: "66", Direction: "Westbound", Name: "Chicago & State", Lat: 41.8966, Lon: -87.6279, Sequence: 1},
		{StopID: "66W-2", RouteID: "66", Direction: "Westbound", Name: "Chicago & Halsted", Lat: 41.8963, Lon: -87.6474, Sequence: 2},
		{StopID: "66W-3", RouteID: "66", Direction: "Westbound", Name: "Chicago & Ashland", Lat: 41.8957, Lon: -87.6637, Sequence: 3},
		{StopID: "66W-4", RouteID: "66", Direction: "Westbound", Name: "Chicago & Western", Lat: 41.8957, Lon: -87.6876, Sequence: 4},
		{StopID: "66W-5", RouteID: "66", Direction: "Westbound", Name: "Chicago & Kedzie", Lat: 41.8957, Lon: -87.7054, Sequence: 5},
		{StopID: "66W-6", RouteID: "66", Direction: "Westbound", Name: "Chicago & Pulaski", Lat: 41.8957, Lon: -87.7241, Sequence: 6},
		{StopID: "66W-7", RouteID: "66", Direction: "Westbound", Name: "Chicago & Cicero", Lat: 41.8957, Lon: -87.7455, Sequence: 7},
		{StopID: "66W-8", RouteID: "66", Direction: "Westbound", Name: "Chicago & Austin", Lat: 41.8957, Lon: -87.8065, Sequence: 8},
	},
}
