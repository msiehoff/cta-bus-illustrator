package app

import (
	"os"
	"testing"
	"time"
)

func TestPipelineConfigFromEnv(t *testing.T) {
	t.Setenv("PIPELINE_ROUTES", "8, 66, 77")
	t.Setenv("PIPELINE_POLL_INTERVAL", "45s")

	cfg := PipelineConfigFromEnv()
	if len(cfg.RouteIDs) != 3 || cfg.RouteIDs[0] != "8" || cfg.RouteIDs[2] != "77" {
		t.Fatalf("unexpected routes: %v", cfg.RouteIDs)
	}
	if cfg.PollInterval != 45*time.Second {
		t.Fatalf("unexpected poll interval: %v", cfg.PollInterval)
	}
}

func TestPipelineEnabledFromEnv(t *testing.T) {
	os.Unsetenv("PIPELINE_ENABLED")
	if PipelineEnabledFromEnv() {
		t.Fatal("expected false by default")
	}

	t.Setenv("PIPELINE_ENABLED", "true")
	if !PipelineEnabledFromEnv() {
		t.Fatal("expected true")
	}
}
