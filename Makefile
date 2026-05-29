.PHONY: all build test clean dev models

all: build

build:
	go build -o bin/speaking_hearts ./cmd/...

dev:
	go run ./cmd/server/main.go

models:
	go run scripts/download_models.go

test:
	go test ./...

clean:
	rm -rf bin/
	go clean
