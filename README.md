# MOEX MCP Server

MCP-сервер для доступа к данным Московской биржи (MOEX ISS API) через протокол [Model Context Protocol](https://modelcontextprotocol.io/).

## Инструменты

| Инструмент | Описание |
|-----------|----------|
| `moex_ticker` | Текущая котировка тикера (цена, изменение, объём) |
| `moex_candles` | Исторические свечи OHLCV (1m, 3m, 6m, 1y) |
| `moex_fundamentals` | Фундаментальные данные эмитента (ISIN, номинал, объём выпуска) |
| `moex_search` | Поиск бумаг по тикеру, ISIN или названию |

## Запуск

### Локально

```bash
go build -o moex-mcp .
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./moex-mcp
```

### Docker

```bash
docker build -t moex-mcp .
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | docker run -i moex-mcp
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
├── main.go                    # Точка входа (stdio JSON-RPC)
├── internal/
│   ├── moex/
│   │   ├── client.go          # HTTP-клиент MOEX ISS API
│   │   └── client_test.go     # 8 тестов с мок-сервером
│   └── mcp/
│       ├── server.go          # MCP JSON-RPC обработчик
│       └── server_test.go     # 10 тестов MCP-протокола
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
