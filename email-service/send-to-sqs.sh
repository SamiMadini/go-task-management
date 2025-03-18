#!/bin/bash

if ! command -v aws &> /dev/null; then
    echo "AWS CLI is not installed. Please install it and try again."
    exit 1
fi

ENDPOINT_URL=${AWS_ENDPOINT_URL:-"http://localhost:9324"}
QUEUE_NAME=${SQS_QUEUE_NAME:-"go-email-service-queue"}
TASK_ID=${1:-"test-task-id"}
CORRELATION_ID=${2:-"test-correlation-id"}
REGION=${AWS_REGION:-"us-east-1"}

QUEUE_URL="${ENDPOINT_URL}/queue/${QUEUE_NAME}"

MESSAGE="{\"taskId\":\"${TASK_ID}\",\"correlationId\":\"${CORRELATION_ID}\"}"

echo "Using SQS endpoint: ${ENDPOINT_URL}"
echo "Sending message to SQS queue: ${QUEUE_URL}"
echo "Message body: ${MESSAGE}"

aws --endpoint-url="${ENDPOINT_URL}" sqs send-message \
    --queue-url "${QUEUE_URL}" \
    --message-body "${MESSAGE}" \
    --region ${REGION}

echo "Message sent successfully!" 