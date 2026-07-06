import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from tools.code_analyzer import analyze_file


def test_analyze_file_counts_trailing_newline_as_existing_line(tmp_path):
    path = tmp_path / "sample.py"
    path.write_text("line1\nline2\nline3\n", encoding="utf-8")

    result = analyze_file(str(path))

    assert result.lines == 3


def test_analyze_file_counts_empty_file_as_zero_lines(tmp_path):
    path = tmp_path / "empty.py"
    path.write_text("", encoding="utf-8")

    result = analyze_file(str(path))

    assert result.lines == 0