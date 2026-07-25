package main

import (
	"log"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/database"
	"nusagizi_be/internal/handlers"
	"nusagizi_be/internal/middleware"

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

	// 4. JWT middleware
	jwtMiddleware, err := middleware.NewMiddleware(jwtValidator)
	if err != nil {
		log.Fatal("Failed to create JWT middleware:", err)
	}

	// 5. Setup Gin
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// 6. Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "NusaGizi API is running!",
			"status":   "success",
			"database": "connected",
		})
	})

	// 7. Protected routes (require valid Auth0 JWT)
	protected := router.Group("/")
	protected.Use(middleware.GinMiddleware(jwtMiddleware, pool))
	{
		// Legacy — onboarding sets role in Auth0 + DB
		protected.POST("/onboarding", handlers.OnboardingHandler(pool, cfg))

		// Modul 1: Auth & User Profile
		protected.GET("/users/:user_id", handlers.GetUserHandler(pool))
		protected.PATCH("/users/:user_id", handlers.UpdateUserHandler(pool))
		protected.DELETE("/users/:user_id", handlers.DeleteUserHandler(pool))

		// Modul 2: Child Profile
		protected.POST("/children", handlers.CreateChildHandler(pool))
		protected.GET("/children/:child_id", handlers.GetChildHandler(pool))
		protected.PATCH("/children/:child_id", handlers.UpdateChildHandler(pool))
		protected.DELETE("/children/:child_id", handlers.DeleteChildHandler(pool))
		protected.GET("/mother-profiles/:mother_profile_id/children", handlers.GetChildrenHandler(pool))
	}

	router.Run(":" + cfg.Port)
}
