package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ResolvePipelineConfig(ctx context.Context, provider PipelineRouteProvider) (PipelineConfig, error) {
	cfg := PipelineConfig{
		PollInterval: defaultPollInterval(),
	}

	if v := strings.TrimSpace(os.Getenv("PIPELINE_ROUTES")); v != "" {
		cfg.RouteIDs = splitCommaList(v)
		return cfg, nil
	}

	if provider == nil {
		return cfg, fmt.Errorf("PIPELINE_ROUTES unset and no route provider configured")
	}

	routeIDs, err := provider.GetRouteIDs(ctx)
	if err != nil {
		return cfg, fmt.Errorf("resolve pipeline routes: %w", err)
	}
	if len(routeIDs) == 0 {
		return cfg, fmt.Errorf("no routes found for pipeline")
	}

	cfg.RouteIDs = routeIDs
	return cfg, nil
}

func defaultPollInterval() time.Duration {
	interval := 30 * time.Second
	if v := strings.TrimSpace(os.Getenv("PIPELINE_POLL_INTERVAL")); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}
	return interval
}

func PipelineEnabledFromEnv() bool {
	return envBool("PIPELINE_ENABLED", false)
}

func PipelineUseFakeCTAFromEnv() bool {
	return envBool("PIPELINE_USE_FAKE_CTA", false)
}

func envBool(key string, defaultVal bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

func splitCommaList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
