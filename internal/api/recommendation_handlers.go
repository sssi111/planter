package api

import (
	"encoding/json"
	"net/http"

	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// handleSaveQuestionnaire handles the save questionnaire request
func (a *API) handleSaveQuestionnaire(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.QuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := utils.Validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ValidationErrorMessage(err))
		return
	}

	// Get the authenticated user ID from the context if available
	var userID *uuid.UUID
	authUserID, err := middleware.GetUserID(r.Context())
	if err == nil {
		userID = &authUserID
	}

	// Save the questionnaire
	questionnaire, err := a.recommendationService.SaveQuestionnaire(r.Context(), userID, &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save questionnaire")
		return
	}

	// Get recommendations
	plants, err := a.recommendationService.GetRecommendations(r.Context(), questionnaire.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	if len(plants) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No plants found matching the criteria")
		return
	}

	// Respond with the best matching plant (first in the list)
	utils.RespondWithJSON(w, http.StatusCreated, plants[0])
}

// handleGetRecommendations handles the get recommendations request
func (a *API) handleGetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get the questionnaire ID from the URL
	vars := mux.Vars(r)
	questionnaireID, err := uuid.Parse(vars["questionnaireId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	// Get the recommendations
	plants, err := a.recommendationService.GetRecommendations(r.Context(), questionnaireID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	// Respond with the recommended plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}

// handleSaveDetailedQuestionnaire handles the save detailed questionnaire request
func (a *API) handleSaveDetailedQuestionnaire(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.DetailedQuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := utils.Validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ValidationErrorMessage(err))
		return
	}

	// Get the authenticated user ID from the context if available
	var userID *uuid.UUID
	authUserID, err := middleware.GetUserID(r.Context())
	if err == nil {
		userID = &authUserID
	}

	// Save the detailed questionnaire
	questionnaire, err := a.recommendationService.SaveDetailedQuestionnaire(r.Context(), userID, &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save questionnaire")
		return
	}

	// Get recommendations
	plants, err := a.recommendationService.GetRecommendations(r.Context(), questionnaire.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	if len(plants) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No plants found matching the criteria")
		return
	}

	// Respond with the best matching plant (first in the list)
	utils.RespondWithJSON(w, http.StatusCreated, plants[0])
}

// handleCreateChatSession handles the create chat session request
func (a *API) handleCreateChatSession(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Create a new chat session
	session, err := a.recommendationService.CreateChatSession(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create chat session")
		return
	}

	// Respond with the chat session
	utils.RespondWithJSON(w, http.StatusCreated, session)
}

// handleGetChatSessions handles the get chat sessions request
func (a *API) handleGetChatSessions(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get all chat sessions for the user
	sessions, err := a.recommendationService.GetChatSessionsByUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get chat sessions")
		return
	}

	// Respond with the chat sessions
	utils.RespondWithJSON(w, http.StatusOK, sessions)
}

// handleGetChatSession handles the get chat session request
func (a *API) handleGetChatSession(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the chat session ID from the URL
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	// Get the chat session
	session, err := a.recommendationService.GetChatSession(r.Context(), sessionID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get chat session")
		return
	}

	// Check if the user owns the session
	if session.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	// Respond with the chat session
	utils.RespondWithJSON(w, http.StatusOK, session)
}

// handleSendChatMessage handles the send chat message request
func (a *API) handleSendChatMessage(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the chat session ID from the URL
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	// Parse the request body
	var req models.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := utils.Validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ValidationErrorMessage(err))
		return
	}

	// Send the chat message
	message, err := a.recommendationService.SendChatMessage(r.Context(), sessionID, userID, req.Message)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send chat message")
		return
	}

	// Respond with the chat message
	utils.RespondWithJSON(w, http.StatusOK, models.ChatResponse{Message: *message})
}

// handleGetChatMessages handles the get chat messages request
func (a *API) handleGetChatMessages(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the chat session ID from the URL
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	// Get the chat messages
	messages, err := a.recommendationService.GetChatMessages(r.Context(), sessionID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get chat messages")
		return
	}

	// Respond with the chat messages
	utils.RespondWithJSON(w, http.StatusOK, messages)
}