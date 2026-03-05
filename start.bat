@echo off
echo 🚀 Starting Observability Platform...

REM 1. Start infrastructure
echo 📦 Starting ClickHouse and Kafka...
docker compose up clickhouse kafka -d

REM Wait for ClickHouse to be ready
echo ⏳ Waiting for ClickHouse...
timeout /t 10 /nobreak >nul

REM 2. Run migrations
echo 🗃️ Running database migrations...
cat migrations\001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations\002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations\003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability

REM 3. Start all services
echo 🌐 Starting all services...
docker compose up

echo ✅ Platform ready!
echo 📊 Dashboard: http://localhost:3000
echo 🔍 API: http://localhost:8080
echo 🗄️ ClickHouse: http://localhost:8123
pause
