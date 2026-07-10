package app

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func PipelineConfigFromEnv() PipelineConfig {
	routes := []string{"8", "66"}
	if v := strings.TrimSpace(os.Getenv("PIPELINE_ROUTES")); v != "" {
		routes = splitCommaList(v)
	}

	interval := 30 * time.Second
	if v := strings.TrimSpace(os.Getenv("PIPELINE_POLL_INTERVAL")); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			interval = d
		}
	}

	return PipelineConfig{
		RouteIDs:     routes,
		PollInterval: interval,
	}
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
