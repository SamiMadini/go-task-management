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
        commons.WriteJSON(w, http.StatusOK, GetAllTaskSystemEventsResponse{
            Events: []TaskSystemEventResponse{},
        })
        return
    }

    response := GetAllTaskSystemEventsResponse{
        Events: make([]TaskSystemEventResponse, len(events)),
    }

    for i, event := range events {
        response.Events[i] = TaskSystemEventResponse{
            ID:            event.ID,
            TaskId:        event.TaskId,
            CorrelationId: event.CorrelationId,
            Origin:        event.Origin,
            Action:        event.Action,
            Message:       event.Message,
            JsonData:      event.JsonData,
            EmitAt:        event.EmitAt,
            CreatedAt:     event.CreatedAt,
        }
    }

    commons.WriteJSON(w, http.StatusOK, response)
}
