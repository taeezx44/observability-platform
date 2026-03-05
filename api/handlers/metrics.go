package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/taeezx44/observability-platform/collector/storage"
)

type MetricsHandler struct {
	storage *storage.ClickHouseStorage
}

func NewMetricsHandler(storage *storage.ClickHouseStorage) *MetricsHandler {
	return &MetricsHandler{storage: storage}
}

func (h *MetricsHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/metrics", h.GetMetrics).Methods("GET")
	router.HandleFunc("/api/metrics/names", h.GetMetricNames).Methods("GET")
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := storage.MetricsQuery{
		From: time.Now().Add(-1 * time.Hour), // Default: last hour
		To:   time.Now(),
	}

	if name := r.URL.Query().Get("name"); name != "" {
		query.Name = name
	}

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			query.From = from
		}
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			query.To = to
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			query.Limit = limit
		}
	}

	// Parse label filters
	query.Labels = make(map[string]string)
	for key, values := range r.URL.Query() {
		if strings.HasPrefix(key, "label.") {
			labelKey := strings.TrimPrefix(key, "label.")
			if len(values) > 0 {
				query.Labels[labelKey] = values[0]
			}
		}
	}

	metrics, err := h.storage.GetMetrics(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *MetricsHandler) GetMetricNames(w http.ResponseWriter, r *http.Request) {
	names, err := h.storage.QueryDistinctNames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if names == nil {
		names = []string{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
