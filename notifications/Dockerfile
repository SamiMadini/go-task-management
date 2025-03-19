FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY commons/ /app/commons/
COPY gateway/ /app/gateway/
COPY notifications/ /app/notifications/
COPY email-service/src/ /app/email-service/src/
COPY go.work go.work.sum ./

WORKDIR /app/notifications

RUN cd /app && go mod download && cd /app/notifications && go build -o notifications

EXPOSE 2000

ENV AWS_REGION=us-east-1
ENV AWS_ENDPOINT_URL=http://elasticmq:9324
ENV AWS_ACCESS_KEY_ID=localstack
ENV AWS_SECRET_ACCESS_KEY=localstack
ENV SQS_QUEUE_NAME=go-email-service-queue

CMD ["./notifications"]
