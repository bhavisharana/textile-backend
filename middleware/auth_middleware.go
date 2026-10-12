package middleware

import (
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid Authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// Check permission dynamically based on route and HTTP method
		reqPath := c.Request.URL.Path
		reqMethod := c.Request.Method
		requiredPerm := getRequiredPermission(reqPath, reqMethod)

		if requiredPerm != "" {
			// Super admin username bypasses all checks
			if claims.Username == "admin" {
				c.Next()
				return
			}

			var user models.User
			if err := database.DB.First(&user, claims.UserID).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "User not found",
				})
				c.Abort()
				return
			}

			if user.RoleID == nil {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Access denied: No role assigned to your account",
				})
				c.Abort()
				return
			}

			var role models.Role
			if err := database.DB.Preload("Permissions").First(&role, *user.RoleID).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Access denied: Assigned role not found",
				})
				c.Abort()
				return
			}

			// Admin role bypasses all checks
			if role.Name == "Admin" {
				c.Next()
				return
			}

			// Check if required permission exists in role's permissions
			hasPermission := false
			for _, perm := range role.Permissions {
				if perm.Name == requiredPerm {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Access denied: Missing required permission (" + requiredPerm + ")",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func getRequiredPermission(path string, method string) string {
	path = strings.TrimPrefix(path, "/api")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	module := parts[0]
	switch module {
	case "users", "roles", "qualities", "labdips", "parties", "orders":
		// valid resource modules
	case "permissions":
		return "roles.view"
	case "dashboard", "auth":
		return ""
	default:
		return ""
	}

	action := "view"
	if method == http.MethodGet {
		action = "view"
	} else if strings.HasSuffix(path, "/delete") || method == http.MethodDelete {
		action = "delete"
	} else if strings.HasSuffix(path, "/update") || strings.HasSuffix(path, "/status") || method == http.MethodPut || method == http.MethodPatch {
		action = "edit"
	} else if method == http.MethodPost {
		action = "create"
	}

	return module + "." + action
}
