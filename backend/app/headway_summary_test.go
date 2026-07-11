package app_test

import (
	"math"
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func TestSummarizeMinutes(t *testing.T) {
	stats := app.SummarizeMinutes([]float64{10, 12, 14, 20})
	if stats.Count != 4 {
		t.Fatalf("count = %d", stats.Count)
	}
	if math.Abs(stats.MeanMinutes-14) > 1e-9 {
		t.Errorf("mean = %v want 14", stats.MeanMinutes)
	}
	if math.Abs(stats.MedianMinutes-13) > 1e-9 {
		t.Errorf("median = %v want 13", stats.MedianMinutes)
	}
	if stats.AvgWaitMinutes != stats.MeanMinutes/2 {
		t.Errorf("avg wait = %v", stats.AvgWaitMinutes)
	}
	if stats.CV <= 0 {
		t.Errorf("expected positive CV, got %v", stats.CV)
	}
}

func TestSummarizeMinutesEmpty(t *testing.T) {
	stats := app.SummarizeMinutes(nil)
	if stats.Count != 0 || stats.MeanMinutes != 0 {
		t.Fatalf("unexpected %+v", stats)
	}
}

func TestSummarizeHeadwaysByStopAndMeanOfMeans(t *testing.T) {
	base := time.Date(2026, 7, 10, 8, 0, 0, 0, time.UTC)
	headways := []business.Headway{
		{StopID: "a", StopName: "A", RouteID: "8", Direction: "NB", Timestamp: base, HeadwayMinutes: 10},
		{StopID: "a", StopName: "A", RouteID: "8", Direction: "NB", Timestamp: base.Add(time.Hour), HeadwayMinutes: 14},
		{StopID: "b", StopName: "B", RouteID: "8", Direction: "NB", Timestamp: base, HeadwayMinutes: 20},
		{StopID: "b", StopName: "B", RouteID: "8", Direction: "NB", Timestamp: base.Add(time.Hour), HeadwayMinutes: 20},
	}

	byStop := app.SummarizeHeadwaysByStop(headways)
	if len(byStop) != 2 {
		t.Fatalf("stops = %d", len(byStop))
	}

	// Stop B mean 20 > stop A mean 12 → B first
	if byStop[0].StopID != "b" {
		t.Errorf("expected stop b first by mean, got %s", byStop[0].StopID)
	}

	agg := app.MeanOfStopMeans(byStop)
	// mean of means: (12 + 20) / 2 = 16
	if math.Abs(agg.MeanMinutes-16) > 1e-9 {
		t.Errorf("mean of means = %v want 16", agg.MeanMinutes)
	}
	if agg.Count != 4 {
		t.Errorf("count = %d", agg.Count)
	}
}
