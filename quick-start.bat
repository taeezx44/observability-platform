@echo off
echo 🚀 Quick Start - Observability Platform
echo.

REM Check if Docker Desktop is running
echo 🐳 Checking Docker status...
docker info >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker Desktop is not running
    echo Please start Docker Desktop manually and wait for it to be ready
    echo Then run this script again
    echo.
    pause
    exit /b 1
)

echo ✅ Docker Desktop is running!
echo.

REM Start infrastructure
echo 📦 Starting ClickHouse and Kafka...
docker compose up clickhouse kafka -d

echo ⏳ Waiting 15 seconds for services to start...
timeout /t 15 /nobreak >nul

echo.
echo 🗃️ Running database migrations...
cat migrations\001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul
cat migrations\002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul  
cat migrations\003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul

echo.
echo 🌐 Starting all services...
docker compose up --build

echo.
echo 🎉 Platform started successfully!
echo.
echo 📊 Dashboard: http://localhost:3000
echo 🔍 API: http://localhost:8080
echo 🗄️ ClickHouse: http://localhost:8123
echo.
echo 🛑 To stop: docker compose down
echo.
pause
