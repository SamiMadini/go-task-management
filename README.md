# Task Management System with Notification Microservice

### Docs
#### System Architecture

The Task Management System is built with a microservices architecture, providing a scalable and maintainable solution with notification capabilities.

![System Architecture](/frontend/public/images/full1.png)

#### Database Schema

The system uses PostgreSQL for data persistence with the following schema:

![Database Schema](/frontend/public//images/full2.png)

#### User Interface

The frontend is built with Next.js and Tailwind CSS, providing an intuitive user experience:

![User Interface](/frontend/public//images/full3.png)

## Technology Stack

- **API Gateway**: Custom implementation in Go (Golang)
- **Frontend**: TypeScript, Next.js, Tailwind CSS
- **Database**: PostgreSQL
- **Communication**: gRPC for service-to-service communication
- **Containerization**: Docker
- **AWS**: Lambda, EKS, CDK


### API

- RESTful API
- Postgres database for persistence
- Swagger documentation
- Routes:
  - GET /api/_health - Health check route

  - GET /api/tasks - List all tasks
  - POST /api/tasks - Create a new task
  - GET /api/tasks/{id} - Get task details
  - PUT /api/tasks/{id} - Update a task
  - DELETE /api/tasks/{id} - Delete a task

  - GET /api/notifications
  - POST /api/notifications/{id}/read
  - DELETE /api/notifications/{id}

  - GET /api/task-system-events

### Notification Microservice

- Event-driven communication with gRPC
- Event-driven communication with AWS SQS
- Two use cases:
  - InApp notifications
  - Email notifications

## Dependencies

- [Swagger](https://github.com/swaggo/swag) - API documentation
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [uuid](https://github.com/google/uuid) - UUID generation

## Access to Swagger API Documentation
http://localhost:3012/swagger/



## Installation

1. Init multi-modules

```bash
make init-multi-modules
```

2. Install dependencies

```bash
make install
```

3. Install Go gRPC

With brew
```bash
brew install protobuf
```

Or with Go
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Then update your PATH so that the protoc compiler can find the plugins:

> echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
> source ~/.zshrc

5. Install Air to have live reload

```bash
go install github.com/air-verse/air@latest
```

6. Start the application

```bash
make dc-start-with-build
```

The server will start on port :8080 by default.
