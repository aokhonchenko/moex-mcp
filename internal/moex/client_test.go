package moex

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockMOEXServer возвращает httptest.Server с мокированными ответами MOEX ISS.
func mockMOEXServer() *httptest.Server {
	mux := http.NewServeMux()

	// /iss/engines/stock/markets/shares/boards/TQBR/securities/SBER.json
	mux.HandleFunc("/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "SHORTNAME", "PREVLEGALCLOSEPRICE", "ISSUESIZE"},
				"data":    [][]interface{}{{"SBER", "Сбербанк", 250.5, 22586908915}},
			},
			"marketdata": map[string]interface{}{
				"columns": []string{"SECID", "LAST", "CHANGE", "LASTCHANGEPRCNT", "VALTODAY"},
				"data":    [][]interface{}{{"SBER", 255.3, 4.8, 1.92, 15000000000}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// /iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/candles.json
	mux.HandleFunc("/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/candles.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"candles": map[string]interface{}{
				"columns": []string{"open", "close", "high", "low", "volume", "begin"},
				"data": [][]interface{}{
					{250.0, 255.0, 258.0, 248.0, 1000000, "2026-07-01 00:00:00"},
					{255.0, 260.0, 262.0, 253.0, 1200000, "2026-07-02 00:00:00"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// /iss/securities/SBER.json
	mux.HandleFunc("/iss/securities/SBER.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "ISIN", "ISSUESIZE", "FACEVALUE", "FACEUNIT", "ISSUEDATE", "SECTYPE", "EMITTER_NAME"},
				"data":    [][]interface{}{{"SBER", "RU0009029540", 22586908915.0, 3.0, "SUR", "2007-07-20", "1", "ПАО Сбербанк"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// /iss/securities.json?q=SBER
	mux.HandleFunc("/iss/securities.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "SHORTNAME", "SECTYPE", "ISIN"},
				"data": [][]interface{}{
					{"SBER", "Сбербанк", "1", "RU0009029540"},
					{"SBERP", "Сбербанк-п", "2", "RU0009029557"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	return httptest.NewServer(mux)
}

func TestGetTicker(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	data, err := c.GetTicker("SBER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Symbol != "SBER" {
		t.Errorf("expected SBER, got %s", data.Symbol)
	}
	if data.Name != "Сбербанк" {
		t.Errorf("expected Сбербанк, got %s", data.Name)
	}
	if data.Price != 255.3 {
		t.Errorf("expected 255.3, got %f", data.Price)
	}
	if data.Change != 4.8 {
		t.Errorf("expected 4.8, got %f", data.Change)
	}
	if data.Volume != 15000000000 {
		t.Errorf("expected 15000000000, got %d", data.Volume)
	}
}

func TestGetTickerNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"securities":  map[string]interface{}{"columns": []string{}, "data": [][]interface{}{}},
			"marketdata":  map[string]interface{}{"columns": []string{}, "data": [][]interface{}{}},
		})
	}))
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	_, err := c.GetTicker("XXXX")
	if err == nil {
		t.Fatal("expected error for missing ticker")
	}
}

func TestGetCandles(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	candles, err := c.GetCandles("SBER", "3m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(candles) != 2 {
		t.Fatalf("expected 2 candles, got %d", len(candles))
	}

	if candles[0].Open != 250.0 {
		t.Errorf("expected open 250, got %f", candles[0].Open)
	}
	if candles[0].Close != 255.0 {
		t.Errorf("expected close 255, got %f", candles[0].Close)
	}
	if candles[0].Date != "2026-07-01 00:00:00" {
		t.Errorf("expected date 2026-07-01, got %s", candles[0].Date)
	}
}

func TestGetFundamentals(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	fund, err := c.GetFundamentals("SBER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fund.ISIN != "RU0009029540" {
		t.Errorf("expected ISIN RU0009029540, got %s", fund.ISIN)
	}
	if fund.Currency != "SUR" {
		t.Errorf("expected SUR, got %s", fund.Currency)
	}
	if fund.EmitterName != "ПАО Сбербанк" {
		t.Errorf("expected ПАО Сбербанк, got %s", fund.EmitterName)
	}
}

func TestSearchSecurities(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	results, err := c.SearchSecurities("SBER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Symbol != "SBER" {
		t.Errorf("expected SBER, got %s", results[0].Symbol)
	}
	if results[1].Symbol != "SBERP" {
		t.Errorf("expected SBERP, got %s", results[1].Symbol)
	}
}

func TestHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	_, err := c.GetTicker("SBER")
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
}

func TestNewClientDefaults(t *testing.T) {
	c := NewClient("")
	if c.board != "TQBR" {
		t.Errorf("expected TQBR, got %s", c.board)
	}
	if c.BaseURL != "https://iss.moex.com" {
		t.Errorf("expected https://iss.moex.com, got %s", c.BaseURL)
	}
}

func TestNewClientCustomBoard(t *testing.T) {
	c := NewClient("TQBS")
	if c.board != "TQBS" {
		t.Errorf("expected TQBS, got %s", c.board)
	}
}
