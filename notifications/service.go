package main

import (
	"log"
	"sync"

	commons "sama/go-task-management/commons"
)

type NotificationService struct {
	taskRepository              commons.TaskRepositoryInterface
	taskSystemEventRepository   commons.TaskSystemEventRepositoryInterface
	inAppNotificationRepository commons.InAppNotificationRepositoryInterface
}

func NewNotificationService(taskRepository commons.TaskRepositoryInterface, taskSystemEventRepository commons.TaskSystemEventRepositoryInterface, inAppNotificationRepository commons.InAppNotificationRepositoryInterface) *NotificationService {
	return &NotificationService{taskRepository: taskRepository, taskSystemEventRepository: taskSystemEventRepository, inAppNotificationRepository: inAppNotificationRepository}
}

func (service *NotificationService) Handle(taskId string, correlationId string, notificationTypes []string) error {
	// for _, t := range notificationTypes {
	// 	log.Println("Notification Type:", t)
	// }
	
	task, err := service.taskRepository.GetByID(taskId)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	
	inAppNotificationParams := commons.InAppNotification{
		Title: task.Title,
		Description: task.Description,
	}

	_, err = service.inAppNotificationRepository.Create(inAppNotificationParams)
	if err != nil {
		log.Println("InAppNotification creation error: ", err)
		return err
	}

	wg.Add(4)

	go func() {
		defer wg.Done()
		eventInAppNotificationCreated := commons.TaskSystemEvent{
			TaskId: taskId,
			CorrelationId: correlationId,
			Origin: "Notification Service",
			Action: "notification:db:in-app-notification-created",
			Message: "In-app notification created in database",
			JsonData: "{}",
		}
		service.taskSystemEventRepository.Create(eventInAppNotificationCreated, 6)
	}()

	go func() {
		defer wg.Done()
		eventEventEmailTaskCreated := commons.TaskSystemEvent{
			TaskId: taskId,
			CorrelationId: correlationId,
			Origin: "Notification Service",
			Action: "notification:event:email-task-created",
			Message: "Email event sent",
			JsonData: "{}",
		}
		service.taskSystemEventRepository.Create(eventEventEmailTaskCreated, 9)
	}()

	// TEMP -> move it into email service
	go func() {
		defer wg.Done()
		eventEmailDbCreated := commons.TaskSystemEvent{
			TaskId: taskId,
			CorrelationId: correlationId,
			Origin: "Email Service",
			Action: "email:db:email-created",
			Message: "Email successfully created in database",
			JsonData: "{}",
		}
		service.taskSystemEventRepository.Create(eventEmailDbCreated, 11)
	}()

	// TEMP -> move it into email service
	go func() {
		defer wg.Done()
		eventEmailDeliverySent := commons.TaskSystemEvent{
			TaskId: taskId,
			CorrelationId: correlationId,
			Origin: "Email Service",
			Action: "email:third-party:email-delivery-sent",
			Message: "Email sent for delivery",
			JsonData: "{}",
		}
		service.taskSystemEventRepository.Create(eventEmailDeliverySent, 14)
	}()

	wg.Wait()

	task.EmailSent = true
	task.InAppSent = true

	err = service.taskRepository.Update(task)
	if err != nil {
		return err
	}

	return nil
}
