FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY commons/ /app/commons/
COPY gateway/ /app/gateway/
COPY notifications/ /app/notifications/
COPY email-service/src/ /app/email-service/src/
COPY go.work go.work.sum ./

WORKDIR /app/gateway

RUN cd /app && go mod download && cd /app/gateway && go build -o gateway

EXPOSE 8080

CMD ["./gateway"]
