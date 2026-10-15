ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: build dev test vet sqlc docker migrate-new migrate-up migrate-down
.DEFAULT_GOAL := build

APP_NAME := server
BUILD_DIR := bin

SQLC_DIR := internal/adapter/storage/sqlite
MIGRATIONS_DIR := $(SQLC_DIR)/migrations
DATABASE_PATH ?= data/app.db

MIGRATE := go run -tags sqlite github.com/golang-migrate/migrate/v4/cmd/migrate@v4.20.1
SQLC := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

dev:
	@go run github.com/air-verse/air@v1.67.4

test:
	go test -race ./...

vet:
	go vet ./...

# Regenerate the sqlc code in $(SQLC_DIR)/queries
sqlc:
	$(SQLC) generate -f $(SQLC_DIR)/sqlc.yaml

docker:
	docker build -t $(APP_NAME) .

# Usage: make migrate-new name=add_orders
migrate-new:
	@test -n "$(name)" || { echo "usage: make migrate-new name=add_orders"; exit 1; }
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "sqlite://$(DATABASE_PATH)" up

# Roll back the last migration only
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "sqlite://$(DATABASE_PATH)" down 1
