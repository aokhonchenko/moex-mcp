package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aokhonchenko/moex-mcp/internal/moex"
)

// Server — HTTP-обёртка над MOEX-клиентом.
// Предоставляет REST API для тех же данных, что и MCP-инструменты.
type Server struct {
	client *moex.Client
	logger *log.Logger
}

// NewServer создаёт HTTP-сервер.
func NewServer(client *moex.Client, logger *log.Logger) *Server {
	return &Server{client: client, logger: logger}
}

// Handler возвращает настроенный http.Handler.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/ticker/", s.handleTicker)
	mux.HandleFunc("/api/candles/", s.handleCandles)
	mux.HandleFunc("/api/fundamentals/", s.handleFundamentals)
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/sectors", s.handleSectors)
	mux.HandleFunc("/api/index/", s.handleIndex)
	mux.HandleFunc("/api/cache/stats", s.handleCacheStats)
	mux.HandleFunc("/api/cache/clear", s.handleCacheClear)

	return withCORS(withJSON(mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "moex-mcp",
		"version": "0.3.0",
	})
}

func (s *Server) handleTicker(w http.ResponseWriter, r *http.Request) {
	symbol := extractPathParam(r.URL.Path, "/api/ticker/")
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}

	data, err := s.client.GetTicker(symbol)
	if err != nil {
		s.logger.Printf("ticker error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleCandles(w http.ResponseWriter, r *http.Request) {
	symbol := extractPathParam(r.URL.Path, "/api/candles/")
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "3m"
	}

	data, err := s.client.GetCandles(symbol, period)
	if err != nil {
		s.logger.Printf("candles error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleFundamentals(w http.ResponseWriter, r *http.Request) {
	symbol := extractPathParam(r.URL.Path, "/api/fundamentals/")
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}

	data, err := s.client.GetFundamentals(symbol)
	if err != nil {
		s.logger.Printf("fundamentals error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	data, err := s.client.SearchSecurities(query)
	if err != nil {
		s.logger.Printf("search error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleSectors(w http.ResponseWriter, r *http.Request) {
	data, err := s.client.GetSectors()
	if err != nil {
		s.logger.Printf("sectors error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	symbol := extractPathParam(r.URL.Path, "/api/index/")
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}

	data, err := s.client.GetIndex(symbol)
	if err != nil {
		s.logger.Printf("index error: %v", err)
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := s.client.CacheStats()
	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) handleCacheClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "use POST or DELETE to clear cache")
		return
	}
	s.client.ClearCache()
	writeJSON(w, http.StatusOK, map[string]string{"status": "cache cleared"})
}

// --- helpers ---

// extractPathParam извлекает параметр пути после prefix.
// Например: extractPathParam("/api/ticker/SBER", "/api/ticker/") → "SBER"
func extractPathParam(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	param := strings.TrimPrefix(path, prefix)
	// Убираем возможный trailing slash
	param = strings.TrimSuffix(param, "/")
	return param
}

func withJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
