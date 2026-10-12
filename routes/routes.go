package routes

import (
	"net/http"

	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Textile API is running",
			})
		})

		authController := controllers.NewAuthController()
		qualityController := controllers.NewQualityController()
		labdipController := controllers.NewLabdipController()
		partyController := controllers.NewPartyController()
		orderController := controllers.NewOrderController()
		dashboardController := controllers.NewDashboardController()
		userController := controllers.NewUserController()
		roleController := controllers.NewRoleController()
		permissionController := controllers.NewPermissionController()

		// Authentication & Verification Routes (Public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)

			// Protected Auth Routes
			protectedAuth := authGroup.Group("")
			protectedAuth.Use(middleware.AuthMiddleware())
			{
				protectedAuth.GET("/me", authController.GetMe)
			}
		}

		// Protected Application Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Dashboard Statistics (Protected)
			protected.GET("/dashboard/stats", dashboardController.GetDashboardStats)

			// Users Management (Protected)
			users := protected.Group("/users")
			{
				users.GET("", userController.GetUsers)
				users.GET("/:id", userController.GetUserByID)
				users.POST("", userController.CreateUser)
				users.POST("/:id/update", userController.UpdateUser)
				users.POST("/:id/delete", userController.DeleteUser)
			}

			// Roles Management (Protected)
			roles := protected.Group("/roles")
			{
				roles.GET("", roleController.GetRoles)
				roles.GET("/:id", roleController.GetRoleByID)
				roles.POST("", roleController.CreateRole)
				roles.POST("/:id/update", roleController.UpdateRole)
				roles.POST("/:id/delete", roleController.DeleteRole)
			}

			// Permissions Management (Protected)
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", permissionController.GetPermissions)
			}

			// Qualities Management (Protected)
			qualities := protected.Group("/qualities")
			{
				qualities.GET("", qualityController.GetQualities)
				qualities.GET("/:id", qualityController.GetQualityByID)
				qualities.POST("", qualityController.CreateQuality)
				qualities.POST("/:id/update", qualityController.UpdateQuality)
				qualities.POST("/:id/delete", qualityController.DeleteQuality)
			}

			// Labdips Management (Protected)
			labdips := protected.Group("/labdips")
			{
				labdips.GET("", labdipController.GetLabdips)
				labdips.GET("/:id", labdipController.GetLabdipByID)
				labdips.POST("", labdipController.CreateLabdip)
				labdips.POST("/:id/update", labdipController.UpdateLabdip)
				labdips.POST("/:id/delete", labdipController.DeleteLabdip)
			}

			// Parties Management (Protected)
			parties := protected.Group("/parties")
			{
				parties.GET("", partyController.GetParties)
				parties.GET("/:id", partyController.GetPartyByID)
				parties.POST("", partyController.CreateParty)
				parties.POST("/:id/update", partyController.UpdateParty)
				parties.POST("/:id/delete", partyController.DeleteParty)
			}

			// Orders Management (Protected)
			orders := protected.Group("/orders")
			{
				orders.GET("", orderController.GetOrders)
				orders.GET("/:id", orderController.GetOrderByID)
				orders.POST("", orderController.CreateOrder)
				orders.POST("/:id/update", orderController.UpdateOrder)
				orders.POST("/:id/status", orderController.UpdateOrderStatus)
				orders.POST("/:id/delete", orderController.DeleteOrder)
			}
		}
	}

	return r
}
