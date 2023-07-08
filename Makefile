BINARY_COMPILE=braces-compile
BINARY_VM=braces-vm
BINARY_INTROSPECT=braces-introspect
TEST?='./...'
GOFMT_FILES?=$$(find . -name '*.go')
GOLANGCI_LINT_VERSION?='v1.52.2'
IN_CI ?= false

.PHONY: build build-compile build-vm test lint format check-format vet clean tidy install-tools repl build-introspect

build: tidy build-compile build-vm build-introspect 

build-compile: build-proto
	@echo "Building $(BINARY_COMPILE)..."
	@go build -o target/$(BINARY_COMPILE) ./cmd/$(BINARY_COMPILE)

build-vm:
	@echo "Building $(BINARY_VM)..."
	@go build -o target/$(BINARY_VM) ./cmd/$(BINARY_VM)

build-introspect:
	@echo "Building $(BINARY_INTROSPECT)..."
	@go build -o target/$(BINARY_INTROSPECT) ./cmd/$(BINARY_INTROSPECT)

test:
	@echo "Running tests..."
	@if [ $(IN_CI) = false ]; then \
		gotestsum $(TEST); \
	else \
		go test $(TEST); \
	fi

lint:
	@echo "Running linters..."
	@golangci-lint run --tests=false --timeout=5m

staticcheck:
	@echo "Running staticcheck..."
	@staticcheck ./...

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

install-tools: install-linter install-gotestsum  install-staticcheck

install-linter:
	@if ! command -v golangci-lint &> /dev/null; then \
	  echo "Installing linter..."; \
	  go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi 

install-gotestsum:
	@if [ $(IN_CI) = false ]; then \
    if ! command -v gotestsum &> /dev/null; then \
		  echo "Installing gotestsum..."; \
		  go install gotest.tools/gotestsum; \
		fi \
	fi

install-staticcheck:
	@if ! command -v staticcheck &> /dev/null; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi

repl: build
	./target/braces-vm repl

introspect: build 
	./target/braces-introspect compiler

