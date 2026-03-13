package main

import (
	"log"

	"github.com/msiehoff/cta-bus-illustrator/backend/api"
	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/storage/fake"
)

func main() {
	routeRepo := &fake.RouteRepo{}
	routeService := app.NewRouteService(routeRepo)

	a := api.New(routeService)
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
