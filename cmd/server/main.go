package main

import (
	"database/sql"
	"time"

	"github.com/your-project/api"
	"github.com/your-project/impl"
	"github.com/your-project/jobs"
	"github.com/your-project/services"
)

func main() {
	// ... existing initialization code ...

	// Create repositories
	userRepo := impl.NewUserRepository(db)
	plantRepo := impl.NewPlantRepository(db)
	shopRepo := impl.NewShopRepository(db)
	notificationRepo := impl.NewNotificationRepository(db)

	// Create services
	userService := services.NewUserService(userRepo)
	plantService := services.NewPlantService(plantRepo)
	shopService := services.NewShopService(shopRepo)
	notificationService := services.NewNotificationService(notificationRepo, plantRepo)

	// Create and start background jobs
	wateringJob := jobs.NewWateringNotificationsJob(notificationService, 1*time.Hour)
	wateringJob.Start()
	defer wateringJob.Stop()

	// Create API
	api := api.NewAPI(
		userService,
		plantService,
		shopService,
		notificationService,
	)

	// ... existing server code ...
} 