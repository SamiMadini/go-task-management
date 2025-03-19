package main

const (
	defaultRegion    = "us-east-1"
	defaultQueueName = "go-email-service-queue"
)

type Config struct {
	AWSEndpoint string
	AWSRegion   string
	QueueName   string
}
