package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Simulate some initial load
	go simulateLoad()

	mux := http.NewServeMux()

	// Root route - Portfolio friendly
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Demo API Service Running"))
	})

	// API endpoints
	mux.HandleFunc("/info", rootHandler)
	mux.HandleFunc("/api/users", usersHandler)
	mux.HandleFunc("/api/orders", ordersHandler)
	mux.HandleFunc("/api/products", productsHandler)
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	// Middleware for metrics
	handler := metricsMiddleware(mux)

	log.Printf("Demo API starting on port %s", port)
	log.Printf("Metrics available at http://localhost:%s/metrics", port)
	log.Printf("Health check at http://localhost:%s/api/health", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		// Log the request (in real app, this would go to metrics system)
		log.Printf("Request: %s %s - Status: %d - Duration: %.3fs",
			r.Method, r.URL.Path, rw.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(50)+10) * time.Millisecond)
	fmt.Fprintf(w, `{"message": "Demo API Service", "timestamp": "%s", "version": "1.0.0"}`, time.Now().Format(time.RFC3339))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(100)+20) * time.Millisecond)

	users := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com", "status": "active"},
		{"id": 2, "name": "Bob", "email": "bob@example.com", "status": "active"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com", "status": "inactive"},
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"users": %d, "data": %v}`, len(users), users)
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(150)+30) * time.Millisecond)

	orders := []map[string]interface{}{
		{"id": 1001, "user_id": 1, "amount": 99.99, "status": "completed"},
		{"id": 1002, "user_id": 2, "amount": 149.99, "status": "pending"},
		{"id": 1003, "user_id": 3, "amount": 79.99, "status": "shipped"},
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"orders": %d, "data": %v}`, len(orders), orders)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(80)+15) * time.Millisecond)

	products := []map[string]interface{}{
		{"id": "p1", "name": "Laptop", "price": 999.99, "stock": 50},
		{"id": "p2", "name": "Mouse", "price": 29.99, "stock": 200},
		{"id": "p3", "name": "Keyboard", "price": 79.99, "stock": 100},
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"products": %d, "data": %v}`, len(products), products)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	if rand.Float32() < 0.1 { // 10% chance of unhealthy
		status = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"status": "%s",
		"timestamp": "%s",
		"uptime": "%s",
		"version": "1.0.0",
		"checks": {
			"database": "ok",
			"cache": "ok",
			"external_api": "ok"
		}
	}`, status, time.Now().Format(time.RFC3339), time.Since(time.Now()).String())
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	now := time.Now()

	// Prometheus-style metrics
	metrics := []string{
		"# HELP demo_api_requests_total Total HTTP requests",
		"# TYPE demo_api_requests_total counter",
		fmt.Sprintf("demo_api_requests_total{method=\"GET\",endpoint=\"/api/users\"} %d", rand.Intn(1000)+500),
		fmt.Sprintf("demo_api_requests_total{method=\"GET\",endpoint=\"/api/orders\"} %d", rand.Intn(800)+300),
		fmt.Sprintf("demo_api_requests_total{method=\"GET\",endpoint=\"/api/products\"} %d", rand.Intn(600)+200),

		"# HELP demo_api_request_duration_seconds HTTP request duration",
		"# TYPE demo_api_request_duration_seconds histogram",
		fmt.Sprintf("demo_api_request_duration_seconds_bucket{le=\"0.1\"} %d", rand.Intn(800)+200),
		fmt.Sprintf("demo_api_request_duration_seconds_bucket{le=\"0.5\"} %d", rand.Intn(1200)+800),
		fmt.Sprintf("demo_api_request_duration_seconds_bucket{le=\"1.0\"} %d", rand.Intn(1500)+1200),
		fmt.Sprintf("demo_api_request_duration_seconds_bucket{le=\"+Inf\"} %d", rand.Intn(2000)+1500),
		fmt.Sprintf("demo_api_request_duration_seconds_sum %.2f", rand.Float64()*100+50),
		fmt.Sprintf("demo_api_request_duration_seconds_count %d", rand.Intn(2000)+1500),

		"# HELP demo_api_active_connections Active connections",
		"# TYPE demo_api_active_connections gauge",
		fmt.Sprintf("demo_api_active_connections %d", rand.Intn(50)+10),

		"# HELP demo_api_cpu_usage CPU usage percentage",
		"# TYPE demo_api_cpu_usage gauge",
		fmt.Sprintf("demo_api_cpu_usage %.2f", rand.Float64()*80+20),

		"# HELP demo_api_memory_usage_bytes Memory usage in bytes",
		"# TYPE demo_api_memory_usage_bytes gauge",
		fmt.Sprintf("demo_api_memory_usage_bytes %d", rand.Intn(100000000)+50000000),
	}

	w.Write([]byte(fmt.Sprintf("# Generated at %s\n\n%s", now.Format(time.RFC3339),
		fmt.Sprintf("%s\n", metrics))))
}

func simulateLoad() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Simulate varying load patterns
			log.Printf("Load simulation: Active connections: %d, CPU: %.1f%%",
				rand.Intn(50)+10, rand.Float64()*80+20)
		}
	}
}
