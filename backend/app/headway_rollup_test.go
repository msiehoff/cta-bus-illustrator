package app_test

import (
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func TestComputeObservedHeadways(t *testing.T) {
	base := time.Date(2026, 7, 10, 8, 0, 0, 0, time.UTC)
	arrivals := []business.Arrival{
		{StopID: "s1", RouteID: "8", Direction: "Northbound", VehicleID: "a", Timestamp: base},
		{StopID: "s1", RouteID: "8", Direction: "Northbound", VehicleID: "b", Timestamp: base.Add(10 * time.Minute)},
		{StopID: "s1", RouteID: "8", Direction: "Northbound", VehicleID: "c", Timestamp: base.Add(25 * time.Minute)},
		{StopID: "s2", RouteID: "8", Direction: "Northbound", VehicleID: "a", Timestamp: base.Add(2 * time.Minute)},
	}

	got := app.ComputeObservedHeadways(arrivals)
	if len(got) != 2 {
		t.Fatalf("expected 2 headways (s1 has 2 gaps; s2 has 1 arrival), got %d", len(got))
	}
	if got[0].HeadwayMinutes != 10 {
		t.Errorf("first gap: got %.1f want 10", got[0].HeadwayMinutes)
	}
	if got[1].HeadwayMinutes != 15 {
		t.Errorf("second gap: got %.1f want 15", got[1].HeadwayMinutes)
	}
	if got[0].FromVehicleID != "a" || got[0].ToVehicleID != "b" {
		t.Errorf("unexpected vehicles: %+v", got[0])
	}
}

func TestComputeObservedHeadwaysSkipsSameVehicle(t *testing.T) {
	base := time.Date(2026, 7, 10, 8, 0, 0, 0, time.UTC)
	arrivals := []business.Arrival{
		{StopID: "s1", RouteID: "66", Direction: "Westbound", VehicleID: "1123", Timestamp: base},
		{StopID: "s1", RouteID: "66", Direction: "Westbound", VehicleID: "1123", Timestamp: base.Add(6 * time.Second)},
		{StopID: "s1", RouteID: "66", Direction: "Westbound", VehicleID: "8690", Timestamp: base.Add(7 * time.Minute)},
		{StopID: "s1", RouteID: "66", Direction: "Westbound", VehicleID: "1351", Timestamp: base.Add(36 * time.Minute)},
	}

	got := app.ComputeObservedHeadways(arrivals)
	if len(got) != 2 {
		t.Fatalf("expected 2 headways (same-vehicle gap skipped), got %d", len(got))
	}
	if got[0].FromVehicleID != "1123" || got[0].ToVehicleID != "8690" {
		t.Errorf("first gap vehicles: %+v", got[0])
	}
	if got[1].FromVehicleID != "8690" || got[1].ToVehicleID != "1351" {
		t.Errorf("second gap vehicles: %+v", got[1])
	}
}

func TestBuildPersistedSummaries(t *testing.T) {
	base := time.Date(2026, 7, 10, 8, 0, 0, 0, time.UTC)
	headways := []business.Headway{
		{StopID: "a", RouteID: "8", Direction: "NB", Timestamp: base, HeadwayMinutes: 10},
		{StopID: "a", RouteID: "8", Direction: "NB", Timestamp: base.Add(time.Hour), HeadwayMinutes: 14},
		{StopID: "b", RouteID: "8", Direction: "NB", Timestamp: base, HeadwayMinutes: 20},
	}
	serviceDate := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)
	start, end := app.ServiceDateBounds(serviceDate)

	got := app.BuildPersistedSummaries(headways, serviceDate, start, end)
	var stops, routeDir, serviceDay int
	for _, s := range got {
		switch s.Grain {
		case business.HeadwayGrainStop:
			stops++
		case business.HeadwayGrainRouteDirection:
			routeDir++
		case business.HeadwayGrainServiceDay:
			serviceDay++
		}
	}
	if stops != 2 {
		t.Fatalf("stop rows = %d", stops)
	}
	if routeDir != 2 { // pooled + equal_stop
		t.Fatalf("route_direction rows = %d", routeDir)
	}
	if serviceDay != 2 {
		t.Fatalf("service_day rows = %d", serviceDay)
	}
}

func TestServiceDateBounds(t *testing.T) {
	d, err := app.ParseServiceDate("2026-07-10")
	if err != nil {
		t.Fatal(err)
	}
	start, end := app.ServiceDateBounds(d)
	if start.Format("2006-01-02") != "2026-07-10" {
		t.Errorf("start date = %s", start.Format("2006-01-02"))
	}
	if end.Sub(start) != 24*time.Hour {
		t.Errorf("expected 24h window, got %v", end.Sub(start))
	}
}
