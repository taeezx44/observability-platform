@echo off
echo ⚛️ Building React Frontend...
echo.

cd frontend

REM Check if node_modules exists
if not exist "node_modules" (
    echo 📦 Installing dependencies...
    call npm install
    if %ERRORLEVEL% neq 0 (
        echo ❌ npm install failed
        pause
        exit /b 1
    )
)

echo 🔨 Building frontend...
call npm run build
if %ERRORLEVEL% neq 0 (
    echo ❌ Build failed
    pause
    exit /b 1
)

echo ✅ Frontend built successfully!
echo 📁 Build output: dist/
echo.

cd ..
pause
