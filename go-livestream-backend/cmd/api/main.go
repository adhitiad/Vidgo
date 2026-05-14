package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-livestream-backend/internal/config"
	"go-livestream-backend/internal/controller"
	"go-livestream-backend/internal/domain"
	"go-livestream-backend/internal/middleware"
	"go-livestream-backend/internal/repository"
	"go-livestream-backend/internal/usecase"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db := config.ConnectDatabase(cfg)

	// Init repositories
	userRepo := repository.NewUserRepository(db)

	// Init usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)

	// Init controllers
	authController := controller.NewAuthController(authUsecase)

	// Setup router
	r := gin.Default()

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to your profile",
				"user_id": userID,
				"role":    role,
			})
		})

		// Admin only route
		adminRoutes := protected.Group("/admin")
		adminRoutes.Use(middleware.RoleMiddleware(domain.RoleAdmin))
		{
			adminRoutes.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Admin Dashboard"})
			})
		}

		// Agency or Admin route
		agencyRoutes := protected.Group("/agency")
		agencyRoutes.Use(middleware.RoleMiddleware(domain.RoleAdmin, domain.RoleAgency))
		{
			agencyRoutes.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Agency Dashboard"})
			})
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
