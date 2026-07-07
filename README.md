# MOEX MCP Server

MCP-сервер для доступа к данным Московской биржи (MOEX ISS API) через протокол [Model Context Protocol](https://modelcontextprotocol.io/).

Поддерживает два режима работы:
- **stdio** — JSON-RPC 2.0 (для Claude Desktop, Cursor и других MCP-хостов)
- **http** — REST API (для интеграции с другими сервисами, например foundation-finance)

## Инструменты (MCP)

| Инструмент | Описание |
|-----------|----------|
| `moex_ticker` | Текущая котировка тикера (цена, изменение, объём) |
| `moex_candles` | Исторические свечи OHLCV (1m, 3m, 6m, 1y) |
| `moex_fundamentals` | Фундаментальные данные эмитента (ISIN, номинал, объём выпуска) |
| `moex_search` | Поиск бумаг по тикеру, ISIN или названию |

## REST API endpoints (http-режим)

| Endpoint | Описание |
|----------|----------|
| `GET /api/health` | Health check |
| `GET /api/ticker/{symbol}` | Котировка тикера |
| `GET /api/candles/{symbol}?period=3m` | Свечи (1m, 3m, 6m, 1y) |
| `GET /api/fundamentals/{symbol}` | Фундаментальные данные |
| `GET /api/search?q=SBER` | Поиск бумаг |

## Запуск

### Локально (stdio)

```bash
go build -o moex-mcp .
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./moex-mcp
```

### Локально (HTTP)

```bash
go build -o moex-mcp .
./moex-mcp -mode=http -addr=:8081
```

### Docker (stdio)

```bash
docker build -t moex-mcp .
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | docker run -i moex-mcp
```

### Docker (HTTP)

```bash
docker build -t moex-mcp .
docker run -p 8081:8081 moex-mcp -mode=http -addr=:8081
```

### Claude Desktop / Cursor

Добавьте в конфиг MCP:

```json
{
  "mcpServers": {
    "moex": {
      "command": "path/to/moex-mcp",
      "args": []
    }
  }
}
```

## Тесты

```bash
go test ./... -v
```

## Структура

```
moex-mcp/
├── main.go                          # Точка входа (stdio + http)
├── internal/
│   ├── moex/
│   │   ├── client.go                # HTTP-клиент MOEX ISS API
│   │   └── client_test.go           # 8 тестов с мок-сервером
│   ├── mcp/
│   │   ├── server.go                # MCP JSON-RPC обработчик
│   │   └── server_test.go           # 10 тестов MCP-протокола
│   └── httpserver/
│       ├── server.go                # REST API обёртка
│       └── server_test.go           # 10 тестов HTTP API
├── Dockerfile
├── go.mod
└── README.md
```

## Источники данных

- [MOEX ISS API Reference](https://iss.moex.com/iss/reference/)
- Доска: TQBR (основной режим торгов акциями)
- User-Agent: `moex-mcp/0.1`

## Лицензия

MIT
