"""Repair loop for failed validation checks inside a transactional session."""

from __future__ import annotations

from pathlib import Path
from typing import Callable, Sequence

from scripts.command_runners import CommandResult


DEFAULT_REPAIR_ATTEMPTS = 2
CommandRunner = Callable[[Sequence[str], Path], CommandResult]
FailureDetails = Callable[[CommandResult], str]
RepairSession = Callable[[Path, str, CommandRunner], None]
WorktreeCheck = Callable[[Path], None]
Diagnostic = Callable[[str], None]


def check_failure_context(
    check_command: Sequence[str],
    details: str,
    attempt: int,
    max_attempts: int,
) -> str:
    command = " ".join(check_command)
    return (
        "# Диагностика упавших проверок\n\n"
        f"- Попытка исправления: {attempt}/{max_attempts}\n"
        f"- Команда проверки: `{command}`\n\n"
        "Проверки упали. Это не финальный результат сессии: исправь ошибки "
        "в текущем временном worktree и снова обнови обязательные файлы сессии.\n\n"
        "## Вывод проверок\n\n"
        "```text\n"
        f"{details}\n"
        "```\n"
    )


def write_check_failure_context(
    worktree: Path,
    check_command: Sequence[str],
    details: str,
    attempt: int,
    max_attempts: int,
) -> Path:
    state_dir = worktree / "state"
    state_dir.mkdir(parents=True, exist_ok=True)
    content = check_failure_context(check_command, details, attempt, max_attempts)

    failure_path = state_dir / "check_failure.md"
    failure_path.write_text(content, encoding="utf-8")

    external_path = state_dir / "external_messages.md"
    previous = external_path.read_text(encoding="utf-8") if external_path.exists() else ""
    separator = "\n\n" if previous.strip() else ""
    external_path.write_text(f"{previous}{separator}{content}", encoding="utf-8")
    return failure_path


def run_validation_with_repairs(
    worktree: Path,
    agent_command: str,
    check_command: Sequence[str],
    *,
    repair_attempts: int,
    runner: CommandRunner,
    failure_details: FailureDetails,
    run_repair_session: RepairSession,
    ensure_required_session_files: WorktreeCheck,
    ensure_file_size_policy: WorktreeCheck,
    diagnostic: Diagnostic,
) -> None:
    attempts_used = 0
    while True:
        result = runner(check_command, worktree)
        if result.returncode == 0:
            return

        details = failure_details(result) or "see streamed output above"
        if attempts_used >= repair_attempts:
            raise RuntimeError(f"checks failed after repair attempts: {details}")

        attempts_used += 1
        diagnostic(
            f"validation checks failed; running repair attempt "
            f"{attempts_used}/{repair_attempts}"
        )
        context_path = write_check_failure_context(
            worktree,
            check_command,
            details,
            attempts_used,
            repair_attempts,
        )
        diagnostic(f"wrote check failure context: {context_path.relative_to(worktree)}")

        run_repair_session(worktree, agent_command, runner)
        ensure_required_session_files(worktree)
        ensure_file_size_policy(worktree)