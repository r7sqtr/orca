.PHONY: build run clean test fmt lint install uninstall

BINARY := orca
VERSION := 0.1.0
PREFIX := /usr/local

build:
	go build -ldflags "-s -w" -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	go clean

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

install: build
	install -m 755 $(BINARY) $(PREFIX)/bin/$(BINARY)

uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

.DEFAULT_GOAL := build
