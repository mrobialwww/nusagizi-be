package main

import (
	"log"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/database"
	"nusagizi_be/internal/middleware"

	// "nusagizi_be/internal/handlers"

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

    // 2. Koneksi DB (pakai Connect yang sudah kamu buat)
    var pool *pgxpool.Pool
    pool, err = database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer pool.Close()

    // 3. Buat JWT validator
    jwtValidator, err := auth.NewValidator(cfg.Auth0Domain, cfg.Auth0Audience)
    if err != nil {
        log.Fatal("Failed to create JWT validator:", err)
    }

    // 4. Buat JWT middleware
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
						"message":  "Todo API is running well!",
						"status":   "success",
						"database": "connected",
					})
				})

    // 7. Protected routes
    protected := router.Group("/")
    protected.Use(middleware.GinMiddleware(jwtMiddleware, pool))
    {
        protected.GET("/profile", func(c *gin.Context) {
            user, _ := c.Get("user")
            c.JSON(200, gin.H{
                "message": "Hello from private endpoint!",
                "user":    user,
            })
        })
    }

	router.Run(":" + cfg.Port)
}