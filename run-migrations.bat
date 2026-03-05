@echo off
echo 🗃️ Running database migrations...
echo.

REM Wait for ClickHouse to be ready
echo ⏳ Waiting for ClickHouse...
timeout /t 10 /nobreak >nul

REM Run migrations using PowerShell
echo 📊 Running metrics migration...
powershell -Command "Get-Content migrations\001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability"

echo 📝 Running logs migration...
powershell -Command "Get-Content migrations\002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability"

echo 🔍 Running traces migration...
powershell -Command "Get-Content migrations\003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability"

echo.
echo ✅ All migrations completed!
echo.
pause
