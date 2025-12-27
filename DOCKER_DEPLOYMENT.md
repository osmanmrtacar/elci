# Docker Deployment Guide

## Quick Start

### 1. Using Docker Compose (Recommended)

```bash
# Create .env file from example
cp .env.docker.example .env

# Edit .env and add your API credentials
nano .env

# Start both frontend and backend
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

**Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health checks:
  - Frontend: http://localhost:3000/health
  - Backend: http://localhost:8080/health

---

## Individual Container Deployment

### Backend Only

```bash
# Build
docker build -t sosyal-backend ./backend

# Run
docker run -d \
  --name sosyal-backend \
  -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e TIKTOK_CLIENT_KEY=your_key \
  -e TIKTOK_CLIENT_SECRET=your_secret \
  -e X_CLIENT_ID=your_x_id \
  -e X_CLIENT_SECRET=your_x_secret \
  -e INSTAGRAM_APP_ID=your_instagram_id \
  -e INSTAGRAM_APP_SECRET=your_instagram_secret \
  -e JWT_SECRET=your_jwt_secret_min_32_chars \
  -v sosyal-data:/app/data \
  sosyal-backend
```

### Frontend Only

```bash
# Build (with API URL)
docker build \
  --build-arg VITE_API_BASE_URL=https://yourdomain.com \
  -t sosyal-frontend \
  ./frontend

# Run
docker run -d \
  --name sosyal-frontend \
  -p 3000:80 \
  sosyal-frontend
```

---

## VPS Deployment with Cloudflare

### Prerequisites
- VPS with Docker installed
- Domain name (e.g., `yourdomain.com`)
- Cloudflare account with domain added

### Step 1: Update DNS Records
In Cloudflare dashboard:
- Add `A` record: `@` → Your VPS IP (Orange cloud enabled for DDoS protection)
- Add `A` record: `api` → Your VPS IP (Orange cloud enabled)

### Step 2: Clone & Configure

```bash
# SSH into your VPS
ssh user@your-vps-ip

# Clone repository
git clone <your-repo-url>
cd sosyal

# Create environment file
cp .env.docker.example .env

# Edit with your production values
nano .env
```

**Important environment variables:**
```bash
# Use your domain
TIKTOK_REDIRECT_URI=https://yourdomain.com/api/v1/auth/tiktok/callback
X_REDIRECT_URI=https://yourdomain.com/api/v1/auth/x/callback
INSTAGRAM_REDIRECT_URI=https://yourdomain.com/api/v1/auth/instagram/callback
VITE_API_BASE_URL=https://api.yourdomain.com

# CORS - add your domain
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://api.yourdomain.com

# Strong JWT secret (generate with: openssl rand -base64 32)
JWT_SECRET=<your-secure-secret>

# Set to production
ENVIRONMENT=production
```

### Step 3: Setup Nginx Reverse Proxy (on VPS)

```bash
# Install nginx
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx

# Create nginx config
sudo nano /etc/nginx/sites-available/sosyal
```

**Nginx configuration:**
```nginx
# Frontend - yourdomain.com
server {
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}

# Backend - api.yourdomain.com
server {
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

**Enable site:**
```bash
sudo ln -s /etc/nginx/sites-available/sosyal /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 4: Get SSL Certificates

```bash
# Get certificates for both domains
sudo certbot --nginx -d yourdomain.com -d api.yourdomain.com

# Auto-renewal is set up automatically
```

### Step 5: Start Application

```bash
# Start with docker-compose
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 6: Update OAuth Redirect URIs

Update your OAuth app settings:
- **TikTok Developer Portal**: Update redirect URI to `https://yourdomain.com/api/v1/auth/tiktok/callback`
- **X Developer Portal**: Update redirect URI to `https://yourdomain.com/api/v1/auth/x/callback`
- **Facebook/Instagram Developer**: Update redirect URI to `https://yourdomain.com/api/v1/auth/instagram/callback`

---

## Development with Docker

For development with hot reload:

### Frontend (with Vite dev server)
```bash
cd frontend

# Run dev server (not Docker)
npm run dev

# Access at http://localhost:3000
```

### Backend
```bash
cd backend

# Run locally (not Docker)
go run ./cmd/server

# Or use Docker with live reload (requires air or similar)
```

---

## Updating the Application

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## Monitoring & Maintenance

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Check Health
```bash
# Backend
curl http://localhost:8080/health

# Frontend
curl http://localhost:3000/health
```

### Database Backup
```bash
# Backup database
docker cp sosyal-backend:/app/data/sosyal.db ./backup-$(date +%Y%m%d).db

# Restore database
docker cp ./backup-20231227.db sosyal-backend:/app/data/sosyal.db
docker-compose restart backend
```

### Cleanup
```bash
# Remove stopped containers
docker-compose down

# Remove with volumes (WARNING: deletes database!)
docker-compose down -v

# Clean up unused images
docker image prune -a
```

---

## Troubleshooting

### Container won't start
```bash
# Check logs
docker-compose logs backend
docker-compose logs frontend

# Check if ports are in use
sudo lsof -i :8080
sudo lsof -i :3000
```

### Database issues
```bash
# Enter container
docker exec -it sosyal-backend sh

# Check database
ls -la /app/data/
```

### Rebuild from scratch
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

---

## Security Checklist

- [ ] Change default JWT_SECRET
- [ ] Use strong, unique API secrets
- [ ] Enable Cloudflare proxy (orange cloud)
- [ ] Enable SSL/TLS (certbot)
- [ ] Set CORS_ALLOWED_ORIGINS to your domain only
- [ ] Set ENVIRONMENT=production
- [ ] Regular database backups
- [ ] Monitor logs for suspicious activity
- [ ] Keep Docker images updated

---

## Performance Optimization

### For Production:
1. **Cloudflare Settings:**
   - Enable caching
   - Enable Brotli compression
   - Enable HTTP/2 and HTTP/3
   - Set up rate limiting

2. **Docker Resources:**
   ```yaml
   # Add to docker-compose.yml services
   deploy:
     resources:
       limits:
         cpus: '0.5'
         memory: 512M
       reservations:
         cpus: '0.25'
         memory: 256M
   ```

3. **Database:**
   - Regular VACUUM
   - Monitor size
   - Consider PostgreSQL for heavy load

---

## Support

For issues or questions, check:
- Application logs: `docker-compose logs`
- Health endpoints: `/health`
- GitHub issues (if applicable)
