package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func TestAdminAuth_LoginAndSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := &AdminAuth{
		username: "admin",
		password: "secret",
		secret:   []byte("test-secret"),
	}

	router := gin.New()
	router.POST("/admin/login", func(c *gin.Context) {
		var req adminLoginRequest
		_ = c.ShouldBindJSON(&req)
		token, err := auth.Login(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		auth.setSessionCookie(c, token)
		c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: true, Username: req.Username})
	})
	router.GET("/admin/session", func(c *gin.Context) {
		if auth.ValidateToken(auth.sessionToken(c)) {
			c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: true, Username: auth.username})
			return
		}
		c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: false})
	})

	loginReq := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"username":"admin","password":"secret"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d", loginRec.Code)
	}

	cookies := loginRec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	sessionReq := httptest.NewRequest(http.MethodGet, "/admin/session", nil)
	sessionReq.AddCookie(cookies[0])
	sessionRec := httptest.NewRecorder()
	router.ServeHTTP(sessionRec, sessionReq)

	if sessionRec.Code != http.StatusOK {
		t.Fatalf("session status = %d", sessionRec.Code)
	}
}

func TestAdminPipelineStatusRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := &AdminAuth{
		username: "admin",
		password: "secret",
		secret:   []byte("test-secret"),
	}
	arrivalRepo := &fake.ArrivalRepo{}
	api := New(Options{
		RouteService: app.NewRouteService(&fake.RouteRepo{}, &fake.RidershipRepo{}),
		ArrivalRepo:  arrivalRepo,
		AdminAuth:    auth,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/pipeline/status", nil)
	rec := httptest.NewRecorder()
	api.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAdminListArrivals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := &AdminAuth{
		username: "admin",
		password: "secret",
		secret:   []byte("test-secret"),
	}
	arrivalRepo := &fake.ArrivalRepo{}
	_ = arrivalRepo.SaveArrival(t.Context(), business.Arrival{
		StopID: "1", RouteID: "8", Direction: "Northbound", VehicleID: "100",
	})

	api := New(Options{
		RouteService: app.NewRouteService(&fake.RouteRepo{}, &fake.RidershipRepo{}),
		ArrivalRepo:  arrivalRepo,
		AdminAuth:    auth,
	})

	token, err := auth.Login("admin", "secret")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/arrivals", nil)
	req.AddCookie(&http.Cookie{Name: adminSessionCookie, Value: token})
	rec := httptest.NewRecorder()
	api.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}
