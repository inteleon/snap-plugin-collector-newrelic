all: build

build: test
	@echo "Building..."
	@go build

test: deps
	@echo "Running tests..."
	@go test ./...

deps:
	@echo "Fetching dependencies..."
	@go get

install: build
	@echo "Installing..."
	@go install
