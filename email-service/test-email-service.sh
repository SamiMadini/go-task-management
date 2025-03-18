#!/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
  "taskId": "test-task-id",
  "correlationId": "test-correlation-id"
}' http://localhost:8080

echo ""
echo "Message sent to email service"

# Optional: Send a message to the SQS queue (requires awscli)
# aws --endpoint-url=http://localhost:9324 sqs send-message \
#     --queue-url http://localhost:9324/queue/go-email-service-queue \
#     --message-body '{"taskId": "test-task-id", "correlationId": "test-correlation-id"}' 