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
    {
        "type": "function",
        "function": {
            "name": "read_lines",
            "description": "Read a 1-based line range from a UTF-8 text file inside the session root.",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string", "description": "Relative path inside the session root."},
                    "start_line": {"type": "integer", "description": "First line to read, 1-based."},
                    "line_count": {"type": "integer", "description": "Number of lines to read."},
                },
                "required": ["path", "start_line", "line_count"],
                "additionalProperties": False,
            },
        },
    },
    {
        "type": "function",
        "function": {
            "name": "replace_text",
            "description": "Replace an exact text fragment in a UTF-8 file without rewriting the whole file.",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string", "description": "Relative path inside the session root."},
                    "old": {"type": "string", "description": "Exact text fragment to replace."},
                    "new": {"type": "string", "description": "Replacement text."},
                    "expected_replacements": {
                        "type": "integer",
                        "description": "Required number of replacements; defaults to 1.",
                    },
                },
                "required": ["path", "old", "new"],
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


def read_lines(root: Path, path: str, start_line: int, line_count: int) -> dict[str, Any]:
    if start_line < 1:
        raise ToolError("start_line must be at least 1")
    if line_count < 1:
        raise ToolError("line_count must be at least 1")
    target = safe_path(root, path)
    if not target.exists():
        raise ToolError(f"file does not exist: {path}")
    if not target.is_file():
        raise ToolError(f"path is not a file: {path}")

    lines = target.read_text(encoding="utf-8").splitlines()
    start_index = start_line - 1
    selected = lines[start_index : start_index + line_count]
    numbered = [f"{start_line + index}: {line}" for index, line in enumerate(selected)]
    return {
        "path": path,
        "start_line": start_line,
        "end_line": start_line + len(selected) - 1 if selected else start_line - 1,
        "total_lines": len(lines),
        "content": "\n".join(numbered),
    }


def replace_text(root: Path, path: str, old: str, new: str, expected_replacements: int = 1) -> dict[str, Any]:
    if not old:
        raise ToolError("old text must be non-empty")
    if expected_replacements < 1:
        raise ToolError("expected_replacements must be at least 1")
    target = safe_path(root, path)
    if not target.exists():
        raise ToolError(f"file does not exist: {path}")
    if not target.is_file():
        raise ToolError(f"path is not a file: {path}")

    content = target.read_text(encoding="utf-8")
    replacements = content.count(old)
    if replacements != expected_replacements:
        raise ToolError(
            f"expected {expected_replacements} replacement(s), found {replacements}: {path}"
        )
    updated = content.replace(old, new, expected_replacements)
    target.write_text(updated, encoding="utf-8")
    return {"path": path, "replacements": replacements, "bytes": len(updated.encode("utf-8"))}


def call_tool(root: Path, name: str, arguments: dict[str, Any]) -> dict[str, Any]:
    if name == "read_file":
        return read_file(root, str(arguments.get("path", "")))
    if name == "write_file":
        return write_file(root, str(arguments.get("path", "")), str(arguments.get("content", "")))
    if name == "read_lines":
        return read_lines(
            root,
            str(arguments.get("path", "")),
            int(arguments.get("start_line", 1)),
            int(arguments.get("line_count", 1)),
        )
    if name == "replace_text":
        return replace_text(
            root,
            str(arguments.get("path", "")),
            str(arguments.get("old", "")),
            str(arguments.get("new", "")),
            int(arguments.get("expected_replacements", 1)),
        )
    raise ToolError(f"unknown tool: {name}")


def tool_result_json(payload: dict[str, Any]) -> str:
    return json.dumps(payload, ensure_ascii=False)