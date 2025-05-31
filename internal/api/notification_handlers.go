package api

import (
    "net/http"
    "strconv"

    "github.com/anpanovv/planter/internal/middleware"
    "github.com/anpanovv/planter/internal/services"
    "github.com/anpanovv/planter/internal/utils"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

// handleGetUserNotifications handles the get user notifications request
func (a *API) handleGetUserNotifications(w http.ResponseWriter, r *http.Request) {
    // Get the authenticated user ID from the context
    userID, err := middleware.GetUserID(r.Context())
    if err != nil {
        utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Get pagination parameters
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

    // Get notifications
    response, err := a.notificationService.GetUserNotifications(r.Context(), userID, page, pageSize)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get notifications")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, response)
}

// handleMarkNotificationAsRead handles the mark notification as read request
func (a *API) handleMarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
    // Get the authenticated user ID from the context
    userID, err := middleware.GetUserID(r.Context())
    if err != nil {
        utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Get the notification ID from the URL
    vars := mux.Vars(r)
    notificationID, err := uuid.Parse(vars["notificationId"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid notification ID")
        return
    }

    // Mark as read
    err = a.notificationService.MarkAsRead(r.Context(), notificationID, userID)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark notification as read")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
} 