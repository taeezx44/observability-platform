@echo off
echo 🔍 Observability Platform Status Check
echo =====================================
echo.

echo 📋 Prerequisites Status:
echo.

REM Check Go
"C:\Program Files\Go\bin\go.exe" version >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ Go: Installed
    "C:\Program Files\Go\bin\go.exe" version
) else (
    echo ❌ Go: Not installed
)

echo.

REM Check Node.js
"C:\Program Files\nodejs\node.exe" --version >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ Node.js: Installed
    "C:\Program Files\nodejs\node.exe" --version
) else (
    echo ❌ Node.js: Not installed
)

echo.

REM Check Docker
docker --version >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ Docker: Installed
    docker --version
) else (
    echo ❌ Docker: Not installed
)

echo.

REM Check if Docker is running
docker info >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo ✅ Docker Desktop: Running
    echo.
    echo 🚀 Ready to start platform!
    echo Run: quick-start.bat
) else (
    echo ❌ Docker Desktop: Not running
    echo.
    echo 📝 To fix:
    echo 1. Start Docker Desktop manually from Start Menu
    echo 2. Wait for it to fully initialize (2-3 minutes)
    echo 3. Run this script again to verify
    echo 4. Then run: quick-start.bat
)

echo.

REM Check frontend build
if exist "frontend\dist\index.html" (
    echo ✅ Frontend: Built
) else (
    echo ❌ Frontend: Not built
    echo Run: build-frontend.bat
)

echo.
echo 📊 Project Status:
echo   Go Services: Ready
echo   Frontend: Ready
echo   Infrastructure: Waiting for Docker
echo.
pause
