BINARY_COMPILE=braces-compile
BINARY_VM=braces-vm
TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go')

.PHONY: build build-compile build-vm test lint format clean

build: build-compile build-vm

build-compile:
	@echo "Building $(BINARY_COMPILE)..."
	@go build -o target/$(BINARY_COMPILE) ./cmd/braces-compile

build-vm:
	@echo "Building $(BINARY_VM)..."
	@go build -o target/$(BINARY_VM) ./cmd/braces-vm

test:
	@echo "Running tests..."
	@go test -v $(TEST)

lint:
	@echo "Running linters..."
	@golangci-lint run 

format:
	@echo "Formatting code..."
	@gofmt -w $(GOFMT_FILES)

vet:
	@echo "Vetting code..."
	@go vet ./... 

clean:
	@echo "Cleaning up..."
	@rm -rf target/*

