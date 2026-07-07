import io
import json
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


def test_session_transaction_command_uses_uv_and_streams_output(monkeypatch):
    calls = []

    class FakeProcess:
        def __init__(self):
            self.stdout = io.StringIO("first line\nsecond line\n")

        def poll(self):
            if self.stdout.tell() >= len(self.stdout.getvalue()):
                return 0
            return None

        def wait(self):
            return 0

    def fake_popen(cmd, cwd, stdout, stderr, text, encoding, errors, env):
        calls.append(
            {
                "cmd": cmd,
                "cwd": cwd,
                "stdout": stdout,
                "stderr": stderr,
                "text": text,
                "encoding": encoding,
                "errors": errors,
                "env": env,
            }
        )
        return FakeProcess()

    monkeypatch.setattr(session_server.subprocess, "Popen", fake_popen)
    streamed = []

    success, output = session_server.run_session_transaction(streamed.append)

    assert success is True
    assert output == "first line\nsecond line\n"
    assert streamed == ["first line", "second line"]
    assert len(calls) == 1
    call = calls[0]
    assert call["cmd"] == ["uv", "run", "python", "scripts/session_transaction.py"]
    assert call["cwd"] == str(session_server.PROJECT_ROOT)
    assert call["stdout"] is session_server.subprocess.PIPE
    assert call["stderr"] is session_server.subprocess.STDOUT
    assert call["text"] is True
    assert call["encoding"] == "utf-8"
    assert call["errors"] == "replace"
    assert call["env"]["PYTHONUTF8"] == "1"
    assert call["env"]["PYTHONIOENCODING"] == "utf-8"


def test_create_server_uses_threading_http_server():
    httpd = session_server.create_server("127.0.0.1", 0)
    try:
        assert isinstance(httpd, ThreadingHTTPServer)
    finally:
        httpd.server_close()


def test_session_start_is_not_blocked_by_open_sse(monkeypatch):
    reset_server_state()
    called = threading.Event()

    def fake_run_session_transaction(on_output=None):
        if on_output is not None:
            on_output("streamed output")
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

def test_frontend_streams_session_log_and_renders_markdown_without_last_session_scroll():
    html = (session_server.STATIC_DIR / "index.html").read_text(encoding="utf-8")

    assert "addEventListener('session_log'" in html
    assert "function renderMarkdown" in html
    assert "setLastSessionMarkdown" in html
    assert ".last-session::-webkit-scrollbar" not in html
    last_session_css = html.split(".last-session {", 1)[1].split("}", 1)[0]
    assert "overflow-y: auto" not in last_session_css
    assert "overflow: visible" in last_session_css
