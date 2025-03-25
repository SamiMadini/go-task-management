package handlers

import (
	"log"
	"net/http"
	"time"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/middleware"
)

type NotificationHandler struct {
	*BaseHandler
}

func NewNotificationHandler(base *BaseHandler) *NotificationHandler {
	return &NotificationHandler{BaseHandler: base}
}

type InAppNotificationResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsRead      bool      `json:"is_read"`
	ReadAt      time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAllInAppNotificationsResponse struct {
	InAppNotifications []InAppNotificationResponse `json:"in_app_notifications"`
}

type UpdateOnReadRequest struct {
	IsRead bool `json:"is_read"`
}

func (r *UpdateOnReadRequest) Validate() error {
	// TODO: Add validation
	return nil
}

type UpdateOnReadResponse struct {
	Success bool `json:"success"`
}

// @Summary Get all in-app notifications
// @Description Retrieves all in-app notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} GetAllInAppNotificationsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notifications, err := h.inAppNotificationRepository.GetByUserID(userID)
	if err != nil {
		log.Printf("Failed to get notifications for user %s: %v", userID, err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	if len(notifications) == 0 {
		h.respondWithJSON(w, http.StatusOK, GetAllInAppNotificationsResponse{
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

	h.respondWithJSON(w, http.StatusOK, response)
}

// @Summary Mark notification as read
// @Description Updates a notification's read status
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Param request body UpdateOnReadRequest true "Read status update"
// @Success 200 {object} UpdateOnReadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "notification ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, "invalid notification ID format")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching notification %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, "notification not found")
		return
	}

	if notification.UserID != userID {
		h.respondWithError(w, http.StatusForbidden, "You don't have permission to update this notification")
		return
	}

	var params UpdateOnReadRequest
	if err := h.decodeJSON(r, &params); err != nil {
		log.Printf("Invalid read update request: %v", err)
		h.respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := params.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.inAppNotificationRepository.UpdateOnRead(id, params.IsRead); err != nil {
		log.Printf("Error updating notification read status: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update notification")
		return
	}

	h.respondWithJSON(w, http.StatusOK, UpdateOnReadResponse{
		Success: true,
	})
}

// @Summary Delete a notification
// @Description Deletes a notification from the system
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, "notification ID is required")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, "invalid notification ID format")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		log.Printf("Error fetching notification %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, "notification not found")
		return
	}

	if notification.UserID != userID {
		h.respondWithError(w, http.StatusForbidden, "You don't have permission to delete this notification")
		return
	}

	if err := h.inAppNotificationRepository.Delete(id); err != nil {
		log.Printf("Error deleting notification: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
