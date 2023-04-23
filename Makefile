# Main Makefile for the Go-Braces project

# Variables
PACKAGES = compiler vm repl debugger

# Build all packages
all: $(PACKAGES)

# Build individual packages
$(PACKAGES):
	$(MAKE) -C $@ build

# Run tests for all packages
test:
	for pkg in $(PACKAGES); do \
		$(MAKE) -C $$pkg test; \
	done

# Lint all packages
lint:
	for pkg in $(PACKAGES); do \
		$(MAKE) -C $$pkg lint; \
	done

# Clean up
clean:
	rm -rf bin

tidy:
	@for pkg in $(PACKAGES); do \
		echo "Tidying $$pkg..."; \
		cd $$pkg && go mod tidy && cd ..; \
	done

install-tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2



.PHONY: all test lint clean tidy install-tools $(PACKAGES)
