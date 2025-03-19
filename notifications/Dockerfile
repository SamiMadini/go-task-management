FROM golang:1.24-alpine

# We still need build tools for some dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Create directory structure
RUN mkdir -p /app/gateway /app/commons /app/notifications

# Copy go.work files to root
COPY go.work go.work.sum ./

# Copy module files maintaining directory structure
COPY gateway/go.mod gateway/go.sum /app/gateway/
COPY commons/ /app/commons/
COPY notifications/go.mod /app/notifications/

# Set working directory to notifications for building
WORKDIR /app/notifications

# Copy notifications source code
COPY notifications/ ./

# Build the application
RUN cd /app && go mod download && cd /app/notifications && go build -o notifications

# Expose port
EXPOSE 2000

# Set AWS environment variables with defaults
ENV AWS_REGION=us-east-1
ENV AWS_ENDPOINT_URL=http://elasticmq:9324
ENV AWS_ACCESS_KEY_ID=localstack
ENV AWS_SECRET_ACCESS_KEY=localstack
ENV SQS_QUEUE_NAME=go-email-service-queue

# Run the binary
CMD ["./notifications"] 