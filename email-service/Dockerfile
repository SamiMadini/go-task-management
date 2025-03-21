FROM golang:1.24-bullseye AS deps

WORKDIR /deps

# Copy only the files needed for downloading dependencies
COPY commons/go.mod commons/go.sum ./commons/
COPY email-service/src/go.mod email-service/src/go.sum ./email-service/src/

# Download dependencies
RUN cd /deps/email-service/src && go mod download

FROM golang:1.24-bullseye AS builder

WORKDIR /build

# Copy dependencies from deps stage
COPY --from=deps /go/pkg /go/pkg

# Copy source code
COPY commons/ /build/commons/
COPY email-service/src/ /build/email-service/src/

WORKDIR /build/email-service/src

# Build the application
RUN go build -o bootstrap

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /build/email-service/src/bootstrap /app/bootstrap

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=tasks

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

CMD ["/app/bootstrap"]
