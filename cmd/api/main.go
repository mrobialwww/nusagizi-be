package main

import (
	"log"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/database"
	"nusagizi_be/internal/handlers"
	"nusagizi_be/internal/middleware"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Load config
	var cfg *config.Config
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// 2. Connect to DB
	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// 3. JWT validator
	jwtValidator, err := auth.NewValidator(cfg.Auth0Domain, cfg.Auth0Audience)
	if err != nil {
		log.Fatal("Failed to create JWT validator:", err)
	}

	// 4. JWT middleware initialization
	jwtMiddleware, err := middleware.NewMiddleware(jwtValidator)
	if err != nil {
		log.Fatal("Failed to create JWT middleware:", err)
	}

	// 5. DI Wiring: Repositories -> Services -> Handlers

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	motherRepo := repository.NewMotherProfileRepository(pool)
	childRepo := repository.NewChildRepository(pool)
	caregiverRepo := repository.NewCaregiverRepository(pool)

	// Services
	userSvc := services.NewUserService(userRepo, motherRepo, caregiverRepo)
	childSvc := services.NewChildService(childRepo, motherRepo, caregiverRepo)

	// Handlers
	onboardHandler := handlers.NewOnboardingHandler(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userSvc)
	childHandler := handlers.NewChildHandler(childSvc)

	// 6. Router Setup

	// 7. Setup Gin
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// 8. Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "NusaGizi API is running!",
			"status":   "success",
			"database": "connected",
		})
	})

	// 9. Protected routes (require valid Auth0 JWT)
	protected := router.Group("/")
	// Note: GinMiddleware now requires userRepo instead of pool
	protected.Use(middleware.GinMiddleware(jwtMiddleware, userRepo))
	{
		// Legacy — onboarding sets role in Auth0 + DB
		protected.POST("/onboarding", onboardHandler.Handle)

		// Modul 1: Auth & User Profile
		protected.GET("/users", userHandler.Get)
		protected.PATCH("/users", userHandler.Update)
		protected.DELETE("/users", userHandler.Delete)
		protected.GET("/mother-profiles", userHandler.GetMotherProfile)
		protected.GET("/caregiver-profiles", userHandler.GetCaregiverProfile)

		// Modul 2: Child Profile
		protected.POST("/children", childHandler.Create)
		protected.GET("/children/:child_id", childHandler.Get)
		protected.PATCH("/children/:child_id", childHandler.Update)
		protected.DELETE("/children/:child_id", childHandler.Delete)
		protected.GET("/children", childHandler.ListByMother)
	}

	router.Run(":" + cfg.Port)
}
