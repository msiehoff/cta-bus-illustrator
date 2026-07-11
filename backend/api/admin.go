package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
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
			protected.GET("/headways", a.handleAdminListHeadways)
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
