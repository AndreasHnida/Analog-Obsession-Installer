@echo off
go build -ldflags="-H windowsgui -s -w" -o AOInstaller.exe .
if %errorlevel% neq 0 (
    echo Build failed.
    exit /b %errorlevel%
)
echo Build successful: AOInstaller.exe
