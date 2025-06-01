package main

import (
	"log"
	"os"
	"time"

	"github.com/anpanovv/planter/internal/api"
	"github.com/anpanovv/planter/internal/config"
	"github.com/anpanovv/planter/internal/db"
	"github.com/anpanovv/planter/internal/jobs"
	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/repository/impl"
	"github.com/anpanovv/planter/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to the database
	database, err := db.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Apply database schema
	schema, err := os.ReadFile("scripts/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}
	if _, err := database.Exec(string(schema)); err != nil {
		log.Fatalf("Failed to apply database schema: %v", err)
	}
	log.Println("Database schema applied successfully")

	// Create repositories
	userRepo := impl.NewUserRepository(database)
	plantRepo := impl.NewPlantRepository(database)
	shopRepo := impl.NewShopRepository(database)
	recommendationRepo := impl.NewRecommendationRepository(database)
	notificationRepo := impl.NewNotificationRepository(database)

	// Create auth middleware
	auth := middleware.NewAuth(cfg.Auth.JWTSecret)

	// Create services
	authService := services.NewAuthService(userRepo, auth)
	userService := services.NewUserService(userRepo)
	plantService := services.NewPlantService(plantRepo)
	shopService := services.NewShopService(shopRepo)
	recommendationService := services.NewRecommendationService(
		recommendationRepo,
		plantRepo,
		cfg.YandexGPT.APIKey,
		cfg.YandexGPT.Model,
	)
	notificationService := services.NewNotificationService(notificationRepo, plantRepo)

	// Create and start background jobs
	log.Println("Initializing watering notifications job...")
	wateringJob := jobs.NewWateringNotificationsJob(notificationService, 1*time.Minute)
	wateringJob.Start()
	defer wateringJob.Stop()
	log.Println("Watering notifications job started successfully")

	// Create API
	api := api.New(
		authService,
		userService,
		plantService,
		shopService,
		recommendationService,
		notificationService,
		auth,
	)

	// Start the API server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	err = api.Start(cfg)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}