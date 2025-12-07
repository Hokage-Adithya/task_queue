package main

import (
	"context"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func parseUpstashURL(upstashURL string) (host string, password string, err error) {
	// Parse Upstash REST URL: https://joint-swine-8107.upstash.io
	// Convert to: joint-swine-8107.upstash.io:6379 with token as password
	u, err := url.Parse(upstashURL)
	if err != nil {
		return "", "", err
	}

	// Extract hostname from https://host.upstash.io -> host.upstash.io
	host = u.Host
	if host == "" {
		host = strings.TrimPrefix(u.String(), "https://")
		host = strings.TrimPrefix(host, "http://")
	}

	return host + ":6379", os.Getenv("UPSTASH_REDIS_REST_TOKEN"), nil
}

func InitRedis() {
	// Check for Upstash first
	upstashURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisURL := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if upstashURL != "" {
		// Use Upstash
		var err error
		redisURL, redisPassword, err = parseUpstashURL(upstashURL)
		if err != nil {
			redisURL = "localhost:6379"
		}
	}

	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	opts := &redis.Options{
		Addr: redisURL,
	}

	if redisPassword != "" {
		opts.Password = redisPassword
	}

	redisClient = redis.NewClient(opts)

	// Test connection with timeout but don't block startup
	go testRedisConnection()
}

func testRedisConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		// Log but continue - workers will retry
	}
}
