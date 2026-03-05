#!/bin/bash

echo "🚀 Starting Observability Platform..."

# 1. Start infrastructure
echo "📦 Starting ClickHouse and Kafka..."
docker compose up clickhouse kafka -d

# Wait for ClickHouse to be ready
echo "⏳ Waiting for ClickHouse..."
sleep 10

# 2. Run migrations
echo "🗃️ Running database migrations..."
cat migrations/001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations/002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations/003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability

# 3. Start all services
echo "🌐 Starting all services..."
docker compose up

echo "✅ Platform ready!"
echo "📊 Dashboard: http://localhost:3000"
echo "🔍 API: http://localhost:8080"
echo "🗄️ ClickHouse: http://localhost:8123"
