package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// MockLogParser represents a log parser for testing
type MockLogParser struct {
	entries []LogEntry
}

func NewMockLogParser() *MockLogParser {
	return &MockLogParser{
		entries: []LogEntry{},
	}
}

func (p *MockLogParser) ParseJSONLog(logLine string) (*LogEntry, error) {
	var entry LogEntry
	if err := json.Unmarshal([]byte(logLine), &entry); err != nil {
		return nil, err
	}

	p.entries = append(p.entries, entry)
	return &entry, nil
}

func (p *MockLogParser) ParsePlainTextLog(logLine string) *LogEntry {
	// Simple plain text parser for testing
	// Format: [TIMESTAMP] LEVEL SERVICE MESSAGE
	parts := strings.SplitN(logLine, " ", 4)
	if len(parts) < 4 {
		return nil
	}

	// Remove brackets from timestamp
	timestampStr := strings.Trim(parts[0], "[]")
	timestamp, err := time.Parse("2006-01-02T15:04:05", timestampStr)
	if err != nil {
		timestamp = time.Now()
	}

	entry := LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Service:   parts[2],
		Message:   parts[3],
	}

	p.entries = append(p.entries, entry)
	return &entry
}

func (p *MockLogParser) GetEntries() []LogEntry {
	return p.entries
}

func (p *MockLogParser) FilterByLevel(level string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range p.entries {
		if entry.Level == level {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (p *MockLogParser) FilterByService(service string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range p.entries {
		if entry.Service == service {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// TestJSONLogParsing tests parsing JSON log entries
func TestJSONLogParsing(t *testing.T) {
	parser := NewMockLogParser()

	jsonLog := `{
		"timestamp": "2024-01-15T10:30:00Z",
		"level": "INFO",
		"service": "api",
		"message": "Request processed successfully",
		"labels": {
			"request_id": "req-123",
			"user_id": "user-456"
		}
	}`

	entry, err := parser.ParseJSONLog(jsonLog)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", entry.Level)
	}

	if entry.Service != "api" {
		t.Errorf("Expected service 'api', got '%s'", entry.Service)
	}

	if entry.Message != "Request processed successfully" {
		t.Errorf("Expected message 'Request processed successfully', got '%s'", entry.Message)
	}

	if entry.Fields["request_id"] != "req-123" {
		t.Errorf("Expected request_id 'req-123', got '%s'", entry.Fields["request_id"])
	}
}

// TestPlainTextLogParsing tests parsing plain text log entries
func TestPlainTextLogParsing(t *testing.T) {
	parser := NewMockLogParser()

	plainLog := "[2024-01-15T10:30:00] INFO api Request processed successfully"

	entry := parser.ParsePlainTextLog(plainLog)
	if entry == nil {
		t.Fatal("Expected non-nil entry")
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", entry.Level)
	}

	if entry.Service != "api" {
		t.Errorf("Expected service 'api', got '%s'", entry.Service)
	}

	if entry.Message != "Request processed successfully" {
		t.Errorf("Expected message 'Request processed successfully', got '%s'", entry.Message)
	}
}

// TestInvalidJSONLog tests parsing invalid JSON
func TestInvalidJSONLog(t *testing.T) {
	parser := NewMockLogParser()

	invalidJSON := `{"timestamp": "2024-01-15T10:30:00Z", "level": "INFO"`

	_, err := parser.ParseJSONLog(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestInvalidPlainTextLog tests parsing invalid plain text
func TestInvalidPlainTextLog(t *testing.T) {
	parser := NewMockLogParser()

	invalidLog := "invalid log format"

	entry := parser.ParsePlainTextLog(invalidLog)
	if entry != nil {
		t.Error("Expected nil entry for invalid format, got non-nil")
	}
}

// TestMultipleLogEntries tests parsing multiple log entries
func TestMultipleLogEntries(t *testing.T) {
	parser := NewMockLogParser()

	// Add JSON logs
	jsonLogs := []string{
		`{"timestamp": "2024-01-15T10:30:00Z", "level": "INFO", "service": "api", "message": "Request 1"}`,
		`{"timestamp": "2024-01-15T10:31:00Z", "level": "ERROR", "service": "api", "message": "Request 2"}`,
		`{"timestamp": "2024-01-15T10:32:00Z", "level": "WARN", "service": "database", "message": "Connection slow"}`,
	}

	for _, log := range jsonLogs {
		_, err := parser.ParseJSONLog(log)
		if err != nil {
			t.Errorf("Failed to parse JSON log: %v", err)
		}
	}

	// Add plain text logs
	plainLogs := []string{
		"[2024-01-15T10:33:00] INFO worker Task completed",
		"[2024-01-15T10:34:00] ERROR worker Task failed",
	}

	for _, log := range plainLogs {
		entry := parser.ParsePlainTextLog(log)
		if entry == nil {
			t.Error("Failed to parse plain text log")
		}
	}

	entries := parser.GetEntries()
	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}
}

// TestLogFiltering tests filtering logs by level and service
func TestLogFiltering(t *testing.T) {
	parser := NewMockLogParser()

	// Add test logs
	logs := []string{
		`{"timestamp": "2024-01-15T10:30:00Z", "level": "INFO", "service": "api", "message": "Info message"}`,
		`{"timestamp": "2024-01-15T10:31:00Z", "level": "ERROR", "service": "api", "message": "Error message"}`,
		`{"timestamp": "2024-01-15T10:32:00Z", "level": "INFO", "service": "database", "message": "DB info"}`,
		`{"timestamp": "2024-01-15T10:33:00Z", "level": "ERROR", "service": "database", "message": "DB error"}`,
	}

	for _, log := range logs {
		_, err := parser.ParseJSONLog(log)
		if err != nil {
			t.Errorf("Failed to parse log: %v", err)
		}
	}

	// Test level filtering
	infoLogs := parser.FilterByLevel("INFO")
	if len(infoLogs) != 2 {
		t.Errorf("Expected 2 INFO logs, got %d", len(infoLogs))
	}

	errorLogs := parser.FilterByLevel("ERROR")
	if len(errorLogs) != 2 {
		t.Errorf("Expected 2 ERROR logs, got %d", len(errorLogs))
	}

	// Test service filtering
	apiLogs := parser.FilterByService("api")
	if len(apiLogs) != 2 {
		t.Errorf("Expected 2 api logs, got %d", len(apiLogs))
	}

	dbLogs := parser.FilterByService("database")
	if len(dbLogs) != 2 {
		t.Errorf("Expected 2 database logs, got %d", len(dbLogs))
	}
}

// TestLogTimestamps tests timestamp parsing and validation
func TestLogTimestamps(t *testing.T) {
	parser := NewMockLogParser()

	// Test valid timestamp
	validLog := `{"timestamp": "2024-01-15T10:30:00Z", "level": "INFO", "service": "api", "message": "Test"}`
	entry, err := parser.ParseJSONLog(validLog)
	if err != nil {
		t.Errorf("Failed to parse valid timestamp: %v", err)
	}

	if entry.Timestamp.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", entry.Timestamp.Year())
	}

	// Test plain text timestamp parsing
	plainLog := "[2024-01-15T10:30:00] INFO api Test"
	entry = parser.ParsePlainTextLog(plainLog)
	if entry == nil {
		t.Fatal("Failed to parse plain text log")
	}

	if entry.Timestamp.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", entry.Timestamp.Year())
	}
}

// BenchmarkJSONLogParsing benchmarks JSON log parsing
func BenchmarkJSONLogParsing(b *testing.B) {
	parser := NewMockLogParser()

	jsonLog := `{
		"timestamp": "2024-01-15T10:30:00Z",
		"level": "INFO",
		"service": "api",
		"message": "Request processed successfully",
		"labels": {
			"request_id": "req-123",
			"user_id": "user-456"
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseJSONLog(jsonLog)
		if err != nil {
			b.Errorf("Failed to parse JSON log: %v", err)
		}
	}
}

// BenchmarkPlainTextLogParsing benchmarks plain text log parsing
func BenchmarkPlainTextLogParsing(b *testing.B) {
	parser := NewMockLogParser()

	plainLog := "[2024-01-15T10:30:00] INFO api Request processed successfully"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entry := parser.ParsePlainTextLog(plainLog)
		if entry == nil {
			b.Error("Failed to parse plain text log")
		}
	}
}

// BenchmarkLogFiltering benchmarks log filtering
func BenchmarkLogFiltering(b *testing.B) {
	parser := NewMockLogParser()

	// Add many entries for realistic benchmarking
	for i := 0; i < 1000; i++ {
		log := fmt.Sprintf(`{"timestamp": "2024-01-15T10:30:00Z", "level": "%s", "service": "api", "message": "Message %d"}`,
			map[bool]string{true: "INFO", false: "ERROR"}[i%2 == 0], i)
		_, err := parser.ParseJSONLog(log)
		if err != nil {
			b.Errorf("Failed to parse log: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.FilterByLevel("INFO")
	}
}
