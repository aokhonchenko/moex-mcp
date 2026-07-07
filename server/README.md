# Сервер управления сессиями ai-lives

Веб-интерфейс для управления автономными сессиями агента.

## Запуск

```bash
# Windows
server.bat

# Или напрямую
uv run python server/server.py --port 11000
```

Откройте http://127.0.0.1:11000 в браузере.

## Возможности

- **Запуск сессии** — кнопка «Запустить сессию» вызывает `uv run python scripts/session_transaction.py`
- **Автосессия** — toggle-кнопка включает цикл: запуск → пауза 30 сек → повтор
- **Real-time обновления** — SSE-поток обновляет статус и `last_session.md` без перезагрузки
- **Лог событий** — история запусков и ошибок

## API

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/` | Веб-дашборд |
| GET | `/api/last-session` | Содержимое `state/last_session.md` |
| GET | `/api/status` | Статус (idle/running/auto) |
| GET | `/api/events` | SSE-поток событий |
| POST | `/api/session/start` | Запуск одной сессии |
| POST | `/api/auto/toggle` | Вкл/выкл автосессии |

## Архитектура

- Python stdlib (ThreadingHTTPServer + threading) — без внешних зависимостей
- SSE (Server-Sent Events) для real-time обновлений
- Потокобезопасное глобальное состояние
