package main

import (
	"log"
	"net/http"
	"time"

	commons "sama/go-task-management/commons"
)

// @Summary Get all in-app notifications
// @Description Retrieves all in-app notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} GetAllInAppNotificationsResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications [get]
func (h *handler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == "" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notifications, err := h.inAppNotificationRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get notifications for user %s: %v", userID, err)
		commons.InternalServerErrorHandler(w)
		return
	}

	if len(notifications) == 0 {
		commons.WriteJSON(w, http.StatusOK, GetAllInAppNotificationsResponse{
			InAppNotifications: []InAppNotificationResponse{},
		})
		return
	}

	response := GetAllInAppNotificationsResponse{
		InAppNotifications: make([]InAppNotificationResponse, len(notifications)),
	}

	for i, n := range notifications {
		var readAt time.Time
		if n.ReadAt != nil {
			readAt = *n.ReadAt
		}
		response.InAppNotifications[i] = InAppNotificationResponse{
			ID:          n.ID,
			Title:       n.Title,
			Description: n.Description,
			IsRead:      n.IsRead,
			ReadAt:      readAt,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	commons.WriteJSON(w, http.StatusOK, response)
}

// @Summary Mark notification as read
// @Description Updates a notification's read status
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} UpdateOnReadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id}/read [put]
func (h *handler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == "" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		commons.WriteJSONError(w, http.StatusBadRequest, "invalid notification ID format")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching notification %s: %v", id, err)
		commons.WriteJSONError(w, http.StatusNotFound, "notification not found")
		return
	}

	// Check if user has permission to update this notification
	if notification.UserID != userID {
		commons.WriteJSONError(w, http.StatusForbidden, "You don't have permission to update this notification")
		return
	}

	var params UpdateOnReadRequest
	if err := commons.ReadJSON(r, &params); err != nil {
		log.Printf("Invalid read update request: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		commons.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.inAppNotificationRepository.UpdateOnRead(id, params.IsRead); err != nil {
		log.Printf("Error updating notification read status: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	commons.WriteJSON(w, http.StatusOK, UpdateOnReadResponse{
		Success: true,
	})
}

// @Summary Delete a notification
// @Description Deletes a notification from the system
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id} [delete]
func (h *handler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r)
	if userID == "" {
		commons.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		commons.WriteJSONError(w, http.StatusBadRequest, "invalid notification ID format")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching notification %s: %v", id, err)
		commons.WriteJSONError(w, http.StatusNotFound, "notification not found")
		return
	}

	// Check if user has permission to delete this notification
	if notification.UserID != userID {
		commons.WriteJSONError(w, http.StatusForbidden, "You don't have permission to delete this notification")
		return
	}

	if err := h.inAppNotificationRepository.Delete(id); err != nil {
		log.Printf("Error deleting notification: %v", err)
		commons.InternalServerErrorHandler(w)
		return
	}

	commons.WriteJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
