package app

import (
	"context"
	"log"
	"math"

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
	// stopsByRouteDir is keyed by "routeID:direction".
	stopsByRouteDir map[string][]business.Stop
	// patternDir maps CTA pattern ID → normalized direction (e.g. "Eastbound").
	patternDir map[int]string
	repo       ArrivalRepository
}

// NewArrivalDetector creates a detector backed by the given repository.
// Stops must be pre-loaded via LoadStops before processing pings.
func NewArrivalDetector(repo ArrivalRepository) *ArrivalDetector {
	return &ArrivalDetector{
		vehicleStates:   make(map[string]*business.VehicleState),
		stopsByRouteDir: make(map[string][]business.Stop),
		patternDir:      make(map[int]string),
		repo:            repo,
	}
}

// LoadStops caches the stops for a route+direction so the detector can do
// proximity checks without hitting the CTA API on every ping.
func (d *ArrivalDetector) LoadStops(routeID, direction string, stops []business.Stop) {
	d.stopsByRouteDir[routeKey(routeID, direction)] = stops
}

// LoadPatterns caches pattern ID → direction mappings used to resolve vehicle direction.
func (d *ArrivalDetector) LoadPatterns(patterns map[int]string) {
	for pid, dir := range patterns {
		d.patternDir[pid] = dir
	}
}

// ProcessPing evaluates a single vehicle ping and saves an Arrival if the vehicle
// has entered the radius of a new stop (subject to cooldown).
func (d *ArrivalDetector) ProcessPing(ctx context.Context, ping business.VehiclePing) {
	state, ok := d.vehicleStates[ping.VehicleID]
	if !ok {
		state = &business.VehicleState{}
		d.vehicleStates[ping.VehicleID] = state
		Debugf("arrival_detector: tracking new vehicle=%s route=%s pid=%d dir=%q",
			ping.VehicleID, ping.RouteID, ping.PatternID, ping.Direction)
	}

	direction := d.resolveDirection(ping)
	if direction == "" {
		Debugf("arrival_detector: skip vehicle=%s route=%s pid=%d — cannot resolve direction",
			ping.VehicleID, ping.RouteID, ping.PatternID)
		return
	}

	stops, ok := d.stopsByRouteDir[routeKey(ping.RouteID, direction)]
	if !ok || len(stops) == 0 {
		Debugf("arrival_detector: skip vehicle=%s route=%s dir=%s — no stops loaded",
			ping.VehicleID, ping.RouteID, direction)
		return
	}

	nearest, dist := nearestStop(stops, ping.Lat, ping.Lon)
	if nearest == nil {
		return
	}

	if dist > ArrivalRadiusMeters {
		if dist < 200 {
			Debugf("arrival_detector: skip vehicle=%s route=%s dir=%s nearest=%s(%s) dist=%.1fm > radius=%.0fm (near miss)",
				ping.VehicleID, ping.RouteID, direction, nearest.StopID, nearest.Name, dist, ArrivalRadiusMeters)
		}
		return
	}

	if nearest.StopID == state.LastStopID && !state.LastArrivalTime.IsZero() {
		elapsedMin := ping.Timestamp.Sub(state.LastArrivalTime).Minutes()
		if elapsedMin < ArrivalCooldownMinutes {
			Debugf("arrival_detector: skip vehicle=%s stop=%s — cooldown (%.1fmin < %.0fmin)",
				ping.VehicleID, nearest.StopID, elapsedMin, ArrivalCooldownMinutes)
			return
		}
	}

	arrival := business.Arrival{
		StopID:    nearest.StopID,
		RouteID:   ping.RouteID,
		Direction: direction,
		VehicleID: ping.VehicleID,
		Timestamp: ping.Timestamp,
	}

	if err := d.repo.SaveArrival(ctx, arrival); err != nil {
		log.Printf("arrival_detector: failed to save arrival vehicle=%s stop=%s: %v",
			ping.VehicleID, nearest.StopID, err)
		return
	}

	log.Printf("arrival_detector: arrival recorded vehicle=%s route=%s dir=%s stop=%s(%s) dist=%.1fm pid=%d",
		ping.VehicleID, ping.RouteID, direction, nearest.StopID, nearest.Name, dist, ping.PatternID)

	state.LastStopID = nearest.StopID
	state.LastArrivalTime = ping.Timestamp
	state.LastLat = ping.Lat
	state.LastLon = ping.Lon
	state.LastTimestamp = ping.Timestamp
}

// resolveDirection prefers an explicit ping direction, then pattern ID lookup.
func (d *ArrivalDetector) resolveDirection(ping business.VehiclePing) string {
	if ping.Direction != "" {
		return business.NormalizeDirection(ping.Direction)
	}
	if ping.PatternID != 0 {
		if dir, ok := d.patternDir[ping.PatternID]; ok {
			return dir
		}
		Debugf("arrival_detector: unknown pattern pid=%d route=%s", ping.PatternID, ping.RouteID)
	}
	return ""
}

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
