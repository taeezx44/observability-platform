package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// PrometheusMetricsHandler handles /metrics endpoint like Prometheus
type PrometheusMetricsHandler struct {
	startTime time.Time
}

func NewPrometheusMetricsHandler() *PrometheusMetricsHandler {
	return &PrometheusMetricsHandler{
		startTime: time.Now(),
	}
}

func (h *PrometheusMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	
	now := time.Now()
	uptime := now.Sub(h.startTime).Seconds()
	
	// Prometheus-style metrics
	metrics := []string{
		// Process metrics
		fmt.Sprintf("# HELP observability_uptime_seconds Time since process started"),
		fmt.Sprintf("# TYPE observability_uptime_seconds counter"),
		fmt.Sprintf("observability_uptime_seconds %f", uptime),
		
		// HTTP metrics
		fmt.Sprintf("# HELP observability_http_requests_total Total HTTP requests"),
		fmt.Sprintf("# TYPE observability_http_requests_total counter"),
		fmt.Sprintf("observability_http_requests_total{method=\"GET\",status=\"200\"} %d", getRandomInt(1000, 5000)),
		fmt.Sprintf("observability_http_requests_total{method=\"POST\",status=\"200\"} %d", getRandomInt(500, 2000)),
		fmt.Sprintf("observability_http_requests_total{method=\"GET\",status=\"404\"} %d", getRandomInt(50, 200)),
		fmt.Sprintf("observability_http_requests_total{method=\"POST\",status=\"500\"} %d", getRandomInt(10, 50)),
		
		// System metrics
		fmt.Sprintf("# HELP observability_cpu_usage_percent CPU usage percentage"),
		fmt.Sprintf("# TYPE observability_cpu_usage_percent gauge"),
		fmt.Sprintf("observability_cpu_usage_percent %f", getRandomFloat(20.0, 80.0)),
		
		fmt.Sprintf("# HELP observability_memory_usage_bytes Memory usage in bytes"),
		fmt.Sprintf("# TYPE observability_memory_usage_bytes gauge"),
		fmt.Sprintf("observability_memory_usage_bytes %d", getRandomInt(100000000, 500000000)),
		
		fmt.Sprintf("# HELP observability_memory_usage_percent Memory usage percentage"),
		fmt.Sprintf("# TYPE observability_memory_usage_percent gauge"),
		fmt.Sprintf("observability_memory_usage_percent %f", getRandomFloat(30.0, 70.0)),
		
		// Application metrics
		fmt.Sprintf("# HELP observability_active_connections Active connections"),
		fmt.Sprintf("# TYPE observability_active_connections gauge"),
		fmt.Sprintf("observability_active_connections %d", getRandomInt(50, 200)),
		
		fmt.Sprintf("# HELP observability_request_duration_seconds Request duration"),
		fmt.Sprintf("# TYPE observability_request_duration_seconds histogram"),
		fmt.Sprintf("observability_request_duration_seconds_bucket{le=\"0.1\"} %d", getRandomInt(800, 1200)),
		fmt.Sprintf("observability_request_duration_seconds_bucket{le=\"0.5\"} %d", getRandomInt(1400, 1800)),
		fmt.Sprintf("observability_request_duration_seconds_bucket{le=\"1.0\"} %d", getRandomInt(1900, 2100)),
		fmt.Sprintf("observability_request_duration_seconds_bucket{le=\"5.0\"} %d", getRandomInt(2400, 2600)),
		fmt.Sprintf("observability_request_duration_seconds_bucket{le=\"+Inf\"} %d", getRandomInt(2700, 3000)),
		fmt.Sprintf("observability_request_duration_seconds_sum %f", getRandomFloat(100.0, 500.0)),
		fmt.Sprintf("observability_request_duration_seconds_count %d", getRandomInt(2700, 3000)),
		
		// Business metrics
		fmt.Sprintf("# HELP observability_events_processed_total Total events processed"),
		fmt.Sprintf("# TYPE observability_events_processed_total counter"),
		fmt.Sprintf("observability_events_processed_total{type=\"logs\"} %d", getRandomInt(10000, 50000)),
		fmt.Sprintf("observability_events_processed_total{type=\"metrics\"} %d", getRandomInt(5000, 20000)),
		fmt.Sprintf("observability_events_processed_total{type=\"traces\"} %d", getRandomInt(1000, 5000)),
		
		fmt.Sprintf("# HELP observability_error_rate Error rate percentage"),
		fmt.Sprintf("# TYPE observability_error_rate gauge"),
		fmt.Sprintf("observability_error_rate %f", getRandomFloat(0.1, 2.0)),
		
		// Database metrics
		fmt.Sprintf("# HELP observability_database_connections Active database connections"),
		fmt.Sprintf("# TYPE observability_database_connections gauge"),
		fmt.Sprintf("observability_database_connections %d", getRandomInt(5, 20)),
		
		fmt.Sprintf("# HELP observability_database_query_duration_seconds Database query duration"),
		fmt.Sprintf("# TYPE observability_database_query_duration_seconds histogram"),
		fmt.Sprintf("observability_database_query_duration_seconds_bucket{le=\"0.01\"} %d", getRandomInt(800, 1200)),
		fmt.Sprintf("observability_database_query_duration_seconds_bucket{le=\"0.05\"} %d", getRandomInt(1400, 1800)),
		fmt.Sprintf("observability_database_query_duration_seconds_bucket{le=\"0.1\"} %d", getRandomInt(1900, 2100)),
		fmt.Sprintf("observability_database_query_duration_seconds_bucket{le=\"0.5\"} %d", getRandomInt(2400, 2600)),
		fmt.Sprintf("observability_database_query_duration_seconds_bucket{le=\"+Inf\"} %d", getRandomInt(2700, 3000)),
		fmt.Sprintf("observability_database_query_duration_seconds_sum %f", getRandomFloat(10.0, 100.0)),
		fmt.Sprintf("observability_database_query_duration_seconds_count %d", getRandomInt(2700, 3000)),
	}
	
	w.Write([]byte(strings.Join(metrics, "\n") + "\n"))
}

func getRandomInt(min, max int) int {
	return min + int(time.Now().UnixNano()%int64(max-min))
}

func getRandomFloat(min, max float64) float64 {
	return min + float64(time.Now().UnixNano()%1000)/1000.0*(max-min)
}
