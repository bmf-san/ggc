# Makefile for Go project

APP_NAME=ggc

.PHONY: all build run clean test lint

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