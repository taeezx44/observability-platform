CREATE TABLE IF NOT EXISTS spans (
    trace_id   String,
    span_id    String,
    parent_id  String,
    service    LowCardinality(String),
    operation  String,
    start_time DateTime64(3),
    end_time   DateTime64(3),
    status     LowCardinality(String), -- OK/ERROR
    tags       Map(String, String),
    logs       Array(Map(String, String))
) ENGINE = MergeTree()
PARTITION BY toDate(start_time)
ORDER BY     (trace_id, span_id)
TTL          start_time + INTERVAL 7 DAY;

-- Get trace by ID with all spans
SELECT *
FROM spans
WHERE trace_id = {trace_id:String}
ORDER BY start_time;

-- Find slow traces
SELECT trace_id, service, operation, count() as span_count, 
       max(end_time - start_time) as duration_ms
FROM spans
WHERE start_time >= now() - INTERVAL 1 HOUR
GROUP BY trace_id, service, operation
HAVING duration_ms > 1000
ORDER BY duration_ms DESC
LIMIT 50;
