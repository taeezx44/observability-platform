package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yourname/observability/collector/scraper"
	"github.com/yourname/observability/collector/storage"
)

func main() {
	// Configuration
	clickhouseURL := getEnv("CLICKHOUSE_URL", "clickhouse:9000")
	scrapeInterval := getEnvDuration("SCRAPE_INTERVAL", 15*time.Second)
	targetsEnv := getEnv("SCRAPE_TARGETS", "http://localhost:8080/metrics")
	
	// Parse targets
	targets := strings.Split(targetsEnv, ",")
	for i, target := range targets {
		targets[i] = strings.TrimSpace(target)
	}

	// Connect to ClickHouse
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s/%s", 
		getEnv("CLICKHOUSE_USER", "admin"),
		getEnv("CLICKHOUSE_PASSWORD", "secret"),
		clickhouseURL,
		getEnv("CLICKHOUSE_DB", "observability"),
	)

	storage, err := storage.NewClickHouseStorage(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer storage.Close()

	log.Printf("Connected to ClickHouse at %s", clickhouseURL)
	log.Printf("Scraping %d targets every %v", len(targets), scrapeInterval)

	// Create and run scraper
	s := &scraper.Scraper{
		Targets:  targets,
		Interval: scrapeInterval,
		Storage:  storage,
	}

	// Start scraping
	s.Run()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
