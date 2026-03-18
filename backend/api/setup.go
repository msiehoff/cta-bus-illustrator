package api

import (
	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
)

type API struct {
	router       *gin.Engine
	routeService *app.RouteService
	ctaDataSrc   app.RouteSegmentDataSource
}

func New(routeService *app.RouteService, ctaDataSrc app.RouteSegmentDataSource) *API {
	a := &API{
		router:       gin.Default(),
		routeService: routeService,
		ctaDataSrc:   ctaDataSrc,
	}
	a.registerRoutes()
	return a
}

func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}
