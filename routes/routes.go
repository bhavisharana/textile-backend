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

		// Public Auth Routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)

			// Protected Auth Routes
			protected := authGroup.Group("")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.GET("/me", authController.GetMe)
			}
		}

		// Protected Application Routes
		appGroup := api.Group("")
		appGroup.Use(middleware.AuthMiddleware())
		{
			// Qualities CRUD
			appGroup.GET("/qualities", qualityController.GetQualities)
			appGroup.GET("/qualities/:id", qualityController.GetQualityByID)
			appGroup.POST("/qualities", qualityController.CreateQuality)
			appGroup.PUT("/qualities/:id", qualityController.UpdateQuality)
			appGroup.DELETE("/qualities/:id", qualityController.DeleteQuality)

			// Labdips CRUD
			appGroup.GET("/labdips", labdipController.GetLabdips)
			appGroup.GET("/labdips/:id", labdipController.GetLabdipByID)
			appGroup.POST("/labdips", labdipController.CreateLabdip)
			appGroup.PUT("/labdips/:id", labdipController.UpdateLabdip)
			appGroup.DELETE("/labdips/:id", labdipController.DeleteLabdip)
		}
	}

	return r
}
