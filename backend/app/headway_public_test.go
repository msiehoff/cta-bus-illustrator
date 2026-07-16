package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func TestAggregatePeriodStats(t *testing.T) {
	day1 := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	rows := []business.HeadwaySummary{
		{
			ServiceDate: day1,
			HeadwaySummaryStats: business.HeadwaySummaryStats{
				Count: 10, MeanMinutes: 10, MedianMinutes: 9, StdDevMinutes: 2, CV: 0.2, AvgWaitMinutes: 5,
			},
		},
		{
			ServiceDate: day2,
			HeadwaySummaryStats: business.HeadwaySummaryStats{
				Count: 30, MeanMinutes: 14, MedianMinutes: 13, StdDevMinutes: 4, CV: 0.4, AvgWaitMinutes: 7,
			},
		},
	}

	got := app.AggregatePeriodStats(rows)
	if got.Count != 40 {
		t.Fatalf("count = %d", got.Count)
	}
	if got.DaysWithData != 2 {
		t.Fatalf("days = %d", got.DaysWithData)
	}
	// (10*10 + 14*30) / 40 = 13
	if got.MeanMinutes < 12.9 || got.MeanMinutes > 13.1 {
		t.Fatalf("mean = %v", got.MeanMinutes)
	}
}

func TestHeadwayPublicListRoutes(t *testing.T) {
	repo := &fake.HeadwaySummaryRepo{}
	day1 := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	_ = repo.InsertBatch(context.Background(), []business.HeadwaySummary{
		{
			ServiceDate: day1, Grain: business.HeadwayGrainRoute, Method: business.HeadwayMethodEqualStop,
			RouteID: "8", RouteName: "Halsted",
			HeadwaySummaryStats: business.HeadwaySummaryStats{Count: 20, MeanMinutes: 10, MedianMinutes: 10, AvgWaitMinutes: 5, CV: 0.3},
		},
		{
			ServiceDate: day2, Grain: business.HeadwayGrainRoute, Method: business.HeadwayMethodEqualStop,
			RouteID: "8", RouteName: "Halsted",
			HeadwaySummaryStats: business.HeadwaySummaryStats{Count: 20, MeanMinutes: 12, MedianMinutes: 12, AvgWaitMinutes: 6, CV: 0.4},
		},
		{
			ServiceDate: day2, Grain: business.HeadwayGrainRoute, Method: business.HeadwayMethodEqualStop,
			RouteID: "22", RouteName: "Clark",
			HeadwaySummaryStats: business.HeadwaySummaryStats{Count: 15, MeanMinutes: 8, MedianMinutes: 8, AvgWaitMinutes: 4, CV: 0.2},
		},
	})

	svc := app.NewHeadwayPublicService(repo)
	routes, period, err := svc.ListRoutes(context.Background(), 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(routes) != 2 {
		t.Fatalf("routes = %d", len(routes))
	}
	if period.DaysWithData != 2 {
		t.Fatalf("period days = %d", period.DaysWithData)
	}
	// Route 8 should rank first (higher median after weighting).
	if routes[0].RouteID != "8" {
		t.Fatalf("expected route 8 first, got %s", routes[0].RouteID)
	}

	detail, err := svc.GetRoute(context.Background(), "8", 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(detail.Series) != 2 {
		t.Fatalf("series len = %d", len(detail.Series))
	}
	if detail.RouteName != "Halsted" {
		t.Fatalf("name = %q", detail.RouteName)
	}
}
