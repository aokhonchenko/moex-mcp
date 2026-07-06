# Переносимый запуск

Основной ручной запуск для автономного режима: `scripts/session_transaction.py`. `scripts/run_session.py` оставлен для отладки сборки промпта.

Проверить сборку промпта без запуска агента:

```bash
python scripts/run_session.py --dry-run
```

Отладочно запустить одну сессию без транзакционной оболочки:

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

Сессии предполагается запускать вручную. Для автономного режима используйте `scripts/session_transaction.py`, для отладки сборки промпта - `scripts/run_session.py`.


## Ручной цикл

1. При необходимости обновите `GLOBAL_TARGET.md`, `state/external_messages.md` или ответьте на вопросы в `state/questions/*.md`.
2. Запустите транзакционную сессию вручную:

```bash
python scripts/session_transaction.py --agent-command 'your-agent --cwd "{ROOT}" --prompt-file "{PROMPT_FILE}"'
```

3. Просмотрите результат сессии и новые вопросы агента.
4. Повторите цикл, когда хотите дать агенту следующий ход.

Если перед запуском изменены только `GLOBAL_TARGET.md`, `state/external_messages.md` или `state/questions/*.md`, runner сам сделает checkpoint-коммит этих человеческих входов. Если изменены другие файлы, запуск будет остановлен.

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
