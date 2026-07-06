from pathlib import Path

from scripts import sleep_memory


def write_question(path: Path, status: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(
        f"# Вопрос\n\nСтатус: {status}\n\n## Ответ создателя\n\nТекст.\n",
        encoding="utf-8",
    )


def test_classify_question_detects_statuses(tmp_path):
    open_question = tmp_path / "open.md"
    answered_question = tmp_path / "answered.md"
    closed_question = tmp_path / "closed.md"
    unknown_question = tmp_path / "unknown.md"
    write_question(open_question, "open")
    write_question(answered_question, "answered")
    write_question(closed_question, "closed")
    unknown_question.write_text("# Без статуса\n", encoding="utf-8")

    assert sleep_memory.classify_question(open_question) == "open"
    assert sleep_memory.classify_question(answered_question) == "answered"
    assert sleep_memory.classify_question(closed_question) == "closed"
    assert sleep_memory.classify_question(unknown_question) == "unknown"


def test_archive_closed_questions_moves_only_closed_questions(tmp_path):
    root = tmp_path
    questions = root / "state" / "questions"
    write_question(questions / "0001-open.md", "open")
    write_question(questions / "0002-answered.md", "answered")
    write_question(questions / "0003-closed.md", "closed")

    archived = sleep_memory.archive_closed_questions(root, sleep_memory.datetime(2026, 7, 6, tzinfo=sleep_memory.timezone.utc))

    assert len(archived) == 1
    assert archived[0].name == "0003-closed.md"
    assert not (questions / "0003-closed.md").exists()
    assert (questions / "0001-open.md").exists()
    assert (questions / "0002-answered.md").exists()


def test_run_sleep_writes_report_history_and_last_session(tmp_path):
    root = tmp_path
    (root / "logs").mkdir()
    (root / "logs" / "history.md").write_text("# История\n", encoding="utf-8")
    (root / "state").mkdir()
    write_question(root / "state" / "questions" / "0001-closed.md", "closed")

    report_path = sleep_memory.run_sleep(root)

    assert report_path.exists()
    assert (root / "state" / "sleep" / "last_sleep.md").exists()
    assert "Сон завершён" in (root / "state" / "sleep" / "last_sleep.md").read_text(encoding="utf-8")
    assert "Последняя сессия была сном" in (root / "state" / "last_session.md").read_text(encoding="utf-8")
    assert "Закрытые вопросы перенесены" in (root / "logs" / "history.md").read_text(encoding="utf-8")