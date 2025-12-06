package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Payload      string    `json:"payload"`
	Status       string    `json:"status"`
	Priority     int       `json:"priority"` // 1=low, 5=high
	ScheduledFor time.Time `json:"scheduled_for,omitempty"`
	RetryCount   int       `json:"retry_count"`
	MaxRetries   int       `json:"max_retries"`
	Webhook      string    `json:"webhook,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
	Error        string    `json:"error,omitempty"`
}

type QueueStats struct {
	TotalTasks      int64 `json:"total_tasks"`
	PendingCount    int64 `json:"pending_count"`
	ProcessingCount int64 `json:"processing_count"`
	CompletedCount  int64 `json:"completed_count"`
	FailedCount     int64 `json:"failed_count"`
	ScheduledCount  int64 `json:"scheduled_count"`
	WorkerCount     int   `json:"worker_count"`
}

// CreateTask: Add task to Redis queue and store metadata
func CreateTask(task *Task) error {
	ctx := context.Background()

	// Store task metadata in Redis hash
	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if err := redisClient.HSet(ctx, "tasks", task.ID, taskData).Err(); err != nil {
		return err
	}

	// Add task ID to queue
	if err := redisClient.LPush(ctx, "task_queue", task.ID).Err(); err != nil {
		return err
	}

	log.Printf("Task %s created and queued", task.ID)
	return nil
}

// GetTaskByID: Fetch task metadata from Redis hash
func GetTaskByID(taskID string) (*Task, error) {
	ctx := context.Background()
	result, err := redisClient.HGet(ctx, "tasks", taskID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, err
	}

	var task Task
	if err := json.Unmarshal([]byte(result), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTaskStatus: Update task status in Redis hash
func UpdateTaskStatus(taskID string, status string) error {
	ctx := context.Background()

	task, err := GetTaskByID(taskID)
	if err != nil {
		return err
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	taskData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if err := redisClient.HSet(ctx, "tasks", taskID, taskData).Err(); err != nil {
		return err
	}

	log.Printf("Task %s status updated to %s", taskID, status)
	return nil
}

// PublishTaskEvent: Publish event via Redis Pub/Sub
func PublishTaskEvent(taskID string, event string) {
	if redisClient == nil {
		log.Printf("⚠️  Redis client not initialized, skipping event publish for %s:%s", taskID, event)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ PANIC in PublishTaskEvent: %v", r)
		}
	}()

	ctx := context.Background()
	message := taskID + ":" + event
	if err := redisClient.Publish(ctx, "task_events", message).Err(); err != nil {
		log.Printf("Error publishing event: %v", err)
	}
}

// ListAllTasks: Get all tasks from Redis hash
func ListAllTasks() ([]*Task, error) {
	ctx := context.Background()
	result, err := redisClient.HGetAll(ctx, "tasks").Result()
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, taskData := range result {
		var task Task
		if err := json.Unmarshal([]byte(taskData), &task); err != nil {
			continue
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// GetQueueStats: Get statistics about queue
func GetQueueStats() (*QueueStats, error) {
	ctx := context.Background()

	// Count pending tasks
	pendingCount, err := redisClient.LLen(ctx, "task_queue").Result()
	if err != nil {
		return nil, err
	}

	// Count all tasks
	totalCount, err := redisClient.HLen(ctx, "tasks").Result()
	if err != nil {
		return nil, err
	}

	// Count by status
	tasks, err := ListAllTasks()
	if err != nil {
		return nil, err
	}

	processingCount := int64(0)
	completedCount := int64(0)

	for _, task := range tasks {
		switch task.Status {
		case "processing":
			processingCount++
		case "completed":
			completedCount++
		}
	}

	failedCount := int64(0)
	scheduledCount := int64(0)
	for _, task := range tasks {
		if task.Status == "failed" {
			failedCount++
		}
		if !task.ScheduledFor.IsZero() && task.ScheduledFor.After(time.Now()) {
			scheduledCount++
		}
	}

	return &QueueStats{
		TotalTasks:      totalCount,
		PendingCount:    pendingCount,
		ProcessingCount: processingCount,
		CompletedCount:  completedCount,
		FailedCount:     failedCount,
		ScheduledCount:  scheduledCount,
		WorkerCount:     3,
	}, nil
}

// SendWebhookNotification: POST task completion to webhook URL
func SendWebhookNotification(task *Task) {
	if task.Webhook == "" {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ PANIC in webhook: %v", r)
			}
		}()

		payload, _ := json.Marshal(task)
		log.Printf("🔗 Sending webhook to %s for task %s", task.Webhook, task.ID)

		// Simple HTTP POST (in production, use proper HTTP client with timeout)
		// For now, just log it
		log.Printf("🔗 Webhook payload: %s", string(payload))
	}()
}

// CanRetryTask: Check if task should be retried
func CanRetryTask(task *Task) bool {
	return task.Status == "failed" && task.RetryCount < task.MaxRetries
}

// RetryTask: Put task back in queue for retry
func RetryTask(taskID string) error {
	task, err := GetTaskByID(taskID)
	if err != nil {
		return err
	}

	if !CanRetryTask(task) {
		return nil
	}

	task.RetryCount++
	task.Status = "pending"
	task.Error = ""
	task.UpdatedAt = time.Now()

	ctx := context.Background()
	taskData, _ := json.Marshal(task)

	if err := redisClient.HSet(ctx, "tasks", taskID, taskData).Err(); err != nil {
		return err
	}

	if err := redisClient.LPush(ctx, "task_queue", taskID).Err(); err != nil {
		return err
	}

	log.Printf("📋 Task %s retrying (attempt %d/%d)", taskID, task.RetryCount, task.MaxRetries)
	return nil
}

// ProcessScheduledTasks: Move scheduled tasks to queue when ready
func ProcessScheduledTasks() {
	ctx := context.Background()
	tasks, err := ListAllTasks()
	if err != nil {
		return
	}

	for _, task := range tasks {
		if !task.ScheduledFor.IsZero() && task.ScheduledFor.Before(time.Now()) && task.Status == "scheduled" {
			task.Status = "pending"
			task.UpdatedAt = time.Now()
			taskData, _ := json.Marshal(task)
			redisClient.HSet(ctx, "tasks", task.ID, taskData)
			redisClient.LPush(ctx, "task_queue", task.ID)
			log.Printf("⏰ Scheduled task %s moved to queue", task.ID)
		}
	}
}
