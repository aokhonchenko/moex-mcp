package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

// mockMOEXServer возвращает httptest.Server, эмулирующий MOEX ISS API.
func mockMOEXServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		// Ticker: /iss/engines/stock/markets/shares/boards/TQBR/securities/SBER.json
		case r.URL.Path == "/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER.json":
			fmt.Fprint(w, `{
				"securities":{"columns":["SECID","SHORTNAME","ISSUESIZE","PREVLEGALCLOSEPRICE"],"data":[["SBER","Сбербанк",21586974000,297.98]]},
				"marketdata":{"columns":["LAST","CHANGE","LASTCHANGEPRCNT","VALTODAY"],"data":[[300.5,2.52,0.8456,5800000000]]}
			}`)

		// Candles: /iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/candles.json
		case r.URL.Path == "/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/candles.json":
			fmt.Fprint(w, `{
				"candles":{"columns":["open","close","high","low","volume","begin"],"data":[
					[295.0,300.5,302.0,294.0,5000000,"2026-07-01 00:00:00"],
					[300.5,298.0,301.0,297.0,4500000,"2026-07-02 00:00:00"]
				]}
			}`)

		// Fundamentals: /iss/securities/SBER.json
		case r.URL.Path == "/iss/securities/SBER.json":
			fmt.Fprint(w, `{
				"securities":{"columns":["SECID","ISIN","ISSUESIZE","FACEVALUE","FACEUNIT","ISSUEDATE","SECTYPE","EMITTER_NAME"],"data":[
					["SBER","RU0009029540",21586974000,3,"RUR","1996-01-01","1","ПАО Сбербанк"]
				]}
			}`)

		// Search: /iss/securities.json?q=SBER
		case r.URL.Path == "/iss/securities.json":
			q := r.URL.Query().Get("q")
			if q == "SBER" {
				fmt.Fprint(w, `{
					"securities":{"columns":["SECID","SHORTNAME","SECTYPE","ISIN"],"data":[
						["SBER","Сбербанк","1","RU0009029540"],
						["SBERP","Сбербанк-п","2","RU0009029557"]
					]}
				}`)
			} else {
				fmt.Fprint(w, `{"securities":{"columns":["SECID","SHORTNAME","SECTYPE","ISIN"],"data":[]}}`)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":"not found"}`)
		}
	}))
}

func newTestServer(mockURL string) (*Server, *httptest.Server) {
	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	client := moex.NewClient("")
	client.BaseURL = mockURL

	srv := NewServer(client, logger)
	ts := httptest.NewServer(srv.Handler())
	return srv, ts
}

func TestHealth(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)

	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %s", body["status"])
	}
	if body["service"] != "moex-mcp" {
		t.Errorf("expected service=moex-mcp, got %s", body["service"])
	}
}

func TestTicker(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/ticker/SBER")
	if err != nil {
		t.Fatalf("ticker request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var data moex.TickerData
	json.NewDecoder(resp.Body).Decode(&data)

	if data.Symbol != "SBER" {
		t.Errorf("expected symbol SBER, got %s", data.Symbol)
	}
	if data.Price != 300.5 {
		t.Errorf("expected price 300.5, got %f", data.Price)
	}
	if data.Name != "Сбербанк" {
		t.Errorf("expected name Сбербанк, got %s", data.Name)
	}
}

func TestTickerNotFound(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/ticker/UNKNOWN")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", resp.StatusCode)
	}
}

func TestCandles(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/candles/SBER?period=1m")
	if err != nil {
		t.Fatalf("candles request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var candles []moex.OHLCV
	json.NewDecoder(resp.Body).Decode(&candles)

	if len(candles) != 2 {
		t.Fatalf("expected 2 candles, got %d", len(candles))
	}
	if candles[0].Open != 295.0 {
		t.Errorf("expected open 295.0, got %f", candles[0].Open)
	}
}

func TestCandlesDefaultPeriod(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	// Без параметра period — должен использовать 3m
	resp, err := http.Get(ts.URL + "/api/candles/SBER")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestFundamentals(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/fundamentals/SBER")
	if err != nil {
		t.Fatalf("fundamentals request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var fund moex.FundamentalData
	json.NewDecoder(resp.Body).Decode(&fund)

	if fund.Symbol != "SBER" {
		t.Errorf("expected symbol SBER, got %s", fund.Symbol)
	}
	if fund.ISIN != "RU0009029540" {
		t.Errorf("expected ISIN RU0009029540, got %s", fund.ISIN)
	}
	if fund.FaceValue != 3 {
		t.Errorf("expected face value 3, got %f", fund.FaceValue)
	}
}

func TestSearch(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/search?q=SBER")
	if err != nil {
		t.Fatalf("search request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var results []moex.SearchResult
	json.NewDecoder(resp.Body).Decode(&results)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Symbol != "SBER" {
		t.Errorf("expected first result SBER, got %s", results[0].Symbol)
	}
}

func TestSearchMissingQuery(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/search")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCORS(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	// OPTIONS preflight
	req, _ := http.NewRequest("OPTIONS", ts.URL+"/api/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected 200 for OPTIONS, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS header *")
	}
}

func TestContentType(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("expected application/json; charset=utf-8, got %s", ct)
	}
}
