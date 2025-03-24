package main

import (
	"context"
	"fmt"
	"sync"

	commons "sama/go-task-management/commons"
)

type InAppNotificationService struct {
	taskRepository              commons.TaskRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
}

func NewInAppNotificationService(
	taskRepo commons.TaskRepositoryInterface,
	eventRepo commons.TaskSystemEventRepositoryInterface,
	notifRepo commons.InAppNotificationRepositoryInterface,
) *InAppNotificationService {
	return &InAppNotificationService{
		taskRepository:              taskRepo,
		taskSystemEventRepository:   eventRepo,
		inAppNotificationRepository: notifRepo,
	}
}

func (s *InAppNotificationService) Handle(ctx context.Context, taskID, correlationID string, _ []string) error {
	task, err := s.taskRepository.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if err := s.createInAppNotification(ctx, &task); err != nil {
		return fmt.Errorf("failed to create in-app notification: %w", err)
	}

	if err := s.processNotificationEvents(ctx, taskID, correlationID); err != nil {
		return fmt.Errorf("failed to process notification events: %w", err)
	}

	return s.updateTaskStatus(ctx, &task)
}

func (s *InAppNotificationService) createInAppNotification(_ context.Context, task *commons.Task) error {
	// Create notification for the task creator
	creatorNotification := commons.InAppNotification{
		UserID:      task.CreatorID,
		Title:       task.Title,
		Description: task.Description,
	}

	_, err := s.inAppNotificationRepository.Create(creatorNotification)
	if err != nil {
		return fmt.Errorf("failed to create notification for creator: %w", err)
	}

	// Create notification for the assignee if exists
	if task.AssigneeID != nil {
		assigneeNotification := commons.InAppNotification{
			UserID:      *task.AssigneeID,
			Title:       task.Title,
			Description: task.Description,
		}

		_, err := s.inAppNotificationRepository.Create(assigneeNotification)
		if err != nil {
			return fmt.Errorf("failed to create notification for assignee: %w", err)
		}
	}

	return nil
}

func (s *InAppNotificationService) processNotificationEvents(ctx context.Context, taskID, correlationID string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)

	go s.createSystemEvent(ctx, &wg, errChan, taskID, correlationID, "Notification Service",
		"notification:db:in-app-notification-created", "In-app notification created in database", 6)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("error processing notification events: %w", err)
		}
	}

	return nil
}

func (s *InAppNotificationService) createSystemEvent(_ context.Context, wg *sync.WaitGroup, errChan chan<- error,
	taskID, correlationID, origin, action, message string, priority int) {
	defer wg.Done()

	event := commons.TaskSystemEvent{
		TaskId:        taskID,
		CorrelationId: correlationID,
		Origin:        origin,
		Action:        action,
		Message:       message,
		JsonData:      "{}",
	}

	_, err := s.taskSystemEventRepository.Create(event, priority)
	if err != nil {
		errChan <- fmt.Errorf("failed to create system event: %w", err)
	}
}

func (s *InAppNotificationService) updateTaskStatus(_ context.Context, task *commons.Task) error {
	task.InAppSent = true

	err := s.taskRepository.Update(*task)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}
