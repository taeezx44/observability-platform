@echo off
echo 🔄 Restarting Docker Desktop...
echo.

echo 1️⃣ Stopping Docker Desktop...
taskkill /f /im "Docker Desktop.exe" 2>nul

echo 2️⃣ Waiting for complete shutdown...
timeout /t 10 /nobreak >nul

echo 3️⃣ Starting Docker Desktop...
start "" "Docker Desktop"

echo 4️⃣ Waiting for initialization...
echo This may take 2-3 minutes...

set DOCKER_READY=0
for /l %%i in (1,1,36) do (
    docker info >nul 2>&1
    if !ERRORLEVEL! equ 0 (
        set DOCKER_READY=1
        echo ✅ Docker Desktop is ready!
        goto :docker_ready
    )
    echo Waiting... %%i/36 (30 seconds each)
    timeout /t 30 /nobreak >nul
)

:docker_ready
if !DOCKER_READY! equ 0 (
    echo ⚠️ Docker Desktop may still be starting
    echo Please check manually and run start-now.bat
) else (
    echo 🚀 Docker is ready! Running platform...
    call start-now.bat
)

pause
