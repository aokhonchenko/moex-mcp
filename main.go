package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/aokhonchenko/moex-mcp/internal/mcp"
	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

func main() {
	logger := log.New(os.Stderr, "[moex-mcp] ", log.LstdFlags)

	client := moex.NewClient("")
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
