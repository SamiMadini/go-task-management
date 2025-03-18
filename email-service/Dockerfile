FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY src/ .

RUN go mod tidy
RUN GOCACHE=/tmp GOARCH=amd64 GOOS=linux go build -o /app/bootstrap

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/bootstrap /app/bootstrap

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=tasks

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

CMD ["/app/bootstrap"]
