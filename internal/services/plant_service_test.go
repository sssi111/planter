package services

import (
	"context"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlantRepository is a mock implementation of the PlantRepository interface
type MockPlantRepository struct {
	mock.Mock
}

func (m *MockPlantRepository) GetAll(ctx context.Context) ([]*models.Plant, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Plant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) Search(ctx context.Context, query string) ([]*models.Plant, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) GetFavorites(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) AddToFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantRepository) RemoveFromFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantRepository) MarkAsWatered(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantRepository) GetUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (*models.UserPlant, error) {
	args := m.Called(ctx, userID, plantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserPlant), args.Error(1)
}

func (m *MockPlantRepository) GetUserPlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) AddUserPlant(ctx context.Context, userPlant *models.UserPlant) error {
	args := m.Called(ctx, userPlant)
	return args.Error(0)
}

func (m *MockPlantRepository) UpdateUserPlant(ctx context.Context, userPlant *models.UserPlant) error {
	args := m.Called(ctx, userPlant)
	return args.Error(0)
}

func (m *MockPlantRepository) RemoveUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantRepository) IsFavorite(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, plantID)
	return args.Bool(0), args.Error(1)
}

// TestPlantService_GetAllPlants tests the GetAllPlants method of the PlantService
func TestPlantService_GetAllPlants(t *testing.T) {
	// Create a mock plant repository
	mockPlantRepo := new(MockPlantRepository)

	// Create test plants
	plant1 := &models.Plant{
		ID:          uuid.New(),
		Name:        "Plant 1",
		Description: "Description 1",
	}
	plant2 := &models.Plant{
		ID:          uuid.New(),
		Name:        "Plant 2",
		Description: "Description 2",
	}
	plants := []*models.Plant{plant1, plant2}

	// Set up the mock expectations
	mockPlantRepo.On("GetAll", mock.Anything).Return(plants, nil)

	// Create the plant service
	plantService := NewPlantService(mockPlantRepo)

	// Test the GetAllPlants method
	result, err := plantService.GetAllPlants(context.Background())

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected plants
	assert.Equal(t, plants, result)

	// Verify that all expectations were met
	mockPlantRepo.AssertExpectations(t)
}

// TestPlantService_GetPlant tests the GetPlant method of the PlantService
func TestPlantService_GetPlant(t *testing.T) {
	// Create a mock plant repository
	mockPlantRepo := new(MockPlantRepository)

	// Create a test plant
	plantID := uuid.New()
	plant := &models.Plant{
		ID:          plantID,
		Name:        "Test Plant",
		Description: "Test Description",
	}

	// Set up the mock expectations
	mockPlantRepo.On("GetByID", mock.Anything, plantID).Return(plant, nil)

	// Create the plant service
	plantService := NewPlantService(mockPlantRepo)

	// Test the GetPlant method
	result, err := plantService.GetPlant(context.Background(), plantID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected plant
	assert.Equal(t, plant, result)

	// Verify that all expectations were met
	mockPlantRepo.AssertExpectations(t)
}