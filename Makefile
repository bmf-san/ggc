# Makefile for Go project

APP_NAME=ggc

.PHONY: all build run clean test lint cover

all: build

build:
	go build -o $(APP_NAME) main.go

run: build
	./$(APP_NAME)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(APP_NAME)

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out