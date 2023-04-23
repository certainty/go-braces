# Makefile for the repl package

# Variables
EXECUTABLES = braces-repl 

# Build the package
build:
	for exe in $(EXECUTABLES); do \
		go build -o ./target/bin/$(notdir $(CURDIR))/$$exe ./cmd/$$exe; \
	done

# Run tests
test:
	go test -v ./...

# Lint the package
lint:
	golint ./...

# Clean up
clean:
	rm -rf ./target/bin/$(notdir $(CURDIR))

.PHONY: build test lint clean
