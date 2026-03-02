.PHONY: help build clean test dev-backend dev-frontend cloudron-build cloudron-install

# Default target
help:
	@echo "ISOMan - Linux ISO Download Manager"
	@echo ""
	@echo "Available targets:"
	@echo "  make dev-backend       - Run backend in development mode"
	@echo "  make dev-frontend      - Run frontend in development mode"
	@echo "  make build            - Build frontend and backend"
	@echo "  make cloudron-build   - Build Cloudron package"
	@echo "  make cloudron-install - Install package to Cloudron (requires CLOUDRON_LOCATION)"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make test             - Run tests"

# Development
dev-backend:
	@echo "Starting backend server..."
	cd backend && go run main.go

dev-frontend:
	@echo "Starting frontend dev server..."
	cd ui && bun run dev

# Build
build: build-frontend build-backend

build-frontend:
	@echo "Building frontend..."
	cd ui && bun install && bun run build

build-backend:
	@echo "Building backend..."
	cd backend && go build -o server .

# Cloudron
cloudron-build:
	@echo "Building Cloudron package..."
	cloudron build

cloudron-install:
	@echo "Installing Cloudron package..."
	@test -n "$(CLOUDRON_LOCATION)" || (echo "CLOUDRON_LOCATION is required (example: isoman.example.com)" && exit 1)
	cloudron install --location $(CLOUDRON_LOCATION)

# Testing
test:
	@echo "Running backend tests..."
	cd backend && go test ./...

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ui/dist
	rm -f backend/server
	rm -rf backend/data
