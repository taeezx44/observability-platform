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
		port = "8082"
	}

	log.Printf("Demo Worker starting on port %s", port)
	log.Printf("Health check at http://localhost:%s/health", port)

	// Start background workers
	go startJobProcessor()
	go startMetricsGenerator()
	go startHealthCheck(port)

	// Keep the service running
	select {}
}

func startJobProcessor() {
	log.Println("Starting job processor...")

	for {
		jobType := getRandomJobType()
		processingTime := time.Duration(rand.Intn(2000)+500) * time.Millisecond

		log.Printf("Processing %s job (will take %v)", jobType, processingTime)

		// Simulate job processing
		time.Sleep(processingTime)

		// Random success/failure
		if rand.Float32() < 0.95 { // 95% success rate
			log.Printf("✅ %s job completed successfully", jobType)
		} else {
			log.Printf("❌ %s job failed", jobType)
		}

		// Wait between jobs
		time.Sleep(time.Duration(rand.Intn(3000)+1000) * time.Millisecond)
	}
}

func startMetricsGenerator() {
	log.Println("Starting metrics generator...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate random metrics
			cpuUsage := rand.Float64()*80 + 20
			memoryUsage := rand.Float64()*70 + 30
			queueDepth := rand.Intn(100) + 10
			processedJobs := rand.Intn(50) + 20

			log.Printf("📊 Metrics - CPU: %.1f%%, Memory: %.1f%%, Queue: %d, Jobs/sec: %d",
				cpuUsage, memoryUsage, queueDepth, processedJobs)
		}
	}
}

func startHealthCheck(port string) {
	log.Println("Starting health check server...")

	// Simple HTTP server for health checks
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "healthy"
		if rand.Float32() < 0.05 { // 5% chance of unhealthy
			status = "unhealthy"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"status": "%s",
			"timestamp": "%s",
			"service": "demo-worker",
			"version": "1.0.0",
			"metrics": {
				"jobs_processed": %d,
				"queue_depth": %d,
				"cpu_usage": %.2f,
				"memory_usage": %.2f
			}
		}`, status, time.Now().Format(time.RFC3339),
			rand.Intn(10000)+5000, rand.Intn(100)+10,
			rand.Float64()*80+20, rand.Float64()*70+30)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		metrics := []string{
			"# HELP demo_worker_jobs_processed_total Total jobs processed",
			"# TYPE demo_worker_jobs_processed_total counter",
			fmt.Sprintf("demo_worker_jobs_processed_total %d", rand.Intn(10000)+5000),

			"# HELP demo_worker_queue_depth Current queue depth",
			"# TYPE demo_worker_queue_depth gauge",
			fmt.Sprintf("demo_worker_queue_depth %d", rand.Intn(100)+10),

			"# HELP demo_worker_cpu_usage CPU usage percentage",
			"# TYPE demo_worker_cpu_usage gauge",
			fmt.Sprintf("demo_worker_cpu_usage %.2f", rand.Float64()*80+20),

			"# HELP demo_worker_memory_usage Memory usage percentage",
			"# TYPE demo_worker_memory_usage gauge",
			fmt.Sprintf("demo_worker_memory_usage %.2f", rand.Float64()*70+30),

			"# HELP demo_worker_job_duration_seconds Job processing duration",
			"# TYPE demo_worker_job_duration_seconds histogram",
			fmt.Sprintf("demo_worker_job_duration_seconds_bucket{le=\"1.0\"} %d", rand.Intn(800)+200),
			fmt.Sprintf("demo_worker_job_duration_seconds_bucket{le=\"2.0\"} %d", rand.Intn(1200)+800),
			fmt.Sprintf("demo_worker_job_duration_seconds_bucket{le=\"5.0\"} %d", rand.Intn(1500)+1200),
			fmt.Sprintf("demo_worker_job_duration_seconds_bucket{le=\"+Inf\"} %d", rand.Intn(2000)+1500),
			fmt.Sprintf("demo_worker_job_duration_seconds_sum %.2f", rand.Float64()*100+50),
			fmt.Sprintf("demo_worker_job_duration_seconds_count %d", rand.Intn(2000)+1500),
		}

		w.Write([]byte(fmt.Sprintf("# Generated at %s\n\n%s",
			time.Now().Format(time.RFC3339),
			fmt.Sprintf("%s\n", metrics))))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Health check server failed: %v", err)
	}
}

func getRandomJobType() string {
	jobTypes := []string{
		"data_processing",
		"email_sending",
		"image_resize",
		"report_generation",
		"cache_warmup",
		"cleanup_task",
		"backup_job",
		"notification_push",
	}

	return jobTypes[rand.Intn(len(jobTypes))]
}
