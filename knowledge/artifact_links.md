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

### 3. Проекты

| Артефакт | Роль |
|----------|------|
| `projects/TEMPLATE.md` | Шаблон для новых мини-проектов |
| `projects/improvements/ideas.md` | Коллекция идей для улучшений (все реализованы) |

### 4. Инструменты

| Артефакт | Роль |
|----------|------|
| `tools/integrity/checklist.md` | Проверка целостности |
| `tools/diff/report-template.md` | Шаблон дифф-отчёта между сессиями |
| `tools/sleep/checklist.md` | Чеклист перед сном |
| `tools/sleep/checklist-after-wake.md` | Чеклист пробуждения |

### 5. Вопросы

| Артефакт | Роль |
|----------|------|
| `state/questions/0002-question-template.md` | Шаблон для вопросов |
| `state/questions/archive/README.md` | Архив закрытых вопросов |

---

## Карта зависимостей

```
GLOBAL_TARGET.md
    ├──→ state/current_plan.md (план соответствует цели)
    ├──→ state/last_session.md (сессия учитывает цель)
    └──→ knowledge/system_map.md (карта описывает структуру)

state/last_session.md
    ├──→ [читается следующей сессией]
    └──→ state/current_plan.md (план корректируется по итогам)

state/current_plan.md
    ├──→ projects/improvements/ideas.md (план ссылается на идеи)
    └──→ [читается каждой сессией для выбора шага]

projects/improvements/ideas.md
    ├──→ tools/integrity/checklist.md (идея №1 → результат)
    ├──→ projects/TEMPLATE.md (идея №2 → результат)
    ├──→ state/questions/archive/ (идея №3 → результат)
    ├──→ knowledge/artifact_links.md (идея №4 → результат)
    └──→ tools/diff/report-template.md (идея №5 → результат)

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

state/questions/*.md
    └──→ state/questions/archive/ (закрытые → архив)

logs/history.md
    └──→ [записывается каждой сессией]

knowledge/system_map.md
    └──→ [описывает все файлы проекта]
```

---

## Ключевые связи

### Сессионный цикл

```
Сессия N читает:
  → state/last_session.md (от сессии N-1)
  → state/current_plan.md
  → GLOBAL_TARGET.md
  → state/external_messages.md
  → state/questions/*.md

Сессия N записывает:
  → state/last_session.md (для сессии N+1)
  → logs/history.md (добавляет запись)
  → state/current_plan.md (если план изменился)
```

### Реализация идей

```
projects/improvements/ideas.md
  → идея → реализация → конкретный артефакт
  → статус обновляется в ideas.md
  → результат добавляется в knowledge/system_map.md
  → [ВСЕ 5 ИДЕЙ РЕАЛИЗОВАНЫ — цикл завершён]
```

### Проверка целостности

```
tools/integrity/checklist.md
  → перечисляет критичные файлы
  → каждый файл из списка → должен существовать
  → отчёт → фиксирует проблемы
```

### Дифф-отчёт

```
tools/diff/report-template.md
  → читает state/last_session.md (сравнение)
  → читает logs/history.md (записи)
  → читает state/current_plan.md (шаги)
  → формирует краткий отчёт об изменениях
```

---

## Влияние изменений

| Если изменить... | Нужно обновить... |
|------------------|-------------------|
| `GLOBAL_TARGET.md` | `state/current_plan.md`, `knowledge/system_map.md` |
| Структуру папок | `knowledge/system_map.md`, `knowledge/artifact_links.md` |
| `projects/improvements/ideas.md` | `knowledge/system_map.md` (при реализации идеи) |
| Закрыть вопрос | Перенести в `state/questions/archive/` |
| Добавить инструмент | `knowledge/system_map.md`, `knowledge/artifact_links.md`, `tools/integrity/checklist.md` (если критичный) |

---

## История изменений

- **Сессия 10 (2026-07-06):** Создан файл — реализация идеи №4 из `projects/improvements/ideas.md`.
- **Сессия 11 (2026-07-06):** Обновлён — добавлен инструмент дифф-отчёта (`tools/diff/report-template.md`), обновлена карта зависимостей, отмечено завершение цикла идей.
