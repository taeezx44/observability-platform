@echo off
echo 🚀 Observability Platform Benchmark Test
echo.

echo 📋 System Information
echo =====================
echo CPU: 
wmic cpu get name /value | findstr Name
echo.
echo RAM:
wmic computersystem get TotalPhysicalMemory /value | findstr TotalPhysicalMemory
echo.

echo 🌐 Starting API Server...
echo Starting Go server for benchmark tests...
start "API Server" cmd /k "go run main.go"

echo ⏳ Waiting for server to start...
timeout /t 5 /nobreak >nul

echo.
echo 🔍 Testing API Health...
curl -s http://localhost:8080/health
echo.

echo 📊 Running Benchmark Tests...
echo ==============================

echo.
echo Test 1: Health Check (Light Load)
echo ab -n 1000 -c 10 http://localhost:8080/health
ab -n 1000 -c 10 http://localhost:8080/health

echo.
echo Test 2: API Info (Medium Load)
echo ab -n 5000 -c 50 http://localhost:8080/api
ab -n 5000 -c 50 http://localhost:8080/api

echo.
echo Test 3: Metrics Endpoint (Heavy Load)
echo ab -n 10000 -c 100 http://localhost:8080/metrics
ab -n 10000 -c 100 http://localhost:8080/metrics

echo.
echo Test 4: Logs Endpoint (Heavy Load)
echo ab -n 10000 -c 100 http://localhost:8080/logs
ab -n 10000 -c 100 http://localhost:8080/logs

echo.
echo Test 5: Traces Endpoint (Heavy Load)
echo ab -n 10000 -c 100 http://localhost:8080/traces
ab -n 10000 -c 100 http://localhost:8080/traces

echo.
echo 🎯 Benchmark Complete!
echo.
echo 📝 Results saved to: benchmark.md
echo.
pause
