package controllers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

func (rc *RoleController) GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Order("id asc").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch roles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

func (rc *RoleController) GetRoleByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
		})
		return
	}

	var role models.Role
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Role not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    role,
	})
}

func (rc *RoleController) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Role name is required",
		})
		return
	}

	// Check name uniqueness
	var existing models.Role
	if err := database.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "A role with this name already exists",
		})
		return
	}

	// Find permissions if provided
	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: permissions,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create role",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Role created successfully",
		"data":    role,
	})
}

func (rc *RoleController) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
		})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Role not found",
		})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Role name is required",
		})
		return
	}

	// Check name uniqueness for other roles
	var existing models.Role
	if err := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Another role with this name already exists",
		})
		return
	}

	role.Name = req.Name
	role.Description = req.Description

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update role",
		})
		return
	}

	// Update permissions association
	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}
	database.DB.Model(&role).Association("Permissions").Replace(permissions)

	// Fetch updated role with permissions
	database.DB.Preload("Permissions").First(&role, id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role updated successfully",
		"data":    role,
	})
}

func (rc *RoleController) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
		})
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Role not found",
		})
		return
	}

	// Protect default Admin role
	if role.Name == "Admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "The system Admin role cannot be deleted",
		})
		return
	}

	// Unassign this role from any users so their accounts remain completely safe and intact
	database.DB.Model(&models.User{}).Where("role_id = ?", id).Update("role_id", nil)

	// Clear associations
	database.DB.Model(&role).Association("Permissions").Clear()

	if err := database.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role deleted successfully",
	})
}
