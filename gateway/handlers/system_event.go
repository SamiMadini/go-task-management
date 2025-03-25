package handlers

import (
	"log"
	"net/http"
)

type SystemEventHandler struct {
	*BaseHandler
}

func NewSystemEventHandler(base *BaseHandler) *SystemEventHandler {
	return &SystemEventHandler{BaseHandler: base}
}

type GetAllTaskSystemEventsResponse struct {
	Events []TaskSystemEventResponse `json:"events"`
}

// @Summary Get all task system events
// @Description Retrieves all system events related to tasks
// @Tags system-events
// @Accept json
// @Produce json
// @Success 200 {object} GetAllTaskSystemEventsResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-system-events [get]
func (h *SystemEventHandler) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.taskSystemEventRepository.GetAll()
	if err != nil {
		log.Printf("Failed to get all task system events: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch system events")
		return
	}

	if len(events) == 0 {
		h.respondWithJSON(w, http.StatusOK, GetAllTaskSystemEventsResponse{
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

	h.respondWithJSON(w, http.StatusOK, response)
}
