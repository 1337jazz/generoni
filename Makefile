ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: build dev
.DEFAULT_GOAL := build

APP_NAME := server
BUILD_DIR := bin

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

dev:
	@go run github.com/air-verse/air@v1.67.4
