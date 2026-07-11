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
	headwayRollup  *app.HeadwayRollup
	adminAuth      *AdminAuth
	jobTokenAuth   *JobTokenAuth
}

type Options struct {
	RouteService   *app.RouteService
	CtaDataSrc     app.RouteSegmentDataSource
	PipelineRunner *app.PipelineRunner
	ArrivalRepo    app.ArrivalRepository
	HeadwayRollup  *app.HeadwayRollup
	AdminAuth      *AdminAuth
	JobTokenAuth   *JobTokenAuth
}

func New(opts Options) *API {
	a := &API{
		router:         gin.Default(),
		routeService:   opts.RouteService,
		ctaDataSrc:     opts.CtaDataSrc,
		pipelineRunner: opts.PipelineRunner,
		arrivalRepo:    opts.ArrivalRepo,
		headwayRollup:  opts.HeadwayRollup,
		adminAuth:      opts.AdminAuth,
		jobTokenAuth:   opts.JobTokenAuth,
	}
	a.registerRoutes()
	a.registerAdminRoutes()
	return a
}

func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}
