@echo off
REM server.bat — запуск сервера управления сессиями ai-lives
REM Порт по умолчанию: 11000
REM Использование: server.bat [порт]

setlocal

set PYTHONUTF8=1
set PYTHONIOENCODING=utf-8
chcp 65001 >nul

set PORT=11000
if not "%~1"=="" set PORT=%~1

echo [server.bat] Запуск сервера на порту %PORT%...
echo [server.bat] Откройте http://127.0.0.1:%PORT% в браузере
echo [server.bat] Ctrl+C для остановки
echo.

cd /d "%~dp0"
uv run python "server\server.py" --port %PORT%

endlocal
