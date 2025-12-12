# Makefile for todo-ai-assistant Go project

BINARY_NAME=todo

.PHONY: build test clean

build:
	go build -o $(BINARY_NAME) ./...

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
