package controllers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   *uint  `json:"role_id"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	RoleID   *uint  `json:"role_id"`
}

func attachUserRoles(users []models.User) []models.User {
	var roles []models.Role
	database.DB.Find(&roles)
	roleMap := make(map[uint]models.Role)
	for _, r := range roles {
		roleMap[r.ID] = r
	}

	for i := range users {
		if users[i].RoleID != nil {
			if r, ok := roleMap[*users[i].RoleID]; ok {
				roleCopy := r
				users[i].Role = &roleCopy
			}
		}
	}
	return users
}

func attachUserRole(user *models.User) {
	if user.RoleID != nil {
		var role models.Role
		if err := database.DB.Preload("Permissions").First(&role, *user.RoleID).Error; err == nil {
			user.Role = &role
		}
	}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Order("id asc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch users",
		})
		return
	}

	users = attachUserRoles(users)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	attachUserRole(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Name, username, and a password (min 6 chars) are required",
		})
		return
	}

	// Check unique username
	var existing models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username already taken",
		})
		return
	}

	// Verify role exists if role_id provided
	if req.RoleID != nil {
		var role models.Role
		if err := database.DB.First(&role, *req.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Selected role does not exist",
			})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to encrypt password",
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(hashedPassword),
		RoleID:   req.RoleID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	attachUserRole(&user)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Name and username are required",
		})
		return
	}

	// Check unique username for others
	var existing models.User
	if err := database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username already taken by another user",
		})
		return
	}

	// Verify role exists if role_id provided
	if req.RoleID != nil {
		var role models.Role
		if err := database.DB.First(&role, *req.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Selected role does not exist",
			})
			return
		}
	}

	user.Name = req.Name
	user.Username = req.Username
	user.RoleID = req.RoleID

	if req.Password != "" {
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Password must be at least 6 characters",
			})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to encrypt password",
			})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	attachUserRole(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Prevent deleting primary admin
	if user.Username == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot delete primary admin user",
		})
		return
	}

	// Prevent self deletion
	currentUserID, exists := c.Get("userID")
	if exists {
		if currID, ok := currentUserID.(uint); ok && currID == user.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You cannot delete your own user account",
			})
			return
		}
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}
