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
AI_API_KEY=replace-with-your-api-key
AI_BASE_URL=https://your-openai-compatible-endpoint.example/v1
```

Выберите модель в `config/project.toml`:

```toml
[mini_swe_agent]
model = "openai/your-model-name"
custom_llm_provider = "openai"
```

`AI_API_KEY` и `AI_BASE_URL` не коммитятся. `config/project.toml` коммитится, потому что это проектная настройка.

## Ручной цикл

1. При необходимости обновите `GLOBAL_TARGET.md`, `state/external_messages.md` или ответьте на вопросы в `state/questions/*.md`.
2. Запустите транзакционную сессию вручную:

```bash
uv run python scripts/session_transaction.py
```

3. Просмотрите результат сессии и новые вопросы агента.
4. Повторите цикл, когда хотите дать агенту следующий ход.

Если перед запуском изменены только `GLOBAL_TARGET.md`, `state/external_messages.md` или `state/questions/*.md`, runner сам сделает checkpoint-коммит этих человеческих входов. Если изменены другие файлы, запуск будет остановлен.

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

Проверить wrapper без реального запуска модели нельзя без валидных `AI_API_KEY`, `AI_BASE_URL` и рабочей модели. Ошибки подключения будут откатываться транзакционным runner'ом.

## Тесты

```bash
uv run python -m pytest
```

Порог покрытия задан в `pyproject.toml`: 90%.