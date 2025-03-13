init-multi-modules:
	go work init ./commons ./gateway ./notifications

install-dependencies: install-dependencies-commons install-dependencies-gateway install-dependencies-notifications

install-dependencies-commons:
	cd ./commons && go mod tidy

install-dependencies-gateway:
	cd ./gateway && go mod tidy

install-dependencies-notifications:
	cd ./notifications && go mod tidy

start-dev: dev-start-gateway

dev-start-gateway:
	cd ./gateway && air

dev-start-notifications-service:
	cd ./notifications && go run main.go


dc-start:
	docker compose up --build

dc-stop:
	docker compose down

dc-restart:
	docker compose down
	docker compose up --build

dc-build:
	docker compose build

dc-build-gateway:
	docker compose build gateway

dc-build-notifications:
	docker compose build notifications
