"""File tools available to the local autonomous agent."""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any


class ToolError(RuntimeError):
    """Raised when a local tool call is invalid or unsafe."""


FORBIDDEN_PARTS = {".git", ".venv", "__pycache__", ".pytest_cache", "runs"}


TOOL_SCHEMAS: list[dict[str, Any]] = [
    {
        "type": "function",
        "function": {
            "name": "read_file",
            "description": "Read a UTF-8 text file inside the session root.",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string", "description": "Relative path inside the session root."},
                },
                "required": ["path"],
                "additionalProperties": False,
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "write_file",
            "description": "Write a UTF-8 text file inside the session root, creating parent directories.",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string", "description": "Relative path inside the session root."},
                    "content": {"type": "string", "description": "Full file content to write."},
                },
                "required": ["path", "content"],
                "additionalProperties": False,
            },
        },
    },
]


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


def read_file(root: Path, path: str) -> dict[str, Any]:
    target = safe_path(root, path)
    if not target.exists():
        raise ToolError(f"file does not exist: {path}")
    if not target.is_file():
        raise ToolError(f"path is not a file: {path}")
    return {"path": path, "content": target.read_text(encoding="utf-8")}


def write_file(root: Path, path: str, content: str) -> dict[str, Any]:
    target = safe_path(root, path)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content, encoding="utf-8")
    return {"path": path, "bytes": len(content.encode("utf-8"))}


def call_tool(root: Path, name: str, arguments: dict[str, Any]) -> dict[str, Any]:
    if name == "read_file":
        return read_file(root, str(arguments.get("path", "")))
    if name == "write_file":
        return write_file(root, str(arguments.get("path", "")), str(arguments.get("content", "")))
    raise ToolError(f"unknown tool: {name}")


def tool_result_json(payload: dict[str, Any]) -> str:
    return json.dumps(payload, ensure_ascii=False)