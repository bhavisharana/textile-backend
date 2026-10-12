package seeders

import (
	"log"

	"backend/database"
	"backend/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	// First seed permissions
	allPerms := SeedPermissions()

	// Seed or find Admin Role
	var adminRole models.Role
	err := database.DB.Preload("Permissions").Where("name = ?", "Admin").First(&adminRole).Error
	if err != nil {
		adminRole = models.Role{
			Name:        "Admin",
			Description: "Full administrative access with all permissions",
			Permissions: allPerms,
		}
		if err := database.DB.Create(&adminRole).Error; err != nil {
			log.Fatalf("Failed to seed admin role: %v", err)
		}
		log.Println("Admin role created with full permissions")
	} else {
		// Update permissions association for admin role in case new permissions were added
		database.DB.Model(&adminRole).Association("Permissions").Replace(allPerms)
	}

	var user models.User
	err = database.DB.Where("username = ?", "admin").First(&user).Error

	if err != nil {
		// Create admin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}

		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Name:     "System Administrator",
			RoleID:   &adminRole.ID,
		}

		if err := database.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}

		log.Println("Static admin user seeded successfully (Username: admin, Password: admin123, Role: Admin)")
	} else {
		// If user exists, ensure Name and RoleID are assigned
		updates := map[string]interface{}{}
		if user.RoleID == nil {
			updates["role_id"] = adminRole.ID
		}
		if user.Name == "" {
			updates["name"] = "System Administrator"
		}
		if len(updates) > 0 {
			database.DB.Model(&user).Updates(updates)
			log.Println("Admin user updated with Role and Name")
		}
		log.Println("Admin user ready in database")
	}
}
