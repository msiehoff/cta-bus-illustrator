package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func halstedStops() []business.Stop {
	return []business.Stop{
		{StopID: "8N-1", RouteID: "8", Direction: "Northbound", Name: "Halsted & 35th", Lat: 41.8309, Lon: -87.6441, Sequence: 1},
		{StopID: "8N-2", RouteID: "8", Direction: "Northbound", Name: "Halsted & Cermak", Lat: 41.8529, Lon: -87.6471, Sequence: 2},
		{StopID: "8N-3", RouteID: "8", Direction: "Northbound", Name: "Halsted & Roosevelt", Lat: 41.8671, Lon: -87.6474, Sequence: 3},
	}
}

func TestArrivalDetector_DetectsArrivalWithinRadius(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())

	// Ping within 40m of "Halsted & Cermak" (stop 8N-2)
	ping := business.VehiclePing{
		VehicleID: "v1", RouteID: "8", Direction: "Northbound",
		Lat: 41.8530, Lon: -87.6471, // ~11m north of stop
		Timestamp: time.Now(),
	}
	d.ProcessPing(context.Background(), ping)

	arrivals := repo.All()
	if len(arrivals) != 1 {
		t.Fatalf("expected 1 arrival, got %d", len(arrivals))
	}
	if arrivals[0].StopID != "8N-2" {
		t.Errorf("expected stop 8N-2, got %s", arrivals[0].StopID)
	}
	if arrivals[0].VehicleID != "v1" {
		t.Errorf("expected vehicle v1, got %s", arrivals[0].VehicleID)
	}
}

func TestArrivalDetector_NoDuplicateWithinCooldown(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())

	now := time.Now()
	// Same vehicle, same stop, two pings 10 seconds apart (within 3-min cooldown)
	ping := business.VehiclePing{
		VehicleID: "v1", RouteID: "8", Direction: "Northbound",
		Lat: 41.8530, Lon: -87.6471, Timestamp: now,
	}
	d.ProcessPing(context.Background(), ping)

	ping.Timestamp = now.Add(10 * time.Second)
	d.ProcessPing(context.Background(), ping)

	if len(repo.All()) != 1 {
		t.Errorf("expected 1 arrival (cooldown should suppress duplicate), got %d", len(repo.All()))
	}
}

func TestArrivalDetector_RecordsAfterCooldownExpires(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())

	now := time.Now()
	ping := business.VehiclePing{
		VehicleID: "v1", RouteID: "8", Direction: "Northbound",
		Lat: 41.8530, Lon: -87.6471, Timestamp: now,
	}
	d.ProcessPing(context.Background(), ping)

	// Second ping at same stop but 5 minutes later — past the 3-min cooldown
	ping.Timestamp = now.Add(5 * time.Minute)
	d.ProcessPing(context.Background(), ping)

	if len(repo.All()) != 2 {
		t.Errorf("expected 2 arrivals after cooldown expired, got %d", len(repo.All()))
	}
}

func TestArrivalDetector_IgnoresPingOutsideRadius(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())

	// ~500m away from any stop
	ping := business.VehiclePing{
		VehicleID: "v1", RouteID: "8", Direction: "Northbound",
		Lat: 41.8700, Lon: -87.6900, Timestamp: time.Now(),
	}
	d.ProcessPing(context.Background(), ping)

	if len(repo.All()) != 0 {
		t.Errorf("expected no arrivals for out-of-range ping, got %d", len(repo.All()))
	}
}

func TestArrivalDetector_MultipleVehiclesIndependent(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())

	now := time.Now()
	for _, vid := range []string{"v1", "v2", "v3"} {
		d.ProcessPing(context.Background(), business.VehiclePing{
			VehicleID: vid, RouteID: "8", Direction: "Northbound",
			Lat: 41.8530, Lon: -87.6471, Timestamp: now,
		})
	}

	if len(repo.All()) != 3 {
		t.Errorf("expected 3 arrivals (one per vehicle), got %d", len(repo.All()))
	}
}

func TestArrivalDetector_SkipsUnknownRoute(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	// No stops loaded for route "99"

	d.ProcessPing(context.Background(), business.VehiclePing{
		VehicleID: "v1", RouteID: "99", Direction: "Northbound",
		Lat: 41.8530, Lon: -87.6471, Timestamp: time.Now(),
	})

	if len(repo.All()) != 0 {
		t.Errorf("expected no arrivals for unknown route, got %d", len(repo.All()))
	}
}

func TestArrivalDetector_MatchesWithoutPingDirection(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())
	d.LoadPatterns(map[int]string{801: "Northbound"})

	// Real CTA pings have no direction field — resolve via pattern ID.
	ping := business.VehiclePing{
		VehicleID: "v1", RouteID: "8", PatternID: 801, Direction: "",
		Lat: 41.8530, Lon: -87.6471, Timestamp: time.Now(),
	}
	d.ProcessPing(context.Background(), ping)

	arrivals := repo.All()
	if len(arrivals) != 1 {
		t.Fatalf("expected 1 arrival via pattern ID, got %d", len(arrivals))
	}
	if arrivals[0].Direction != "Northbound" {
		t.Errorf("expected Northbound, got %q", arrivals[0].Direction)
	}
}

func TestArrivalDetector_SkipsUnknownPattern(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())
	d.LoadPatterns(map[int]string{801: "Northbound"})

	d.ProcessPing(context.Background(), business.VehiclePing{
		VehicleID: "v1", RouteID: "8", PatternID: 999, Direction: "",
		Lat: 41.8530, Lon: -87.6471, Timestamp: time.Now(),
	})

	if len(repo.All()) != 0 {
		t.Errorf("expected no arrivals for unknown pattern, got %d", len(repo.All()))
	}
}

func TestArrivalDetector_DoesNotCrossDirections(t *testing.T) {
	repo := &fake.ArrivalRepo{}
	d := app.NewArrivalDetector(repo)
	d.LoadStops("8", "Northbound", halstedStops())
	d.LoadStops("8", "Southbound", []business.Stop{
		{StopID: "8S-2", RouteID: "8", Direction: "Southbound", Name: "Halsted & Cermak SB", Lat: 41.8529, Lon: -87.6471, Sequence: 1},
	})
	d.LoadPatterns(map[int]string{802: "Southbound"})

	// Bus is on southbound pattern; should match the SB stop, not NB.
	d.ProcessPing(context.Background(), business.VehiclePing{
		VehicleID: "v1", RouteID: "8", PatternID: 802,
		Lat: 41.8530, Lon: -87.6471, Timestamp: time.Now(),
	})

	arrivals := repo.All()
	if len(arrivals) != 1 {
		t.Fatalf("expected 1 arrival, got %d", len(arrivals))
	}
	if arrivals[0].StopID != "8S-2" || arrivals[0].Direction != "Southbound" {
		t.Errorf("expected southbound stop, got stop=%s dir=%s", arrivals[0].StopID, arrivals[0].Direction)
	}
}

