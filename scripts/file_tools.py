"""Registry for directory-based tools available to the local autonomous agent."""

from __future__ import annotations

import importlib
import json
from pathlib import Path
from types import ModuleType
from typing import Any

from scripts.agent_tools._shared import ToolError, safe_path
from scripts.agent_tools.read_file.tool import read_file
from scripts.agent_tools.read_lines.tool import read_lines
from scripts.agent_tools.replace_text.tool import replace_text
from scripts.agent_tools.write_file.tool import write_file

TOOLS_PACKAGE = "scripts.agent_tools"
TOOLS_ROOT = Path(__file__).resolve().parent / "agent_tools"
REQUIRED_TOOL_FUNCTIONS = ("schema", "passport", "handle")


def discover_tool_modules() -> list[ModuleType]:
    modules: list[ModuleType] = []
    for child in sorted(TOOLS_ROOT.iterdir(), key=lambda path: path.name):
        if not child.is_dir() or child.name.startswith("_"):
            continue
        if not (child / "tool.py").is_file():
            continue
        module = importlib.import_module(f"{TOOLS_PACKAGE}.{child.name}.tool")
        for function_name in REQUIRED_TOOL_FUNCTIONS:
            if not callable(getattr(module, function_name, None)):
                raise ToolError(f"tool {child.name} does not define callable {function_name}()")
        modules.append(module)
    return modules


def schema_tool_name(schema: dict[str, Any]) -> str:
    try:
        name = schema["function"]["name"]
    except KeyError as exc:
        raise ToolError(f"tool schema is missing function.name: {schema}") from exc
    if not isinstance(name, str) or not name:
        raise ToolError(f"tool schema has invalid function.name: {schema}")
    return name


def load_tools() -> dict[str, ModuleType]:
    tools: dict[str, ModuleType] = {}
    for module in discover_tool_modules():
        schema = module.schema()
        name = schema_tool_name(schema)
        if name in tools:
            raise ToolError(f"duplicate tool name: {name}")
        tools[name] = module
    return tools


def tool_schemas() -> list[dict[str, Any]]:
    return [module.schema() for module in load_tools().values()]


def tool_passport() -> str:
    lines = [module.passport().strip() for module in load_tools().values()]
    return "\n".join(line for line in lines if line)


def call_tool(root: Path, name: str, arguments: dict[str, Any]) -> dict[str, Any]:
    module = load_tools().get(name)
    if module is None:
        raise ToolError(f"unknown tool: {name}")
    return module.handle(root, arguments)


def tool_result_json(payload: dict[str, Any]) -> str:
    return json.dumps(payload, ensure_ascii=False)


TOOL_SCHEMAS: list[dict[str, Any]] = tool_schemas()
TOOL_PASSPORT = tool_passport()
