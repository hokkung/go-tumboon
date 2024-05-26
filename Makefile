build:
	@echo "Building..."
	@go build ./...
	@echo "Finished Building..."

generate:
	@echo "Generaing code..."
	@go generate ./...
	@wire ./...
	@echo "Finished generaing code..."

run: all
	@go run ./cmd/make_permit_runner/main.go

test:
	@echo "Running test..."
	@go test ./...
	@echo "Tests completed..."

all: generate build test 

.PHONY: build generate run all test
