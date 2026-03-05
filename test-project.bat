@echo off
echo 🔍 Testing Observability Platform Project Structure...
echo.

REM Check if all required directories exist
echo 📁 Checking project structure...
if not exist "collector" echo ❌ Missing collector directory & pause & exit /b 1
if not exist "api" echo ❌ Missing api directory & pause & exit /b 1  
if not exist "alerting" echo ❌ Missing alerting directory & pause & exit /b 1
if not exist "frontend" echo ❌ Missing frontend directory & pause & exit /b 1
if not exist "migrations" echo ❌ Missing migrations directory & pause & exit /b 1
echo ✅ All directories found

REM Check key files
echo.
echo 📄 Checking key files...
if not exist "go.mod" echo ❌ Missing go.mod & pause & exit /b 1
if not exist "docker-compose.yml" echo ❌ Missing docker-compose.yml & pause & exit /b 1
if not exist "README.md" echo ❌ Missing README.md & pause & exit /b 1
echo ✅ Key files found

REM Check collector structure
echo.
echo 🔧 Checking collector structure...
if not exist "collector\cmd\main.go" echo ❌ Missing collector main.go & pause & exit /b 1
if not exist "collector\scraper\scraper.go" echo ❌ Missing scraper.go & pause & exit /b 1
if not exist "collector\storage\clickhouse.go" echo ❌ Missing storage clickhouse.go & pause & exit /b 1
if not exist "collector\logger\parser.go" echo ❌ Missing logger parser.go & pause & exit /b 1
if not exist "collector\tracer\span.go" echo ❌ Missing tracer span.go & pause & exit /b 1
echo ✅ Collector structure OK

REM Check API structure
echo.
echo 🌐 Checking API structure...
if not exist "api\cmd\main.go" echo ❌ Missing API main.go & pause & exit /b 1
if not exist "api\handlers\metrics.go" echo ❌ Missing metrics handler & pause & exit /b 1
if not exist "api\handlers\logs.go" echo ❌ Missing logs handler & pause & exit /b 1
if not exist "api\handlers\traces.go" echo ❌ Missing traces handler & pause & exit /b 1
echo ✅ API structure OK

REM Check alerting structure
echo.
echo 🚨 Checking alerting structure...
if not exist "alerting\main.go" echo ❌ Missing alerting main.go & pause & exit /b 1
if not exist "alerting\engine.go" echo ❌ Missing alerting engine & pause & exit /b 1
if not exist "alerting\rules.yaml" echo ❌ Missing alerting rules & pause & exit /b 1
echo ✅ Alerting structure OK

REM Check frontend structure
echo.
echo ⚛️ Checking frontend structure...
if not exist "frontend\package.json" echo ❌ Missing package.json & pause & exit /b 1
if not exist "frontend\src\App.jsx" echo ❌ Missing App.jsx & pause & exit /b 1
if not exist "frontend\src\pages\Dashboard.jsx" echo ❌ Missing Dashboard.jsx & pause & exit /b 1
if not exist "frontend\src\pages\Logs.jsx" echo ❌ Missing Logs.jsx & pause & exit /b 1
if not exist "frontend\src\pages\Traces.jsx" echo ❌ Missing Traces.jsx & pause & exit /b 1
if not exist "frontend\index.html" echo ❌ Missing index.html & pause & exit /b 1
echo ✅ Frontend structure OK

REM Check migrations
echo.
echo 🗃️ Checking migrations...
if not exist "migrations\001_metrics.sql" echo ❌ Missing metrics migration & pause & exit /b 1
if not exist "migrations\002_logs.sql" echo ❌ Missing logs migration & pause & exit /b 1
if not exist "migrations\003_traces.sql" echo ❌ Missing traces migration & pause & exit /b 1
echo ✅ Migrations OK

REM Check Docker files
echo.
echo 🐳 Checking Docker files...
if not exist "collector\Dockerfile" echo ❌ Missing collector Dockerfile & pause & exit /b 1
if not exist "api\Dockerfile" echo ❌ Missing API Dockerfile & pause & exit /b 1
if not exist "alerting\Dockerfile" echo ❌ Missing alerting Dockerfile & pause & exit /b 1
if not exist "frontend\Dockerfile" echo ❌ Missing frontend Dockerfile & pause & exit /b 1
echo ✅ Docker files OK

echo.
echo 🎉 PROJECT STRUCTURE TEST PASSED!
echo.
echo 📋 Next steps to run the platform:
echo 1. Install Docker Desktop
echo 2. Install Go 1.21+
echo 3. Install Node.js 18+
echo 4. Run: start.bat
echo.
echo 🌐 Access points when running:
echo - Dashboard: http://localhost:3000
echo - API: http://localhost:8080
echo - ClickHouse: http://localhost:8123
echo.
pause
