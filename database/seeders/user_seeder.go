package seeders

import (
	"log"

	"backend/database"
	"backend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}

		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
		}

		if err := database.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}

		log.Println("Static admin user seeded successfully (Username: admin, Password: admin123)")
	} else {
		log.Println("Admin user already exists in database")
	}
}
