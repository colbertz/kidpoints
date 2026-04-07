package main

import (
	"log"

	"points-app/config"
	"points-app/database"
	"points-app/router"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg.DBFilePath); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	r := router.Setup()
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
