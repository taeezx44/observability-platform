package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// MockMetricsHandler represents a metrics handler for testing
type MockMetricsHandler struct {
	metrics []Metric
}

type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

func NewMockMetricsHandler() *MockMetricsHandler {
	return &MockMetricsHandler{
		metrics: []Metric{
			{
				Name:      "cpu_usage",
				Value:     75.5,
				Timestamp: time.Now(),
				Labels:    map[string]string{"service": "api", "instance": "localhost:8080"},
			},
			{
				Name:      "memory_usage",
				Value:     60.2,
				Timestamp: time.Now(),
				Labels:    map[string]string{"service": "api", "instance": "localhost:8080"},
			},
		},
	}
}

func (h *MockMetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.metrics)
}

func (h *MockMetricsHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *MockMetricsHandler) AddMetric(w http.ResponseWriter, r *http.Request) {
	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.metrics = append(h.metrics, metric)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)
}

// TestMetricsEndpoint tests the metrics API endpoint
func TestMetricsEndpoint(t *testing.T) {
	handler := NewMockMetricsHandler()

	req := httptest.NewRequest("GET", "/api/metrics", nil)
	w := httptest.NewRecorder()
	handler.GetMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var metrics []Metric
	if err := json.NewDecoder(w.Body).Decode(&metrics); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(metrics))
	}

	if metrics[0].Name != "cpu_usage" {
		t.Errorf("Expected first metric name 'cpu_usage', got '%s'", metrics[0].Name)
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	handler := NewMockMetricsHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handler.GetHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

// TestAddMetricEndpoint tests adding a new metric
func TestAddMetricEndpoint(t *testing.T) {
	handler := NewMockMetricsHandler()

	newMetric := Metric{
		Name:      "disk_usage",
		Value:     45.8,
		Timestamp: time.Now(),
		Labels:    map[string]string{"service": "api", "instance": "localhost:8080"},
	}

	metricJSON, _ := json.Marshal(newMetric)
	req := httptest.NewRequest("POST", "/api/metrics", bytes.NewReader(metricJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.AddMetric(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response Metric
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Name != "disk_usage" {
		t.Errorf("Expected metric name 'disk_usage', got '%s'", response.Name)
	}

	// Verify metric was added
	if len(handler.metrics) != 3 {
		t.Errorf("Expected 3 metrics after adding, got %d", len(handler.metrics))
	}
}

// TestInvalidMetric tests handling of invalid metric data
func TestInvalidMetric(t *testing.T) {
	handler := NewMockMetricsHandler()

	invalidJSON := `{"name": "test", "value": "invalid"}`
	req := httptest.NewRequest("POST", "/api/metrics", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.AddMetric(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestMetricsFiltering tests filtering metrics by labels
func TestMetricsFiltering(t *testing.T) {
	handler := NewMockMetricsHandler()

	// Add metrics with different labels
	handler.metrics = append(handler.metrics, Metric{
		Name:      "cpu_usage",
		Value:     80.0,
		Timestamp: time.Now(),
		Labels:    map[string]string{"service": "database", "instance": "localhost:5432"},
	})

	req := httptest.NewRequest("GET", "/api/metrics?service=database", nil)
	w := httptest.NewRecorder()
	handler.GetMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var metrics []Metric
	if err := json.NewDecoder(w.Body).Decode(&metrics); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	// Should only return database metrics
	for _, metric := range metrics {
		if metric.Labels["service"] != "database" {
			t.Errorf("Expected only database metrics, got metric with service '%s'", metric.Labels["service"])
		}
	}
}

// BenchmarkMetricsEndpoint benchmarks the metrics endpoint
func BenchmarkMetricsEndpoint(b *testing.B) {
	handler := NewMockMetricsHandler()

	// Add more metrics for realistic benchmarking
	for i := 0; i < 1000; i++ {
		handler.metrics = append(handler.metrics, Metric{
			Name:      "test_metric",
			Value:     float64(i),
			Timestamp: time.Now(),
			Labels:    map[string]string{"service": "test", "instance": "localhost:8080"},
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/metrics", nil)
		w := httptest.NewRecorder()
		handler.GetMetrics(w, req)
	}
}

// BenchmarkAddMetric benchmarks adding metrics
func BenchmarkAddMetric(b *testing.B) {
	handler := NewMockMetricsHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metric := Metric{
			Name:      "benchmark_metric",
			Value:     float64(i),
			Timestamp: time.Now(),
			Labels:    map[string]string{"service": "benchmark", "instance": "localhost:8080"},
		}

		metricJSON, _ := json.Marshal(metric)
		req := httptest.NewRequest("POST", "/api/metrics", bytes.NewReader(metricJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handler.AddMetric(w, req)
	}
}
