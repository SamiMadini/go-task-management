package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

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

func getSQSClient(ctx context.Context) (*sqs.Client, error) {
	log.Println("Configuring SQS client...")

	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
		log.Printf("No AWS_REGION specified, defaulting to %s", region)
	}

	log.Printf("SQS Client: Endpoint=%s, Region=%s", endpoint, region)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			log.Printf("Using custom endpoint: %s", endpoint)
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}

		log.Printf("Using default AWS endpoint resolution for region: %s", region)
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	log.Println("Loading AWS config...")
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Printf("ERROR: Failed to load AWS config: %v", err)
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	log.Println("Successfully created SQS client")
	return sqs.NewFromConfig(cfg), nil
}

func sendToSQS(ctx context.Context, message string) error {
	log.Printf("Attempting to send message to SQS: %s", message)

	client, err := getSQSClient(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get SQS client: %v", err)
		return err
	}

	queueName := "go-email-service-queue"
	if name := os.Getenv("SQS_QUEUE_NAME"); name != "" {
		queueName = name
	}

	log.Printf("Getting URL for queue: %s", queueName)

	queueURLResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Printf("ERROR: Failed to get queue URL: %v", err)
		return fmt.Errorf("failed to get queue URL: %w", err)
	}

	log.Printf("Queue URL: %s", *queueURLResp.QueueUrl)

	_, err = client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    queueURLResp.QueueUrl,
		MessageBody: aws.String(message),
	})
	if err != nil {
		log.Printf("ERROR: Failed to send message to SQS: %v", err)
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	log.Printf("Message sent to SQS queue: %s", *queueURLResp.QueueUrl)
	return nil
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

	wg.Add(5)

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

	go func() {
		defer wg.Done()

		emailEvent := struct {
			TaskId        string `json:"taskId"`
			CorrelationId string `json:"correlationId"`
		}{
			TaskId:        taskId,
			CorrelationId: correlationId,
		}

		jsonBytes, err := json.Marshal(emailEvent)
		if err != nil {
			log.Printf("Error marshaling email event: %v", err)
			return
		}

		err = sendToSQS(context.Background(), string(jsonBytes))
		if err != nil {
			log.Printf("Error sending to SQS: %v", err)
			return
		}

		log.Printf("Successfully sent task notification to SQS queue for email processing: %s", taskId)
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
