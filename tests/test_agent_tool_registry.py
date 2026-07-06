import sys

import pytest

from scripts import file_tools


def tool_names():
    return {schema["function"]["name"] for schema in file_tools.TOOL_SCHEMAS}


def test_registry_discovers_tool_schemas_and_passport():
    names = tool_names()

    assert {
        "read_file",
        "read_lines",
        "replace_text",
        "write_file",
        "run_command",
        "run_pytest",
        "run_python_script",
    } <= names
    assert "`run_command" in file_tools.TOOL_PASSPORT
    assert "`run_pytest" in file_tools.TOOL_PASSPORT


def test_discovered_command_tool_executes_inside_root(tmp_path):
    command = f'"{sys.executable}" -c "print(123)"'

    result = file_tools.call_tool(tmp_path, "run_command", {"command": command, "timeout": 10})

    assert result["returncode"] == 0
    assert result["stdout"].strip() == "123"
    assert result["stderr"] == ""


def test_discovered_command_tool_rejects_cwd_outside_root(tmp_path):
    with pytest.raises(file_tools.ToolError, match="cwd escapes"):
        file_tools.call_tool(
            tmp_path,
            "run_command",
            {"command": "echo nope", "cwd": str(tmp_path.parent), "timeout": 10},
        )


def test_discovered_run_python_script_tool(tmp_path):
    script = tmp_path / "hello.py"
    script.write_text("print('hello from script')\n", encoding="utf-8")

    result = file_tools.call_tool(
        tmp_path,
        "run_python_script",
        {"script_path": "hello.py", "timeout": 10},
    )

    assert result["returncode"] == 0
    assert result["stdout"].strip() == "hello from script"


def test_discovered_run_pytest_tool(tmp_path):
    tests_dir = tmp_path / "tests"
    tests_dir.mkdir()
    (tests_dir / "test_sample.py").write_text("def test_ok():\n    assert True\n", encoding="utf-8")

    result = file_tools.call_tool(
        tmp_path,
        "run_pytest",
        {"test_path": "tests/test_sample.py", "args": ["-q"], "timeout": 30},
    )

    assert result["returncode"] == 0
    assert "1 passed" in result["stdout"]
