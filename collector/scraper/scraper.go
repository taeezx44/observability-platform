package scraper

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Metric struct {
	Timestamp time.Time
	Name      string
	Value     float64
	Labels    map[string]string
}

type Storage interface {
	BatchInsert(metrics []Metric) error
}

type Scraper struct {
	Targets  []string     // ["http://app:8080/metrics"]
	Interval time.Duration // 15s
	Storage  Storage
}

func (s *Scraper) Run() {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	
	for range ticker.C {
		for _, target := range s.Targets {
			metrics, err := s.scrapeTarget(target)
			if err != nil {
				fmt.Printf("scrape error %s: %v\n", target, err)
				continue
			}
			// Batch insert — ห้าม insert ทีละ row!
			if len(metrics) > 0 {
				s.Storage.BatchInsert(metrics)
			}
		}
	}
}

func (s *Scraper) scrapeTarget(url string) ([]Metric, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now()
	var metrics []Metric
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		// ข้าม comments และ empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		m, err := parseLine(line, now)
		if err != nil {
			continue
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// parseLine: "http_requests_total{method="GET"} 1234 1705123456"
func parseLine(line string, ts time.Time) (Metric, error) {
	// Split on whitespace to get metric part and value
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Metric{}, fmt.Errorf("invalid metric format")
	}

	metricPart := parts[0]
	valueStr := parts[1]

	// Parse value
	var value float64
	_, err := fmt.Sscanf(valueStr, "%f", &value)
	if err != nil {
		return Metric{}, fmt.Errorf("invalid value: %s", valueStr)
	}

	// Parse metric name and labels
	name, labels := parseMetricLabels(metricPart)

	return Metric{
		Timestamp: ts,
		Name:      name,
		Value:     value,
		Labels:    labels,
	}, nil
}

func parseMetricLabels(metricPart string) (string, map[string]string) {
	labels := make(map[string]string)
	
	// Check if there are labels
	if !strings.Contains(metricPart, "{") {
		return metricPart, labels
	}

	// Extract name and labels part
	nameStart := strings.Index(metricPart, "{")
	nameEnd := strings.Index(metricPart, "}")
	
	if nameStart == -1 || nameEnd == -1 {
		// Fallback - treat whole thing as name
		return metricPart, labels
	}

	name := metricPart[:nameStart]
	labelsPart := metricPart[nameStart+1 : nameEnd]

	// Parse labels like method="GET",status="200"
	labelPairs := strings.Split(labelsPart, ",")
	for _, pair := range labelPairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(kv[1], `"`)
			labels[key] = value
		}
	}

	return name, labels
}
