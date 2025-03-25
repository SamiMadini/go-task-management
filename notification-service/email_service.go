package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	commons "sama/go-task-management/commons"
)

type EmailEvent struct {
	TaskID        string `json:"taskId"`
	CorrelationID string `json:"correlationId"`
}

type EmailNotificationService struct {
	taskRepository            commons.TaskRepositoryInterface
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface
	sqsClient                SQSClientInterface
}

func NewEmailNotificationService(
	taskRepo commons.TaskRepositoryInterface,
	eventRepo commons.TaskSystemEventRepositoryInterface,
	sqsClient SQSClientInterface,
) *EmailNotificationService {
	return &EmailNotificationService{
		taskRepository:            taskRepo,
		taskSystemEventRepository: eventRepo,
		sqsClient:                sqsClient,
	}
}

func (s *EmailNotificationService) Handle(ctx context.Context, taskID, correlationID string, _ []string) error {
	task, err := s.taskRepository.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if err := s.processNotificationEvents(ctx, taskID, correlationID); err != nil {
		return fmt.Errorf("failed to process notification events: %w", err)
	}

	return s.updateTaskStatus(ctx, &task)
}

func (s *EmailNotificationService) processNotificationEvents(ctx context.Context, taskID, correlationID string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go s.createSystemEvent(ctx, &wg, errChan, taskID, correlationID, "Notification Service",
		"notification:event:email-task-created", "Email event sent", 9)

	go func() {
		defer wg.Done()
		if err := s.sendEmailNotification(ctx, taskID, correlationID); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("error processing notification events: %w", err)
		}
	}

	return nil
}

func (s *EmailNotificationService) createSystemEvent(_ context.Context, wg *sync.WaitGroup, errChan chan<- error,
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

func (s *EmailNotificationService) sendEmailNotification(ctx context.Context, taskID, correlationID string) error {
	emailEvent := EmailEvent{
		TaskID:        taskID,
		CorrelationID: correlationID,
	}

	jsonBytes, err := json.Marshal(emailEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal email event: %w", err)
	}

	if err := s.sqsClient.SendMessage(ctx, string(jsonBytes)); err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	log.Printf("Successfully sent task notification to SQS queue for email processing: %s", taskID)
	return nil
}

func (s *EmailNotificationService) updateTaskStatus(_ context.Context, task *commons.Task) error {
	task.EmailSent = true

	err := s.taskRepository.Update(*task)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}
