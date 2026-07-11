package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
)

type headwayRunRequest struct {
	ServiceDate string `json:"service_date"`
}

func (a *API) registerHeadwayAdminRoutes(admin *gin.RouterGroup) {
	if a.headwayRollup == nil {
		return
	}
	if a.adminAuth == nil && a.jobTokenAuth == nil {
		return
	}

	headway := admin.Group("", a.headwayAuthMiddleware())
	{
		headway.POST("/headways/run", a.handleAdminHeadwayRun)
		headway.GET("/headways/runs", a.handleAdminListHeadwayRuns)
	}
}

func (a *API) handleAdminHeadwayRun(c *gin.Context) {
	var req headwayRunRequest
	_ = c.ShouldBindJSON(&req)

	serviceDate := app.YesterdayServiceDate()
	if req.ServiceDate != "" {
		parsed, err := app.ParseServiceDate(req.ServiceDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		serviceDate = parsed
	}

	result, err := a.headwayRollup.Run(c.Request.Context(), serviceDate, headwayTriggerFromContext(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, HeadwayJobRunResponseFrom(result.Run))
		return
	}

	c.JSON(http.StatusOK, HeadwayJobRunResponseFrom(result.Run))
}

func (a *API) handleAdminListHeadwayRuns(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 50, 200)
	offset := parseIntQuery(c, "offset", 0, 1_000_000)

	runs, err := a.headwayRollup.ListJobRuns(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := ListHeadwayJobRunsResponse{
		Runs:   make([]HeadwayJobRunResponse, len(runs)),
		Limit:  limit,
		Offset: offset,
	}
	for i, run := range runs {
		resp.Runs[i] = HeadwayJobRunResponseFrom(run)
	}
	c.JSON(http.StatusOK, resp)
}
