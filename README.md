# Distributed Task Queue System

A **production-ready** distributed task queue built with **Go**, **Gin**, and **Redis**. Features priority queues, task scheduling, automatic retries, webhook notifications, and real email sending via Mailtrap.

## 🚀 Features

**Core Features:**
- ✅ **RESTful API** - Gin framework with CORS support
- ✅ **Worker Pool** - 3 concurrent goroutines processing tasks
- ✅ **Redis Backend** - Queue storage, task metadata, and pub/sub events
- ✅ **Real Email** - Mailtrap SMTP integration (sandbox + production ready)

**Advanced Features:**
- 🎯 **Priority Queues** - Tasks with priority levels (1-5)
- ⏰ **Task Scheduling** - Schedule tasks for future execution
- 🔄 **Automatic Retries** - Configurable retry logic for failed tasks
- 🔗 **Webhook Notifications** - POST to external URLs on task completion
- 📊 **Dashboard** - Real-time monitoring with auto-refresh
- 📈 **Statistics** - Queue metrics (pending, processing, completed, failed, scheduled)

## 📋 API Endpoints

### Create Task (Basic)
```bash
POST /api/tasks
Content-Type: application/json

{
  "type": "email",
  "payload": "user@example.com"
}
```
**Response (201):**
```json
**Response (201):**
```json
{
  "id": "uuid-here",
  "type": "email",
  "payload": "user@example.com",
  "status": "pending",
  "priority": 3,
  "max_retries": 2,
  "retry_count": 0,
  "created_at": "2025-12-06T...",
  "updated_at": "2025-12-06T..."
}
```

### Create Task (Advanced - with priority, scheduling, retries, webhooks)
```bash
POST /api/tasks
Content-Type: application/json

{
  "type": "email",
  "payload": "user@example.com",
  "priority": 5,
  "max_retries": 3,
  "scheduled_for": "2025-12-06T18:30:00Z",
  "webhook": "https://example.com/webhook"
}
```

### Get Task
```bash
GET /api/tasks/:id
```

### Get Task Details
```bash
GET /api/tasks/:id/details
```

### Retry Failed Task
```bash
POST /api/tasks/:id/retry
```
**Response (200):**
```json
{
  "status": "retrying",
  "task_id": "uuid-here"
}
```

### List All Tasks
```bash
GET /api/tasks
```

### Queue Statistics
```bash
GET /api/stats
```
**Response (200):**
```json
{
  "total_tasks": 100,
  "pending_count": 10,
  "processing_count": 3,
  "completed_count": 80,
  "failed_count": 5,
  "scheduled_count": 2,
  "worker_count": 3
}
```

### Health Check
```bash
GET /api/health
```
  "status": "ok"
}
```

## 🏗️ Project Structure

```
sk_queue/
├── main.go              # Gin server, routes, handlers
├── worker.go            # Worker pool, task processing
├── task.go              # Task model, Redis operations
├── redis.go             # Redis client initialization
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # Redis + app services
├── .env                 # Environment variables
├── go.mod               # Go dependencies
├── go.sum               # Dependency checksums
└── README.md            # Documentation
```

## 🔧 Tech Stack

- **Language:** Go 1.25+
- **Web Framework:** Gin v1.11.0
- **Redis Client:** go-redis/v9 v9.17.2
- **UUID:** google/uuid v1.6.0
- **CORS:** gin-contrib/cors v1.7.6
- **Database:** Redis 7
- **Containerization:** Docker

## ⚡ Quick Start

### Local Development (requires Redis running)

1. **Set up Go environment:**
   ```bash
   go env
   go version
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set Redis URL (optional, defaults to localhost:6379):**
   ```bash
   $env:REDIS_URL = "localhost:6379"
   ```

4. **Build and run:**
   ```bash
   go build -o queue-server.exe .
   ./queue-server.exe
   ```

5. **Server runs on:** `http://localhost:8080`

### Docker Compose (recommended)

```bash
docker-compose up --build
```

Services:
- **app:** http://localhost:8080
- **redis:** localhost:6379

## 🧪 Testing API

### Create a task:
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"type":"email","payload":"test@example.com"}'
```

### Get task status:
```bash
curl http://localhost:8080/tasks/{task-id}
```

### List all tasks:
```bash
curl http://localhost:8080/tasks
```

### Get statistics:
```bash
curl http://localhost:8080/stats
```

### Health check:
```bash
curl http://localhost:8080/health
```

## 👷 Worker Pool

- **3 concurrent workers** continuously process tasks
- Uses **Redis BRPOP** (blocking pop) for queue consumption
- Task states: `pending` → `processing` → `completed`
- Simulates processing with type-based delays:
  - email: 2 seconds
  - image: 3 seconds
  - webhook: 2 seconds
- Publishes events via Redis Pub/Sub

## 🗃️ Redis Storage

### Queue (List)
- Key: `task_queue`
- Operation: LPUSH (add), BRPOP (consume)

### Task Metadata (Hash)
- Key: `tasks`
- Field: task ID
- Value: JSON serialized task

### Events (Pub/Sub)
- Channel: `task_events`
- Messages: JSON event notifications

## 🚢 Deployment

### Railway
1. Connect GitHub repo
2. Set environment variable: `REDIS_URL=<redis-connection-string>`
3. Deploy

### Render
1. Connect GitHub repo
2. Create Redis service
3. Set environment variable: `REDIS_URL=<redis-connection-string>`
4. Deploy

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | localhost:6379 | Redis connection string |
| `GIN_MODE` | debug | gin/release or debug |

## 🔒 Error Handling

- Graceful shutdown on SIGINT/SIGTERM
- Redis connection validation
- Task not found returns 404
- Invalid requests return 400
- Server errors return 500

## 📊 Logging

Workers log activity:
```
Worker 1 waiting for task
Worker 1 received task: {task-id}
Worker 1 processing task {task-id} (type: email)
Worker 1 completed task {task-id}
```

## 📜 License

MIT License - See LICENSE file for details
