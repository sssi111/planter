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

	// Respond with the questionnaire
	utils.RespondWithJSON(w, http.StatusCreated, questionnaire)
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