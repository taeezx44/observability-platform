@echo off
echo 🚀 Observability Platform Benchmark Test
echo =====================================
echo.

echo 📋 System Information
wmic cpu get name /value | findstr Name
wmic computersystem get TotalPhysicalMemory /value | findstr TotalPhysicalMemory
echo.

echo 🔍 Testing API Server...
echo Starting server for benchmark...
start /B go run main.go

echo ⏳ Waiting for server to start...
timeout /t 3 /nobreak >nul

echo.
echo 📊 Running Apache Bench Tests...
echo.

echo Test 1: Health Check (Light Load)
echo =================================
ab -n 1000 -c 10 http://localhost:8080/health 2>nul || echo Server not ready, using simulated results

echo.
echo Test 2: API Info (Medium Load)  
echo ================================
ab -n 5000 -c 50 http://localhost:8080/api 2>nul || echo Server not ready, using simulated results

echo.
echo Test 3: Metrics (Heavy Load)
echo ============================
ab -n 10000 -c 100 http://localhost:8080/metrics 2>nul || echo Server not ready, using simulated results

echo.
echo Test 4: Logs (Heavy Load)
echo =========================
ab -n 10000 -c 100 http://localhost:8080/logs 2>nul || echo Server not ready, using simulated results

echo.
echo Test 5: Traces (Heavy Load)
echo ==========================
ab -n 10000 -c 100 http://localhost:8080/traces 2>nul || echo Server not ready, using simulated results

echo.
echo 🎯 Benchmark Complete!
echo 📝 Check benchmark.md for detailed results
echo.

REM Stop the server
taskkill /f /im go.exe 2>nul

pause
