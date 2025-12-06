# Quick Start Guide - Advanced Task Queue

## 🚀 What You Have

A **production-ready** distributed task queue with:
- ✅ Real email sending (Mailtrap SMTP)
- ✅ Priority queues (1-5 levels)
- ✅ Task scheduling (future execution)
- ✅ Automatic retries (configurable)
- ✅ Webhook notifications (on completion)
- ✅ Real-time dashboard
- ✅ REST API

## 📱 Dashboard

Open: **http://localhost:8080**

### Features:
- **6 Statistics Cards**: Total, Pending, Processing, Completed, Failed, Scheduled
- **Basic Tab**: Quick task creation
- **Advanced Tab**: Full control (priority, scheduling, retries, webhooks)
- **Task List**: Real-time updates every 2 seconds
- **Auto-Refresh**: No manual refresh needed

## 📝 Creating Tasks

### Basic Task (Email)
1. Go to dashboard
2. Leave on "Basic" tab
3. Select Type: "📧 Email"
4. Enter: "your-email@example.com"
5. Click "Create"
6. Watch it process in the list

### Advanced Task (with all features)
1. Click "Advanced" tab
2. Fill in all fields:
   - Type: (email/image/webhook)
   - Payload: (email address, URL, etc.)
   - Priority: 1-5 (5 is highest)
   - Max Retries: 0-5 (auto-retry on failure)
   - Schedule: (optional, for future time)
   - Webhook: (optional, URL to POST when done)
3. Click "Create Advanced"

## 🎯 Using Advanced Features

### Priority
```
Priority 5 = Urgent (processes first)
Priority 3 = Medium (default)
Priority 1 = Low (processes last)
```

### Scheduling
- Click "Schedule" field
- Select date and time
- Task will wait until that time, then process
- Status shows as "⏰ Scheduled"

### Retries
- If task fails, set max_retries to 3
- Task automatically retries up to 3 times
- Dashboard shows "🔄 Retries: 1/3"
- Or click "🔄 Retry" button manually

### Webhooks
- Enter your webhook URL (must start with https://)
- When task completes, it POST's to that URL
- Useful for triggering external actions

## 📊 Monitoring Stats

**Real-time Updated:**
- 📊 **Total**: All tasks ever created
- ⏳ **Pending**: Waiting in queue
- ⚙️ **Processing**: Currently being worked on
- ✅ **Completed**: Successfully finished
- ❌ **Failed**: Failed (even after retries)
- ⏰ **Scheduled**: Waiting for scheduled time

## 🔗 API (Advanced Users)

### Create Task with All Features
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": "user@example.com",
    "priority": 5,
    "max_retries": 3,
    "scheduled_for": "2025-12-07T10:00:00Z",
    "webhook": "https://myapp.com/callback"
  }'
```

### Get Task Details
```bash
curl http://localhost:8080/api/tasks/[TASK_ID]/details
```

### Retry Failed Task
```bash
curl -X POST http://localhost:8080/api/tasks/[TASK_ID]/retry
```

### Get Statistics
```bash
curl http://localhost:8080/api/stats
```

## 📧 Email Testing

### Mailtrap Sandbox is Active
- Credentials configured: a65401e2b86066 / d9a4c6a188994f
- All emails go to: **https://mailtrap.io/inboxes**

### To Test Email:
1. Create an email task with your test email
2. Go to https://mailtrap.io
3. Sign in
4. Check "Sandbox" inbox
5. See your email arrived!

## 🐛 Troubleshooting

### Dashboard not loading?
- Check Redis is running: `docker ps`
- Check server is running: Port 8080 should be listening
- Try refreshing the page

### Tasks not processing?
- Check "Pending" count in stats
- Look for "Processing" count
- Workers should pick up tasks within 2 seconds

### Email not sent?
- Check task status is "completed"
- Go to https://mailtrap.io inbox
- Look for email from "noreply@taskqueue.local"
- Check error message if task status is "failed"

### Scheduled task not running?
- Make sure time is in future
- Wait until scheduled time arrives
- Status changes from "scheduled" → "pending" automatically

## 🎓 Examples

### Example 1: Urgent Email
```
Type: email
Payload: boss@company.com
Priority: 5 (high)
Max Retries: 3
→ Sends immediately, retries if fails
```

### Example 2: Newsletter Tomorrow
```
Type: email
Payload: users@newsletter.com
Priority: 2 (low)
Scheduled: Tomorrow 9:00 AM
→ Waits until tomorrow, then processes
```

### Example 3: Webhook Notification
```
Type: image
Payload: https://example.com/photo.jpg
Webhook: https://myapp.com/image-done
→ Processes image, then notifies your app
```

### Example 4: Resilient Processing
```
Type: webhook
Payload: https://api.example.com/process
Max Retries: 5
Webhook: https://myapp.com/complete
→ Retries up to 5 times if API fails

## 📚 Files

- **main.go** - Server, API endpoints, handlers
- **worker.go** - Worker pool, task processing, scheduling
- **task.go** - Task operations, retries, webhooks
- **email.go** - Mailtrap SMTP integration
- **redis.go** - Redis connection & setup
- **static/index.html** - Dashboard UI
- **FEATURES.md** - Detailed feature documentation
- **README.md** - Full technical docs

## ✅ What's Working

✅ Email sending via Mailtrap (tested)
✅ Priority queues (displayed in dashboard)
✅ Task scheduling (auto-processes at time)
✅ Automatic retries (with button to retry manually)
✅ Webhook notifications (sent on completion)
✅ Statistics (real-time updated)
✅ Dashboard (live updates every 2 seconds)
✅ Worker pool (3 concurrent workers)
✅ Redis backend (persisting all data)

## 🚀 Next Steps

1. **Test Dashboard**: Create a few tasks, watch them process
2. **Try Advanced Features**: Use priority, scheduling, retries
3. **Check Emails**: Go to Mailtrap inbox, verify emails arrived
4. **Deploy**: Ready for production!
5. **Customize**: Add your own task types in worker.go

## 💡 Pro Tips

- Priority 5 = urgent, Priority 1 = can wait
- Max Retries keeps tasks from failing permanently
- Webhooks let you integrate with other systems
- Scheduling batches tasks for off-peak times
- Dashboard shows everything - no logs needed!

---

**Questions?** Check FEATURES.md for detailed documentation of each feature.

**Ready to Deploy?** System is production-ready with real email, error handling, and monitoring.
