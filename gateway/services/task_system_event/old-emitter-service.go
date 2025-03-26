package task_system_event

import (
	"encoding/json"
	"time"

	"sama/go-task-management/commons"

	"github.com/google/uuid"
)

const (
	EventOriginGateway = "gateway"
)

type EventAction string

const (
	EventActionCreated  EventAction = "CREATED"
	EventActionUpdated  EventAction = "UPDATED"
	EventActionDeleted  EventAction = "DELETED"
	EventActionAssigned EventAction = "ASSIGNED"
)

type EventEmitter interface {
	EmitEvent(event commons.TaskSystemEvent) error
}

type eventService struct {
	emitter EventEmitter
	logger  commons.Logger
}

func NewEmitterService(emitter EventEmitter, logger commons.Logger) *eventService {
	return &eventService{
		emitter: emitter,
		logger:  logger,
	}
}

func (s *eventService) EmitTaskCreated(task commons.Task, correlationID string) error {
	event := buildTaskSystemEvent(
		task.ID,
		correlationID,
		EventOriginGateway,
		string(EventActionCreated),
		"Task created",
		task,
	)
	return s.emitter.EmitEvent(event)
}

func (s *eventService) EmitTaskUpdated(task commons.Task, correlationID string, changes map[string]interface{}) error {
	event := buildTaskSystemEvent(
		task.ID,
		correlationID,
		EventOriginGateway,
		string(EventActionUpdated),
		"Task updated",
		changes,
	)
	return s.emitter.EmitEvent(event)
}

func (s *eventService) EmitTaskDeleted(taskID string, correlationID string) error {
	event := buildTaskSystemEvent(
		taskID,
		correlationID,
		EventOriginGateway,
		string(EventActionDeleted),
		"Task deleted",
		nil,
	)
	return s.emitter.EmitEvent(event)
}

func (s *eventService) EmitTaskAssigned(task commons.Task, correlationID string, assigneeID *string) error {
	data := map[string]interface{}{
		"assignee_id": assigneeID,
	}
	event := buildTaskSystemEvent(
		task.ID,
		correlationID,
		EventOriginGateway,
		string(EventActionAssigned),
		"Task assigned",
		data,
	)
	return s.emitter.EmitEvent(event)
}

func buildTaskSystemEvent(taskId, correlationId, origin, action, message string, data interface{}) commons.TaskSystemEvent {
	var jsonData string
	if data != nil {
		if jsonBytes, err := json.Marshal(data); err == nil {
			jsonData = string(jsonBytes)
		}
	}

	now := time.Now()
	return commons.TaskSystemEvent{
		ID:            uuid.New().String(),
		TaskId:        taskId,
		CorrelationId: correlationId,
		Origin:        origin,
		Action:        action,
		Message:       message,
		JsonData:      jsonData,
		EmitAt:        now,
		CreatedAt:     now,
	}
}
