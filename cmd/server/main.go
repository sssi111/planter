package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anpanovv/planter/internal/api"
	"github.com/anpanovv/planter/internal/db"
	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/repository/impl"
	"github.com/anpanovv/planter/internal/jobs"
	"github.com/anpanovv/planter/internal/services"
)

func main() {
	// Initialize database
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create repositories
	userRepo := impl.NewUserRepository(database)
	plantRepo := impl.NewPlantRepository(database)
	shopRepo := impl.NewShopRepository(database)
	notificationRepo := impl.NewNotificationRepository(database)

	// Create services
	userService := services.NewUserService(userRepo)
	plantService := services.NewPlantService(plantRepo)
	shopService := services.NewShopService(shopRepo)
	notificationService := services.NewNotificationService(notificationRepo, plantRepo)

	// Create and start background jobs
	wateringJob := jobs.NewWateringNotificationsJob(notificationService, 1*time.Hour)
	wateringJob.Start()
	defer wateringJob.Stop()

	// Create auth middleware first
	authMiddleware := middleware.NewAuth("development-secret-key") // TODO: Replace with config value
	
	// Create additional services
	authService := services.NewAuthService(userRepo, authMiddleware)
	recommendationService := services.NewRecommendationService(
		impl.NewRecommendationRepository(database),
		plantRepo,
		"", // yandexGPT API key
		"", // yandexGPT model
	)

	// Create and start API server
	apiHandler := api.New(
		authService,
		userService,
		plantService,
		shopService,
		recommendationService,
		notificationService,
		authMiddleware,
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: apiHandler.Handler(),
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
} 