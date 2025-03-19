package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	commons "sama/go-task-management/commons"
)

type EmailNotificationEvent struct {
	TaskId        string `json:"taskId"`
	CorrelationId string `json:"correlationId"`
}

type MessageHandler struct {
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface
}

func NewMessageHandler(eventRepo commons.TaskSystemEventRepositoryInterface) *MessageHandler {
	return &MessageHandler{
		taskSystemEventRepository: eventRepo,
	}
}

func (h *MessageHandler) HandleMessage(ctx context.Context, message []byte) error {
	var event EmailNotificationEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf(
		"Processing message for task: %s, correlation: %s",
		event.TaskId,
		event.CorrelationId,
	)

	if err := h.createEmailCreatedEvent(ctx, event); err != nil {
		return err
	}

	if err := h.createEmailDeliveryEvent(ctx, event); err != nil {
		return err
	}

	return nil
}

func (h *MessageHandler) createEmailCreatedEvent(ctx context.Context, event EmailNotificationEvent) error {
	return h.createSystemEvent(
		ctx,
		event.TaskId,
		event.CorrelationId,
		"Email Service",
		"email:db:email-created",
		"Email successfully created in database",
		11,
	)
}

func (h *MessageHandler) createEmailDeliveryEvent(ctx context.Context, event EmailNotificationEvent) error {
	return h.createSystemEvent(
		ctx,
		event.TaskId,
		event.CorrelationId,
		"Email Service",
		"email:third-party:email-delivery-sent",
		"Email sent for delivery",
		14,
	)
}

func (h *MessageHandler) createSystemEvent(_ context.Context, taskID, correlationID, origin, action, message string, priority int) error {
	event := commons.TaskSystemEvent{
		TaskId:        taskID,
		CorrelationId: correlationID,
		Origin:        origin,
		Action:        action,
		Message:       message,
		JsonData:      "{}",
	}

	_, err := h.taskSystemEventRepository.Create(event, priority)
	if err != nil {
		return fmt.Errorf("failed to create system event: %w", err)
	}
	return nil
}
