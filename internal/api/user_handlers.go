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

// handleGetUser handles the get user request
func (a *API) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["userId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get the authenticated user ID from the context
	authUserID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if the user is requesting their own data
	if userID != authUserID {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	// Get the user
	user, err := a.userService.GetUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Respond with the user
	utils.RespondWithJSON(w, http.StatusOK, user)
}

// handleUpdateUser handles the update user request
func (a *API) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["userId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get the authenticated user ID from the context
	authUserID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if the user is updating their own data
	if userID != authUserID {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	// Parse the request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set the user ID
	user.ID = userID

	// Update the user
	updatedUser, err := a.userService.UpdateUser(r.Context(), &user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Respond with the updated user
	utils.RespondWithJSON(w, http.StatusOK, updatedUser)
}