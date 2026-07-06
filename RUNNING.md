# Переносимый запуск

Основной runner эксперимента: `scripts/run_session.py`.

Проверить сборку промпта без запуска агента:

```bash
python scripts/run_session.py --dry-run
```

Запустить одну сессию через внешний CLI-агент:

```bash
export AI_AGENT_COMMAND='your-agent --cwd "{ROOT}" --prompt-file "{PROMPT_FILE}"'
python scripts/run_session.py
```

Плейсхолдеры в `AI_AGENT_COMMAND`:

- `{ROOT}` - корень эксперимента;
- `{PROMPT_FILE}` - собранный файл `state/active_prompt.md`;
- `{SESSION}` - номер текущей сессии.

Если агент читает промпт из stdin:

```bash
export AI_AGENT_COMMAND='cat "{PROMPT_FILE}" | your-agent --cwd "{ROOT}"'
python scripts/run_session.py
```

Для автоматизации используйте любой планировщик: `cron`, `systemd timer`, Task Scheduler, launchd, CI-расписание или внешний оркестратор. Скрипт запуска остаётся Python-скриптом и не требует привязки к конкретной операционной системе.

Новые изменения должны ориентироваться на scripts/session_transaction.py для автономного режима и на scripts/run_session.py для отладки.

## Тесты

Проверка тестов и покрытия:

```bash
python -m pytest
```

Порог покрытия задан в `pyproject.toml`: 90%.

## Атомарный запуск

Для реального автономного режима используйте транзакционную оболочку:

```bash
python scripts/session_transaction.py --agent-command 'your-agent --cwd "{ROOT}" --prompt-file "{PROMPT_FILE}"'
```

Что делает `session_transaction.py`:

1. Проверяет, что основной проект является Git-репозиторием.
2. Проверяет, что основной worktree чистый.
3. Берёт `.session.lock`, чтобы две сессии не шли одновременно.
4. Создаёт временный `git worktree` в соседней директории `<project>-runs/session-NNNN`.
5. Запускает `scripts/run_session.py` внутри временного worktree.
6. Запускает проверки, по умолчанию `python -m pytest`.
7. Если всё успешно, коммитит изменения сессии и применяет их в основной проект через `git merge --ff-only`.
8. Если агент или проверки падают, удаляет временный worktree и ветку сессии. Основная директория остаётся в состоянии до сессии.

Обычный `scripts/run_session.py` полезен для отладки и ручных запусков, но он пишет прямо в текущую директорию. Для свойства "или сессия завершилась целиком, или её будто не было" используйте только `scripts/session_transaction.py`.

Перед первым транзакционным запуском проект должен быть Git-репозиторием с базовым коммитом:

```bash
git init
git add -A
git commit -m "initial autonomous experiment scaffold"
```
