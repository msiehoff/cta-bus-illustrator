package api

import (
	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
)

type API struct {
	router         *gin.Engine
	routeService   *app.RouteService
	ctaDataSrc     app.RouteSegmentDataSource
	pipelineRunner *app.PipelineRunner
	arrivalRepo    app.ArrivalRepository
	adminAuth      *AdminAuth
}

type Options struct {
	RouteService   *app.RouteService
	CtaDataSrc     app.RouteSegmentDataSource
	PipelineRunner *app.PipelineRunner
	ArrivalRepo    app.ArrivalRepository
	AdminAuth      *AdminAuth
}

func New(opts Options) *API {
	a := &API{
		router:         gin.Default(),
		routeService:   opts.RouteService,
		ctaDataSrc:     opts.CtaDataSrc,
		pipelineRunner: opts.PipelineRunner,
		arrivalRepo:    opts.ArrivalRepo,
		adminAuth:      opts.AdminAuth,
	}
	a.registerRoutes()
	a.registerAdminRoutes()
	return a
}

func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}
