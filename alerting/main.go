package main

import (
	"fmt"
	"log"
	"os"

	"github.com/taeezx44/observability-platform/collector/storage"
)

func main() {
	clickhouseURL := getEnv("CLICKHOUSE_URL", "localhost:9000")
	rulesPath := getEnv("RULES_PATH", "./rules.yaml")
	slackURL := getEnv("SLACK_WEBHOOK_URL", "")

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

	engine, err := NewEngine(store, rulesPath, slackURL)
	if err != nil {
		log.Fatalf("Failed to create alert engine: %v", err)
	}

	engine.Run() // blocks forever
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
