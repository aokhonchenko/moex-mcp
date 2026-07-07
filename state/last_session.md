# Сообщение будущей сессии (сессия 86)

## Что было сделано в сессии 86

**foundation-finance: персистентность новостей (NER-модуль).**

### foundation-finance (1 коммит: `31189ae`)

1. **Персистентность хранилища новостей** — `news.Store` теперь поддерживает `NewPersistentStore(filePath)` с JSON-файлом (по аналогии с portfolio)
2. **main.go** — используется `NEWS_FILE` env (по умолчанию `data/news.json`), аналогично `PORTFOLIO_FILE`
3. **4 новых теста** — `TestPersistentStoreSaveAndLoad`, `TestPersistentStoreDeletePersisted`, `TestPersistentStoreNonexistentFile`, `TestPersistentStoreNextID`
4. **Все мутирующие методы** (`Add`, `Delete`, `SaveAnalysis`) вызывают `save()` после изменений

### Проверки

- Go тесты: все PASS (~274 foundation-finance + 66 moex-mcp + 4 новых)
- `git push origin main` — `31189ae`

### Что важно для следующей сессии

1. **Docker Compose** — `docker compose up --build` должен запустить 3 сервиса (moex-mcp, app, Ollama локально). Новости теперь сохраняются в `data/news.json` (volume `app-data`).
2. **NER-сервер** — модуль новостей уже создан (session 86 origin `dab3717`), но можно улучшить: админка для редактирования связей сущностей, bulk-импорт новостей.
3. **moex-mcp: LLM-интеграция** — подключение к Claude Desktop / Cursor (MCP stdio режим)
4. **moex-mcp: кэш-статистика на фронтенде** — показывать hits/misses из moex-mcp
5. **Фронтенд: пустые графики** — из external_messages: на свечном графике, RSI, MACD и других индикаторах пусто. Нужно проверить и исправить.
