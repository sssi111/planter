package api

import (
	"encoding/json"
	"net/http"

	"github.com/anpanovv/planter/internal/middleware"
	"github.com/anpanovv/planter/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// handleGetAllPlants handles the get all plants request
func (a *API) handleGetAllPlants(w http.ResponseWriter, r *http.Request) {
	// Get all plants
	plants, err := a.plantService.GetAllPlants(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get plants")
		return
	}

	// Respond with the plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}

// handleGetPlant handles the get plant request
func (a *API) handleGetPlant(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the plant
	plant, err := a.plantService.GetPlant(r.Context(), plantID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Plant not found")
		return
	}

	// Respond with the plant
	utils.RespondWithJSON(w, http.StatusOK, plant)
}

// handleSearchPlants handles the search plants request
func (a *API) handleSearchPlants(w http.ResponseWriter, r *http.Request) {
	// Get the query parameter
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}

	// Search for plants
	plants, err := a.plantService.SearchPlants(r.Context(), query)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to search plants")
		return
	}

	// Respond with the plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}

// handleGetFavoritePlants handles the get favorite plants request
func (a *API) handleGetFavoritePlants(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the favorite plants
	plants, err := a.plantService.GetFavoritePlants(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get favorite plants")
		return
	}

	// Respond with the plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}

// handleAddToFavorites handles the add to favorites request
func (a *API) handleAddToFavorites(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Add to favorites
	err = a.plantService.AddToFavorites(r.Context(), userID, plantID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add to favorites")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Added to favorites"})
}

// handleRemoveFromFavorites handles the remove from favorites request
func (a *API) handleRemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Remove from favorites
	err = a.plantService.RemoveFromFavorites(r.Context(), userID, plantID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to remove from favorites")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Removed from favorites"})
}

// handleMarkAsWatered handles the mark as watered request
func (a *API) handleMarkAsWatered(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Mark as watered
	plant, err := a.plantService.MarkAsWatered(r.Context(), userID, plantID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark as watered")
		return
	}

	// Respond with the updated plant
	utils.RespondWithJSON(w, http.StatusOK, plant)
}

// handleGetUserPlants handles the get user plants request
func (a *API) handleGetUserPlants(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the user plants
	plants, err := a.plantService.GetUserPlants(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user plants")
		return
	}

	// Respond with the plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}

// handleAddUserPlant handles the add user plant request
func (a *API) handleAddUserPlant(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse the request body
	var req struct {
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Add the plant to the user's collection
	err = a.plantService.AddUserPlant(r.Context(), userID, plantID, req.Location)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add user plant")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Plant added to collection"})
}

// handleUpdateUserPlant handles the update user plant request
func (a *API) handleUpdateUserPlant(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse the request body
	var req struct {
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update the user plant
	err = a.plantService.UpdateUserPlant(r.Context(), userID, plantID, req.Location)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user plant")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Plant updated"})
}

// handleRemoveUserPlant handles the remove user plant request
func (a *API) handleRemoveUserPlant(w http.ResponseWriter, r *http.Request) {
	// Get the plant ID from the URL
	vars := mux.Vars(r)
	plantID, err := uuid.Parse(vars["plantId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid plant ID")
		return
	}

	// Get the authenticated user ID from the context
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Remove the user plant
	err = a.plantService.RemoveUserPlant(r.Context(), userID, plantID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to remove user plant")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Plant removed from collection"})
}