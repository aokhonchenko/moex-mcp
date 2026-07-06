"""Shared helpers for directory-based agent tools."""

from __future__ import annotations

import os
import subprocess
import sys
import time
from pathlib import Path
from typing import Any, Sequence


class ToolError(RuntimeError):
    """Raised when a local tool call is invalid or unsafe."""


FORBIDDEN_PARTS = {".git", ".venv", "__pycache__", ".pytest_cache", "runs"}
MAX_TIMEOUT_SECONDS = 300.0
MAX_OUTPUT_CHARS = 20000


def safe_path(root: Path, path: str) -> Path:
    if not path or path.strip() != path:
        raise ToolError("path must be a non-empty relative path without surrounding whitespace")
    requested = Path(path)
    if requested.is_absolute():
        raise ToolError("absolute paths are not allowed")
    if any(part in FORBIDDEN_PARTS for part in requested.parts):
        raise ToolError(f"path contains a forbidden part: {path}")
    resolved_root = root.resolve()
    resolved = (resolved_root / requested).resolve()
    try:
        resolved.relative_to(resolved_root)
    except ValueError as exc:
        raise ToolError(f"path escapes the session root: {path}") from exc
    return resolved


def safe_cwd(root: Path, cwd: str | None) -> Path:
    resolved_root = root.resolve()
    if cwd is None or not str(cwd).strip():
        return resolved_root
    requested = Path(str(cwd))
    resolved = requested.resolve() if requested.is_absolute() else (resolved_root / requested).resolve()
    try:
        resolved.relative_to(resolved_root)
    except ValueError as exc:
        raise ToolError(f"cwd escapes the session root: {cwd}") from exc
    if not resolved.exists():
        raise ToolError(f"cwd does not exist: {cwd}")
    if not resolved.is_dir():
        raise ToolError(f"cwd is not a directory: {cwd}")
    return resolved


def timeout_seconds(value: Any, default: float = 120.0) -> float:
    if value in (None, ""):
        return default
    try:
        timeout = float(value)
    except (TypeError, ValueError) as exc:
        raise ToolError(f"timeout must be a number: {value}") from exc
    if timeout <= 0:
        raise ToolError("timeout must be positive")
    if timeout > MAX_TIMEOUT_SECONDS:
        raise ToolError(f"timeout must be <= {int(MAX_TIMEOUT_SECONDS)} seconds")
    return timeout


def subprocess_env() -> dict[str, str]:
    env = os.environ.copy()
    env.pop("VIRTUAL_ENV", None)
    return env


def normalize_args(value: Any, field: str) -> list[str]:
    if value in (None, ""):
        return []
    if not isinstance(value, list):
        raise ToolError(f"{field} must be an array of strings")
    return [str(item) for item in value]


def trim_output(value: str) -> tuple[str, bool]:
    if len(value) <= MAX_OUTPUT_CHARS:
        return value, False
    return value[:MAX_OUTPUT_CHARS], True


def command_result(
    command: str | Sequence[str],
    cwd: Path,
    timeout: float,
    *,
    shell: bool = False,
) -> dict[str, Any]:
    start = time.monotonic()
    try:
        completed = subprocess.run(
            command,
            cwd=cwd,
            env=subprocess_env(),
            text=True,
            encoding="utf-8",
            errors="replace",
            capture_output=True,
            timeout=timeout,
            shell=shell,
        )
        stdout, stdout_truncated = trim_output(completed.stdout or "")
        stderr, stderr_truncated = trim_output(completed.stderr or "")
        return {
            "command": command if isinstance(command, str) else list(command),
            "cwd": str(cwd),
            "returncode": completed.returncode,
            "stdout": stdout,
            "stderr": stderr,
            "stdout_truncated": stdout_truncated,
            "stderr_truncated": stderr_truncated,
            "duration_seconds": round(time.monotonic() - start, 3),
        }
    except subprocess.TimeoutExpired as exc:
        stdout, stdout_truncated = trim_output(exc.stdout or "")
        stderr, stderr_truncated = trim_output(exc.stderr or "")
        return {
            "command": command if isinstance(command, str) else list(command),
            "cwd": str(cwd),
            "returncode": -1,
            "stdout": stdout,
            "stderr": stderr or f"command timed out after {timeout} seconds",
            "stdout_truncated": stdout_truncated,
            "stderr_truncated": stderr_truncated,
            "duration_seconds": round(time.monotonic() - start, 3),
            "timed_out": True,
        }
    except OSError as exc:
        raise ToolError(f"failed to start command: {exc}") from exc


def python_executable() -> str:
    return sys.executable
