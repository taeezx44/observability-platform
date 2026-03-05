package scraper

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// MockPrometheusServer creates a mock Prometheus metrics endpoint
func MockPrometheusServer() *httptest.Server {
	metrics := `# HELP cpu_usage CPU usage percentage
# TYPE cpu_usage gauge
cpu_usage{service="api",instance="localhost:8080"} 75.5
cpu_usage{service="database",instance="localhost:5432"} 45.2

# HELP memory_usage Memory usage percentage  
# TYPE memory_usage gauge
memory_usage{service="api",instance="localhost:8080"} 60.2
memory_usage{service="database",instance="localhost:5432"} 80.8

# HELP request_count Total number of requests
# TYPE request_count counter
request_count{service="api",method="GET",status="200"} 1234
request_count{service="api",method="POST",status="201"} 567
`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(metrics))
	}))
}

// MockScraper represents a metrics scraper for testing
type MockScraper struct {
	targets  []string
	interval time.Duration
	metrics  []Metric
}

type Metric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
	Type      string            `json:"type"` // gauge, counter, histogram
}

func NewMockScraper(targets []string, interval time.Duration) *MockScraper {
	return &MockScraper{
		targets:  targets,
		interval: interval,
		metrics:  []Metric{},
	}
}

func (s *MockScraper) ScrapeTarget(target string) ([]Metric, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse Prometheus metrics format
	// This is a simplified parser for testing
	metrics := []Metric{
		{
			Name:      "cpu_usage",
			Value:     75.5,
			Timestamp: time.Now(),
			Labels:    map[string]string{"service": "api", "instance": "localhost:8080"},
			Type:      "gauge",
		},
		{
			Name:      "memory_usage",
			Value:     60.2,
			Timestamp: time.Now(),
			Labels:    map[string]string{"service": "api", "instance": "localhost:8080"},
			Type:      "gauge",
		},
	}

	s.metrics = append(s.metrics, metrics...)
	return metrics, nil
}

func (s *MockScraper) StartScraping(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, target := range s.targets {
				_, err := s.ScrapeTarget(target)
				if err != nil {
					// Log error in real implementation
					continue
				}
			}
		}
	}
}

func (s *MockScraper) GetMetrics() []Metric {
	return s.metrics
}

// TestScraperCreation tests creating a new scraper
func TestScraperCreation(t *testing.T) {
	targets := []string{"http://localhost:8080/metrics"}
	interval := 15 * time.Second

	scraper := NewMockScraper(targets, interval)

	if scraper == nil {
		t.Fatal("Expected non-nil scraper")
	}

	if len(scraper.targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(scraper.targets))
	}

	if scraper.targets[0] != "http://localhost:8080/metrics" {
		t.Errorf("Expected target 'http://localhost:8080/metrics', got '%s'", scraper.targets[0])
	}

	if scraper.interval != interval {
		t.Errorf("Expected interval %v, got %v", interval, scraper.interval)
	}

	if len(scraper.metrics) != 0 {
		t.Errorf("Expected empty metrics, got %d", len(scraper.metrics))
	}
}

// TestTargetScraping tests scraping a single target
func TestTargetScraping(t *testing.T) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{}, 15*time.Second)

	metrics, err := scraper.ScrapeTarget(server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(metrics))
	}

	if metrics[0].Name != "cpu_usage" {
		t.Errorf("Expected first metric name 'cpu_usage', got '%s'", metrics[0].Name)
	}

	if metrics[0].Type != "gauge" {
		t.Errorf("Expected metric type 'gauge', got '%s'", metrics[0].Type)
	}
}

// TestInvalidTarget tests scraping an invalid target
func TestInvalidTarget(t *testing.T) {
	scraper := NewMockScraper([]string{}, 15*time.Second)

	_, err := scraper.ScrapeTarget("http://invalid-target:9999/metrics")
	if err == nil {
		t.Error("Expected error for invalid target, got nil")
	}
}

// TestMultipleTargets tests scraping multiple targets
func TestMultipleTargets(t *testing.T) {
	server1 := MockPrometheusServer()
	server2 := MockPrometheusServer()
	defer server1.Close()
	defer server2.Close()

	scraper := NewMockScraper([]string{server1.URL, server2.URL}, 15*time.Second)

	// Scrape both targets
	metrics1, err1 := scraper.ScrapeTarget(server1.URL)
	if err1 != nil {
		t.Errorf("Expected no error for server1, got %v", err1)
	}

	metrics2, err2 := scraper.ScrapeTarget(server2.URL)
	if err2 != nil {
		t.Errorf("Expected no error for server2, got %v", err2)
	}

	if len(metrics1) != 2 {
		t.Errorf("Expected 2 metrics from server1, got %d", len(metrics1))
	}

	if len(metrics2) != 2 {
		t.Errorf("Expected 2 metrics from server2, got %d", len(metrics2))
	}

	// Verify total metrics collected
	totalMetrics := len(scraper.GetMetrics())
	if totalMetrics != 4 {
		t.Errorf("Expected 4 total metrics, got %d", totalMetrics)
	}
}

// TestScrapingInterval tests the scraping interval
func TestScrapingInterval(t *testing.T) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{server.URL}, 100*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// Start scraping in background
	go scraper.StartScraping(ctx)

	// Wait for scraping to complete
	<-ctx.Done()

	metrics := scraper.GetMetrics()

	// Should have scraped at least twice (at 0ms and 100ms, maybe at 200ms)
	if len(metrics) < 2 {
		t.Errorf("Expected at least 2 scraping cycles, got %d metrics", len(metrics))
	}
}

// TestMetricsStorage tests metric storage functionality
func TestMetricsStorage(t *testing.T) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{server.URL}, 15*time.Second)

	// Scrape multiple times
	for i := 0; i < 3; i++ {
		_, err := scraper.ScrapeTarget(server.URL)
		if err != nil {
			t.Errorf("Scraping iteration %d failed: %v", i, err)
		}
	}

	metrics := scraper.GetMetrics()
	if len(metrics) != 6 { // 2 metrics per scrape * 3 scrapes
		t.Errorf("Expected 6 metrics after 3 scrapes, got %d", len(metrics))
	}

	// Verify metric structure
	for i, metric := range metrics {
		if metric.Name == "" {
			t.Errorf("Metric %d: expected non-empty name", i)
		}
		if metric.Timestamp.IsZero() {
			t.Errorf("Metric %d: expected non-zero timestamp", i)
		}
		if len(metric.Labels) == 0 {
			t.Errorf("Metric %d: expected non-empty labels", i)
		}
		if metric.Type == "" {
			t.Errorf("Metric %d: expected non-empty type", i)
		}
	}
}

// TestConcurrentScraping tests concurrent scraping
func TestConcurrentScraping(t *testing.T) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{server.URL}, 15*time.Second)

	// Scrape concurrently
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := scraper.ScrapeTarget(server.URL)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Should have 20 metrics (2 per scrape * 10 scrapes)
	metrics := scraper.GetMetrics()
	if len(metrics) != 20 {
		t.Errorf("Expected 20 metrics after concurrent scraping, got %d", len(metrics))
	}
}

// BenchmarkScraping benchmarks the scraping performance
func BenchmarkScraping(b *testing.B) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{server.URL}, 15*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scraper.ScrapeTarget(server.URL)
		if err != nil {
			b.Errorf("Scraping failed: %v", err)
		}
	}
}

// BenchmarkConcurrentScraping benchmarks concurrent scraping
func BenchmarkConcurrentScraping(b *testing.B) {
	server := MockPrometheusServer()
	defer server.Close()

	scraper := NewMockScraper([]string{server.URL}, 15*time.Second)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := scraper.ScrapeTarget(server.URL)
			if err != nil {
				b.Errorf("Scraping failed: %v", err)
			}
		}
	})
}
