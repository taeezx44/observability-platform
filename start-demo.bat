@echo off
echo 🚀 Starting Observability Platform Demo
echo =====================================
echo.

echo 📋 Demo Services:
echo   - Demo API (port 8081)
echo   - Demo Worker (port 8082) 
echo   - Demo Database (port 8083)
echo   - Observability API (port 8080)
echo   - Frontend (port 3000)
echo   - ClickHouse (port 8123/9000)
echo   - Kafka (port 9092)
echo.

echo 🔧 Starting demo services...
echo.

REM Start demo services in background
echo Starting Demo API...
start "Demo API" cmd /k "cd demo/api && go run main.go"

echo Starting Demo Worker...
start "Demo Worker" cmd /k "cd demo/worker && go run main.go"

echo Starting Demo Database...
start "Demo Database" cmd /k "cd demo/database && go run main.go"

echo ⏳ Waiting for demo services to start...
timeout /t 5 /nobreak >nul

echo.
echo 🔍 Starting observability platform...
echo.

REM Start main platform
echo Starting Observability API...
start "Observability API" cmd /k "cd api/cmd && go run main.go"

echo ⏳ Waiting for API to start...
timeout /t 3 /nobreak >nul

echo.
echo 🌐 Starting frontend...
echo.

cd frontend
start "Frontend" cmd /k "npm run dev"

echo.
echo 🎉 Demo Platform Starting!
echo.
echo 📊 Access Points:
echo   🌐 Frontend: http://localhost:3000
echo   🔍 API: http://localhost:8080
echo   📈 Metrics: http://localhost:8080/metrics
echo   🔌 WebSocket: ws://localhost:8080/ws/logs
echo   🏥 Health: http://localhost:8080/health
echo.
echo 🎯 Demo Services:
echo   📡 Demo API: http://localhost:8081
echo   ⚙️ Demo Worker: http://localhost:8082
echo   🗄️ Demo Database: http://localhost:8083
echo.
echo 📊 Demo Metrics:
echo   📈 API Metrics: http://localhost:8081/metrics
echo   ⚙️ Worker Metrics: http://localhost:8082/metrics
echo   🗄️ DB Metrics: http://localhost:8083/metrics
echo.
echo 🎯 The platform will monitor all demo services!
echo.
echo 🛑 To stop: Close all terminal windows
echo.
pause
