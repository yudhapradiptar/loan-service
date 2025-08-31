package main

import (
	"log"
	"os"

	"loan-service/internal/config"
	"loan-service/internal/database"
	"loan-service/internal/server"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	cfg := config.New()

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "migrate":
			if err := database.RunMigrations(cfg); err != nil {
				log.Fatal("Error running migrations:", err)
			}
			log.Println("Migrations completed successfully")
			return
		case "migrate-down":
			if err := database.RollbackMigration(cfg); err != nil {
				log.Fatal("Error rolling back migrations:", err)
			}
			log.Println("Migration rollback completed successfully")
			return
		}
	}

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
