package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	routes, err := a.routeService.GetRoutes()
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
