# Makefile for Docker operations

# Stop and remove all containers
stop-all:
	@echo "Stopping and removing all Docker containers..."
	docker-compose down

# Remove all Docker containers, networks, and volumes
clean:
	@echo "Removing all Docker containers, networks, and volumes..."
	docker-compose down -v

# Build all Docker images
build:
	@echo "Building Docker images..."
	docker-compose build

# Start all containers in detached mode
up:
	@echo "Starting all containers in detached mode..."
	docker-compose up -d

# Show logs of all containers
logs:
	@echo "Showing logs of all containers..."
	docker-compose logs

# Show status of containers
status:
	@echo "Showing status of containers..."
	docker-compose ps

# Remove all unused Docker resources (use with caution)
prune:
	@echo "Removing all unused Docker resources..."
	docker system prune -af

# Default target
all: build up

.PHONY: stop-all clean build up logs status prune all
