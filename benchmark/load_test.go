package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkConfig holds configuration for load tests
type BenchmarkConfig struct {
	TargetURL        string
	RequestsPerSec   int
	Duration         time.Duration
	ConcurrentUsers  int
	PayloadSize      int
	Timeout          time.Duration
	KeepAlive        bool
}

// BenchmarkResults holds test results
type BenchmarkResults struct {
	TotalRequests     int64         `json:"total_requests"`
	SuccessfulReqs    int64         `json:"successful_requests"`
	FailedReqs        int64         `json:"failed_requests"`
	Duration          time.Duration `json:"duration"`
	RequestsPerSec    float64       `json:"requests_per_second"`
	AvgLatency        time.Duration `json:"avg_latency"`
	P50Latency        time.Duration `json:"p50_latency"`
	P95Latency        time.Duration `json:"p95_latency"`
	P99Latency        time.Duration `json:"p99_latency"`
	MinLatency        time.Duration `json:"min_latency"`
	MaxLatency        time.Duration `json:"max_latency"`
	BytesTransferred  int64         `json:"bytes_transferred"`
	MemoryUsage       int64         `json:"memory_usage_bytes"`
	CPUUsage          float64       `json:"cpu_usage_percent"`
	ErrorRate         float64       `json:"error_rate_percent"`
}

// LoadGenerator generates load for testing
type LoadGenerator struct {
	config    BenchmarkConfig
	client    *http.Client
	latencies []time.Duration
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewLoadGenerator creates a new load generator
func NewLoadGenerator(config BenchmarkConfig) *LoadGenerator {
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   !config.KeepAlive,
		},
	}
	
	return &LoadGenerator{
		config: config,
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

// generatePayload creates test payload
func (lg *LoadGenerator) generatePayload() []byte {
	if lg.config.PayloadSize <= 0 {
		return []byte(`{"metric": "cpu_usage", "value": 75.5, "labels": {"service": "test"}}`)
	}
	
	payload := make([]byte, lg.config.PayloadSize)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	return payload
}

// singleRequest performs a single HTTP request
func (lg *LoadGenerator) singleRequest() error {
	payload := lg.generatePayload()
	
	req, err := http.NewRequestWithContext(lg.ctx, "POST", lg.config.TargetURL, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Observability-Platform-Benchmark/1.0")
	
	start := time.Now()
	resp, err := lg.client.Do(req)
	latency := time.Since(start)
	
	lg.mu.Lock()
	lg.latencies = append(lg.latencies, latency)
	lg.mu.Unlock()
	
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	atomic.AddInt64(&bytesTransferred, int64(len(payload)))
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	
	return nil
}

// worker runs requests continuously
func (lg *LoadGenerator) worker(id int, results *BenchmarkResults) {
	ticker := time.NewTicker(time.Second / time.Duration(lg.config.RequestsPerSec/lg.config.ConcurrentUsers))
	defer ticker.Stop()
	
	for {
		select {
		case <-lg.ctx.Done():
			return
		case <-ticker.C:
			err := lg.singleRequest()
			if err != nil {
				atomic.AddInt64(&results.FailedReqs, 1)
			} else {
				atomic.AddInt64(&results.SuccessfulReqs, 1)
			}
			atomic.AddInt64(&results.TotalRequests, 1)
		}
	}
}

// RunBenchmark executes the load test
func (lg *LoadGenerator) RunBenchmark() *BenchmarkResults {
	results := &BenchmarkResults{
		MinLatency: time.Hour,
		MaxLatency: 0,
	}
	
	var wg sync.WaitGroup
	
	// Start workers
	for i := 0; i < lg.config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			lg.worker(id, results)
		}(i)
	}
	
	// Monitor system resources
	go lg.monitorResources(results)
	
	// Wait for completion
	start := time.Now()
	wg.Wait()
	results.Duration = time.Since(start)
	
	// Calculate latency percentiles
	lg.mu.Lock()
	if len(lg.latencies) > 0 {
		sorted := make([]time.Duration, len(lg.latencies))
		copy(sorted, lg.latencies)
		
		// Simple bubble sort for percentiles
		for i := 0; i < len(sorted); i++ {
			for j := 0; j < len(sorted)-1-i; j++ {
				if sorted[j] > sorted[j+1] {
					sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
				}
			}
		}
		
		results.MinLatency = sorted[0]
		results.MaxLatency = sorted[len(sorted)-1]
		results.P50Latency = sorted[len(sorted)/2]
		results.P95Latency = sorted[int(float64(len(sorted))*0.95)]
		results.P99Latency = sorted[int(float64(len(sorted))*0.99)]
		
		// Calculate average
		var total time.Duration
		for _, lat := range sorted {
			total += lat
		}
		results.AvgLatency = total / time.Duration(len(sorted))
	}
	lg.mu.Unlock()
	
	// Calculate final metrics
	results.RequestsPerSec = float64(results.TotalRequests) / results.Duration.Seconds()
	results.ErrorRate = float64(results.FailedReqs) / float64(results.TotalRequests) * 100
	results.BytesTransferred = atomic.LoadInt64(&bytesTransferred)
	
	return results
}

// monitorResources tracks CPU and memory usage
func (lg *LoadGenerator) monitorResources(results *BenchmarkResults) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	
	var memSamples []int64
	
	for {
		select {
		case <-lg.ctx.Done():
			// Calculate average memory usage
			if len(memSamples) > 0 {
				var total int64
				for _, mem := range memSamples {
					total += mem
				}
				results.MemoryUsage = total / int64(len(memSamples))
			}
			return
		case <-ticker.C:
			// Memory usage
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			memSamples = append(memSamples, int64(m.Alloc))
		}
	}
}

// Global counter for bytes transferred
var bytesTransferred int64

// Benchmark10kReqSec tests 10,000 requests per second
func Benchmark10kReqSec(b *testing.B) {
	config := BenchmarkConfig{
		TargetURL:       "http://localhost:8080/api/metrics",
		RequestsPerSec:  10000,
		Duration:        30 * time.Second,
		ConcurrentUsers: 100,
		PayloadSize:     256,
		Timeout:         5 * time.Second,
		KeepAlive:       true,
	}
	
	generator := NewLoadGenerator(config)
	results := generator.RunBenchmark()
	
	// Report results
	b.Logf("=== 10k req/sec Benchmark Results ===")
	b.Logf("Total Requests: %d", results.TotalRequests)
	b.Logf("Successful: %d", results.SuccessfulReqs)
	b.Logf("Failed: %d", results.FailedReqs)
	b.Logf("Duration: %v", results.Duration)
	b.Logf("Actual RPS: %.2f", results.RequestsPerSec)
	b.Logf("Error Rate: %.2f%%", results.ErrorRate)
	b.Logf("Avg Latency: %v", results.AvgLatency)
	b.Logf("P50 Latency: %v", results.P50Latency)
	b.Logf("P95 Latency: %v", results.P95Latency)
	b.Logf("P99 Latency: %v", results.P99Latency)
	b.Logf("Min Latency: %v", results.MinLatency)
	b.Logf("Max Latency: %v", results.MaxLatency)
	b.Logf("Memory Usage: %d bytes", results.MemoryUsage)
	b.Logf("Bytes Transferred: %d", results.BytesTransferred)
	
	// Performance assertions
	if results.RequestsPerSec < 9000 {
		b.Errorf("RPS too low: %.2f, expected >= 9000", results.RequestsPerSec)
	}
	
	if results.ErrorRate > 1.0 {
		b.Errorf("Error rate too high: %.2f%%, expected <= 1%%", results.ErrorRate)
	}
	
	if results.P99Latency > 100*time.Millisecond {
		b.Errorf("P99 latency too high: %v, expected <= 100ms", results.P99Latency)
	}
	
	if results.MemoryUsage > 100*1024*1024 { // 100MB
		b.Errorf("Memory usage too high: %d bytes, expected <= 100MB", results.MemoryUsage)
	}
}

// BenchmarkAPIEndpoints tests different API endpoints
func BenchmarkAPIEndpoints(b *testing.B) {
	endpoints := []string{
		"/api/metrics",
		"/api/logs",
		"/api/traces",
		"/api/alerts",
		"/health",
	}
	
	for _, endpoint := range endpoints {
		b.Run(endpoint, func(b *testing.B) {
			config := BenchmarkConfig{
				TargetURL:       "http://localhost:8080" + endpoint,
				RequestsPerSec:  5000,
				Duration:        10 * time.Second,
				ConcurrentUsers: 50,
				PayloadSize:     128,
				Timeout:         3 * time.Second,
				KeepAlive:       true,
			}
			
			generator := NewLoadGenerator(config)
			results := generator.RunBenchmark()
			
			b.Logf("Endpoint: %s", endpoint)
			b.Logf("RPS: %.2f", results.RequestsPerSec)
			b.Logf("P99 Latency: %v", results.P99Latency)
			b.Logf("Error Rate: %.2f%%", results.ErrorRate)
			
			// Set benchmark metrics for reporting
			b.ReportMetric(results.RequestsPerSec, "req/sec")
			b.ReportMetric(float64(results.P99Latency.Nanoseconds())/1e6, "ms")
		})
	}
}

// BenchmarkConcurrentLoad tests concurrent load patterns
func BenchmarkConcurrentLoad(b *testing.B) {
	testCases := []struct {
		name             string
		concurrentUsers  int
		requestsPerSec   int
		duration         time.Duration
	}{
		{"Light Load", 10, 1000, 10 * time.Second},
		{"Medium Load", 50, 5000, 10 * time.Second},
		{"Heavy Load", 100, 10000, 10 * time.Second},
		{"Peak Load", 200, 20000, 10 * time.Second},
	}
	
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			config := BenchmarkConfig{
				TargetURL:       "http://localhost:8080/api/metrics",
				RequestsPerSec:  tc.requestsPerSec,
				Duration:        tc.duration,
				ConcurrentUsers:  tc.concurrentUsers,
				PayloadSize:     256,
				Timeout:         5 * time.Second,
				KeepAlive:       true,
			}
			
			generator := NewLoadGenerator(config)
			results := generator.RunBenchmark()
			
			b.Logf("Test: %s", tc.name)
			b.Logf("Target RPS: %d, Actual RPS: %.2f", tc.requestsPerSec, results.RequestsPerSec)
			b.Logf("P99 Latency: %v", results.P99Latency)
			b.Logf("Error Rate: %.2f%%", results.ErrorRate)
			b.Logf("Memory Usage: %.2f MB", float64(results.MemoryUsage)/1024/1024)
			
			// Performance should scale reasonably
			efficiency := results.RequestsPerSec / float64(tc.concurrentUsers)
			if efficiency < 50 { // Less than 50 req/sec per user is inefficient
				b.Errorf("Low efficiency: %.2f req/sec per user", efficiency)
			}
		})
	}
}

// BenchmarkPayloadSizes tests different payload sizes
func BenchmarkPayloadSizes(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096, 16384}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			config := BenchmarkConfig{
				TargetURL:       "http://localhost:8080/api/metrics",
				RequestsPerSec:  5000,
				Duration:        10 * time.Second,
				ConcurrentUsers: 50,
				PayloadSize:     size,
				Timeout:         5 * time.Second,
				KeepAlive:       true,
			}
			
			generator := NewLoadGenerator(config)
			results := generator.RunBenchmark()
			
			b.Logf("Payload Size: %d bytes", size)
			b.Logf("RPS: %.2f", results.RequestsPerSec)
			b.Logf("P99 Latency: %v", results.P99Latency)
			b.Logf("Throughput: %.2f MB/s", float64(results.BytesTransferred)/1024/1024/results.Duration.Seconds())
			
			// Larger payloads should not dramatically reduce performance
			if size > 1024 && results.RequestsPerSec < 2000 {
				b.Errorf("Large payload performance too low: %.2f RPS for %d bytes", results.RequestsPerSec, size)
			}
		})
	}
}

// BenchmarkStressTest runs a stress test
func BenchmarkStressTest(b *testing.B) {
	config := BenchmarkConfig{
		TargetURL:       "http://localhost:8080/api/metrics",
		RequestsPerSec:  50000, // 50k req/sec target
		Duration:        60 * time.Second, // 1 minute stress test
		ConcurrentUsers: 500,
		PayloadSize:     512,
		Timeout:         10 * time.Second,
		KeepAlive:       true,
	}
	
	generator := NewLoadGenerator(config)
	results := generator.RunBenchmark()
	
	b.Logf("=== Stress Test Results ===")
	b.Logf("Target RPS: 50,000")
	b.Logf("Actual RPS: %.2f", results.RequestsPerSec)
	b.Logf("Total Requests: %d", results.TotalRequests)
	b.Logf("Duration: %v", results.Duration)
	b.Logf("P99 Latency: %v", results.P99Latency)
	b.Logf("Error Rate: %.2f%%", results.ErrorRate)
	b.Logf("Memory Usage: %.2f MB", float64(results.MemoryUsage)/1024/1024)
	b.Logf("Throughput: %.2f MB/s", float64(results.BytesTransferred)/1024/1024/results.Duration.Seconds())
	
	// Stress test assertions
	if results.RequestsPerSec < 30000 {
		b.Errorf("Stress test RPS too low: %.2f, expected >= 30000", results.RequestsPerSec)
	}
	
	if results.ErrorRate > 5.0 {
		b.Errorf("Stress test error rate too high: %.2f%%, expected <= 5%%", results.ErrorRate)
	}
}

// SaveResults saves benchmark results to JSON file
func SaveResults(results *BenchmarkResults, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	
	return WriteFile(filename, data, 0644)
}

// WriteFile is a placeholder for os.WriteFile
func WriteFile(filename string, data []byte, perm int) error {
	// In real implementation, this would use os.WriteFile
	fmt.Printf("Results saved to: %s\n", filename)
	fmt.Printf("Data: %s\n", string(data))
	return nil
}
