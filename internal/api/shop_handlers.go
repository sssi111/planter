package api

import (
	"net/http"

	"github.com/anpanovv/planter/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// handleGetAllShops handles the get all shops request
func (a *API) handleGetAllShops(w http.ResponseWriter, r *http.Request) {
	// Get all shops
	shops, err := a.shopService.GetAllShops(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get shops")
		return
	}

	// Respond with the shops
	utils.RespondWithJSON(w, http.StatusOK, shops)
}

// handleGetShop handles the get shop request
func (a *API) handleGetShop(w http.ResponseWriter, r *http.Request) {
	// Get the shop ID from the URL
	vars := mux.Vars(r)
	shopID, err := uuid.Parse(vars["shopId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	// Get the shop
	shop, err := a.shopService.GetShop(r.Context(), shopID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Shop not found")
		return
	}

	// Respond with the shop
	utils.RespondWithJSON(w, http.StatusOK, shop)
}

// handleGetShopPlants handles the get shop plants request
func (a *API) handleGetShopPlants(w http.ResponseWriter, r *http.Request) {
	// Get the shop ID from the URL
	vars := mux.Vars(r)
	shopID, err := uuid.Parse(vars["shopId"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	// Get the shop plants
	plants, err := a.shopService.GetShopPlants(r.Context(), shopID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get shop plants")
		return
	}

	// Respond with the plants
	utils.RespondWithJSON(w, http.StatusOK, plants)
}