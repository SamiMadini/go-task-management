package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	commons "sama/go-task-management/commons"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var emailService *EmailService

func init() {
	logFile, err := os.OpenFile("/tmp/email-service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Warning: Unable to open log file: %v", err)
	} else {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}
	
	log.Println("======= Email Service Starting Up =======")
	// log.Printf("Environment variables:")
	// for _, env := range os.Environ() {
	// 	log.Println(env)
	// }

	database, err := commons.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	eventRepo := commons.NewPostgresTaskSystemEventRepository(database)
	emailService = NewEmailService(eventRepo)
}

type EmailNotificationEvent struct {
	TaskId        string `json:"taskId"`
	CorrelationId string `json:"correlationId"`
}

type EmailService struct {
	taskSystemEventRepository commons.TaskSystemEventRepositoryInterface
}

func NewEmailService(eventRepo commons.TaskSystemEventRepositoryInterface) *EmailService {
	return &EmailService{
		taskSystemEventRepository: eventRepo,
	}
}

func (s *EmailService) createSystemEvent(_ context.Context, taskID, correlationID, origin, action, message string, priority int) error {
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
		return fmt.Errorf("failed to create system event: %w", err)
	}
	return nil
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


func processQueueMessages(ctx context.Context) {
	log.Println("Starting SQS message processor...")
	
	go func() {
		retryCount := 0
		backoff := 5 * time.Second
		maxBackoff := 30 * time.Second
		maxRetries := 10
		
		for retryCount < maxRetries {
			client, err := getSQSClient(ctx)
			if err != nil {
				log.Printf("Attempt %d: Failed to create SQS client: %v. Retrying in %v", 
					retryCount+1, err, backoff)
				time.Sleep(backoff)
				backoff = time.Duration(float64(backoff) * 1.5)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				retryCount++
				continue
			}

			queueName := "go-email-service-queue"
			if name := os.Getenv("SQS_QUEUE_NAME"); name != "" {
				queueName = name
			}
			
			log.Printf("Getting URL for queue: %s (attempt %d)", queueName, retryCount+1)

			queueURLResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				log.Printf("Attempt %d: Failed to get queue URL: %v. Retrying in %v", 
					retryCount+1, err, backoff)
				time.Sleep(backoff)
				backoff = time.Duration(float64(backoff) * 1.5)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				retryCount++
				continue
			}
			
			queueURL := queueURLResp.QueueUrl
			log.Printf("Successfully connected to SQS queue: %s", *queueURL)
			
			pollSQSQueue(ctx, client, queueURL)
			return
		}
		
		log.Printf("Failed to connect to SQS after %d attempts. HTTP endpoint will still work.", maxRetries)
	}()
}

func pollSQSQueue(ctx context.Context, client *sqs.Client, queueURL *string) {
	log.Printf("Starting to poll SQS queue: %s", *queueURL)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in SQS polling: %v", r)
			time.Sleep(5 * time.Second)
			go pollSQSQueue(ctx, client, queueURL)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping SQS polling")
			return
		default:
			log.Println("Polling for messages...")
			try_receive_message(ctx, client, queueURL)

			time.Sleep(5 * time.Second)
		}
	}
}

func try_receive_message(ctx context.Context, client *sqs.Client, queueURL *string) {
	receiveCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	resp, err := client.ReceiveMessage(receiveCtx, &sqs.ReceiveMessageInput{
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	
	if err != nil {
		log.Printf("Error receiving messages: %v", err)
		return
	}

	log.Printf("Received %d messages", len(resp.Messages))
	
	for _, msg := range resp.Messages {
		log.Printf("Processing message: %s", *msg.MessageId)
		log.Printf("Message body: %s", *msg.Body)
		
		err := handleRequest(ctx, []byte(*msg.Body))
		if err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		_, err = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("Error deleting message: %v", err)
		}
	}
}

func handleSQSEvent(ctx context.Context, event events.SQSEvent) error {
	log.Printf("Received SQS event with %d records", len(event.Records))
	
	for i, record := range event.Records {
		log.Printf("Processing record %d: %s", i, record.MessageId)
		log.Printf("Record body: %s", record.Body)
		
		err := handleRequest(ctx, []byte(record.Body))
		if err != nil {
			log.Printf("Error processing SQS event: %v", err)
			return err
		}
	}
	return nil
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	log.Printf("Handling request with payload: %s", string(event))
	
	var emailNotificationEvent EmailNotificationEvent
	if err := json.Unmarshal(event, &emailNotificationEvent); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf("Email service processing message")

	// log.Printf("PostgreSQL Configuration: Host=%s, Port=%s, User=%s, DB=%s",
	// 	os.Getenv("POSTGRES_HOST"),
	// 	os.Getenv("POSTGRES_PORT"),
	// 	os.Getenv("POSTGRES_USER"),
	// 	os.Getenv("POSTGRES_DB"))

	log.Printf("Request parameters: taskId=%s, correlationId=%s",
		emailNotificationEvent.TaskId, emailNotificationEvent.CorrelationId)

	if err := emailService.createSystemEvent(
		ctx,
		emailNotificationEvent.TaskId,
		emailNotificationEvent.CorrelationId,
		"Email Service",
		"email:db:email-created",
		"Email successfully created in database",
		11,
	); err != nil {
		log.Printf("Error creating email created event: %v", err)
		return err
	}

	if err := emailService.createSystemEvent(
		ctx,
		emailNotificationEvent.TaskId,
		emailNotificationEvent.CorrelationId,
		"Email Service",
		"email:third-party:email-delivery-sent",
		"Email sent for delivery",
		14,
	); err != nil {
		log.Printf("Error creating email delivery event: %v", err)
		return err
	}

	return nil
}

func startHTTPServer(ctx context.Context) {
	log.Println("Starting HTTP server...")
	
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	
	log.Printf("HTTP server will listen on port: %s", port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received %s request from %s", r.Method, r.RemoteAddr)
			
			if r.Method == http.MethodGet {
				log.Println("Health check request received")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Email service is running"))
				return
			}
			
			if r.Method != http.MethodPost {
				log.Printf("Method not allowed: %s", r.Method)
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var payload json.RawMessage
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				log.Printf("Invalid request body: %v", err)
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			
			log.Printf("Received payload: %s", string(payload))

			err := handleRequest(r.Context(), payload)
			if err != nil {
				log.Printf("Error handling request: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			jsonBytes, _ := json.Marshal(payload)
			err = sendToSQS(r.Context(), string(jsonBytes))
			if err != nil {
				log.Printf("Warning: Failed to send to SQS: %v", err)
			}

			log.Println("Request processed successfully")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Message processed successfully"))
		}),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Email service HTTP server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	go processQueueMessages(ctx)
	<-stop

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in main: %v", r)
		}
	}()
	
	log.Println("Starting email service...")
	ctx := context.Background()
	
	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		log.Println("Running in AWS Lambda environment")
		lambda.Start(func(ctx context.Context, event events.SQSEvent) error {
			return handleSQSEvent(ctx, event)
		})
	} else {
		log.Println("Starting local development environment")
		startHTTPServer(ctx)
	}
}
