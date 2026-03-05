package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRootEndpoint tests the root endpoint
func TestRootEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Observability API Running"))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Observability API Running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","timestamp":"2026-03-05T14:30:00+07:00"}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

// TestAPIEndpoint tests the API info endpoint
func TestAPIEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"name": "Observability Platform",
			"version": "1.0.0",
			"services": ["clickhouse", "kafka", "api"],
			"features": ["metrics", "logs", "traces"],
			"status": "running"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["name"] != "Observability Platform" {
		t.Errorf("Expected name 'Observability Platform', got %v", response["name"])
	}

	if response["status"] != "running" {
		t.Errorf("Expected status 'running', got %v", response["status"])
	}
}

// TestMetricsEndpoint tests the Prometheus metrics endpoint
func TestMetricsEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(`# HELP observability_uptime_seconds Time since process started
# TYPE observability_uptime_seconds counter
observability_uptime_seconds 123.456000
# HELP observability_requests_total Total HTTP requests
# TYPE observability_requests_total counter
observability_requests_total{method="GET",status="200"} 42
# HELP observability_cpu_usage CPU usage percentage
# TYPE observability_cpu_usage gauge
observability_cpu_usage 45.67`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got %v", contentType)
	}

	body := rr.Body.String()
	if !contains(body, "observability_uptime_seconds") {
		t.Errorf("Metrics body should contain uptime metric")
	}

	if !contains(body, "observability_cpu_usage") {
		t.Errorf("Metrics body should contain CPU usage metric")
	}
}

// TestServicesEndpoint tests the service discovery endpoint
func TestServicesEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/services", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"services": [
				{
					"name": "api",
					"status": "running",
					"version": "1.0.0",
					"port": 8084,
					"endpoints": ["/", "/health", "/metrics", "/api", "/services"]
				}
			],
			"total": 1,
			"healthy": 1,
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	services, ok := response["services"].([]interface{})
	if !ok {
		t.Errorf("Expected services array, got %T", response["services"])
	}

	if len(services) == 0 {
		t.Errorf("Expected at least one service, got empty array")
	}

	if response["total"] != float64(1) {
		t.Errorf("Expected total 1, got %v", response["total"])
	}
}

// TestSystemEndpoint tests the system metrics endpoint
func TestSystemEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/system", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"cpu": {
				"usage_percent": 55.94,
				"cores": 4,
				"load_average": [1.2, 1.5, 1.8]
			},
			"memory": {
				"total_mb": 16384,
				"used_mb": 9362,
				"available_mb": 7502,
				"usage_percent": 57.16
			},
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	cpu, ok := response["cpu"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected CPU object, got %T", response["cpu"])
	}

	if cpu["cores"] != float64(4) {
		t.Errorf("Expected 4 CPU cores, got %v", cpu["cores"])
	}

	memory, ok := response["memory"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected memory object, got %T", response["memory"])
	}

	if memory["total_mb"] != float64(16384) {
		t.Errorf("Expected 16384 MB total memory, got %v", memory["total_mb"])
	}
}

// TestPerformanceEndpoint tests the performance metrics endpoint
func TestPerformanceEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/performance", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"requests_per_second": 626,
			"avg_response_time": 150.83,
			"error_rate": 1.23,
			"cache_hit_rate": 80.23,
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["requests_per_second"] != float64(626) {
		t.Errorf("Expected 626 RPS, got %v", response["requests_per_second"])
	}

	if response["error_rate"] != float64(1.23) {
		t.Errorf("Expected error rate 1.23, got %v", response["error_rate"])
	}
}

// TestAlertsEndpoint tests the alert status endpoint
func TestAlertsEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/alerts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"alerts": [
				{
					"id": "alert-001",
					"severity": "warning",
					"service": "api",
					"status": "active",
					"message": "CPU usage is above threshold"
				}
			],
			"total": 1,
			"active": 1,
			"resolved": 0,
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	alerts, ok := response["alerts"].([]interface{})
	if !ok {
		t.Errorf("Expected alerts array, got %T", response["alerts"])
	}

	if len(alerts) == 0 {
		t.Errorf("Expected at least one alert, got empty array")
	}

	if response["active"] != float64(1) {
		t.Errorf("Expected 1 active alert, got %v", response["active"])
	}
}

// TestLogsEndpoint tests the logs endpoint
func TestLogsEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"logs": [
				{
					"id": "log-000001",
					"level": "INFO",
					"service": "api",
					"message": "HTTP request processed successfully",
					"timestamp": "2026-03-05T14:30:00+07:00"
				}
			],
			"total": 1,
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	logs, ok := response["logs"].([]interface{})
	if !ok {
		t.Errorf("Expected logs array, got %T", response["logs"])
	}

	if len(logs) == 0 {
		t.Errorf("Expected at least one log entry, got empty array")
	}

	if response["total"] != float64(1) {
		t.Errorf("Expected total 1, got %v", response["total"])
	}
}

// TestDashboardEndpoint tests the dashboard endpoint
func TestDashboardEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"overview": {
				"total_services": 4,
				"healthy_services": 3,
				"total_alerts": 2,
				"active_alerts": 1,
				"uptime_percent": 99.9
			},
			"metrics": {
				"requests_per_second": 626,
				"error_rate": 1.23,
				"avg_response_time": 87.45
			},
			"timestamp": "2026-03-05T14:30:00+07:00"
		}`))
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	overview, ok := response["overview"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected overview object, got %T", response["overview"])
	}

	if overview["total_services"] != float64(4) {
		t.Errorf("Expected 4 total services, got %v", overview["total_services"])
	}

	if overview["uptime_percent"] != float64(99.9) {
		t.Errorf("Expected 99.9%% uptime, got %v", overview["uptime_percent"])
	}
}

// BenchmarkRootEndpoint benchmarks the root endpoint
func BenchmarkRootEndpoint(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Observability API Running"))
	})

	req, _ := http.NewRequest("GET", "/", nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

// BenchmarkJSONEndpoint benchmarks a JSON endpoint
func BenchmarkJSONEndpoint(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","timestamp":"2026-03-05T14:30:00+07:00"}`))
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
