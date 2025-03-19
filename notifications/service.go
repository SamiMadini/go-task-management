package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	commons "sama/go-task-management/commons"
)

const (
	defaultRegion    = "us-east-1"
	defaultQueueName = "go-email-service-queue"
)

type Config struct {
	AWSEndpoint string
	AWSRegion   string
	QueueName   string
}

type EmailEvent struct {
	TaskID        string `json:"taskId"`
	CorrelationID string `json:"correlationId"`
}
type NotificationService struct {
	taskRepository              commons.TaskRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
	sqsClient                   SQSClientInterface
}

func NewNotificationService(
	taskRepo commons.TaskRepositoryInterface,
	eventRepo commons.TaskSystemEventRepositoryInterface,
	notifRepo commons.InAppNotificationRepositoryInterface,
	sqsClient SQSClientInterface,
) *NotificationService {
	return &NotificationService{
		taskRepository:              taskRepo,
		taskSystemEventRepository:   eventRepo,
		inAppNotificationRepository: notifRepo,
		sqsClient:                   sqsClient,
	}
}

func (s *NotificationService) Handle(ctx context.Context, taskID, correlationID string, notificationTypes []string) error {
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

func (s *NotificationService) createInAppNotification(_ context.Context, task *commons.Task) error {
	inAppNotification := commons.InAppNotification{
		Title:       task.Title,
		Description: task.Description,
	}

	_, err := s.inAppNotificationRepository.Create(inAppNotification)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) processNotificationEvents(ctx context.Context, taskID, correlationID string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(3)

	go s.createSystemEvent(ctx, &wg, errChan, taskID, correlationID, "Notification Service", 
		"notification:db:in-app-notification-created", "In-app notification created in database", 6)

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

func (s *NotificationService) createSystemEvent(_ context.Context, wg *sync.WaitGroup, errChan chan<- error,
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

func (s *NotificationService) sendEmailNotification(ctx context.Context, taskID, correlationID string) error {
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

func (s *NotificationService) updateTaskStatus(_ context.Context, task *commons.Task) error {
	task.EmailSent = true
	task.InAppSent = true

	err := s.taskRepository.Update(*task)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}
