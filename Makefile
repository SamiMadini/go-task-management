init-multi-modules:
	go work init ./commons ./gateway ./notification-service

install-dependencies: install-dependencies-commons install-dependencies-gateway install-dependencies-notifications

install-dependencies-commons:
	cd ./commons && go mod tidy

install-dependencies-gateway:
	cd ./gateway && go mod tidy

install-dependencies-notifications:
	cd ./notification-service && go mod tidy

start-dev: dev-start-gateway

dev-start-gateway:
	cd ./gateway && air

dev-start-notifications-service:
	cd ./notification-service && go run main.go


dc-start:
	docker compose up

dc-start-with-build:
	docker compose up --build

dc-stop:
	docker compose down

dc-restart:
	docker compose down
	docker compose up

dc-build:
	docker compose build

dc-build-gateway:
	docker compose build gateway

dc-build-notifications:
	docker compose build notifications

dc-reset-db:
	docker volume rm go-task-management_postgres-data


open-psql:
	psql -h localhost -p 5433 -U postgres -d tasks

.PHONY: lint lint-fix install-linters

# Install golangci-lint if not already installed
install-linters:
	@if ! command -v golangci-lint &> /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.56.2; \
		cp $$(go env GOPATH)/bin/golangci-lint /usr/local/bin/; \
	fi

# Run linter on all modules
lint:
	docker run --rm -v $(PWD):/app -w /app/gateway golangci/golangci-lint:v1.64.6 golangci-lint run ./...
	docker run --rm -v $(PWD):/app -w /app/notification-service golangci/golangci-lint:v1.64.6 golangci-lint run ./...
	docker run --rm -v $(PWD):/app -w /app/email-service golangci/golangci-lint:v1.64.6 golangci-lint run ./...
	docker run --rm -v $(PWD):/app -w /app/commons golangci/golangci-lint:v1.64.6 golangci-lint run ./...

# Run linter with auto-fix on all modules
lint-fix:
	docker run --rm -v $(PWD):/app -w /app/gateway golangci/golangci-lint:v1.64.6 golangci-lint run --fix ./...
	docker run --rm -v $(PWD):/app -w /app/notification-service golangci/golangci-lint:v1.64.6 golangci-lint run --fix ./...
	docker run --rm -v $(PWD):/app -w /app/email-service golangci/golangci-lint:v1.64.6 golangci-lint run --fix ./...
	docker run --rm -v $(PWD):/app -w /app/commons golangci/golangci-lint:v1.64.6 golangci-lint run --fix ./...

# Run linter on a specific package
lint-pkg:
	@if [ "$(pkg)" = "" ]; then \
		echo "Error: pkg is not set. Usage: make lint-pkg pkg=gateway"; \
		exit 1; \
	fi
	docker run --rm -v $(PWD):/app -w /app/$(pkg) golangci/golangci-lint:v1.64.6 golangci-lint run ./...
	