package in_app_notification

import (
	"context"
	"time"

	"sama/go-task-management/commons"

	"github.com/google/uuid"
)

type Repository interface {
	Create(notification commons.InAppNotification) (commons.InAppNotification, error)
	GetByUserID(userID string) ([]commons.InAppNotification, error)
	UpdateOnRead(id string, isRead bool) error
	Delete(id string) error
}

type Service struct {
	logger           commons.Logger
	notificationRepo Repository
}

func NewService(logger commons.Logger, notificationRepo Repository) *Service {
	return &Service{
		logger:          logger,
		notificationRepo: notificationRepo,
	}
}

func (s *Service) CreateNotification(ctx context.Context, input CreateNotificationInput) (*NotificationResponse, error) {
	notification := commons.InAppNotification{
		ID:          uuid.New().String(),
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Message,
		IsRead:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdNotification, err := s.notificationRepo.Create(notification)
	if err != nil {
		return nil, err
	}

	response := toNotificationResponse(createdNotification)
	return &response, nil
}

func (s *Service) GetUserNotifications(ctx context.Context, userID string) ([]NotificationResponse, error) {
	notifications, err := s.notificationRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := make([]NotificationResponse, len(notifications))
	for i, notification := range notifications {
		response[i] = toNotificationResponse(notification)
	}

	return response, nil
}

func (s *Service) MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error {
	notifications, err := s.notificationRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if notification.ID == notificationID {
			return s.notificationRepo.UpdateOnRead(notificationID, true)
		}
	}

	return commons.ErrForbidden
}

func (s *Service) DeleteNotification(ctx context.Context, notificationID string, userID string) error {
	notifications, err := s.notificationRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if notification.ID == notificationID {
			return s.notificationRepo.Delete(notificationID)
		}
	}

	return commons.ErrForbidden
}
