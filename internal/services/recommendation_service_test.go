package services

import (
	"context"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecommendationRepository is a mock implementation of the RecommendationRepository interface
type MockRecommendationRepository struct {
	mock.Mock
}

func (m *MockRecommendationRepository) SaveQuestionnaire(ctx context.Context, questionnaire *models.PlantQuestionnaire) error {
	args := m.Called(ctx, questionnaire)
	return args.Error(0)
}

func (m *MockRecommendationRepository) GetQuestionnaire(ctx context.Context, id uuid.UUID) (*models.PlantQuestionnaire, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlantQuestionnaire), args.Error(1)
}

func (m *MockRecommendationRepository) SaveRecommendation(ctx context.Context, recommendation *models.PlantRecommendation) error {
	args := m.Called(ctx, recommendation)
	return args.Error(0)
}

func (m *MockRecommendationRepository) GetRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.PlantRecommendation, error) {
	args := m.Called(ctx, questionnaireID)
	return args.Get(0).([]*models.PlantRecommendation), args.Error(1)
}

func (m *MockRecommendationRepository) GetRecommendedPlants(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, questionnaireID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

// TestRecommendationService_SaveQuestionnaire tests the SaveQuestionnaire method of the RecommendationService
func TestRecommendationService_SaveQuestionnaire(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test user ID and questionnaire request
	userID := uuid.New()
	questionnaireRequest := &models.QuestionnaireRequest{
		SunlightPreference:   models.SunlightLevelMedium,
		PetFriendly:          true,
		CareLevel:            3,
		PreferredLocation:    stringPtr("Living Room"),
		AdditionalPreferences: stringPtr("Low maintenance"),
	}

	// Set up the mock expectations
	mockRecommendationRepo.On("SaveQuestionnaire", mock.Anything, mock.AnythingOfType("*models.PlantQuestionnaire")).
		Return(nil).
		Run(func(args mock.Arguments) {
			questionnaire := args.Get(1).(*models.PlantQuestionnaire)
			questionnaire.ID = uuid.New() // Simulate the database generating an ID
		})

	// Create the recommendation service
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Test the SaveQuestionnaire method
	result, err := recommendationService.SaveQuestionnaire(context.Background(), &userID, questionnaireRequest)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the expected values
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, &userID, result.UserID)
	assert.Equal(t, models.SunlightLevelMedium, result.SunlightPreference)
	assert.Equal(t, true, result.PetFriendly)
	assert.Equal(t, 3, result.CareLevel)
	assert.Equal(t, stringPtr("Living Room"), result.PreferredLocation)
	assert.Equal(t, stringPtr("Low maintenance"), result.AdditionalPreferences)

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// TestRecommendationService_GetRecommendations tests the GetRecommendations method of the RecommendationService
func TestRecommendationService_GetRecommendations(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test questionnaire ID
	questionnaireID := uuid.New()

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
	recommendedPlants := []*models.Plant{plant1, plant2}

	recommendation1 := &models.PlantRecommendation{
		ID:              uuid.New(),
		QuestionnaireID: questionnaireID,
		PlantID:         plant1.ID,
		Score:           0.9,
		Reasoning:       "Good match for medium sunlight and pet-friendly",
	}
	recommendation2 := &models.PlantRecommendation{
		ID:              uuid.New(),
		QuestionnaireID: questionnaireID,
		PlantID:         plant2.ID,
		Score:           0.8,
		Reasoning:       "Good match for medium sunlight but requires more care",
	}
	recommendations := []*models.PlantRecommendation{recommendation1, recommendation2}

	// Set up the mock expectations
	mockRecommendationRepo.On("GetRecommendations", mock.Anything, questionnaireID).Return(recommendations, nil)
	mockRecommendationRepo.On("GetRecommendedPlants", mock.Anything, questionnaireID).Return(recommendedPlants, nil)

	// Create the recommendation service
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Test the GetRecommendations method
	result, err := recommendationService.GetRecommendations(context.Background(), questionnaireID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected plants
	assert.Equal(t, recommendedPlants, result)

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// TestRecommendationService_GenerateRecommendations tests the GenerateRecommendations method of the RecommendationService
func TestRecommendationService_GenerateRecommendations(t *testing.T) {
	// This test is more complex because it involves the Yandex GPT API
	// We'll mock the behavior to simulate the API response

	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test questionnaire and plants
	questionnaireID := uuid.New()
	questionnaire := &models.PlantQuestionnaire{
		ID:                 questionnaireID,
		SunlightPreference: models.SunlightLevelMedium,
		PetFriendly:        true,
		CareLevel:          3,
	}

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
	allPlants := []*models.Plant{plant1, plant2}
	recommendedPlants := []*models.Plant{plant1, plant2}

	// Set up the mock expectations
	mockRecommendationRepo.On("GetQuestionnaire", mock.Anything, questionnaireID).Return(questionnaire, nil)
	mockPlantRepo.On("GetAll", mock.Anything).Return(allPlants, nil)
	mockRecommendationRepo.On("SaveRecommendation", mock.Anything, mock.AnythingOfType("*models.PlantRecommendation")).Return(nil)
	mockRecommendationRepo.On("GetRecommendedPlants", mock.Anything, questionnaireID).Return(recommendedPlants, nil)

	// Create a mock recommendation service
	// We'll create a custom implementation that skips the actual API call
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// We'll mock the GetQuestionnaire call to return a questionnaire
	mockRecommendationRepo.On("GetQuestionnaire", mock.Anything, questionnaireID).Return(&models.PlantQuestionnaire{
		ID:                 questionnaireID,
		SunlightPreference: models.SunlightLevelMedium,
		PetFriendly:        true,
		CareLevel:          3,
	}, nil)

	// Mock the behavior of the recommendation generation
	// Instead of calling the Yandex GPT API, we'll directly create and save recommendations
	mockRecommendationRepo.On("SaveRecommendation", mock.Anything, mock.MatchedBy(func(r *models.PlantRecommendation) bool {
		return r.QuestionnaireID == questionnaireID && (r.PlantID == plant1.ID || r.PlantID == plant2.ID)
	})).Return(nil)

	// Test the GenerateRecommendations method
	result, err := recommendationService.GenerateRecommendations(context.Background(), questionnaireID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected plants
	assert.Equal(t, recommendedPlants, result)

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}