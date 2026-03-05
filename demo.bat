@echo off
REM Observability Platform Demo Script (Windows)
REM This script sets up and runs a complete demo of the observability platform

echo 🚀 Observability Platform Demo Setup
echo ====================================
echo.

REM Check prerequisites
echo 🔍 Checking prerequisites...
docker --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker is not installed. Please install Docker Desktop first.
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker Compose is not installed. Please install Docker Compose first.
    pause
    exit /b 1
)

go version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Go is not installed. Please install Go first.
    pause
    exit /b 1
)

npm --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Node.js/npm is not installed. Please install Node.js first.
    pause
    exit /b 1
)

echo ✅ All prerequisites are installed!
echo.

REM Clean up any existing containers
echo 🧹 Cleaning up existing containers...
docker-compose down -v >nul 2>&1
docker system prune -f >nul 2>&1

REM Start infrastructure
echo 📦 Starting infrastructure (ClickHouse + Kafka)...
docker-compose up -d clickhouse kafka

echo ⏳ Waiting for infrastructure to be ready...
timeout /t 15 /nobreak >nul

REM Check if ClickHouse is ready
echo 🔍 Checking ClickHouse connectivity...
:check_clickhouse
docker-compose exec -T clickhouse clickhouse-client --query "SELECT 1" >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo Waiting for ClickHouse...
    timeout /t 2 /nobreak >nul
    goto check_clickhouse
)

echo ✅ ClickHouse is ready!

REM Run database migrations
echo 🗃️ Running database migrations...
cat migrations\001_metrics.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability
cat migrations\002_logs.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability
cat migrations\003_traces.sql | docker-compose exec -T clickhouse clickhouse-client --database=observability

echo ✅ Database migrations completed!

REM Build and start services
echo 🔨 Building and starting services...
docker-compose build
docker-compose up -d

echo ⏳ Waiting for services to start...
timeout /t 20 /nobreak >nul

REM Check service health
echo 🔍 Checking service health...

REM Check API service
echo Checking API...
set /a attempt=1
:check_api
curl -f http://localhost:8080/health >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ API is ready!
    goto check_frontend
)
if %attempt% geq 30 (
    echo ❌ API failed to start within expected time
    echo 📋 Container logs:
    docker-compose logs api
    pause
    exit /b 1
)
echo ⏳ Waiting for API... (attempt %attempt%/30)
timeout /t 2 /nobreak >nul
set /a attempt+=1
goto check_api

:check_frontend
echo Checking Frontend...
set /a attempt=1
:check_frontend_loop
curl -f http://localhost:3000 >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ Frontend is ready!
    goto generate_data
)
if %attempt% geq 30 (
    echo ❌ Frontend failed to start within expected time
    echo 📋 Container logs:
    docker-compose logs frontend
    pause
    exit /b 1
)
echo ⏳ Waiting for Frontend... (attempt %attempt%/30)
timeout /t 2 /nobreak >nul
set /a attempt+=1
goto check_frontend_loop

:generate_data
REM Generate some demo data
echo 📊 Generating demo data...
timeout /t 5 /nobreak >nul

REM Generate metrics
echo Generating sample metrics...
for /l %%i in (1,1,10) do (
    curl -X POST http://localhost:8080/api/metrics -H "Content-Type: application/json" -d "{\"name\": \"cpu_usage\", \"value\": 75, \"labels\": {\"service\": \"demo-api\", \"instance\": \"localhost:8080\"}}" >nul 2>&1
    curl -X POST http://localhost:8080/api/metrics -H "Content-Type: application/json" -d "{\"name\": \"memory_usage\", \"value\": 60, \"labels\": {\"service\": \"demo-api\", \"instance\": \"localhost:8080\"}}" >nul 2>&1
    timeout /t 1 /nobreak >nul
)

REM Generate logs
echo Generating sample logs...
for /l %%i in (1,1,5) do (
    curl -X POST http://localhost:8080/api/logs -H "Content-Type: application/json" -d "{\"timestamp\": \"2024-01-15T10:30:00Z\", \"level\": \"INFO\", \"service\": \"demo-api\", \"message\": \"Processing request %%i\", \"fields\": {\"request_id\": \"req-%%i\", \"user_id\": \"user-%%i\"}}" >nul 2>&1
    timeout /t 1 /nobreak >nul
)

echo ✅ Demo data generated!

REM Display access information
echo.
echo 🎉 Observability Platform Demo is Ready!
echo ==========================================
echo.
echo 📊 Access Points:
echo   🌐 Frontend Dashboard: http://localhost:3000
echo   🔍 API Server:        http://localhost:8080
echo   📈 Metrics Endpoint:  http://localhost:8080/metrics
echo   🏥 Health Check:      http://localhost:8080/health
echo   📋 API Documentation: http://localhost:8080/swagger
echo.
echo 🗄️ Infrastructure:
echo   📊 ClickHouse:        http://localhost:8123
echo   🔄 Kafka:             localhost:9092
echo.
echo 📱 Demo Features:
echo   ✅ Real-time metrics collection
echo   ✅ Live log streaming
echo   ✅ Interactive dashboard
echo   ✅ Alert management
echo   ✅ Service topology
echo.
echo 🛠️ Management Commands:
echo   📋 View logs:         docker-compose logs -f
echo   🛑 Stop demo:         docker-compose down
echo   🧹 Clean up:          docker-compose down -v ^&^& docker system prune -f
echo.
echo 📸 Taking screenshot in 10 seconds...
timeout /t 10 /nobreak >nul

echo.
echo 🎯 Demo Tips:
echo   1. Open http://localhost:3000 in your browser
echo   2. Navigate through different sections (Metrics, Logs, Traces, Alerts)
echo   3. Watch real-time updates in the dashboard
echo   4. Try the filtering and search features
echo   5. Check the alert management panel
echo.
echo 📞 For support, check the logs: docker-compose logs -f
echo.
echo 🎉 Enjoy the demo!
echo.
echo Press any key to open the dashboard in your browser...
pause >nul

REM Open browser
start http://localhost:3000

echo.
echo 📸 Please manually capture a screenshot of the dashboard.
echo The demo will continue running until you stop it.
echo.
echo To stop the demo, run: docker-compose down
echo.
