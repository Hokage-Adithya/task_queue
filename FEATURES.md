# Advanced Task Queue - New Features

## 🎯 Priority Queues

Workers process tasks based on priority levels (1-5, with 5 being highest priority).

**Example:**
```bash
POST /api/tasks
{
  "type": "email",
  "payload": "urgent@example.com",
  "priority": 5  # High priority
}
```

**Implementation:**
- Priority stored in Task struct
- Workers display priority in logs: "Worker 1 processing task (priority: 5/5)"
- Dashboard shows priority badge for each task

---

## ⏰ Task Scheduling

Schedule tasks to run at a future time. Tasks automatically move from "scheduled" to "pending" when their time arrives.

**Example:**
```bash
POST /api/tasks
{
  "type": "email",
  "payload": "schedule@example.com",
  "scheduled_for": "2025-12-07T10:00:00Z"  # ISO 8601 format
}
```

**How It Works:**
- Task created with status "scheduled"
- Background goroutine checks every 5 seconds
- When time arrives, task automatically queued for processing
- `ProcessScheduledTasks()` runs continuously

**API Response:**
- Task shows status "scheduled" with scheduled time
- Stats include `scheduled_count` for monitoring

---

## 🔄 Automatic Retries

Failed tasks automatically retry up to `max_retries` times.

**Example:**
```bash
POST /api/tasks
{
  "type": "email",
  "payload": "retry@example.com",
  "max_retries": 3  # Up to 3 automatic retries
}
```

**Features:**
- `retry_count` tracks current attempt
- Failed tasks show retry progress: "🔄 Retries: 1/3"
- Manual retry via `POST /api/tasks/:id/retry`
- Error message stored in task for debugging
- Tasks move to "failed" status only after max retries exceeded

**Dashboard UI:**
- Shows retry badge with progress
- "Retry" button appears for failed tasks with retries remaining
- Error message displayed in red

---

## 🔗 Webhook Notifications

POST task details to a webhook URL when task completes.

**Example:**
```bash
POST /api/tasks
{
  "type": "email",
  "payload": "webhook@example.com",
  "webhook": "https://myapp.com/task-completed"  # Notified here
}
```

**How It Works:**
- When task completes successfully, webhook triggered
- Full task details sent as JSON payload
- Runs asynchronously (doesn't block worker)
- Error handling with panic recovery

**Payload Example:**
```json
{
  "id": "uuid-here",
  "type": "email",
  "payload": "webhook@example.com",
  "status": "completed",
  "priority": 3,
  "webhook": "https://myapp.com/task-completed",
  "completed_at": "2025-12-06T12:34:56Z"
}
```

**Dashboard:**
- Webhook URL shown with 🔗 badge
- Only sent on successful completion

---

## 📊 Enhanced Statistics

New metrics added to track advanced features:

```bash
GET /api/stats
```

**Response:**
```json
{
  "total_tasks": 150,
  "pending_count": 12,
  "processing_count": 3,
  "completed_count": 128,
  "failed_count": 5,        // NEW - failed task tracking
  "scheduled_count": 2,     // NEW - upcoming scheduled tasks
  "worker_count": 3
}
```

**Dashboard Shows:**
- 6 stat cards (was 4 before)
- Failed and scheduled count monitoring
- Real-time updates every 2 seconds

---

## 🎨 Enhanced Dashboard

### Tabs for Task Creation
- **Basic Tab**: Quick task creation (type + payload)
- **Advanced Tab**: Full control (priority, scheduling, retries, webhooks)

### Task Display Enhanced
- Priority badge (⭐ 1-5)
- Status color coding:
  - Blue: Pending
  - Green: Completed
  - Red: Failed
  - Orange: Scheduled
- Retry progress indicator
- Error messages in red
- Webhook indicator 🔗

### Real-time Stats
- 6 cards showing all metrics
- Auto-refresh every 2 seconds
- Mobile responsive design

---

## 🔧 Task Structure (Enhanced)

```go
type Task struct {
    ID           string    // UUID
    Type         string    // email, image, webhook
    Payload      string    // Task data
    Status       string    // pending, processing, completed, failed, scheduled
    Priority     int       // 1-5 (NEW)
    ScheduledFor time.Time // Future execution time (NEW)
    RetryCount   int       // Current retry attempt (NEW)
    MaxRetries   int       // Max retry attempts (NEW)
    Webhook      string    // Callback URL (NEW)
    CreatedAt    time.Time
    UpdatedAt    time.Time
    CompletedAt  time.Time // Task finish time (NEW)
    Error        string    // Error message (NEW)
}
```

---

## 📝 Worker Enhancements

Workers now:
1. Log retry progress: "Worker 1 retrying task (attempt 1/3)"
2. Handle failures gracefully with retry logic
3. Send webhook notifications on success
4. Process scheduled tasks in background
5. Show priority in log messages

**Log Example:**
```
Worker 1 processing task (type: email, priority: 5, attempt: 1/3)
Worker 1 completed task abc123... ✅
Worker 2 task xyz789 failed: Connection timeout ⚠️
Worker 2 retrying task xyz789...
🔗 Sending webhook to https://myapp.com/callback for task abc123...
⏰ Scheduled task def456 moved to queue
```

---

## 🚀 New API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/tasks` | Create task (basic or advanced) |
| GET | `/api/tasks` | List all tasks |
| GET | `/api/tasks/:id` | Get task by ID |
| GET | `/api/tasks/:id/details` | Get full task details |
| **POST** | **`/api/tasks/:id/retry`** | **Retry failed task (NEW)** |
| GET | `/api/stats` | Queue statistics (enhanced) |
| GET | `/api/health` | Health check |

---

## 💡 Usage Examples

### Create high-priority email task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": "urgent@example.com",
    "priority": 5
  }'
```

### Schedule task for 1 hour from now
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": "scheduled@example.com",
    "scheduled_for": "'$(date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%SZ)'"
  }'
```

### Create task with automatic retries + webhook
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": "robust@example.com",
    "priority": 3,
    "max_retries": 3,
    "webhook": "https://myapp.com/task-complete"
  }'
```

### Manually retry a failed task
```bash
curl -X POST http://localhost:8080/api/tasks/[TASK_ID]/retry
```

---

## 📚 Implementation Details

**New Functions Added:**

| Function | Purpose |
|----------|---------|
| `CanRetryTask()` | Check if task should retry |
| `RetryTask()` | Put task back in queue |
| `ProcessScheduledTasks()` | Move scheduled tasks to queue |
| `SendWebhookNotification()` | POST to webhook URL |
| `scheduledTaskProcessor()` | Background scheduler |
| `retryTaskHandler()` | API endpoint handler |
| `taskDetailsHandler()` | Enhanced task endpoint |

**Modified Functions:**

- `worker()` - Added retry logic, webhook calls
- `startWorkerPool()` - Added scheduler goroutine
- `CreateTask()` - Handles scheduled tasks
- `GetQueueStats()` - Counts failed/scheduled tasks

---

## ✅ Testing the Features

### 1. Test Priority
```bash
# Create 3 tasks with different priorities
POST /api/tasks {"type":"email","payload":"low@example.com","priority":1}
POST /api/tasks {"type":"email","payload":"high@example.com","priority":5}
POST /api/tasks {"type":"email","payload":"medium@example.com","priority":3}
# High priority task processes first (visible in stats)
```

### 2. Test Scheduling
```bash
# Schedule 1 minute in future
POST /api/tasks {"type":"email","payload":"future@example.com","scheduled_for":"2025-12-06T16:35:00Z"}
# Check stats - shows 1 scheduled
# Wait 1 minute - task moves to pending automatically
```

### 3. Test Retries
```bash
# Create with max_retries=2
POST /api/tasks {"type":"email","payload":"retry@example.com","max_retries":2}
# Simulate failure in worker (or wait for natural failure)
# Check dashboard - shows "🔄 Retries: 0/2"
# Click Retry button
# Task reprocesses
```

### 4. Test Webhooks
```bash
# Set up webhook receiver first (e.g., webhook.site)
POST /api/tasks {"type":"email","payload":"webhook@example.com","webhook":"https://webhook.site/[ID]"}
# Task completes
# Check webhook.site - receives POST with task data
```

---

## 🎯 Production Readiness

✅ **Ready for Production:**
- Error handling with panic recovery
- Async webhook calls (non-blocking)
- Retry logic prevents cascading failures
- Scheduled task processor runs independently
- Task status tracking for monitoring
- Comprehensive logging
- CORS enabled for cross-origin requests

🔧 **Considerations for Production:**
- Replace Mailtrap with production SMTP
- Add database persistence (PostgreSQL/MongoDB)
- Implement rate limiting
- Add authentication/authorization
- Set up monitoring/alerting
- Configure backup Redis
- Implement graceful shutdown

---

## 📈 What's Next?

Consider adding:
1. **Dead Letter Queue** - For permanently failed tasks
2. **Task Batching** - Process multiple related tasks
3. **Delayed Processing** - Exponential backoff retry
4. **Task Filtering** - Query by status, type, priority
5. **Admin Panel** - Delete/cancel tasks
6. **Metrics Export** - Prometheus format
7. **Task Dependencies** - Chain task execution
8. **Rate Limiting** - Per-user/IP rate limits
