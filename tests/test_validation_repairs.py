from pathlib import Path

import pytest

from scripts import validation_repairs
from scripts.command_runners import CommandResult


class RepairRunner:
    def __init__(self, failures_before_success: int = 0, fail_forever: bool = False):
        self.failures_before_success = failures_before_success
        self.fail_forever = fail_forever
        self.commands = []
        self.repairs = 0

    def __call__(self, args, cwd):
        self.commands.append((list(args), Path(cwd)))
        if self.fail_forever or self.failures_before_success > 0:
            if self.failures_before_success > 0:
                self.failures_before_success -= 1
            return CommandResult(1, "", "tests failed")
        return CommandResult(0, "tests ok", "")


def failure_details(result: CommandResult) -> str:
    return result.stderr or result.stdout


def repair_session(worktree: Path, agent_command: str, runner):
    runner.repairs += 1
    (worktree / "state").mkdir(exist_ok=True)
    (worktree / "logs").mkdir(exist_ok=True)
    (worktree / "state" / "last_session.md").write_text("repaired\n", encoding="utf-8")
    (worktree / "logs" / "history.md").write_text("repaired\n", encoding="utf-8")


def assert_required_files(worktree: Path) -> None:
    assert (worktree / "state" / "last_session.md").exists()
    assert (worktree / "logs" / "history.md").exists()


def no_file_size_errors(worktree: Path) -> None:
    return None


def test_run_validation_with_repairs_reruns_agent_after_failed_checks(tmp_path):
    runner = RepairRunner(failures_before_success=1)
    messages = []

    validation_repairs.run_validation_with_repairs(
        tmp_path,
        "agent --ok",
        ["check"],
        repair_attempts=2,
        runner=runner,
        failure_details=failure_details,
        run_repair_session=repair_session,
        ensure_required_session_files=assert_required_files,
        ensure_file_size_policy=no_file_size_errors,
        diagnostic=messages.append,
    )

    assert runner.repairs == 1
    assert runner.commands == [(["check"], tmp_path), (["check"], tmp_path)]
    context = (tmp_path / "state" / "check_failure.md").read_text(encoding="utf-8")
    assert "Диагностика упавших проверок" in context
    assert "tests failed" in context
    external = (tmp_path / "state" / "external_messages.md").read_text(encoding="utf-8")
    assert context in external
    assert any("repair attempt 1/2" in message for message in messages)


def test_run_validation_with_repairs_raises_after_exhausted_attempts(tmp_path):
    runner = RepairRunner(fail_forever=True)

    with pytest.raises(RuntimeError, match="after repair attempts"):
        validation_repairs.run_validation_with_repairs(
            tmp_path,
            "agent --ok",
            ["check"],
            repair_attempts=1,
            runner=runner,
            failure_details=failure_details,
            run_repair_session=repair_session,
            ensure_required_session_files=assert_required_files,
            ensure_file_size_policy=no_file_size_errors,
            diagnostic=lambda message: None,
        )

    assert runner.repairs == 1
    assert (tmp_path / "state" / "check_failure.md").exists()