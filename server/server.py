#!/usr/bin/env python3
"""
Сервер управления сессиями ai-lives.

Запуск: uv run python server/server.py [--port 11000]
Endpoints:
  GET  /                    — веб-дашборд
  GET  /api/last-session    — содержимое state/last_session.md
  GET  /api/status          — статус сервера (idle/running/auto)
  GET  /api/events          — SSE-поток для real-time обновлений
  POST /api/session/start   — запуск одной сессии
  POST /api/auto/toggle     — переключение автосессии
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any

# Корень проекта — родительская директория относительно server/
PROJECT_ROOT = Path(__file__).resolve().parents[1]
STATE_DIR = PROJECT_ROOT / "state"
STATIC_DIR = Path(__file__).resolve().parent / "static"
SESSION_COMMAND = ["uv", "run", "python", "scripts/session_transaction.py"]


def utf8_subprocess_env() -> dict[str, str]:
    env = os.environ.copy()
    env["PYTHONUTF8"] = "1"
    env["PYTHONIOENCODING"] = "utf-8"
    return env

# Глобальное состояние
class ServerState:
    def __init__(self) -> None:
        self.lock = threading.Lock()
        self.status: str = "idle"  # idle | running | auto
        self.auto_enabled: bool = False
        self.auto_thread: threading.Thread | None = None
        self.last_error: str = ""
        self.last_run_time: str = ""
        self.session_count: int = 0
        self._sse_clients: list[Any] = []
        self._sse_lock = threading.Lock()

    def broadcast(self, event: str, data: dict) -> None:
        """Отправить SSE-событие всем подключённым клиентам."""
        msg = f"event: {event}\ndata: {json.dumps(data, ensure_ascii=False)}\n\n"
        with self._sse_lock:
            dead = []
            for client in self._sse_clients:
                try:
                    client.write(msg.encode("utf-8"))
                    client.flush()
                except Exception:
                    dead.append(client)
            for c in dead:
                self._sse_clients.remove(c)

    def add_sse_client(self, wfile: Any) -> None:
        with self._sse_lock:
            self._sse_clients.append(wfile)

    def remove_sse_client(self, wfile: Any) -> None:
        with self._sse_lock:
            if wfile in self._sse_clients:
                self._sse_clients.remove(wfile)


STATE = ServerState()


def read_last_session() -> str:
    """Прочитать state/last_session.md."""
    path = STATE_DIR / "last_session.md"
    if path.exists():
        return path.read_text(encoding="utf-8")
    return "# Сообщение будущей сессии\n\nФайл пока не создан."


def run_session_transaction() -> tuple[bool, str]:
    """Запустить session_transaction.py и вернуть (success, output)."""
    try:
        result = subprocess.run(
            SESSION_COMMAND,
            cwd=str(PROJECT_ROOT),
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            env=utf8_subprocess_env(),
            timeout=600,  # 10 минут на сессию
        )
        output = result.stdout + result.stderr
        if result.returncode == 0:
            return True, output
        return False, output
    except subprocess.TimeoutExpired:
        return False, "Таймаут: сессия превысила 10 минут"
    except Exception as e:
        return False, str(e)


def auto_session_loop() -> None:
    """Цикл автосессий: запуск → пауза 30с → повтор."""
    while True:
        with STATE.lock:
            if not STATE.auto_enabled:
                STATE.status = "idle"
                break
            STATE.status = "running"

        STATE.broadcast("status", {"status": "running", "message": "Запуск сессии..."})

        success, output = run_session_transaction()

        with STATE.lock:
            STATE.session_count += 1
            STATE.last_run_time = time.strftime("%Y-%m-%d %H:%M:%S")
            if not success:
                STATE.last_error = output[-500:] if len(output) > 500 else output

        # Перечитать last_session.md после сессии
        last_session = read_last_session()

        STATE.broadcast("session_done", {
            "success": success,
            "count": STATE.session_count,
            "time": STATE.last_run_time,
            "last_session": last_session,
            "error": STATE.last_error if not success else "",
        })

        # Пауза 30 секунд перед следующей сессией
        for i in range(30):
            with STATE.lock:
                if not STATE.auto_enabled:
                    break
            time.sleep(1)

        with STATE.lock:
            if not STATE.auto_enabled:
                STATE.status = "idle"
                break
            STATE.status = "auto"


class SessionHandler(BaseHTTPRequestHandler):
    """HTTP-обработчик для управления сессиями."""

    def log_message(self, format: str, *args: Any) -> None:
        # Подавляем стандартные логи для чистоты
        pass

    def _send_json(self, data: dict, status: int = 200) -> None:
        body = json.dumps(data, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()
        self.wfile.write(body)

    def _send_html(self, html: str) -> None:
        body = html.encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "text/html; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _send_sse_headers(self) -> None:
        self.send_response(200)
        self.send_header("Content-Type", "text/event-stream; charset=utf-8")
        self.send_header("Cache-Control", "no-cache")
        self.send_header("Connection", "keep-alive")
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()

    def do_OPTIONS(self) -> None:
        self.send_response(204)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type")
        self.end_headers()

    def do_GET(self) -> None:
        if self.path == "/":
            self._serve_index()
        elif self.path == "/api/last-session":
            self._api_last_session()
        elif self.path == "/api/status":
            self._api_status()
        elif self.path == "/api/events":
            self._api_events()
        else:
            self.send_error(404)

    def do_POST(self) -> None:
        if self.path == "/api/session/start":
            self._api_session_start()
        elif self.path == "/api/auto/toggle":
            self._api_auto_toggle()
        else:
            self.send_error(404)

    def _serve_index(self) -> None:
        index_path = STATIC_DIR / "index.html"
        if index_path.exists():
            html = index_path.read_text(encoding="utf-8")
            self._send_html(html)
        else:
            self._send_html("<h1>index.html не найден</h1>")

    def _api_last_session(self) -> None:
        content = read_last_session()
        self._send_json({"content": content})

    def _api_status(self) -> None:
        with STATE.lock:
            self._send_json({
                "status": STATE.status,
                "auto_enabled": STATE.auto_enabled,
                "session_count": STATE.session_count,
                "last_run_time": STATE.last_run_time,
                "last_error": STATE.last_error,
            })

    def _api_events(self) -> None:
        """SSE-поток для real-time обновлений."""
        self._send_sse_headers()
        STATE.add_sse_client(self.wfile)

        # Отправляем начальное состояние
        with STATE.lock:
            status_data = {
                "status": STATE.status,
                "auto_enabled": STATE.auto_enabled,
                "session_count": STATE.session_count,
                "last_run_time": STATE.last_run_time,
            }
        msg = f"event: status\ndata: {json.dumps(status_data, ensure_ascii=False)}\n\n"
        try:
            self.wfile.write(msg.encode("utf-8"))
            self.wfile.flush()
        except Exception:
            STATE.remove_sse_client(self.wfile)
            return

        # Держим соединение открытым
        try:
            while True:
                time.sleep(1)
                # Проверяем что клиент ещё жив
                try:
                    self.wfile.write(b": keepalive\n\n")
                    self.wfile.flush()
                except Exception:
                    break
        except Exception:
            pass
        finally:
            STATE.remove_sse_client(self.wfile)

    def _api_session_start(self) -> None:
        with STATE.lock:
            if STATE.status == "running":
                self._send_json({"ok": False, "error": "Сессия уже выполняется"}, 409)
                return
            STATE.status = "running"

        STATE.broadcast("status", {"status": "running", "message": "Запуск сессии..."})

        # Запускаем в отдельном потоке, чтобы не блокировать ответ
        def run_and_report():
            success, output = run_session_transaction()
            with STATE.lock:
                STATE.session_count += 1
                STATE.last_run_time = time.strftime("%Y-%m-%d %H:%M:%S")
                if not success:
                    STATE.last_error = output[-500:] if len(output) > 500 else output
                if not STATE.auto_enabled:
                    STATE.status = "idle"
                else:
                    STATE.status = "auto"

            last_session = read_last_session()
            STATE.broadcast("session_done", {
                "success": success,
                "count": STATE.session_count,
                "time": STATE.last_run_time,
                "last_session": last_session,
                "error": STATE.last_error if not success else "",
            })

        t = threading.Thread(target=run_and_report, daemon=True)
        t.start()

        self._send_json({"ok": True, "message": "Сессия запущена"})

    def _api_auto_toggle(self) -> None:
        with STATE.lock:
            STATE.auto_enabled = not STATE.auto_enabled
            enabled = STATE.auto_enabled

            if enabled:
                STATE.status = "auto"
                # Запускаем цикл автосессий
                if STATE.auto_thread is None or not STATE.auto_thread.is_alive():
                    STATE.auto_thread = threading.Thread(target=auto_session_loop, daemon=True)
                    STATE.auto_thread.start()
            else:
                # Цикл остановится сам при следующей проверке
                pass

        STATE.broadcast("status", {
            "status": "auto" if enabled else "idle",
            "auto_enabled": enabled,
            "message": "Автосессия " + ("включена" if enabled else "выключена"),
        })

        self._send_json({"ok": True, "auto_enabled": enabled})

def create_server(host: str, port: int) -> ThreadingHTTPServer:
    """Создать многопоточный HTTP-сервер, чтобы SSE не блокировал API."""
    return ThreadingHTTPServer((host, port), SessionHandler)


def main() -> int:
    parser = argparse.ArgumentParser(description="Сервер управления сессиями ai-lives")
    parser.add_argument("--port", type=int, default=11000, help="Порт (по умолчанию 11000)")
    parser.add_argument("--host", default="127.0.0.1", help="Хост (по умолчанию 127.0.0.1)")
    args = parser.parse_args()

    server = create_server(args.host, args.port)
    print(f"[server] Запущен на http://{args.host}:{args.port}")
    print(f"[server] Корень проекта: {PROJECT_ROOT}")
    print(f"[server] Ctrl+C для остановки")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n[server] Остановка...")
        with STATE.lock:
            STATE.auto_enabled = False
        server.shutdown()

    return 0


if __name__ == "__main__":
    sys.exit(main())
