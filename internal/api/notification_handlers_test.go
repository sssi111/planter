package api

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/anpanovv/planter/internal/middleware"
    "github.com/anpanovv/planter/internal/models"
    "github.com/anpanovv/planter/internal/services"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of the notification service
type MockNotificationService struct {
    mock.Mock
}

func (m *MockNotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, page, pageSize int) (*models.NotificationResponse, error) {
    args := m.Called(ctx, userID, page, pageSize)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.NotificationResponse), args.Error(1)
}

func (m *MockNotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
    args := m.Called(ctx, notificationID, userID)
    return args.Error(0)
}

func (m *MockNotificationService) CheckAndCreateWateringNotifications(ctx context.Context) error {
    args := m.Called(ctx)
    return args.Error(0)
}

func TestHandleGetUserNotifications(t *testing.T) {
    // Create test data
    userID := uuid.New()
    notifications := []*models.Notification{
        {
            ID:      uuid.New(),
            UserID:  userID,
            Type:    models.NotificationTypeWatering,
            Message: "Test notification",
            Plant: &models.Plant{
                ID:   uuid.New(),
                Name: "Test Plant",
            },
        },
    }
    response := &models.NotificationResponse{
        Notifications: notifications,
        Total:        1,
    }

    // Create mocks
    mockService := new(MockNotificationService)
    mockService.On("GetUserNotifications", mock.Anything, userID, 1, 10).Return(response, nil)

    // Create API instance
    api := &API{
        notificationService: mockService,
    }

    // Create request
    req := httptest.NewRequest("GET", "/notifications?page=1&pageSize=10", nil)
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call handler
    api.handleGetUserNotifications(rr, req)

    // Assert response
    assert.Equal(t, http.StatusOK, rr.Code)

    var result models.NotificationResponse
    err := json.Unmarshal(rr.Body.Bytes(), &result)
    assert.NoError(t, err)
    assert.Equal(t, 1, len(result.Notifications))
    assert.Equal(t, notifications[0].ID, result.Notifications[0].ID)

    mockService.AssertExpectations(t)
}

func TestHandleMarkNotificationAsRead(t *testing.T) {
    // Create test data
    userID := uuid.New()
    notificationID := uuid.New()

    // Create mocks
    mockService := new(MockNotificationService)
    mockService.On("MarkAsRead", mock.Anything, notificationID, userID).Return(nil)

    // Create API instance
    api := &API{
        notificationService: mockService,
    }

    // Create request
    req := httptest.NewRequest("POST", "/notifications/"+notificationID.String()+"/read", nil)
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

    // Add URL parameters
    vars := map[string]string{
        "notificationId": notificationID.String(),
    }
    req = mux.SetURLVars(req, vars)

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call handler
    api.handleMarkNotificationAsRead(rr, req)

    // Assert response
    assert.Equal(t, http.StatusOK, rr.Code)

    var result map[string]string
    err := json.Unmarshal(rr.Body.Bytes(), &result)
    assert.NoError(t, err)
    assert.Equal(t, "Notification marked as read", result["message"])

    mockService.AssertExpectations(t)
}

func TestHandleGetUserNotifications_Unauthorized(t *testing.T) {
    // Create API instance
    api := &API{
        notificationService: new(MockNotificationService),
    }

    // Create request without user ID in context
    req := httptest.NewRequest("GET", "/notifications", nil)
    rr := httptest.NewRecorder()

    // Call handler
    api.handleGetUserNotifications(rr, req)

    // Assert response
    assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleMarkNotificationAsRead_InvalidID(t *testing.T) {
    // Create test data
    userID := uuid.New()

    // Create API instance
    api := &API{
        notificationService: new(MockNotificationService),
    }

    // Create request with invalid notification ID
    req := httptest.NewRequest("POST", "/notifications/invalid-uuid/read", nil)
    req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

    // Add URL parameters
    vars := map[string]string{
        "notificationId": "invalid-uuid",
    }
    req = mux.SetURLVars(req, vars)

    rr := httptest.NewRecorder()

    // Call handler
    api.handleMarkNotificationAsRead(rr, req)

    // Assert response
    assert.Equal(t, http.StatusBadRequest, rr.Code)
} 