.PHONY: all build vet test lint clean

all: build vet test

build:
	go build ./...

vet:
	go vet ./...

test:
	go test ./... -count=1

test-race:
	go test -race ./... -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf tmp/
