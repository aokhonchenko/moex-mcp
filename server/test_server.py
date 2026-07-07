"""Smoke-тесты для server/server.py."""
import json
import threading
import time
import urllib.request
from http.server import HTTPServer

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from server.server import SessionHandler


def test_server():
    server = HTTPServer(('127.0.0.1', 11002), SessionHandler)
    t = threading.Thread(target=server.serve_forever, daemon=True)
    t.start()
    time.sleep(0.3)

    # GET /
    resp = urllib.request.urlopen('http://127.0.0.1:11002/')
    html = resp.read().decode()
    assert 'ai-lives' in html
    print('PASS: GET /')

    # GET /api/status
    resp = urllib.request.urlopen('http://127.0.0.1:11002/api/status')
    data = json.loads(resp.read())
    assert data['status'] == 'idle'
    print(f'PASS: GET /api/status -> {data}')

    # GET /api/last-session
    resp = urllib.request.urlopen('http://127.0.0.1:11002/api/last-session')
    data = json.loads(resp.read())
    assert 'content' in data
    print(f'PASS: GET /api/last-session ({len(data["content"])} chars)')

    # POST /api/auto/toggle (on)
    req = urllib.request.Request('http://127.0.0.1:11002/api/auto/toggle', method='POST', data=b'')
    resp = urllib.request.urlopen(req)
    data = json.loads(resp.read())
    assert data['auto_enabled'] is True
    print(f'PASS: POST /api/auto/toggle ON')

    # POST /api/auto/toggle (off)
    req = urllib.request.Request('http://127.0.0.1:11002/api/auto/toggle', method='POST', data=b'')
    resp = urllib.request.urlopen(req)
    data = json.loads(resp.read())
    assert data['auto_enabled'] is False
    print(f'PASS: POST /api/auto/toggle OFF')

    server.shutdown()
    print('\nAll 5 tests passed!')


if __name__ == '__main__':
    test_server()
