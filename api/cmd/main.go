package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/taeezx44/observability-platform/api/handlers"
	"github.com/taeezx44/observability-platform/collector/storage"
)

func main() {
	port := getEnv("PORT", "8080")
	clickhouseURL := getEnv("CLICKHOUSE_URL", "localhost:9000")

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s/%s",
		getEnv("CLICKHOUSE_USER", "admin"),
		getEnv("CLICKHOUSE_PASSWORD", "secret"),
		clickhouseURL,
		getEnv("CLICKHOUSE_DB", "observability"),
	)

	store, err := storage.NewClickHouseStorage(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer store.Close()

	router := mux.NewRouter()

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Register all handlers
	handlers.NewMetricsHandler(store).RegisterRoutes(router)
	handlers.NewLogsHandler(store).RegisterRoutes(router)
	handlers.NewTracesHandler(store).RegisterRoutes(router)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// WebSocket for live metrics push
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()
		log.Printf("WebSocket client connected")

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metrics, err := store.GetMetrics(storage.MetricsQuery{
				From:  time.Now().Add(-5 * time.Minute),
				To:    time.Now(),
				Limit: 200,
			})
			if err != nil {
				log.Printf("WS metrics error: %v", err)
				continue
			}
			if err := conn.WriteJSON(metrics); err != nil {
				log.Printf("WS write failed: %v", err)
				return
			}
		}
	})

	addr := ":" + port
	log.Printf("API server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
