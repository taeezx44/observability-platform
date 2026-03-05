package tracer

import (
	"time"
)

type Span struct {
	TraceID   string            // ผูก request ทั้งหมด
	SpanID    string            // unique ต่อ operation
	ParentID  string            // "" = root span
	Service   string
	Operation string
	StartTime time.Time
	EndTime   time.Time
	Status    string            // "OK" | "ERROR"
	Tags      map[string]string // {"db.type":"postgres"}
	Logs      []SpanLog
}

func (s *Span) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

func (s *Span) DurationMs() float64 {
	return float64(s.Duration().Nanoseconds()) / 1e6
}

func (s *Span) IsRoot() bool {
	return s.ParentID == ""
}

func (s *Span) IsSlow(threshold time.Duration) bool {
	return s.Duration() > threshold
}

type SpanLog struct {
	Timestamp time.Time
	Fields    map[string]interface{}
}

type SpanWithDepth struct {
	Span  Span
	Depth int
}

// BuildTree จัด spans เป็น tree structure
func BuildTree(spans []Span) map[string][]Span {
	tree := make(map[string][]Span)
	for _, s := range spans {
		tree[s.ParentID] = append(tree[s.ParentID], s)
	}
	return tree // tree[""] = root spans
}

// FlattenTree converts tree to flat list with depth info for waterfall
func FlattenTree(tree map[string][]Span, parentID string, depth int) []SpanWithDepth {
	var result []SpanWithDepth
	
	for _, span := range tree[parentID] {
		result = append(result, SpanWithDepth{
			Span:  span,
			Depth: depth,
		})
		
		// Recursively add children
		children := FlattenTree(tree, span.SpanID, depth+1)
		result = append(result, children...)
	}
	
	return result
}
