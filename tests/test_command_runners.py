import io

import pytest

from scripts import command_runners


def test_streaming_runner_streams_output_and_sanitizes_env(monkeypatch, capsys, tmp_path):
    calls = []

    class FakeProcess:
        def __init__(self):
            self.stdout = io.StringIO("first line\nsecond line\n")

        def poll(self):
            return 0

        def wait(self):
            return 0

    def fake_popen(args, cwd, env, text, stdout, stderr, bufsize):
        calls.append((args, cwd, env, text, stdout, stderr, bufsize))
        return FakeProcess()

    monkeypatch.setenv("VIRTUAL_ENV", "C:/old/project/.venv")
    monkeypatch.setattr(command_runners.subprocess, "Popen", fake_popen)

    result = command_runners.streaming_runner(["cmd", "arg"], tmp_path)

    captured = capsys.readouterr()
    assert result == command_runners.CommandResult(0, "first line\nsecond line\n", "")
    assert "[cmd]" in captured.out
    assert "first line" in captured.out
    args, cwd, env, text, stdout, stderr, bufsize = calls[0]
    assert args == ["cmd", "arg"]
    assert cwd == tmp_path
    assert "VIRTUAL_ENV" not in env
    assert text is True
    assert stdout is command_runners.subprocess.PIPE
    assert stderr is command_runners.subprocess.STDOUT
    assert bufsize == 1


def test_streaming_runner_reports_missing_stdout(monkeypatch, tmp_path):
    class FakeProcess:
        stdout = None

    def fake_popen(args, cwd, env, text, stdout, stderr, bufsize):
        return FakeProcess()

    monkeypatch.setattr(command_runners.subprocess, "Popen", fake_popen)

    with pytest.raises(command_runners.CommandExecutionError, match="did not expose stdout"):
        command_runners.run_streaming_command(["cmd"], tmp_path)


def test_streaming_runner_prints_heartbeat_for_silent_process(monkeypatch, capsys, tmp_path):
    class FakeProcess:
        def __init__(self):
            self.stdout = io.StringIO("")
            self.poll_count = 0

        def poll(self):
            self.poll_count += 1
            if self.poll_count == 1:
                return None
            return 0

        def wait(self):
            return 0

    def fake_popen(args, cwd, env, text, stdout, stderr, bufsize):
        return FakeProcess()

    times = iter([0.0, 31.0])
    monkeypatch.setattr(command_runners.subprocess, "Popen", fake_popen)
    monkeypatch.setattr(command_runners.time, "monotonic", lambda: next(times))

    result = command_runners.run_streaming_command(["cmd"], tmp_path, heartbeat_seconds=30)

    captured = capsys.readouterr()
    assert result == command_runners.CommandResult(0, "", "")
    assert "[wait] command is still running" in captured.out