.PHONY: help up down reset env sqlc swag deps test build run dev stop logs

COMPOSE := docker compose
BACKEND := backend
ENV_FILE := $(BACKEND)/.env

help:
	@echo "Planeweave — локальный запуск"
	@echo ""
	@echo "  make dev      — PostgreSQL + Swagger + API"
	@echo "  make up       — только PostgreSQL (docker compose up -d)"
	@echo "  make run      — Swagger + API (PostgreSQL должен быть запущен)"
	@echo "  make down     — остановить PostgreSQL"
	@echo "  make reset    — удалить volume PostgreSQL и поднять заново"
	@echo "  make swag     — сгенерировать Swagger из кода (swaggo)"
	@echo "  make sqlc     — сгенерировать код из SQL (sqlc)"
	@echo "  make test     — go test ./..."
	@echo "  make build    — собрать бинарник backend/bin/server"
	@echo "  make deps     — go mod tidy"
	@echo "  make logs     — логи PostgreSQL"
	@echo ""
	@echo "  Swagger UI:  http://localhost:8080/swagger/index.html"

up: env
	$(COMPOSE) up -d
	@echo "PostgreSQL: localhost:5432 (planeweave / planeweave)"

down:
	$(COMPOSE) down

reset: down
	$(COMPOSE) down -v
	$(COMPOSE) up -d
	@echo "База пересоздана"

env:
	@test -f $(ENV_FILE) || cp $(BACKEND)/.env.example $(ENV_FILE)

swag:
	cd $(BACKEND) && go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/server/main.go -o docs/swag --parseInternal --parseDependency
	cp $(BACKEND)/docs/swag/swagger.json $(BACKEND)/api/openapi.json

sqlc:
	cd $(BACKEND) && sqlc generate

deps:
	cd $(BACKEND) && go mod tidy

test: deps swag
	cd $(BACKEND) && go test ./...

build: deps sqlc swag
	cd $(BACKEND) && go build -o bin/server ./cmd/server

run: env swag
	cd $(BACKEND) && go run ./cmd/server

dev: up deps swag
	cd $(BACKEND) && go run ./cmd/server

stop: down

logs:
	$(COMPOSE) logs -f postgres
