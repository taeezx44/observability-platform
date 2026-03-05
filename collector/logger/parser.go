package logger

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"time"
)

type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
}

// รองรับ 2 format: JSON และ plaintext
var plainRegex = regexp.MustCompile(
	`(\d{4}-\d{2}-\d{2}T[\d:]+Z?)\s+(ERROR|WARN|INFO|DEBUG)\s+(.*)`,
)

func ParseLine(line, service string) LogEntry {
	// ลอง parse JSON ก่อน
	var entry LogEntry
	if json.Unmarshal([]byte(line), &entry) == nil {
		entry.Service = service
		return entry
	}

	// fallback: plain text regex
	if m := plainRegex.FindStringSubmatch(line); m != nil {
		ts, _ := time.Parse(time.RFC3339, m[1])
		return LogEntry{
			Timestamp: ts,
			Level:     m[2],
			Service:   service,
			Message:   m[3],
		}
	}

	// ถ้า parse ไม่ได้ เก็บ raw
	return LogEntry{
		Timestamp: time.Now(),
		Level:     "UNKNOWN",
		Service:   service,
		Message:   line,
	}
}

// TailFile — คล้าย `tail -f`
func TailFile(path, service string, out chan<- LogEntry) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Seek(0, os.SEEK_END) // เริ่มจาก end ของ file
	scanner := bufio.NewScanner(f)
	
	for scanner.Scan() {
		out <- ParseLine(scanner.Text(), service)
	}
	
	return scanner.Err()
}
