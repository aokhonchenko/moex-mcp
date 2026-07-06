# Инструмент частичных правок (apply_patch)

**Создан:** сессия 35 (2026-07-06)
**Файл:** `src/tools/apply_patch.py`
**Тесты:** `tests/test_apply_patch.py` (~30 тестов)

## Назначение

Инструмент для внесения точечных изменений в файлы без перезаписи всего содержимого.
Решает проблему неоптимальной работы с моделью: вместо записи всего файла целиком
можно применить только нужное изменение.

## Операции

| Операция | Описание |
|----------|----------|
| `--replace <старый> <новый>` | Замена первого вхождения текста |
| `--replace-all <старый> <новый>` | Замена всех вхождений текста |
| `--regex <шаблон> <замена>` | Замена по регулярному выражению |
| `--insert-after <текст> <строка>` | Вставка текста после строки |
| `--insert-before <текст> <строка>` | Вставка текста перед строкой |
| `--delete <текст>` | Удаление строк, содержащих текст |
| `--delete-range <start> <end>` | Удаление диапазона строк (1-based) |
| `--section <имя> <содержимое>` | Замена секции markdown (##) |

## Общие опции

- `--dry-run` — показать изменения без записи
- `--help` — справка

## Примеры использования

```bash
# Замена текста
python src/tools/apply_patch.py file.py --replace 'old_func' 'new_func'

# Замена по regex
python src/tools/apply_patch.py file.py --regex 'def (\w+)' 'def renamed_\1'

# Вставка после строки
python src/tools/apply_patch.py file.py --insert-after 'import os' 'import sys'

# Удаление строк с TODO
python src/tools/apply_patch.py file.py --delete 'TODO'

# Замена секции markdown
python src/tools/apply_patch.py file.md --section 'Контекст' 'Новый текст'

# Предпросмотр без записи
python src/tools/apply_patch.py file.py --replace 'foo' 'bar' --dry-run
```

## Программное использование

```python
from apply_patch import replace_text, replace_regex, insert_after_line

# Замена текста
result = replace_text('file.py', 'old', 'new')
if result.applied:
    print(f"Изменено: {result.changes}")

# Замена по regex
result = replace_regex('file.py', r'def (\w+)', r'def renamed_\1')

# Вставка после строки
result = insert_after_line('file.py', 'import os', 'import sys')

# Сухой прогон
result = replace_text('file.py', 'old', 'new', dry_run=True)
print(result.preview)  # diff-подобный вывод
```

## API

### `PatchResult`
- `path: str` — путь к файлу
- `applied: bool` — успешно ли применён патч
- `operation: str` — тип операции
- `changes: int` — количество изменений
- `preview: str` — предпросмотр изменений
- `error: Optional[str]` — описание ошибки

### Функции

Все функции возвращают `PatchResult` и поддерживают `dry_run`:

- `replace_text(filepath, old_text, new_text, count=1, dry_run=False)`
- `replace_regex(filepath, pattern, replacement, count=0, dry_run=False)`
- `insert_after_line(filepath, target, text, dry_run=False)`
- `insert_before_line(filepath, target, text, dry_run=False)`
- `delete_lines(filepath, target, dry_run=False)`
- `delete_line_range(filepath, start, end, dry_run=False)`
- `replace_section(filepath, section_name, new_content, dry_run=False)`

## Связи

- Создан в ответ на просьбу создателя: "ты все еще пишешь файлы целиком. собери себе инструмент для частичных правок типа apply"
- Дополняет `reader.py` (точечное чтение) — теперь есть и точечное чтение, и точечная запись
- Должен быть интегрирован в сессионный цикл для замены `write_file` на частичные правки
