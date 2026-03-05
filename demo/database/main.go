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
		port = "8083"
	}

	log.Printf("Demo Database starting on port %s", port)
	log.Printf("Health check at http://localhost:%s/health", port)

	// Start background processes
	go startDatabaseOperations()
	go startMetricsGenerator()
	go startHealthCheck(port)

	// Keep the service running
	select {}
}

func startDatabaseOperations() {
	log.Println("Starting database operations...")

	for {
		operation := getRandomOperation()
		processingTime := time.Duration(rand.Intn(100)+10) * time.Millisecond

		log.Printf("🗄️ Executing %s (will take %v)", operation, processingTime)

		// Simulate database operation
		time.Sleep(processingTime)

		// Random success/failure
		if rand.Float32() < 0.98 { // 98% success rate
			log.Printf("✅ %s completed successfully", operation)
		} else {
			log.Printf("❌ %s failed - connection timeout", operation)
		}

		// Wait between operations
		time.Sleep(time.Duration(rand.Intn(2000)+500) * time.Millisecond)
	}
}

func startMetricsGenerator() {
	log.Println("Starting database metrics generator...")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate random database metrics
			connections := rand.Intn(50) + 10
			queriesPerSec := rand.Intn(1000) + 500
			slowQueries := rand.Intn(10) + 1
			cacheHitRate := rand.Float64()*30 + 70 // 70-100%
			diskUsage := rand.Float64()*60 + 40    // 40-100%

			log.Printf("🗄️ DB Metrics - Connections: %d, Queries/sec: %d, Slow queries: %d, Cache hit: %.1f%%, Disk: %.1f%%",
				connections, queriesPerSec, slowQueries, cacheHitRate, diskUsage)
		}
	}
}

func startHealthCheck(port string) {
	log.Println("Starting database health check server...")

	// Simple HTTP server for health checks
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Demo Database Service Running"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "healthy"
		if rand.Float32() < 0.03 { // 3% chance of unhealthy
			status = "unhealthy"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "%s",
			"timestamp": "%s",
			"service": "demo-database",
			"version": "1.0.0",
			"database": {
				"type": "postgresql",
				"version": "14.0",
				"uptime": "%s",
				"connections": {
					"active": %d,
					"idle": %d,
					"max": %d
				},
				"performance": {
					"queries_per_second": %d,
					"slow_queries": %d,
					"cache_hit_rate": %.2f,
					"disk_usage": %.2f
				}
			}
		}`, status, time.Now().Format(time.RFC3339), time.Since(time.Now()).String(),
			rand.Intn(30)+10, rand.Intn(20)+5, 100,
			rand.Intn(1000)+500, rand.Intn(10)+1, rand.Float64()*30+70, rand.Float64()*60+40)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		metrics := []string{
			"# HELP demo_db_connections_active Active database connections",
			"# TYPE demo_db_connections_active gauge",
			fmt.Sprintf("demo_db_connections_active %d", rand.Intn(30)+10),

			"# HELP demo_db_queries_per_second Queries per second",
			"# TYPE demo_db_queries_per_second gauge",
			fmt.Sprintf("demo_db_queries_per_second %d", rand.Intn(1000)+500),

			"# HELP demo_db_slow_queries_total Slow queries count",
			"# TYPE demo_db_slow_queries_total counter",
			fmt.Sprintf("demo_db_slow_queries_total %d", rand.Intn(100)+50),

			"# HELP demo_db_cache_hit_rate Cache hit rate percentage",
			"# TYPE demo_db_cache_hit_rate gauge",
			fmt.Sprintf("demo_db_cache_hit_rate %.2f", rand.Float64()*30+70),

			"# HELP demo_db_disk_usage Disk usage percentage",
			"# TYPE demo_db_disk_usage gauge",
			fmt.Sprintf("demo_db_disk_usage %.2f", rand.Float64()*60+40),

			"# HELP demo_db_query_duration_seconds Query execution duration",
			"# TYPE demo_db_query_duration_seconds histogram",
			fmt.Sprintf("demo_db_query_duration_seconds_bucket{le=\"0.01\"} %d", rand.Intn(800)+200),
			fmt.Sprintf("demo_db_query_duration_seconds_bucket{le=\"0.05\"} %d", rand.Intn(1200)+800),
			fmt.Sprintf("demo_db_query_duration_seconds_bucket{le=\"0.1\"} %d", rand.Intn(1500)+1200),
			fmt.Sprintf("demo_db_query_duration_seconds_bucket{le=\"0.5\"} %d", rand.Intn(1800)+1500),
			fmt.Sprintf("demo_db_query_duration_seconds_bucket{le=\"+Inf\"} %d", rand.Intn(2000)+1800),
			fmt.Sprintf("demo_db_query_duration_seconds_sum %.2f", rand.Float64()*50+20),
			fmt.Sprintf("demo_db_query_duration_seconds_count %d", rand.Intn(2000)+1800),
		}

		w.Write([]byte(fmt.Sprintf("# Generated at %s\n\n%s",
			time.Now().Format(time.RFC3339),
			fmt.Sprintf("%s\n", metrics))))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Health check server failed: %v", err)
	}
}

func getRandomOperation() string {
	operations := []string{
		"SELECT * FROM users WHERE id = ?",
		"INSERT INTO logs (message, level) VALUES (?, ?)",
		"UPDATE metrics SET value = ? WHERE name = ?",
		"DELETE FROM sessions WHERE expires < NOW()",
		"CREATE INDEX idx_users_email ON users(email)",
		"ANALYZE TABLE metrics",
		"VACUUM FULL logs",
		"REINDEX TABLE traces",
		"BEGIN TRANSACTION",
		"COMMIT TRANSACTION",
		"ROLLBACK TRANSACTION",
		"CHECKPOINT",
	}

	return operations[rand.Intn(len(operations))]
}
