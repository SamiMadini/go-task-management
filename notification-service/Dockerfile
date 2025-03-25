FROM golang:1.24-alpine AS deps

RUN apk add --no-cache gcc musl-dev

WORKDIR /deps

# Copy only the files needed for downloading dependencies
COPY go.work go.work.sum ./
COPY commons/go.mod commons/go.sum ./commons/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY email-service/src/go.mod email-service/src/go.sum ./email-service/src/

# Download dependencies
RUN go mod download

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /go/pkg /go/pkg

# Copy go.mod and go.sum files for all modules
COPY commons/go.mod commons/go.sum ./commons/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY email-service/src/go.mod email-service/src/go.sum ./email-service/src/
COPY go.work go.work.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY commons/ /app/commons/
COPY gateway/ /app/gateway/
COPY notification-service/ /app/notification-service/
COPY email-service/src/ /app/email-service/src/

# Set the working directory to the notification service
WORKDIR /app/notification-service

# Build the application
RUN go build -o notification-service

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=builder /app/notification-service/notification-service /app/

EXPOSE 2000

ENV AWS_REGION=us-east-1
ENV AWS_ENDPOINT_URL=http://elasticmq:9324
ENV AWS_ACCESS_KEY_ID=localstack
ENV AWS_SECRET_ACCESS_KEY=localstack
ENV SQS_QUEUE_NAME=go-email-service-queue

CMD ["./notification-service"]
