package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve the frontend static files
	fs := http.FileServer(http.Dir("frontend/dist/"))
	http.Handle("/", fs)

	// Root route for API - Portfolio friendly
	http.HandleFunc("/api-root", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Observability API Running"))
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","service":"observability-platform"}`)
	})

	// API info endpoint
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"name": "Observability Platform",
			"version": "1.0.0",
			"services": ["collector", "api", "alerting", "frontend"],
			"features": ["metrics", "logs", "traces", "alerting"],
			"status": "running"
		}`)
	})

	log.Printf("Observability Platform starting on port %s", port)
	log.Printf("Frontend available at: http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("API info: http://localhost:%s/api", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
