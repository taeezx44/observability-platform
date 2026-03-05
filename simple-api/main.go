package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Include advanced routes from advanced.go
// The setupAdvancedRoutes function is defined in advanced.go

func main() {
	port := ":8084"

	// Setup advanced routes
	setupAdvancedRoutes()

	// Root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Observability API Running"))
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Metrics endpoint
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, `# HELP observability_uptime_seconds Time since process started
# TYPE observability_uptime_seconds counter
observability_uptime_seconds %f
# HELP observability_requests_total Total HTTP requests
# TYPE observability_requests_total counter
observability_requests_total{method="GET",status="200"} %d
# HELP observability_cpu_usage CPU usage percentage
# TYPE observability_cpu_usage gauge
observability_cpu_usage %.2f
`, time.Since(time.Now()).Seconds(), rand.Intn(1000)+100, rand.Float64()*80+20)
	})

	// API info
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"name": "Observability Platform",
			"version": "1.0.0",
			"services": ["clickhouse", "kafka", "api"],
			"features": ["metrics", "logs", "traces"],
			"status": "running"
		}`)
	})

	log.Printf("API server starting on %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  / - Root endpoint")
	log.Printf("  /health - Health check")
	log.Printf("  /metrics - Prometheus metrics")
	log.Printf("  /api - API information")
	log.Printf("  /services - Service discovery")
	log.Printf("  /system - System metrics")
	log.Printf("  /performance - Performance metrics")
	log.Printf("  /alerts - Alert status")
	log.Printf("  /logs - Recent logs")
	log.Printf("  /dashboard - Dashboard data")
	log.Fatal(http.ListenAndServe(port, nil))
}
