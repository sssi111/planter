package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPlantService is a mock implementation of the plant service
type MockPlantService struct {
	mock.Mock
}

func (m *MockPlantService) GetAllPlants(ctx context.Context) ([]*models.Plant, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantService) GetPlant(ctx context.Context, plantID uuid.UUID) (*models.Plant, error) {
	args := m.Called(ctx, plantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

func (m *MockPlantService) SearchPlants(ctx context.Context, query string) ([]*models.Plant, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantService) GetFavoritePlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantService) AddToFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantService) RemoveFromFavorites(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantService) MarkAsWatered(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) (*models.Plant, error) {
	args := m.Called(ctx, userID, plantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

func (m *MockPlantService) GetUserPlants(ctx context.Context, userID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockPlantService) AddUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID, location string) error {
	args := m.Called(ctx, userID, plantID, location)
	return args.Error(0)
}

func (m *MockPlantService) UpdateUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID, location string) error {
	args := m.Called(ctx, userID, plantID, location)
	return args.Error(0)
}

func (m *MockPlantService) RemoveUserPlant(ctx context.Context, userID uuid.UUID, plantID uuid.UUID) error {
	args := m.Called(ctx, userID, plantID)
	return args.Error(0)
}

func (m *MockPlantService) CreatePlant(ctx context.Context, plant *models.Plant, careInstructions *models.CareInstructions) (*models.Plant, error) {
	args := m.Called(ctx, plant, careInstructions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

// TestAPI is a test implementation of the API
type TestAPI struct {
	plantService *MockPlantService
}

// handleAdminCreatePlant is a test implementation of the handleAdminCreatePlant handler
func (a *TestAPI) handleAdminCreatePlant(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req AdminPlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a new plant model
	plant := &models.Plant{
		Name:           req.Name,
		ScientificName: req.ScientificName,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		Price:          req.Price,
		ShopID:         req.ShopID,
	}

	// Create the plant
	createdPlant, err := a.plantService.CreatePlant(r.Context(), plant, &req.CareInstructions)
	if err != nil {
		http.Error(w, "Failed to create plant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created plant
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPlant)
}

// TestHandleAdminCreatePlant tests the handleAdminCreatePlant function
func TestHandleAdminCreatePlant(t *testing.T) {
	// Create a mock plant service
	mockPlantService := new(MockPlantService)

	// Create a test API
	api := &TestAPI{
		plantService: mockPlantService,
	}

	// Create a test request
	reqBody := AdminPlantRequest{
		Name:           "Test Plant",
		ScientificName: "Testus Plantus",
		Description:    "A test plant",
		ImageURL:       "https://example.com/test-plant.jpg",
		CareInstructions: models.CareInstructions{
			WateringFrequency:   7,
			Sunlight:            models.SunlightLevelMedium,
			Temperature:         models.TemperatureRange{Min: 18, Max: 24},
			Humidity:            models.HumidityLevelMedium,
			SoilType:            "Well-draining",
			FertilizerFrequency: 30,
			AdditionalNotes:     "Keep away from direct sunlight",
		},
	}

	// Expected result
	expectedPlant := &models.Plant{
		ID:               uuid.New(),
		Name:             "Test Plant",
		ScientificName:   "Testus Plantus",
		Description:      "A test plant",
		ImageURL:         "https://example.com/test-plant.jpg",
		CareInstructions: reqBody.CareInstructions,
	}

	// Set up the mock expectations
	mockPlantService.On("CreatePlant", mock.Anything, mock.AnythingOfType("*models.Plant"), mock.AnythingOfType("*models.CareInstructions")).
		Return(expectedPlant, nil)

	// Create a request
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/admin/plants", bytes.NewBuffer(reqBodyBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	api.handleAdminCreatePlant(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse the response
	var response models.Plant
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, expectedPlant.ID, response.ID)
	assert.Equal(t, expectedPlant.Name, response.Name)
	assert.Equal(t, expectedPlant.ScientificName, response.ScientificName)
	assert.Equal(t, expectedPlant.Description, response.Description)
	assert.Equal(t, expectedPlant.ImageURL, response.ImageURL)
	assert.Equal(t, expectedPlant.CareInstructions.WateringFrequency, response.CareInstructions.WateringFrequency)
	assert.Equal(t, expectedPlant.CareInstructions.Sunlight, response.CareInstructions.Sunlight)
	assert.Equal(t, expectedPlant.CareInstructions.Humidity, response.CareInstructions.Humidity)
	assert.Equal(t, expectedPlant.CareInstructions.SoilType, response.CareInstructions.SoilType)
	assert.Equal(t, expectedPlant.CareInstructions.FertilizerFrequency, response.CareInstructions.FertilizerFrequency)
	assert.Equal(t, expectedPlant.CareInstructions.AdditionalNotes, response.CareInstructions.AdditionalNotes)

	// Verify that all expectations were met
	mockPlantService.AssertExpectations(t)
}
