@echo off
echo 🐳 Starting Observability Platform with Docker
echo =============================================
echo.

echo 📋 Services to start:
echo   - ClickHouse Database (ports 8123, 9000)
echo   - Kafka (port 9092)
echo   - Collector Service
echo   - API Service (port 8080)
echo   - Alerting Service
echo   - Frontend (port 3000)
echo.

echo 🔧 Checking Docker...
docker --version
if %ERRORLEVEL% neq 0 (
    echo ❌ Docker not found! Please install Docker Desktop first.
    pause
    exit /b 1
)

echo.
echo 🐳 Docker Compose found!
docker-compose --version

echo.
echo 🛑 Stopping any existing containers...
docker-compose down

echo.
echo 🏗️ Building Docker images...
docker-compose build

echo.
echo 🚀 Starting all services...
docker-compose up -d

echo.
echo ⏳ Waiting for services to start...
timeout /t 10 /nobreak >nul

echo.
echo 📊 Checking service status...
docker-compose ps

echo.
echo 🌐 Access Points:
echo   🌐 Frontend: http://localhost:3000
echo   🔍 API: http://localhost:8080
echo   📈 Metrics: http://localhost:8080/metrics
echo   🏥 Health: http://localhost:8080/health
echo   🗄️ ClickHouse: http://localhost:8123
echo   📡 Kafka: localhost:9092
echo.

echo 📊 Additional Endpoints:
echo   🔌 WebSocket: ws://localhost:8080/ws/logs
echo   📋 API Info: http://localhost:8080/api
echo.

echo 🎯 Demo Services (if using demo-compose):
echo   📡 Demo API: http://localhost:8081
echo   ⚙️ Demo Worker: http://localhost:8082
echo   🗄️ Demo Database: http://localhost:8083
echo.

echo 🔍 Checking logs...
echo.
echo 📊 API Logs:
docker-compose logs api | tail -10

echo.
echo 🗄️ ClickHouse Logs:
docker-compose logs clickhouse | tail -5

echo.
echo 🎉 Observability Platform is running!
echo.
echo 🛑 To stop: docker-compose down
echo 📊 To view logs: docker-compose logs -f
echo 🔄 To restart: docker-compose restart
echo.
pause
