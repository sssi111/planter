package impl

import (
    "context"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/anpanovv/planter/internal/db"
    "github.com/anpanovv/planter/internal/models"
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"
)

func setupNotificationTest(t *testing.T) (*NotificationRepository, sqlmock.Sqlmock, func()) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock DB: %v", err)
    }

    sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
    db := &db.DB{DB: sqlxDB}
    repo := NewNotificationRepository(db)

    return repo, mock, func() {
        mockDB.Close()
    }
}

func TestNotificationRepository_Create(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    notification := &models.Notification{
        UserID:  uuid.New(),
        PlantID: uuid.New(),
        Type:    models.NotificationTypeWatering,
        Message: "Test notification",
        IsRead:  false,
    }

    mock.ExpectExec("INSERT INTO notifications").
        WithArgs(notification.UserID, notification.PlantID, notification.Type, notification.Message, notification.IsRead).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err := repo.Create(context.Background(), notification)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetUserNotifications(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    userID := uuid.New()
    expectedTotal := 1
    expectedNotification := &models.Notification{
        ID:      uuid.New(),
        UserID:  userID,
        PlantID: uuid.New(),
        Type:    models.NotificationTypeWatering,
        Message: "Test notification",
        IsRead:  false,
        Plant: &models.Plant{
            ID:   uuid.New(),
            Name: "Test Plant",
        },
    }

    // Expect count query
    mock.ExpectQuery("SELECT COUNT").
        WithArgs(userID).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

    // Expect notifications query
    rows := sqlmock.NewRows([]string{
        "id", "user_id", "plant_id", "type", "message", "is_read", "created_at", "updated_at",
        "plant.id", "plant.name", "plant.scientific_name", "plant.image_url",
    }).AddRow(
        expectedNotification.ID, expectedNotification.UserID, expectedNotification.PlantID,
        expectedNotification.Type, expectedNotification.Message, expectedNotification.IsRead,
        time.Now(), time.Now(),
        expectedNotification.Plant.ID, expectedNotification.Plant.Name,
        "Scientific Name", "image.jpg",
    )

    mock.ExpectQuery("SELECT n.*, p.id").
        WithArgs(userID, 10, 0).
        WillReturnRows(rows)

    notifications, total, err := repo.GetUserNotifications(context.Background(), userID, 0, 10)
    assert.NoError(t, err)
    assert.Equal(t, expectedTotal, total)
    assert.Len(t, notifications, 1)
    assert.Equal(t, expectedNotification.ID, notifications[0].ID)
    assert.Equal(t, expectedNotification.Plant.ID, notifications[0].Plant.ID)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetUserNotifications_NullPlant(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    userID := uuid.New()
    expectedTotal := 1
    expectedNotification := &models.Notification{
        ID:      uuid.New(),
        UserID:  userID,
        PlantID: uuid.New(),
        Type:    models.NotificationTypeWatering,
        Message: "Test notification",
        IsRead:  false,
    }

    // Expect count query
    mock.ExpectQuery("SELECT COUNT").
        WithArgs(userID).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

    // Expect notifications query with NULL plant fields
    rows := sqlmock.NewRows([]string{
        "id", "user_id", "plant_id", "type", "message", "is_read", "created_at", "updated_at",
        "plant.id", "plant.name", "plant.scientific_name", "plant.image_url",
    }).AddRow(
        expectedNotification.ID, expectedNotification.UserID, expectedNotification.PlantID,
        expectedNotification.Type, expectedNotification.Message, expectedNotification.IsRead,
        time.Now(), time.Now(),
        nil, nil, nil, nil,
    )

    mock.ExpectQuery("SELECT n.id").
        WithArgs(userID, 10, 0).
        WillReturnRows(rows)

    notifications, total, err := repo.GetUserNotifications(context.Background(), userID, 0, 10)
    assert.NoError(t, err)
    assert.Equal(t, expectedTotal, total)
    assert.Len(t, notifications, 1)
    assert.Equal(t, expectedNotification.ID, notifications[0].ID)
    assert.Nil(t, notifications[0].Plant)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_MarkAsRead(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    notificationID := uuid.New()
    userID := uuid.New()

    mock.ExpectExec("UPDATE notifications").
        WithArgs(notificationID, userID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err := repo.MarkAsRead(context.Background(), notificationID, userID)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetUnreadWateringNotifications(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    expectedNotification := &models.Notification{
        ID:      uuid.New(),
        UserID:  uuid.New(),
        PlantID: uuid.New(),
        Type:    models.NotificationTypeWatering,
        Message: "Test notification",
        IsRead:  false,
        Plant: &models.Plant{
            ID:   uuid.New(),
            Name: "Test Plant",
        },
    }

    rows := sqlmock.NewRows([]string{
        "id", "user_id", "plant_id", "type", "message", "is_read", "created_at", "updated_at",
        "plant.id", "plant.name", "plant.scientific_name", "plant.image_url",
    }).AddRow(
        expectedNotification.ID, expectedNotification.UserID, expectedNotification.PlantID,
        expectedNotification.Type, expectedNotification.Message, expectedNotification.IsRead,
        time.Now(), time.Now(),
        expectedNotification.Plant.ID, expectedNotification.Plant.Name,
        "Scientific Name", "image.jpg",
    )

    mock.ExpectQuery("SELECT n.*, p.id").
        WithArgs(models.NotificationTypeWatering).
        WillReturnRows(rows)

    notifications, err := repo.GetUnreadWateringNotifications(context.Background())
    assert.NoError(t, err)
    assert.Len(t, notifications, 1)
    assert.Equal(t, expectedNotification.ID, notifications[0].ID)
    assert.Equal(t, expectedNotification.Plant.ID, notifications[0].Plant.ID)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNotificationRepository_GetUnreadWateringNotifications_NullPlant(t *testing.T) {
    repo, mock, cleanup := setupNotificationTest(t)
    defer cleanup()

    expectedNotification := &models.Notification{
        ID:      uuid.New(),
        UserID:  uuid.New(),
        PlantID: uuid.New(),
        Type:    models.NotificationTypeWatering,
        Message: "Test notification",
        IsRead:  false,
    }

    rows := sqlmock.NewRows([]string{
        "id", "user_id", "plant_id", "type", "message", "is_read", "created_at", "updated_at",
        "plant.id", "plant.name", "plant.scientific_name", "plant.image_url",
    }).AddRow(
        expectedNotification.ID, expectedNotification.UserID, expectedNotification.PlantID,
        expectedNotification.Type, expectedNotification.Message, expectedNotification.IsRead,
        time.Now(), time.Now(),
        nil, nil, nil, nil,
    )

    mock.ExpectQuery("SELECT n.id").
        WithArgs(models.NotificationTypeWatering).
        WillReturnRows(rows)

    notifications, err := repo.GetUnreadWateringNotifications(context.Background())
    assert.NoError(t, err)
    assert.Len(t, notifications, 1)
    assert.Equal(t, expectedNotification.ID, notifications[0].ID)
    assert.Nil(t, notifications[0].Plant)
    assert.NoError(t, mock.ExpectationsWereMet())
} 