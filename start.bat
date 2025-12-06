@echo off
REM Production startup script for Windows

REM Set production environment
set GIN_MODE=release
set REDIS_URL=localhost:6379

REM Build if needed
if not exist queue-server.exe (
    echo Building application...
    go build -o queue-server.exe .
)

REM Start the server
echo Starting Task Queue Server...
queue-server.exe

pause
