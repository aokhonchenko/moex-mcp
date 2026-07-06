# Переносимый ручной запуск

Проект использует `uv` и `mini-swe-agent`.

## Сборка окружения

```bash
uv sync
```

`mini-swe-agent` установлен как обычная зависимость проекта и зафиксирован в `uv.lock`.

## Настройка модели

Скопируйте пример окружения:

```bash
cp .env.example .env
```

Заполните `.env`:

```bash
AI_API_KEY=
AI_BASE_URL=https://your-openai-compatible-endpoint.example/v1
AI_MODEL=openai/your-model-name
```

Модель, URL и опциональный ключ находятся в `.env` и не коммитятся. `config/project.toml` коммитится и хранит только не-секретные параметры запуска `mini-swe-agent`.

## Ручной цикл

1. При необходимости обновите `GLOBAL_TARGET.md`, `state/external_messages.md` или ответьте на вопросы в `state/questions/*.md`.
2. Запустите транзакционную сессию вручную:

```bash
uv run python scripts/session_transaction.py
```

3. Просмотрите результат сессии и новые вопросы агента.
4. Повторите цикл, когда хотите дать агенту следующий ход.

Если перед запуском изменены только `GLOBAL_TARGET.md`, `state/external_messages.md` или `state/questions/*.md`, runner сам сделает checkpoint-коммит этих человеческих входов. Если изменены другие файлы, запуск будет остановлен.


## Сон

Сон не запускается отдельным внешним режимом. Вы запускаете обычную сессию:

```bash
uv run python scripts/session_transaction.py
```

Если агент решит, что накопилась усталость или нужно очистить память, он сам вызовет `uv run python scripts/sleep_memory.py --root .` внутри этой обычной сессии. Результат сна сохраняется в `state/sleep/last_sleep.md` и `state/sleep/reports/`.

## Что запускается внутри

По умолчанию `scripts/session_transaction.py` вызывает:

```bash
uv run python scripts/run_mini_agent.py --root "{ROOT}" --prompt-file "{PROMPT_FILE}"
```

`run_mini_agent.py` использует Python API `mini-swe-agent`, а не интерактивный CLI `mini`, потому что транзакционный runner выполняет агента non-interactively.

Если нужно заменить агента, передайте команду явно:

```bash
uv run python scripts/session_transaction.py --agent-command 'your-agent --cwd "{ROOT}" --prompt-file "{PROMPT_FILE}"'
```

## Отладка

Проверить сборку prompt без запуска модели:

```bash
uv run python scripts/run_session.py --dry-run
```

Проверить wrapper без реального запуска модели нельзя без `AI_BASE_URL`, `AI_MODEL` и рабочей модели. `AI_API_KEY` может быть пустым, если endpoint это допускает; wrapper подставит технический dummy credential для LiteLLM/OpenAI SDK. Ошибки подключения будут откатываться транзакционным runner'ом.

## Тесты

```bash
uv run python -m pytest
```

Порог покрытия задан в `pyproject.toml`: 90%.