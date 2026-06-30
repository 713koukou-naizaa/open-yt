# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Project variables
BINARY_NAME=open-yt
BINARY_DIR=./bin
CMD_PATH=./cmd/open-yt

.PHONY: all
all: build

.PHONY: build
build:
	@echo "==> Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "==> Binary available at $(BINARY_DIR)/$(BINARY_NAME)"

# Run application passing arguments to it
# Example: make run ARGS="search hello world"
.PHONY: run
run: build
	@$(BINARY_DIR)/$(BINARY_NAME) $(ARGS)

.PHONY: clean
clean:
	@echo "==> Cleaning..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

.PHONY: help
help:
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@echo '  all        Builds the application (default).'
	@echo '  build      Builds the application binary into the bin/ directory.'
	@echo '  run        Runs the application. Use ARGS to pass arguments. e.g., make run ARGS="search cats"'
	@echo '  clean      Removes the build artifacts.'
	@echo '  help       Shows this help message.'