package main

import (
	"context"
	"crypto/tls"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func parseUpstashURL(upstashURL string) (host string, password string, err error) {
	// Parse Upstash REST URL: https://joint-swine-8107.upstash.io
	// Upstash uses TCP protocol on port 6379 with TLS
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

	// Upstash default port is 6379 with TLS
	return host + ":6379", os.Getenv("UPSTASH_REDIS_REST_TOKEN"), nil
}

func InitRedis() {
	// Check for Upstash first
	upstashURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisURL := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	useUpstash := false

	if upstashURL != "" {
		// Use Upstash
		var err error
		redisURL, redisPassword, err = parseUpstashURL(upstashURL)
		if err == nil {
			useUpstash = true
		} else {
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

	// Enable TLS for Upstash
	if useUpstash {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
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
