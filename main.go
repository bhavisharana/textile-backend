package main

import (
	"log"

	"backend/database"
	"backend/database/seeders"
	"backend/routes"
)

func main() {
	// Initialize Database
	database.InitDB()

	// Run Seeders
	seeders.SeedUser()

	// Setup Router & Start Server
	router := routes.SetupRouter()

	log.Println("Starting Textile Admin API server on :8081...")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}