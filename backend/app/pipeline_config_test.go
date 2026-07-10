package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type stubRouteProvider struct {
	routeIDs []string
	err      error
}

func (s stubRouteProvider) GetRouteIDs(context.Context) ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.routeIDs, nil
}

type stubRidershipRepo struct {
	month   time.Time
	byMonth map[business.RidershipType]map[string]*business.RidershipRecord
}

func (s stubRidershipRepo) GetLatestMonth() (time.Time, error) { return s.month, nil }
func (s stubRidershipRepo) GetAvailableMonths() ([]time.Time, error) {
	return []time.Time{s.month}, nil
}
func (s stubRidershipRepo) GetByMonth(_ time.Time, ridershipType business.RidershipType) (map[string]*business.RidershipRecord, error) {
	return s.byMonth[ridershipType], nil
}
func (s stubRidershipRepo) GetAllByRoute(string) ([]business.RidershipRecord, error) {
	return nil, nil
}
func (s stubRidershipRepo) GetSystemTotals() ([]business.RidershipRecord, error) {
	return nil, nil
}
func (s stubRidershipRepo) UpsertBatch([]business.RidershipRecord) error { return nil }

func TestResolvePipelineConfig_UsesEnvOverride(t *testing.T) {
	t.Setenv("PIPELINE_ROUTES", "8, 66, 77")
	t.Setenv("PIPELINE_POLL_INTERVAL", "45s")

	cfg, err := ResolvePipelineConfig(t.Context(), stubRouteProvider{routeIDs: []string{"99"}})
	if err != nil {
		t.Fatalf("ResolvePipelineConfig: %v", err)
	}
	if len(cfg.RouteIDs) != 3 || cfg.RouteIDs[0] != "8" || cfg.RouteIDs[2] != "77" {
		t.Fatalf("unexpected routes: %v", cfg.RouteIDs)
	}
	if cfg.PollInterval != 45*time.Second {
		t.Fatalf("unexpected poll interval: %v", cfg.PollInterval)
	}
}

func TestResolvePipelineConfig_UsesProviderWhenEnvUnset(t *testing.T) {
	t.Setenv("PIPELINE_ROUTES", "")

	cfg, err := ResolvePipelineConfig(t.Context(), stubRouteProvider{routeIDs: []string{"8", "66"}})
	if err != nil {
		t.Fatalf("ResolvePipelineConfig: %v", err)
	}
	if len(cfg.RouteIDs) != 2 || cfg.RouteIDs[0] != "8" {
		t.Fatalf("unexpected routes: %v", cfg.RouteIDs)
	}
}

func TestRidershipRouteProvider_GetRouteIDs(t *testing.T) {
	month := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	repo := stubRidershipRepo{
		month: month,
		byMonth: map[business.RidershipType]map[string]*business.RidershipRecord{
			business.RidershipTypeWeekday: {
				"8":  {RouteExternalID: "8"},
				"66": {RouteExternalID: "66"},
			},
			business.RidershipTypeSaturday: {
				"66": {RouteExternalID: "66"},
			},
			business.RidershipTypeSunday: {
				"77": {RouteExternalID: "77"},
			},
		},
	}

	provider := NewRidershipRouteProvider(repo)
	routeIDs, err := provider.GetRouteIDs(t.Context())
	if err != nil {
		t.Fatalf("GetRouteIDs: %v", err)
	}
	if len(routeIDs) != 3 {
		t.Fatalf("expected 3 routes, got %v", routeIDs)
	}
}

func TestResolvePipelineConfig_ProviderError(t *testing.T) {
	t.Setenv("PIPELINE_ROUTES", "")

	_, err := ResolvePipelineConfig(t.Context(), stubRouteProvider{err: errors.New("boom")})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPipelineEnabledFromEnv(t *testing.T) {
	t.Setenv("PIPELINE_ENABLED", "")
	if PipelineEnabledFromEnv() {
		t.Fatal("expected false by default")
	}

	t.Setenv("PIPELINE_ENABLED", "true")
	if !PipelineEnabledFromEnv() {
		t.Fatal("expected true")
	}
}
