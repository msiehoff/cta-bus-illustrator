package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
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
		v1.POST("/routes/import-segments", a.handleImportRouteSegments)
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

func (a *API) handleImportRouteSegments(c *gin.Context) {
	if err := a.routeService.ImportRouteSegments(c.Request.Context(), a.ctaDataSrc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// resolveMonth returns the parsed month from the query string, or falls back to
// the latest available ridership month. Returns the zero time if neither is available.
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
		// No ridership data available — return zero time so routes are still
		// returned but with nil ridership fields.
		return time.Time{}, nil
	}
	return latest, nil
}

// resolveRidershipType validates the type param and defaults to weekday.
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
