package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yourname/observability/collector/storage"
	"gopkg.in/yaml.v3"
)

// ── Rule types ────────────────────────────────────────────────────────────────

type AlertRule struct {
	Name       string        `yaml:"name"`
	MetricName string        `yaml:"metric"`
	Condition  string        `yaml:"condition"` // ">" | "<" | "==" | "!="
	Threshold  float64       `yaml:"threshold"`
	For        time.Duration `yaml:"for"`
	Severity   string        `yaml:"severity"` // critical | warning | info
	SlackURL   string        `yaml:"slack_url"`
}

type RulesConfig struct {
	Rules []AlertRule `yaml:"rules"`
}

// ── Alert Engine ──────────────────────────────────────────────────────────────

type Engine struct {
	rules   []AlertRule
	store   *storage.ClickHouseStorage
	firing  map[string]time.Time // rule name → when it first fired
	slackURL string
}

func NewEngine(store *storage.ClickHouseStorage, rulesPath, slackURL string) (*Engine, error) {
	e := &Engine{
		store:    store,
		firing:   make(map[string]time.Time),
		slackURL: slackURL,
	}

	if err := e.loadRules(rulesPath); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engine) loadRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read rules: %w", err)
	}
	var cfg RulesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse rules: %w", err)
	}
	e.rules = cfg.Rules
	log.Printf("Loaded %d alert rules", len(e.rules))
	return nil
}

func (e *Engine) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Printf("Alert engine running, evaluating %d rules every 30s", len(e.rules))

	for range ticker.C {
		for _, rule := range e.rules {
			e.evaluate(rule)
		}
	}
}

func (e *Engine) evaluate(rule AlertRule) {
	val, err := e.latestValue(rule.MetricName)
	if err != nil {
		log.Printf("Error fetching metric %s: %v", rule.MetricName, err)
		return
	}

	isFiring := e.check(val, rule.Condition, rule.Threshold)

	if isFiring {
		first, seen := e.firing[rule.Name]
		if !seen {
			// First time firing — record time
			e.firing[rule.Name] = time.Now()
			log.Printf("[PENDING] %s: %.2f %s %.2f", rule.Name, val, rule.Condition, rule.Threshold)
			return
		}
		// Fire alert if it's been firing longer than 'for' duration
		if time.Since(first) >= rule.For {
			log.Printf("[FIRING] %s: %.2f %s %.2f (severity: %s)", rule.Name, val, rule.Condition, rule.Threshold, rule.Severity)
			e.notify(rule, val)
		}
	} else {
		if _, wasFiring := e.firing[rule.Name]; wasFiring {
			log.Printf("[RESOLVED] %s", rule.Name)
		}
		delete(e.firing, rule.Name)
	}
}

func (e *Engine) latestValue(metricName string) (float64, error) {
	metrics, err := e.store.GetMetrics(storage.MetricsQuery{
		Name:  metricName,
		From:  time.Now().Add(-2 * time.Minute),
		To:    time.Now(),
		Limit: 1,
	})
	if err != nil {
		return 0, err
	}
	if len(metrics) == 0 {
		return 0, fmt.Errorf("no data for metric %s", metricName)
	}
	return metrics[0].Value, nil
}

func (e *Engine) check(val float64, cond string, threshold float64) bool {
	switch cond {
	case ">":
		return val > threshold
	case "<":
		return val < threshold
	case "==":
		return val == threshold
	case "!=":
		return val != threshold
	case ">=":
		return val >= threshold
	case "<=":
		return val <= threshold
	}
	return false
}

func (e *Engine) notify(rule AlertRule, val float64) {
	emoji := "🚨"
	if rule.Severity == "warning" {
		emoji = "⚠️"
	} else if rule.Severity == "info" {
		emoji = "ℹ️"
	}

	msg := fmt.Sprintf("%s *[%s] %s*\nMetric `%s` is `%.2f` (threshold: `%s %.2f`)",
		emoji, rule.Severity, rule.Name,
		rule.MetricName, val, rule.Condition, rule.Threshold,
	)

	// Use rule-specific slack URL or fallback to global
	slackURL := rule.SlackURL
	if slackURL == "" {
		slackURL = e.slackURL
	}
	if slackURL == "" {
		log.Printf("No Slack URL configured for rule %s", rule.Name)
		return
	}

	payload, _ := json.Marshal(map[string]string{"text": msg})
	resp, err := http.Post(slackURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Slack notification sent for %s (status: %d)", rule.Name, resp.StatusCode)
}
