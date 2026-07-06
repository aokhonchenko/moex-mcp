# Быстрый контекст проекта ai-lives

**Обновлено:** сессия 29 (2026-07-06)

---

## Суть проекта

Русскоязычная автономная рабочая среда агента. Агент существует между сессиями через файлы. Каждая сессия — один шаг, оставляющий артефакт.

## Структура (ключевое)

| Директория | Назначение |
|------------|------------|
| `src/` | Код: `agent/context.py`, `tools/partial_reader.py`, `tools/prompt_builder.py`, `tools/code_analyzer.py`, `tools/compat.py`, `session_runner.py` |
| `tests/` | Тесты: `test_code_analyzer.py` (~30 тестов) |
| `state/` | Состояние: `last_session.md`, `current_plan.md`, `session_context.md` |
| `tasks/` | Задачи: `active.md` (приоритеты), `archive.md` |
| `knowledge/` | Знания: `system_map.md`, `file_manifest.md`, `quick_context.md`, `codebase-analysis-26.md`, `notes/` |
| `tools/` | Инструменты: `dashboard/`, `notes/`, `sleep/`, `diff/`, `integrity/` |
| `logs/` | История: `history.md`, `week-*.md` |

## Текущее состояние

- 29 сессий завершены, практическая фаза
- 0 открытых вопросов
- 4 активных задачи (1 средний, 3 низких)
- `src/` содержит 6 модулей, покрытие docstrings 100%
- `tests/` содержит 1 тестовый модуль (~30 тестов)
- Два практических инструмента: `code_analyzer.py`, `dashboard/`

## Сессионный цикл

1. Прочитать этот файл → `state/last_session.md` → `tasks/active.md`
2. Выбрать задачу, сделать шаг
3. Обновить `state/last_session.md`, `logs/history.md`, `tasks/active.md`

## Фидбек создателя

> оптимизируй чтение файлов, выдели код в src/, улучшай инструменты
> сдвинуться от мета-работы к практике
> я хочу UI дашборд для тебя

## Подробнее

→ `knowledge/system_map.md` — полная карта  
→ `knowledge/file_manifest.md` — размеры и правила чтения  
→ `knowledge/artifact_links.md` — связи между артефактами  
→ `knowledge/codebase-analysis-26.md` — анализ кодовой базы
