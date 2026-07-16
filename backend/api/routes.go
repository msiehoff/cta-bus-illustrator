package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/cta"
)

var validRidershipTypes = map[string]business.RidershipType{
	string(business.RidershipTypeWeekday):  business.RidershipTypeWeekday,
	string(business.RidershipTypeSaturday): business.RidershipTypeSaturday,
	string(business.RidershipTypeSunday):   business.RidershipTypeSunday,
}

func (a *API) registerRoutes() {
	v1 := a.router.Group("/api/v1")
	{
		v1.GET("/health", a.handleHealth)
		v1.GET("/routes", a.handleGetRoutes)
		v1.GET("/routes/comparison", a.handleGetRoutesComparison)
		v1.POST("/routes/import-segments", a.handleImportRouteSegments)
		v1.POST("/routes/:externalId/segments", a.handleImportRouteSegmentsJSON)
		v1.GET("/ridership/months", a.handleGetRidershipMonths)
		v1.GET("/ridership/system", a.handleGetSystemRidership)
		v1.GET("/ridership/routes/:externalId", a.handleGetRouteRidership)
		v1.POST("/ridership/import", a.handleImportRidership)
		a.registerHeadwayPublicRoutes(v1)
	}
}

func (a *API) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) handleGetRoutes(c *gin.Context) {
	month, err := resolveMonth(c.Query("month"), func() (time.Time, error) {
		return a.routeService.GetLatestRidershipMonth()
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ridershipType, err := resolveRidershipType(c.Query("type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routes, err := a.routeService.GetRoutesForMonth(month, ridershipType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toGetRoutesResponse(routes))
}

func (a *API) handleGetRoutesComparison(c *gin.Context) {
	ridershipType, err := resolveRidershipType(c.Query("type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	month, err := resolveMonth(c.Query("month"), func() (time.Time, error) {
		return a.routeService.GetLatestRidershipMonth()
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := a.routeService.GetRoutesComparison(ridershipType, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toRoutesComparisonResponse(result))
}

func (a *API) handleImportRouteSegments(c *gin.Context) {
	if err := a.routeService.ImportRouteSegmentsFromSrc(c.Request.Context(), a.ctaDataSrc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) handleImportRouteSegmentsJSON(c *gin.Context) {
	externalID := c.Param("externalId")
	var body cta.GetRoutePatternResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	segments, err := cta.SegmentsFromPatternResponse(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.routeService.ImportRouteSegments(c.Request.Context(), externalID, segments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "segments": len(segments)})
}

func resolveMonth(monthParam string, getLatest func() (time.Time, error)) (time.Time, error) {
	if monthParam != "" {
		parsed, err := time.Parse("2006-01", monthParam)
		if err != nil {
			return time.Time{}, err
		}
		return time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}

	latest, err := getLatest()
	if err != nil {
		return time.Time{}, nil
	}
	return latest, nil
}

func resolveRidershipType(typeParam string) (business.RidershipType, error) {
	if typeParam == "" {
		return business.RidershipTypeWeekday, nil
	}
	t, ok := validRidershipTypes[typeParam]
	if !ok {
		return "", fmt.Errorf("invalid type %q, must be one of: weekday, saturday, sunday", typeParam)
	}
	return t, nil
}
