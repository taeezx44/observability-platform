# Observability Platform Makefile
# Professional build and deployment automation

.PHONY: help build test docker-build docker-up docker-down benchmark demo clean install

# Default target
help:
	@echo "Observability Platform - Build & Deployment Commands"
	@echo ""
	@echo "🚀 Development:"
	@echo "  make install     - Install dependencies"
	@echo "  make build       - Build all services"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run platform locally"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Start all services"
	@echo "  make docker-down  - Stop all services"
	@echo "  make docker-logs  - View logs"
	@echo ""
	@echo "📊 Benchmarking:"
	@echo "  make benchmark   - Run performance tests"
	@echo "  make demo        - Start demo services"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make status      - Check system status"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	cd api && go mod tidy
	cd collector && go mod tidy
	cd alerting && go mod tidy
	cd demo/api && go mod tidy
	cd demo/worker && go mod tidy
	cd demo/database && go mod tidy
	cd simple-api && go mod tidy
	cd frontend && npm install

# Build all services
build:
	@echo "🔨 Building all services..."
	cd api/cmd && go build -o ../../bin/api .
	cd collector/cmd && go build -o ../../bin/collector .
	cd alerting && go build -o ../bin/alerting .
	cd simple-api && go build -o ../bin/api .
	cd frontend && npm run build

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
	cd api && go test ./...
	cd collector && go test ./...
	cd alerting && go test ./...

# Run platform locally
run:
	@echo "🚀 Starting observability platform locally..."
	./bin/api

# Build Docker images
docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build

# Start all services
docker-up:
	@echo "🚀 Starting all services..."
	docker-compose up -d

# Stop all services
docker-down:
	@echo "🛑 Stopping all services..."
	docker-compose down

# View Docker logs
docker-logs:
	docker-compose logs -f

# Run benchmark tests
benchmark:
	@echo "📊 Running benchmark tests..."
	./run-benchmark.bat

# Start demo services
demo:
	@echo "🎯 Starting demo services..."
	./start-demo.bat

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/dist/
	docker system prune -f

# Check system status
status:
	@echo "🔍 Checking system status..."
	./status-check.bat

# Quick start (install + build + run)
quick-start: install build run

# Full demo (install + build + docker-up)
full-demo: install docker-build docker-up
	@echo "🎉 Full demo started!"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔍 API: http://localhost:8080"
	@echo "📊 Metrics: http://localhost:8080/metrics"

# Production build
prod-build:
	@echo "🏭 Production build..."
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api ./api/cmd
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/collector ./collector/cmd
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/alerting ./alerting
	cd frontend && npm run build

# Development workflow
dev: install build test
	@echo "✅ Development setup complete!"

# CI/CD pipeline
ci: test docker-build
	@echo "✅ CI/CD pipeline complete!"
