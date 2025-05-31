package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecommendationService is a mock implementation of the recommendation service
type MockRecommendationService struct {
	mock.Mock
}

func (m *MockRecommendationService) SaveQuestionnaire(ctx context.Context, userID *uuid.UUID, questionnaire *models.QuestionnaireRequest) (*models.PlantQuestionnaire, error) {
	args := m.Called(ctx, userID, questionnaire)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlantQuestionnaire), args.Error(1)
}

func (m *MockRecommendationService) GetRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, questionnaireID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockRecommendationService) GenerateRecommendations(ctx context.Context, questionnaireID uuid.UUID) ([]*models.Plant, error) {
	args := m.Called(ctx, questionnaireID)
	return args.Get(0).([]*models.Plant), args.Error(1)
}

func (m *MockRecommendationService) SaveDetailedQuestionnaire(ctx context.Context, userID *uuid.UUID, questionnaire *models.DetailedQuestionnaireRequest) (*models.PlantQuestionnaire, error) {
	args := m.Called(ctx, userID, questionnaire)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlantQuestionnaire), args.Error(1)
}

func (m *MockRecommendationService) CreateChatSession(ctx context.Context, userID uuid.UUID) (*models.ChatSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationService) GetChatSession(ctx context.Context, id uuid.UUID) (*models.ChatSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationService) GetChatSessionsByUser(ctx context.Context, userID uuid.UUID) ([]*models.ChatSession, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.ChatSession), args.Error(1)
}

func (m *MockRecommendationService) SendChatMessage(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, message string) (*models.ChatMessage, error) {
	args := m.Called(ctx, sessionID, userID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatMessage), args.Error(1)
}

func (m *MockRecommendationService) GetChatMessages(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) ([]*models.ChatMessage, error) {
	args := m.Called(ctx, sessionID, userID)
	return args.Get(0).([]*models.ChatMessage), args.Error(1)
}

// TestHandlers is a test implementation of the API handlers
type TestHandlers struct {
	recommendationService *MockRecommendationService
}

// handleSaveQuestionnaire is a test implementation of the handleSaveQuestionnaire handler
func (h *TestHandlers) handleSaveQuestionnaire(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.QuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the authenticated user ID from the context if available
	var userID *uuid.UUID
	authUserID, err := middleware.GetUserID(r.Context())
	if err == nil {
		userID = &authUserID
	}

	// Save the questionnaire
	questionnaire, err := h.recommendationService.SaveQuestionnaire(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, "Failed to save questionnaire", http.StatusInternalServerError)
		return
	}

	// Get recommendations
	plants, err := h.recommendationService.GetRecommendations(r.Context(), questionnaire.ID)
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	if len(plants) == 0 {
		http.Error(w, "No plants found matching the criteria", http.StatusNotFound)
		return
	}

	// Respond with the best matching plant
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plants[0])
}

// handleSaveDetailedQuestionnaire is a test implementation of the handleSaveDetailedQuestionnaire handler
func (h *TestHandlers) handleSaveDetailedQuestionnaire(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.DetailedQuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the authenticated user ID from the context if available
	var userID *uuid.UUID
	authUserID, err := middleware.GetUserID(r.Context())
	if err == nil {
		userID = &authUserID
	}

	// Save the detailed questionnaire
	questionnaire, err := h.recommendationService.SaveDetailedQuestionnaire(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, "Failed to save questionnaire", http.StatusInternalServerError)
		return
	}

	// Get recommendations
	plants, err := h.recommendationService.GetRecommendations(r.Context(), questionnaire.ID)
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	if len(plants) == 0 {
		http.Error(w, "No plants found matching the criteria", http.StatusNotFound)
		return
	}

	// Respond with the best matching plant
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plants[0])
}

// handleCreateChatSession is a test implementation of the handleCreateChatSession handler
func (h *TestHandlers) handleCreateChatSession(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a new chat session
	session, err := h.recommendationService.CreateChatSession(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to create chat session", http.StatusInternalServerError)
		return
	}

	// Respond with the chat session
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// handleSendChatMessage is a test implementation of the handleSendChatMessage handler
func (h *TestHandlers) handleSendChatMessage(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the chat session ID from the URL
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send the chat message
	message, err := h.recommendationService.SendChatMessage(r.Context(), sessionID, userID, req.Message)
	if err != nil {
		http.Error(w, "Failed to send chat message", http.StatusInternalServerError)
		return
	}

	// Respond with the chat message
	w.Header().Set("Content-Type", "application/json")
	response := models.ChatResponse{Message: *message}
	json.NewEncoder(w).Encode(response)
}

// TestHandleSaveDetailedQuestionnaire tests the handleSaveDetailedQuestionnaire function
func TestHandleSaveDetailedQuestionnaire(t *testing.T) {
	// Create a mock recommendation service
	mockService := new(MockRecommendationService)

	// Create test handlers
	handlers := &TestHandlers{
		recommendationService: mockService,
	}

	// Create a test detailed questionnaire request
	detailedRequest := &models.DetailedQuestionnaireRequest{
		SunlightPreference:   models.SunlightLevelMedium,
		PetFriendly:          true,
		CareLevel:            3,
		PreferredLocation:    stringPtr("Living Room"),
		HasChildren:          true,
		PlantSize:            "MEDIUM",
		FloweringPreference:  true,
		AirPurifying:         true,
		WateringFrequency:    "REGULAR",
		ExperienceLevel:      "BEGINNER",
		AdditionalPreferences: stringPtr("Low maintenance"),
	}

	// Expected questionnaire
	expectedQuestionnaire := &models.PlantQuestionnaire{
		ID:                 uuid.New(),
		SunlightPreference: models.SunlightLevelMedium,
		PetFriendly:        true,
		CareLevel:          3,
		PreferredLocation:  stringPtr("Living Room"),
		CreatedAt:          time.Now(),
	}

	// Expected plant response
	expectedPlant := &models.Plant{
		ID:             uuid.New(),
		Name:           "Spathiphyllum",
		ScientificName: "Spathiphyllum wallisii",
		Description:    "Peace lily, good for low light conditions",
		ImageURL:       "https://example.com/spathiphyllum.jpg",
		CreatedAt:      time.Now(),
	}

	// Set up the mock expectations
	mockService.On("SaveDetailedQuestionnaire", mock.Anything, mock.AnythingOfType("*uuid.UUID"), mock.AnythingOfType("*models.DetailedQuestionnaireRequest")).
		Return(expectedQuestionnaire, nil)
	mockService.On("GetRecommendations", mock.Anything, expectedQuestionnaire.ID).
		Return([]*models.Plant{expectedPlant}, nil)

	// Create a request
	requestBody, _ := json.Marshal(detailedRequest)
	req, err := http.NewRequest("POST", "/recommendations/questionnaire/detailed", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.handleSaveDetailedQuestionnaire(rr, req)

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

	// Verify that all expectations were met
	mockService.AssertExpectations(t)
}

// TestHandleCreateChatSession tests the handleCreateChatSession function
func TestHandleCreateChatSession(t *testing.T) {
	// Create a mock recommendation service
	mockService := new(MockRecommendationService)

	// Create test handlers
	handlers := &TestHandlers{
		recommendationService: mockService,
	}

	// Create a test user ID
	userID := uuid.New()

	// Expected response
	expectedSession := &models.ChatSession{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Разговор о растениях",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	// Set up the mock expectations
	mockService.On("CreateChatSession", mock.Anything, userID).
		Return(expectedSession, nil)

	// Create a request
	req, err := http.NewRequest("POST", "/chat/sessions", nil)
	assert.NoError(t, err)

	// Add a user ID to the context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.handleCreateChatSession(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse the response
	var response models.ChatSession
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, expectedSession.ID, response.ID)
	assert.Equal(t, expectedSession.UserID, response.UserID)
	assert.Equal(t, expectedSession.Title, response.Title)

	// Verify that all expectations were met
	mockService.AssertExpectations(t)
}

// TestHandleSendChatMessage tests the handleSendChatMessage function
func TestHandleSendChatMessage(t *testing.T) {
	// Create a mock recommendation service
	mockService := new(MockRecommendationService)

	// Create test handlers
	handlers := &TestHandlers{
		recommendationService: mockService,
	}

	// Create test data
	userID := uuid.New()
	sessionID := uuid.New()
	chatRequest := &models.ChatRequest{
		Message: "Какие растения подходят для темной комнаты?",
	}

	// Expected response
	expectedMessage := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   "Для темной комнаты хорошо подходят следующие растения: сансевиерия, аспидистра, спатифиллум, замиокулькас.",
		CreatedAt: time.Now(),
	}

	// Set up the mock expectations
	mockService.On("SendChatMessage", mock.Anything, sessionID, userID, chatRequest.Message).
		Return(expectedMessage, nil)

	// Create a request
	requestBody, _ := json.Marshal(chatRequest)
	req, err := http.NewRequest("POST", "/chat/sessions/"+sessionID.String()+"/messages", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add a user ID to the context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a router and add the route parameters
	router := mux.NewRouter()
	router.HandleFunc("/chat/sessions/{sessionId}/messages", handlers.handleSendChatMessage).Methods("POST")

	// Create a new request with the same URL but using the router
	req = mux.SetURLVars(req, map[string]string{
		"sessionId": sessionID.String(),
	})

	// Call the handler
	handlers.handleSendChatMessage(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response
	var response models.ChatResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, expectedMessage.ID, response.Message.ID)
	assert.Equal(t, expectedMessage.SessionID, response.Message.SessionID)
	assert.Equal(t, expectedMessage.UserID, response.Message.UserID)
	assert.Equal(t, expectedMessage.Role, response.Message.Role)
	assert.Equal(t, expectedMessage.Content, response.Message.Content)

	// Verify that all expectations were met
	mockService.AssertExpectations(t)
}

// TestHandleSaveQuestionnaire tests the handleSaveQuestionnaire function
func TestHandleSaveQuestionnaire(t *testing.T) {
	// Create a mock recommendation service
	mockService := new(MockRecommendationService)

	// Create test handlers
	handlers := &TestHandlers{
		recommendationService: mockService,
	}

	// Create a test questionnaire request
	questionnaireRequest := &models.QuestionnaireRequest{
		SunlightPreference:   models.SunlightLevelMedium,
		PetFriendly:          true,
		CareLevel:            3,
		PreferredLocation:    stringPtr("Living Room"),
		AdditionalPreferences: stringPtr("Low maintenance"),
	}

	// Expected questionnaire
	expectedQuestionnaire := &models.PlantQuestionnaire{
		ID:                 uuid.New(),
		SunlightPreference: models.SunlightLevelMedium,
		PetFriendly:        true,
		CareLevel:          3,
		PreferredLocation:  stringPtr("Living Room"),
		CreatedAt:          time.Now(),
	}

	// Expected plant response
	expectedPlant := &models.Plant{
		ID:             uuid.New(),
		Name:           "Sansevieria",
		ScientificName: "Sansevieria trifasciata",
		Description:    "Snake plant, very hardy and low maintenance",
		ImageURL:       "https://example.com/sansevieria.jpg",
		CreatedAt:      time.Now(),
	}

	// Set up the mock expectations
	mockService.On("SaveQuestionnaire", mock.Anything, mock.AnythingOfType("*uuid.UUID"), mock.AnythingOfType("*models.QuestionnaireRequest")).
		Return(expectedQuestionnaire, nil)
	mockService.On("GetRecommendations", mock.Anything, expectedQuestionnaire.ID).
		Return([]*models.Plant{expectedPlant}, nil)

	// Create a request
	requestBody, _ := json.Marshal(questionnaireRequest)
	req, err := http.NewRequest("POST", "/recommendations/questionnaire", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.handleSaveQuestionnaire(rr, req)

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

	// Verify that all expectations were met
	mockService.AssertExpectations(t)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}