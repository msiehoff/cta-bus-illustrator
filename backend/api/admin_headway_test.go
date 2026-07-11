package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func TestHeadwayRunWithJobToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	arrivalRepo := &fake.ArrivalRepo{}
	base := time.Date(2026, 7, 10, 8, 0, 0, 0, app.ChicagoLocation())
	_ = arrivalRepo.SaveArrival(t.Context(), business.Arrival{
		StopID: "s1", RouteID: "8", Direction: "Northbound", VehicleID: "a", Timestamp: base,
	})
	_ = arrivalRepo.SaveArrival(t.Context(), business.Arrival{
		StopID: "s1", RouteID: "8", Direction: "Northbound", VehicleID: "b", Timestamp: base.Add(12 * time.Minute),
	})

	rollup := app.NewHeadwayRollup(arrivalRepo, &fake.HeadwayRepo{}, &fake.HeadwayJobRunRepo{})
	api := New(Options{
		RouteService:  app.NewRouteService(&fake.RouteRepo{}, &fake.RidershipRepo{}),
		ArrivalRepo:   arrivalRepo,
		HeadwayRollup: rollup,
		JobTokenAuth:  &JobTokenAuth{token: "test-job-token"},
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/headways/run",
		strings.NewReader(`{"service_date":"2026-07-10"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-job-token")
	rec := httptest.NewRecorder()
	api.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var resp HeadwayJobRunResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != string(business.HeadwayJobSuccess) {
		t.Fatalf("status = %s", resp.Status)
	}
	if resp.TriggeredBy != string(business.HeadwayTriggerCron) {
		t.Fatalf("triggeredBy = %s", resp.TriggeredBy)
	}
	if resp.HeadwaysWritten != 1 {
		t.Fatalf("headwaysWritten = %d want 1", resp.HeadwaysWritten)
	}
}

func TestHeadwayRunWithAdminSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := &AdminAuth{
		username: "admin",
		password: "secret",
		secret:   []byte("test-secret"),
	}
	rollup := app.NewHeadwayRollup(&fake.ArrivalRepo{}, &fake.HeadwayRepo{}, &fake.HeadwayJobRunRepo{})
	api := New(Options{
		RouteService:  app.NewRouteService(&fake.RouteRepo{}, &fake.RidershipRepo{}),
		HeadwayRollup: rollup,
		AdminAuth:     auth,
	})

	token, err := auth.Login("admin", "secret")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/headways/run",
		strings.NewReader(`{"service_date":"2026-07-10"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: adminSessionCookie, Value: token})
	rec := httptest.NewRecorder()
	api.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var resp HeadwayJobRunResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.TriggeredBy != string(business.HeadwayTriggerAdmin) {
		t.Fatalf("triggeredBy = %s", resp.TriggeredBy)
	}
}

func TestHeadwayRunUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rollup := app.NewHeadwayRollup(&fake.ArrivalRepo{}, &fake.HeadwayRepo{}, &fake.HeadwayJobRunRepo{})
	api := New(Options{
		RouteService:  app.NewRouteService(&fake.RouteRepo{}, &fake.RidershipRepo{}),
		HeadwayRollup: rollup,
		JobTokenAuth:  &JobTokenAuth{token: "secret"},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/headways/run", nil)
	rec := httptest.NewRecorder()
	api.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
