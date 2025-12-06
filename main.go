package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// loadEnv loads environment variables from .env file if it exists
func loadEnv() {
	envFile := ".env"
	data, err := os.ReadFile(envFile)
	if err != nil {
		return // .env file not found, that's ok
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

func main() {
	loadEnv()

	InitRedis()
	defer redisClient.Close()

	InitEmailSender()

	r := gin.New() // Use New() instead of Default() to avoid built-in recovery issues
	r.Use(gin.Logger())

	// Add explicit recovery middleware with logging
	r.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("🚨 PANIC in request: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
	})

	// CORS middleware
	r.Use(cors.Default())

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// API routes
	r.POST("/api/tasks", createTaskHandler)
	r.GET("/api/tasks/:id", getTaskHandler)
	r.GET("/api/tasks/:id/details", taskDetailsHandler)
	r.POST("/api/tasks/:id/retry", retryTaskHandler)
	r.GET("/api/tasks", listTasksHandler)
	r.GET("/api/stats", statsHandler)
	r.GET("/api/health", healthHandler)

	// Start worker pool
	go startWorkerPool()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Server listening on http://0.0.0.0:8080")
	log.Println("API endpoints:")
	log.Println("  POST   /api/tasks - Create task (with priority, scheduling, webhook, retries)")
	log.Println("  GET    /api/tasks - List all tasks")
	log.Println("  GET    /api/tasks/:id - Get task")
	log.Println("  GET    /api/tasks/:id/details - Get task details")
	log.Println("  POST   /api/tasks/:id/retry - Retry failed task")
	log.Println("  GET    /api/stats - Queue statistics")
	log.Println("  GET    /api/health - Health check")

	// For testing, just run the server without signal handling
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Server error: %v", err)
	}
} // Handler: Create new task
func createTaskHandler(c *gin.Context) {
	log.Println("🔵 createTaskHandler START")

	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 PANIC in createTaskHandler: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	log.Println("🔵 Reading body...")
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	log.Printf("🔵 Body: %s", string(bodyBytes))

	var req struct {
		Type         string    `json:"type"`
		Payload      string    `json:"payload"`
		Priority     int       `json:"priority"`
		ScheduledFor time.Time `json:"scheduled_for,omitempty"`
		Webhook      string    `json:"webhook,omitempty"`
		MaxRetries   int       `json:"max_retries"`
	}

	log.Println("🔵 Parsing JSON...")
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		log.Printf("JSON parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	log.Printf("🔵 Parsed: type=%s, payload=%s, priority=%d, scheduled=%v",
		req.Type, req.Payload, req.Priority, req.ScheduledFor)

	if req.Type == "" || req.Payload == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type and payload required"})
		return
	}

	// Set defaults
	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 3 // medium priority
	}
	if req.MaxRetries == 0 {
		req.MaxRetries = 2
	}

	taskID := uuid.New().String()
	status := "pending"
	if !req.ScheduledFor.IsZero() && req.ScheduledFor.After(time.Now()) {
		status = "scheduled"
	}

	task := Task{
		ID:           taskID,
		Type:         req.Type,
		Payload:      req.Payload,
		Status:       status,
		Priority:     req.Priority,
		ScheduledFor: req.ScheduledFor,
		Webhook:      req.Webhook,
		MaxRetries:   req.MaxRetries,
		RetryCount:   0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	log.Printf("📝 Creating task: %s (type: %s, priority: %d, status: %s)", taskID, req.Type, req.Priority, status)

	// Call CreateTask
	err = CreateTask(&task)
	log.Printf("✅ CreateTask returned (err=%v)", err)

	if err != nil {
		log.Printf("❌ Error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Task created, publishing event for %s", taskID)
	PublishTaskEvent(taskID, "created")
	log.Printf("✅ Event published, returning response")
	c.JSON(http.StatusCreated, task)
	log.Println("🔵 createTaskHandler END")
} // Handler: Get task by ID
func getTaskHandler(c *gin.Context) {
	taskID := c.Param("id")
	task, err := GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Handler: List all tasks
func listTasksHandler(c *gin.Context) {
	tasks, err := ListAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Handler: Queue statistics
func statsHandler(c *gin.Context) {
	stats, err := GetQueueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Handler: Retry failed task
func retryTaskHandler(c *gin.Context) {
	taskID := c.Param("id")
	err := RetryTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "retrying", "task_id": taskID})
}

// Handler: Get task by ID with full details
func taskDetailsHandler(c *gin.Context) {
	taskID := c.Param("id")
	task, err := GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// Handler: Health check
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
