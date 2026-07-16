package api

import (
	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
)

type API struct {
	router             *gin.Engine
	routeService       *app.RouteService
	ctaDataSrc         app.RouteSegmentDataSource
	pipelineRunner     *app.PipelineRunner
	arrivalRepo        app.ArrivalRepository
	headwayRepo        app.HeadwayRepository
	headwaySummaryRepo app.HeadwaySummaryRepository
	headwayRollup      *app.HeadwayRollup
	headwayPublic      *app.HeadwayPublicService
	adminAuth          *AdminAuth
	jobTokenAuth       *JobTokenAuth
}

type Options struct {
	RouteService       *app.RouteService
	CtaDataSrc         app.RouteSegmentDataSource
	PipelineRunner     *app.PipelineRunner
	ArrivalRepo        app.ArrivalRepository
	HeadwayRepo        app.HeadwayRepository
	HeadwaySummaryRepo app.HeadwaySummaryRepository
	HeadwayRollup      *app.HeadwayRollup
	HeadwayPublic      *app.HeadwayPublicService
	AdminAuth          *AdminAuth
	JobTokenAuth       *JobTokenAuth
}

func New(opts Options) *API {
	a := &API{
		router:             gin.Default(),
		routeService:       opts.RouteService,
		ctaDataSrc:         opts.CtaDataSrc,
		pipelineRunner:     opts.PipelineRunner,
		arrivalRepo:        opts.ArrivalRepo,
		headwayRepo:        opts.HeadwayRepo,
		headwaySummaryRepo: opts.HeadwaySummaryRepo,
		headwayRollup:      opts.HeadwayRollup,
		headwayPublic:      opts.HeadwayPublic,
		adminAuth:          opts.AdminAuth,
		jobTokenAuth:       opts.JobTokenAuth,
	}
	if a.headwayPublic == nil && opts.HeadwaySummaryRepo != nil {
		a.headwayPublic = app.NewHeadwayPublicService(opts.HeadwaySummaryRepo)
	}
	a.registerRoutes()
	a.registerAdminRoutes()
	return a
}

func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}
