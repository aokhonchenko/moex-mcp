#!/usr/bin/env python3
"""Run the local minimal ai-lives agent."""

from __future__ import annotations

import argparse
import json
import sys
import tomllib
from pathlib import Path
from typing import Any

if __package__ in (None, ""):
    sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from scripts.file_tools import TOOL_SCHEMAS, ToolError, call_tool, tool_result_json
from scripts.llm_client import LlmClientError, OpenAICompatibleClient


DEFAULT_SETTINGS = {
    "step_limit": 30,
    "request_timeout_seconds": 120,
    "temperature": 0.2,
}


class AgentError(RuntimeError):
    """Raised when the local agent cannot complete a session."""


def diagnostic(message: str) -> None:
    print(f"[agent] {message}", flush=True)


def read_settings(path: Path) -> dict[str, Any]:
    if not path.exists():
        return dict(DEFAULT_SETTINGS)
    data = tomllib.loads(path.read_text(encoding="utf-8"))
    settings = dict(DEFAULT_SETTINGS)
    settings.update(data.get("agent", {}))
    return settings


def read_task(prompt_file: Path) -> str:
    if not prompt_file.exists():
        raise AgentError(f"Prompt file does not exist: {prompt_file}")
    return prompt_file.read_text(encoding="utf-8")


def system_message() -> str:
    return """Ты локальный автономный агент проекта ai-lives.

Работай только через доступные инструменты:
- read_file: прочитать UTF-8 файл внутри корня сессии;
- write_file: записать UTF-8 файл внутри корня сессии.

У тебя нет shell-инструмента. Если нужен новый инструмент, создай его как файл проекта через write_file,
но текущую сессию всё равно завершай доступными средствами.

Обязательные результаты каждой успешной сессии:
1. Обновить state/last_session.md.
2. Добавить запись в logs/history.md.
3. Если менялся план, обновить state/current_plan.md.

Все пользовательские артефакты пиши на русском языке. Для завершения ответь обычным финальным
сообщением без вызова инструментов. Если endpoint не поддерживает tool calling, можно вернуть ровно
JSON-объект одного из видов:
{"tool":"read_file","path":"state/last_session.md"}
{"tool":"write_file","path":"state/last_session.md","content":"..."}
{"final":"краткий итог"}
"""


def initial_messages(task: str) -> list[dict[str, Any]]:
    return [
        {"role": "system", "content": system_message()},
        {"role": "user", "content": task},
    ]


def parse_tool_arguments(raw: str) -> dict[str, Any]:
    try:
        parsed = json.loads(raw or "{}")
    except json.JSONDecodeError as exc:
        raise AgentError(f"tool arguments are not valid JSON: {raw}") from exc
    if not isinstance(parsed, dict):
        raise AgentError("tool arguments must be a JSON object")
    return parsed


def assistant_message_for_history(message: dict[str, Any]) -> dict[str, Any]:
    stored = {"role": "assistant", "content": message.get("content") or ""}
    if message.get("tool_calls"):
        stored["tool_calls"] = message["tool_calls"]
    return stored


def execute_native_tool_calls(root: Path, message: dict[str, Any]) -> list[dict[str, Any]]:
    results = []
    for tool_call in message.get("tool_calls") or []:
        function = tool_call.get("function") or {}
        name = function.get("name", "")
        arguments = parse_tool_arguments(function.get("arguments", "{}"))
        diagnostic(f"tool: {name} {arguments}")
        result = call_tool(root, name, arguments)
        results.append(
            {
                "role": "tool",
                "tool_call_id": tool_call.get("id", name),
                "content": tool_result_json({"ok": True, "result": result}),
            }
        )
    return results


def parse_text_protocol(content: str) -> dict[str, Any] | None:
    stripped = content.strip()
    if not stripped.startswith("{") or not stripped.endswith("}"):
        return None
    try:
        parsed = json.loads(stripped)
    except json.JSONDecodeError:
        return None
    return parsed if isinstance(parsed, dict) else None


def execute_text_protocol(root: Path, parsed: dict[str, Any]) -> dict[str, Any] | None:
    if "final" in parsed:
        return None
    name = str(parsed.get("tool", ""))
    if not name:
        return None
    arguments = {key: value for key, value in parsed.items() if key != "tool"}
    diagnostic(f"text tool: {name} {arguments}")
    result = call_tool(root, name, arguments)
    return {
        "role": "user",
        "content": tool_result_json({"tool": name, "ok": True, "result": result}),
    }


def run_agent(root: Path, prompt_file: Path, settings_path: Path) -> int:
    settings = read_settings(settings_path)
    task = read_task(prompt_file)
    client = OpenAICompatibleClient.from_environment(
        timeout=int(settings["request_timeout_seconds"]),
        temperature=float(settings["temperature"]),
    )
    diagnostic(f"root: {root}")
    diagnostic(f"prompt: {prompt_file}")
    diagnostic(f"model: {client.model}")
    diagnostic(f"api_base: {client.base_url}")
    diagnostic(f"step_limit: {settings['step_limit']}")

    messages = initial_messages(task)
    for step in range(1, int(settings["step_limit"]) + 1):
        diagnostic(f"step {step}: requesting model")
        message = client.complete(messages, TOOL_SCHEMAS)
        messages.append(assistant_message_for_history(message))
        tool_calls = message.get("tool_calls") or []
        if tool_calls:
            diagnostic(f"step {step}: model requested {len(tool_calls)} tool call(s)")
            messages.extend(execute_native_tool_calls(root, message))
            continue

        content = message.get("content") or ""
        parsed = parse_text_protocol(content)
        if parsed is not None:
            observation = execute_text_protocol(root, parsed)
            if observation is not None:
                messages.append(observation)
                continue
            final = str(parsed.get("final", "")).strip()
        else:
            final = content.strip()

        diagnostic(f"finished: {final or '<empty final response>'}")
        return 0

    raise AgentError(f"step limit exceeded: {settings['step_limit']}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run the local ai-lives agent.")
    parser.add_argument("--root", required=True, help="Session worktree root.")
    parser.add_argument("--prompt-file", required=True, help="Prompt file generated by run_session.py.")
    parser.add_argument(
        "--settings",
        default="config/project.toml",
        help="Project settings path, relative to --root unless absolute.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    root = Path(args.root).resolve()
    prompt_file = Path(args.prompt_file)
    if not prompt_file.is_absolute():
        prompt_file = root / prompt_file
    settings_path = Path(args.settings)
    if not settings_path.is_absolute():
        settings_path = root / settings_path

    try:
        return run_agent(root, prompt_file, settings_path)
    except (AgentError, ToolError, LlmClientError) as exc:
        print(f"agent session failed: {exc}", file=sys.stderr, flush=True)
        return 1


if __name__ == "__main__":
    sys.exit(main())