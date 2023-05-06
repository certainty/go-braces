BINARY_COMPILE=braces-compile
BINARY_VM=braces-vm
TEST?='./...'
GOFMT_FILES?=$$(find . -name '*.go')
GOLANGCI_LINT_VERSION?='v1.52.2'

.PHONY: build build-compile build-vm test lint format check-format vet clean tidy install-tools install-proto-tools repl build-proto

build: tidy build-compile build-vm

build-compile:
	@echo "Building $(BINARY_COMPILE)..."
	@go build -o target/$(BINARY_COMPILE) ./cmd/braces-compile

build-vm:
	@echo "Building $(BINARY_VM)..."
	@go build -o target/$(BINARY_VM) ./cmd/braces-vm

build-proto:
	@echo "Building grpc services..."
	@protoc --go_out=./internal/introspection/service --go-grpc_out=./internal/introspection/service ./internal/introspection/service/proto/compiler_introspection.proto

test:
	@echo "Running tests..."
	@gotestsum $(TEST)

lint:
	@echo "Running linters..."
	@golangci-lint run --tests=false --timeout=5m

format:
	@echo "Formatting code..."
	@gofmt -w $(GOFMT_FILES)

check-format:
	@echo "Checking code format..."
	@gofmt -d . |grep -q '^' && (echo "Not formatted correctly"; exit 1) || exit 0

vet:
	@echo "Vetting code..."
	@go vet ./... 

clean:
	@echo "Cleaning up..."
	@rm -rf target/*

tidy:
	@echo "Tidying up..."
	@go mod tidy

install-tools:
	@echo "Installing tools..."
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go get gotest.tools/gotestsum  
	@go install gotest.tools/gotestsum  

install-proto-tools:
	@echo "Installing protobuf tools..."
	@brew install protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

repl: build
	./target/braces-vm repl

