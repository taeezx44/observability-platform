package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Advanced endpoints for portfolio demonstration
func setupAdvancedRoutes() {

	// Service discovery endpoint
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		services := []map[string]interface{}{
			{
				"name":       "api",
				"status":     "running",
				"version":    "1.0.0",
				"port":       8084,
				"uptime":     time.Since(time.Now()).String(),
				"endpoints":  []string{"/", "/health", "/metrics", "/api", "/services"},
				"last_check": time.Now().Format(time.RFC3339),
			},
			{
				"name":       "clickhouse",
				"status":     "healthy",
				"version":    "24.0",
				"port":       8123,
				"uptime":     "2h 15m",
				"endpoints":  []string{"http://localhost:8123"},
				"last_check": time.Now().Format(time.RFC3339),
			},
			{
				"name":       "kafka",
				"status":     "healthy",
				"version":    "7.6.0",
				"port":       9092,
				"uptime":     "2h 15m",
				"endpoints":  []string{"localhost:9092"},
				"last_check": time.Now().Format(time.RFC3339),
			},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"services":  services,
			"total":     len(services),
			"healthy":   3,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// System metrics endpoint
	http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		system := map[string]interface{}{
			"cpu": map[string]interface{}{
				"usage_percent": rand.Float64()*80 + 20,
				"cores":         4,
				"load_average":  []float64{1.2, 1.5, 1.8},
			},
			"memory": map[string]interface{}{
				"total_mb":      16384,
				"used_mb":       rand.Intn(8000) + 4000,
				"available_mb":  rand.Intn(8000) + 4000,
				"usage_percent": rand.Float64()*70 + 30,
			},
			"disk": map[string]interface{}{
				"total_gb":      500,
				"used_gb":       rand.Intn(200) + 100,
				"available_gb":  rand.Intn(300) + 200,
				"usage_percent": rand.Float64()*60 + 20,
			},
			"network": map[string]interface{}{
				"bytes_sent":       rand.Int63n(1000000) + 500000,
				"bytes_received":   rand.Int63n(1000000) + 500000,
				"packets_sent":     rand.Intn(10000) + 5000,
				"packets_received": rand.Intn(10000) + 5000,
			},
			"uptime":    time.Since(time.Now()).String(),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(system)
	})

	// Performance metrics endpoint
	http.HandleFunc("/performance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		performance := map[string]interface{}{
			"requests_per_second":   rand.Intn(1000) + 500,
			"avg_response_time":     rand.Float64()*200 + 50,
			"error_rate":            rand.Float64() * 5,
			"throughput_mb_per_sec": rand.Float64()*10 + 5,
			"active_connections":    rand.Intn(100) + 50,
			"cache_hit_rate":        rand.Float64()*30 + 70,
			"database_queries": map[string]interface{}{
				"per_second":   rand.Intn(500) + 200,
				"avg_time_ms":  rand.Float64()*100 + 10,
				"slow_queries": rand.Intn(10),
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(performance)
	})

	// Alert status endpoint
	http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		alerts := []map[string]interface{}{
			{
				"id":        "alert-001",
				"severity":  "warning",
				"service":   "api",
				"metric":    "cpu_usage",
				"threshold": 80,
				"current":   rand.Float64()*30 + 70,
				"status":    "active",
				"message":   "CPU usage is above threshold",
				"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			},
			{
				"id":        "alert-002",
				"severity":  "info",
				"service":   "database",
				"metric":    "connections",
				"threshold": 50,
				"current":   rand.Intn(30) + 20,
				"status":    "resolved",
				"message":   "Database connections normalized",
				"timestamp": time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
			},
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"alerts":    alerts,
			"total":     len(alerts),
			"active":    1,
			"resolved":  1,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Logs endpoint
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		levels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
		services := []string{"api", "database", "kafka", "collector"}

		var logs []map[string]interface{}
		for i := 0; i < 50; i++ {
			logs = append(logs, map[string]interface{}{
				"id":        fmt.Sprintf("log-%06d", i+1),
				"timestamp": time.Now().Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
				"level":     levels[rand.Intn(len(levels))],
				"service":   services[rand.Intn(len(services))],
				"message":   getRandomLogMessage(),
				"trace_id":  fmt.Sprintf("trace_%x", rand.Int63()),
				"span_id":   fmt.Sprintf("span_%x", rand.Int63()),
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"logs":      logs,
			"total":     len(logs),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Dashboard data endpoint
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dashboard := map[string]interface{}{
			"overview": map[string]interface{}{
				"total_services":   4,
				"healthy_services": 3,
				"total_alerts":     2,
				"active_alerts":    1,
				"uptime_percent":   99.9,
			},
			"metrics": map[string]interface{}{
				"requests_per_second": rand.Intn(1000) + 500,
				"error_rate":          rand.Float64() * 2,
				"avg_response_time":   rand.Float64()*100 + 50,
				"throughput":          rand.Float64()*50 + 25,
			},
			"top_services": []map[string]interface{}{
				{"name": "api", "requests": rand.Intn(10000) + 5000, "errors": rand.Intn(100)},
				{"name": "database", "queries": rand.Intn(5000) + 2000, "slow": rand.Intn(50)},
				{"name": "kafka", "messages": rand.Intn(20000) + 10000, "lag": rand.Intn(100)},
			},
			"recent_alerts": []map[string]interface{}{
				{"severity": "warning", "service": "api", "message": "High CPU usage"},
				{"severity": "info", "service": "database", "message": "Connection pool optimized"},
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(dashboard)
	})

	log.Println("Advanced routes registered:")
	log.Println("  /services - Service discovery")
	log.Println("  /system - System metrics")
	log.Println("  /performance - Performance metrics")
	log.Println("  /alerts - Alert status")
	log.Println("  /logs - Recent logs")
	log.Println("  /dashboard - Dashboard data")
}

func getRandomLogMessage() string {
	messages := []string{
		"HTTP request processed successfully",
		"Database query executed",
		"Cache hit for key",
		"Authentication successful",
		"Background job completed",
		"Metrics collected and stored",
		"Alert rule evaluated",
		"Configuration reloaded",
		"Health check passed",
		"Connection established",
		"Request timeout occurred",
		"Memory usage threshold reached",
		"Disk space warning",
		"Network latency increased",
		"Service restarted automatically",
	}

	return messages[rand.Intn(len(messages))]
}
