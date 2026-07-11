package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type adminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) registerAdminRoutes() {
	admin := a.router.Group("/api/v1/admin")

	if a.adminAuth != nil {
		admin.POST("/login", a.handleAdminLogin)
		admin.POST("/logout", a.adminAuth.Middleware(), a.handleAdminLogout)
		admin.GET("/session", a.handleAdminSession)

		protected := admin.Group("", a.adminAuth.Middleware())
		{
			protected.GET("/pipeline/status", a.handleAdminPipelineStatus)
			protected.GET("/arrivals", a.handleAdminListArrivals)
			protected.GET("/headways/summary", a.handleAdminHeadwaySummary)
			protected.GET("/headways", a.handleAdminListHeadways)
			protected.GET("/headway-summaries", a.handleAdminListHeadwaySummaries)
		}
	}

	a.registerHeadwayAdminRoutes(admin)
}

func (a *API) handleAdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.adminAuth.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	a.adminAuth.setSessionCookie(c, token)
	c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: true, Username: req.Username})
}

func (a *API) handleAdminLogout(c *gin.Context) {
	a.adminAuth.clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) handleAdminSession(c *gin.Context) {
	if a.adminAuth.ValidateToken(a.adminAuth.sessionToken(c)) {
		c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: true, Username: a.adminAuth.username})
		return
	}
	c.JSON(http.StatusOK, AdminSessionResponse{Authenticated: false})
}

func (a *API) handleAdminPipelineStatus(c *gin.Context) {
	var status app.PipelineStatus
	enabled := a.pipelineRunner != nil
	if enabled {
		status = a.pipelineRunner.Status()
	}

	var arrivalCount int64
	if a.arrivalRepo != nil {
		count, err := a.arrivalRepo.CountArrivals(c.Request.Context(), app.ArrivalFilter{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		arrivalCount = count
	}

	c.JSON(http.StatusOK, toPipelineStatusResponse(status, enabled, arrivalCount))
}

func (a *API) handleAdminListArrivals(c *gin.Context) {
	if a.arrivalRepo == nil {
		c.JSON(http.StatusOK, ListArrivalsResponse{
			Arrivals: []ArrivalResponse{},
			Total:    0,
			Limit:    50,
			Offset:   0,
		})
		return
	}

	limit := parseIntQuery(c, "limit", 50, 200)
	offset := parseIntQuery(c, "offset", 0, 1_000_000)
	filter := app.ArrivalFilter{
		RouteID:   c.Query("route"),
		Direction: c.Query("direction"),
		Stop:      c.Query("stop"),
		VehicleID: c.Query("vehicle"),
		SortAsc:   c.Query("sort") == "asc",
		Limit:     limit,
		Offset:    offset,
	}

	if date := c.Query("date"); date != "" {
		serviceDate, err := app.ParseServiceDate(date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		start, end := app.ServiceDateBounds(serviceDate)
		filter.From = &start
		filter.To = &end
	}

	arrivals, err := a.arrivalRepo.ListArrivals(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := a.arrivalRepo.CountArrivals(c.Request.Context(), app.ArrivalFilter{
		RouteID:   filter.RouteID,
		Direction: filter.Direction,
		Stop:      filter.Stop,
		VehicleID: filter.VehicleID,
		From:      filter.From,
		To:        filter.To,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ListArrivalsResponse{
		Arrivals: make([]ArrivalResponse, len(arrivals)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}
	for i, arrival := range arrivals {
		resp.Arrivals[i] = toArrivalResponse(arrival)
	}

	c.JSON(http.StatusOK, resp)
}

func (a *API) handleAdminListHeadways(c *gin.Context) {
	if a.headwayRepo == nil {
		c.JSON(http.StatusOK, ListHeadwaysResponse{
			Headways: []HeadwayResponse{},
			Total:    0,
			Limit:    50,
			Offset:   0,
		})
		return
	}

	limit := parseIntQuery(c, "limit", 50, 200)
	offset := parseIntQuery(c, "offset", 0, 1_000_000)
	filter := app.HeadwayListFilter{
		RouteID:   c.Query("route"),
		Direction: c.Query("direction"),
		Stop:      c.Query("stop"),
		VehicleID: c.Query("vehicle"),
		SortAsc:   c.Query("sort") == "asc",
		Limit:     limit,
		Offset:    offset,
	}

	if date := c.Query("date"); date != "" {
		serviceDate, err := app.ParseServiceDate(date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		start, end := app.ServiceDateBounds(serviceDate)
		filter.From = &start
		filter.To = &end
	}

	headways, err := a.headwayRepo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := a.headwayRepo.Count(c.Request.Context(), app.HeadwayListFilter{
		RouteID:   filter.RouteID,
		Direction: filter.Direction,
		Stop:      filter.Stop,
		VehicleID: filter.VehicleID,
		From:      filter.From,
		To:        filter.To,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ListHeadwaysResponse{
		Headways: make([]HeadwayResponse, len(headways)),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}
	for i, h := range headways {
		resp.Headways[i] = toHeadwayResponse(h)
	}

	c.JSON(http.StatusOK, resp)
}

func (a *API) handleAdminHeadwaySummary(c *gin.Context) {
	date := c.Query("date")
	vehicle := c.Query("vehicle")

	// Prefer persisted summaries when we have a service date and no vehicle filter
	// (vehicle is only on raw observed gaps).
	if date != "" && vehicle == "" && a.headwayRollup != nil {
		serviceDate, err := app.ParseServiceDate(date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		serviceDateOnly := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, time.UTC)
		pooled, equal, byStop, ok, err := a.headwayRollup.LoadStoredSummary(c.Request.Context(), app.HeadwaySummaryFilter{
			ServiceDate: serviceDateOnly,
			RouteID:     c.Query("route"),
			Direction:   c.Query("direction"),
			Stop:        c.Query("stop"),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if ok {
			c.JSON(http.StatusOK, toHeadwaySummaryResponse(pooled, equal, byStop, "stored"))
			return
		}
	}

	// Fallback: compute on read from observed headways (no date yet / job not run / vehicle filter).
	if a.headwayRepo == nil {
		c.JSON(http.StatusOK, HeadwaySummaryResponse{
			ByStop: []HeadwayStopSummaryResponse{},
			Source: "computed",
		})
		return
	}

	filter := app.HeadwayListFilter{
		RouteID:   c.Query("route"),
		Direction: c.Query("direction"),
		Stop:      c.Query("stop"),
		VehicleID: vehicle,
	}
	if date != "" {
		serviceDate, err := app.ParseServiceDate(date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		start, end := app.ServiceDateBounds(serviceDate)
		filter.From = &start
		filter.To = &end
	}

	headways, err := a.headwayRepo.ListAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	byStop := app.SummarizeHeadwaysByStop(headways)
	c.JSON(http.StatusOK, toHeadwaySummaryResponse(
		app.SummarizeHeadways(headways),
		app.MeanOfStopMeans(byStop),
		byStop,
		"computed",
	))
}

func toHeadwaySummaryResponse(
	pooled, equal business.HeadwaySummaryStats,
	byStop []business.HeadwayStopSummary,
	source string,
) HeadwaySummaryResponse {
	resp := HeadwaySummaryResponse{
		Pooled:          toSummaryStatsResponse(pooled),
		EqualStopWeight: toSummaryStatsResponse(equal),
		ByStop:          make([]HeadwayStopSummaryResponse, len(byStop)),
		Source:          source,
	}
	for i, s := range byStop {
		resp.ByStop[i] = HeadwayStopSummaryResponse{
			StopID:                      s.StopID,
			StopName:                    s.StopName,
			RouteID:                     s.RouteID,
			RouteName:                   s.RouteName,
			Direction:                   s.Direction,
			HeadwaySummaryStatsResponse: toSummaryStatsResponse(s.HeadwaySummaryStats),
		}
	}
	return resp
}

func (a *API) handleAdminListHeadwaySummaries(c *gin.Context) {
	if a.headwaySummaryRepo == nil {
		c.JSON(http.StatusOK, ListHeadwaySummariesResponse{
			Summaries: []HeadwaySummaryRowResponse{},
			Total:     0,
			Limit:     50,
			Offset:    0,
		})
		return
	}

	limit := parseIntQuery(c, "limit", 50, 200)
	offset := parseIntQuery(c, "offset", 0, 1_000_000)
	filter := app.HeadwaySummaryFilter{
		Grain:     c.Query("grain"),
		Method:    c.Query("method"),
		RouteID:   c.Query("route"),
		Direction: c.Query("direction"),
		Stop:      c.Query("stop"),
		SortAsc:   c.Query("sort") == "asc",
		Limit:     limit,
		Offset:    offset,
	}

	if date := c.Query("date"); date != "" {
		serviceDate, err := app.ParseServiceDate(date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter.ServiceDate = time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, time.UTC)
	}

	summaries, err := a.headwaySummaryRepo.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := a.headwaySummaryRepo.Count(c.Request.Context(), app.HeadwaySummaryFilter{
		ServiceDate: filter.ServiceDate,
		Grain:       filter.Grain,
		Method:      filter.Method,
		RouteID:     filter.RouteID,
		Direction:   filter.Direction,
		Stop:        filter.Stop,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ListHeadwaySummariesResponse{
		Summaries: make([]HeadwaySummaryRowResponse, len(summaries)),
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}
	for i, s := range summaries {
		resp.Summaries[i] = toHeadwaySummaryRowResponse(s)
	}
	c.JSON(http.StatusOK, resp)
}

func toHeadwaySummaryRowResponse(s business.HeadwaySummary) HeadwaySummaryRowResponse {
	return HeadwaySummaryRowResponse{
		ServiceDate:    s.ServiceDate.UTC().Format("2006-01-02"),
		Grain:          s.Grain,
		Method:         s.Method,
		StopID:         s.StopID,
		StopName:       s.StopName,
		RouteID:        s.RouteID,
		RouteName:      s.RouteName,
		Direction:      s.Direction,
		Count:          s.Count,
		MeanMinutes:    round2(s.MeanMinutes),
		MedianMinutes:  round2(s.MedianMinutes),
		StdDevMinutes:  round2(s.StdDevMinutes),
		CV:             round3(s.CV),
		AvgWaitMinutes: round2(s.AvgWaitMinutes),
	}
}
