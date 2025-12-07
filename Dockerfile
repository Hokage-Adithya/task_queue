# Multi-stage Dockerfile
FROM golang:1.25.5 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o queue-server .

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/queue-server .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./queue-server"]
