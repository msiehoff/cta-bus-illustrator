package app

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

const (
	// ArrivalRadiusMeters is how close a vehicle must be to a stop to count as an arrival.
	ArrivalRadiusMeters = 40.0

	// ArrivalCooldownMinutes prevents duplicate arrivals for the same (vehicle, stop) pair.
	// Buses may dwell at stops or GPS jitter can repeatedly place them within the radius.
	ArrivalCooldownMinutes = 3.0
)

// ArrivalDetector is Stage 2 of the pipeline.
// It consumes VehiclePings, maintains per-vehicle state, and emits Arrivals when a
// vehicle is detected arriving at a stop it has not recently visited.
//
// It must be driven by a single goroutine — arrival detection is stateful and
// order-sensitive. Do not call ProcessPing concurrently.
type ArrivalDetector struct {
	vehicleStates map[string]*business.VehicleState // keyed by VehicleID
	stops         map[string][]business.Stop        // keyed by "routeID:direction"
	repo          ArrivalRepository
}

// NewArrivalDetector creates a detector backed by the given repository.
// Stops must be pre-loaded via LoadStops before processing pings.
func NewArrivalDetector(repo ArrivalRepository) *ArrivalDetector {
	return &ArrivalDetector{
		vehicleStates: make(map[string]*business.VehicleState),
		stops:         make(map[string][]business.Stop),
		repo:          repo,
	}
}

// LoadStops caches the stops for a route+direction so the detector can do
// proximity checks without hitting the CTA API on every ping.
func (d *ArrivalDetector) LoadStops(routeID, direction string, stops []business.Stop) {
	d.stops[routeKey(routeID, direction)] = stops
}

// ProcessPing evaluates a single vehicle ping and saves an Arrival if the vehicle
// has entered the radius of a new stop (subject to cooldown).
func (d *ArrivalDetector) ProcessPing(ctx context.Context, ping business.VehiclePing) {
	state, ok := d.vehicleStates[ping.VehicleID]
	if !ok {
		state = &business.VehicleState{}
		d.vehicleStates[ping.VehicleID] = state
	}

	stops, ok := d.stops[routeKey(ping.RouteID, ping.Direction)]
	if !ok {
		// Stops not yet loaded for this route/direction — skip silently.
		return
	}

	nearest, dist := nearestStop(stops, ping.Lat, ping.Lon)
	if nearest == nil {
		return
	}

	if dist > ArrivalRadiusMeters {
		// Not close enough to any stop.
		return
	}

	if nearest.StopID == state.LastStopID {
		// Already recorded this stop for this vehicle; check cooldown.
		if time.Since(state.LastArrivalTime).Minutes() < ArrivalCooldownMinutes {
			return
		}
	}

	arrival := business.Arrival{
		StopID:    nearest.StopID,
		RouteID:   ping.RouteID,
		Direction: ping.Direction,
		VehicleID: ping.VehicleID,
		Timestamp: ping.Timestamp,
	}

	if err := d.repo.SaveArrival(ctx, arrival); err != nil {
		log.Printf("arrival_detector: failed to save arrival vehicle=%s stop=%s: %v",
			ping.VehicleID, nearest.StopID, err)
		return
	}

	log.Printf("arrival_detector: arrival recorded vehicle=%s route=%s dir=%s stop=%s(%s) dist=%.1fm",
		ping.VehicleID, ping.RouteID, ping.Direction, nearest.StopID, nearest.Name, dist)

	state.LastStopID = nearest.StopID
	state.LastArrivalTime = ping.Timestamp
	state.LastLat = ping.Lat
	state.LastLon = ping.Lon
	state.LastTimestamp = ping.Timestamp
}

// nearestStop returns the closest stop to (lat, lon) and the distance to it in metres.
func nearestStop(stops []business.Stop, lat, lon float64) (*business.Stop, float64) {
	var nearest *business.Stop
	minDist := math.MaxFloat64

	for i := range stops {
		d := haversineMeters(lat, lon, stops[i].Lat, stops[i].Lon)
		if d < minDist {
			minDist = d
			nearest = &stops[i]
		}
	}
	return nearest, minDist
}

// haversineMeters returns the great-circle distance between two lat/lon points in metres.
func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6_371_000.0
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func routeKey(routeID, direction string) string {
	return routeID + ":" + direction
}
