# Быстрый контекст проекта ai-lives

**Обновлено:** сессия 33 (2026-07-06)

---

## Суть проекта

Русскоязычная автономная рабочая среда агента. Агент существует между сессиями через файлы. Каждая сессия — один шаг, оставляющий артефакт.

## Структура (ключевое)

| Директория | Назначение |
|------------|------------|
| `src/` | Код: `agent/context.py`, `tools/partial_reader.py`, `tools/prompt_builder.py`, `tools/code_analyzer.py`, `tools/compat.py`, `tools/reader.py`, `session_runner.py` |
| `tests/` | Тесты: `test_code_analyzer.py`, `test_compat.py`, `test_partial_reader.py`, `test_reader.py` |
| `state/` | Состояние: `last_session.md`, `current_plan.md`, `session_context.md` |
| `tasks/` | Задачи: `active.md` (приоритеты), `archive.md` |
| `knowledge/` | Знания: `system_map.md`, `file_manifest.md`, `quick_context.md`, `codebase-analysis-26.md`, `notes/` |
| `tools/` | Инструменты: `dashboard/`, `notes/`, `sleep/`, `diff/`, `integrity/` |
| `logs/` | История: `history.md`, `week-*.md` |

## Текущее состояние

- 33 сессии завершены, практическая фаза
- 1 открытый вопрос (`state/questions/0032-project-structure.md`)
- 2 активных задачи (1 средний, 1 низкий)
- `src/` содержит 7 модулей, покрытие docstrings 100%
- `tests/` содержит 4 тестовых модуля (~90+ тестов)
- Практические инструменты: `code_analyzer.py`, `reader.py`, `dashboard/`

## Сессионный цикл

1. Прочитать этот файл → `state/last_session.md` → `tasks/active.md`
2. Выбрать задачу, сделать шаг
3. Обновить `state/last_session.md`, `logs/history.md`, `tasks/active.md`

## Фидбек создателя

> оптимизируй чтение файлов, выдели код в src/, улучшай инструменты
> сдвинуться от мета-работы к практике
> я хочу UI дашборд для тебя

## Новый инструмент: точечное чтение (`src/tools/reader.py`)

Создан для решения проблемы неоптимального чтения больших файлов целиком.

**Возможности:**
- `--lines START END` — чтение диапазона строк (1-based)
- `--head N` — первые N строк
- `--tail N` — последние N строк
- `--func ИМЯ` — определение функции из Python-файла (с декораторами)
- `--class ИМЯ` — определение класса из Python-файла
- `--pattern REGEX` — строки по регулярному выражению с контекстом
- `--section ИМЯ` — секция markdown по заголовку `## ИМЯ`
- `--info` — метаданные файла без чтения содержимого

**Исправления сессии 33:**
- Исправлены 0-based/1-based индексы в `read_func()` и `read_class()`
- Добавлено поле `error: Optional[str]` в `ReadResult`
- Все функции чтения теперь заполняют `error` при ошибках

**Примеры использования:**
```bash
# Прочитать только функцию analyze_file без загрузки всего файла
python src/tools/reader.py src/tools/code_analyzer.py --func analyze_file

# Прочитать секцию плана без загрузки всего плана
python src/tools/reader.py state/current_plan.md --section 'Следующий разумный шаг'

# Посмотреть метаданные файла
python src/tools/reader.py src/tools/code_analyzer.py --info
```

## Подробнее

→ `knowledge/system_map.md` — полная карта  
→ `knowledge/file_manifest.md` — размеры и правила чтения  
→ `knowledge/artifact_links.md` — связи между артефактами  
→ `knowledge/codebase-analysis-26.md` — анализ кодовой базы
