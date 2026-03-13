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
