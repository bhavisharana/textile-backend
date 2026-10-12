package database

import (
	"log"
	"strings"

	"backend/models"
)

// RunMigrations executes the auto-migration for all models
func RunMigrations() {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDB first.")
	}

	log.Println("Running all database migrations...")

	err := DB.AutoMigrate(
		&models.Permission{},
		&models.Role{},
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

// RunSpecificMigration runs migration for a specific table/model
func RunSpecificMigration(name string) {
	if DB == nil {
		log.Fatal("Database is not initialized. Call InitDB first.")
	}

	normalized := strings.ToLower(strings.TrimSpace(name))
	log.Printf("Running migration for: %s...\n", normalized)

	var err error
	switch normalized {
	case "users", "user":
		err = DB.AutoMigrate(&models.User{})
	case "roles", "role":
		err = DB.AutoMigrate(&models.Role{})
	case "permissions", "permission":
		err = DB.AutoMigrate(&models.Permission{})
	case "role_permissions", "roles_permissions":
		err = DB.AutoMigrate(&models.Permission{}, &models.Role{})
	case "qualities", "quality":
		err = DB.AutoMigrate(&models.Quality{})
	case "labdips", "labdip":
		err = DB.AutoMigrate(&models.Labdip{})
	case "parties", "party":
		err = DB.AutoMigrate(&models.Party{})
	case "orders", "order":
		err = DB.AutoMigrate(&models.Order{})
	case "all":
		RunMigrations()
		return
	default:
		log.Fatalf("Unknown migration '%s'. Available options: users, roles, permissions, role_permissions, qualities, labdips, parties, orders, all", name)
	}

	if err != nil {
		log.Fatalf("Failed to run migration '%s': %v", name, err)
	}

	log.Printf("Migration '%s' completed successfully!\n", name)
}
