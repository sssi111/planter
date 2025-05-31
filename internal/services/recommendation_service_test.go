package services

import (
	"context"
	"testing"
	"time"

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

func (m *MockRecommendationRepository) SaveDetailedQuestionnaire(ctx context.Context, questionnaire *models.DetailedQuestionnaireRequest) (*models.PlantQuestionnaire, error) {
	args := m.Called(ctx, questionnaire)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlantQuestionnaire), args.Error(1)
}

func (m *MockRecommendationRepository) CreateChatSession(ctx context.Context, userID uuid.UUID, title string) (*models.ChatSession, error) {
	args := m.Called(ctx, userID, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationRepository) GetChatSession(ctx context.Context, id uuid.UUID) (*models.ChatSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationRepository) GetChatSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.ChatSession, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationRepository) SaveChatMessage(ctx context.Context, message *models.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockRecommendationRepository) GetChatMessages(ctx context.Context, sessionID uuid.UUID) ([]*models.ChatMessage, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]*models.ChatMessage), args.Error(1)
}

func (m *MockRecommendationRepository) UpdateChatSessionLastUsed(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// TestRecommendationService_SaveQuestionnaire tests the SaveQuestionnaire method of the RecommendationService
func TestRecommendationService_SaveQuestionnaire(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test user ID and questionnaire request
	userID := uuid.New()
	questionnaireRequest := &models.QuestionnaireRequest{
		SunlightPreference:    models.SunlightLevelMedium,
		PetFriendly:           true,
		CareLevel:             3,
		PreferredLocation:     stringPtr("Living Room"),
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

// TestRecommendationService_SaveDetailedQuestionnaire tests the SaveDetailedQuestionnaire method
func TestRecommendationService_SaveDetailedQuestionnaire(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test user ID and detailed questionnaire request
	userID := uuid.New()
	detailedQuestionnaireRequest := &models.DetailedQuestionnaireRequest{
		SunlightPreference:    models.SunlightLevelMedium,
		PetFriendly:           true,
		CareLevel:             3,
		PreferredLocation:     stringPtr("Living Room"),
		HasChildren:           true,
		PlantSize:             "MEDIUM",
		FloweringPreference:   true,
		AirPurifying:          true,
		WateringFrequency:     "REGULAR",
		ExperienceLevel:       "BEGINNER",
		AdditionalPreferences: stringPtr("Low maintenance"),
	}

	// Expected plant questionnaire that would be created
	expectedQuestionnaire := &models.PlantQuestionnaire{
		ID:                 uuid.New(),
		UserID:             &userID,
		SunlightPreference: models.SunlightLevelMedium,
		PetFriendly:        true,
		CareLevel:          3,
		PreferredLocation:  stringPtr("Living Room"),
	}

	// Set up the mock expectations
	mockRecommendationRepo.On("SaveQuestionnaire", mock.Anything, mock.AnythingOfType("*models.PlantQuestionnaire")).
		Return(nil).
		Run(func(args mock.Arguments) {
			questionnaire := args.Get(1).(*models.PlantQuestionnaire)
			questionnaire.ID = expectedQuestionnaire.ID // Simulate the database generating an ID
		})

	// Create the recommendation service
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Test the SaveDetailedQuestionnaire method
	result, err := recommendationService.SaveDetailedQuestionnaire(context.Background(), &userID, detailedQuestionnaireRequest)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the expected values
	assert.Equal(t, expectedQuestionnaire.ID, result.ID)
	assert.Equal(t, &userID, result.UserID)
	assert.Equal(t, models.SunlightLevelMedium, result.SunlightPreference)
	assert.Equal(t, true, result.PetFriendly)
	assert.Equal(t, 3, result.CareLevel)
	assert.Equal(t, stringPtr("Living Room"), result.PreferredLocation)

	// Check that AdditionalPreferences contains all the detailed information
	assert.NotNil(t, result.AdditionalPreferences)
	assert.Contains(t, *result.AdditionalPreferences, "Размер растения: MEDIUM")
	assert.Contains(t, *result.AdditionalPreferences, "Цветущее: true")
	assert.Contains(t, *result.AdditionalPreferences, "Очищающее воздух: true")
	assert.Contains(t, *result.AdditionalPreferences, "Частота полива: REGULAR")
	assert.Contains(t, *result.AdditionalPreferences, "Уровень опыта: BEGINNER")
	assert.Contains(t, *result.AdditionalPreferences, "Есть дети: true")
	assert.Contains(t, *result.AdditionalPreferences, "Low maintenance")

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// TestRecommendationService_CreateChatSession tests the CreateChatSession method
func TestRecommendationService_CreateChatSession(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create a test user ID
	userID := uuid.New()

	// Expected chat session
	expectedSession := &models.ChatSession{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Разговор о растениях",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Set up the mock expectations
	mockRecommendationRepo.On("CreateChatSession", mock.Anything, userID, "Разговор о растениях").
		Return(expectedSession, nil)

	// Create the recommendation service
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Test the CreateChatSession method
	result, err := recommendationService.CreateChatSession(context.Background(), userID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the expected values
	assert.Equal(t, expectedSession.ID, result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "Разговор о растениях", result.Title)

	// Verify that the in-memory session was initialized with a system message
	sessionMessages, ok := recommendationService.chatSessions[result.ID]
	assert.True(t, ok)
	assert.Equal(t, 1, len(sessionMessages))
	assert.Equal(t, "system", sessionMessages[0].Role)
	assert.Contains(t, sessionMessages[0].Text, "эксперт по растениям")

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// MockRecommendationServiceWithResponse is a mock implementation of the RecommendationService
// that returns a fixed response for the callYandexGPTAPI method
type MockRecommendationServiceWithResponse struct {
	*RecommendationService
	fixedResponse string
}

// callYandexGPTAPI is a mock implementation that returns a fixed response
func (m *MockRecommendationServiceWithResponse) callYandexGPTAPI(ctx context.Context, prompt string, messages []Message) (string, error) {
	return m.fixedResponse, nil
}

// TestRecommendationService_SendChatMessage tests the SendChatMessage method
func TestRecommendationService_SendChatMessage(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create test data
	userID := uuid.New()
	sessionID := uuid.New()
	userMessage := "Какие растения подходят для темной комнаты?"
	
	// Expected chat session
	session := &models.ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     "Разговор о растениях",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Expected assistant response
	assistantResponse := "Для темной комнаты хорошо подходят следующие растения: сансевиерия, аспидистра, спатифиллум, замиокулькас."

	// Set up the mock expectations
	mockRecommendationRepo.On("GetChatSession", mock.Anything, sessionID).Return(session, nil)
	mockRecommendationRepo.On("SaveChatMessage", mock.Anything, mock.MatchedBy(func(m *models.ChatMessage) bool {
		return m.SessionID == sessionID && m.UserID == userID && m.Role == "user" && m.Content == userMessage
	})).Return(nil)
	mockRecommendationRepo.On("GetChatMessages", mock.Anything, sessionID).Return([]*models.ChatMessage{}, nil)
	mockRecommendationRepo.On("SaveChatMessage", mock.Anything, mock.MatchedBy(func(m *models.ChatMessage) bool {
		return m.SessionID == sessionID && m.UserID == userID && m.Role == "assistant"
	})).Return(nil)
	mockRecommendationRepo.On("UpdateChatSessionLastUsed", mock.Anything, sessionID).Return(nil)

	// Create a base recommendation service
	baseService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Create a mock service that returns a fixed response
	mockService := &MockRecommendationServiceWithResponse{
		RecommendationService: baseService,
		fixedResponse:         assistantResponse,
	}

	// Initialize the in-memory session
	mockService.chatSessions[sessionID] = []Message{
		{
			Role: "system",
			Text: "Ты - эксперт по растениям. Помогай пользователям с вопросами о выращивании, уходе и выборе растений. Отвечай на русском языке.",
		},
	}

	// Test the SendChatMessage method
	result, err := mockService.SendChatMessage(context.Background(), sessionID, userID, userMessage)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the expected values
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "assistant", result.Role)
	assert.Equal(t, assistantResponse, result.Content)

	// Verify that the in-memory session was updated
	sessionMessages, ok := mockService.chatSessions[sessionID]
	assert.True(t, ok)
	assert.Equal(t, 3, len(sessionMessages)) // system + user + assistant
	assert.Equal(t, "system", sessionMessages[0].Role)
	assert.Equal(t, "user", sessionMessages[1].Role)
	assert.Equal(t, userMessage, sessionMessages[1].Text)
	assert.Equal(t, "assistant", sessionMessages[2].Role)
	assert.Equal(t, assistantResponse, sessionMessages[2].Text)

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}

// TestRecommendationService_GetChatMessages tests the GetChatMessages method
func TestRecommendationService_GetChatMessages(t *testing.T) {
	// Create mock repositories
	mockRecommendationRepo := new(MockRecommendationRepository)
	mockPlantRepo := new(MockPlantRepository)

	// Create test data
	userID := uuid.New()
	sessionID := uuid.New()

	// Expected chat session
	session := &models.ChatSession{
		ID:        sessionID,
		UserID:    userID,
		Title:     "Разговор о растениях",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Expected chat messages
	message1 := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   "Какие растения подходят для темной комнаты?",
		CreatedAt: time.Now(),
	}
	message2 := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   "Для темной комнаты хорошо подходят следующие растения: сансевиерия, аспидистра, спатифиллум, замиокулькас.",
		CreatedAt: time.Now(),
	}
	expectedMessages := []*models.ChatMessage{message1, message2}

	// Set up the mock expectations
	mockRecommendationRepo.On("GetChatSession", mock.Anything, sessionID).Return(session, nil)
	mockRecommendationRepo.On("GetChatMessages", mock.Anything, sessionID).Return(expectedMessages, nil)

	// Create the recommendation service
	recommendationService := NewRecommendationService(
		mockRecommendationRepo,
		mockPlantRepo,
		"test-api-key",
		"test-model",
	)

	// Test the GetChatMessages method
	result, err := recommendationService.GetChatMessages(context.Background(), sessionID, userID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result has the expected values
	assert.Equal(t, expectedMessages, result)

	// Verify that all expectations were met
	mockRecommendationRepo.AssertExpectations(t)
	mockPlantRepo.AssertExpectations(t)
}
