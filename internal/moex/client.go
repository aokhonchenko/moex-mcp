package moex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client — HTTP-клиент для MOEX ISS API.
type Client struct {
	client  *http.Client
	board   string
	BaseURL string // экспортируем для тестов
}

// NewClient создаёт клиент MOEX ISS.
func NewClient(board string) *Client {
	if board == "" {
		board = "TQBR"
	}
	return &Client{
		client:  &http.Client{Timeout: 15 * time.Second},
		board:   board,
		BaseURL: "https://iss.moex.com",
	}
}

// TickerData — рыночная котировка.
type TickerData struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volume        int64   `json:"volume"`
	MarketCap     int64   `json:"market_cap"`
	UpdatedAt     string  `json:"updated_at"`
}

// OHLCV — свеча.
type OHLCV struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// FundamentalData — фундаментальные данные.
type FundamentalData struct {
	Symbol      string  `json:"symbol"`
	ISIN        string  `json:"isin"`
	IssueSize   int64   `json:"issue_size"`
	FaceValue   float64 `json:"face_value"`
	Currency    string  `json:"currency"`
	IssueDate   string  `json:"issue_date"`
	SecType     string  `json:"sec_type"`
	EmitterName string  `json:"emitter_name"`
}

// SearchResult — результат поиска.
type SearchResult struct {
	Symbol  string `json:"symbol"`
	Name    string `json:"name"`
	SecType string `json:"sec_type"`
	ISIN    string `json:"isin"`
}

// SectorGroup — группа бумаг по сектору.
type SectorGroup struct {
	SectorID   string       `json:"sector_id"`
	SectorName string       `json:"sector_name"`
	Count      int          `json:"count"`
	AvgChange  float64      `json:"avg_change"`
	Items      []SectorInfo `json:"items"`
}

// SectorInfo — информация об одной бумаге в секторе.
type SectorInfo struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	Volume    int64   `json:"volume"`
}

// issResponse — обобщённый ответ ISS.
type issResponse struct {
	Securities struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"securities"`
	Marketdata struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"marketdata"`
	Candles struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"candles"`
}

// sectorNames — маппинг SECTORID MOEX на читаемые названия секторов.
var sectorNames = map[string]string{
	"Financial":    "Финансовый сектор",
	"Oil and Gas":  "Нефтегазовый сектор",
	"Consumer":     "Потребительский сектор",
	"Electric":     "Электроэнергетика",
	"Metals":       "Металлургия",
	"Telecom":      "Телекоммуникации",
	"Chemicals":    "Химическая промышленность",
	"Construction": "Строительство",
	"Transport":    "Транспорт",
	"Technology":   "Технологии",
	"Real Estate":  "Недвижимость",
	"Other":        "Прочие",
}

// GetTicker получает текущую котировку.
func (c *Client) GetTicker(symbol string) (*TickerData, error) {
	url := fmt.Sprintf(
		"%s/iss/engines/stock/markets/shares/boards/%s/securities/%s.json?iss.meta=off",
		c.BaseURL, c.board, symbol,
	)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Securities.Data) == 0 {
		return nil, fmt.Errorf("нет данных для тикера %s", symbol)
	}

	sec := columnsToMap(resp.Securities.Columns, resp.Securities.Data[0])

	var md map[string]interface{}
	if len(resp.Marketdata.Data) > 0 {
		md = columnsToMap(resp.Marketdata.Columns, resp.Marketdata.Data[0])
	}

	t := &TickerData{
		Symbol:    symbol,
		Name:      getString(sec, "SHORTNAME"),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if md != nil {
		t.Price = getFloat(md, "LAST")
		t.Change = getFloat(md, "CHANGE")
		t.ChangePercent = getFloat(md, "LASTCHANGEPRCNT")
		t.Volume = getInt64(md, "VALTODAY")
	} else {
		t.Price = getFloat(sec, "PREVLEGALCLOSEPRICE")
	}

	t.MarketCap = getInt64(sec, "ISSUESIZE")
	return t, nil
}

// GetCandles получает исторические свечи.
func (c *Client) GetCandles(symbol, period string) ([]OHLCV, error) {
	till := time.Now()
	var from time.Time

	switch period {
	case "1m":
		from = till.AddDate(0, -1, 0)
	case "3m":
		from = till.AddDate(0, -3, 0)
	case "6m":
		from = till.AddDate(0, -6, 0)
	case "1y":
		from = till.AddDate(-1, 0, 0)
	default:
		from = till.AddDate(0, -3, 0)
	}

	url := fmt.Sprintf(
		"%s/iss/engines/stock/markets/shares/boards/%s/securities/%s/candles.json?from=%s&till=%s&interval=24&iss.meta=off",
		c.BaseURL, c.board, symbol,
		from.Format("2006-01-02"), till.Format("2006-01-02"),
	)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Candles.Data) == 0 {
		return nil, fmt.Errorf("нет свечей для %s", symbol)
	}

	candles := make([]OHLCV, 0, len(resp.Candles.Data))
	for _, row := range resp.Candles.Data {
		m := columnsToMap(resp.Candles.Columns, row)
		candles = append(candles, OHLCV{
			Date:   getString(m, "BEGIN"),
			Open:   getFloat(m, "OPEN"),
			High:   getFloat(m, "HIGH"),
			Low:    getFloat(m, "LOW"),
			Close:  getFloat(m, "CLOSE"),
			Volume: getInt64(m, "VOLUME"),
		})
	}

	return candles, nil
}

// GetFundamentals получает фундаментальные данные.
func (c *Client) GetFundamentals(symbol string) (*FundamentalData, error) {
	url := fmt.Sprintf("%s/iss/securities/%s.json?iss.meta=off", c.BaseURL, symbol)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Securities.Data) == 0 {
		return nil, fmt.Errorf("нет данных для %s", symbol)
	}

	var secMap map[string]interface{}
	for _, row := range resp.Securities.Data {
		m := columnsToMap(resp.Securities.Columns, row)
		if getString(m, "SECID") == symbol {
			secMap = m
			break
		}
	}
	if secMap == nil {
		secMap = columnsToMap(resp.Securities.Columns, resp.Securities.Data[0])
	}

	return &FundamentalData{
		Symbol:      symbol,
		ISIN:        getString(secMap, "ISIN"),
		IssueSize:   getInt64(secMap, "ISSUESIZE"),
		FaceValue:   getFloat(secMap, "FACEVALUE"),
		Currency:    getString(secMap, "FACEUNIT"),
		IssueDate:   getString(secMap, "ISSUEDATE"),
		SecType:     getString(secMap, "SECTYPE"),
		EmitterName: getString(secMap, "EMITTER_NAME"),
	}, nil
}

// SearchSecurities ищет бумаги по запросу.
func (c *Client) SearchSecurities(query string) ([]SearchResult, error) {
	url := fmt.Sprintf(
		"%s/iss/securities.json?q=%s&limit=20&iss.meta=off&iss.only=securities",
		c.BaseURL, query,
	)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(resp.Securities.Data))
	for _, row := range resp.Securities.Data {
		m := columnsToMap(resp.Securities.Columns, row)
		secid := getString(m, "SECID")
		if secid == "" {
			continue
		}
		results = append(results, SearchResult{
			Symbol:  secid,
			Name:    getString(m, "SHORTNAME"),
			SecType: getString(m, "SECTYPE"),
			ISIN:    getString(m, "ISIN"),
		})
	}

	return results, nil
}

// GetSectors возвращает секторальную аналитику — группировку бумаг по секторам MOEX.
func (c *Client) GetSectors() ([]SectorGroup, error) {
	url := fmt.Sprintf(
		"%s/iss/engines/stock/markets/shares/boards/%s/securities.json?iss.meta=off&iss.only=securities,marketdata",
		c.BaseURL, c.board,
	)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	// Собираем marketdata в map по SECID
	mdBySymbol := make(map[string]map[string]interface{})
	for _, row := range resp.Marketdata.Data {
		m := columnsToMap(resp.Marketdata.Columns, row)
		secid := getString(m, "SECID")
		if secid != "" {
			mdBySymbol[secid] = m
		}
	}

	// Группируем бумаги по секторам
	sectorMap := make(map[string]*SectorGroup)

	for _, row := range resp.Securities.Data {
		m := columnsToMap(resp.Securities.Columns, row)
		secid := getString(m, "SECID")
		if secid == "" {
			continue
		}
		secType := getString(m, "SECTYPE")
		if secType != "1" && secType != "2" {
			continue // только акции и паи
		}

		sectorID := getString(m, "SECTORID")
		if sectorID == "" {
			sectorID = "Other"
		}

		group, ok := sectorMap[sectorID]
		if !ok {
			name := sectorID
			if n, exists := sectorNames[sectorID]; exists {
				name = n
			}
			group = &SectorGroup{
				SectorID:   sectorID,
				SectorName: name,
				Items:      make([]SectorInfo, 0),
			}
			sectorMap[sectorID] = group
		}

		info := SectorInfo{
			Symbol: secid,
			Name:   getString(m, "SHORTNAME"),
		}

		if md, ok := mdBySymbol[secid]; ok {
			info.Price = getFloat(md, "LAST")
			info.Change = getFloat(md, "CHANGE")
			info.ChangePct = getFloat(md, "LASTCHANGEPRCNT")
			info.Volume = getInt64(md, "VALTODAY")
		}

		group.Items = append(group.Items, info)
		group.Count++
	}

	// Считаем среднее изменение по секторам
	result := make([]SectorGroup, 0, len(sectorMap))
	for _, group := range sectorMap {
		if group.Count > 0 {
			var sumChange float64
			for _, item := range group.Items {
				sumChange += item.ChangePct
			}
			group.AvgChange = sumChange / float64(group.Count)
		}
		result = append(result, *group)
	}

	return result, nil
}

// doGet выполняет GET-запрос и декодирует JSON.
func (c *Client) doGet(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("User-Agent", "moex-mcp/0.1")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса к MOEX ISS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MOEX ISS вернул статус %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("ошибка декодирования: %w", err)
	}

	return nil
}

// columnsToMap объединяет columns и data-строку в map.
// Нормализует имена колонок к uppercase, т.к. MOEX ISS API возвращает
// разный регистр для разных эндпоинтов (UPPERCASE для /boards/.../securities/,
// lowercase для /iss/securities.json?q=).
func columnsToMap(columns []string, data []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(columns))
	for i, col := range columns {
		if i < len(data) {
			m[strings.ToUpper(col)] = data[i]
		}
	}
	return m
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case json.Number:
			f, _ := val.Float64()
			return f
		}
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val)
		case json.Number:
			i, _ := val.Int64()
			return i
		}
	}
	return 0
}
