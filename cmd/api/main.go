package main

import (
	"context"
	"log"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/auth/auth0client"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/database"
	"nusagizi_be/internal/handlers"
	"nusagizi_be/internal/infrastructure/storage/r2"
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

	// Auth0 Clients for OTP
	idTokenVerifier, err := auth0client.NewIDTokenVerifier(context.Background(), cfg.Auth0Domain, cfg.Auth0ClientID)
	if err != nil {
		log.Fatal("Failed to create IDTokenVerifier:", err)
	}
	managementClient := auth0client.NewManagementClient(cfg.Auth0Domain, cfg.M2MClientID, cfg.M2MClientSecret)

	// Cloudflare R2 Client Init
	r2Client, err := r2.New(context.Background(),
		cfg.R2AccountID, cfg.R2AccessKeyID, cfg.R2SecretAccessKey, cfg.R2Bucket,
	)
	if err != nil {
		log.Fatal("Failed to connect to R2:", err)
	}

	// 5. DI Wiring: Repositories -> Services -> Handlers

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	motherRepo := repository.NewMotherProfileRepository(pool)
	childRepo := repository.NewChildRepository(pool)
	caregiverRepo := repository.NewCaregiverRepository(pool)
	childDevRepo := repository.NewChildDevelopmentRepository(pool)
	menuRepo := repository.NewMenuRepository(pool)
	childNutriRepo := repository.NewChildNutritionRepository(pool)
	childGrowthRepo := repository.NewChildGrowthRepository(pool)
	dashboardRepo := repository.NewDashboardRepository(pool)
	medicalRepo := repository.NewMedicalRepository(pool)
	photosContactsRepo := repository.NewPhotosContactsRepository(pool)
	notificationRepo := repository.NewNotificationRepository(pool)
	checkinRepo := repository.NewCheckinRepository(pool)

	// Services
	emailOTPSvc := services.NewEmailOTPService(idTokenVerifier, managementClient)
	imageSvc := services.NewImageService(r2Client, motherRepo, childRepo, caregiverRepo)
	userSvc := services.NewUserService(userRepo, motherRepo, caregiverRepo, cfg.R2PublicURL)
	childSvc := services.NewChildService(childRepo, motherRepo, caregiverRepo, cfg.R2PublicURL)
	childDevSvc := services.NewChildDevelopmentService(childDevRepo, childRepo, motherRepo, cfg.R2PublicURL)
	menuSvc := services.NewMenuService(cfg, menuRepo, motherRepo, childRepo, childNutriRepo)
	childNutriSvc := services.NewChildNutritionService(childNutriRepo, childRepo, motherRepo, caregiverRepo, cfg.R2PublicURL)
	childGrowthSvc := services.NewChildGrowthService(childGrowthRepo, childRepo, motherRepo)
	caregiverSvc := services.NewCaregiverService(caregiverRepo, motherRepo)
	dashboardSvc := services.NewDashboardService(dashboardRepo, motherRepo)
	medicalSvc := services.NewMedicalService(medicalRepo, motherRepo, childRepo)
	photosContactsSvc := services.NewPhotosContactsService(photosContactsRepo, motherRepo, childRepo, caregiverRepo, cfg.R2PublicURL)
	notificationSvc := services.NewNotificationService(notificationRepo)
	checkinSvc := services.NewCheckinService(checkinRepo, caregiverRepo)

	// Handlers
	emailOTPHandler := handlers.NewConfirmEmailHandler(emailOTPSvc)
	onboardHandler := handlers.NewOnboardingHandler(userRepo, cfg)
	imageHandler := handlers.NewImageHandler(imageSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	childHandler := handlers.NewChildHandler(childSvc)
	childDevHandler := handlers.NewChildDevelopmentHandler(childDevSvc)
	menuHandler := handlers.NewMenuHandler(menuSvc)
	childNutriHandler := handlers.NewChildNutritionHandler(childNutriSvc)
	childGrowthHandler := handlers.NewChildGrowthHandler(childGrowthSvc)
	caregiverHandler := handlers.NewCaregiverHandler(caregiverSvc)
	dashboardHandler := handlers.NewDashboardHandler(dashboardSvc)
	medicalHandler := handlers.NewMedicalHandler(medicalSvc)
	photosContactsHandler := handlers.NewPhotosContactsHandler(photosContactsSvc)
	notificationHandler := handlers.NewNotificationHandler(notificationSvc)
	checkinHandler := handlers.NewCheckinHandler(checkinSvc)

	// 6. Setup Gin
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// 7. Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "NusaGizi API is running!",
			"status":   "success",
			"database": "connected",
		})
	})
	router.POST("/confirm-email", emailOTPHandler.Handle)

	// 8. Protected routes (require valid Auth0 JWT)
	protected := router.Group("/")
	protected.Use(middleware.GinMiddleware(jwtMiddleware, userRepo))
	{
		// Legacy — onboarding sets role in Auth0 + DB
		protected.POST("/onboarding", onboardHandler.Handle)

		// Modul 1: Auth & User Profile
		protected.GET("/users", userHandler.GetUserProfile)
		protected.PATCH("/users", userHandler.UpdateUserProfile)
		protected.DELETE("/users", userHandler.DeleteUserProfile)

		protected.POST("/mother-profiles", userHandler.CreateMotherProfile)
		protected.GET("/mother-profiles", userHandler.GetMotherProfile)

		protected.POST("/caregiver-profiles", userHandler.CreateCaregiverProfile)
		protected.GET("/caregiver-profiles", userHandler.GetCaregiverProfile)

		// Modul 1.5: Upload Photo (R2)
		protected.POST("/images/presign-upload", imageHandler.PresignUpload)
		protected.POST("/images/confirm", imageHandler.Confirm)

		// Modul 11: Checkin
		protected.POST("/checkin/generate", checkinHandler.Generate)
		protected.POST("/checkin/validate", checkinHandler.Validate)

		// Modul 2: Child Profile
		protected.POST("/children", childHandler.CreateChildProfile)
		protected.GET("/children", childHandler.GetListChild)
		protected.GET("/children/:child_id", childHandler.GetChild)
		protected.PATCH("/children/:child_id", childHandler.UpdateChildProfile)
		protected.DELETE("/children/:child_id", childHandler.DeleteChildProfile)
		protected.GET("/children/:child_id/profile", childHandler.GetChildProfile)

		// Modul 3: Child Growth
		protected.GET("/children/:child_id/growth-reports/latest", childGrowthHandler.GetLatestGrowthReport)
		protected.GET("/children/:child_id/growth-analyses", childGrowthHandler.GetGrowthAnalyses)
		protected.GET("/children/:child_id/growth-reports", childGrowthHandler.GetGrowthReports)
		protected.POST("/children/:child_id/growth-reports", childGrowthHandler.CreateGrowthReport)
		protected.PATCH("/growth-reports/:child_growth_report_id", childGrowthHandler.UpdateGrowthReport)
		protected.DELETE("/growth-reports/:child_growth_report_id", childGrowthHandler.DeleteGrowthReport)

		// Modul 4: Child Development (KPSP)
		protected.GET("/children/:child_id/development-reports/latest", childDevHandler.GetLatestDevelopmentReport)
		protected.GET("/children/:child_id/development-reports", childDevHandler.GetDevelopmentReports)
		protected.GET("/development-reports/:child_development_report_id", childDevHandler.GetDevelopmentReportByID)
		protected.GET("/assessment-kpsp-questions", childDevHandler.GetKPSPQuestions)
		protected.POST("/children/:child_id/development-reports", childDevHandler.CreateDevelopmentReport)
		protected.PATCH("/children/:child_id/development-reports/:child_development_report_id", childDevHandler.UpdateDevelopmentReport)
		protected.GET("/checklist-milestone-tasks", childDevHandler.GetChecklistMilestoneTasks)
		protected.GET("/children/:child_id/development-reports/:report_id/recommendations", childDevHandler.GetRecommendations)
		protected.PATCH("/children/:child_id/checklist-milestone-progress", childDevHandler.UpdateChecklistMilestone)
		protected.DELETE("/development-reports/:child_development_report_id", childDevHandler.DeleteDevelopmentReport)

		// Modul 5: Child Nutrition & Menu
		protected.GET("/children/:child_id/nutrition/today", childNutriHandler.GetTodayNutritionReport)
		protected.POST("/children/:child_id/nutrition/reuse-recipe", childNutriHandler.ReuseRecipe)
		protected.GET("/children/:child_id/nutrition-reports/:report_id", childNutriHandler.GetReportMenuByID)
		protected.PATCH("/recipes/:recipe_id/complete", childNutriHandler.UpdateRecipeCompleteStatus)
		protected.PATCH("/recipes/:recipe_id/bookmark", childNutriHandler.UpdateRecipeBookmarkStatus)
		protected.GET("/recipes/:recipe_id", childNutriHandler.GetRecipeDetail)
		protected.PATCH("/recipes/main-ingredients/priority", childNutriHandler.SwapMainIngredientPriority)
		protected.GET("/children/:child_id/recipes/bookmarked", childNutriHandler.GetBookmarkedRecipes)
		protected.GET("/children/:child_id/nutrition-reports", childNutriHandler.GetNutritionReportsByMonth)
		protected.GET("/menu/daily-shop", childNutriHandler.GetDailyShopIngredients)
		protected.GET("/children/:child_id/nutrition-reports/today/shopping", childNutriHandler.GetTodayMenuShopping)
		protected.POST("/menu/generate", menuHandler.GenerateMenu)

		// Modul 6: Photos & Contacts
		protected.GET("/mother-profiles/contacts", photosContactsHandler.GetContacts)
		protected.POST("/contacts", photosContactsHandler.AddContact)
		protected.DELETE("/contacts/:contact_id", photosContactsHandler.DeleteContact)
		protected.GET("/mother-profiles/child-photos", photosContactsHandler.GetMotherChildPhotos)
		protected.GET("/contacts/:contact_id/child-photos", photosContactsHandler.GetContactChildPhotos)
		protected.GET("/mother-profiles/child-photos/all", photosContactsHandler.GetAllChildPhotos)
		protected.GET("/child-photos/:child_photo_id", photosContactsHandler.GetPhotoDetail)
		protected.POST("/children/:child_id/photos/mother", photosContactsHandler.AddPhotoMother)
		protected.POST("/children/:child_id/photos/caregiver", photosContactsHandler.AddPhotoCaregiver)
		protected.PATCH("/children/:child_id/photos/:child_photo_id", photosContactsHandler.UpdatePhoto)
		protected.DELETE("/child-photos/:child_photo_id", photosContactsHandler.DeletePhoto)

		// Modul 7: Caregiver
		protected.GET("/mother-profiles/caregiver-engagements", caregiverHandler.GetCaregiverEngagements)
		protected.GET("/mother-profiles/caregiver-engagements/revoked", caregiverHandler.GetCaregiverEngagementsRevoked)
		protected.DELETE("/caregiver-engagements/:caregiver_engagement_id", caregiverHandler.DeleteCaregiverEngagement)

		// Modul 8: Medical
		protected.GET("/medical-notes", medicalHandler.GetMedicalNotes)
		protected.GET("/medical-notes/:medical_note_id", medicalHandler.GetMedicalNoteDetail)
		protected.POST("/medical-notes", medicalHandler.CreateMedicalNote)
		protected.PATCH("/medical-notes/:medical_note_id", medicalHandler.UpdateMedicalNote)
		protected.DELETE("/medical-notes/:medical_note_id", medicalHandler.DeleteMedicalNote)

		// Modul 9: Dashboard
		protected.GET("/dashboard/children", dashboardHandler.GetDashboardSummary)

		// Modul 10: Notifications
		protected.GET("/notifications", notificationHandler.GetNotifications)
		protected.GET("/notifications/latest", notificationHandler.GetLatestNotification)
	}

	router.Run(":" + cfg.Port)
}
