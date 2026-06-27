@echo off
go build -o artisan.exe artisan.go
if %errorlevel% neq 0 exit /b %errorlevel%
echo Build successful: artisan.exe
