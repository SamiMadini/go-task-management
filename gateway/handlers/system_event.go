package handlers

import (
	"net/http"
	"sama/go-task-management/gateway/handlers/constants"
	"sama/go-task-management/gateway/middleware"
	task_system_event "sama/go-task-management/gateway/services/task_system_event"
)

type TaskSystemEventHandler struct {
	*BaseHandler
	taskEventService *task_system_event.Service
}

func NewTaskSystemEventHandler(base *BaseHandler, taskEventService *task_system_event.Service) *TaskSystemEventHandler {
	return &TaskSystemEventHandler{BaseHandler: base, taskEventService: taskEventService}
}

type GetAllTaskSystemEventsResponse struct {
	Events []TaskSystemEventResponse `json:"events"`
}

// @Summary Get all task system events
// @Description Retrieves all system events related to tasks
// @Tags system-events
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse{data=GetAllTaskSystemEventsResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-system-events [get]
func (h *TaskSystemEventHandler) GetAllTaskSystemEvents(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		h.respondWithError(w, http.StatusUnauthorized, constants.ErrCodeUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	events, err := h.taskEventService.GetAll()
	if err != nil {
		h.logger.Printf("Failed to get all task system events: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, constants.ErrCodeInternal, "Failed to fetch system events", err.Error())
		return
	}

	if len(events) == 0 {
		h.respondWithJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Data: GetAllTaskSystemEventsResponse{
				Events: []TaskSystemEventResponse{},
			},
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

	h.respondWithJSON(w, http.StatusOK, StandardResponse{
		Success: true,
		Data:    response,
	})
}
