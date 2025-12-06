package main

import (
	"context"
	"log"
	"time"
)

func startWorkerPool() {
	// Start worker goroutines
	for i := 1; i <= 3; i++ {
		go worker(i)
	}
	// Start scheduled task processor
	go scheduledTaskProcessor()
}

func scheduledTaskProcessor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ProcessScheduledTasks()
	}
}

func worker(id int) {
	for {
		log.Printf("Worker %d waiting for task", id)

		// BRPOP from Redis queue (blocking pop, 0 timeout = wait forever)
		result, err := redisClient.BRPop(context.Background(), 0, "task_queue").Result()
		if err != nil {
			log.Printf("Worker %d error: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(result) < 2 {
			continue
		}

		taskID := result[1]
		log.Printf("Worker %d received task: %s", id, taskID)

		// Fetch task metadata
		task, err := GetTaskByID(taskID)
		if err != nil {
			log.Printf("Worker %d failed to get task %s: %v", id, taskID, err)
			continue
		}

		// Update status to processing
		UpdateTaskStatus(taskID, "processing")
		log.Printf("Worker %d processing task %s (type: %s, priority: %d, attempt: %d/%d)",
			id, taskID, task.Type, task.Priority, task.RetryCount+1, task.MaxRetries+1)

		// Execute task based on type
		var execErr error
		delay := getTaskDelay(task.Type)
		time.Sleep(delay)

		switch task.Type {
		case "email":
			execErr = SendEmail(
				task.Payload,
				"Task Queue Notification",
				"Your email task has been processed by the task queue system!",
			)
		case "image":
			log.Printf("Worker %d processing image: %s", id, task.Payload)
		case "webhook":
			log.Printf("Worker %d calling webhook: %s", id, task.Payload)
		default:
			log.Printf("Worker %d unknown task type: %s", id, task.Type)
		}

		// Handle completion or retry
		if execErr == nil {
			UpdateTaskStatus(taskID, "completed")
			task.CompletedAt = time.Now()
			log.Printf("Worker %d completed task %s ✅", id, taskID)

			// Send webhook notification on success
			SendWebhookNotification(task)
		} else {
			log.Printf("Worker %d task %s failed: %v ⚠️", id, taskID, execErr)
			task.Error = execErr.Error()

			// Try retry if available
			if CanRetryTask(task) {
				log.Printf("Worker %d retrying task %s...", id, taskID)
				RetryTask(taskID)
			} else {
				UpdateTaskStatus(taskID, "failed")
				log.Printf("Worker %d marked task %s as failed (no more retries)", id, taskID)
			}
		}

		// Publish event
		PublishTaskEvent(taskID, "completed")
	}
}

func getTaskDelay(taskType string) time.Duration {
	switch taskType {
	case "email":
		return 2 * time.Second
	case "image":
		return 3 * time.Second
	case "webhook":
		return 2 * time.Second
	default:
		return 2 * time.Second
	}
}
