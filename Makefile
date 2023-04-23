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

.PHONY: all test lint clean $(PACKAGES)
