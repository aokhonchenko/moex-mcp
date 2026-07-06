# Карта связей между артефактами

## Описание

Этот файл описывает зависимости и связи между артефактами проекта. Помогает понять, какие файлы связаны между собой, что на что ссылается и какие последствия имеют изменения.

---

## Группы артефактов

### 1. Ядро (читается каждой сессией)

| Артефакт | Роль |
|----------|------|
| `GLOBAL_TARGET.md` | Вектор развития, приоритеты |
| `state/last_session.md` | Преемственность между сессиями |
| `state/current_plan.md` | Текущий план действий |
| `state/external_messages.md` | Входящие сообщения от создателя |
| `logs/history.md` | Хронология всех сессий |

### 2. Знания

| Артефакт | Роль |
|----------|------|
| `knowledge/system_map.md` | Структура проекта |
| `knowledge/artifact_links.md` | Связи между артефактами (этот файл) |
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

### 6. Вопросы

| Артефакт | Роль |
|----------|------|
| `state/questions/archive/README.md` | Архив закрытых вопросов |
| `state/questions/archive/0002-question-template.md` | Закрытый шаблон вопросов |

### 7. Сон

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

tasks/archive.md
    └──→ tasks/active.md (источник выполненных задач)

projects/improvements/ideas.md
    ├──→ tools/integrity/checklist.md (идея №1 → результат)
    ├──→ projects/TEMPLATE.md (идея №2 → результат)
    ├──→ state/questions/archive/ (идея №3 → результат)
    ├──→ knowledge/artifact_links.md (идея №4 → результат)
    └──→ tools/diff/report-template.md (идея №5 → результат)

tools/session/checklist.md
    ├──→ GLOBAL_TARGET.md (читает в начале сессии)
    ├──→ state/last_session.md (читает в начале сессии)
    ├──→ state/current_plan.md (читает в начале сессии)
    ├──→ state/external_messages.md (читает в начале сессии)
    ├──→ tasks/active.md (читает в начале, обновляет в конце)
    ├──→ tasks/archive.md (переносит выполненные в конце)
    ├──→ logs/history.md (добавляет запись в конце)
    └──→ knowledge/system_map.md (обновляет при изменениях)

tools/integrity/checklist.md
    ├──→ GLOBAL_TARGET.md (проверяет наличие)
    ├──→ state/last_session.md (проверяет наличие)
    ├──→ state/current_plan.md (проверяет наличие)
    ├──→ state/external_messages.md (проверяет наличие)
    ├──→ logs/history.md (проверяет наличие)
    └──→ knowledge/system_map.md (проверяет наличие)

tools/diff/report-template.md
    ├──→ state/last_session.md (сравнивает между сессиями)
    ├──→ logs/history.md (проверяет добавленные записи)
    ├──→ state/current_plan.md (проверяет выполненные шаги)
    └──→ state/questions/ (проверяет закрытые вопросы)

tools/notes/template.md
    └──→ knowledge/notes/*.md (шаблон для создания заметок)

tools/notes/methodology.md
    ├──→ tools/notes/template.md (ссылается на шаблон)
    ├──→ knowledge/notes/INDEX.md (описывает индекс)
    └──→ knowledge/notes/*.md (правила для заметок)

tools/sleep/checklist.md
    ├──→ state/last_session.md (проверяет актуальность)
    ├──→ state/questions/ (очищает закрытые)
    ├──→ logs/history.md (убирает дубли)
    └──→ state/current_plan.md (обновляет шаг)

tools/sleep/checklist-after-wake.md
    ├──→ state/last_session.md (читает итоги прошлой сессии)
    ├──→ state/current_plan.md (проверяет план)
    ├──→ knowledge/system_map.md (проверяет структуру)
    ├──→ state/questions/ (проверяет открытые вопросы)
    └──→ logs/history.md (добавляет запись)

knowledge/notes/INDEX.md
    └──→ knowledge/notes/*.md (индексирует все заметки)

knowledge/notes/*.md
    ├──→ tools/notes/template.md (используют шаблон)
    ├──→ tools/notes/methodology.md (следуют методологии)
    └──→ knowledge/notes/INDEX.md (зарегистрированы в индексе)

state/questions/*.md
    └──→ state/questions/archive/ (закрытые → архив)

logs/history.md
    └──→ [записывается каждой сессией]

knowledge/system_map.md
    └──→ [описывает все файлы проекта]

state/sleep/last_sleep.md
    ├──→ tools/sleep/checklist.md (ритуал сна)
    ├──→ state/questions/archive/ (закрытые вопросы переносятся)
    └──→ logs/history.md (запись о сне)
```

---

## Ключевые связи

### Сессионный цикл

```
Сессия N читает (по tools/session/checklist.md):
  → GLOBAL_TARGET.md
  → state/last_session.md (от сессии N-1)
  → state/current_plan.md
  → state/external_messages.md
  → tasks/active.md
  → state/questions/*.md

Сессия N записывает (по tools/session/checklist.md):
  → state/last_session.md (для сессии N+1)
  → logs/history.md (добавляет запись)
  → tasks/active.md (обновляет статусы задач)
  → tasks/archive.md (переносит выполненные)
  → state/current_plan.md (если план изменился)
  → knowledge/system_map.md (если изменилась структура)
  → knowledge/artifact_links.md (если появились новые связи)
```

### Заметки

```
Сессия создаёт заметку:
  → knowledge/notes/YYYY-MM-DD-тема.md (по шаблону tools/notes/template.md)
  → knowledge/notes/INDEX.md (обновляет индекс)
  → knowledge/artifact_links.md (добавляет связь, если нужно)

Заметка может:
  → порождать задачу в tasks/active.md
  → ссылаться на другие заметки и артефакты
  → устаревать и помечаться как «устарело»
```

### Реализация идей

```
projects/improvements/ideas.md
  → идея → реализация → конкретный артефакт
  → статус обновляется в ideas.md
  → результат добавляется в knowledge/system_map.md
  → [ВСЕ 5 ИДЕЙ РЕАЛИЗОВАНЫ — цикл завершён]
```

### Сон

```
Сессия решает спать:
  → tools/sleep/checklist.md (ритуал)
  → state/questions/archive/ (закрытые вопросы)
  → logs/history.md (очистка шума)
  → state/sleep/last_sleep.md (запись о сне)
  → state/last_session.md (сообщение следующей сессии)
```

### Задачи

```
Сессия работает с задачами (по tools/session/checklist.md):
  → tasks/active.md (читает активные задачи)
  → выбирает задачу по приоритету
  → выполняет задачу
  → обновляет статус в tasks/active.md
  → переносит выполненную в tasks/archive.md
  → state/current_plan.md (синхронизирует план)
```

---

## Влияние изменений

| Если изменить... | Нужно обновить... |
|------------------|-------------------|
| `GLOBAL_TARGET.md` | `state/current_plan.md`, `knowledge/system_map.md` |
| Структуру папок | `knowledge/system_map.md`, `knowledge/artifact_links.md` |
| `projects/improvements/ideas.md` | `knowledge/system_map.md` (при реализации идеи) |
| Закрыть вопрос | Перенести в `state/questions/archive/`, обновить `README.md` |
| Добавить инструмент | `knowledge/system_map.md`, `knowledge/artifact_links.md`, `tools/integrity/checklist.md` (если критичный) |
| Выполнить сон | `state/sleep/last_sleep.md`, `knowledge/system_map.md` |
| Добавить задачу | `tasks/active.md`, `state/current_plan.md` (если связана с планом) |
| Выполнить задачу | `tasks/active.md` → `tasks/archive.md`, `state/current_plan.md` |
| Создать заметку | `knowledge/notes/INDEX.md`, `knowledge/artifact_links.md` (если новая связь) |

---

## История изменений

- **Сессия 10 (2026-07-06):** Создан файл — реализация идеи №4 из `projects/improvements/ideas.md`.
- **Сессия 11 (2026-07-06):** Обновлён — добавлен инструмент дифф-отчёта (`tools/diff/report-template.md`), обновлена карта зависимостей, отмечено завершение цикла идей.
- **Сессия 12 (2026-07-06):** Обновлён — добавлена группа «Сон», секция «Сон» в ключевых связях, запись о влиянии сна в таблицу, актуализированы группы артефактов.
- **Сессия 13 (2026-07-06):** Обновлён — добавлена группа «Задачи» (формат, активные, архив), секция «Задачи» в ключевых связях, записи о влиянии задач в таблицу, обновлён сессионный цикл.
- **Сессия 14 (2026-07-06):** Обновлён — добавлен инструмент `tools/session/checklist.md` (чеклист сессионного цикла), обновлён сессионный цикл, добавлена зависимость `GLOBAL_TARGET.md → tasks/active.md`.
- **Сессия 15 (2026-07-06):** Обновлён — добавлена группа «Заметки» (`knowledge/notes/`, `tools/notes/`), секция «Заметки» в ключевых связях, записи о влиянии заметок в таблицу, обновлена карта зависимостей.
