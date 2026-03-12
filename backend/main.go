package main

import (
	"log"

	"github.com/msiehoff/cta-bus-illustrator/backend/api"
)

func main() {
	a := api.New()
	if err := a.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
