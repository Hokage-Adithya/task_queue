
# Railway Deployment Guide

## Quick Deploy to Railway (Easiest - 5 min)

### 1. Create Railway Account
- Go to https://railway.app
- Sign up with GitHub
- Link your GitHub account

### 2. Connect Repository
```bash
cd d:\sk_queue
git init
git add .
git commit -m "Initial commit: Advanced Task Queue"
git remote add origin https://github.com/YOUR_USERNAME/sk_queue.git
git push -u origin main
```

### 3. Deploy on Railway
- Open https://railway.app/dashboard
- Click "New Project"
- Select "Deploy from GitHub"
- Select your repository
- Railway auto-detects Go project
- Configure environment variables:
  - MAILTRAP_HOST=sandbox.smtp.mailtrap.io
  - MAILTRAP_PORT=2525
  - MAILTRAP_USERNAME=a65401e2b86066
  - MAILTRAP_PASSWORD=d9a4c6a188994f
  - MAILTRAP_FROM=noreply@taskqueue.local

### 4. Add Redis
- In Railway dashboard: "Add Service"
- Select "Redis"
- Auto-configures REDIS_URL

### 5. Deploy
- Click "Deploy"
- Wait 2-3 minutes
- Get public URL from Railway dashboard

---

## Docker Hub Deployment

### 1. Build and Push Image
```bash
docker build -t YOUR_USERNAME/task-queue:latest .
docker login
docker push YOUR_USERNAME/task-queue:latest
```

### 2. Run Anywhere
```bash
docker run -p 8080:8080 \
  -e MAILTRAP_HOST=sandbox.smtp.mailtrap.io \
  -e MAILTRAP_PORT=2525 \
  -e MAILTRAP_USERNAME=a65401e2b86066 \
  -e MAILTRAP_PASSWORD=d9a4c6a188994f \
  YOUR_USERNAME/task-queue:latest
```

---

## Heroku Deployment

### 1. Install Heroku CLI
```bash
# Windows
choco install heroku-cli

# Or download from https://devcenter.heroku.com/articles/heroku-cli
```

### 2. Login
```bash
heroku login
```

### 3. Create App
```bash
heroku create your-app-name
```

### 4. Add Redis
```bash
heroku addons:create heroku-redis:premium-0
```

### 5. Set Environment Variables
```bash
heroku config:set MAILTRAP_HOST=sandbox.smtp.mailtrap.io
heroku config:set MAILTRAP_PORT=2525
heroku config:set MAILTRAP_USERNAME=a65401e2b86066
heroku config:set MAILTRAP_PASSWORD=d9a4c6a188994f
```

### 6. Deploy
```bash
git push heroku main
```

---

## AWS Deployment

### 1. Create EC2 Instance
- Ubuntu 22.04 LTS
- t3.small or larger
- Security group: Allow 8080, 443

### 2. SSH and Setup
```bash
ssh -i your-key.pem ubuntu@your-instance.ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### 3. Clone and Deploy
```bash
git clone https://github.com/YOUR_USERNAME/sk_queue.git
cd sk_queue
docker-compose up -d
go build -o queue-server
./queue-server
```

---

## Environment Variables Required

```env
MAILTRAP_HOST=sandbox.smtp.mailtrap.io
MAILTRAP_PORT=2525
MAILTRAP_USERNAME=a65401e2b86066
MAILTRAP_PASSWORD=d9a4c6a188994f
MAILTRAP_FROM=noreply@taskqueue.local
REDIS_URL=localhost:6379
GIN_MODE=release
```

---

## Production Checklist

- [ ] Mailtrap switched to production (if needed)
- [ ] Database persistence configured
- [ ] Redis backup enabled
- [ ] Environment variables set
- [ ] HTTPS/SSL configured
- [ ] Monitoring setup
- [ ] Logging configured
- [ ] Backups scheduled

---

## Testing Deployment

After deployment, test:

```bash
# Health check
curl https://your-deployed-url.com/api/health

# Create task
curl -X POST https://your-deployed-url.com/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"type":"email","payload":"test@example.com"}'

# Check stats
curl https://your-deployed-url.com/api/stats

# Access dashboard
open https://your-deployed-url.com
```

---

## Recommended: Railway

**Why Railway?**
✅ Easiest setup (5 min)
✅ Free tier available
✅ Auto-deploys from GitHub
✅ Redis included
✅ Custom domains
✅ Good for development and production

**Cost**: Free tier sufficient for testing, ~$5/month for small production use
