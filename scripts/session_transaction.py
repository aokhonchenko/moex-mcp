#!/usr/bin/env python3
"""Run one autonomous session as an all-or-nothing Git transaction."""

from __future__ import annotations

import argparse
import os
import sys
from contextlib import contextmanager
from pathlib import Path
from typing import Callable, Iterator, Sequence

if __package__ in (None, ""):
    sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from scripts import run_session
from scripts.command_runners import (
    CommandExecutionError,
    CommandResult,
    default_runner,
    streaming_runner,
)


class TransactionError(RuntimeError):
    """Raised when a session cannot be applied atomically."""


CommandRunner = Callable[[Sequence[str], Path], CommandResult]


SENSITIVE_DIRS = {".git", ".pytest_cache", "__pycache__", ".venv", "runs"}
GENERATED_LONG_FILES = {"uv.lock", "poetry.lock", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"}
HUMAN_INPUT_FILES = {"GLOBAL_TARGET.md", "state/external_messages.md"}


def diagnostic(message: str) -> None:
    print(f"[session] {message}", flush=True)


def run_checked(
    runner: CommandRunner,
    args: Sequence[str],
    cwd: Path,
    action: str,
) -> CommandResult:
    result = runner(args, cwd)
    if result.returncode != 0:
        details = (result.stderr or result.stdout).strip()
        if details:
            raise TransactionError(f"{action}: {details}")
        raise TransactionError(f"{action}: command failed with exit code {result.returncode}")
    return result


def git(runner: CommandRunner, root: Path, *args: str) -> CommandResult:
    return run_checked(runner, ["git", *args], root, f"git {' '.join(args)}")


def ensure_git_repo(root: Path, runner: CommandRunner = default_runner) -> None:
    result = git(runner, root, "rev-parse", "--show-toplevel")
    repo_root = Path(result.stdout.strip()).resolve()
    if repo_root != root.resolve():
        raise TransactionError(f"expected git root {root}, got {repo_root}")


def ensure_clean_worktree(root: Path, runner: CommandRunner = default_runner) -> None:
    result = git(runner, root, "status", "--porcelain")
    if result.stdout.strip():
        raise TransactionError("main worktree must be clean before transactional session")






def parse_env_line(line: str) -> tuple[str, str] | None:
    stripped = line.strip()
    if not stripped or stripped.startswith("#"):
        return None
    if stripped.startswith("export "):
        stripped = stripped.removeprefix("export ").strip()
    if "=" not in stripped:
        return None
    key, value = stripped.split("=", 1)
    key = key.strip()
    value = value.strip().strip('"').strip("'")
    if not key:
        return None
    return key, value


def load_dotenv(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}

    loaded: dict[str, str] = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        parsed = parse_env_line(line)
        if parsed is None:
            continue
        key, value = parsed
        os.environ.setdefault(key, value)
        loaded[key] = os.environ[key]
    return loaded

def parse_porcelain_paths(status: str) -> list[tuple[str, str]]:
    changes: list[tuple[str, str]] = []
    for line in status.splitlines():
        if not line:
            continue
        code = line[:2]
        path = line[3:].strip().strip('"')
        if " -> " in path:
            raise TransactionError("renamed files are not allowed as human input")
        changes.append((code, path.replace("\\", "/")))
    return changes


def is_human_input_change(code: str, path: str) -> bool:
    if "D" in code:
        return False
    if path in HUMAN_INPUT_FILES:
        return True
    return path.startswith("state/questions/") and path.endswith(".md")


def checkpoint_human_input(root: Path, runner: CommandRunner = default_runner) -> bool:
    result = git(runner, root, "status", "--porcelain")
    changes = parse_porcelain_paths(result.stdout)
    if not changes:
        return False

    rejected = [path for code, path in changes if not is_human_input_change(code, path)]
    if rejected:
        joined = ", ".join(rejected)
        raise TransactionError(f"main worktree has non-human changes: {joined}")

    paths = [path for _, path in changes]
    git(runner, root, "add", "--", *paths)
    git(runner, root, "commit", "-m", "record human input before session")
    return True


def current_branch(root: Path, runner: CommandRunner = default_runner) -> str:
    result = git(runner, root, "branch", "--show-current")
    branch = result.stdout.strip()
    if not branch:
        raise TransactionError("main worktree must be on a branch, not detached HEAD")
    return branch


def session_id(session: int) -> str:
    return f"{session:04d}"


@contextmanager
def lock_file(root: Path) -> Iterator[Path]:
    lock_path = root / ".session.lock"
    flags = os.O_CREAT | os.O_EXCL | os.O_WRONLY
    try:
        fd = os.open(lock_path, flags)
    except FileExistsError as exc:
        raise TransactionError(f"session lock already exists: {lock_path}") from exc

    try:
        with os.fdopen(fd, "w", encoding="utf-8") as handle:
            handle.write(f"pid={os.getpid()}\n")
        yield lock_path
    finally:
        try:
            lock_path.unlink()
        except FileNotFoundError:
            pass


def default_runs_dir(root: Path) -> Path:
    return root / "runs"


def create_worktree(
    root: Path,
    runs_dir: Path,
    session: int,
    runner: CommandRunner = default_runner,
) -> tuple[Path, str]:
    sid = session_id(session)
    branch = f"session/{sid}"
    worktree = runs_dir / f"session-{sid}"
    if worktree.exists():
        raise TransactionError(f"worktree already exists: {worktree}")

    runs_dir.mkdir(parents=True, exist_ok=True)
    git(runner, root, "worktree", "add", "-b", branch, str(worktree), "HEAD")
    return worktree, branch


def remove_worktree_and_branch(
    root: Path,
    worktree: Path,
    branch: str,
    runner: CommandRunner = default_runner,
    force_branch: bool = True,
) -> None:
    runner(["git", "worktree", "remove", "--force", str(worktree)], root)
    delete_flag = "-D" if force_branch else "-d"
    runner(["git", "branch", delete_flag, branch], root)
    try:
        worktree.parent.rmdir()
    except OSError:
        pass


def run_inner_session(
    worktree: Path,
    agent_command: str,
    runner: CommandRunner = default_runner,
) -> None:
    script = worktree / "scripts" / "run_session.py"
    args = [sys.executable, str(script), "--agent-command", agent_command]
    result = runner(args, worktree)
    if result.returncode != 0:
        details = (result.stderr or result.stdout).strip()
        raise TransactionError(f"agent session failed: {details}")


def run_checks(
    worktree: Path,
    check_command: Sequence[str],
    runner: CommandRunner = default_runner,
) -> None:
    result = runner(check_command, worktree)
    if result.returncode != 0:
        details = (result.stderr or result.stdout).strip()
        raise TransactionError(f"checks failed: {details}")


def ensure_required_session_files(worktree: Path) -> None:
    required = [
        worktree / "state" / "last_session.md",
        worktree / "logs" / "history.md",
    ]
    missing = [str(path.relative_to(worktree)) for path in required if not path.exists()]
    if missing:
        raise TransactionError(f"required session files are missing: {', '.join(missing)}")


def find_oversized_files(root: Path, max_lines: int = 500) -> list[Path]:
    oversized: list[Path] = []
    for path in root.rglob("*"):
        if not path.is_file():
            continue
        if any(part in SENSITIVE_DIRS for part in path.parts):
            continue
        if path.name in GENERATED_LONG_FILES:
            continue
        try:
            line_count = sum(1 for _ in path.open("r", encoding="utf-8"))
        except UnicodeDecodeError:
            continue
        if line_count > max_lines:
            oversized.append(path)
    return oversized


def ensure_file_size_policy(worktree: Path, max_lines: int = 500) -> None:
    oversized = find_oversized_files(worktree, max_lines=max_lines)
    if oversized:
        relative = ", ".join(str(path.relative_to(worktree)) for path in oversized)
        raise TransactionError(f"files exceed {max_lines} lines and must be decomposed: {relative}")


def ensure_session_changed_worktree(
    worktree: Path,
    runner: CommandRunner = default_runner,
) -> None:
    result = git(runner, worktree, "status", "--porcelain")
    if not result.stdout.strip():
        raise TransactionError("session produced no tracked or untracked changes")


def commit_session(
    worktree: Path,
    session: int,
    runner: CommandRunner = default_runner,
) -> str:
    git(runner, worktree, "add", "-A")
    git(runner, worktree, "commit", "-m", f"session {session_id(session)}")
    result = git(runner, worktree, "rev-parse", "HEAD")
    return result.stdout.strip()


def apply_session_commit(
    root: Path,
    branch: str,
    runner: CommandRunner = default_runner,
) -> None:
    git(runner, root, "merge", "--ff-only", branch)




def default_agent_command() -> str:
    return 'uv run python scripts/run_mini_agent.py --root "{ROOT}" --prompt-file "{PROMPT_FILE}"'

def run_transaction(
    root: Path,
    agent_command: str,
    runs_dir: Path | None = None,
    check_command: Sequence[str] | None = None,
    runner: CommandRunner = default_runner,
) -> str:
    if not agent_command.strip():
        agent_command = default_agent_command()

    root = root.resolve()
    runs_dir = (runs_dir or default_runs_dir(root)).resolve()
    check_command = check_command or [sys.executable, "-m", "pytest"]

    diagnostic(f"root: {root}")
    diagnostic(f"runs dir: {runs_dir}")

    with lock_file(root):
        diagnostic("lock acquired")
        diagnostic("checking repository")
        ensure_git_repo(root, runner)

        loaded_env = load_dotenv(root / ".env")
        if loaded_env:
            diagnostic(f"loaded .env keys: {', '.join(sorted(loaded_env))}")
        else:
            diagnostic(".env not found or empty")

        diagnostic("checking human input changes")
        if checkpoint_human_input(root, runner):
            diagnostic("human input checkpoint committed")

        ensure_clean_worktree(root, runner)
        branch_name = current_branch(root, runner)
        diagnostic(f"main branch: {branch_name}")

        session = run_session.read_counter(root / "state" / "session_counter.txt")
        diagnostic(f"session: {session_id(session)}")
        worktree: Path | None = None
        branch = ""
        applied = False

        try:
            diagnostic("creating temporary worktree")
            worktree, branch = create_worktree(root, runs_dir, session, runner)
            diagnostic(f"temporary branch: {branch}")

            diagnostic("running agent session")
            run_inner_session(worktree, agent_command, runner)

            diagnostic("checking required session files")
            ensure_required_session_files(worktree)

            diagnostic("checking file size policy")
            ensure_file_size_policy(worktree)

            diagnostic("running validation checks")
            run_checks(worktree, check_command, runner)

            diagnostic("checking produced changes")
            ensure_session_changed_worktree(worktree, runner)

            diagnostic("committing session changes")
            commit_hash = commit_session(worktree, session, runner)

            diagnostic("applying session commit to main worktree")
            apply_session_commit(root, branch, runner)
            applied = True
            diagnostic(f"applied commit: {commit_hash}")
            return commit_hash
        finally:
            if worktree is not None and branch:
                diagnostic("cleaning temporary worktree")
                remove_worktree_and_branch(
                    root,
                    worktree,
                    branch,
                    runner,
                    force_branch=not applied,
                )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Run one autonomous session in a temporary git worktree."
    )
    parser.add_argument(
        "--root",
        default=str(Path(__file__).resolve().parents[1]),
        help="Main project root. Defaults to the parent of scripts/.",
    )
    parser.add_argument(
        "--agent-command",
        default=os.environ.get("AI_AGENT_COMMAND", ""),
        help="Agent command passed to scripts/run_session.py inside the worktree. Defaults to uv-run mini-swe-agent wrapper.",
    )
    parser.add_argument(
        "--runs-dir",
        default="",
        help="Directory for temporary worktrees. Defaults to <root>/runs.",
    )
    parser.add_argument(
        "--check-command",
        nargs="+",
        default=[sys.executable, "-m", "pytest"],
        help="Command used to validate the worktree before applying changes.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    root = Path(args.root)
    runs_dir = Path(args.runs_dir) if args.runs_dir else None

    try:
        commit_hash = run_transaction(
            root=root,
            agent_command=args.agent_command,
            runs_dir=runs_dir,
            check_command=args.check_command,
            runner=streaming_runner,
        )
    except (TransactionError, CommandExecutionError) as exc:
        print(f"Transaction failed: {exc}", file=sys.stderr)
        return 1

    print(f"Transaction applied: {commit_hash}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
