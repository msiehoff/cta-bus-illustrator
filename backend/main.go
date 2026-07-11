package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/msiehoff/cta-bus-illustrator/backend/api"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/migrations"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/cta"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
	pgstore "github.com/msiehoff/cta-bus-illustrator/backend/storage/postgres"
	"gorm.io/gorm"
)

func main() {
	app.InitLogLevel()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var (
		routeRepo     app.RouteRepository
		ridershipRepo app.RidershipRepository
		arrivalRepo   app.ArrivalRepository
		db            *gorm.DB
	)

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		var err error
		db, err = pgstore.Connect(dsn)
		if err != nil {
			log.Printf("db connection error: %v", err)
			routeRepo = &fake.RouteRepo{}
			ridershipRepo = &fake.RidershipRepo{}
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("failed to get underlying sql.DB: %v", err)
			} else if err := migrations.Run(sqlDB); err != nil {
				log.Printf("migration error: %v", err)
			}

			routeRepo = pgstore.NewRouteRepo(db)
			ridershipRepo = pgstore.NewRidershipRepo(db)
			arrivalRepo = pgstore.NewArrivalRepo(db)
			log.Println("using postgres repository")
		}
	} else {
		log.Println("DATABASE_URL not set — using fake repository")
		routeRepo = &fake.RouteRepo{}
		ridershipRepo = &fake.RidershipRepo{}
	}

	if arrivalRepo == nil {
		arrivalRepo = &fake.ArrivalRepo{}
	}

	ctaAPIKey := os.Getenv("CTA_API_KEY")
	ctaDataSrc := cta.NewRouteSegmentDataSource(cta.NewClient(ctaAPIKey))

	var pipelineRunner *app.PipelineRunner
	if app.PipelineEnabledFromEnv() {
		var routeProvider app.PipelineRouteProvider
		if db != nil {
			routeProvider = app.NewRidershipRouteProvider(ridershipRepo)
		} else {
			routeProvider = fake.NewPipelineRouteProvider()
		}
		pipelineRunner = startPipeline(ctx, ctaAPIKey, routeProvider, arrivalRepo, db)
	}

	adminAuth, adminEnabled := api.AdminAuthFromEnv()
	if adminEnabled {
		log.Println("admin UI enabled")
	} else {
		log.Println("admin UI disabled — set ADMIN_USERNAME and ADMIN_PASSWORD to enable")
	}

	jobTokenAuth, jobTokenEnabled := api.JobTokenAuthFromEnv()
	if jobTokenEnabled {
		log.Println("headway job token auth enabled")
	}

	var headwayRepo app.HeadwayRepository
	var headwayJobRunRepo app.HeadwayJobRunRepository
	var headwaySummaryRepo app.HeadwaySummaryRepository
	if db != nil {
		headwayRepo = pgstore.NewHeadwayRepo(db)
		headwayJobRunRepo = pgstore.NewHeadwayJobRunRepo(db)
		headwaySummaryRepo = pgstore.NewHeadwaySummaryRepo(db)
	} else {
		headwayRepo = &fake.HeadwayRepo{}
		headwayJobRunRepo = &fake.HeadwayJobRunRepo{}
		headwaySummaryRepo = &fake.HeadwaySummaryRepo{}
	}
	headwayRollup := app.NewHeadwayRollup(arrivalRepo, headwayRepo, headwaySummaryRepo, headwayJobRunRepo)

	routeService := app.NewRouteService(routeRepo, ridershipRepo)
	a := api.New(api.Options{
		RouteService:   routeService,
		CtaDataSrc:     ctaDataSrc,
		PipelineRunner: pipelineRunner,
		ArrivalRepo:    arrivalRepo,
		HeadwayRepo:    headwayRepo,
		HeadwayRollup:  headwayRollup,
		AdminAuth:      adminAuth,
		JobTokenAuth:   jobTokenAuth,
	})
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func startPipeline(
	ctx context.Context,
	ctaAPIKey string,
	routeProvider app.PipelineRouteProvider,
	arrivalRepo app.ArrivalRepository,
	db *gorm.DB,
) *app.PipelineRunner {
	var vehicleClient app.CTAVehicleClient
	if app.PipelineUseFakeCTAFromEnv() {
		log.Println("pipeline: using fake CTA client")
		vehicleClient = fake.NewCTAClient()
	} else {
		if ctaAPIKey == "" {
			log.Println("pipeline: CTA_API_KEY not set — pipeline disabled")
			return nil
		}
		vehicleClient = cta.NewVehicleClient(cta.NewClient(ctaAPIKey))
	}

	var stopRepo app.StopRepository
	if db != nil {
		stopRepo = pgstore.NewStopRepo(db)
	}

	cfg, err := app.ResolvePipelineConfig(ctx, routeProvider)
	if err != nil {
		log.Printf("pipeline: config error: %v", err)
		return nil
	}
	log.Printf("pipeline: monitoring %d routes", len(cfg.RouteIDs))

	runner := app.NewPipelineRunner(vehicleClient, arrivalRepo, stopRepo, cfg)

	go func() {
		if err := runner.Run(ctx); err != nil && err != context.Canceled {
			log.Printf("pipeline error: %v", err)
		}
	}()

	return runner
}
