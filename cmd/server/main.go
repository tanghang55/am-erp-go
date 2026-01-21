package main

import (
	"fmt"
	"log"

	"am-erp-go/internal/infrastructure/bootstrap"
)

func main() {
	app, err := bootstrap.Build()
	if err != nil {
		log.Fatalf("Failed to bootstrap app: %v", err)
	}

	addr := fmt.Sprintf(":%s", app.Config.Server.Port)
	log.Printf("Starting server on %s", addr)

	if err := app.Engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
