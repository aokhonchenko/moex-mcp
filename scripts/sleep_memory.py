#!/usr/bin/env python3
"""Sleep memory maintenance tool called by the agent during a normal session."""

from __future__ import annotations

import argparse
import re
import shutil
import sys
from datetime import datetime, timezone
from pathlib import Path


CLOSED_STATUS_RE = re.compile(r"^\s*(?:status|статус)\s*:\s*(?:closed|закрыт|закрыто)\s*$", re.IGNORECASE | re.MULTILINE)
ANSWERED_STATUS_RE = re.compile(r"^\s*(?:status|статус)\s*:\s*(?:answered|отвечен|отвечено)\s*$", re.IGNORECASE | re.MULTILINE)
OPEN_STATUS_RE = re.compile(r"^\s*(?:status|статус)\s*:\s*(?:open|открыт|открыто)\s*$", re.IGNORECASE | re.MULTILINE)


def question_files(root: Path) -> list[Path]:
    questions_dir = root / "state" / "questions"
    if not questions_dir.exists():
        return []
    return sorted(
        path
        for path in questions_dir.glob("*.md")
        if path.is_file() and path.name != "README.md"
    )


def classify_question(path: Path) -> str:
    text = path.read_text(encoding="utf-8")
    if CLOSED_STATUS_RE.search(text):
        return "closed"
    if ANSWERED_STATUS_RE.search(text):
        return "answered"
    if OPEN_STATUS_RE.search(text):
        return "open"
    return "unknown"


def archive_closed_questions(root: Path, now: datetime) -> list[Path]:
    archived: list[Path] = []
    archive_dir = root / "state" / "questions" / "archive" / now.strftime("%Y-%m-%d")

    for path in question_files(root):
        if classify_question(path) != "closed":
            continue
        archive_dir.mkdir(parents=True, exist_ok=True)
        target = archive_dir / path.name
        if target.exists():
            stem = path.stem
            suffix = path.suffix
            target = archive_dir / f"{stem}-{now.strftime('%H%M%S')}{suffix}"
        shutil.move(str(path), str(target))
        archived.append(target)

    return archived


def build_sleep_report(root: Path, archived: list[Path], now: datetime) -> str:
    remaining = [(path, classify_question(path)) for path in question_files(root)]
    open_questions = [path for path, status in remaining if status == "open"]
    answered_questions = [path for path, status in remaining if status == "answered"]
    unknown_questions = [path for path, status in remaining if status == "unknown"]

    def format_paths(paths: list[Path]) -> str:
        if not paths:
            return "- нет\n"
        return "".join(f"- `{path.relative_to(root)}`\n" for path in paths)

    return f"""# Отчёт сна

Дата: {now.strftime('%Y-%m-%d %H:%M:%S %z')}

## Что очищено

Закрытые вопросы перенесены в архив:

{format_paths(archived)}
## Что осталось активным

Открытые вопросы:

{format_paths(open_questions)}
Вопросы с ответом, которые ещё нужно учесть или закрыть:

{format_paths(answered_questions)}
Вопросы без распознанного статуса:

{format_paths(unknown_questions)}
## Итог

Сон завершён. Память очищена от закрытых вопросов, активные вопросы оставлены на месте.
"""


def write_sleep_artifacts(root: Path, report: str, now: datetime) -> Path:
    sleep_dir = root / "state" / "sleep"
    reports_dir = sleep_dir / "reports"
    reports_dir.mkdir(parents=True, exist_ok=True)

    report_path = reports_dir / f"{now.strftime('%Y-%m-%d-%H%M%S')}.md"
    report_path.write_text(report, encoding="utf-8")
    (sleep_dir / "last_sleep.md").write_text(report, encoding="utf-8")

    last_session = root / "state" / "last_session.md"
    last_session.write_text(
        "# Сообщение будущей сессии\n\n"
        "Последняя сессия была сном: закрытые вопросы архивированы, активные вопросы оставлены в `state/questions/`. "
        "Начни следующую рабочую сессию с чтения `GLOBAL_TARGET.md`, `state/current_plan.md` и `state/sleep/last_sleep.md`.\n",
        encoding="utf-8",
    )

    history = root / "logs" / "history.md"
    with history.open("a", encoding="utf-8") as handle:
        handle.write(
            f"\n## Сон - {now.strftime('%Y-%m-%d %H:%M:%S %z')}\n\n"
            f"- Отчёт: `{report_path.relative_to(root)}`\n"
            "- Закрытые вопросы перенесены в архив.\n"
        )

    return report_path


def run_sleep(root: Path) -> Path:
    now = datetime.now(timezone.utc).astimezone()
    archived = archive_closed_questions(root, now)
    report = build_sleep_report(root, archived, now)
    return write_sleep_artifacts(root, report, now)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run ai-lives sleep memory maintenance inside a normal session.")
    parser.add_argument("--root", required=True, help="Session worktree root.")
    parser.add_argument("--prompt-file", default="", help="Accepted for command compatibility; not used.")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    root = Path(args.root).resolve()
    try:
        report_path = run_sleep(root)
    except Exception as exc:
        print(f"sleep session failed: {exc}", file=sys.stderr)
        return 1

    print(f"Sleep session complete: {report_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())