# 🐳 Docker Setup Guide

## 🚀 Quick Start

### Option 1: Simple Local Run (Recommended)
```bash
# Start the main API
cd simple-api
go run main.go

# Access at:
# - API: http://localhost:8084/
# - Metrics: http://localhost:8084/metrics
# - Health: http://localhost:8084/health
# - API Info: http://localhost:8084/api
```

### Option 2: Full Docker Setup
```bash
# Use simple Docker compose
docker-compose -f docker-compose.simple.yml up -d

# Access points:
# - Frontend: http://localhost:3000
# - API: http://localhost:8080
# - ClickHouse: http://localhost:8123
# - Kafka: localhost:9092
```

## 📊 Running Services

### ✅ Currently Running:
- **Simple API** (Port 8084) - ✅ ACTIVE
  - Root: http://localhost:8084/ → "Observability API Running"
  - Metrics: http://localhost:8084/metrics (Prometheus-style)
  - Health: http://localhost:8084/health
  - API Info: http://localhost:8084/api

### 🎯 Demo Services (When Disk Space Available):
- **Demo API** (Port 8081)
- **Demo Worker** (Port 8082)  
- **Demo Database** (Port 8083)

## 🔧 Docker Services

### Core Infrastructure:
- **ClickHouse** - Time-series database
- **Kafka** - Message streaming
- **API Service** - REST API with metrics
- **Frontend** - Web dashboard

### Features:
- Prometheus-style metrics
- Health checks
- Service discovery
- Auto-restart
- Volume persistence

## 📋 Access Points

### 🌐 Web Interfaces:
- **Frontend**: http://localhost:3000
- **ClickHouse UI**: http://localhost:8123
- **API Root**: http://localhost:8084/

### 📊 API Endpoints:
- **Metrics**: http://localhost:8084/metrics
- **Health**: http://localhost:8084/health
- **API Info**: http://localhost:8084/api

### 🗄️ Database:
- **ClickHouse HTTP**: http://localhost:8123
- **ClickHouse Native**: localhost:9000
- **Kafka**: localhost:9092

## 🛠️ Troubleshooting

### Port Conflicts:
- Main API uses port 8084 (instead of 8080)
- Demo services use 8081-8083
- Frontend uses 3000

### Disk Space Issues:
- Demo services require disk space for Go compilation
- Use simple API for minimal setup
- Docker images may require cleanup

### Docker Issues:
- Ensure Docker Desktop is running
- Check disk space in Docker
- Use `docker system prune` if needed

## 🎯 Portfolio Demo

### Quick Demo:
1. Start simple API: `cd simple-api && go run main.go`
2. Open browser: http://localhost:8084/
3. Show "Observability API Running" ✅
4. Show metrics: http://localhost:8084/metrics
5. Show health: http://localhost:8084/health

### Full Demo:
1. Run `start-demo.bat` for all services
2. Access all endpoints
3. Show real-time metrics
4. Show WebSocket logs

## 🏆 Success Metrics

✅ **API Running**: http://localhost:8084/  
✅ **Metrics Available**: http://localhost:8084/metrics  
✅ **Health Checks**: http://localhost:8084/health  
✅ **Portfolio Ready**: Clean, professional responses  
✅ **Production Style**: Prometheus metrics format  

## 📝 Notes

- Simple API is lightweight and fast
- Docker setup provides full infrastructure
- Demo services show realistic workloads
- All services have portfolio-friendly root routes
