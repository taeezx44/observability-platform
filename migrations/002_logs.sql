CREATE TABLE IF NOT EXISTS logs (
    timestamp DateTime64(3),
    level     LowCardinality(String),  -- ERROR/WARN/INFO/DEBUG
    service   LowCardinality(String),
    message   String,
    fields    Map(String, String),

    -- Full-text search index
    INDEX idx_message message TYPE tokenbf_v1(32768, 3, 0)
          GRANULARITY 4
) ENGINE = MergeTree()
PARTITION BY toDate(timestamp)
ORDER BY     (service, level, timestamp)
TTL          timestamp + INTERVAL 7 DAY;

-- Search errors in last hour
SELECT timestamp, service, message
FROM  logs
WHERE level = 'ERROR'
  AND timestamp >= now() - INTERVAL 1 HOUR
  AND message LIKE '%database%'
ORDER BY timestamp DESC
LIMIT 100;
