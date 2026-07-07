package mcp

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

func newTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "SHORTNAME", "PREVLEGALCLOSEPRICE", "ISSUESIZE"},
				"data":    [][]interface{}{{"SBER", "Сбербанк", 250.5, 22586908915}},
			},
			"marketdata": map[string]interface{}{
				"columns": []string{"SECID", "LAST", "CHANGE", "LASTCHANGEPRCNT", "VALTODAY"},
				"data":    [][]interface{}{{"SBER", 255.3, 4.8, 1.92, 15000000000}},
			},
		})
	})

	mux.HandleFunc("/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/candles.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"candles": map[string]interface{}{
				"columns": []string{"open", "close", "high", "low", "volume", "begin"},
				"data":    [][]interface{}{{250.0, 255.0, 258.0, 248.0, 1000000, "2026-07-01 00:00:00"}},
			},
		})
	})

	mux.HandleFunc("/iss/securities/SBER.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "ISIN", "ISSUESIZE", "FACEVALUE", "FACEUNIT", "ISSUEDATE", "SECTYPE", "EMITTER_NAME"},
				"data":    [][]interface{}{{"SBER", "RU0009029540", 22586908915.0, 3.0, "SUR", "2007-07-20", "1", "ПАО Сбербанк"}},
			},
		})
	})

	mux.HandleFunc("/iss/securities.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "SHORTNAME", "SECTYPE", "ISIN"},
				"data":    [][]interface{}{{"SBER", "Сбербанк", "1", "RU0009029540"}},
			},
		})
	})

	srv := httptest.NewServer(mux)

	client := moex.NewClient("")
	client.BaseURL = srv.URL

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	server := NewServer(client, logger)

	return server, srv
}

func makeRequest(method string, params interface{}) *JSONRPCRequest {
	rawParams, _ := json.Marshal(params)
	id := json.RawMessage(`"1"`)
	return &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  rawParams,
	}
}

func TestInitialize(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("initialize", nil)
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}
	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("expected protocolVersion 2024-11-05, got %v", result["protocolVersion"])
	}
}

func TestToolsList(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/list", nil)
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}

	result := resp.Result.(map[string]interface{})
	tools := result["tools"].([]map[string]interface{})

	if len(tools) != 4 {
		t.Fatalf("expected 4 tools, got %d", len(tools))
	}

	names := make(map[string]bool)
	for _, tool := range tools {
		names[tool["name"].(string)] = true
	}

	for _, expected := range []string{"moex_ticker", "moex_candles", "moex_fundamentals", "moex_search"} {
		if !names[expected] {
			t.Errorf("missing tool: %s", expected)
		}
	}
}

func TestToolCallTicker(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/call", map[string]interface{}{
		"name":      "moex_ticker",
		"arguments": map[string]interface{}{"symbol": "SBER"},
	})
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}

	result := resp.Result.(map[string]interface{})
	content := result["content"].([]map[string]interface{})
	text := content[0]["text"].(string)

	if text == "" {
		t.Fatal("expected non-empty text")
	}
}

func TestToolCallUnknown(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/call", map[string]interface{}{
		"name":      "unknown_tool",
		"arguments": map[string]interface{}{},
	})
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
	if resp.Error.Code != -32602 {
		t.Errorf("expected code -32602, got %d", resp.Error.Code)
	}
}

func TestUnknownMethod(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("unknown/method", nil)
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error == nil {
		t.Fatal("expected error")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("expected code -32601, got %d", resp.Error.Code)
	}
}

func TestPing(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("ping", nil)
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}
}

func TestNotificationsInitialized(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("notifications/initialized", nil)
	resp := s.HandleRequest(req)

	if resp != nil {
		t.Fatal("expected nil response for notification")
	}
}

func TestToolCallCandles(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/call", map[string]interface{}{
		"name":      "moex_candles",
		"arguments": map[string]interface{}{"symbol": "SBER", "period": "3m"},
	})
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}

	result := resp.Result.(map[string]interface{})
	if result["isError"] != nil {
		t.Fatal("expected success, got error")
	}
}

func TestToolCallFundamentals(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/call", map[string]interface{}{
		"name":      "moex_fundamentals",
		"arguments": map[string]interface{}{"symbol": "SBER"},
	})
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}
}

func TestToolCallSearch(t *testing.T) {
	s, srv := newTestServer(t)
	defer srv.Close()

	req := makeRequest("tools/call", map[string]interface{}{
		"name":      "moex_search",
		"arguments": map[string]interface{}{"query": "SBER"},
	})
	resp := s.HandleRequest(req)

	if resp == nil {
		t.Fatal("expected response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Message)
	}
}
