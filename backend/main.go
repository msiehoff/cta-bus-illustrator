package main

import (
	"log"
	"os"

	"github.com/msiehoff/cta-bus-illustrator/backend/api"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/migrations"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
	pgstore "github.com/msiehoff/cta-bus-illustrator/backend/storage/postgres"
)

func main() {
	var routeRepo app.RouteRepository

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		db, err := pgstore.Connect(dsn)
		if err != nil {
			log.Fatalf("db connection error: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get underlying sql.DB: %v", err)
		}

		if err := migrations.Run(sqlDB); err != nil {
			log.Fatalf("migration error: %v", err)
		}

		routeRepo = pgstore.NewRouteRepo(db)
		log.Println("using postgres repository")
	} else {
		log.Println("DATABASE_URL not set — using fake repository")
		routeRepo = &fake.RouteRepo{}
	}

	routeService := app.NewRouteService(routeRepo)
	a := api.New(routeService)
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
