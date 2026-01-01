# Deployment Guide

Panduan lengkap untuk deploy Starter Template Go Fiber menggunakan Docker Compose.

## 📋 Prerequisites

Sebelum deploy, pastikan sudah terinstall:

- **Docker** (v20.10+)
- **Docker Compose** (v2.0+)
- **Git** (untuk clone repository)
- **OpenSSL** (untuk generate SSL certificate)

Cek versi:
```bash
docker --version
docker compose version
git --version
openssl version
```

## 🚀 Quick Start (Development)

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/starter-gofiber.git
cd starter-gofiber
```

### 2. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env file
nano .env
```

### 3. Generate SSL Certificate

```bash
cd assets/certs
openssl genpkey -algorithm RSA -out certificate.pem -pkeyopt rsa_keygen_bits:4096
cd ../..
```

### 4. Run dengan Docker Compose

```bash
# Development mode
docker compose up -d

# Lihat logs
docker compose logs -f app

# Stop
docker compose down
```

Aplikasi akan berjalan di: `http://localhost:3000`

## 🔧 Environment Variables

### Required Variables

```bash
# .env
ENV_TYPE="dev"                    # dev | prod | test

# Database Configuration
DB_HOST="db"                      # Service name di docker-compose
DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="your_secure_password"
DB_NAME="starter_gofiber"

# JWT Configuration
LOCATION_CERT="assets/certs/certificate.pem"
JWT_EXPIRY="15"                   # Access token expiry (minutes)
REFRESH_TOKEN_EXPIRY="10080"      # Refresh token expiry (minutes = 7 days)

# Application
APP_PORT="3000"
APP_NAME="Starter Go Fiber"

# CORS
CORS_ORIGIN="http://localhost:3000,http://localhost:5173"

# Rate Limiting
RATE_LIMIT_MAX="100"              # Max requests
RATE_LIMIT_WINDOW="60"            # Per X seconds
```

### Optional Variables

```bash
# Email Service (for production)
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USER="your-email@gmail.com"
SMTP_PASSWORD="your-app-password"

# File Upload
MAX_UPLOAD_SIZE="10485760"        # 10MB in bytes
UPLOAD_PATH="./public/uploads"

# Logging
LOG_LEVEL="info"                  # debug | info | warn | error
```

## 🐳 Docker Deployment

### Development Mode

File: `assets/docker/docker-compose.yml`

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: assets/docker/Dockerfile
    ports:
      - "3000:3000"
    volumes:
      - .:/app                    # Hot reload enabled
      - /app/tmp                  # Exclude tmp directory
    environment:
      - ENV_TYPE=dev
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: starter_gofiber
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

**Commands**:
```bash
# Start
docker compose up -d

# Rebuild
docker compose up -d --build

# View logs
docker compose logs -f

# Stop
docker compose down

# Stop and remove volumes
docker compose down -v
```

### Production Mode

File: `assets/docker/prod/docker-compose.yml`

```yaml
version: '3.8'

services:
  app:
    build:
      context: ../../..
      dockerfile: assets/docker/prod/Dockerfile
    ports:
      - "3000:3000"
    environment:
      - ENV_TYPE=prod
    env_file:
      - ../../../.env
    depends_on:
      - db
    restart: always
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backup:/backup          # Backup directory
    restart: always
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
      - ../../../public:/usr/share/nginx/html/public:ro
    depends_on:
      - app
    restart: always
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  postgres_data:
```

**Production Commands**:
```bash
# Update .env for production
nano .env
# Set: ENV_TYPE="prod"

# Copy production configs
cp assets/docker/prod/docker-compose.yml docker-compose.prod.yml
cp assets/docker/prod/Dockerfile Dockerfile.prod

# Build and run
docker compose -f docker-compose.prod.yml up -d --build

# Monitor
docker compose -f docker-compose.prod.yml logs -f

# Stop
docker compose -f docker-compose.prod.yml down
```

## 📦 Dockerfile Configurations

### Development Dockerfile

File: `assets/docker/Dockerfile`

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Run with air
CMD ["air", "-c", ".air.toml"]
```

### Production Dockerfile

File: `assets/docker/prod/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates curl

# Copy binary from builder
COPY --from=builder /app/main .

# Copy necessary files
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/.env .env

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# Run
CMD ["./main"]
```

## 🌐 Nginx Configuration (Production)

File: `assets/docker/prod/nginx.conf`

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api {
        server app:3000;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

    server {
        listen 80;
        server_name yourdomain.com www.yourdomain.com;

        # Redirect to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name yourdomain.com www.yourdomain.com;

        # SSL Configuration
        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers on;

        # Security Headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

        # Gzip compression
        gzip on;
        gzip_vary on;
        gzip_min_length 1024;
        gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

        # Static files
        location /public {
            alias /usr/share/nginx/html/public;
            expires 30d;
            add_header Cache-Control "public, immutable";
        }

        # API Proxy
        location /api {
            limit_req zone=api_limit burst=20 nodelay;

            proxy_pass http://api;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Health check
        location /health {
            proxy_pass http://api;
            access_log off;
        }
    }
}
```

## 🔐 SSL Certificate Setup

### Let's Encrypt (Free SSL)

```bash
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Certificates will be in:
# /etc/letsencrypt/live/yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/yourdomain.com/privkey.pem

# Copy to project
sudo cp /etc/letsencrypt/live/yourdomain.com/*.pem assets/docker/prod/ssl/

# Auto-renewal (add to crontab)
0 0 1 * * certbot renew --quiet
```

### Self-Signed SSL (Development)

```bash
cd assets/docker/prod/ssl

# Generate private key and certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout privkey.pem \
  -out fullchain.pem \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

## 📊 Database Management

### Database Migration

```bash
# Access database container
docker compose exec db psql -U postgres -d starter_gofiber

# Or use external tool
psql -h localhost -p 5432 -U postgres -d starter_gofiber
```

### Backup Database

```bash
# Backup
docker compose exec db pg_dump -U postgres starter_gofiber > backup_$(date +%Y%m%d).sql

# Restore
docker compose exec -T db psql -U postgres starter_gofiber < backup_20251231.sql
```

### Automated Backup Script

File: `scripts/backup-db.sh`

```bash
#!/bin/bash

BACKUP_DIR="./backup"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="backup_${DATE}.sql"

mkdir -p $BACKUP_DIR

docker compose exec -T db pg_dump -U postgres starter_gofiber > $BACKUP_DIR/$FILENAME

# Keep only last 7 backups
ls -t $BACKUP_DIR/backup_*.sql | tail -n +8 | xargs -r rm

echo "Backup created: $BACKUP_DIR/$FILENAME"
```

Add to crontab:
```bash
# Daily backup at 2 AM
0 2 * * * /path/to/project/scripts/backup-db.sh
```

## 🔄 CI/CD Integration

### GitHub Actions

File: `.github/workflows/deploy.yml`

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Copy files to server
        uses: appleboy/scp-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          source: "."
          target: "/var/www/starter-gofiber"
      
      - name: Deploy with Docker Compose
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /var/www/starter-gofiber
            docker compose -f docker-compose.prod.yml down
            docker compose -f docker-compose.prod.yml up -d --build
```

## 🖥️ VPS Deployment

### Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin

# Create deployment directory
sudo mkdir -p /var/www/starter-gofiber
sudo chown $USER:$USER /var/www/starter-gofiber
```

### Deploy Application

```bash
# Clone repository
cd /var/www/starter-gofiber
git clone https://github.com/yourusername/starter-gofiber.git .

# Setup environment
cp .env.example .env
nano .env  # Edit configuration

# Generate SSL certificate
cd assets/certs
openssl genpkey -algorithm RSA -out certificate.pem -pkeyopt rsa_keygen_bits:4096
cd ../..

# Run production
docker compose -f docker-compose.prod.yml up -d --build
```

### Firewall Configuration

```bash
# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow SSH (jika belum)
sudo ufw allow 22/tcp

# Enable firewall
sudo ufw enable
```

## 📈 Monitoring & Logging

### View Logs

```bash
# All logs
docker compose logs -f

# Specific service
docker compose logs -f app

# Last 100 lines
docker compose logs --tail=100 app

# With timestamps
docker compose logs -f -t app
```

### Health Check

```bash
# Check container health
docker compose ps

# Check application health endpoint
curl http://localhost:3000/health
```

### Resource Monitoring

```bash
# Container stats
docker stats

# Disk usage
docker system df

# Clean unused resources
docker system prune -a
```

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs app

# Check configuration
docker compose config

# Rebuild without cache
docker compose build --no-cache
docker compose up -d
```

### Database Connection Failed

```bash
# Check database is running
docker compose ps db

# Test connection
docker compose exec db psql -U postgres -d starter_gofiber

# Check environment variables
docker compose exec app env | grep DB_
```

### Port Already in Use

```bash
# Find process using port 3000
sudo lsof -i :3000

# Kill process
sudo kill -9 <PID>

# Or change port in .env
APP_PORT=3001
```

## 🔄 Update & Rollback

### Update Application

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.prod.yml up -d --build

# Or with zero-downtime
docker compose -f docker-compose.prod.yml up -d --no-deps --build app
```

### Rollback

```bash
# Revert to previous commit
git log --oneline
git checkout <commit-hash>

# Rebuild
docker compose -f docker-compose.prod.yml up -d --build
```

## 📚 Best Practices

### ✅ DO

1. **Use Environment Variables** - Jangan hardcode credentials
2. **Enable HTTPS** - Gunakan SSL certificate di production
3. **Regular Backups** - Backup database secara berkala
4. **Monitor Logs** - Setup log monitoring dan alerting
5. **Update Dependencies** - Keep Docker images dan dependencies updated
6. **Use Health Checks** - Implement health check endpoints
7. **Limit Resources** - Set memory dan CPU limits di docker-compose

### ❌ DON'T

1. **Expose Database Port** - Jangan expose port 5432 di production
2. **Use Default Passwords** - Ganti semua default passwords
3. **Run as Root** - Jangan run container sebagai root user
4. **Ignore Logs** - Jangan abaikan error logs
5. **Skip Backups** - Jangan lupa backup sebelum major updates

---

**Last Updated**: December 31, 2025
