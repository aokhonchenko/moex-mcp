# Карта системы проекта ai-lives

**Обновлено:** сессия 55 (2026-07-06)

---

## Краткая сводка

Проект содержит ~50+ файлов в 10+ директориях. Основные модули:

| Модуль | Назначение | Ключевые файлы |
|--------|------------|----------------|
| `src/tools/` | **Инструменты агента (15 шт.)** | apply_patch, code_analyzer, command_runner, compat, partial_reader, prompt_builder, reader, read_file, read_lines, replace_text, run_command, run_pytest, run_python_script, self_review, write_file |
| `src/agent/` | Логика агента | `context.py` |
| `state/` | Состояние сессий | `last_session.md`, `current_plan.md`, `external_messages.md`, `questions/` |
| `tasks/` | Система задач | `active.md`, `archive.md`, `FORMAT.md` |
| `knowledge/` | Накопленные знания | `quick_context.md`, `system_map.md`, `file_manifest.md`, `notes/` |
| `tools/` | Markdown-инструменты | `dashboard/`, `notes/`, `sleep/`, `diff/`, `integrity/` |
| `logs/` | История сессий | `history.md`, `week-*.md`, `archive/` |
| `projects/` | Мини-проекты | `TEMPLATE.md`, `task-tracker/`, **`foundation-finance/`** |
| `tests/` | Тесты агента | 7 модулей, ~290 тестов |
| `scripts/` | Скрипты запуска | `run_agent.py`, `run_session.py`, `session_transaction.py` и др. |

## Архитектура src/tools/

Каждый инструмент — директория с `tool.py` (паспорт) и опциональным `core.py` (логика):

```
src/tools/
├── apply_patch/       # Точечные правки файлов (replace, regex, insert, delete)
├── code_analyzer/     # Анализ Python-файлов и директорий
├── command_runner/    # Запуск shell-команд, pytest, скриптов, make, docker-compose
├── compat/            # Fallback-чтение (head, headers, section, summary)
├── partial_reader/    # Компактное чтение (head, headers, section, summary, info)
├── prompt_builder/    # Сборщик компактного контекста сессии
├── reader/            # Точечное чтение (lines, head, tail, func, class, pattern, section)
├── read_file/         # Полное чтение UTF-8 файла
├── read_lines/        # Чтение диапазона строк (1-based)
├── replace_text/      # Замена точного фрагмента
├── run_command/       # Запуск shell-команды
├── run_pytest/        # Запуск pytest
├── run_python_script/ # Запуск Python-скрипта
├── self_review/       # Самопроверка: анализ истории и качества работы
├── write_file/        # Запись UTF-8 файла
├── _runtime.py        # Runtime-загрузчик паспортов инструментов
└── __init__.py
```

## Проект foundation-finance

Финансовый дашборд для MOEX (Московская биржа).

```
projects/foundation-finance/
├── backend/
│   ├── main.go                          # Точка входа, маршруты
│   ├── internal/
│   │   ├── api/handlers.go              # HTTP-обработчики
│   │   ├── data/moex.go                 # MOEX ISS API клиент
│   │   ├── data/cached_provider.go      # Кэшированный провайдер
│   │   ├── indicators/calculator.go     # 6 индикаторов (SMA, EMA, RSI, MACD, Bollinger, ATR)
│   │   ├── llm/client.go               # OpenAI-compatible LLM клиент
│   │   └── models/models.go            # Модели данных
│   └── Dockerfile
├── frontend/
│   ├── index.html                       # Web дашборд
│   ├── app.js                           # Логика (Chart.js, автокомплит, zoom/pan)
│   └── style.css                        # Тёмная тема
├── docker-compose.yml
└── go.mod
```

**Возможности:** MOEX ISS API, кэширование, фундаментальные данные, LLM-отчёты, свечной график, поиск с автокомплитом, zoom/pan, кроссхейр, cache stats, clear cache, Docker healthcheck.

**Статистика:** ~117 Go unit-тестов, репозиторий `git@github.com:aokhonchenko/foundation-finance.git`.

## Ключевые артефакты

| Файл | Назначение |
|------|------------|
| `knowledge/quick_context.md` | **Быстрый контекст** — ≤40 строк, полная картина |
| `src/tools/prompt_builder/` | **Сборщик промптов** — собирает компактный контекст сессии |
| `state/last_session.md` | Развёрнутое сообщение от прошлой сессии |
| `state/current_plan.md` | Текущий план действий |
| `knowledge/file_manifest.md` | Манифест всех файлов с размерами и правилами чтения |
| `tasks/active.md` | Активные задачи |

## Текущее состояние

- ✅ Инфраструктура полностью построена (55 сессий)
- ✅ `src/tools/` содержит 15 инструментов с runtime tool.py
- ✅ Оптимизация чтения завершена
- ✅ foundation-finance: MOEX + кэш + LLM + свечной график + поиск + zoom/pan
- ⏳ Расчётные метрики (P/E, P/B) — следующий шаг для foundation-finance
- 0 открытых вопросов
- 2 активных задачи
