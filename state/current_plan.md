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

## Текущий фокус: финансовый дашборд foundation-finance

Приоритет определён создателем: улучшение агента (завершено) → финансовый дашборд.

### Завершённые шаги foundation-finance

- ✅ **Создан финансовый дашборд `foundation-finance`** — сессия 42
  - Go backend (chi, Yahoo Finance, 6 индикаторов, LLM-клиент)
  - Web frontend (Chart.js, тёмная тема)
  - Docker Compose
  - Репозиторий: `git@github.com:aokhonchenko/foundation-finance.git`

- ✅ **Unit-тесты для `indicators/calculator.go`** — 26 тестов — сессия 43

- ✅ **Замена Yahoo Finance на MOEX ISS API** — сессия 44
  - 15 unit-тестов для data/moex.go, 10 для api/handlers.go
  - Фронтенд: тикеры MOEX (SBER, GAZP, LKOH)

- ✅ **In-memory кэширование MOEX данных** — сессия 45
  - 23 unit-теста (11 cache + 12 cached_provider)

- ✅ **Фундаментальные данные (FundamentalData)** — сессия 46
  - 10 unit-тестов, API + фронтенд

- ✅ **LLM unit-тесты** — сессия 47
  - 20 тестов с мок OpenAI-compatible сервером

- ✅ **Свечной график (candlestick chart)** — сессия 48
  - `GET /api/ticker/{symbol}/candles` endpoint
  - chartjs-chart-financial + luxon + adapter
  - Fallback на line chart
  - Таблица фундаментальных данных с русскими подписями
  - 3 новых API теста
  - Всего Go тестов: 106, Python тестов: 290

- ✅ **Docker Compose тест** — сессия 49
  - Исправлен Dockerfile (пути к бинарнику и фронтенду)
  - `docker-compose build` + `up` + health check + API + фронтенд — всё работает
  - Коммит `35952ec` запушен в `origin/main`

- ✅ **Cache stats endpoint + кнопки быстрого выбора** — сессия 50
  - `GET /api/cache/stats` — CacheStatsResponse (tickers, candles, fundamentals, total)
  - CacheStatsProvider интерфейс + auto-detection
  - Кнопки быстрого выбора тикеров (SBER, GAZP, LKOH, GMKN, ROSN, NVTK, YDEX, TATN)
  - 2 новых Go-теста, всего 108 Go + 290 Python
  - Коммит `8dd4c7a` запушен в `origin/main`

## Следующий разумный шаг (сессия 51)

1. **Отображение cache stats в UI** — показать статистику кэша на дашборде
2. **Docker Compose healthcheck** — добавить в docker-compose.yml
3. **Улучшить свечной график** — тултипы, кроссхейр, зум
4. **Расчётные метрики** — P/E, P/B на основе доступных данных
5. **Кнопка очистки кэша** — `POST /api/cache/clear` + кнопка в UI
