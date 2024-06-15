# Define variables for commands and paths for readability and maintainability
MOCKGEN=mockgen
GO=go
SOURCE=/Users/louiskwezi/Desktop/logger/back-end/models/models.go
DESTINATION=/Users/louiskwezi/Desktop/logger/back-end/mocks
PACKAGE=mocks

# Default target
all: generate test

# Generate mocks
generate:
	$(MOCKGEN) -source=$(SOURCE) -destination=$(DESTINATION) -package=$(PACKAGE)

# Run tests
test:
	$(GO)test ./...

# Clean up generated files
clean:
	rm -f$(DESTINATION)
run:
	go run main.go
# Mark targets not associated with files
.PHONY: all generate test clean

up:
	docker-compose -f /Users/louiskwezi/Desktop/device_logger/database/docker-compose.yml.yaml up -d

down:
	docker-compose -f /Users/louiskwezi/Desktop/device_logger/database/docker-compose.yml.yaml down

status:
	docker-compose -f /Users/louiskwezi/Desktop/device_logger/database/docker-compose.yml.yaml ps

.PHONY: up down