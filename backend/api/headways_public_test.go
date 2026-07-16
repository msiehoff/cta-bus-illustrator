package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func TestPublicHeadwayRoutes(t *testing.T) {
	repo := &fake.HeadwaySummaryRepo{}
	day := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	_ = repo.InsertBatch(context.Background(), []business.HeadwaySummary{
		{
			ServiceDate: day, Grain: business.HeadwayGrainRoute, Method: business.HeadwayMethodEqualStop,
			RouteID: "8", RouteName: "Halsted",
			HeadwaySummaryStats: business.HeadwaySummaryStats{Count: 20, MedianMinutes: 11, AvgWaitMinutes: 5.5, CV: 0.35, MeanMinutes: 11},
		},
		{
			ServiceDate: day, Grain: business.HeadwayGrainServiceDay, Method: business.HeadwayMethodEqualStop,
			HeadwaySummaryStats: business.HeadwaySummaryStats{Count: 100, MedianMinutes: 10, AvgWaitMinutes: 5, CV: 0.3, MeanMinutes: 10},
		},
	})

	a := New(Options{
		HeadwaySummaryRepo: repo,
		HeadwayPublic:      app.NewHeadwayPublicService(repo),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/headways/routes?days=30", nil)
	w := httptest.NewRecorder()
	a.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var list HeadwayRoutesListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Routes) != 1 || list.Routes[0].RouteID != "8" {
		t.Fatalf("unexpected routes: %+v", list.Routes)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/headways/routes/8?days=30", nil)
	w = httptest.NewRecorder()
	a.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("route status = %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/headways/system?days=30", nil)
	w = httptest.NewRecorder()
	a.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("system status = %d body=%s", w.Code, w.Body.String())
	}
}
