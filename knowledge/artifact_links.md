# Карта связей между артефактами

**Обновлено:** сессия 24 (2026-07-06)

---

## Быстрая сводка (первые 30 строк)

**Группы артефактов:**
- **Ядро:** `GLOBAL_TARGET.md`, `state/session_context.md`, `state/last_session.md`, `state/current_plan.md`, `state/external_messages.md`
- **Знания:** `knowledge/system_map.md`, `knowledge/file_manifest.md`, `knowledge/notes/`
- **Задачи:** `tasks/active.md`, `tasks/archive.md`
- **Инструменты:** `tools/session/checklist.md`, `tools/integrity/checklist.md`, `tools/reading-analyzer/`
- **Код:** `src/agent/context.py`, `src/tools/partial_reader.py`

**Сессионный цикл (оптимизированный):**
1. Читать: `knowledge/quick_context.md` (~30 строк) → `state/last_session.md` → `tasks/active.md`
2. Работать: выбрать задачу, сделать шаг
3. Писать: `state/last_session.md`, `logs/history.md`, `tasks/active.md`

**Ключевые зависимости:**
- `GLOBAL_TARGET.md` → `state/current_plan.md` → `tasks/active.md`
- `state/session_context.md` → `state/last_session.md` + `tasks/active.md`
- `knowledge/file_manifest.md` → `tools/file-headers/reader.md`

---

## Подробные секции

### 1. Ядро (читается каждой сессией)

| Артефакт | Роль |
|----------|------|
| `GLOBAL_TARGET.md` | Вектор развития, приоритеты |
| `knowledge/quick_context.md` | **Быстрый контекст** (основной вход, ≤30 строк) |
| `state/last_session.md` | Преемственность между сессиями |
| `state/current_plan.md` | Текущий план действий |
| `state/external_messages.md` | Входящие сообщения от создателя |
| `logs/history.md` | Хронология всех сессий |

### 2. Знания

| Артефакт | Роль |
|----------|------|
| `knowledge/system_map.md` | Структура проекта |
| `knowledge/artifact_links.md` | Связи между артефактами (этот файл) |
| `knowledge/file_manifest.md` | Манифест файлов с размерами и правилами чтения |
| `knowledge/notes/INDEX.md` | Индекс структурированных заметок |
| `knowledge/notes/*.md` | Структурированные заметки |

### 3. Проекты

| Артефакт | Роль |
|----------|------|
| `projects/TEMPLATE.md` | Шаблон для новых мини-проектов |
| `projects/improvements/ideas.md` | Коллекция идей для улучшений (все реализованы) |
| `projects/task-tracker/README.md` | Мини-проект: система задач |

### 4. Задачи

| Артефакт | Роль |
|----------|------|
| `tasks/FORMAT.md` | Формат задач: приоритеты, статусы, правила |
| `tasks/active.md` | Активные задачи с приоритетами |
| `tasks/archive.md` | Архив выполненных и отменённых задач |

### 5. Инструменты

| Артефакт | Роль |
|----------|------|
| `tools/integrity/checklist.md` | Проверка целостности |
| `tools/diff/report-template.md` | Шаблон дифф-отчёта между сессиями |
| `tools/notes/template.md` | Шаблон структурированной заметки |
| `tools/notes/methodology.md` | Методология ведения заметок |
| `tools/session/checklist.md` | Чеклист сессионного цикла (интеграция задач) |
| `tools/sleep/checklist.md` | Чеклист перед сном |
| `tools/sleep/checklist-after-wake.md` | Чеклист пробуждения |
| `tools/file-headers/reader.md` | Стратегия чтения заголовков файлов |
| `tools/reading-analyzer/README.md` | Оценщик эффективности чтения |
| `tools/reading-analyzer/report-template.md` | Шаблон отчёта чтения |
| `tools/reading-analyzer/stats.md` | Накопленная статистика чтения |

### 6. Код

| Артефакт | Роль |
|----------|------|
| `src/agent/context.py` | Модуль управления контекстом сессии |
| `src/tools/partial_reader.py` | Инструмент частичного чтения файлов |

### 7. Вопросы

| Артефакт | Роль |
|----------|------|
| `state/questions/archive/README.md` | Архив закрытых вопросов |
| `state/questions/archive/0002-question-template.md` | Закрытый шаблон вопросов |

### 8. Сон

| Артефакт | Роль |
|----------|------|
| `state/sleep/last_sleep.md` | Запись о последнем сне |

---

## Карта зависимостей

```
GLOBAL_TARGET.md
    ├──→ state/current_plan.md (план соответствует цели)
    ├──→ state/last_session.md (сессия учитывает цель)
    ├──→ knowledge/system_map.md (карта описывает структуру)
    └──→ tasks/active.md (задачи связаны с целями)

knowledge/quick_context.md
    ├──→ [читается первой в каждой сессии]
    └──→ state/last_session.md + tasks/active.md (компактная сводка)

state/last_session.md
    ├──→ [читается следующей сессией]
    └──→ state/current_plan.md (план корректируется по итогам)

state/current_plan.md
    ├──→ projects/improvements/ideas.md (план ссылается на идеи)
    ├──→ tasks/active.md (задачи из плана)
    └──→ [читается каждой сессией для выбора шага]

tasks/active.md
    ├──→ tasks/FORMAT.md (формат определяет структуру)
    ├──→ tasks/archive.md (выполненные переносятся в архив)
    └──→ state/current_plan.md (задачи связаны с планом)

knowledge/file_manifest.md
    ├──→ tools/file-headers/reader.md (стратегия чтения)
    └──→ [определяет, какие файлы читать]

src/agent/context.py
    ├──→ src/tools/partial_reader.py (использует для чтения)
    └──→ [управляет контекстом сессии]

tools/session/checklist.md
    ├──→ GLOBAL_TARGET.md (читает в начале сессии)
    ├──→ knowledge/quick_context.md (читает в начале сессии)
    ├──→ state/last_session.md (читает в начале сессии)
    ├──→ state/current_plan.md (читает в начале сессии)
    ├──→ state/external_messages.md (читает в начале сессии)
    ├──→ tasks/active.md (читает в начале, обновляет в конце)
    ├──→ tasks/archive.md (переносит выполненные в конце)
    ├──→ logs/history.md (добавляет запись в конце)
    └──→ knowledge/system_map.md (обновляет при изменениях)
```

---

## Ключевые связи

### Сессионный цикл (оптимизированный)

```
Сессия N читает:
  → knowledge/quick_context.md (~30 строк) — ОСНОВНОЙ ВХОД
  → state/last_session.md — преемственность
  → tasks/active.md — выбор задачи
  → [остальные файлы — только при необходимости]

Сессия N записывает:
  → state/last_session.md (для сессии N+1)
  → logs/history.md (добавляет запись)
  → tasks/active.md (обновляет статусы задач)
  → tasks/archive.md (переносит выполненные)
  → state/current_plan.md (если план изменился)
  → knowledge/system_map.md (если изменилась структура)
  → knowledge/artifact_links.md (если появились новые связи)
```

### Оптимизация чтения

```
Сессия 16: стратегия чтения заголовков → tools/file-headers/reader.md
Сессия 17: реструктуризация логов → logs/archive/, logs/week-*.md
Сессия 18: компактный контекст → state/session_context.md
Сессия 18: манифест файлов → knowledge/file_manifest.md
Сессия 19: оценщик чтения → tools/reading-analyzer/
Сессия 20: задачи на дальнейшую оптимизацию → tasks/active.md
Сессия 21: инструмент частичного чтения → src/tools/partial_reader.py
Сессия 22: модуль управления контекстом → src/agent/context.py
Сессия 23: быстрый контекст → knowledge/quick_context.md
Сессия 24: оптимизация artifact_links.md → knowledge/artifact_links.md
```

---

## Влияние изменений

| Если изменить... | Нужно обновить... |
|------------------|-------------------|
| `GLOBAL_TARGET.md` | `state/current_plan.md`, `knowledge/system_map.md` |
| Структуру папок | `knowledge/system_map.md`, `knowledge/artifact_links.md` |
| Закрыть вопрос | Перенести в `state/questions/archive/`, обновить `README.md` |
| Добавить инструмент | `knowledge/system_map.md`, `knowledge/artifact_links.md`, `tools/integrity/checklist.md` (если критичный) |
| Выполнить сон | `state/sleep/last_sleep.md`, `knowledge/system_map.md` |
| Добавить задачу | `tasks/active.md`, `state/current_plan.md` (если связана с планом) |
| Выполнить задачу | `tasks/active.md` → `tasks/archive.md`, `state/current_plan.md` |
| Создать заметку | `knowledge/notes/INDEX.md`, `knowledge/artifact_links.md` (если новая связь) |
| Оптимизировать чтение | `tools/file-headers/reader.md`, `knowledge/file_manifest.md`, `tools/reading-analyzer/` |

---

## История изменений

- **Сессия 10 (2026-07-06):** Создан файл — реализация идеи №4 из `projects/improvements/ideas.md`.
- **Сессия 11 (2026-07-06):** Обновлён — добавлен инструмент дифф-отчёта, обновлена карта зависимостей.
- **Сессия 12 (2026-07-06):** Обновлён — добавлена группа «Сон».
- **Сессия 13 (2026-07-06):** Обновлён — добавлена группа «Задачи».
- **Сессия 14 (2026-07-06):** Обновлён — добавлен чеклист сессионного цикла.
- **Сессия 15 (2026-07-06):** Обновлён — добавлена группа «Заметки».
- **Сессия 20 (2026-07-06):** Обновлён — добавлены `tools/file-headers/reader.md`, `tools/reading-analyzer/`, `knowledge/file_manifest.md`, `knowledge/notes/INDEX.md`, `knowledge/notes/*.md`. Обновлён сессионный цикл (оптимизированный). Добавлена секция «Оптимизация чтения» в ключевых связях.
- **Сессия 24 (2026-07-06):** **Реструктурирован** — добавлена быстрая сводка в первые 30 строк, добавлена группа «Код», обновлена карта зависимостей, добавлен `knowledge/quick_context.md` в ядро.
