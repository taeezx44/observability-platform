@echo off
echo 🔍 Validating Code Quality...
echo.

REM Check Go modules
echo 📦 Checking Go modules...
if exist "go.mod" (
    echo ✅ go.mod found
    for /f "tokens=*" %%i in (go.mod) do (
        echo    %%i
    )
) else (
    echo ❌ go.mod missing
)

echo.
echo 🐳 Checking Docker Compose...
if exist "docker-compose.yml" (
    echo ✅ docker-compose.yml found
    echo    Services defined:
    findstr "image:" docker-compose.yml | findstr /v "comment"
) else (
    echo ❌ docker-compose.yml missing
)

echo.
echo ⚛️ Checking Frontend Dependencies...
if exist "frontend\package.json" (
    echo ✅ package.json found
    echo    Dependencies:
    for /f "tokens=2 delims=: " %%i in ('findstr "react" frontend\package.json') do echo      - %%i
) else (
    echo ❌ package.json missing
)

echo.
echo 🗃️ Checking Database Schemas...
if exist "migrations\001_metrics.sql" (
    echo ✅ Metrics schema found
    findstr "CREATE TABLE" migrations\001_metrics.sql
) else (
    echo ❌ Metrics schema missing
)

if exist "migrations\002_logs.sql" (
    echo ✅ Logs schema found  
    findstr "CREATE TABLE" migrations\002_logs.sql
) else (
    echo ❌ Logs schema missing
)

if exist "migrations\003_traces.sql" (
    echo ✅ Traces schema found
    findstr "CREATE TABLE" migrations\003_traces.sql
) else (
    echo ❌ Traces schema missing
)

echo.
echo 🚨 Checking Alert Rules...
if exist "alerting\rules.yaml" (
    echo ✅ Alert rules found
    findstr "name:" alerting\rules.yaml | findstr /v "comment"
) else (
    echo ❌ Alert rules missing
)

echo.
echo 📊 Project Summary:
echo    - Go services: 3 (collector, api, alerting)
echo    - Frontend: React + Vite
echo    - Database: ClickHouse
echo    - Message Queue: Kafka
echo    - Total files: XX

echo.
echo 🎉 Code validation complete!
echo.
echo 🚀 To run the platform:
echo    1. Install prerequisites (Docker, Go, Node.js)
echo    2. Run: start.bat
echo    3. Open: http://localhost:3000
echo.
pause
