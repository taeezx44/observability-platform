package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/yourname/observability/collector/logger"
	"github.com/yourname/observability/collector/scraper"
	"github.com/yourname/observability/collector/tracer"
)

type ClickHouseStorage struct {
	conn clickhouse.Conn
}

func NewClickHouseStorage(dsn string) (*ClickHouseStorage, error) {
	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}
	return &ClickHouseStorage{conn: conn}, nil
}

// BatchInsert satisfies scraper.Storage interface
func (s *ClickHouseStorage) BatchInsert(metrics []scraper.Metric) error {
	return s.BatchInsertMetrics(metrics)
}

func (s *ClickHouseStorage) BatchInsertMetrics(metrics []scraper.Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	batch, err := s.conn.PrepareBatch(context.Background(),
		`INSERT INTO metrics (timestamp, name, value, labels)`)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	for _, m := range metrics {
		if err := batch.Append(m.Timestamp, m.Name, m.Value, m.Labels); err != nil {
			return fmt.Errorf("append metric: %w", err)
		}
	}
	return batch.Send()
}

func (s *ClickHouseStorage) GetMetrics(query MetricsQuery) ([]scraper.Metric, error) {
	q := `SELECT timestamp, name, value, labels FROM metrics WHERE timestamp >= ? AND timestamp <= ?`
	args := []interface{}{query.From, query.To}
	if query.Name != "" {
		q += " AND name = ?"
		args = append(args, query.Name)
	}
	for k, v := range query.Labels {
		q += " AND labels[?] = ?"
		args = append(args, k, v)
	}
	q += " ORDER BY timestamp DESC"
	if query.Limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", query.Limit)
	}

	rows, err := s.conn.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var metrics []scraper.Metric
	for rows.Next() {
		var m scraper.Metric
		if err := rows.Scan(&m.Timestamp, &m.Name, &m.Value, &m.Labels); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *ClickHouseStorage) QueryDistinctNames() ([]string, error) {
	rows, err := s.conn.Query(context.Background(),
		`SELECT DISTINCT name FROM metrics WHERE timestamp >= now() - INTERVAL 24 HOUR ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// ── Logs ─────────────────────────────────────────────────────────────────────

func (s *ClickHouseStorage) BatchInsertLogs(logs []logger.LogEntry) error {
	if len(logs) == 0 {
		return nil
	}
	batch, err := s.conn.PrepareBatch(context.Background(),
		`INSERT INTO logs (timestamp, level, service, message, fields)`)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	for _, l := range logs {
		if err := batch.Append(l.Timestamp, l.Level, l.Service, l.Message, l.Fields); err != nil {
			return fmt.Errorf("append log: %w", err)
		}
	}
	return batch.Send()
}

func (s *ClickHouseStorage) GetLogs(query LogsQuery) ([]logger.LogEntry, error) {
	q := `SELECT timestamp, level, service, message, fields FROM logs WHERE timestamp >= ? AND timestamp <= ?`
	args := []interface{}{query.From, query.To}
	if query.Level != "" {
		q += " AND level = ?"
		args = append(args, query.Level)
	}
	if query.Service != "" {
		q += " AND service = ?"
		args = append(args, query.Service)
	}
	if query.Search != "" {
		q += " AND message ILIKE ?"
		args = append(args, "%"+query.Search+"%")
	}
	q += " ORDER BY timestamp DESC"
	if query.Limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", query.Limit)
	}

	rows, err := s.conn.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []logger.LogEntry
	for rows.Next() {
		var l logger.LogEntry
		if err := rows.Scan(&l.Timestamp, &l.Level, &l.Service, &l.Message, &l.Fields); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// ── Traces ───────────────────────────────────────────────────────────────────

func (s *ClickHouseStorage) BatchInsertSpans(spans []tracer.Span) error {
	if len(spans) == 0 {
		return nil
	}
	batch, err := s.conn.PrepareBatch(context.Background(),
		`INSERT INTO spans (trace_id, span_id, parent_id, service, operation, start_time, end_time, status, tags)`)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	for _, sp := range spans {
		if err := batch.Append(
			sp.TraceID, sp.SpanID, sp.ParentID,
			sp.Service, sp.Operation,
			sp.StartTime, sp.EndTime,
			sp.Status, sp.Tags,
		); err != nil {
			return fmt.Errorf("append span: %w", err)
		}
	}
	return batch.Send()
}

func (s *ClickHouseStorage) GetTrace(traceID string) ([]tracer.Span, error) {
	rows, err := s.conn.Query(context.Background(),
		`SELECT trace_id, span_id, parent_id, service, operation, start_time, end_time, status, tags
		 FROM spans WHERE trace_id = ? ORDER BY start_time`, traceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var spans []tracer.Span
	for rows.Next() {
		var sp tracer.Span
		if err := rows.Scan(&sp.TraceID, &sp.SpanID, &sp.ParentID,
			&sp.Service, &sp.Operation, &sp.StartTime, &sp.EndTime, &sp.Status, &sp.Tags); err != nil {
			return nil, err
		}
		spans = append(spans, sp)
	}
	return spans, nil
}

func (s *ClickHouseStorage) GetSlowTraces(minMs float64, limit int) ([]TracesSummary, error) {
	rows, err := s.conn.Query(context.Background(),
		`SELECT trace_id, service, operation, count() as span_count,
		        dateDiff('ms', min(start_time), max(end_time)) as duration_ms
		 FROM spans WHERE start_time >= now() - INTERVAL 1 HOUR
		 GROUP BY trace_id, service, operation
		 HAVING duration_ms > ?
		 ORDER BY duration_ms DESC LIMIT ?`, minMs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []TracesSummary
	for rows.Next() {
		var t TracesSummary
		if err := rows.Scan(&t.TraceID, &t.Service, &t.Operation, &t.SpanCount, &t.DurationMs); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *ClickHouseStorage) Close() error { return s.conn.Close() }

// ── Query types ───────────────────────────────────────────────────────────────

type MetricsQuery struct {
	Name   string
	Labels map[string]string
	From   time.Time
	To     time.Time
	Limit  int
}

type LogsQuery struct {
	Level   string
	Service string
	Search  string
	From    time.Time
	To      time.Time
	Limit   int
}

type TracesSummary struct {
	TraceID    string  `json:"trace_id"`
	Service    string  `json:"service"`
	Operation  string  `json:"operation"`
	SpanCount  uint64  `json:"span_count"`
	DurationMs float64 `json:"duration_ms"`
}
