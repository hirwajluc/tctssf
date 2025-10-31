# TCTSSF Quick Start Guide

## Improvements Implemented ✅

All 9 security and infrastructure improvements have been successfully implemented:

1. ✅ **Redis-backed session management** - Sessions persist across restarts
2. ✅ **Environment variables (.env)** - No hardcoded credentials
3. ✅ **Restricted CORS** - Configurable allowed origins
4. ✅ **HTTPS/TLS support** - Optional SSL/TLS encryption
5. ✅ **Structured logging (Zap)** - Production-ready logging
6. ✅ **Unit tests** - Test coverage for core functionality
7. ✅ **Swagger/OpenAPI docs** - Interactive API documentation
8. ✅ **Password policy** - Strong password requirements
9. ✅ **Rate limiting** - 100 requests/minute protection

## Quick Start

### 1. Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` with your settings:
```env
# Minimal required configuration
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=tctssf

# Optional: Redis for session persistence
REDIS_URL=redis://localhost:6379/0

# Optional: Restrict CORS
ALLOWED_ORIGINS=http://localhost:3000
```

### 2. Install Dependencies

#### MySQL (Required)
```bash
# Ubuntu/Debian
sudo apt-get install mysql-server

# macOS
brew install mysql
```

#### Redis (Optional - recommended for production)
```bash
# Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis-server

# macOS
brew install redis
brew services start redis
```

### 3. Build and Run

```bash
# Install Go dependencies
go mod download

# Build the application
go build -o tctssf main.go

# Run the application
./tctssf
```

Or run directly:
```bash
go run main.go
```

### 4. Generate Swagger Documentation (First Time Only)

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin

# Generate Swagger docs
swag init
```

### 5. Access the Application

- **Frontend**: http://localhost:3000
- **API Docs**: http://localhost:3000/swagger/index.html
- **API Base**: http://localhost:3000/api

### 6. Default Credentials

```
Superadmin:
  Email: superadmin@tctssf.rw
  Password: admin123

Admin:
  Email: admin@tctssf.rw
  Password: admin123

Treasurer:
  Email: treasurer@tctssf.rw
  Password: treasurer123
```

**⚠️ Change these passwords immediately in production!**

## Testing

Run tests:
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./utils -v
```

## Production Deployment

### 1. Enable TLS

Generate certificates:
```bash
# Create directory
mkdir -p certs

# Self-signed (development)
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key \
  -out certs/server.crt -days 365 -nodes

# Production: Use Let's Encrypt
sudo certbot certonly --standalone -d yourdomain.com
```

Update `.env`:
```env
ENABLE_TLS=true
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key
ENVIRONMENT=production
```

### 2. Configure Redis

For production, always use Redis:
```env
REDIS_URL=redis://localhost:6379/0
```

### 3. Secure Database

```env
DB_PASSWORD=<strong-password>
DB_HOST=127.0.0.1  # Localhost only
```

### 4. Restrict CORS

```env
ALLOWED_ORIGINS=https://yourdomain.com
```

### 5. Run as Service

Create systemd service `/etc/systemd/system/tctssf.service`:
```ini
[Unit]
Description=TCTSSF Server
After=network.target mysql.service redis.service

[Service]
Type=simple
User=tctssf
WorkingDirectory=/opt/tctssf
ExecStart=/opt/tctssf/tctssf
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable tctssf
sudo systemctl start tctssf
sudo systemctl status tctssf
```

## Features Overview

### New Security Features

1. **Session Persistence**: Sessions stored in Redis survive restarts
2. **Environment Config**: All secrets in `.env` file (not in code)
3. **CORS Protection**: Only allowed origins can access API
4. **TLS Encryption**: HTTPS support for production
5. **Structured Logs**: JSON logs for production monitoring
6. **Password Policy**: 8+ chars, upper/lower/number/special required
7. **Rate Limiting**: Anti-brute force (100 req/min)

### API Documentation

Access Swagger UI at `/swagger/index.html` for:
- Complete API reference
- Interactive testing
- Request/response examples
- Authentication guide

### Password Requirements

New passwords must have:
- ✅ Minimum 8 characters
- ✅ At least one uppercase letter
- ✅ At least one lowercase letter
- ✅ At least one number
- ✅ At least one special character (!@#$%^&*)
- ❌ Not a common weak password

Generate secure passwords:
```go
password, err := utils.GenerateSecurePassword(12)
```

## Monitoring

### Logs

Logs are structured (JSON in production):
```bash
# View logs
tail -f /var/log/tctssf/app.log

# or with systemd
journalctl -u tctssf -f
```

### Health Checks

```bash
# Check server is running
curl http://localhost:3000/api/login

# Check Redis connection (appears in startup logs)
# Check database connection (appears in startup logs)
```

## Troubleshooting

### Redis not available
- App will fall back to in-memory sessions
- Warning logged on startup
- Sessions won't persist across restarts

### Database connection failed
- Check MySQL is running: `sudo systemctl status mysql`
- Verify credentials in `.env`
- Ensure database exists: `mysql -u root -p -e "CREATE DATABASE tctssf;"`

### Port already in use
- Change `SERVER_PORT` in `.env`
- Kill process using port: `sudo lsof -ti:3000 | xargs kill`

### Tests failing
- Database tests need MySQL connection
- Other tests should pass without database
- Run: `go test ./middleware ./models ./utils -v`

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 3000 | HTTP server port |
| `SERVER_HOST` | localhost | Server hostname |
| `DB_HOST` | localhost | MySQL hostname |
| `DB_PORT` | 3306 | MySQL port |
| `DB_USER` | root | MySQL username |
| `DB_PASSWORD` | (empty) | MySQL password |
| `DB_NAME` | tctssf | Database name |
| `REDIS_URL` | redis://localhost:6379/0 | Redis connection URL |
| `ALLOWED_ORIGINS` | * | CORS allowed origins |
| `ENABLE_TLS` | false | Enable HTTPS |
| `TLS_CERT_FILE` | ./certs/server.crt | TLS certificate path |
| `TLS_KEY_FILE` | ./certs/server.key | TLS private key path |
| `SESSION_EXPIRATION_HOURS` | 24 | Session lifetime |
| `ENVIRONMENT` | development | Environment (development/production) |

## Support

### Documentation
- Full improvements: [IMPROVEMENTS.md](IMPROVEMENTS.md)
- API docs: http://localhost:3000/swagger/index.html

### Logs
Check logs for errors:
- Development: Console output
- Production: JSON logs

### Common Issues
1. **"No .env file found"**: Create `.env` from `.env.example`
2. **"Failed to connect to database"**: Check MySQL is running and credentials
3. **"Redis not connected"**: Optional - app works without Redis
4. **Rate limited**: Wait 1 minute or reduce request rate

---

**Version**: 1.0.0 (Improved)
**Last Updated**: 2025-10-31
