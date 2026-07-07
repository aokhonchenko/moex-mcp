package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aokhonchenko/moex-mcp/internal/httpserver"
	"github.com/aokhonchenko/moex-mcp/internal/mcp"
	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

func main() {
	mode := flag.String("mode", "stdio", "Режим работы: stdio (MCP JSON-RPC) или http (REST API)")
	addr := flag.String("addr", ":8081", "Адрес HTTP-сервера (только для -mode=http)")
	board := flag.String("board", "", "Торговая площадка MOEX (по умолчанию TQBR)")
	flag.Parse()

	logger := log.New(os.Stderr, "[moex-mcp] ", log.LstdFlags)
	client := moex.NewClient(*board)

	switch *mode {
	case "stdio":
		runStdio(client, logger)
	case "http":
		runHTTP(client, logger, *addr)
	default:
		logger.Fatalf("неизвестный режим: %s (допустимые: stdio, http)", *mode)
	}
}

// runStdio запускает MCP-сервер в режиме stdio (JSON-RPC 2.0 построчно).
func runStdio(client *moex.Client, logger *log.Logger) {
	server := mcp.NewServer(client, logger)

	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	logger.Println("MOEX MCP server started (stdio)")

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				logger.Println("EOF, shutting down")
				return
			}
			logger.Printf("read error: %v", err)
			return
		}

		var req mcp.JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			logger.Printf("invalid JSON: %v", err)
			continue
		}

		resp := server.HandleRequest(&req)
		if resp != nil {
			if err := encoder.Encode(resp); err != nil {
				logger.Printf("write error: %v", err)
				return
			}
		}
	}
}

// runHTTP запускает REST API сервер.
func runHTTP(client *moex.Client, logger *log.Logger, addr string) {
	srv := httpserver.NewServer(client, logger)

	fmt.Printf("MOEX MCP HTTP server starting on %s\n", addr)
	logger.Printf("Endpoints: /api/health, /api/ticker/{symbol}, /api/candles/{symbol}, /api/fundamentals/{symbol}, /api/search?q=")

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		logger.Fatalf("HTTP server error: %v", err)
	}
}
