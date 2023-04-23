# Makefile for the compiler package

# Variables
EXECUTABLES = braces-compiler 
LINTER := golangci-lint

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
	@$(LINTER) run 

# Clean up
clean:
	rm -rf ./target/bin/$(notdir $(CURDIR))

tidy:
	@for pkg in $(PACKAGES); do \
		echo "Tidying $$pkg..."; \
		cd $$pkg && go mod tidy && cd ..; \
	done

.PHONY: build test lint clean
