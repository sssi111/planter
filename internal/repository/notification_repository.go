package repository

import (
    "context"
    "github.com/anpanovv/planter/internal/models"
    "github.com/google/uuid"
)

// NotificationRepository defines the interface for notification operations
type NotificationRepository interface {
    // Create creates a new notification
    Create(ctx context.Context, notification *models.Notification) error

    // GetUserNotifications gets all notifications for a user with pagination
    GetUserNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Notification, int, error)

    // MarkAsRead marks a notification as read
    MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error

    // GetUnreadWateringNotifications gets all unread watering notifications that need to be sent
    GetUnreadWateringNotifications(ctx context.Context) ([]*models.Notification, error)
} 