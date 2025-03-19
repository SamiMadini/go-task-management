package sqs

import (
	"context"
	"fmt"
	"log"
	"time"

	"sama/go-task-management/email-service/src/config"
	"sama/go-task-management/email-service/src/handlers"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSManager struct {
	client    *sqs.Client
	queueURL  *string
	config    *config.Config
	handler   *handlers.MessageHandler
}

func NewSQSManager(cfg *config.Config, handler *handlers.MessageHandler) (*SQSManager, error) {
	ctx := context.Background()
	client, err := getSQSClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	queueURL, err := getQueueURL(ctx, client, cfg.QueueName)
	if err != nil {
		return nil, err
	}

	return &SQSManager{
		client:   client,
		queueURL: queueURL,
		config:   cfg,
		handler:  handler,
	}, nil
}

func getSQSClient(ctx context.Context, cfg *config.Config) (*sqs.Client, error) {
	log.Println("Configuring SQS client...")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.AWSEndpoint != "" {
			return aws.Endpoint{
				URL:           cfg.AWSEndpoint,
				SigningRegion: cfg.AWSRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return sqs.NewFromConfig(awsCfg), nil
}

func getQueueURL(ctx context.Context, client *sqs.Client, queueName string) (*string, error) {
	resp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get queue URL: %w", err)
	}
	return resp.QueueUrl, nil
}

func (m *SQSManager) SendMessage(ctx context.Context, message string) error {
	_, err := m.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    m.queueURL,
		MessageBody: aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}
	return nil
}

func (m *SQSManager) StartPolling(ctx context.Context) {
	log.Printf("Starting to poll SQS queue: %s", *m.queueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Context canceled, stopping SQS polling")
				return
			default:
				m.pollMessages(ctx)
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (m *SQSManager) pollMessages(ctx context.Context) {
	receiveCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := m.client.ReceiveMessage(receiveCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            m.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})

	if err != nil {
		log.Printf("Error receiving messages: %v", err)
		return
	}

	for _, msg := range resp.Messages {
		if err := m.handler.HandleMessage(ctx, []byte(*msg.Body)); err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		if err := m.deleteMessage(ctx, msg.ReceiptHandle); err != nil {
			log.Printf("Error deleting message: %v", err)
		}
	}
}

func (m *SQSManager) deleteMessage(ctx context.Context, receiptHandle *string) error {
	_, err := m.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      m.queueURL,
		ReceiptHandle: receiptHandle,
	})
	return err
}
