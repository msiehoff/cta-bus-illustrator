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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var (
		routeRepo     app.RouteRepository
		ridershipRepo app.RidershipRepository
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
			log.Println("using postgres repository")
		}
	} else {
		log.Println("DATABASE_URL not set — using fake repository")
		routeRepo = &fake.RouteRepo{}
		ridershipRepo = &fake.RidershipRepo{}
	}

	ctaAPIKey := os.Getenv("CTA_API_KEY")
	ctaDataSrc := cta.NewRouteSegmentDataSource(cta.NewClient(ctaAPIKey))

	if app.PipelineEnabledFromEnv() {
		startPipeline(ctx, db, ctaAPIKey)
	}

	routeService := app.NewRouteService(routeRepo, ridershipRepo)
	a := api.New(routeService, ctaDataSrc)
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func startPipeline(ctx context.Context, db *gorm.DB, ctaAPIKey string) {
	var (
		vehicleClient app.CTAVehicleClient
		arrivalRepo   app.ArrivalRepository
		stopRepo      app.StopRepository
	)

	if app.PipelineUseFakeCTAFromEnv() {
		log.Println("pipeline: using fake CTA client")
		vehicleClient = fake.NewCTAClient()
	} else {
		if ctaAPIKey == "" {
			log.Println("pipeline: CTA_API_KEY not set — pipeline disabled")
			return
		}
		vehicleClient = cta.NewVehicleClient(cta.NewClient(ctaAPIKey))
	}

	if db != nil {
		arrivalRepo = pgstore.NewArrivalRepo(db)
		stopRepo = pgstore.NewStopRepo(db)
	} else {
		log.Println("pipeline: no database — using in-memory arrival repo")
		arrivalRepo = &fake.ArrivalRepo{}
	}

	cfg := app.PipelineConfigFromEnv()
	runner := app.NewPipelineRunner(vehicleClient, arrivalRepo, stopRepo, cfg)

	go func() {
		if err := runner.Run(ctx); err != nil && err != context.Canceled {
			log.Printf("pipeline error: %v", err)
		}
	}()
}
