FROM golang:1.24-bullseye AS builder

WORKDIR /build

COPY . ./

WORKDIR /build/email-service/src

RUN go mod download && go build -o bootstrap

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
