package main

import (
    "log"
    "net/http"

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
        commons.WriteJSON(w, http.StatusOK, []commons.InAppNotification{})
        return
    }

    commons.WriteJSON(w, http.StatusOK, notifications)
}

func (h *handler) UpdateOnRead(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
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

    if err := h.inAppNotificationRepository.UpdateOnRead(id, params.IsRead); err != nil {
        log.Printf("Error updating notification %s read status: %v", id, err)
        commons.InternalServerErrorHandler(w)
        return
    }

    commons.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (h *handler) DeleteInAppNotification(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        commons.WriteJSONError(w, http.StatusBadRequest, "notification ID is required")
        return
    }

    _, err := h.inAppNotificationRepository.GetByID(id)
    if err != nil {
        log.Printf("Error fetching notification %s: %v", id, err)
        commons.WriteJSONError(w, http.StatusNotFound, "notification not found")
        return
    }

    if err := h.inAppNotificationRepository.Delete(id); err != nil {
        log.Printf("Error deleting notification %s: %v", id, err)
        commons.InternalServerErrorHandler(w)
        return
    }

    commons.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
