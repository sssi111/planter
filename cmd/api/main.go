package main

import (
	"log"

	"github.com/anpanovv/planter/internal/api"
	"github.com/anpanovv/planter/internal/config"
	"github.com/anpanovv/planter/internal/db"
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

	// Create repositories
	userRepo := impl.NewUserRepository(database)
	plantRepo := impl.NewPlantRepository(database)
	shopRepo := impl.NewShopRepository(database)
	recommendationRepo := impl.NewRecommendationRepository(database)

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

	// Create API
	api := api.New(
		authService,
		userService,
		plantService,
		shopService,
		recommendationService,
		auth,
	)

	// Start the API server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	err = api.Start(cfg)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}