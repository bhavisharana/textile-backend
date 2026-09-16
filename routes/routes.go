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
			// Qualities Management (Protected)
			qualities := protected.Group("/qualities")
			{
				qualities.GET("", qualityController.GetQualities)
				qualities.GET("/:id", qualityController.GetQualityByID)
				qualities.POST("", qualityController.CreateQuality)
				qualities.PUT("/:id", qualityController.UpdateQuality)
				qualities.DELETE("/:id", qualityController.DeleteQuality)
			}

			// Labdips Management (Protected)
			labdips := protected.Group("/labdips")
			{
				labdips.GET("", labdipController.GetLabdips)
				labdips.GET("/:id", labdipController.GetLabdipByID)
				labdips.POST("", labdipController.CreateLabdip)
				labdips.PUT("/:id", labdipController.UpdateLabdip)
				labdips.DELETE("/:id", labdipController.DeleteLabdip)
			}

			// Parties Management (Protected)
			parties := protected.Group("/parties")
			{
				parties.GET("", partyController.GetParties)
				parties.GET("/:id", partyController.GetPartyByID)
				parties.POST("", partyController.CreateParty)
				parties.PUT("/:id", partyController.UpdateParty)
				parties.DELETE("/:id", partyController.DeleteParty)
			}

			// Orders Management (Protected)
			orders := protected.Group("/orders")
			{
				orders.GET("", orderController.GetOrders)
				orders.GET("/:id", orderController.GetOrderByID)
				orders.POST("", orderController.CreateOrder)
				orders.PUT("/:id", orderController.UpdateOrder)
				orders.DELETE("/:id", orderController.DeleteOrder)
			}
		}
	}

	return r
}
