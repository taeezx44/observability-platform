#!/bin/bash

# Observability Platform Demo Script
# This script sets up and runs a complete demo of the observability platform

set -e

echo "🚀 Observability Platform Demo Setup"
echo "===================================="
echo ""

# Check prerequisites
echo "🔍 Checking prerequisites..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ Node.js/npm is not installed. Please install Node.js first."
    exit 1
fi

echo "✅ All prerequisites are installed!"
echo ""

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down -v 2>/dev/null || true
docker system prune -f

# Start infrastructure
echo "📦 Starting infrastructure (ClickHouse + Kafka)..."
docker-compose up -d clickhouse kafka

echo "⏳ Waiting for infrastructure to be ready..."
sleep 15

# Check if ClickHouse is ready
echo "🔍 Checking ClickHouse connectivity..."
until docker-compose exec -T clickhouse clickhouse-client --query "SELECT 1" &>/dev/null; do
    echo "Waiting for ClickHouse..."
    sleep 2
done

echo "✅ ClickHouse is ready!"

# Run database migrations
echo "🗃️ Running database migrations..."
cat migrations/001_metrics.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability
cat migrations/002_logs.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability  
cat migrations/003_traces.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability

echo "✅ Database migrations completed!"

# Build and start services
echo "🔨 Building and starting services..."
docker-compose build
docker-compose up -d

echo "⏳ Waiting for services to start..."
sleep 20

# Check service health
echo "🔍 Checking service health..."
services=("api:8080" "collector:9090" "alerting:9091" "frontend:3000")

for service in "${services[@]}"; do
    service_name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    
    echo "Checking $service_name..."
    max_attempts=30
    attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:$port/health &>/dev/null || \
           curl -f http://localhost:$port &>/dev/null; then
            echo "✅ $service_name is ready!"
            break
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            echo "❌ $service_name failed to start within expected time"
            echo "📋 Container logs:"
            docker-compose logs $service_name
            exit 1
        fi
        
        echo "⏳ Waiting for $service_name... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
done

# Generate some demo data
echo "📊 Generating demo data..."
sleep 5

# Generate metrics
echo "Generating sample metrics..."
for i in {1..10}; do
    curl -X POST http://localhost:8080/api/metrics \
         -H "Content-Type: application/json" \
         -d "{
           \"name\": \"cpu_usage\",
           \"value\": $((60 + RANDOM % 30)),
           \"labels\": {
             \"service\": \"demo-api\",
             \"instance\": \"localhost:8080\"
           }
         }" &>/dev/null || true
    
    curl -X POST http://localhost:8080/api/metrics \
         -H "Content-Type: application/json" \
         -d "{
           \"name\": \"memory_usage\",
           \"value\": $((40 + RANDOM % 40)),
           \"labels\": {
             \"service\": \"demo-api\",
             \"instance\": \"localhost:8080\"
           }
         }" &>/dev/null || true
    
    sleep 0.5
done

# Generate logs
echo "Generating sample logs..."
for i in {1..5}; do
    curl -X POST http://localhost:8080/api/logs \
         -H "Content-Type: application/json" \
         -d "{
           \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
           \"level\": \"INFO\",
           \"service\": \"demo-api\",
           \"message\": \"Processing request $i\",
           \"fields\": {
             \"request_id\": \"req-$i\",
             \"user_id\": \"user-$((RANDOM % 100))\"
           }
         }" &>/dev/null || true
    sleep 0.3
done

echo "✅ Demo data generated!"

# Display access information
echo ""
echo "🎉 Observability Platform Demo is Ready!"
echo "=========================================="
echo ""
echo "📊 Access Points:"
echo "  🌐 Frontend Dashboard: http://localhost:3000"
echo "  🔍 API Server:        http://localhost:8080"
echo "  📈 Metrics Endpoint:  http://localhost:8080/metrics"
echo "  🏥 Health Check:      http://localhost:8080/health"
echo "  📋 API Documentation: http://localhost:8080/swagger"
echo ""
echo "🗄️ Infrastructure:"
echo "  📊 ClickHouse:        http://localhost:8123"
echo "  🔄 Kafka:             localhost:9092"
echo ""
echo "📱 Demo Features:"
echo "  ✅ Real-time metrics collection"
echo "  ✅ Live log streaming"
echo "  ✅ Interactive dashboard"
echo "  ✅ Alert management"
echo "  ✅ Service topology"
echo ""
echo "🛠️ Management Commands:"
echo "  📋 View logs:         docker-compose logs -f"
echo "  🛑 Stop demo:         docker-compose down"
echo "  🧹 Clean up:          docker-compose down -v && docker system prune -f"
echo ""
echo "📸 Taking screenshot in 10 seconds..."
sleep 10

# Take a screenshot if available
if command -v screencapture &> /dev/null; then
    screencapture demo-screenshot.png
    echo "📸 Screenshot saved as demo-screenshot.png"
elif command -v gnome-screenshot &> /dev/null; then
    gnome-screenshot -f demo-screenshot.png
    echo "📸 Screenshot saved as demo-screenshot.png"
elif command -v import &> /dev/null; then
    import -window root demo-screenshot.png
    echo "📸 Screenshot saved as demo-screenshot.png"
else
    echo "📸 Screenshot tool not available. Please manually capture the dashboard."
fi

echo ""
echo "🎯 Demo Tips:"
echo "  1. Open http://localhost:3000 in your browser"
echo "  2. Navigate through different sections (Metrics, Logs, Traces, Alerts)"
echo "  3. Watch real-time updates in the dashboard"
echo "  4. Try the filtering and search features"
echo "  5. Check the alert management panel"
echo ""
echo "📞 For support, check the logs: docker-compose logs -f"
echo ""
echo "🎉 Enjoy the demo!"
