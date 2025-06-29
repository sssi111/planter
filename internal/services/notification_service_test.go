package services

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/anpanovv/planter/internal/models"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockNotificationRepository is a mock implementation of the NotificationRepository interface
type MockNotificationRepository struct {
    mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
    args := m.Called(ctx, notification)
    return args.Error(0)
}

func (m *MockNotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Notification, int, error) {
    args := m.Called(ctx, userID, offset, limit)
    return args.Get(0).([]*models.Notification), args.Int(1), args.Error(2)
}

func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
    args := m.Called(ctx, notificationID, userID)
    return args.Error(0)
}

func (m *MockNotificationRepository) GetUnreadWateringNotifications(ctx context.Context) ([]*models.Notification, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *MockPlantRepository) GetAllUserPlantsForWateringCheck(ctx context.Context) ([]*models.UserPlant, error) {
    args := m.Called(ctx)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*models.UserPlant), args.Error(1)
}

func TestNotificationService_GetUserNotifications(t *testing.T) {
    // Create mocks
    mockNotificationRepo := new(MockNotificationRepository)
    mockPlantRepo := new(MockPlantRepository)

    // Create service
    service := NewNotificationService(mockNotificationRepo, mockPlantRepo)

    // Test data
    ctx := context.Background()
    userID := uuid.New()
    notifications := []*models.Notification{
        {
            ID:      uuid.New(),
            UserID:  userID,
            Type:    models.NotificationTypeWatering,
            Message: "Test notification",
        },
    }
    total := 1

    // Set up expectations
    mockNotificationRepo.On("GetUserNotifications", ctx, userID, 0, 10).Return(notifications, total, nil)

    // Call the service
    response, err := service.GetUserNotifications(ctx, userID, 1, 10)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, notifications, response.Notifications)
    assert.Equal(t, total, response.Total)
    mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsRead(t *testing.T) {
    // Create mocks
    mockNotificationRepo := new(MockNotificationRepository)
    mockPlantRepo := new(MockPlantRepository)

    // Create service
    service := NewNotificationService(mockNotificationRepo, mockPlantRepo)

    // Test data
    ctx := context.Background()
    userID := uuid.New()
    notificationID := uuid.New()

    // Set up expectations
    mockNotificationRepo.On("MarkAsRead", ctx, notificationID, userID).Return(nil)

    // Call the service
    err := service.MarkAsRead(ctx, notificationID, userID)

    // Assert
    assert.NoError(t, err)
    mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_CheckAndCreateWateringNotifications(t *testing.T) {
    // Create mocks
    mockNotificationRepo := new(MockNotificationRepository)
    mockPlantRepo := new(MockPlantRepository)

    // Create service
    service := NewNotificationService(mockNotificationRepo, mockPlantRepo)

    // Test data
    ctx := context.Background()
    userID := uuid.New()
    nextWatering := time.Now().Add(-24 * time.Hour) // Plant needs watering
    
    userPlant := &models.UserPlant{
        ID:           uuid.New(),
        UserID:       userID,
        PlantID:      uuid.New(),
        NextWatering: &nextWatering,
        Plant: &models.Plant{
            ID:   uuid.New(),
            Name: "Test Plant",
        },
    }

    userPlants := []*models.UserPlant{userPlant}

    // Set up expectations
    mockPlantRepo.On("GetAllUserPlantsForWateringCheck", ctx).Return(userPlants, nil)
    mockNotificationRepo.On("Create", ctx, mock.MatchedBy(func(n *models.Notification) bool {
        return n.UserID == userID && n.PlantID == userPlant.PlantID && n.Type == models.NotificationTypeWatering
    })).Return(nil)

    // Call the service
    err := service.CheckAndCreateWateringNotifications(ctx)

    // Assert
    assert.NoError(t, err)
    mockPlantRepo.AssertExpectations(t)
    mockNotificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetUserNotifications_NoNotifications(t *testing.T) {
    // Create mocks
    mockNotificationRepo := new(MockNotificationRepository)
    mockPlantRepo := new(MockPlantRepository)

    // Create service
    service := NewNotificationService(mockNotificationRepo, mockPlantRepo)

    // Test data
    ctx := context.Background()
    userID := uuid.New()

    // Set up expectations - simulate database error or no notifications
    mockNotificationRepo.On("GetUserNotifications", ctx, userID, 0, 10).Return(nil, 0, fmt.Errorf("no notifications found"))

    // Call the service
    response, err := service.GetUserNotifications(ctx, userID, 1, 10)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Empty(t, response.Notifications)
    assert.Equal(t, 0, response.Total)
    mockNotificationRepo.AssertExpectations(t)
} 