package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/taeezx44/observability-platform/collector/storage"
	"github.com/taeezx44/observability-platform/collector/tracer"
)

type TracesHandler struct {
	storage *storage.ClickHouseStorage
}

func NewTracesHandler(s *storage.ClickHouseStorage) *TracesHandler {
	return &TracesHandler{storage: s}
}

func (h *TracesHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/traces", h.GetSlowTraces).Methods("GET")
	router.HandleFunc("/api/traces/{traceId}", h.GetTrace).Methods("GET")
}

// GetTrace returns all spans for a single trace, formatted as a waterfall
func (h *TracesHandler) GetTrace(w http.ResponseWriter, r *http.Request) {
	traceID := mux.Vars(r)["traceId"]
	if traceID == "" {
		http.Error(w, "traceId required", http.StatusBadRequest)
		return
	}

	spans, err := h.storage.GetTrace(traceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build waterfall view
	tree := tracer.BuildTree(spans)
	waterfall := tracer.FlattenTree(tree, "", 0)

	// Calculate offsets relative to root span start
	var minStart int64
	for i, item := range waterfall {
		ms := item.Span.StartTime.UnixMilli()
		if i == 0 || ms < minStart {
			minStart = ms
		}
	}

	type WaterfallSpan struct {
		tracer.SpanWithDepth
		StartOffsetMs float64 `json:"start_offset_ms"`
		DurationMs    float64 `json:"duration_ms"`
	}

	var result []WaterfallSpan
	for _, item := range waterfall {
		result = append(result, WaterfallSpan{
			SpanWithDepth: item,
			StartOffsetMs: float64(item.Span.StartTime.UnixMilli() - minStart),
			DurationMs:    item.Span.DurationMs(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetSlowTraces returns traces that exceed a duration threshold
func (h *TracesHandler) GetSlowTraces(w http.ResponseWriter, r *http.Request) {
	minMs := 500.0
	limit := 50

	if v := r.URL.Query().Get("min_ms"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			minMs = f
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	traces, err := h.storage.GetSlowTraces(minMs, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if traces == nil {
		traces = []storage.TracesSummary{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(traces)
}
