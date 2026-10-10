package main

import (
	"log"
	"os"

	"backend/database"
	"backend/database/seeders"
	"backend/routes"
)

func main() {
	// Initialize Database
	database.InitDB()

	// Check for CLI arguments (e.g. "go run main.go migrate")
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		database.RunMigrations()
		return
	}

	// Run Seeders (Optional: you might want to run seeders via a command too)
	seeders.SeedUser()

	// Setup Router & Start Server
	router := routes.SetupRouter()

	log.Println("Starting Textile Admin API server on :8081...")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
