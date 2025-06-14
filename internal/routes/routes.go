package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/money-advice/receipt-backend/internal/handlers"
	"github.com/money-advice/receipt-backend/internal/middleware"
	"github.com/money-advice/receipt-backend/internal/services"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// Initialize services
	authService := services.NewAuthService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	api := router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/google", authHandler.GoogleAuth)
			auth.GET("/validate", authHandler.ValidateToken)
		}
	}

	// Protected routes (require authentication)
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Add protected routes here
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetString("user_id")
			c.JSON(200, gin.H{
				"message": "This is a protected route",
				"user_id": userID,
			})
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Receipt Backend API is running",
		})
	})
}