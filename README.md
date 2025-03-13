# Task Management System with Notification Microservice

### Task Management API

- RESTful API built with Go and Gin framework
- SQLite database for persistence in local
- Swagger documentation
- 5=6 routes:
  - GET /api/_health - Health check route
  - GET /api/tasks - List all tasks
  - POST /api/tasks - Create a new task
  - GET /api/tasks/{id} - Get task details
  - PUT /api/tasks/{id} - Update a task
  - DELETE /api/tasks/{id} - Delete a task

### Notification Microservice

- Event-driven communication with the API via gRPC
- Two use cases:
  - InApp notifications
  - Email notifications

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Swagger](https://github.com/swaggo/swag) - API documentation
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [uuid](https://github.com/google/uuid) - UUID generation

## Installation

1. Install swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

#### Troubleshooting
To do this temporarily for the current session:

> export PATH=$PATH:$HOME/go/bin

To permanently add it to your PATH, you can edit your shell profile (.bashrc or .zshrc, depending on your shell):

**For Bash**

> echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
> source ~/.bashrc

**For Zsh**

> echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
> source ~/.zshrc

After updating your PATH, verify that the swag command works:

> swag --version

If you want to install swag globally and avoid the PATH issue altogether, you can move the swag binary to a directory like /usr/local/bin which is usually in your PATH by default:

> sudo mv $HOME/go/bin/swag /usr/local/bin/


2. Init multi-modules

```bash
make init-multi-modules
```

3. Install dependencies

```bash
make install
```

4. Install Go gRPC

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
make start
```

The server will start on port :8080 by default.
