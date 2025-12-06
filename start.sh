#!/bin/bash

# Production startup script

# Ensure Redis is running
if ! command -v redis-cli &> /dev/null; then
    echo "Starting Redis..."
    redis-server &
fi

# Wait for Redis
sleep 2

# Build if needed
if [ ! -f queue-server ]; then
    echo "Building application..."
    go build -o queue-server .
fi

# Set production environment
export GIN_MODE=release
export REDIS_URL=${REDIS_URL:-localhost:6379}

# Start the server
echo "Starting Task Queue Server..."
./queue-server
