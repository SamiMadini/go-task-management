#!/bin/bash

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "AWS CLI is not installed. Please install it and try again."
    exit 1
fi

# Default values - always use the environment variable if available
ENDPOINT_URL=${AWS_ENDPOINT_URL:-"http://localhost:9324"}
QUEUE_NAME=${SQS_QUEUE_NAME:-"go-email-service-queue"}
DLQ_QUEUE_NAME="go-email-service-queue-dead-letters"
REGION=${AWS_REGION:-"us-east-1"}

echo "Using SQS endpoint: ${ENDPOINT_URL}"
echo "Creating SQS queues in ElasticMQ..."

# Create the dead letter queue first
aws --endpoint-url="${ENDPOINT_URL}" sqs create-queue \
    --queue-name "${DLQ_QUEUE_NAME}" \
    --region ${REGION}

# Get the DLQ queue URL
DLQ_QUEUE_URL=$(aws --endpoint-url="${ENDPOINT_URL}" sqs get-queue-url \
    --queue-name "${DLQ_QUEUE_NAME}" \
    --region ${REGION} \
    --query 'QueueUrl' \
    --output text)

echo "DLQ Queue URL: ${DLQ_QUEUE_URL}"

# Get the DLQ queue ARN
DLQ_QUEUE_ARN=$(aws --endpoint-url="${ENDPOINT_URL}" sqs get-queue-attributes \
    --queue-url "${DLQ_QUEUE_URL}" \
    --attribute-names QueueArn \
    --region ${REGION} \
    --query 'Attributes.QueueArn' \
    --output text)

echo "DLQ Queue ARN: ${DLQ_QUEUE_ARN}"

# Create the main queue with redrive policy
aws --endpoint-url="${ENDPOINT_URL}" sqs create-queue \
    --queue-name "${QUEUE_NAME}" \
    --attributes "{\"VisibilityTimeout\":\"300\", \"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"${DLQ_QUEUE_ARN}\\\",\\\"maxReceiveCount\\\":\\\"3\\\"}\"}" \
    --region ${REGION}

echo "SQS queues created successfully!"

# List the queues to verify
echo "Available queues:"
aws --endpoint-url="${ENDPOINT_URL}" sqs list-queues --region ${REGION} 