package handlers

import (
	"encoding/json"
	"net/http"

	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/handlers/constants"
	"sama/go-task-management/gateway/middleware"
	in_app_notification "sama/go-task-management/gateway/services/in_app_notification"
)

type InAppNotificationHandler struct {
	*BaseHandler
	inAppNotificationService *in_app_notification.Service
}

func NewInAppNotificationHandler(base *BaseHandler, inAppNotificationService *in_app_notification.Service) *InAppNotificationHandler {
	return &InAppNotificationHandler{
		BaseHandler:              base,
		inAppNotificationService: inAppNotificationService,
	}
}

// @Summary Create a notification
// @Description Creates a new notification for a user
// @Tags notifications
// @Accept json
// @Produce json
// @Param input body in_app_notification.CreateNotificationInput true "Notification details"
// @Success 201 {object} in_app_notification.NotificationResponse
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications [post]
func (h *InAppNotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var input in_app_notification.CreateNotificationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid request payload", err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	input.UserID = userID

	response, err := h.inAppNotificationService.CreateNotification(r.Context(), input)
	if err != nil {
		switch err {
		case commons.ErrInvalidInput:
			h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Invalid input", err.Error())
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", err.Error())
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to create notification", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusCreated, StandardResponse{
		Success: true,
		Data:    response,
	})
}

// @Summary Get user notifications
// @Description Retrieves all notifications for the authenticated user
// @Tags notifications
// @Produce json
// @Success 200 {array} in_app_notification.NotificationResponse
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications [get]
func (h *InAppNotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	notifications, err := h.inAppNotificationService.GetUserNotifications(r.Context(), userID)
	if err != nil {
		switch err {
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", err.Error())
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to get notifications", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    notifications,
	})
}

// @Summary Mark notification as read
// @Description Marks a specific notification as read
// @Tags notifications
// @Param id path string true "Notification ID"
// @Success 200
// @Failure 403 {object} ErrorResponse "Not authorized to access this notification"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications/{id}/read [post]
func (h *InAppNotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := r.PathValue("id")
	if notificationID == "" {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Notification ID is required", "")
		return
	}

	userID := middleware.GetUserIDFromContext(r)

	err := h.inAppNotificationService.MarkNotificationAsRead(r.Context(), notificationID, userID)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "Notification not found", err.Error())
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", err.Error())
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to mark notification as read", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Notification marked as read",
		},
	})
}

// @Summary Delete notification
// @Description Deletes a specific notification
// @Tags notifications
// @Param id path string true "Notification ID"
// @Success 200
// @Failure 403 {object} ErrorResponse "Not authorized to delete this notification"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /notifications/{id} [delete]
func (h *InAppNotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	notificationID := r.PathValue("id")
	if notificationID == "" {
		h.respondWithError(w, http.StatusBadRequest, constants.ErrCodeBadRequest, "Notification ID is required", "")
		return
	}

	userID := middleware.GetUserIDFromContext(r)

	err := h.inAppNotificationService.DeleteNotification(r.Context(), notificationID, userID)
	if err != nil {
		switch err {
		case commons.ErrNotFound:
			h.respondWithError(w, http.StatusNotFound, constants.ErrCodeNotFound, "Notification not found", err.Error())
		case commons.ErrUnauthorized:
			h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", err.Error())
		default:
			h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to delete notification", err.Error())
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data: map[string]string{
			"message": "Notification deleted successfully",
		},
	})
}
