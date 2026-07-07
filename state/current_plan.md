# Текущий план

**Обновлён:** сессия 78 (2026-07-07) — дивиденды + стакан заявок

---

## Статус: практическая фаза, финансовый дашборд

Приоритет определён создателем: улучшение агента (завершено) → финансовый дашборд.

---

## Завершённые мини-проекты (инфраструктура агента, сессии 1–41)

1. **Карта системы и навигация** — system_map, file_manifest, quick_context, artifact_links
2. **Система задач** — `tasks/active.md`, `tasks/archive.md`
3. **Система заметок** — `tools/notes/`, `knowledge/notes/`
4. **Оптимизация чтения** — partial_reader, prompt_builder, reader, compat, session_runner
5. **Инструменты агента (15 шт.)** — apply_patch, code_analyzer, command_runner, self_review, read_file, read_lines, write_file, replace_text, run_command, run_pytest, run_python_script и др.
6. **Тесты** — 9 модулей, ~290 Python-тестов
7. **UI-дашборд** — `tools/dashboard/`
8. **Система сна** — чеклисты, sleep_memory

## Завершённые шаги foundation-finance (сессии 42–54)

| Сессия | Шаг | Коммит |
|--------|-----|--------|
| 42 | Создан дашборд: Go backend (chi) + Web frontend (Chart.js) + Docker Compose | — |
| 43 | Unit-тесты indicators/calculator.go (26 тестов) | — |
| 44 | Замена Yahoo Finance на MOEX ISS API (25 тестов) | — |
| 45 | In-memory кэширование (23 теста) | — |
| 46 | Фундаментальные данные (10 тестов) | — |
| 47 | LLM unit-тесты (20 тестов) | — |
| 48 | Свечной график (candlestick chart) + chartjs-chart-financial | — |
| 49 | Docker Compose тест (build + up + healthcheck) | `35952ec` |
| 50 | Cache stats endpoint + кнопки быстрого выбора | `8dd4c7a` |
| 51 | Cache stats UI + clear cache + Docker healthcheck | `483045c` |
| 52 | (не было сессии 52) | — |
| 53 | Zoom/pan + кроссхейр (chartjs-plugin-zoom + hammerjs) | `637b98e` |
| 54 | Поиск по тикерам с автокомплитом (MOEX ISS /securities?q=) | `b2ee178` |
| 56 | Расчётные метрики (P/B, market_cap, 52-нед. диапазон) | `8f2009c` |
| 57 | Расчётные метрики + система алертов (6 метрик, 5 endpoints) | `b269646` |
| 58 | Портфель (in-memory store + 5 endpoints + UI) | `342120f` |
| 59 | Персистентность портфеля (JSON-файл, 22 теста) | `91ceb07` |
| 60 | Docker Compose volume (app-data) + TestHealth fix | `b925f23` |
| 61 | Секторальная аналитика (sectors endpoint + UI) | — |
| 62 | Экспорт CSV (portfolio + ticker + candles) | — |
| 63 | Push в origin + GetSectors на MOEXProvider (реальные данные) | `5136365` |
| 64 | Кэширование секторов (sectorsCache TTL 10 мин, 4 теста) | `2cb12b9` |
| 65 | Переключатель тёмной/светлой темы (CSS vars + localStorage) | `e8a32e0` |
| 66 | Мобильная адаптивность (responsive layout, 3 breakpoint) | `409d168` |
| 67 | PDF-экспорт (тикеры + портфель + LLM-отчёт) | `365fb42` |
| 70 | Stochastic Oscillator + VWAP (10 новых тестов) | `6fff20a` |
| 71 | Stochastic + VWAP графики на фронтенде + fix HTML duplicates | `f029a0d` |
| 72 | moex-mcp: HTTP REST API режим (internal/httpserver, 10 тестов) | `27b90ac` |
| 73 | moex-mcp + foundation-finance интеграция (MCPProvider, Docker Compose) | `41e2b7c` + `b996063` |
| 74 | Docker Compose тест + moex-mcp Dockerfile fix (Go 1.22 + wget) | `b194b45` |
| 75 | moex-mcp: секторы через секторальные индексы MOEX (10 индексов, 31 тест) | `95826ca` |
| 76 | moex-mcp: in-memory cache (56 тестов) + moex_sectors MCP + индексы на фронтенде | `b9aa3ac` + `e7fe993` |
| 77 | CachedProvider.GetIndex() + indices in cache stats (7 новых тестов) | `882ff5b` |
| 78 | moex-mcp: дивиденды + стакан заявок + интеграция (10+8 тестов) | `ce56694` + `a4ffcfc` |
| 79 | Фронтенд: дивиденды + стакан заявок (таблица + визуализация) | `7567c86` |
| 80 | Docker Compose healthcheck fix (HEAD→GET, удалён version) | `543ab87` |
| 81 | Фронтенд: автозагрузка + min-height canvas + изоляция ошибок рендеринга | `d03eeed` |

**Текущий статус:** ~270 Go тестов (foundation-finance) + 66 Go тестов (moex-mcp), версия фронтенда 1.0.0. Docker Compose работает.

---

## Следующие шаги для foundation-finance

### Высокий приоритет 🔴

1. ~~**Расчётные метрики**~~ ✅ — P/B, market_cap, 52-нед. диапазон (сессия 56)

2. ~~**Система алертов**~~ ✅ — уведомления при достижении пороговых значений индикаторов (сессия 57)

3. ~~**Портфель**~~ ✅ — добавление тикеров в «избранное», сводная таблица (сессия 58)

### Средний приоритет 🟡

3. ~~**Персистентность портфеля**~~ ✅ — сохранение в JSON-файл (сессия 59)
4. ~~**Docker Compose volume**~~ ✅ — named volume app-data для data/portfolio.json (сессия 60)
5. ~~**Секторальная аналитика**~~ ✅ — сравнение тикеров по секторам (сессия 61)
6. ~~**Экспорт отчётов**~~ ✅ — CSV экспорт данных (сессия 62)

### Низкий приоритет 🟢

6. ~~**PDF-экспорт**~~ ✅ — ReportPDF + PortfolioPDF + LLM-отчёт (сессия 67)
7. ~~**Кэширование секторов**~~ ✅ — sectorsCache TTL 10 мин, 4 теста (сессия 64)
8. ~~**Тёмная/светлая тема**~~ ✅ — переключатель с localStorage (сессия 65)
9. ~~**Мобильная адаптивность**~~ ✅ — responsive layout, 3 breakpoint (сессия 66)

**Все шаги плана foundation-finance выполнены! 🎉**

### Дополнительные улучшения (после завершения плана)

| Сессия | Шаг | Коммит |
|--------|-----|--------|
| 70 | Stochastic Oscillator + VWAP (10 новых тестов) | `6fff20a` |

### Следующие шаги для foundation-finance

1. ~~**Отображение Stochastic и VWAP на фронтенде**~~ ✅ — графики %K/%D и VWAP (сессия 71)
2. ~~**Интеграция moex-mcp**~~ ✅ — MCPProvider + Docker Compose (сессия 73)
3. ~~**Docker Compose тест**~~ ✅ — `docker compose up --build` работает (сессия 74)
4. ~~**Исправить маппинг секторов**~~ ✅ — секторы через состав секторальных индексов MOEX (сессия 75)
5. ~~**Индексы MOEX на фронтенде**~~ ✅ — IMOEX, RTSI через moex-mcp (сессия 76)
6. ~~**CachedProvider для GetIndex**~~ ✅ — кэширование индексов в foundation-finance (сессия 77)
7. ~~**moex-mcp: дивиденды**~~ ✅ — endpoint `/api/dividends/{symbol}` + интеграция (сессия 78)
8. ~~**moex-mcp: стакан заявок**~~ ✅ — endpoint `/api/orderbook/{symbol}` + интеграция (сессия 78)
9. ~~**Docker Compose тест**~~ ✅ — healthcheck fix: HEAD→GET, удалён version (сессия 80, `543ab87`)
10. ~~**Фронтенд: дивиденды и стакан**~~ ✅ — таблица дивидендов + визуализация стакана (сессия 79, `7567c86`)

---

## MOEX MCP Server (сессия 68)

| Сессия | Шаг | Коммит |
|--------|-----|--------|
| 68 | Создан MOEX MCP Server: 4 инструмента, 18 тестов, Dockerfile | `f5d3419` |

**Репозиторий:** `git@github.com:aokhonchenko/moex-mcp.git`

### Следующие шаги для moex-mcp

1. ~~**Интеграция с foundation-finance**~~ ✅ — MCPProvider + Docker Compose (сессия 73)
2. ~~**Исправление секторов**~~ ✅ — через состав секторальных индексов MOEX (сессия 75)
3. ~~**Расширение инструментов**~~ ✅ — индексы IMOEX/RTSI (сессия 72) + moex_sectors MCP (сессия 76)
4. ~~**Кэширование**~~ ✅ — in-memory cache с TTL для всех запросов (сессия 76)
5. ~~**MCP-инструмент для секторов**~~ ✅ — `moex_sectors` в JSON-RPC (сессия 76)
6. ~~**Дивиденды**~~ ✅ — endpoint `/api/dividends/{symbol}` + MCP `moex_dividends` (сессия 78)
7. ~~**Стакан заявок**~~ ✅ — endpoint `/api/orderbook/{symbol}` + MCP `moex_orderbook` (сессия 78)
8. **LLM-интеграция** — подключение к Claude Desktop / Cursor

---

## Задачи по агенту (средний приоритет)

- [ ] **Интеграция command_runner.py в сессионный цикл** — запускать тесты после изменений
