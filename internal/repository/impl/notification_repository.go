package impl

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/anpanovv/planter/internal/db"
    "github.com/anpanovv/planter/internal/models"
    "github.com/google/uuid"
)

// NotificationRepository is the implementation of the notification repository
type NotificationRepository struct {
    db *db.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *db.DB) *NotificationRepository {
    return &NotificationRepository{
        db: db,
    }
}

// Create creates a new notification
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO notifications (user_id, plant_id, type, message, is_read)
        VALUES ($1, $2, $3, $4, $5)
    `, notification.UserID, notification.PlantID, notification.Type, notification.Message, notification.IsRead)
    if err != nil {
        return fmt.Errorf("failed to create notification: %w", err)
    }
    return nil
}

// GetUserNotifications gets all notifications for a user with pagination
func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Notification, int, error) {
    // Get total count
    var total int
    err := r.db.GetContext(ctx, &total, `
        SELECT COUNT(*)
        FROM notifications
        WHERE user_id = $1
    `, userID)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get notifications count: %w", err)
    }

    // Get notifications with plants
    rows, err := r.db.QueryxContext(ctx, `
        SELECT n.id, n.user_id, n.plant_id, n.type, n.message, n.is_read, n.created_at, n.updated_at,
               p.id as "plant.id", p.name as "plant.name", 
               p.scientific_name as "plant.scientific_name",
               p.image_url as "plant.image_url"
        FROM notifications n
        LEFT JOIN plants p ON n.plant_id = p.id
        WHERE n.user_id = $1
        ORDER BY n.created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
    }
    defer rows.Close()

    var notifications []*models.Notification
    for rows.Next() {
        var notification models.Notification
        var plantID, plantName, scientificName, imageURL sql.NullString
        err := rows.Scan(
            &notification.ID, &notification.UserID, &notification.PlantID,
            &notification.Type, &notification.Message, &notification.IsRead,
            &notification.CreatedAt, &notification.UpdatedAt,
            &plantID, &plantName, &scientificName, &imageURL,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
        }

        // Only set plant if we have valid plant data
        if plantID.Valid {
            notification.Plant = &models.Plant{
                ID:            uuid.MustParse(plantID.String),
                Name:          plantName.String,
                ScientificName: scientificName.String,
                ImageURL:      imageURL.String,
            }
        }

        notifications = append(notifications, &notification)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("error iterating notifications: %w", err)
    }

    return notifications, total, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
    result, err := r.db.ExecContext(ctx, `
        UPDATE notifications
        SET is_read = true, updated_at = NOW()
        WHERE id = $1 AND user_id = $2
    `, notificationID, userID)
    if err != nil {
        return fmt.Errorf("failed to mark notification as read: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rows == 0 {
        return fmt.Errorf("notification not found or not owned by user")
    }

    return nil
}

// GetUnreadWateringNotifications gets all unread watering notifications that need to be sent
func (r *NotificationRepository) GetUnreadWateringNotifications(ctx context.Context) ([]*models.Notification, error) {
    rows, err := r.db.QueryxContext(ctx, `
        SELECT n.id, n.user_id, n.plant_id, n.type, n.message, n.is_read, n.created_at, n.updated_at,
               p.id as "plant.id", p.name as "plant.name", 
               p.scientific_name as "plant.scientific_name",
               p.image_url as "plant.image_url"
        FROM notifications n
        LEFT JOIN plants p ON n.plant_id = p.id
        WHERE n.type = $1 AND n.is_read = false
        ORDER BY n.created_at DESC
    `, models.NotificationTypeWatering)
    if err != nil {
        return nil, fmt.Errorf("failed to get unread watering notifications: %w", err)
    }
    defer rows.Close()

    var notifications []*models.Notification
    for rows.Next() {
        var notification models.Notification
        var plantID, plantName, scientificName, imageURL sql.NullString
        err := rows.Scan(
            &notification.ID, &notification.UserID, &notification.PlantID,
            &notification.Type, &notification.Message, &notification.IsRead,
            &notification.CreatedAt, &notification.UpdatedAt,
            &plantID, &plantName, &scientificName, &imageURL,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan notification: %w", err)
        }

        // Only set plant if we have valid plant data
        if plantID.Valid {
            notification.Plant = &models.Plant{
                ID:            uuid.MustParse(plantID.String),
                Name:          plantName.String,
                ScientificName: scientificName.String,
                ImageURL:      imageURL.String,
            }
        }

        notifications = append(notifications, &notification)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating notifications: %w", err)
    }

    return notifications, nil
} 