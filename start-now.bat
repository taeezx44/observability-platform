@echo off
echo 🚀 Starting Observability Platform (WSL Fixed!)
echo.

REM Check Docker status
echo 🐳 Checking Docker Desktop...
docker info >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker Desktop is not running
    echo.
    echo 📝 Please start Docker Desktop manually:
    echo 1. Open Start Menu
    echo 2. Search for "Docker Desktop"
    echo 3. Click to open it
    echo 4. Wait 2-3 minutes for initialization
    echo 5. Run this script again
    echo.
    pause
    exit /b 1
)

echo ✅ Docker Desktop is running!
echo.

REM Start infrastructure
echo 📦 Starting ClickHouse and Kafka...
docker compose up clickhouse kafka -d

echo ⏳ Waiting 20 seconds for services to initialize...
timeout /t 20 /nobreak >nul

echo.
echo 🗃️ Running database migrations...
cat migrations\001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul
if %ERRORLEVEL% equ 0 echo ✅ Metrics migration completed
cat migrations\002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul
if %ERRORLEVEL% equ 0 echo ✅ Logs migration completed
cat migrations\003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability 2>nul
if %ERRORLEVEL% equ 0 echo ✅ Traces migration completed

echo.
echo 🌐 Starting all services...
echo This will build and run:
echo   - Collector (metrics scraper)
echo   - API (REST server)
echo   - Alerting (rule engine)
echo   - Frontend (React dashboard)
echo.

docker compose up --build

echo.
echo 🎉 Observability Platform is running!
echo.
echo 📊 Access Points:
echo   🌐 Dashboard: http://localhost:3000
echo   🔍 API: http://localhost:8080
echo   🗄️ ClickHouse: http://localhost:8123
echo.
echo 🛑 To stop: docker compose down
echo.
pause
