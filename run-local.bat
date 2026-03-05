@echo off
echo 🚀 Running Observability Platform (Local Mode - No Docker)
echo.

echo 📋 Running in Local Mode:
echo   - Frontend: Local development server
echo   - API Server: Local with mock data
echo   - No ClickHouse/Kafka (Docker issues)
echo.

REM Start with frontend
echo 🌐 Starting Frontend Development Server...
cd frontend

REM Check if node_modules exists
if not exist "node_modules" (
    echo 📦 Installing dependencies...
    call npm install
)

echo 🚀 Starting frontend on http://localhost:3000
start "Frontend" cmd /k "npm run dev"

cd ..

echo.
echo 🔧 Starting API Server (Mock Mode)...
echo Using mock data since Docker is not available...

REM Create a simple mock API server
echo Creating mock API server...
echo.

REM Start a simple HTTP server for mock API
cd api
start "Mock API" cmd /k "python -m http.server 8080"
cd ..

echo.
echo 🎉 Local Platform Started!
echo.
echo 📊 Access Points:
echo   🌐 Frontend: http://localhost:3000
echo   🔍 Mock API: http://localhost:8080
echo.
echo 📝 Note: This is demo mode with mock data
echo    Real metrics/logs require Docker infrastructure
echo.
echo 🛑 To stop: Close the terminal windows
echo.
pause
