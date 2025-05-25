package services

import (
	"context"
	"testing"

	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetLocations(ctx context.Context, userID uuid.UUID) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepository) AddLocation(ctx context.Context, userID uuid.UUID, location string) error {
	args := m.Called(ctx, userID, location)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveLocation(ctx context.Context, userID uuid.UUID, location string) error {
	args := m.Called(ctx, userID, location)
	return args.Error(0)
}

func (m *MockUserRepository) GetFavoritePlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepository) GetOwnedPlantIDs(ctx context.Context, userID uuid.UUID) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

// TestAuthService_Login tests the Login method of the AuthService
func TestAuthService_Login(t *testing.T) {
	// Create a mock user repository
	mockUserRepo := new(MockUserRepository)

	// Create a test user
	userID := uuid.New()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		ID:                  userID,
		Name:                "Test User",
		Email:               "test@example.com",
		PasswordHash:        string(hashedPassword),
		Language:            models.LanguageRussian,
		NotificationsEnabled: true,
	}

	// Set up the mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	// Create a mock auth middleware
	auth := middleware.NewAuth("test-secret")

	// Create the auth service
	authService := NewAuthService(mockUserRepo, auth)

	// Test the login method
	resp, err := authService.Login(context.Background(), "test@example.com", password)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the response contains a token and the user
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, userID, resp.User.ID)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Empty(t, resp.User.PasswordHash) // Password hash should be hidden

	// Verify that all expectations were met
	mockUserRepo.AssertExpectations(t)
}

// TestAuthService_Register tests the Register method of the AuthService
func TestAuthService_Register(t *testing.T) {
	// Create a mock user repository
	mockUserRepo := new(MockUserRepository)

	// Set up the mock expectations
	mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, assert.AnError)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = uuid.New() // Simulate the database generating an ID
	})

	// Create a mock auth middleware
	auth := middleware.NewAuth("test-secret")

	// Create the auth service
	authService := NewAuthService(mockUserRepo, auth)

	// Test the register method
	resp, err := authService.Register(context.Background(), "Test User", "test@example.com", "password123")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the response contains a token and the user
	assert.NotEmpty(t, resp.Token)
	assert.NotEqual(t, uuid.Nil, resp.User.ID)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Empty(t, resp.User.PasswordHash) // Password hash should be hidden

	// Verify that all expectations were met
	mockUserRepo.AssertExpectations(t)
}