SHELL=/bin/bash
SERVICE_NAME=$(notdir $(shell pwd))
SSH_PRIVATE_KEY := $(shell cat ~/.ssh/id_ed25519)
DB_USER := $(shell grep DB_USER .env | cut -d '=' -f 2)
DB_PASSWORD := $(shell grep DB_PASSWORD .env | cut -d '=' -f 2)
DB_HOST := $(shell grep DB_HOST .env | cut -d '=' -f 2)
DB_PORT := $(shell grep DB_PORT .env | cut -d '=' -f 2)
DB_NAME := $(shell grep DB_NAME .env | cut -d '=' -f 2)

.PHONY: download
download:
	go mod download
	go mod tidy -compat=1.21

.PHONY: swag
swag:
	rm -rf ./docs
	swag init --parseDependency -g internal/infra/api/router/router.go
	swag fmt

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint -v run

.PHONY: fix-imports
fix-imports:
	@folder=$$(basename "$(PWD)"); \
	goimports-reviser -rm-unused -project-name=$$folder  -format ./...

.PHONY: check-imports
check-imports:
	@folder=$$(basename "$(PWD)"); \
	goimports-reviser -list-diff -project-name=$$folder ./...

.PHONY: mocks
mocks:
	mockery --config mockeryconfig.yml --all

.PHONY: requirements-up
requirements-up:
	docker-compose -f docker-compose.requirement.local.yml up -d

.PHONY: requirements-down
requirements-down:
	docker-compose -f docker-compose.requirement.local.yml down

.PHONY: requirements-w-up
requirements-w-up:
	docker-compose -f docker-compose.requirement.local.window.yml up -d

.PHONY: requirements-w-down
requirements-w-down:
	docker-compose -f docker-compose.requirement.local.window.yml down

.PHONY: compose-up
compose-up:
	SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_ed25519)" docker-compose -f docker-compose.yml up --build --remove-orphans

.PHONY: compose-down
compose-down:
	docker-compose -f docker-compose.yml down

.PHONY: docker-build
docker-build:
	docker build -f Dockerfile -t $(SERVICE_NAME):latest --build-arg SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_ed25519)" .

.PHONY: docker-build-aws
docker-build-aws:
	docker build --platform=linux/amd64 -f Dockerfile -t $(SERVICE_NAME):latest --build-arg SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_ed25519)" .

.PHONY: docker-remove
docker-remove:
	docker rmi -f $(SERVICE_NAME):latest

.PHONY: docker-run
docker-run:
	docker run --name $(SERVICE_NAME) --env-file .env.local -e DB_HOST=host.docker.internal -p 3003:3003 $(SERVICE_NAME):latest

.PHONY: docker-full-aws
docker-full-aws:
	docker build --platform=linux/amd64 -f Dockerfile -t $(SERVICE_NAME):latest --build-arg SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_ed25519)" .
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 451472926433.dkr.ecr.us-east-1.amazonaws.com
	docker tag $(SERVICE_NAME):latest 451472926433.dkr.ecr.us-east-1.amazonaws.com/$(SERVICE_NAME):latest
	docker push 451472926433.dkr.ecr.us-east-1.amazonaws.com/$(SERVICE_NAME):latest

.PHONY: deploy
deploy:
	aws ecs update-service --cluster DevCluster --service  $(SERVICE_NAME) --force-new-deployment

.PHONY: replace-path
replace-path:
	@folder=$$(basename "$(PWD)"); \
	LC_CTYPE=C find . -type f -not -name "makefile" -exec sed -i '' 's/\/templates\/go-api/\/backend\/'$$folder'/g' {} \;

.PHONY: start-dev
start-dev:
	ENVIRONMENT=local go run cmd/main.go -v

.PHONY: start-temporal-dev
start-temporal-dev:
	temporal server start-dev --ui-port 8080 --port 7234

.PHONY: start-localstack-dev
start-localstack-dev:
	localstack start

migration-files:
	@read -p "Enter the title: " title; \
	migrate create -ext sql -dir ./db/migrations -seq $$title

migrateup:
	migrate -path ./db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

.PHONY: watch-coverage
watch-coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out