# Observability Platform Makefile

.PHONY: help build test docker-build docker-up docker-down benchmark demo clean install start start-demo migrate setup-docker setup-go status-check

# Default target
help:
	@echo "Observability Platform - Build & Deployment Commands"
	@echo ""
	@echo "🚀 Development:"
	@echo "  make install     - Install dependencies"
	@echo "  make build       - Build all services"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run platform locally"
	@echo "  make start       - Start platform (Docker)"
	@echo "  make start-demo  - Start demo services"
	@echo ""
	@echo "🔧 Setup:"
	@echo "  make setup-docker - Install Docker"
	@echo "  make setup-go    - Install Go"
	@echo "  make migrate     - Run database migrations"
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

# Install Docker
setup-docker:
	@echo "🐳 Installing Docker Desktop..."
	@if command -v docker >/dev/null 2>&1; then \
		echo "✅ Docker is already installed:"; \
		docker --version; \
	else \
		echo "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop/"; \
		echo "After installation, restart your terminal and run: docker --version"; \
	fi

# Install Go
setup-go:
	@echo "🚀 Installing Go 1.21..."
	@if command -v go >/dev/null 2>&1; then \
		echo "✅ Go is already installed:"; \
		go version; \
	else \
		echo "Please install Go from: https://go.dev/dl/"; \
		echo "After installation, restart your terminal and run: go version"; \
	fi

# Run database migrations
migrate:
	@echo "🗃️ Running database migrations..."
	cat migrations/001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability
	cat migrations/002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability
	cat migrations/003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability

# Start platform (replaces start.bat)
start:
	@echo "🚀 Starting Observability Platform..."
	@echo "📦 Starting ClickHouse and Kafka..."
	docker compose up clickhouse kafka -d
	@echo "⏳ Waiting for ClickHouse..."
	sleep 10
	@echo "🗃️ Running database migrations..."
	$(MAKE) migrate
	@echo "🌐 Starting all services..."
	docker compose up
	@echo "✅ Platform ready!"
	@echo "📊 Dashboard: http://localhost:3000"
	@echo "🔍 API: http://localhost:8080"
	@echo "🗄️ ClickHouse: http://localhost:8123"

# Start demo services (replaces start-demo.bat)
start-demo:
	@echo "🚀 Starting Observability Platform Demo"
	@echo "====================================="
	@echo ""
	@echo "📋 Demo Services:"
	@echo "  - Demo API (port 8081)"
	@echo "  - Demo Worker (port 8082)"
	@echo "  - Demo Database (port 8083)"
	@echo "  - Observability API (port 8080)"
	@echo "  - Frontend (port 3000)"
	@echo "  - ClickHouse (port 8123/9000)"
	@echo "  - Kafka (port 9092)"
	@echo ""
	@echo "🔧 Starting demo services..."
	@echo "Starting Demo API..."
	cd demo/api && go run main.go &
	@echo "Starting Demo Worker..."
	cd demo/worker && go run main.go &
	@echo "Starting Demo Database..."
	cd demo/database && go run main.go &
	@echo "⏳ Waiting for demo services to start..."
	sleep 5
	@echo ""
	@echo "🔍 Starting observability platform..."
	@echo "Starting Observability API..."
	cd api/cmd && go run main.go &
	@echo "⏳ Waiting for API to start..."
	sleep 3
	@echo ""
	@echo "🌐 Starting frontend..."
	@echo ""
	cd frontend && npm run dev &
	@echo ""
	@echo "🎉 Demo Platform Starting!"
	@echo ""
	@echo "📊 Access Points:"
	@echo "  🌐 Frontend: http://localhost:3000"
	@echo "  🔍 API: http://localhost:8080"
	@echo "  📈 Metrics: http://localhost:8080/metrics"
	@echo "  🔌 WebSocket: ws://localhost:8080/ws/logs"
	@echo "  🏥 Health: http://localhost:8080/health"
	@echo ""
	@echo "🎯 Demo Services:"
	@echo "  📡 Demo API: http://localhost:8081"
	@echo "  ⚙️ Demo Worker: http://localhost:8082"
	@echo "  🗄️ Demo Database: http://localhost:8083"
	@echo ""
	@echo "📊 Demo Metrics:"
	@echo "  📈 API Metrics: http://localhost:8081/metrics"
	@echo "  ⚙️ Worker Metrics: http://localhost:8082/metrics"
	@echo "  🗄️ DB Metrics: http://localhost:8083/metrics"
	@echo ""
	@echo "🎯 The platform will monitor all demo services!"

# Check system status (replaces status-check.bat)
status-check:
	@echo "🔍 Checking system status..."
	@echo "Docker status:"
	@if command -v docker >/dev/null 2>&1; then \
		docker --version; \
		docker compose ps; \
	else \
		echo "❌ Docker not installed"; \
	fi
	@echo ""
	@echo "Go status:"
	@if command -v go >/dev/null 2>&1; then \
		go version; \
	else \
		echo "❌ Go not installed"; \
	fi
	@echo ""
	@echo "Node.js status:"
	@if command -v npm >/dev/null 2>&1; then \
		npm --version; \
	else \
		echo "❌ Node.js/npm not installed"; \
	fi

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
	@if [ -f "./run-benchmark.sh" ]; then \
		./run-benchmark.sh; \
	else \
		echo "Benchmark script not found. Creating simple benchmark..."; \
		go test -bench=. ./...; \
	fi

# Start demo services
demo:
	@echo "🎯 Starting demo services..."
	$(MAKE) start-demo

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/dist/
	docker system prune -f

# Deep clean (remove junk files)
clean-deep:
	@echo "🧹 Deep cleaning project..."
	$(MAKE) clean
	@echo "Removing temporary files..."
	find . -name "*.tmp" -delete 2>/dev/null || true
	find . -name "*.log" -delete 2>/dev/null || true
	find . -name ".DS_Store" -delete 2>/dev/null || true
	find . -name "Thumbs.db" -delete 2>/dev/null || true
	@echo "Removing backup files..."
	find . -name "*.bak" -delete 2>/dev/null || true
	find . -name "*~" -delete 2>/dev/null || true
	@echo "Deep clean completed!"

# Check system status
status:
	@echo "🔍 Checking system status..."
	$(MAKE) status-check

# Quick start (install + build + run)
quick-start: install build run

# Full demo (install + build + docker-up)
full-demo: install docker-build docker-up
	@echo "🎉 Full demo started!"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔍 API: http://localhost:8080"
	@echo "📊 Metrics: http://localhost:8080/metrics"

# Interactive demo (with data generation)
demo-full:
	@echo "🚀 Starting interactive demo..."
	@if [ -f "./demo.sh" ]; then \
		chmod +x demo.sh && ./demo.sh; \
	elif [ -f "./demo.bat" ]; then \
		./demo.bat; \
	else \
		echo "Demo script not found. Running basic demo..."; \
		$(MAKE) full-demo; \
	fi

# Quick demo (just start services)
demo-quick:
	@echo "⚡ Starting quick demo..."
	docker-compose up -d
	@echo "🎉 Quick demo started!"
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
