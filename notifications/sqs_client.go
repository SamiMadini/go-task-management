package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClientInterface interface {
	SendMessage(ctx context.Context, message string) error
}

type SQSClient struct {
	client    *sqs.Client
	queueName string
	queueURL  *string
}

func NewSQSClient(ctx context.Context, cfg Config) (*SQSClient, error) {
	client, err := configureSQSClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	queueURL, err := getQueueURL(ctx, client, cfg.QueueName)
	if err != nil {
		return nil, err
	}

	return &SQSClient{
		client:    client,
		queueName: cfg.QueueName,
		queueURL:  queueURL,
	}, nil
}

func configureSQSClient(ctx context.Context, cfg Config) (*sqs.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.AWSEndpoint != "" {
			return aws.Endpoint{
				URL:           cfg.AWSEndpoint,
				SigningRegion: cfg.AWSRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.AWSRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return sqs.NewFromConfig(awsCfg), nil
}

func getQueueURL(ctx context.Context, client *sqs.Client, queueName string) (*string, error) {
	result, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get queue URL: %w", err)
	}
	return result.QueueUrl, nil
}

func (c *SQSClient) SendMessage(ctx context.Context, message string) error {
	_, err := c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    c.queueURL,
		MessageBody: aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}
	return nil
}
