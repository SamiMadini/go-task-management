FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

RUN mkdir -p /app/gateway /app/commons /app/notifications

COPY go.work go.work.sum ./

COPY gateway/go.mod gateway/go.sum /app/gateway/
COPY commons/ /app/commons/
COPY notifications/go.mod /app/notifications/

WORKDIR /app/gateway

COPY gateway/ ./

RUN cd /app && go mod download && cd /app/gateway && go build -o gateway

EXPOSE 8080

CMD ["./gateway"]
