package services

import (
	"context"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is already defined in auth_service_test.go

// TestUserService_GetUser tests the GetUser method of the UserService
func TestUserService_GetUser(t *testing.T) {
	// Create a mock user repository
	mockUserRepo := new(MockUserRepository)

	// Create a test user
	userID := uuid.New()
	user := &models.User{
		ID:                  userID,
		Name:                "Test User",
		Email:               "test@example.com",
		Language:            models.LanguageRussian,
		NotificationsEnabled: true,
	}

	// Set up the mock expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Create the user service
	userService := NewUserService(mockUserRepo)

	// Test the GetUser method
	result, err := userService.GetUser(context.Background(), userID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected user
	assert.Equal(t, user, result)

	// Verify that all expectations were met
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_UpdateUser tests the UpdateUser method of the UserService
func TestUserService_UpdateUser(t *testing.T) {
	// Create a mock user repository
	mockUserRepo := new(MockUserRepository)

	// Create a test user
	userID := uuid.New()
	existingUser := &models.User{
		ID:                  userID,
		Name:                "Test User",
		Email:               "test@example.com",
		Language:            models.LanguageRussian,
		NotificationsEnabled: true,
	}

	updatedUser := &models.User{
		ID:                  userID,
		Name:                "Updated User",
		Email:               "test@example.com",
		Language:            models.LanguageEnglish,
		NotificationsEnabled: false,
	}

	// Set up the mock expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	// Create the user service
	userService := NewUserService(mockUserRepo)

	// Test the UpdateUser method
	result, err := userService.UpdateUser(context.Background(), updatedUser)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the updated values
	assert.Equal(t, "Updated User", result.Name)
	assert.Equal(t, models.LanguageEnglish, result.Language)
	assert.Equal(t, false, result.NotificationsEnabled)

	// Verify that all expectations were met
	mockUserRepo.AssertExpectations(t)
}