"""Replication helpers for external project repositories under projects/."""

from __future__ import annotations

import shutil
from pathlib import Path
from typing import Callable, Sequence

from scripts.command_runners import CommandResult, default_runner


class ExternalProjectError(RuntimeError):
    """Raised when an external project cannot be replicated."""


CommandRunner = Callable[[Sequence[str], Path], CommandResult]

EXTERNAL_PROJECTS_DIR = "projects"
EXTERNAL_PROJECT_COPY_IGNORES = {"__pycache__", ".pytest_cache", ".venv", "node_modules", "target"}


def is_git_ignored(root: Path, path: Path, runner: CommandRunner = default_runner) -> bool:
    try:
        relative = path.resolve().relative_to(root.resolve()).as_posix()
    except ValueError:
        return False
    result = runner(["git", "check-ignore", "-q", relative], root)
    return result.returncode == 0


def has_git_metadata(path: Path) -> bool:
    return (path / ".git").exists()


def external_project_dirs(root: Path, runner: CommandRunner = default_runner) -> list[Path]:
    projects_dir = root / EXTERNAL_PROJECTS_DIR
    if not projects_dir.exists():
        return []
    candidates: list[Path] = []
    for path in sorted(projects_dir.iterdir(), key=lambda item: item.name):
        if not path.is_dir():
            continue
        if has_git_metadata(path) or is_git_ignored(root, path, runner):
            candidates.append(path)
    return candidates


def copy_external_project(source: Path, target: Path) -> None:
    ignore = shutil.ignore_patterns(*EXTERNAL_PROJECT_COPY_IGNORES)
    if target.exists() and not target.is_dir():
        raise ExternalProjectError(f"external project target is not a directory: {target}")
    try:
        shutil.copytree(source, target, dirs_exist_ok=True, ignore=ignore)
    except (OSError, shutil.Error) as exc:
        raise ExternalProjectError(f"cannot replicate external project {source} to {target}: {exc}") from exc


def replicate_external_projects(
    worktree: Path,
    root: Path,
    runner: CommandRunner = default_runner,
) -> list[Path]:
    replicated: list[Path] = []
    for source in external_project_dirs(worktree, runner):
        relative = source.relative_to(worktree)
        target = root / relative
        copy_external_project(source, target)
        replicated.append(relative)
    return replicated
