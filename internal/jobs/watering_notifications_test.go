package jobs

import (
    "context"
    "testing"
    "time"

    "github.com/anpanovv/planter/internal/services"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of the notification service
type MockNotificationService struct {
    mock.Mock
}

func (m *MockNotificationService) CheckAndCreateWateringNotifications(ctx context.Context) error {
    args := m.Called(ctx)
    return args.Error(0)
}

func TestWateringNotificationsJob_Start(t *testing.T) {
    // Create mock service
    mockService := new(MockNotificationService)
    mockService.On("CheckAndCreateWateringNotifications", mock.Anything).Return(nil)

    // Create job with short interval for testing
    job := NewWateringNotificationsJob(mockService, 100*time.Millisecond)

    // Start job
    job.Start()

    // Wait for at least one execution
    time.Sleep(150 * time.Millisecond)

    // Stop job
    job.Stop()

    // Assert that the service was called at least once
    mockService.AssertNumberOfCalls(t, "CheckAndCreateWateringNotifications", 1)
}

func TestWateringNotificationsJob_Stop(t *testing.T) {
    // Create mock service
    mockService := new(MockNotificationService)

    // Create job
    job := NewWateringNotificationsJob(mockService, time.Hour)

    // Start and immediately stop
    job.Start()
    job.Stop()

    // Assert that the service was not called
    mockService.AssertNotCalled(t, "CheckAndCreateWateringNotifications")
}

func TestWateringNotificationsJob_CheckAndCreateNotifications(t *testing.T) {
    // Create mock service
    mockService := new(MockNotificationService)
    mockService.On("CheckAndCreateWateringNotifications", mock.Anything).Return(nil)

    // Create job
    job := NewWateringNotificationsJob(mockService, time.Hour)

    // Call check directly
    err := job.checkAndCreateNotifications()

    // Assert
    assert.NoError(t, err)
    mockService.AssertExpectations(t)
}

func TestWateringNotificationsJob_CheckAndCreateNotifications_Error(t *testing.T) {
    // Create mock service with error
    mockService := new(MockNotificationService)
    expectedError := assert.AnError
    mockService.On("CheckAndCreateWateringNotifications", mock.Anything).Return(expectedError)

    // Create job
    job := NewWateringNotificationsJob(mockService, time.Hour)

    // Call check directly
    err := job.checkAndCreateNotifications()

    // Assert
    assert.Error(t, err)
    assert.Equal(t, expectedError, err)
    mockService.AssertExpectations(t)
} 