import json

import pytest

from scripts import file_tools


def test_read_and_write_file_inside_root(tmp_path):
    result = file_tools.write_file(tmp_path, "state/last_session.md", "текст")

    assert result["path"] == "state/last_session.md"
    assert result["bytes"] == len("текст".encode("utf-8"))
    assert file_tools.read_file(tmp_path, "state/last_session.md") == {
        "path": "state/last_session.md",
        "content": "текст",
    }


def test_safe_path_rejects_unsafe_paths(tmp_path):
    with pytest.raises(file_tools.ToolError, match="absolute"):
        file_tools.safe_path(tmp_path, str(tmp_path / "file.txt"))
    with pytest.raises(file_tools.ToolError, match="escapes"):
        file_tools.safe_path(tmp_path, "../outside.txt")
    with pytest.raises(file_tools.ToolError, match="forbidden"):
        file_tools.safe_path(tmp_path, ".git/config")


def test_read_file_reports_missing_file(tmp_path):
    with pytest.raises(file_tools.ToolError, match="does not exist"):
        file_tools.read_file(tmp_path, "missing.md")


def test_call_tool_dispatches_and_rejects_unknown_tool(tmp_path):
    assert file_tools.call_tool(tmp_path, "write_file", {"path": "note.md", "content": "ok"})["path"] == "note.md"
    assert file_tools.call_tool(tmp_path, "read_file", {"path": "note.md"})["content"] == "ok"
    with pytest.raises(file_tools.ToolError, match="unknown tool"):
        file_tools.call_tool(tmp_path, "shell", {})


def test_tool_result_json_preserves_russian_text():
    payload = file_tools.tool_result_json({"content": "привет"})

    assert json.loads(payload)["content"] == "привет"