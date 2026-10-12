package seeders

import (
	"log"

	"backend/database"
	"backend/models"
)

type PermDef struct {
	Name   string
	Module string
	Action string
}

func SeedPermissions() []models.Permission {
	permissions := []PermDef{
		// Parties
		{Name: "parties.view", Module: "parties", Action: "view"},
		{Name: "parties.create", Module: "parties", Action: "create"},
		{Name: "parties.edit", Module: "parties", Action: "edit"},
		{Name: "parties.delete", Module: "parties", Action: "delete"},

		// Qualities
		{Name: "qualities.view", Module: "qualities", Action: "view"},
		{Name: "qualities.create", Module: "qualities", Action: "create"},
		{Name: "qualities.edit", Module: "qualities", Action: "edit"},
		{Name: "qualities.delete", Module: "qualities", Action: "delete"},

		// Labdips
		{Name: "labdips.view", Module: "labdips", Action: "view"},
		{Name: "labdips.create", Module: "labdips", Action: "create"},
		{Name: "labdips.edit", Module: "labdips", Action: "edit"},
		{Name: "labdips.delete", Module: "labdips", Action: "delete"},

		// Orders
		{Name: "orders.view", Module: "orders", Action: "view"},
		{Name: "orders.create", Module: "orders", Action: "create"},
		{Name: "orders.edit", Module: "orders", Action: "edit"},
		{Name: "orders.delete", Module: "orders", Action: "delete"},

		// Users
		{Name: "users.view", Module: "users", Action: "view"},
		{Name: "users.create", Module: "users", Action: "create"},
		{Name: "users.edit", Module: "users", Action: "edit"},
		{Name: "users.delete", Module: "users", Action: "delete"},

		// Roles
		{Name: "roles.view", Module: "roles", Action: "view"},
		{Name: "roles.create", Module: "roles", Action: "create"},
		{Name: "roles.edit", Module: "roles", Action: "edit"},
		{Name: "roles.delete", Module: "roles", Action: "delete"},
	}

	for _, p := range permissions {
		var existing models.Permission
		if err := database.DB.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			database.DB.Create(&models.Permission{
				Name:   p.Name,
				Module: p.Module,
				Action: p.Action,
			})
		}
	}

	log.Println("Permissions seeded successfully")

	var allPerms []models.Permission
	database.DB.Find(&allPerms)
	return allPerms
}
