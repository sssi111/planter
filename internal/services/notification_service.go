package services

import (
    "context"
    "fmt"
    "time"

    "github.com/anpanovv/planter/internal/models"
    "github.com/anpanovv/planter/internal/repository"
    "github.com/google/uuid"
)

// NotificationStats contains statistics about notification processing
type NotificationStats struct {
    UsersProcessed      int
    PlantsNeedingWater int
    NotificationsCreated int
}

// NotificationService handles notification operations
type NotificationService struct {
    notificationRepo repository.NotificationRepository
    plantRepo       repository.PlantRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(notificationRepo repository.NotificationRepository, plantRepo repository.PlantRepository) *NotificationService {
    return &NotificationService{
        notificationRepo: notificationRepo,
        plantRepo:       plantRepo,
    }
}

// GetUserNotifications gets all notifications for a user with pagination
func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, page, pageSize int) (*models.NotificationResponse, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize
    notifications, total, err := s.notificationRepo.GetUserNotifications(ctx, userID, offset, pageSize)
    if err != nil {
        return &models.NotificationResponse{
            Notifications: []*models.Notification{},
            Total:        0,
        }, nil
    }

    return &models.NotificationResponse{
        Notifications: notifications,
        Total:        total,
    }, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
    err := s.notificationRepo.MarkAsRead(ctx, notificationID, userID)
    if err != nil {
        return fmt.Errorf("failed to mark notification as read: %w", err)
    }
    return nil
}

// CheckAndCreateWateringNotifications checks for plants that need watering and creates notifications
func (s *NotificationService) CheckAndCreateWateringNotifications(ctx context.Context) (*NotificationStats, error) {
    stats := &NotificationStats{}
    userSet := make(map[uuid.UUID]struct{})

    // Get all user plants
    userPlants, err := s.plantRepo.GetAllUserPlantsForWateringCheck(ctx)
    if err != nil {
    	return nil, fmt.Errorf("failed to get plants for watering check: %w", err)
    }
   
    now := time.Now()
    for _, userPlant := range userPlants {
    	if userPlant.NextWatering != nil && userPlant.NextWatering.Before(now) {
            stats.PlantsNeedingWater++
            userSet[userPlant.UserID] = struct{}{}

    		// Create notification
    		notification := &models.Notification{
    			UserID:  userPlant.UserID,
    			PlantID: userPlant.PlantID,
    			Type:    models.NotificationTypeWatering,
    			Message: fmt.Sprintf("Пора полить ваше растение %s!", userPlant.Plant.Name),
    			IsRead:  false,
    		}
   
    		err := s.notificationRepo.Create(ctx, notification)
    		if err != nil {
    			return nil, fmt.Errorf("failed to create watering notification: %w", err)
    		}
            stats.NotificationsCreated++
    	}
    }

    stats.UsersProcessed = len(userSet)
    return stats, nil
}