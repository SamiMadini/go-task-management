package adapters

import (
	"sama/go-task-management/commons"
)

type InAppNotificationRepositoryAdapter struct {
	commons.InAppNotificationRepositoryInterface
}

func (a *InAppNotificationRepositoryAdapter) Create(notification commons.InAppNotification) (commons.InAppNotification, error) {
	created, err := a.InAppNotificationRepositoryInterface.Create(notification)
	if err != nil {
		return commons.InAppNotification{}, err
	}
	return created, nil
}
