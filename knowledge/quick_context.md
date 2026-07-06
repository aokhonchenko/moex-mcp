# Быстрый контекст проекта ai-lives

**Обновлено:** сессия 55 (2026-07-06)

---

## Суть проекта

Русскоязычная автономная рабочая среда агента. Агент существует между сессиями через файлы. Каждая сессия — один шаг, оставляющий артефакт.

## Структура (ключевое)

| Директория | Назначение |
|------------|------------|
| `src/tools/` | Инструменты агента: apply_patch, code_analyzer, command_runner, compat, partial_reader, prompt_builder, reader, read_file, read_lines, replace_text, run_command, run_pytest, run_python_script, self_review, write_file |
| `src/agent/` | Логика агента: context.py |
| `tests/` | Тесты: test_apply_patch, test_code_analyzer, test_command_runner, test_compat, test_partial_reader, test_reader, test_self_review |
| `state/` | Состояние: last_session.md, current_plan.md, session_context.md, questions/ |
| `tasks/` | Задачи: active.md, archive.md |
| `knowledge/` | Знания: system_map.md, file_manifest.md, quick_context.md, codebase-analysis-26.md, notes/ |
| `tools/` | Markdown-инструменты: dashboard/, notes/, sleep/, diff/, integrity/ |
| `logs/` | История: history.md, week-*.md, archive/ |
| `projects/` | Мини-проекты: TEMPLATE.md, task-tracker/, **foundation-finance/** |

## Текущее состояние

- 55 сессий завершено, практическая фаза
- 0 открытых вопросов
- 2 активных задачи (command_runner интеграция, foundation-finance)
- `src/tools/` содержит 15 инструментов (все с runtime tool.py)
- `tests/` содержит 7 тестовых модулей (~290 тестов)
- **foundation-finance** — финансовый дашборд (Go backend + Web frontend + Docker Compose)
  - MOEX ISS API, кэширование, фундаментальные данные, LLM, свечной график
  - Поиск по тикерам с автокомплитом, zoom/pan, кроссхейр
  - ~117 Go unit-тестов
  - Репозиторий: git@github.com:aokhonchenko/foundation-finance.git

## Сессионный цикл

1. Прочитать этот файл → `state/last_session.md` → `tasks/active.md`
2. Выбрать задачу, сделать шаг
3. Обновить `state/last_session.md`, `logs/history.md`, `tasks/active.md`

## Фидбек создателя

> оптимизируй чтение файлов, выдели код в src/, улучшай инструменты
> сдвинуться от мета-работы к практике
> я хочу UI дашборд для тебя
> для foundation-finance нужны российские источники (MOEX)

## Подробнее

→ `knowledge/system_map.md` — полная карта  
→ `knowledge/file_manifest.md` — размеры и правила чтения  
→ `knowledge/artifact_links.md` — связи между артефактами
