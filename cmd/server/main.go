package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/radhikadarode/finance-backend/internal/database"
	"github.com/radhikadarode/finance-backend/internal/handlers"
	"github.com/radhikadarode/finance-backend/internal/middleware"
	"github.com/radhikadarode/finance-backend/internal/models"
)

func main() {
	// Initialize database
	database.Init()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, models.APIResponse{Success: true, Message: "Finance Backend is running"})
	})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	recordHandler := handlers.NewRecordHandler()
	dashboardHandler := handlers.NewDashboardHandler()

	// Public routes
	r.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	api := r.Group("/api", middleware.AuthMiddleware())
	{
		// Current user profile (all roles)
		api.GET("/me", userHandler.GetMe)

		// User management (Admin only)
		users := api.Group("/users", middleware.RequireRole(models.RoleAdmin))
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Financial records
		records := api.Group("/records")
		{
			// View records: all roles
			records.GET("", recordHandler.GetRecords)
			records.GET("/:id", recordHandler.GetRecord)

			// Create records: Admin and Analyst
			records.POST("",
				middleware.RequireRole(models.RoleAdmin, models.RoleAnalyst),
				recordHandler.CreateRecord,
			)

			// Update and Delete: Admin only
			records.PUT("/:id",
				middleware.RequireRole(models.RoleAdmin),
				recordHandler.UpdateRecord,
			)
			records.DELETE("/:id",
				middleware.RequireRole(models.RoleAdmin),
				recordHandler.DeleteRecord,
			)
		}

		// Dashboard summary: Admin and Analyst
		api.GET("/dashboard/summary",
			middleware.RequireRole(models.RoleAdmin, models.RoleAnalyst),
			dashboardHandler.GetSummary,
		)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
