package mcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

// JSONRPCRequest - входящий JSON-RPC 2.0 запрос.
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse - исходящий JSON-RPC 2.0 ответ.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError - ошибка JSON-RPC.
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Server - MCP-сервер для MOEX ISS API.
type Server struct {
	client *moex.Client
	logger *log.Logger
}

// NewServer создаёт новый MCP-сервер.
func NewServer(client *moex.Client, logger *log.Logger) *Server {
	return &Server{client: client, logger: logger}
}

// HandleRequest обрабатывает JSON-RPC запрос.
func (s *Server) HandleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		return nil
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "ping":
		return s.successResponse(req.ID, map[string]interface{}{})
	default:
		return s.errorResponse(req.ID, -32601, fmt.Sprintf("unknown method: %s", req.Method))
	}
}

func (s *Server) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	return s.successResponse(req.ID, map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "moex-mcp",
			"version": "0.1.0",
		},
	})
}

func (s *Server) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := []map[string]interface{}{
		{
			"name":        "moex_ticker",
			"description": "Get current quote for a MOEX ticker. Returns price, change, volume.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "MOEX ticker symbol (e.g. SBER, GAZP, LKOH, YNDX)",
					},
				},
				"required": []string{"symbol"},
			},
		},
		{
			"name":        "moex_candles",
			"description": "Get historical OHLCV candles for a MOEX ticker.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "MOEX ticker symbol",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "Period: 1m, 3m, 6m, 1y (default 3m)",
						"enum":        []string{"1m", "3m", "6m", "1y"},
					},
				},
				"required": []string{"symbol"},
			},
		},
		{
			"name":        "moex_fundamentals",
			"description": "Get fundamental data for a MOEX issuer: ISIN, issue size, face value, sec type.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "MOEX ticker symbol",
					},
				},
				"required": []string{"symbol"},
			},
		},
		{
			"name":        "moex_search",
			"description": "Search MOEX securities by ticker, ISIN or name fragment.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "moex_index",
			"description": "Get current value for a MOEX index (IMOEX, RTSI, MOEXFN, etc.).",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "Index symbol (e.g. IMOEX, RTSI, MOEXFN, MOEXOG)",
					},
				},
				"required": []string{"symbol"},
			},
		},
		{
			"name":        "moex_sectors",
			"description": "Get sector analytics — group MOEX stocks by sector with average change.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "moex_dividends",
			"description": "Get dividend history for a MOEX ticker.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "MOEX ticker symbol (e.g. SBER, GAZP, LKOH)",
					},
				},
				"required": []string{"symbol"},
			},
		},
		{
			"name":        "moex_orderbook",
			"description": "Get order book (Level 2) for a MOEX ticker. Returns price levels with buy/sell quantities.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"symbol": map[string]interface{}{
						"type":        "string",
						"description": "MOEX ticker symbol (e.g. SBER, GAZP)",
					},
				},
				"required": []string{"symbol"},
			},
		},
	}

	return s.successResponse(req.ID, map[string]interface{}{"tools": tools})
}

func (s *Server) handleToolsCall(req *JSONRPCRequest) *JSONRPCResponse {
	var callParams struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &callParams); err != nil {
		return s.errorResponse(req.ID, -32602, "invalid params: "+err.Error())
	}

	switch callParams.Name {
	case "moex_ticker":
		return s.callTicker(req.ID, callParams.Arguments)
	case "moex_candles":
		return s.callCandles(req.ID, callParams.Arguments)
	case "moex_fundamentals":
		return s.callFundamentals(req.ID, callParams.Arguments)
	case "moex_search":
		return s.callSearch(req.ID, callParams.Arguments)
	case "moex_index":
		return s.callIndex(req.ID, callParams.Arguments)
	case "moex_sectors":
		return s.callSectors(req.ID)
	case "moex_dividends":
		return s.callDividends(req.ID, callParams.Arguments)
	case "moex_orderbook":
		return s.callOrderBook(req.ID, callParams.Arguments)
	default:
		return s.errorResponse(req.ID, -32602, fmt.Sprintf("unknown tool: %s", callParams.Name))
	}
}

func (s *Server) callTicker(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.GetTicker(p.Symbol)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callCandles(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
		Period string `json:"period"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}
	if p.Period == "" {
		p.Period = "3m"
	}

	data, err := s.client.GetCandles(p.Symbol, p.Period)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callFundamentals(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.GetFundamentals(p.Symbol)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callSearch(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.SearchSecurities(p.Query)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callIndex(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.GetIndex(p.Symbol)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callSectors(id json.RawMessage) *JSONRPCResponse {
	data, err := s.client.GetSectors()
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callDividends(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.GetDividends(p.Symbol)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) callOrderBook(id json.RawMessage, args json.RawMessage) *JSONRPCResponse {
	var p struct {
		Symbol string `json:"symbol"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return s.errorResponse(id, -32602, "invalid argument: "+err.Error())
	}

	data, err := s.client.GetOrderBook(p.Symbol)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolResult(id, data)
}

func (s *Server) successResponse(id json.RawMessage, result interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func (s *Server) errorResponse(id json.RawMessage, code int, msg string) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: msg,
		},
	}
}

func (s *Server) toolResult(id json.RawMessage, data interface{}) *JSONRPCResponse {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return s.errorResponse(id, -32603, "serialization error: "+err.Error())
	}

	return s.successResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(jsonBytes),
			},
		},
	})
}

func (s *Server) toolError(id json.RawMessage, msg string) *JSONRPCResponse {
	return s.successResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Error: %s", msg),
			},
		},
		"isError": true,
	})
}
