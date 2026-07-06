# Текущий план


## Статус: практическая фаза (сессия 40)

- ✅ Создана карта системы (`knowledge/system_map.md`) — сессия 1
- ✅ Обновлены файлы состояния и истории — сессия 1
- ✅ Создан шаблон вопроса создателю — сессия 2
- ✅ Правила именования файлов в `knowledge/` — сессия 3
- ✅ Прототип инструмента сна — сессия 4
- ✅ Чеклист пробуждения — сессия 8
- ✅ Мини-проект «Идеи для улучшений» + инструмент проверки целостности — сессия 7
- ✅ Шаблон мини-проекта — сессия 8
- ✅ Архив закрытых вопросов — сессия 9
- ✅ Карта связей между артефактами — сессия 10
- ✅ Инструмент дифф-отчёта — сессия 11
- ✅ Первый сон — сессия 12
- ✅ Система задач — сессия 13
- ✅ Интеграция задач в сессионный цикл — сессия 14
- ✅ Система заметок — сессия 15
- ✅ Стратегия ленивого чтения + инструмент чтения заголовков — сессия 16
- ✅ Реструктуризация `logs/history.md` — сессия 17
- ✅ Компактный контекст сессии (`state/session_context.md`) — сессия 18
- ✅ Манифест файлов (`knowledge/file_manifest.md`) — сессия 18
- ✅ Обновлён чеклист сессии (переход на контекст) — сессия 18
- ✅ Оптимизирована карта системы (краткая сводка + отдельное дерево) — сессия 18
- ✅ Оценщик чтения (`tools/reading-analyzer/`) — сессия 19
- ✅ Обновлена карта связей — сессия 20
- ✅ **Создан `src/` + `partial_reader.py`** — сессия 21
- ✅ **Создан `src/agent/context.py`** — сессия 22
- ✅ **Создан `knowledge/quick_context.md`** — сессия 23
- ✅ **Оптимизирован `knowledge/artifact_links.md`** — сессия 24
- ✅ **Создан `src/tools/prompt_builder.py`** — сессия 25
- ✅ **Создан `src/session_runner.py`** — сессия 25
- ✅ **Создан `src/tools/code_analyzer.py`** — сессия 26
- ✅ **Первый анализ кодовой базы** — сессия 26
- ✅ **Создан UI-дашборд** (`tools/dashboard/`) — сессия 27
- ✅ **Устранено дублирование fallback-функций** (`src/tools/compat.py`) — сессия 28
- ✅ **Первые тесты проекта** (`tests/test_code_analyzer.py`) — сессия 29
- ✅ **Исправлен тест `test_analyze_self`** — сессия 30
- ✅ **Тесты для `compat.py` и `partial_reader.py`** — сессия 31
- ✅ **Точечное чтение файлов** (`src/tools/reader.py` + `tests/test_reader.py`) — сессия 32
- ✅ **Исправлены баги `reader.py`** — 0-based/1-based индексы, добавлено поле `error` в `ReadResult` — сессия 33
- ✅ **Закрыт вопрос о структуре проекта** — создатель подтвердил текущую структуру — сессия 34
- ✅ **Закрыт вопрос о приоритетах** — приоритет: улучшение агента — сессия 35
- ✅ **Создан инструмент частичных правок** (`src/tools/apply_patch.py`) — сессия 35
- ✅ **Созданы тесты для apply_patch.py** (`tests/test_apply_patch.py`) — сессия 35
- ✅ **Исправлен баг replace_regex** — добавлен `re.MULTILINE` — сессия 36
- ✅ **Создан модуль self-review** (`src/tools/self_review.py`) — сессия 36
- ✅ **Созданы тесты для self-review** (`tests/test_self_review.py`) — сессия 36
- ✅ **Создан инструмент запуска команд** (`src/tools/command_runner.py`) — сессия 37
- ✅ **Созданы тесты для command_runner** (`tests/test_command_runner.py`) — сессия 37

## Завершённые мини-проекты

1. **Система задач** ✅ — `tasks/`, `projects/task-tracker/`
2. **Система заметок** ✅ — `tools/notes/`, `knowledge/notes/`
3. **Оптимизация чтения (полная)** ✅ — `tools/file-headers/`, `state/session_context.md`, `knowledge/file_manifest.md`, `tools/reading-analyzer/`, `knowledge/quick_context.md`, реструктуризация больших файлов
4. **Инструмент частичного чтения** ✅ — `src/tools/partial_reader.py`
5. **Модуль управления контекстом** ✅ — `src/agent/context.py`
6. **Быстрый контекст** ✅ — `knowledge/quick_context.md`
7. **Сборщик промптов** ✅ — `src/tools/prompt_builder.py`, `src/session_runner.py`
8. **Анализатор кода** ✅ — `src/tools/code_analyzer.py`, `knowledge/codebase-analysis-26.md`
9. **UI-дашборд** ✅ — `tools/dashboard/generate.py`, `tools/dashboard/index.html`
10. **Устранение дублирования** ✅ — `src/tools/compat.py`
11. **Первые тесты** ✅ — `tests/test_code_analyzer.py`
12. **Тесты для compat/partial_reader** ✅ — `tests/test_compat.py`, `tests/test_partial_reader.py`
13. **Точечное чтение файлов** ✅ — `src/tools/reader.py`, `tests/test_reader.py`
14. **Исправление багов reader.py** ✅ — 0-based/1-based индексы, поле `error` — сессия 33
15. **Закрытие вопроса о структуре** ✅ — структура подтверждена — сессия 34
16. **Закрытие вопроса о приоритетах** ✅ — приоритет: улучшение агента — сессия 35
17. **Инструмент частичных правок** ✅ — `src/tools/apply_patch.py`, `tests/test_apply_patch.py` — сессия 35
18. **Модуль self-review** ✅ — `src/tools/self_review.py`, `tests/test_self_review.py` — сессия 36
19. **Инструмент запуска команд** ✅ — `src/tools/command_runner.py`, `tests/test_command_runner.py` — сессия 37

## Текущий фокус: улучшение агента

Приоритет определён создателем: улучшение агента. Финансовый дашборд — когда задачи в этой области кончатся.

- ✅ Точечное чтение файлов (`reader.py`, `partial_reader.py`)
- ✅ Инструмент частичных правок (`apply_patch.py`)
- ✅ Модуль self-review (`self_review.py`)
- ✅ Инструмент запуска команд (`command_runner.py`)
- ✅ Интеграция self_review.py в сессионный цикл — запущен и сохранён отчёт в `state/self_review/`
- ✅ Запуск тестов через command_runner — все 275 тестов прошли, включая `test_regex_multiline`
- ✅ **Интеграция точечного чтения в сессионный цикл** — `prompt_builder.py` использует `partial_reader.py` для оптимизированного чтения
- ✅ Интеграция `apply_patch.py` в сессионный цикл (частичные правки вместо полной перезаписи) — сессия 41
- ⏳ Интеграция `command_runner.py` в сессионный цикл (кроме тестов)

## Следующий разумный шаг (сессия 42)

Задачи по улучшению агента в основном завершены. Два варианта:
1. **Интегрировать `command_runner` в сессионный цикл** — если есть конкретные сценарии использования
2. **Начать работу над финансовым дашбордом `foundation-finance`** — создать директорию в `projects/`, подключить git-репозиторий, выбрать стек (Rust/Go + frontend)

- ✅ **Создан финансовый дашборд `foundation-finance`** — сессия 42
  - Go backend (chi, Yahoo Finance, 6 индикаторов, LLM-клиент)
  - Web frontend (Chart.js, тёмная тема)
  - Docker Compose
  - Репозиторий: `git@github.com:aokhonchenko/foundation-finance.git`

- ✅ **Unit-тесты для `indicators/calculator.go`** — 26 тестов, все проходят — сессия 43

- ✅ **Замена Yahoo Finance на MOEX ISS API** — сессия 44
  - Создан `data/moex.go` (MOEX ISS провайдер)
  - 15 unit-тестов для data/moex.go (мок-сервер)
  - 10 unit-тестов для api/handlers.go (мок-провайдер + chi-роутер)
  - main.go переключён на MOEXProvider
  - Фронтенд обновлён: тикеры MOEX (SBER, GAZP, LKOH)

- ✅ **In-memory кэширование MOEX данных** — сессия 45
  - `data/cache.go` — потокобезопасный кэш с TTL, max size, auto-cleanup
  - `data/cached_provider.go` — декоратор Provider с кэшированием
  - 23 новых unit-теста (11 cache + 12 cached_provider)
  - main.go: ticker TTL 2 мин, candles TTL 15 мин
  - Всего Go тестов: 73, Python тестов: 290

## Следующий разумный шаг (сессия 46)

1. **Добавить фундаментальные индикаторы** (P/E, P/B, ROE, дивиденды) — MOEX ISS не отдаёт мультипликаторы, нужен другой источник
2. **Добавить LLM тесты** — мок-сервер для OpenAI-compatible API
3. **Улучшить фронтенд** — свечной график, таблица фундаментальных метрик
4. **Docker Compose тест** — проверить, что `docker-compose up` работает
5. **API endpoint для статистики кэша** — `/api/cache/stats`
