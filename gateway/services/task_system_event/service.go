package task_system_event

import (
	"encoding/json"
	"time"

	"sama/go-task-management/commons"

	"github.com/google/uuid"
)

type TaskSystemEventRepository interface {
	GetAll() ([]commons.TaskSystemEvent, error)
	Create(event commons.TaskSystemEvent, delay int) (commons.TaskSystemEvent, error)
}

type Service struct {
	logger                    commons.Logger
	taskSystemEventRepository TaskSystemEventRepository
}

func NewService(logger commons.Logger, taskSystemEventRepository TaskSystemEventRepository) *Service {
	return &Service{
		logger:                    logger,
		taskSystemEventRepository: taskSystemEventRepository,
	}
}

func (s *Service) GetAll() ([]commons.TaskSystemEvent, error) {
	events, err := s.taskSystemEventRepository.GetAll()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) Create(taskId, correlationId, origin, action, message string, data interface{}, delay int) (*commons.TaskSystemEvent, error) {
	var jsonData string
	if data != nil {
		if jsonBytes, err := json.Marshal(data); err == nil {
			jsonData = string(jsonBytes)
		}
	}

	now := time.Now()

	event := commons.TaskSystemEvent{
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

	createdEvent, err := s.taskSystemEventRepository.Create(event, delay)

	if err != nil {
		return nil, err
	}

	return &createdEvent, nil
}
