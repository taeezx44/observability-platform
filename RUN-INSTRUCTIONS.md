# 🚀 Run Instructions - Observability Platform

## ✅ Status Check Results

**All prerequisites are installed and ready!**

- ✅ **Go 1.26.0** - Installed and working
- ✅ **Node.js v24.13.1** - Installed and working  
- ✅ **Docker 29.2.1** - Installed but needs to be started
- ✅ **Frontend** - Built successfully (dist/ folder ready)

## 🐳 Step 1: Start Docker Desktop

**Manual Steps:**
1. Open **Start Menu**
2. Search for **"Docker Desktop"**
3. Click to open it
4. Wait for it to fully initialize (2-3 minutes)
5. You'll see the Docker whale icon in system tray when ready

## 🚀 Step 2: Run the Platform

Once Docker Desktop is running, execute:

```bash
# In the observability-platform directory:
quick-start.bat
```

This will:
- Start ClickHouse and Kafka
- Run database migrations
- Build and start all services
- Open the dashboard

## 📊 Access Points

When running, you can access:

- **📊 Dashboard**: http://localhost:3000
  - Metrics charts with live updates
  - Log streaming with search
  - Distributed tracing waterfall
  - Real-time WebSocket updates

- **🔍 API**: http://localhost:8080
  - REST endpoints for metrics, logs, traces
  - WebSocket for live updates
  - Health check: `/health`

- **🗄️ ClickHouse**: http://localhost:8123
  - Database interface
  - Query browser
  - Performance metrics

- **📨 Kafka**: localhost:9092
  - Message queue for log streaming
  - Internal service communication

## 🛑 To Stop

```bash
docker compose down
```

## 📋 What's Running

### Services
- **collector** - Prometheus scraper + log agent
- **api** - REST API server with WebSocket
- **alerting** - Rule engine with Slack notifications
- **frontend** - React dashboard (nginx)
- **clickhouse** - Time-series database
- **kafka** - Message queue

### Features
- **Metrics Pipeline** - 15s scraping, live charts
- **Log Collection** - Multi-format parser, full-text search
- **Distributed Tracing** - Waterfall visualization
- **Smart Alerting** - Threshold rules, Slack notifications

## 🔧 Troubleshooting

### Docker Desktop Issues
- Make sure WSL 2 is enabled (Windows will prompt)
- Restart Docker Desktop if it shows errors
- Check system resources (needs 4GB+ RAM)

### Port Conflicts
- Make sure ports 3000, 8080, 8123, 9000, 9092 are free
- Stop other services using these ports if needed

### Build Issues
- Frontend already built (dist/ folder exists)
- Go services compile successfully
- All dependencies resolved

## 🎯 Next Steps

1. **Start Docker Desktop** (manual step)
2. **Run `quick-start.bat`**
3. **Open http://localhost:3000**
4. **Add your Prometheus targets** to start seeing real metrics

## 📈 Monitoring the Platform

The platform monitors itself:
- Collector scraping metrics
- API request latency  
- Database query performance
- Alert rule evaluation status

**🎉 Ready to run! Just start Docker Desktop and execute `quick-start.bat`!**
