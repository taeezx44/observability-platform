package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/taeezx44/observability-platform/collector/storage"
)

type LogsHandler struct {
	storage *storage.ClickHouseStorage
}

func NewLogsHandler(s *storage.ClickHouseStorage) *LogsHandler {
	return &LogsHandler{storage: s}
}

func (h *LogsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/logs", h.GetLogs).Methods("GET")
	router.HandleFunc("/api/logs/services", h.GetServices).Methods("GET")
}

func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := storage.LogsQuery{
		From:  time.Now().Add(-1 * time.Hour),
		To:    time.Now(),
		Limit: 200,
	}

	if v := q.Get("level"); v != "" {
		query.Level = v
	}
	if v := q.Get("service"); v != "" {
		query.Service = v
	}
	if v := q.Get("search"); v != "" {
		query.Search = v
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			query.From = t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			query.To = t
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			query.Limit = n
		}
	}

	logs, err := h.storage.GetLogs(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if logs == nil {
		w.Write([]byte("[]"))
		return
	}
	json.NewEncoder(w).Encode(logs)
}

func (h *LogsHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	// TODO: query distinct services from logs table
	services := []string{"api", "collector", "alerting"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}
