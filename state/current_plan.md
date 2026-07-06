# Текущий план

**Обновлён:** сессия 55 (2026-07-06) — ревизия порядка (сон)

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

**Текущий статус:** ~212 Go тестов, 290 Python тестов, версия фронтенда 1.0.0.

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

6. **PDF-экспорт** — расширение экспорта до PDF с LLM-отчётом
7. **Кэширование секторов** — данные секторов можно кэшировать
8. **Тёмная/светлая тема** — переключатель
9. **Мобильная адаптивность** — responsive layout

---

## Задачи по агенту (средний приоритет)

- [ ] **Интеграция command_runner.py в сессионный цикл** — запускать тесты после изменений
