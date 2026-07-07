@echo off
REM server.bat — запуск сервера управления сессиями ai-lives
REM Порт по умолчанию: 11000
REM Использование: server.bat [порт]

setlocal

set PORT=11000
if not "%~1"=="" set PORT=%~1

echo [server.bat] Запуск сервера на порту %PORT%...
echo [server.bat] Откройте http://127.0.0.1:%PORT% в браузере
echo [server.bat] Ctrl+C для остановки
echo.

python "%~dp0server\server.py" --port %PORT%

endlocal
