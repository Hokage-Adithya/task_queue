# Docker Hub Deployment Guide

## Step 1: Create Docker Hub Account

1. Go to https://hub.docker.com
2. Sign up (free account)
3. Verify email
4. Create a repository:
   - Click "Repositories" → "Create repository"
   - Name: `task-queue`
   - Visibility: Public
   - Click "Create"

---

## Step 2: Build Docker Image

Open PowerShell and run:

```powershell
cd d:\sk_queue

# Build the image locally (replace YOUR_DOCKERHUB_USERNAME)
docker build -t YOUR_DOCKERHUB_USERNAME/task-queue:latest .

# Verify the build
docker images | findstr task-queue
```

**Example:**
```powershell
docker build -t johnsmith/task-queue:latest .
```

---

## Step 3: Login to Docker Hub

```powershell
docker login
```

When prompted:
- Username: Your Docker Hub username
- Password: Your Docker Hub password (or access token)

**Tip:** Create an access token for better security:
1. Go to https://hub.docker.com/settings/security
2. Click "New Access Token"
3. Copy the token
4. Use it as your password in `docker login`

---

## Step 4: Push to Docker Hub

```powershell
docker push YOUR_DOCKERHUB_USERNAME/task-queue:latest
```

**Example:**
```powershell
docker push johnsmith/task-queue:latest
```

Wait for upload to complete. You should see:
```
Digest: sha256:abc123...
Status: Pushed
```

---

## Step 5: Verify on Docker Hub

1. Go to https://hub.docker.com
2. Click your username → "Repositories"
3. Click `task-queue`
4. Verify image is listed

---

## Step 6: Run Your Container

### Option A: Local Testing
```powershell
# Start Redis first
docker run -d --name redis -p 6379:6379 redis:7

# Run your app
docker run -p 8080:8080 `
  -e MAILTRAP_HOST=sandbox.smtp.mailtrap.io `
  -e MAILTRAP_PORT=2525 `
  -e MAILTRAP_USERNAME=a65401e2b86066 `
  -e MAILTRAP_PASSWORD=d9a4c6a188994f `
  -e MAILTRAP_FROM=noreply@taskqueue.local `
  -e REDIS_URL=host.docker.internal:6379 `
  YOUR_DOCKERHUB_USERNAME/task-queue:latest
```

Test it:
```powershell
# Health check
curl http://localhost:8080/api/health

# Create task
$body = '{"type":"email","payload":"test@example.com"}'
Invoke-WebRequest -Uri http://localhost:8080/api/tasks `
  -Method POST `
  -Headers @{"Content-Type"="application/json"} `
  -Body $body
```

### Option B: Run on Any Server

```bash
# On your server
docker run -d --name task-queue -p 8080:8080 \
  -e MAILTRAP_HOST=sandbox.smtp.mailtrap.io \
  -e MAILTRAP_PORT=2525 \
  -e MAILTRAP_USERNAME=a65401e2b86066 \
  -e MAILTRAP_PASSWORD=d9a4c6a188994f \
  -e MAILTRAP_FROM=noreply@taskqueue.local \
  -e REDIS_URL=redis-server:6379 \
  YOUR_DOCKERHUB_USERNAME/task-queue:latest
```

---

## Step 7: Using Docker Compose (Recommended)

Create `docker-compose-production.yml`:

```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  task-queue:
    image: YOUR_DOCKERHUB_USERNAME/task-queue:latest
    ports:
      - "8080:8080"
    environment:
      REDIS_URL: redis:6379
      MAILTRAP_HOST: sandbox.smtp.mailtrap.io
      MAILTRAP_PORT: 2525
      MAILTRAP_USERNAME: a65401e2b86066
      MAILTRAP_PASSWORD: d9a4c6a188994f
      MAILTRAP_FROM: noreply@taskqueue.local
      GIN_MODE: release
    depends_on:
      - redis
    restart: unless-stopped

volumes:
  redis_data:
```

Run it:
```bash
docker-compose -f docker-compose-production.yml up -d
```

---

## Popular Hosting Platforms Using Docker Hub

### 🚀 Option 1: Render.com (Easiest)
1. Sign up at https://render.com
2. Connect Docker Hub
3. Deploy from image: `YOUR_DOCKERHUB_USERNAME/task-queue:latest`
4. Set environment variables
5. Deploy

### 🚀 Option 2: Railway (Still Easy)
1. Sign up at https://railway.app
2. Select "Deploy from Docker Hub"
3. Enter: `YOUR_DOCKERHUB_USERNAME/task-queue:latest`
4. Set environment variables
5. Deploy

### 🚀 Option 3: AWS EC2
1. Launch Ubuntu 22.04 instance
2. Install Docker:
   ```bash
   curl -fsSL https://get.docker.com | sh
   ```
3. Run container:
   ```bash
   docker run -d --name task-queue -p 8080:8080 \
     -e MAILTRAP_HOST=... \
     YOUR_DOCKERHUB_USERNAME/task-queue:latest
   ```

### 🚀 Option 4: DigitalOcean App Platform
1. Sign up at https://www.digitalocean.com
2. Create new app
3. Connect Docker Hub
4. Deploy image: `YOUR_DOCKERHUB_USERNAME/task-queue:latest`
5. Configure environment

### 🚀 Option 5: Azure Container Instances
1. Create Azure account
2. Push to Docker Hub (already done)
3. Deploy using Azure CLI or portal

---

## Environment Variables Required

Always set these when running:

```env
MAILTRAP_HOST=sandbox.smtp.mailtrap.io
MAILTRAP_PORT=2525
MAILTRAP_USERNAME=a65401e2b86066
MAILTRAP_PASSWORD=d9a4c6a188994f
MAILTRAP_FROM=noreply@taskqueue.local
REDIS_URL=redis:6379  # or your Redis host
GIN_MODE=release      # production mode
```

---

## Quick Commands

```powershell
# Build locally
docker build -t YOUR_DOCKERHUB_USERNAME/task-queue:latest .

# Login
docker login

# Push to Docker Hub
docker push YOUR_DOCKERHUB_USERNAME/task-queue:latest

# Run with compose
docker-compose -f docker-compose-production.yml up -d

# View logs
docker logs -f task-queue

# Stop container
docker stop task-queue

# Remove container
docker rm task-queue

# Update image (rebuild & push)
docker build -t YOUR_DOCKERHUB_USERNAME/task-queue:latest .
docker push YOUR_DOCKERHUB_USERNAME/task-queue:latest
```

---

## Troubleshooting

### Image fails to run
```powershell
# Check logs
docker logs task-queue

# Verify image
docker images

# Test build locally first
docker run -it --rm YOUR_DOCKERHUB_USERNAME/task-queue:latest /bin/sh
```

### Redis connection fails
- Ensure Redis container is running: `docker ps`
- Use correct REDIS_URL (default: `redis:6379` in compose)
- On Windows local: use `host.docker.internal:6379`

### Port already in use
```powershell
# Find what's using port 8080
netstat -ano | findstr 8080

# Use different port
docker run -p 9000:8080 ...
```

---

## Dashboard Access

After deploying, access at:
- **Local**: http://localhost:8080
- **Production**: https://your-app-url.com

---

## Production Checklist

- [ ] Image pushed to Docker Hub
- [ ] Docker Hub repository is public (or token-protected)
- [ ] Environment variables configured
- [ ] Redis persistence enabled
- [ ] Logs redirected to stdout (check Docker logs)
- [ ] Health check endpoint working
- [ ] HTTPS configured (if using custom domain)
- [ ] Monitoring setup
- [ ] Backup strategy in place

---

## Next Steps

1. **Deploy now**: Use docker-compose-production.yml locally
2. **Test thoroughly**: Create tasks, check emails, monitor stats
3. **Choose platform**: Render, Railway, AWS, DigitalOcean, etc.
4. **Configure domain**: Point custom domain to your deployment
5. **Setup monitoring**: Use platform's built-in monitoring
