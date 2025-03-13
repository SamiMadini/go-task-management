package main

import (
	"fmt"
	commons "sama/go-task-management/commons"
	"time"
)

type NotificationService struct {
	taskRepository              commons.TaskRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
}

func NewNotificationService(taskRepository commons.TaskRepositoryInterface, inAppNotificationRepository commons.InAppNotificationRepositoryInterface) *NotificationService {
	return &NotificationService{taskRepository: taskRepository, inAppNotificationRepository: inAppNotificationRepository}
}

func (service *NotificationService) Handle(taskId string, notificationTypes []string) error {
	fmt.Println("Task ID:", taskId)
	
	for _, t := range notificationTypes {
		fmt.Println("Notification Type:", t)
	}
	
	fmt.Println("Waiting for 4 seconds")
	time.Sleep(4 * time.Second)
	
	task, err := service.taskRepository.GetByID(taskId)
	if err != nil {
		return err
	}

	inAppNotificationParams := commons.InAppNotification{
		Title: task.Title,
		Description: task.Description,
	}

	inAppNotification, err := service.inAppNotificationRepository.Create(inAppNotificationParams)
	if err != nil {
		return err
	}

	fmt.Println("InAppNotification created:", inAppNotification)

	fmt.Println("--------------------")

	task.EmailSent = true
	task.InAppSent = true

	err = service.taskRepository.Update(task)
	if err != nil {
		return err
	}

	return nil
}
