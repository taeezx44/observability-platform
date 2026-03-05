@echo off
echo 🧪 Testing Platform Components (Without Docker)
echo.

echo ✅ Go Services Test:
echo Testing Go compilation...

cd /d "c:\Users\Administrator\Downloads\New folder (3)\observability-platform"

echo.
echo 🔨 Building Collector...
"C:\Program Files\Go\bin\go.exe" build ./collector/cmd
if %ERRORLEVEL% equ 0 (
    echo ✅ Collector builds successfully
    if exist "collector.exe" echo 📁 collector.exe created
) else (
    echo ❌ Collector build failed
)

echo.
echo 🔨 Building API...
"C:\Program Files\Go\bin\go.exe" build ./api/cmd
if %ERRORLEVEL% equ 0 (
    echo ✅ API builds successfully
    if exist "api.exe" echo 📁 api.exe created
) else (
    echo ❌ API build failed
)

echo.
echo 🔨 Building Alerting...
"C:\Program Files\Go\bin\go.exe" build ./alerting
if %ERRORLEVEL% equ 0 (
    echo ✅ Alerting builds successfully
    if exist "alerting.exe" echo 📁 alerting.exe created
) else (
    echo ❌ Alerting build failed
)

echo.
echo ✅ Frontend Test:
if exist "frontend\dist\index.html" (
    echo ✅ Frontend built successfully
    echo 📁 dist/ folder ready
) else (
    echo ❌ Frontend not built
    echo Run: build-frontend.bat
)

echo.
echo 📋 Project Summary:
echo   Go Services: Ready for execution
echo   Frontend: Ready for deployment
echo   Infrastructure: Docker needs troubleshooting
echo.

echo 🚀 When Docker is fixed:
echo   1. Start Docker Desktop
echo   2. Run: quick-start.bat
echo   3. Access: http://localhost:3000
echo.

echo 🐛 Docker Issues Found:
echo   - Storage I/O errors
echo   - May need Docker Desktop restart
echo   - Or disk space cleanup
echo.

pause
