# Makefile for mcpgolang project

all: build

build:
	@echo "Building project..."
	@go build -o bin/eino-mcp cmd/eino-mcp/main.go
	@go build -o bin/mcp-mysql cmd/mcp-mysql/main.go
