package main

import (
	"log"
	"net/http"
	"time"

	commons "sama/go-task-management/commons"
)

func (h *handler) GetAllInAppNotifications(w http.ResponseWriter, r *http.Request) {
    notifications, err := h.inAppNotificationRepository.GetAll()
    if err != nil {
        log.Printf("Failed to get all notifications: %v", err)
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

func (h *handler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
        return
    }

    if !commons.IsValidUUID(id) {
        commons.WriteJSONError(w, http.StatusBadRequest, "invalid notification ID format")
        return
    }

    _, err := h.inAppNotificationRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching notification %s: %v", id, err)
        commons.WriteJSONError(w, http.StatusNotFound, "notification not found")
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

func (h *handler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
        return
    }

    if !commons.IsValidUUID(id) {
        commons.WriteJSONError(w, http.StatusBadRequest, "invalid notification ID format")
        return
    }

    _, err := h.inAppNotificationRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching notification %s: %v", id, err)
        commons.WriteJSONError(w, http.StatusNotFound, "notification not found")
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
