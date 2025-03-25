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

# Copy source code
COPY commons/ /app/commons/
COPY gateway/ /app/gateway/
COPY notification-service/ /app/notification-service/
COPY email-service/src/ /app/email-service/src/
COPY go.work go.work.sum ./

WORKDIR /app/gateway

# Build the application
RUN go build -o gateway

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=builder /app/gateway/gateway /app/

EXPOSE 8080

CMD ["./gateway"]
