@echo off
echo 🚀 Starting Observability Platform...
echo.

REM Check prerequisites
echo 🔍 Checking prerequisites...

REM Check Go
"C:\Program Files\Go\bin\go.exe" version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Go not found. Please install Go 1.21+
    pause
    exit /b 1
)
echo ✅ Go installed

REM Check Node.js
"C:\Program Files\nodejs\node.exe" --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Node.js not found. Please install Node.js 18+
    pause
    exit /b 1
)
echo ✅ Node.js installed

REM Check Docker
docker --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker not found. Please install Docker Desktop
    pause
    exit /b 1
)
echo ✅ Docker installed

echo.
echo 🐳 Starting Docker Desktop...
echo Please wait for Docker to fully start (this may take 2-3 minutes)
echo.

REM Try to start Docker Desktop
start "" "Docker Desktop" 2>nul

REM Wait for Docker to be ready
echo ⏳ Waiting for Docker to be ready...
set DOCKER_READY=0
for /l %%i in (1,1,30) do (
    docker info >nul 2>&1
    if !ERRORLEVEL! equ 0 (
        set DOCKER_READY=1
        echo ✅ Docker is ready!
        goto :docker_ready
    )
    echo Waiting... %%i/30
    timeout /t 5 /nobreak >nul
)

:docker_ready
if !DOCKER_READY! equ 0 (
    echo ⚠️ Docker may still be starting. Please wait and try again.
    echo You can also start Docker Desktop manually and run this script again.
    pause
    exit /b 1
)

echo.
echo 📦 Starting infrastructure services...
docker compose up clickhouse kafka -d

echo ⏳ Waiting for ClickHouse to be ready...
timeout /t 15 /nobreak >nul

echo.
echo 🗃️ Running database migrations...
cat migrations\001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations\002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations\003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability

echo.
echo 🌐 Building and starting all services...
docker compose up --build

echo.
echo 🎉 Observability Platform is running!
echo.
echo 📊 Access Points:
echo   Dashboard: http://localhost:3000
echo   API: http://localhost:8080
echo   ClickHouse: http://localhost:8123
echo   Kafka: localhost:9092
echo.
echo 🛑 To stop: docker compose down
echo.
pause
