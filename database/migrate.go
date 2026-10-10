package database

import (
	"log"

	"backend/models"
)

// RunMigrations executes the auto-migration for all models
func RunMigrations() {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDB first.")
	}

	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Quality{},
		&models.Labdip{},
		&models.Party{},
		&models.Order{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
