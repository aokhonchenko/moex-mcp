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
    assert file_tools.call_tool(
        tmp_path, "read_lines", {"path": "note.md", "start_line": 1, "line_count": 1}
    )["content"] == "1: ok"
    assert file_tools.call_tool(
        tmp_path, "replace_text", {"path": "note.md", "old": "ok", "new": "done"}
    )["replacements"] == 1
    with pytest.raises(file_tools.ToolError, match="unknown tool"):
        file_tools.call_tool(tmp_path, "shell", {})


def test_read_lines_returns_numbered_range(tmp_path):
    path = tmp_path / "notes.md"
    path.write_text("one\ntwo\nthree\nfour\n", encoding="utf-8")

    result = file_tools.read_lines(tmp_path, "notes.md", 2, 2)

    assert result == {
        "path": "notes.md",
        "start_line": 2,
        "end_line": 3,
        "total_lines": 4,
        "content": "2: two\n3: three",
    }


def test_read_lines_validates_range(tmp_path):
    (tmp_path / "notes.md").write_text("one\n", encoding="utf-8")

    with pytest.raises(file_tools.ToolError, match="start_line"):
        file_tools.read_lines(tmp_path, "notes.md", 0, 1)
    with pytest.raises(file_tools.ToolError, match="line_count"):
        file_tools.read_lines(tmp_path, "notes.md", 1, 0)


def test_replace_text_replaces_exact_fragment(tmp_path):
    path = tmp_path / "notes.md"
    path.write_text("alpha\nbeta\ngamma\n", encoding="utf-8")

    result = file_tools.replace_text(tmp_path, "notes.md", "beta", "delta")

    assert result["path"] == "notes.md"
    assert result["replacements"] == 1
    assert path.read_text(encoding="utf-8") == "alpha\ndelta\ngamma\n"


def test_replace_text_requires_expected_replacements(tmp_path):
    path = tmp_path / "notes.md"
    path.write_text("same same", encoding="utf-8")

    with pytest.raises(file_tools.ToolError, match="expected 1 replacement"):
        file_tools.replace_text(tmp_path, "notes.md", "same", "other")
    with pytest.raises(file_tools.ToolError, match="old text"):
        file_tools.replace_text(tmp_path, "notes.md", "", "other")


def test_read_lines_rejects_missing_and_directory(tmp_path):
    with pytest.raises(file_tools.ToolError, match="does not exist"):
        file_tools.read_lines(tmp_path, "missing.md", 1, 1)

    directory = tmp_path / "folder"
    directory.mkdir()
    with pytest.raises(file_tools.ToolError, match="not a file"):
        file_tools.read_lines(tmp_path, "folder", 1, 1)


def test_replace_text_rejects_invalid_count_missing_and_directory(tmp_path):
    with pytest.raises(file_tools.ToolError, match="expected_replacements"):
        file_tools.replace_text(tmp_path, "missing.md", "old", "new", 0)
    with pytest.raises(file_tools.ToolError, match="does not exist"):
        file_tools.replace_text(tmp_path, "missing.md", "old", "new")

    directory = tmp_path / "folder"
    directory.mkdir()
    with pytest.raises(file_tools.ToolError, match="not a file"):
        file_tools.replace_text(tmp_path, "folder", "old", "new")


def test_tool_result_json_preserves_russian_text():
    payload = file_tools.tool_result_json({"content": "привет"})

    assert json.loads(payload)["content"] == "привет"