.PHONY: all build test clean dev

all: build

build:
	go build -o bin/speaking_hearts ./cmd/...

dev:
	go run ./cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin/
	go clean
