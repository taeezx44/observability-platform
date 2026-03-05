# Self-hosted Observability Platform

Complete observability platform built with Go + React + ClickHouse. Free, open-source, and self-hosted.

## 🚀 Quick Start

```bash
# 1. Start infrastructure
docker compose up clickhouse kafka -d

# 2. Run database migrations
cat migrations/001_metrics.sql | docker exec -i clickhouse clickhouse-client --database=observability
cat migrations/002_logs.sql | docker exec -i clickhouse clickhouse-client --database=observability  
cat migrations/003_traces.sql | docker exec -i clickhouse clickhouse-client --database=observability

# 3. Start all services
docker compose up
```

Open http://localhost:3000 to see the dashboard.

## 📊 Features

### Phase 1: Metrics Pipeline ✅
- Prometheus scraper (15s interval)
- ClickHouse time-series storage
- React dashboard with live charts
- WebSocket real-time updates

### Phase 2: Log Collection ✅
- Multi-format log parser (JSON + plaintext)
- Full-text search with ClickHouse
- Live log streaming
- Level and service filtering

### Phase 3: Distributed Tracing ✅
- OpenTelemetry span storage
- Waterfall visualization
- Slow trace detection
- Service dependency mapping

### Phase 4: Alerting Engine ✅
- Rule-based alerting (threshold, duration)
- Slack webhook notifications
- Alert history and silencing
- Multi-severity support

## 🏗️ Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Frontend  │    │     API     │    │  Collector  │
│   (React)   │◄──►│   (Go)      │◄──►│   (Go)      │
│   :3000     │    │   :8080     │    │   (scraper) │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                           ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐
                   │ ClickHouse  │    │    Kafka    │
                   │   :8123     │    │   :9092     │
                   │   :9000     │    │             │
                   └─────────────┘    └─────────────┘
```

## 📁 Project Structure

```
observability-platform/
├── collector/               # Go: scraper + log agent
│   ├── cmd/main.go         # Main collector entrypoint
│   ├── scraper/scraper.go  # Prometheus scraper
│   ├── logger/parser.go    # Log parser
│   ├── tracer/span.go      # Tracing models
│   └── storage/clickhouse.go # ClickHouse client
├── api/                     # Go: REST API server
│   ├── cmd/main.go         # API server
│   ├── handlers/           # HTTP handlers
│   └── Dockerfile
├── alerting/                # Go: alert engine
│   ├── main.go             # Alert engine
│   ├── engine.go           # Rule evaluation
│   ├── rules.yaml          # Alert rules
│   └── Dockerfile
├── frontend/                # React dashboard
│   ├── src/
│   │   ├── pages/          # Dashboard pages
│   │   └── components/     # Reusable components
│   ├── Dockerfile
│   └── nginx.conf
├── migrations/              # SQL schemas
│   ├── 001_metrics.sql
│   ├── 002_logs.sql
│   └── 003_traces.sql
└── docker-compose.yml
```

## ⚙️ Configuration

### Environment Variables

```bash
# ClickHouse Database
CLICKHOUSE_URL=clickhouse:9000
CLICKHOUSE_USER=admin
CLICKHOUSE_PASSWORD=secret
CLICKHOUSE_DB=observability

# Metrics Scraper
SCRAPE_INTERVAL=15s
SCRAPE_TARGETS=http://app:8080/metrics,http://db:9090/metrics

# API Server
PORT=8080
VITE_API_URL=http://localhost:8080

# Alerting Engine
SLACK_WEBHOOK_URL=https://hooks.slack.com/...
RULES_PATH=./rules.yaml
ALERT_INTERVAL=30s

# Frontend (for development)
VITE_API_URL=http://localhost:8080

# Docker/Production
NODE_ENV=production
```

### Docker Compose Environment

Create `.env` file for Docker:

```bash
# .env
CLICKHOUSE_URL=clickhouse:9000
CLICKHOUSE_USER=admin
CLICKHOUSE_PASSWORD=secret123
CLICKHOUSE_DB=observability

SCRAPE_INTERVAL=15s
SCRAPE_TARGETS=http://api:8080/metrics

SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK
RULES_PATH=./rules.yaml
ALERT_INTERVAL=30s
```

### Render.com Environment

```bash
# Render environment variables
PORT=10000
GO_VERSION=1.26
NODE_VERSION=18
VITE_API_URL=https://observability-platform.onrender.com
```

### Alert Rules (alerting/rules.yaml)

```yaml
rules:
  - name: HighCPU
    metric: cpu_usage_percent
    condition: ">"
    threshold: 85
    for: 2m
    severity: critical

  - name: HighMemory
    metric: memory_usage_percent
    condition: ">"
    threshold: 90
    for: 5m
    severity: warning
```

## 🔧 Development

### Local Development

```bash
# Start dependencies
docker compose up clickhouse kafka -d

# Run migrations
./scripts/migrate.sh

# Start Go services (in separate terminals)
go run ./collector/cmd/main.go
go run ./api/cmd/main.go
go run ./alerting/main.go

# Start frontend
cd frontend && npm install && npm run dev
```

### Adding New Metrics

1. Expose `/metrics` endpoint on your service
2. Add to `SCRAPE_TARGETS` environment variable
3. Metrics will automatically appear in dashboard

### Adding New Alert Rules

1. Edit `alerting/rules.yaml`
2. Restart alerting service
3. Alerts will fire when conditions are met

## 📈 Monitoring the Platform

The platform monitors itself:

- **Collector health**: Scraping metrics, batch insert rates
- **API performance**: Request latency, error rates  
- **Database health**: ClickHouse query performance
- **Alert engine**: Rule evaluation status

## 🛠️ Troubleshooting

### Common Issues

**ClickHouse connection failed**
```bash
# Check ClickHouse is running
docker compose ps clickhouse

# Test connection
docker exec -it clickhouse clickhouse-client --database=observability
```

**No metrics showing**
```bash
# Check scraper logs
docker compose logs collector

# Verify targets are accessible
curl http://your-app:8080/metrics
```

**Alerts not firing**
```bash
# Check alert engine logs
docker compose logs alerting

# Verify rules syntax
docker exec alerting ./alerting -check-rules
```

## 🎯 Performance

- **ClickHouse**: Handles 1M+ metrics/second with proper partitioning
- **Scraping**: 15s intervals, batch inserts for efficiency
- **Frontend**: Real-time updates via WebSocket, 30s refresh
- **Storage**: 30-day TTL for metrics, 7-day for logs/traces

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make your changes
4. Add tests if applicable
5. Submit pull request

## 📄 License

MIT License - feel free to use in commercial projects.

## 🙏 Acknowledgments

- [ClickHouse](https://clickhouse.com/) - Fast analytical database
- [Prometheus](https://prometheus.io/) - Metrics format
- [Recharts](https://recharts.org/) - React chart library
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
