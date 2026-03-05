-- ClickHouse: Columnar storage สำหรับ time-series
CREATE TABLE IF NOT EXISTS metrics (
    timestamp  DateTime64(3),          -- millisecond precision
    name       LowCardinality(String), -- enum-like, เร็วกว่า String
    value      Float64,
    labels     Map(String, String)    -- {"method":"GET","status":"200"}
) ENGINE = MergeTree()
PARTITION BY toDate(timestamp)      -- partition ทีละวัน
ORDER BY     (name, timestamp)        -- index หลัก
TTL          timestamp + INTERVAL 30 DAY  -- ลบข้อมูลเก่า auto
SETTINGS     index_granularity = 8192;

-- Query ตัวอย่าง: avg CPU per minute (เร็วมาก)
SELECT
    toStartOfMinute(timestamp) AS minute,
    avg(value)                  AS avg_val,
    max(value)                  AS max_val
FROM  metrics
WHERE name = 'cpu_usage_percent'
  AND timestamp >= now() - INTERVAL 1 HOUR
GROUP BY minute
ORDER BY minute;
