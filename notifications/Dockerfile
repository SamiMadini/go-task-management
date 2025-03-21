FROM golang:1.24-alpine AS deps

RUN apk add --no-cache gcc musl-dev

WORKDIR /deps

# Copy only the files needed for downloading dependencies
COPY go.work go.work.sum ./
COPY commons/go.mod commons/go.sum ./commons/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY notifications/go.mod notifications/go.sum ./notifications/
COPY email-service/src/go.mod email-service/src/go.sum ./email-service/src/

# Download dependencies
RUN go mod download

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /go/pkg /go/pkg

# Copy source code
COPY commons/ /app/commons/
COPY gateway/ /app/gateway/
COPY notifications/ /app/notifications/
COPY email-service/src/ /app/email-service/src/
COPY go.work go.work.sum ./

WORKDIR /app/notifications

# Build the application
RUN go build -o notifications

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=builder /app/notifications/notifications /app/

EXPOSE 2000

ENV AWS_REGION=us-east-1
ENV AWS_ENDPOINT_URL=http://elasticmq:9324
ENV AWS_ACCESS_KEY_ID=localstack
ENV AWS_SECRET_ACCESS_KEY=localstack
ENV SQS_QUEUE_NAME=go-email-service-queue

CMD ["./notifications"]
