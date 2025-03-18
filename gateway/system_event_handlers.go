package main

import (
    "log"
    "net/http"

    commons "sama/go-task-management/commons"
)

func (h *handler) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
    events, err := h.taskSystemEventRepository.GetAll()
    if err != nil {
        log.Printf("Failed to get all task system events: %v", err)
        commons.InternalServerErrorHandler(w)
        return
    }

    if len(events) == 0 {
        commons.WriteJSON(w, http.StatusOK, []commons.TaskSystemEvent{})
        return
    }

    commons.WriteJSON(w, http.StatusOK, events)
}
