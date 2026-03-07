@echo off
if "%1"=="up" (
    docker compose -f docker-compose.yaml up -d
    echo Database is starting, waiting...
    timeout /t 5 > nul
) else if "%1"=="down" (
    docker compose -f docker-compose.yaml down -v
) else if "%1"=="logs" (
    docker compose -f docker-compose.yaml logs -f postgres
) else (
    echo Usage: db.bat [up^|down^|logs]
)
