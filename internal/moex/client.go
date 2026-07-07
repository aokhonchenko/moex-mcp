
 
package moex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aokhonchenko/moex-mcp/internal/cache"
)

// Client — HTTP-клиент для MOEX ISS API.
type Client struct {
	client  *http.Client
	board   string
	BaseURL string // экспортируем для тестов

	// Кэш маппинга SECID → сектор (из состава секторальных индексов)
	sectorMapping     map[string]string
	sectorMappingTime time.Time
	sectorMappingTTL  time.Duration

	// In-memory кэш для запросов к ISS API
	dataCache *cache.Cache
}

// NewClient создаёт клиент MOEX ISS.
func NewClient(board string) *Client {
	if board == "" {
		board = "TQBR"
	}
	return &Client{
		client:           &http.Client{Timeout: 15 * time.Second},
		board:            board,
		BaseURL:          "https://iss.moex.com",
		sectorMappingTTL: 1 * time.Hour,
		dataCache:        cache.New(5 * time.Minute),
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

// IndexData — данные по индексу MOEX.
type IndexData struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	UpdatedAt     string  `json:"updated_at"`
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

// sectorIndexMap — маппинг секторальных индексов MOEX на читаемые названия.
// Источник: https://iss.moex.com/iss/statistics/engines/stock/markets/index/analytics.json
var sectorIndexMap = map[string]string{
	"MOEXFN": "Финансовый сектор",
	"MOEXOG": "Нефтегазовый сектор",
	"MOEXMM": "Металлургия и добыча",
	"MOEXIT": "Технологии",
	"MOEXRE": "Недвижимость",
	"MOEXCN": "Потребительский сектор",
	"MOEXCH": "Химия и нефтехимия",
	"MOEXTN": "Транспорт",
	"MOEXEU": "Электроэнергетика",
	"MOEXTL": "Телекоммуникации",
}

// issIndexResponse — ответ ISS для /statistics/engines/stock/markets/index/analytics/{INDEXID}.json
type issIndexResponse struct {
	Analytics struct {
		Columns []string        `json:"columns"`
		Data    [][]interface{} `json:"data"`
	} `json:"analytics"`
}

// loadSectorMapping загружает маппинг SECID → сектор из состава секторальных индексов MOEX.
// Кэшируется на sectorMappingTTL (по умолчанию 1 час).
func (c *Client) loadSectorMapping() (map[string]string, error) {
	if c.sectorMapping != nil && time.Since(c.sectorMappingTime) < c.sectorMappingTTL {
		return c.sectorMapping, nil
	}

	mapping := make(map[string]string)

	for indexID, sectorName := range sectorIndexMap {
		url := fmt.Sprintf(
			"%s/iss/statistics/engines/stock/markets/index/analytics/%s.json?iss.meta=off&iss.only=analytics",
			c.BaseURL, indexID,
		)

		var resp issIndexResponse
		if err := c.doGet(url, &resp); err != nil {
			continue // пропускаем индекс, если не удалось загрузить
		}

		for _, row := range resp.Analytics.Data {
			m := columnsToMap(resp.Analytics.Columns, row)
			ticker := getString(m, "TICKER")
			if ticker != "" {
				mapping[ticker] = sectorName
			}
		}
	}

	c.sectorMapping = mapping
	c.sectorMappingTime = time.Now()
	return mapping, nil
}

// GetTicker получает текущую котировку.
func (c *Client) GetTicker(symbol string) (*TickerData, error) {
	cacheKey := fmt.Sprintf("ticker:%s", symbol)
	if cached, ok := c.dataCache.Get(cacheKey); ok {
		return cached.(*TickerData), nil
	}

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

	// Кэшируем на 1 минуту (котировки быстро меняются)
	c.dataCache.SetWithTTL(cacheKey, t, 1*time.Minute)
	return t, nil
}

// GetCandles получает исторические свечи.
func (c *Client) GetCandles(symbol, period string) ([]OHLCV, error) {
	cacheKey := fmt.Sprintf("candles:%s:%s", symbol, period)
	if cached, ok := c.dataCache.Get(cacheKey); ok {
		return cached.([]OHLCV), nil
	}

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

	// Кэшируем на 5 минут (исторические данные не меняются быстро)
	c.dataCache.Set(cacheKey, candles)
	return candles, nil
}

// GetFundamentals получает фундаментальные данные.
func (c *Client) GetFundamentals(symbol string) (*FundamentalData, error) {
	cacheKey := fmt.Sprintf("fundamentals:%s", symbol)
	if cached, ok := c.dataCache.Get(cacheKey); ok {
		return cached.(*FundamentalData), nil
	}

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

	result := &FundamentalData{
		Symbol:      symbol,
		ISIN:        getString(secMap, "ISIN"),
		IssueSize:   getInt64(secMap, "ISSUESIZE"),
		FaceValue:   getFloat(secMap, "FACEVALUE"),
		Currency:    getString(secMap, "FACEUNIT"),
		IssueDate:   getString(secMap, "ISSUEDATE"),
		SecType:     getString(secMap, "SECTYPE"),
		EmitterName: getString(secMap, "EMITTER_NAME"),
	}

	// Кэшируем на 1 час (фундаментальные данные меняются редко)
	c.dataCache.SetWithTTL(cacheKey, result, 1*time.Hour)
	return result, nil
	return result, nil
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
// Секторы определяются через состав секторальных индексов (MOEXFN, MOEXOG и т.д.),
// т.к. MOEX ISS API не возвращает SECTORID для бумаг на доске TQBR.
func (c *Client) GetSectors() ([]SectorGroup, error) {
	// Загружаем маппинг SECID → сектор из секторальных индексов
	sectorMapping, err := c.loadSectorMapping()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки секторов: %w", err)
	}

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

		// Определяем сектор через маппинг секторальных индексов
		sectorID := sectorMapping[secid]
		if sectorID == "" {
			sectorID = "Прочие"
		}

		group, ok := sectorMap[sectorID]
		if !ok {
			group = &SectorGroup{
				SectorID:   sectorID,
				SectorName: sectorID,
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

// GetIndex получает данные по индексу MOEX (IMOEX, RTSI и др.).
// Использует эндпоинт /iss/engines/stock/markets/index/securities/{symbol}.json
func (c *Client) GetIndex(symbol string) (*IndexData, error) {
	cacheKey := fmt.Sprintf("index:%s", symbol)
	if cached, ok := c.dataCache.Get(cacheKey); ok {
		return cached.(*IndexData), nil
	}

	url := fmt.Sprintf(
		"%s/iss/engines/stock/markets/index/securities/%s.json?iss.meta=off",
		c.BaseURL, symbol,
	)

	var resp issResponse
	if err := c.doGet(url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Securities.Data) == 0 {
		return nil, fmt.Errorf("нет данных для индекса %s", symbol)
	}

	sec := columnsToMap(resp.Securities.Columns, resp.Securities.Data[0])

	idx := &IndexData{
		Symbol:    symbol,
		Name:      getString(sec, "SHORTNAME"),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if len(resp.Marketdata.Data) > 0 {
		md := columnsToMap(resp.Marketdata.Columns, resp.Marketdata.Data[0])
		idx.Value = getFloat(md, "LAST")
		idx.Change = getFloat(md, "CHANGE")
		idx.ChangePercent = getFloat(md, "LASTCHANGEPRCNT")
		idx.High = getFloat(md, "HIGH")
		idx.Low = getFloat(md, "LOW")
		idx.Open = getFloat(md, "OPEN")
	} else {
		idx.Value = getFloat(sec, "PREVLEGALCLOSEPRICE")
	}

	// Кэшируем на 1 минуту (индексы обновляются в реальном времени)
	c.dataCache.SetWithTTL(cacheKey, idx, 1*time.Minute)
	return idx, nil
}

// CacheStats возвращает статистику кэша данных.
func (c *Client) CacheStats() cache.Stats {
	return c.dataCache.Stats()
}

// ClearCache очищает кэш данных.
func (c *Client) ClearCache() {
	c.dataCache.Clear()
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
