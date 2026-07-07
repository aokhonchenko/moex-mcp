import json
import subprocess
import threading
import urllib.request
from http.server import ThreadingHTTPServer

from server import server as session_server


def reset_server_state():
    with session_server.STATE.lock:
        session_server.STATE.status = "idle"
        session_server.STATE.auto_enabled = False
        session_server.STATE.last_error = ""
        session_server.STATE.last_run_time = ""
        session_server.STATE.session_count = 0
    with session_server.STATE._sse_lock:
        session_server.STATE._sse_clients.clear()


def test_session_transaction_command_uses_uv(monkeypatch):
    calls = []

    def fake_run(cmd, cwd, capture_output, text, timeout):
        calls.append(
            {
                "cmd": cmd,
                "cwd": cwd,
                "capture_output": capture_output,
                "text": text,
                "timeout": timeout,
            }
        )
        return subprocess.CompletedProcess(cmd, 0, "ok", "")

    monkeypatch.setattr(session_server.subprocess, "run", fake_run)

    success, output = session_server.run_session_transaction()

    assert success is True
    assert output == "ok"
    assert calls == [
        {
            "cmd": ["uv", "run", "python", "scripts/session_transaction.py"],
            "cwd": str(session_server.PROJECT_ROOT),
            "capture_output": True,
            "text": True,
            "timeout": 600,
        }
    ]


def test_create_server_uses_threading_http_server():
    httpd = session_server.create_server("127.0.0.1", 0)
    try:
        assert isinstance(httpd, ThreadingHTTPServer)
    finally:
        httpd.server_close()


def test_session_start_is_not_blocked_by_open_sse(monkeypatch):
    reset_server_state()
    called = threading.Event()

    def fake_run_session_transaction():
        called.set()
        return True, "ok"

    monkeypatch.setattr(session_server, "run_session_transaction", fake_run_session_transaction)

    httpd = session_server.create_server("127.0.0.1", 0)
    port = httpd.server_address[1]
    thread = threading.Thread(target=httpd.serve_forever, daemon=True)
    thread.start()
    sse_response = None

    try:
        sse_response = urllib.request.urlopen(f"http://127.0.0.1:{port}/api/events", timeout=2)
        request = urllib.request.Request(
            f"http://127.0.0.1:{port}/api/session/start",
            data=b"",
            method="POST",
        )
        response = urllib.request.urlopen(request, timeout=2)
        payload = json.loads(response.read().decode("utf-8"))

        assert payload["ok"] is True
        assert called.wait(2)
    finally:
        if sse_response is not None:
            sse_response.close()
        httpd.shutdown()
        httpd.server_close()
        thread.join(timeout=2)
        reset_server_state()
