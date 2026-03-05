package alerting

import (
	"testing"
	"time"
)

// MockAlert represents a test alert
type MockAlert struct {
	ID        string    `json:"id"`
	Severity  string    `json:"severity"`
	Service   string    `json:"service"`
	Metric    string    `json:"metric"`
	Threshold float64   `json:"threshold"`
	Current   float64   `json:"current"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AlertEngine represents a simple alerting engine for testing
type AlertEngine struct {
	alerts []MockAlert
}

// NewAlertEngine creates a new alert engine
func NewAlertEngine() *AlertEngine {
	return &AlertEngine{
		alerts: make([]MockAlert, 0),
	}
}

// CheckThreshold checks if a metric exceeds threshold and creates alert
func (ae *AlertEngine) CheckThreshold(service, metric string, current, threshold float64) *MockAlert {
	if current > threshold {
		alert := MockAlert{
			ID:        generateAlertID(),
			Severity:  "warning",
			Service:   service,
			Metric:    metric,
			Threshold: threshold,
			Current:   current,
			Status:    "active",
			Message:   metric + " is above threshold",
			Timestamp: time.Now(),
		}
		ae.alerts = append(ae.alerts, alert)
		return &alert
	}
	return nil
}

// GetActiveAlerts returns all active alerts
func (ae *AlertEngine) GetActiveAlerts() []MockAlert {
	var active []MockAlert
	for _, alert := range ae.alerts {
		if alert.Status == "active" {
			active = append(active, alert)
		}
	}
	return active
}

// ResolveAlert resolves an alert by ID
func (ae *AlertEngine) ResolveAlert(alertID string) bool {
	for i, alert := range ae.alerts {
		if alert.ID == alertID {
			ae.alerts[i].Status = "resolved"
			ae.alerts[i].Message = alert.Metric + " has normalized"
			return true
		}
	}
	return false
}

// GetAlertStats returns alert statistics
func (ae *AlertEngine) GetAlertStats() map[string]int {
	stats := map[string]int{
		"total":    len(ae.alerts),
		"active":   0,
		"resolved": 0,
	}
	
	for _, alert := range ae.alerts {
		if alert.Status == "active" {
			stats["active"]++
		} else if alert.Status == "resolved" {
			stats["resolved"]++
		}
	}
	
	return stats
}

// generateAlertID generates a unique alert ID
func generateAlertID() string {
	return "alert-" + time.Now().Format("20060102-150405")
}

// TestAlertEngineCreation tests creating a new alert engine
func TestAlertEngineCreation(t *testing.T) {
	ae := NewAlertEngine()
	
	if ae == nil {
		t.Fatal("Expected non-nil alert engine")
	}
	
	if len(ae.alerts) != 0 {
		t.Errorf("Expected empty alerts slice, got %d alerts", len(ae.alerts))
	}
}

// TestThresholdCheck tests threshold checking functionality
func TestThresholdCheck(t *testing.T) {
	ae := NewAlertEngine()
	
	// Test normal value (should not create alert)
	alert := ae.CheckThreshold("api", "cpu_usage", 50.0, 80.0)
	if alert != nil {
		t.Errorf("Expected nil alert for normal value, got alert: %+v", alert)
	}
	
	// Test high value (should create alert)
	alert = ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	if alert == nil {
		t.Error("Expected alert for high value, got nil")
	} else {
		if alert.Service != "api" {
			t.Errorf("Expected service 'api', got '%s'", alert.Service)
		}
		if alert.Metric != "cpu_usage" {
			t.Errorf("Expected metric 'cpu_usage', got '%s'", alert.Metric)
		}
		if alert.Current != 85.0 {
			t.Errorf("Expected current value 85.0, got %f", alert.Current)
		}
		if alert.Threshold != 80.0 {
			t.Errorf("Expected threshold 80.0, got %f", alert.Threshold)
		}
		if alert.Status != "active" {
			t.Errorf("Expected status 'active', got '%s'", alert.Status)
		}
	}
}

// TestMultipleAlerts tests handling multiple alerts
func TestMultipleAlerts(t *testing.T) {
	ae := NewAlertEngine()
	
	// Create multiple alerts
	ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	ae.CheckThreshold("database", "memory_usage", 90.0, 85.0)
	ae.CheckThreshold("kafka", "disk_usage", 95.0, 90.0)
	
	activeAlerts := ae.GetActiveAlerts()
	if len(activeAlerts) != 3 {
		t.Errorf("Expected 3 active alerts, got %d", len(activeAlerts))
	}
	
	// Check alert details
	for i, alert := range activeAlerts {
		if alert.Status != "active" {
			t.Errorf("Alert %d: Expected status 'active', got '%s'", i, alert.Status)
		}
	}
}

// TestAlertResolution tests alert resolution
func TestAlertResolution(t *testing.T) {
	ae := NewAlertEngine()
	
	// Create an alert
	alert := ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	if alert == nil {
		t.Fatal("Expected alert to be created")
	}
	
	alertID := alert.ID
	
	// Verify alert is active
	activeAlerts := ae.GetActiveAlerts()
	if len(activeAlerts) != 1 {
		t.Errorf("Expected 1 active alert, got %d", len(activeAlerts))
	}
	
	// Resolve the alert
	resolved := ae.ResolveAlert(alertID)
	if !resolved {
		t.Error("Expected alert to be resolved")
	}
	
	// Verify alert is no longer active
	activeAlerts = ae.GetActiveAlerts()
	if len(activeAlerts) != 0 {
		t.Errorf("Expected 0 active alerts after resolution, got %d", len(activeAlerts))
	}
}

// TestAlertStats tests alert statistics
func TestAlertStats(t *testing.T) {
	ae := NewAlertEngine()
	
	// Initially should have no alerts
	stats := ae.GetAlertStats()
	if stats["total"] != 0 {
		t.Errorf("Expected 0 total alerts, got %d", stats["total"])
	}
	if stats["active"] != 0 {
		t.Errorf("Expected 0 active alerts, got %d", stats["active"])
	}
	if stats["resolved"] != 0 {
		t.Errorf("Expected 0 resolved alerts, got %d", stats["resolved"])
	}
	
	// Create alerts
	ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	ae.CheckThreshold("database", "memory_usage", 90.0, 85.0)
	
	stats = ae.GetAlertStats()
	if stats["total"] != 2 {
		t.Errorf("Expected 2 total alerts, got %d", stats["total"])
	}
	if stats["active"] != 2 {
		t.Errorf("Expected 2 active alerts, got %d", stats["active"])
	}
	if stats["resolved"] != 0 {
		t.Errorf("Expected 0 resolved alerts, got %d", stats["resolved"])
	}
	
	// Resolve one alert
	alert := ae.GetActiveAlerts()[0]
	ae.ResolveAlert(alert.ID)
	
	stats = ae.GetAlertStats()
	if stats["total"] != 2 {
		t.Errorf("Expected 2 total alerts, got %d", stats["total"])
	}
	if stats["active"] != 1 {
		t.Errorf("Expected 1 active alert, got %d", stats["active"])
	}
	if stats["resolved"] != 1 {
		t.Errorf("Expected 1 resolved alert, got %d", stats["resolved"])
	}
}

// TestAlertIDGeneration tests alert ID generation
func TestAlertIDGeneration(t *testing.T) {
	ae := NewAlertEngine()
	
	// Create alerts quickly to test ID uniqueness
	alert1 := ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	time.Sleep(1 * time.Second) // Ensure different timestamps
	alert2 := ae.CheckThreshold("database", "memory_usage", 90.0, 85.0)
	
	if alert1.ID == alert2.ID {
		t.Error("Expected different alert IDs, got same ID")
	}
	
	if len(alert1.ID) < 10 {
		t.Errorf("Alert ID too short: %s", alert1.ID)
	}
	
	if alert1.ID[:6] != "alert-" {
		t.Errorf("Alert ID should start with 'alert-', got: %s", alert1.ID[:6])
	}
}

// TestAlertSeverity tests alert severity levels
func TestAlertSeverity(t *testing.T) {
	ae := NewAlertEngine()
	
	// Test different threshold levels to simulate different severities
	warningAlert := ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	if warningAlert == nil {
		t.Fatal("Expected warning alert to be created")
	}
	
	if warningAlert.Severity != "warning" {
		t.Errorf("Expected severity 'warning', got '%s'", warningAlert.Severity)
	}
	
	// In a real implementation, you might have different severity levels
	// For this test, we'll just verify the basic structure
	if warningAlert.Message == "" {
		t.Error("Expected non-empty alert message")
	}
}

// BenchmarkAlertCreation benchmarks alert creation
func BenchmarkAlertCreation(b *testing.B) {
	ae := NewAlertEngine()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	}
}

// BenchmarkAlertRetrieval benchmarks alert retrieval
func BenchmarkAlertRetrieval(b *testing.B) {
	ae := NewAlertEngine()
	
	// Create some alerts first
	for i := 0; i < 100; i++ {
		ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ae.GetActiveAlerts()
	}
}

// BenchmarkAlertStats benchmarks alert statistics calculation
func BenchmarkAlertStats(b *testing.B) {
	ae := NewAlertEngine()
	
	// Create some alerts first
	for i := 0; i < 100; i++ {
		ae.CheckThreshold("api", "cpu_usage", 85.0, 80.0)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ae.GetAlertStats()
	}
}
