package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
)

func (a *API) registerHeadwayPublicRoutes(v1 *gin.RouterGroup) {
	if a.headwayPublic == nil {
		return
	}
	v1.GET("/headways/routes", a.handleGetHeadwayRoutes)
	v1.GET("/headways/routes/:externalId", a.handleGetHeadwayRoute)
	v1.GET("/headways/system", a.handleGetHeadwaySystem)
}

type HeadwayPeriodStatsResponse struct {
	Count          int     `json:"count"`
	MeanMinutes    float64 `json:"meanMinutes"`
	MedianMinutes  float64 `json:"medianMinutes"`
	StdDevMinutes  float64 `json:"stdDevMinutes"`
	CV             float64 `json:"cv"`
	AvgWaitMinutes float64 `json:"avgWaitMinutes"`
	DaysWithData   int     `json:"daysWithData"`
	PeriodStart    string  `json:"periodStart,omitempty"`
	PeriodEnd      string  `json:"periodEnd,omitempty"`
}

type HeadwayRoutePeriodResponse struct {
	RouteID   string `json:"routeId"`
	RouteName string `json:"routeName,omitempty"`
	HeadwayPeriodStatsResponse
}

type HeadwayDayPointResponse struct {
	ServiceDate    string  `json:"serviceDate"`
	MedianMinutes  float64 `json:"medianMinutes"`
	AvgWaitMinutes float64 `json:"avgWaitMinutes"`
	CV             float64 `json:"cv"`
	Count          int     `json:"count"`
}

type HeadwayRoutesListResponse struct {
	Period HeadwayPeriodStatsResponse  `json:"period"`
	Routes []HeadwayRoutePeriodResponse `json:"routes"`
	Method string                      `json:"method"`
	Grain  string                      `json:"grain"`
	Days   int                         `json:"days"`
}

type HeadwayRouteDetailResponse struct {
	Route  HeadwayRoutePeriodResponse `json:"route"`
	Series []HeadwayDayPointResponse  `json:"series"`
	Method string                     `json:"method"`
	Grain  string                     `json:"grain"`
	Days   int                        `json:"days"`
}

type HeadwaySystemResponse struct {
	Period       HeadwayPeriodStatsResponse   `json:"period"`
	Series       []HeadwayDayPointResponse    `json:"series"`
	LongestWaits []HeadwayRoutePeriodResponse `json:"longestWaits"`
	Method       string                       `json:"method"`
	Grain        string                       `json:"grain"`
	Days         int                          `json:"days"`
}

func (a *API) handleGetHeadwayRoutes(c *gin.Context) {
	days := parseHeadwayDays(c)
	routes, period, err := a.headwayPublic.ListRoutes(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, HeadwayRoutesListResponse{
		Period: toPeriodStatsResponse(period),
		Routes: toRoutePeriodResponses(routes),
		Method: "equal_stop",
		Grain:  "route",
		Days:   app.ResolvePeriodDays(days),
	})
}

func (a *API) handleGetHeadwayRoute(c *gin.Context) {
	routeID := c.Param("externalId")
	days := parseHeadwayDays(c)
	detail, err := a.headwayPublic.GetRoute(c.Request.Context(), routeID, days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, HeadwayRouteDetailResponse{
		Route:  toRoutePeriodResponse(detail.HeadwayRoutePeriod),
		Series: toDayPointResponses(detail.Series),
		Method: "equal_stop",
		Grain:  "route",
		Days:   app.ResolvePeriodDays(days),
	})
}

func (a *API) handleGetHeadwaySystem(c *gin.Context) {
	days := parseHeadwayDays(c)
	overview, err := a.headwayPublic.GetSystem(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, HeadwaySystemResponse{
		Period:       toPeriodStatsResponse(overview.HeadwayPeriodStats),
		Series:       toDayPointResponses(overview.Series),
		LongestWaits: toRoutePeriodResponses(overview.LongestWaits),
		Method:       "equal_stop",
		Grain:        "service_day",
		Days:         app.ResolvePeriodDays(days),
	})
}

func parseHeadwayDays(c *gin.Context) int {
	raw := c.Query("days")
	if raw == "" {
		return app.DefaultHeadwayPeriodDays
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return app.DefaultHeadwayPeriodDays
	}
	return n
}

func toPeriodStatsResponse(s app.HeadwayPeriodStats) HeadwayPeriodStatsResponse {
	resp := HeadwayPeriodStatsResponse{
		Count:          s.Count,
		MeanMinutes:    round2(s.MeanMinutes),
		MedianMinutes:  round2(s.MedianMinutes),
		StdDevMinutes:  round2(s.StdDevMinutes),
		CV:             round3(s.CV),
		AvgWaitMinutes: round2(s.AvgWaitMinutes),
		DaysWithData:   s.DaysWithData,
	}
	if !s.PeriodStart.IsZero() {
		resp.PeriodStart = s.PeriodStart.UTC().Format("2006-01-02")
	}
	if !s.PeriodEnd.IsZero() {
		resp.PeriodEnd = s.PeriodEnd.UTC().Format("2006-01-02")
	}
	return resp
}

func toRoutePeriodResponse(r app.HeadwayRoutePeriod) HeadwayRoutePeriodResponse {
	return HeadwayRoutePeriodResponse{
		RouteID:                    r.RouteID,
		RouteName:                  r.RouteName,
		HeadwayPeriodStatsResponse: toPeriodStatsResponse(r.HeadwayPeriodStats),
	}
}

func toRoutePeriodResponses(routes []app.HeadwayRoutePeriod) []HeadwayRoutePeriodResponse {
	out := make([]HeadwayRoutePeriodResponse, len(routes))
	for i, r := range routes {
		out[i] = toRoutePeriodResponse(r)
	}
	return out
}

func toDayPointResponses(points []app.HeadwayDayPoint) []HeadwayDayPointResponse {
	out := make([]HeadwayDayPointResponse, len(points))
	for i, p := range points {
		out[i] = HeadwayDayPointResponse{
			ServiceDate:    p.ServiceDate.UTC().Format("2006-01-02"),
			MedianMinutes:  round2(p.MedianMinutes),
			AvgWaitMinutes: round2(p.AvgWaitMinutes),
			CV:             round3(p.CV),
			Count:          p.Count,
		}
	}
	return out
}
