FROM golang:1.24-alpine

# Install build dependencies for SQLite
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

# Set working directory to gateway for building
WORKDIR /app/gateway

# Copy gateway source code
COPY gateway/ ./

# Build the application with CGO enabled
RUN cd /app && CGO_ENABLED=1 go mod download && cd /app/gateway && CGO_ENABLED=1 go build -o gateway

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./gateway"] 