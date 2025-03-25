package handlers

import (
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

func (r *UpdateOnReadRequest) Validate() []ValidationError {
	return nil
}

type UpdateOnReadResponse struct {
	Success bool `json:"success"`
}

// @Summary Get all in-app notifications
// @Description Retrieves all in-app notifications for the current user
// @Description
// @Description Error scenarios:
// @Description - Unauthorized access: Returns 401 with UNAUTHORIZED code
// @Description - Database error: Returns 500 with INTERNAL_ERROR code
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} GetAllInAppNotificationsResponse "List of notifications"
// @Failure 401 {object} ErrorResponse "Authentication required"
// @Failure 500 {object} ErrorResponse "Failed to fetch notifications"
// @Router /notifications [get]
func (h *NotificationHandler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	notifications, err := h.inAppNotificationRepository.GetByUserID(userID)
	if err != nil {
		h.logger.Printf("Failed to get notifications for user %s: %v", userID, err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to fetch notifications", err.Error())
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
// @Description
// @Description Error scenarios:
// @Description - Invalid UUID format: Returns 400 with BAD_REQUEST code
// @Description - Notification not found: Returns 404 with NOT_FOUND code
// @Description - Unauthorized access: Returns 403 with FORBIDDEN code (not the recipient)
// @Description - Invalid request body: Returns 400 with VALIDATION_ERROR code
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID (UUID format)"
// @Param request body UpdateOnReadRequest true "Read status update"
// @Success 200 {object} UpdateOnReadResponse "Status updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid notification ID or request format"
// @Failure 401 {object} ErrorResponse "Authentication required"
// @Failure 403 {object} ErrorResponse "Not authorized to update this notification"
// @Failure 404 {object} ErrorResponse "Notification not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Notification ID is required", "Path parameter 'id' is missing")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid notification ID format", "Notification ID must be a valid UUID")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		h.logger.Printf("Error fetching notification %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, ErrCodeNotFound, "Notification not found", err.Error())
		return
	}

	if notification.UserID != userID {
		h.respondWithError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied", "You don't have permission to update this notification")
		return
	}

	var req UpdateOnReadRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid request body", err.Error())
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		h.respondWithValidationErrors(w, validationErrors)
		return
	}

	if err := h.inAppNotificationRepository.UpdateOnRead(id, req.IsRead); err != nil {
		h.logger.Printf("Error updating notification read status: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to update notification", err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, UpdateOnReadResponse{
		Success: true,
	})
}

// @Summary Delete a notification
// @Description Deletes a notification from the system
// @Description
// @Description Error scenarios:
// @Description - Invalid UUID format: Returns 400 with BAD_REQUEST code
// @Description - Notification not found: Returns 404 with NOT_FOUND code
// @Description - Unauthorized deletion: Returns 403 with FORBIDDEN code (not the recipient)
// @Description - Database error: Returns 500 with INTERNAL_ERROR code
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID (UUID format)"
// @Success 200 {object} map[string]bool "Deletion status"
// @Failure 400 {object} ErrorResponse "Invalid notification ID format"
// @Failure 401 {object} ErrorResponse "Authentication required"
// @Failure 403 {object} ErrorResponse "Not authorized to delete this notification"
// @Failure 404 {object} ErrorResponse "Notification not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Notification ID is required", "Path parameter 'id' is missing")
		return
	}

	if !commons.IsValidUUID(id) {
		h.respondWithError(w, http.StatusBadRequest, ErrCodeBadRequest, "Invalid notification ID format", "Notification ID must be a valid UUID")
		return
	}

	notification, err := h.inAppNotificationRepository.GetByID(id)
	if err != nil {
		h.logger.Printf("Error fetching notification %s: %v", id, err)
		h.respondWithError(w, http.StatusNotFound, ErrCodeNotFound, "Notification not found", err.Error())
		return
	}

	if notification.UserID != userID {
		h.respondWithError(w, http.StatusForbidden, ErrCodeForbidden, "Access denied", "You don't have permission to delete this notification")
		return
	}

	if err := h.inAppNotificationRepository.Delete(id); err != nil {
		h.logger.Printf("Error deleting notification: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, ErrCodeInternal, "Failed to delete notification", err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]bool{
		"success": true,
	})
}
