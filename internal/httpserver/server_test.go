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

		// Sectors: /iss/engines/stock/markets/shares/boards/TQBR/securities.json
		case r.URL.Path == "/iss/engines/stock/markets/shares/boards/TQBR/securities.json":
			fmt.Fprint(w, `{
				"securities":{"columns":["SECID","SHORTNAME","SECTYPE","SECTORID"],"data":[
					["SBER","Сбербанк","1",null],
					["GAZP","Газпром","1",null],
					["LKOH","Лукойл","1",null]
				]},
				"marketdata":{"columns":["SECID","LAST","CHANGE","LASTCHANGEPRCNT","VALTODAY"],"data":[
					["SBER",300.5,2.52,0.8456,5800000000],
					["GAZP",180.0,-1.5,-0.8264,3200000000],
					["LKOH",6500.0,50.0,0.7752,2100000000]
				]}
			}`)

		// Sector index analytics: MOEXFN
		case r.URL.Path == "/iss/statistics/engines/stock/markets/index/analytics/MOEXFN.json":
			fmt.Fprint(w, `{
				"analytics":{"columns":["indexid","tradedate","ticker","shortnames","secids","weight"],"data":[
					["MOEXFN","2026-07-06","SBER","Сбербанк","SBER",13.08]
				]}
			}`)

		// Sector index analytics: MOEXOG
		case r.URL.Path == "/iss/statistics/engines/stock/markets/index/analytics/MOEXOG.json":
			fmt.Fprint(w, `{
				"analytics":{"columns":["indexid","tradedate","ticker","shortnames","secids","weight"],"data":[
					["MOEXOG","2026-07-06","GAZP","Газпром","GAZP",20.25],
					["MOEXOG","2026-07-06","LKOH","Лукойл","LKOH",21.43]
				]}
			}`)

		// Dividends: /iss/securities/SBER/dividends.json
		case r.URL.Path == "/iss/securities/SBER/dividends.json":
			fmt.Fprint(w, `{
				"dividends":{"columns":["SECID","ISIN","registrydate","value","currency","boardid"],"data":[
					["SBER","RU0009029540","2026-07-10",25.0,"SUR","TQBR"],
					["SBER","RU0009029540","2025-07-10",18.7,"SUR","TQBR"]
				]}
			}`)

		// Orderbook: /iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/orderbook.json
		case r.URL.Path == "/iss/engines/stock/markets/shares/boards/TQBR/securities/SBER/orderbook.json":
			fmt.Fprint(w, `{
				"orderbook":{"columns":["PRICE","QUANTITY","BUYSELL"],"data":[
					[300.0,100,"B"],
					[300.0,50,"S"],
					[299.5,200,"B"],
					[300.5,150,"S"]
				]}
			}`)

		// Index: /iss/engines/stock/markets/index/securities/IMOEX.json
		case r.URL.Path == "/iss/engines/stock/markets/index/securities/IMOEX.json":
			fmt.Fprint(w, `{
				"securities":{"columns":["SECID","SHORTNAME","PREVLEGALCLOSEPRICE"],"data":[["IMOEX","Индекс Мосбиржи",2800.5]]},
				"marketdata":{"columns":["SECID","LAST","CHANGE","LASTCHANGEPRCNT","HIGH","LOW","OPEN"],"data":[["IMOEX",2815.3,14.8,0.53,2820.0,2790.0,2800.5]]}
			}`)

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

func TestSectors(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/sectors")
	if err != nil {
		t.Fatalf("sectors request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var sectors []moex.SectorGroup
	json.NewDecoder(resp.Body).Decode(&sectors)

	if len(sectors) != 2 {
		t.Fatalf("expected 2 sectors, got %d", len(sectors))
	}

	// Проверяем что есть Финансовый и Нефтегазовый
	sectorMap := make(map[string]moex.SectorGroup)
	for _, s := range sectors {
		sectorMap[s.SectorName] = s
	}

	fin, ok := sectorMap["Финансовый сектор"]
	if !ok {
		t.Fatal("expected Финансовый сектор")
	}
	if fin.Count != 1 {
		t.Errorf("expected 1 item in Финансовый сектор, got %d", fin.Count)
	}

	oil, ok := sectorMap["Нефтегазовый сектор"]
	if !ok {
		t.Fatal("expected Нефтегазовый сектор")
	}
	if oil.Count != 2 {
		t.Errorf("expected 2 items in Нефтегазовый сектор, got %d", oil.Count)
	}
}

func TestIndex(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/index/IMOEX")
	if err != nil {
		t.Fatalf("index request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var idx moex.IndexData
	json.NewDecoder(resp.Body).Decode(&idx)

	if idx.Symbol != "IMOEX" {
		t.Errorf("expected symbol IMOEX, got %s", idx.Symbol)
	}
	if idx.Value != 2815.3 {
		t.Errorf("expected value 2815.3, got %f", idx.Value)
	}
	if idx.Change != 14.8 {
		t.Errorf("expected change 14.8, got %f", idx.Change)
	}
	if idx.High != 2820.0 {
		t.Errorf("expected high 2820.0, got %f", idx.High)
	}
}

func TestIndexEmpty(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/index/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCacheStats(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/cache/stats")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var stats struct {
		Size   int   `json:"size"`
		Hits   int64 `json:"hits"`
		Misses int64 `json:"misses"`
	}
	json.NewDecoder(resp.Body).Decode(&stats)

	// Кэш пустой — size=0
	if stats.Size != 0 {
		t.Errorf("expected size 0, got %d", stats.Size)
	}
}

func TestCacheStatsAfterTicker(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	// Делаем запрос тикера — заполняет кэш
	resp, err := http.Get(ts.URL + "/api/ticker/SBER")
	if err != nil {
		t.Fatalf("ticker request failed: %v", err)
	}
	resp.Body.Close()

	// Проверяем статистику кэша
	resp, err = http.Get(ts.URL + "/api/cache/stats")
	if err != nil {
		t.Fatalf("stats request failed: %v", err)
	}
	defer resp.Body.Close()

	var stats struct {
		Size   int   `json:"size"`
		Hits   int64 `json:"hits"`
		Misses int64 `json:"misses"`
	}
	json.NewDecoder(resp.Body).Decode(&stats)

	if stats.Size < 1 {
		t.Errorf("expected size >= 1 after ticker request, got %d", stats.Size)
	}
}

func TestCacheClear(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	// Заполняем кэш
	resp, _ := http.Get(ts.URL + "/api/ticker/SBER")
	resp.Body.Close()

	// Очищаем кэш
	req, _ := http.NewRequest("POST", ts.URL+"/api/cache/clear", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("clear request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// Проверяем что кэш пуст
	resp, _ = http.Get(ts.URL + "/api/cache/stats")
	defer resp.Body.Close()

	var stats struct {
		Size int `json:"size"`
	}
	json.NewDecoder(resp.Body).Decode(&stats)

	if stats.Size != 0 {
		t.Errorf("expected size 0 after clear, got %d", stats.Size)
	}
}

func TestDividends(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/dividends/SBER")
	if err != nil {
		t.Fatalf("dividends request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var divs []moex.DividendData
	json.NewDecoder(resp.Body).Decode(&divs)

	if len(divs) != 2 {
		t.Fatalf("expected 2 dividends, got %d", len(divs))
	}
	if divs[0].Value != 25.0 {
		t.Errorf("expected value 25.0, got %f", divs[0].Value)
	}
}

func TestDividendsEmpty(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/dividends/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestOrderBook(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/orderbook/SBER")
	if err != nil {
		t.Fatalf("orderbook request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var ob moex.OrderBookData
	json.NewDecoder(resp.Body).Decode(&ob)

	if ob.Symbol != "SBER" {
		t.Errorf("expected SBER, got %s", ob.Symbol)
	}
	if len(ob.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ob.Entries))
	}
}

func TestOrderBookEmpty(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/orderbook/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCacheClearMethodNotAllowed(t *testing.T) {
	mock := mockMOEXServer()
	defer mock.Close()

	_, ts := newTestServer(mock.URL)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/cache/clear")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}
