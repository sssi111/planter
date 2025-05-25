package api

import (
	"encoding/json"
	"net/http"

	"github.com/anpanovv/planter/internal/models"
	"github.com/anpanovv/planter/internal/utils"
)

// handleLogin handles the login request
func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := utils.Validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ValidationErrorMessage(err))
		return
	}

	// Login the user
	resp, err := a.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Respond with the token and user
	utils.RespondWithJSON(w, http.StatusOK, resp)
}

// handleRegister handles the registration request
func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := utils.Validate.Struct(req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.ValidationErrorMessage(err))
		return
	}

	// Register the user
	resp, err := a.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Respond with the token and user
	utils.RespondWithJSON(w, http.StatusCreated, resp)
}