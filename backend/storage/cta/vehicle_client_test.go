package cta

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func TestVehicleClient_GetVehicles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getvehicles" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("rt"); got != "8,66" {
			t.Fatalf("expected rt=8,66, got %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"bustime-response": map[string]interface{}{
				"vehicle": []map[string]interface{}{
					{
						"vid": "1234", "tmstmp": "20240709 14:30:00",
						"lat": "41.8781", "lon": "-87.6298",
						"rt": "8", "pid": 801, "stsd": "2024-07-09",
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.BaseURL = server.URL
	vc := NewVehicleClient(client)

	pings, err := vc.GetVehicles(t.Context(), []string{"8", "66"})
	if err != nil {
		t.Fatalf("GetVehicles: %v", err)
	}
	if len(pings) != 1 {
		t.Fatalf("expected 1 ping, got %d", len(pings))
	}

	p := pings[0]
	if p.VehicleID != "1234" || p.RouteID != "8" || p.Direction != "" || p.PatternID != 801 {
		t.Errorf("unexpected ping: %+v", p)
	}
	if p.Lat != 41.8781 || p.Lon != -87.6298 {
		t.Errorf("unexpected coords: %f, %f", p.Lat, p.Lon)
	}
}

func TestVehicleClient_GetStops(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getstops" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"bustime-response": map[string]interface{}{
				"stops": []map[string]interface{}{
					{"stpid": "100", "stpnm": "Halsted & Madison", "lat": 41.88, "lon": -87.64},
					{"stpid": "101", "stpnm": "Halsted & Monroe", "lat": 41.89, "lon": -87.64},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.BaseURL = server.URL
	vc := NewVehicleClient(client)

	stops, err := vc.GetStops(t.Context(), "8", "Northbound")
	if err != nil {
		t.Fatalf("GetStops: %v", err)
	}
	if len(stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(stops))
	}
	if stops[0].StopID != "100" || stops[0].Sequence != 1 {
		t.Errorf("unexpected first stop: %+v", stops[0])
	}
	if stops[1].Sequence != 2 {
		t.Errorf("expected sequence 2, got %d", stops[1].Sequence)
	}

	want := business.Stop{StopID: "100", RouteID: "8", Direction: "Northbound", Name: "Halsted & Madison", Lat: 41.88, Lon: -87.64, Sequence: 1}
	if stops[0] != want {
		t.Errorf("stop mismatch:\n got: %+v\nwant: %+v", stops[0], want)
	}
}

func TestVehicleClient_GetPatterns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getpatterns" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"bustime-response": map[string]interface{}{
				"ptr": []map[string]interface{}{
					{"pid": 100, "ln": 1000, "rtdir": "Eastbound", "pt": []any{}},
					{"pid": 200, "ln": 1000, "rtdir": "WEST", "pt": []any{}},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.BaseURL = server.URL
	vc := NewVehicleClient(client)

	patterns, err := vc.GetPatterns(t.Context(), "66")
	if err != nil {
		t.Fatalf("GetPatterns: %v", err)
	}
	if patterns[100] != "Eastbound" {
		t.Errorf("expected Eastbound for pid 100, got %q", patterns[100])
	}
	if patterns[200] != "Westbound" {
		t.Errorf("expected Westbound for pid 200, got %q", patterns[200])
	}
}

