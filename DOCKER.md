# Docker Deployment Guide

This guide explains how to run the Cloud Storage application using Docker.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.8+
- At least 2GB RAM
- At least 5GB free disk space

## Quick Start

### Option 1: Start All Services (Recommended)

```bash
# Start backend and frontend
make docker-up

# Or start with full Nginx proxy
make docker-full
```

Access the application at:
- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080
- **Nginx Proxy** (full mode): http://localhost (port 80)

### Option 2: Start Backend Only

```bash
make docker-backend
```

This starts:
- PostgreSQL database (port 5432)
- Redis cache (port 6379)
- Go backend application (port 8080)

### Option 3: Start Frontend Only

```bash
make docker-frontend
```

This starts:
- React frontend application (port 80)
- Requires backend to be running separately

## Environment Variables

The following environment variables are configured in `docker-compose.yml`:

### Backend (app)

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `production` | Application environment |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `SERVER_PORT` | `8080` | Server port |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `cloud_storage` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `password` | Database password |
| `REDIS_HOST` | `redis` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | `redispassword` | Redis password |
| `JWT_SECRET` | `your-secret-key-change-this-in-production` | JWT secret key |
| `STORAGE_PATH` | `/app/storage/uploads` | Upload storage path |
| `MAX_UPLOAD_SIZE` | `104857600` | Max upload size in bytes (100MB) |

### Frontend (web)

Environment variables are configured in `web/.env`:
- `VITE_API_BASE_URL`: Backend API URL (default: `http://localhost:8080/api/v1`)

## Docker Profiles

The project uses Docker Compose profiles to control which services to start:

| Profile | Services |
|---------|----------|
| `backend` | postgres, redis, app, adminer, redis-commander |
| `frontend` | web |
| `full` | nginx (proxy) |

## Available Commands

```bash
# Start all services
make docker-up

# Start backend only
make docker-backend

# Start frontend only
make docker-frontend

# Start with full Nginx proxy
make docker-full

# Stop all services
make docker-down

# Stop and remove volumes
make docker-clean

# Rebuild frontend
docker-compose build --no-cache web

# View logs
make logs              # Backend logs
make logs-db           # Database logs
make logs-redis        # Redis logs

# Access database shell
make db-shell

# Access Redis shell
make redis-shell

# Access application shell
make app-shell

# Backup database
make db-backup

# Restore database
make db-restore
```

## Development Workflow

### Running in Development Mode

For development with hot-reload:

**Backend (Go):**
```bash
cd cmd/server
go run .
```

**Frontend (React):**
```bash
cd web
npm install
npm run dev
```

### Building for Production

**Backend:**
```bash
make docker-build
```

**Frontend:**
```bash
cd web
npm run build
```

## Data Persistence

The following Docker volumes are created for data persistence:

| Volume | Purpose |
|---------|---------|
| `postgres_data` | PostgreSQL database data |
| `redis_data` | Redis cache data |
| `uploads_data` | User uploaded files |
| `temp_data` | Temporary upload files |
| `logs_data` | Application logs |

## Health Checks

All services include health checks:

- **PostgreSQL**: Every 10s, checks `pg_isready`
- **Redis**: Every 10s, checks `redis-cli ping`
- **App**: Every 30s, checks `/health` endpoint

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:
```bash
# Change port in docker-compose.yml
# Find and modify: "8080:8080" to "9090:8080"
```

### Database Connection Issues

```bash
# Check database logs
make logs-db

# Enter database shell
make db-shell
```

### Redis Connection Issues

```bash
# Check Redis logs
make logs-redis

# Enter Redis shell
make redis-shell
```

### Rebuild Services

```bash
# Stop and rebuild backend
docker-compose down
docker-compose build --no-cache app
docker-compose up -d

# Rebuild frontend
docker-compose build --no-cache web
docker-compose up -d web
```

### Clean Everything

```bash
# Remove all containers, volumes, and images
make docker-clean
docker system prune -a --volumes
```

## Production Deployment

For production deployment:

1. **Change passwords** in `docker-compose.yml`:
   - `POSTGRES_PASSWORD`
   - `REDIS_PASSWORD`
   - `JWT_SECRET`

2. **Configure storage**:
   - Adjust `STORAGE_PATH` if using external storage
   - Ensure sufficient disk space

3. **Set up SSL**:
   - Use Nginx with Let's Encrypt for HTTPS
   - Update `CORS_ALLOW_ORIGINS` to your domain

4. **Database backups**:
   ```bash
   # Set up automated backups
   crontab -e "0 2 * * * docker-compose exec postgres pg_dump -U postgres cloud_storage > /backup/db_$(date +\%Y\%m\%d).sql"
   ```

## Monitoring

### Access Adminer (Database GUI)
- URL: http://localhost:8081
- Server: postgres
- Username: postgres
- Password: password

### Access Redis Commander (Redis GUI)
- URL: http://localhost:8082
- Host: local:redis:6379
- Password: redispassword

## Security Considerations

1. **Change default passwords** in production
2. **Use secrets management** for sensitive data
3. **Enable firewall rules** to restrict access
4. **Set up SSL/TLS** for production
5. **Regular updates** of Docker images
6. **Monitor logs** for suspicious activity

## Architecture

```
┌─────────────┐
│   Nginx     │ (Optional - Port 80/443)
│  (Proxy)    │
└──────┬──────┘
       │
       ├──────────────────┐
       │                  │
┌──────▼──────┐   ┌──────▼──────┐
│   Frontend   │   │   Backend    │
│  (React)    │   │    (Go)      │
│   Port 80   │   │  Port 8080   │
└──────────────┘   └──────┬──────┘
                         │
                ┌────────┴────────┐
                │                 │
         ┌──────▼──────┐   ┌───▼──────┐
         │  PostgreSQL  │   │  Redis   │
         │  Port 5432   │   │ Port 6379 │
         └───────────────┘   └──────────┘
```

## Support

For issues or questions:
- Check logs: `make logs`
- Review docker-compose.yml configuration
- Verify network connectivity
- Check Docker version compatibility
