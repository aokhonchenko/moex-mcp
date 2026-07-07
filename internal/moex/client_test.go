 

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

	// /iss/engines/stock/markets/shares/boards/TQBR/securities.json (для GetSectors)
	mux.HandleFunc("/iss/engines/stock/markets/shares/boards/TQBR/securities.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"securities": map[string]interface{}{
				"columns": []string{"SECID", "SHORTNAME", "SECTYPE", "SECTORID"},
				"data": [][]interface{}{
					{"SBER", "Сбербанк", "1", nil},
					{"GAZP", "Газпром", "1", nil},
					{"LKOH", "ЛУКОЙЛ", "1", nil},
				},
			},
			"marketdata": map[string]interface{}{
				"columns": []string{"SECID", "LAST", "CHANGE", "LASTCHANGEPRCNT", "VALTODAY"},
				"data": [][]interface{}{
					{"SBER", 295.0, -3.0, -1.0, 3000000000},
					{"GAZP", 180.0, 2.0, 1.1, 2000000000},
					{"LKOH", 7200.0, 50.0, 0.7, 1500000000},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// /iss/statistics/engines/stock/markets/index/analytics/MOEXFN.json
	mux.HandleFunc("/iss/statistics/engines/stock/markets/index/analytics/MOEXFN.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"analytics": map[string]interface{}{
				"columns": []string{"indexid", "tradedate", "ticker", "shortnames", "secids", "weight"},
				"data": [][]interface{}{
					{"MOEXFN", "2026-07-06", "SBER", "Сбербанк", "SBER", 13.08},
					{"MOEXFN", "2026-07-06", "VTBR", "ВТБ", "VTBR", 13.39},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// /iss/statistics/engines/stock/markets/index/analytics/MOEXOG.json
	mux.HandleFunc("/iss/statistics/engines/stock/markets/index/analytics/MOEXOG.json", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"analytics": map[string]interface{}{
				"columns": []string{"indexid", "tradedate", "ticker", "shortnames", "secids", "weight"},
				"data": [][]interface{}{
					{"MOEXOG", "2026-07-06", "GAZP", "Газпром", "GAZP", 20.25},
					{"MOEXOG", "2026-07-06", "LKOH", "ЛУКОЙЛ", "LKOH", 21.43},
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

func TestGetSectors(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	groups, err := c.GetSectors()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) == 0 {
		t.Fatal("expected at least one sector group")
	}

	// Проверяем, что SBER попал в "Финансовый сектор"
	var financialGroup *SectorGroup
	for i := range groups {
		if groups[i].SectorName == "Финансовый сектор" {
			financialGroup = &groups[i]
			break
		}
	}
	if financialGroup == nil {
		t.Fatal("expected 'Финансовый сектор' group")
	}

	foundSBER := false
	for _, item := range financialGroup.Items {
		if item.Symbol == "SBER" {
			foundSBER = true
			if item.Price != 295.0 {
				t.Errorf("expected SBER price 295.0, got %f", item.Price)
			}
			break
		}
	}
	if !foundSBER {
		t.Error("expected SBER in 'Финансовый сектор'")
	}

	// Проверяем, что GAZP и LKOH попали в "Нефтегазовый сектор"
	var oilGroup *SectorGroup
	for i := range groups {
		if groups[i].SectorName == "Нефтегазовый сектор" {
			oilGroup = &groups[i]
			break
		}
	}
	if oilGroup == nil {
		t.Fatal("expected 'Нефтегазовый сектор' group")
	}
	if oilGroup.Count != 2 {
		t.Errorf("expected 2 items in oil sector, got %d", oilGroup.Count)
	}
}

func TestSectorMappingCache(t *testing.T) {
	srv := mockMOEXServer()
	defer srv.Close()

	c := NewClient("")
	c.BaseURL = srv.URL

	// Первый вызов — загружает маппинг
	_, err := c.GetSectors()
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}

	if c.sectorMapping == nil {
		t.Fatal("expected sectorMapping to be cached")
	}
	if c.sectorMapping["SBER"] != "Финансовый сектор" {
		t.Errorf("expected SBER in Финансовый сектор, got %s", c.sectorMapping["SBER"])
	}
	if c.sectorMapping["GAZP"] != "Нефтегазовый сектор" {
		t.Errorf("expected GAZP in Нефтегазовый сектор, got %s", c.sectorMapping["GAZP"])
	}
}

func TestNewClientCustomBoard(t *testing.T) {
	c := NewClient("TQBS")
	if c.board != "TQBS" {
		t.Errorf("expected TQBS, got %s", c.board)
	}
}
