# RDB SLS - RDB Staff Loan System

A comprehensive savings and loan management system built with Go and Fiber framework.

## 🎉 Status: Production Ready!

All 9 security and infrastructure improvements have been implemented and tested.

## ✨ Key Features

- **3-Stage Loan Approval** - Treasurer → Vice President → President
- **Salary Deduction Integration** - CSV generation and processing for HR
- **Member Management** - Complete member lifecycle management
- **Savings Tracking** - Monthly commitments and balances
- **Social Contributions** - Fixed monthly social fund contributions
- **Role-Based Access Control** - Member, Admin, Superadmin, Treasurer roles

## 🔒 Security Features (NEW!)

1. ✅ **Redis Session Management** - Persistent sessions
2. ✅ **Environment Configuration** - No hardcoded secrets
3. ✅ **CORS Protection** - Configurable allowed origins
4. ✅ **HTTPS/TLS Support** - Production-ready encryption
5. ✅ **Structured Logging** - Zap logger for production
6. ✅ **Password Policies** - Strong password requirements
7. ✅ **Rate Limiting** - 100 req/min anti-brute force
8. ✅ **API Documentation** - Swagger/OpenAPI integration
9. ✅ **Unit Tests** - Comprehensive test coverage

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- MySQL 5.7+
- Redis (optional, recommended)

### Installation

```bash
# 1. Clone or navigate to project
cd /root/go-projects/tctssf

# 2. Copy environment file
cp .env.example .env

# 3. Configure database (edit .env)
nano .env

# 4. Generate Swagger docs (first time only)
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
swag init

# 5. Build
go build -o tctssf main.go

# 6. Run
./tctssf
```

### Access

- **Frontend**: http://localhost:3000
- **API Docs**: http://localhost:3000/swagger/index.html
- **API**: http://localhost:3000/api

### Default Credentials

```
Superadmin: superadmin@rdbsls.rw / admin123
Admin:      admin@rdbsls.rw      / admin123
Treasurer:  treasurer@rdbsls.rw  / treasurer123
```

**⚠️ Change these in production!**

## 📚 Documentation

### Quick References

| Document | Purpose |
|----------|---------|
| [QUICK_START.md](QUICK_START.md) | 5-minute setup guide |
| [IMPROVEMENTS.md](IMPROVEMENTS.md) | All improvements detailed |
| [FINAL_STATUS.md](FINAL_STATUS.md) | Complete status report |

### API Documentation

| Document | Purpose |
|----------|---------|
| [SWAGGER_SETUP.md](SWAGGER_SETUP.md) | Swagger documentation guide |
| [SWAGGER_ROUTE_FIX.md](SWAGGER_ROUTE_FIX.md) | Routing fix explanation |

### Interactive Docs

Access Swagger UI for interactive API testing:
```
http://localhost:3000/swagger/index.html
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verify Swagger
./verify_swagger.sh
```

## 🏗️ Project Structure

```
tctssf/
├── main.go                 # Application entry point
├── config/                 # Configuration
│   ├── database.go        # Database connection
│   ├── env.go             # Environment loader
│   ├── redis.go           # Redis session store
│   └── logger.go          # Zap logger
├── controllers/           # Request handlers
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── loan_controller.go
│   ├── admin_controller.go
│   └── treasurer_controller.go
├── middleware/            # HTTP middleware
│   └── auth.go           # Authentication & authorization
├── models/                # Data models
│   ├── user.go
│   ├── loan.go
│   ├── auth.go
│   └── treasurer.go
├── routes/                # Route definitions
│   └── routes.go
├── utils/                 # Utilities
│   └── password.go       # Password validation & generation
├── docs/                  # Generated Swagger docs
├── frontend/              # Frontend HTML/CSS/JS
└── tests/                 # Test files
```

## 🔧 Configuration

### Environment Variables

```env
# Server
SERVER_PORT=3000
SERVER_HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=tctssf

# Redis (optional)
REDIS_URL=redis://localhost:6379/0

# Security
ALLOWED_ORIGINS=http://localhost:3000
ENABLE_TLS=false
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key

# Environment
ENVIRONMENT=development
```

See [`.env.example`](.env.example) for full configuration options.

## 🚢 Production Deployment

### 1. Enable TLS

```bash
# Generate certificate
mkdir -p certs
openssl req -x509 -newkey rsa:4096 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -days 365 -nodes
```

Update `.env`:
```env
ENABLE_TLS=true
```

### 2. Configure Redis

```env
REDIS_URL=redis://localhost:6379/0
```

### 3. Restrict CORS

```env
ALLOWED_ORIGINS=https://yourdomain.com
```

### 4. Set Production Mode

```env
ENVIRONMENT=production
```

### 5. Run as Service

See [QUICK_START.md](QUICK_START.md) for systemd service setup.

## 📊 Database Schema

- **users** - User accounts and roles
- **savings_accounts** - Member savings tracking
- **transactions** - All financial transactions
- **loans** - Loan applications and approvals
- **loan_repayments** - Repayment schedules
- **salary_deduction_lists** - Monthly deduction batches
- **salary_deduction_items** - Individual deductions

See [tctssf_db.sql](tctssf_db.sql) for complete schema.

## 🔐 Security

### Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character
- Not a common weak password

### Rate Limiting

- 100 requests per minute per IP
- Automatic throttling on exceeded limit

### Session Management

- Redis-backed sessions (persistent)
- Fallback to in-memory storage
- Secure token generation

## 🛠️ Development

### Adding New Endpoints

1. Create controller method with Swagger annotations:

```go
// GetProfile retrieves user profile
// @Summary Get user profile
// @Tags User
// @Security BearerAuth
// @Success 200 {object} models.User
// @Router /profile [get]
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
    // Implementation
}
```

2. Register route in `routes/routes.go`

3. Regenerate Swagger docs:

```bash
swag init
go build -o tctssf main.go
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./middleware -v

# With race detection
go test -race ./...
```

## 📝 API Endpoints

### Public
- POST `/api/login` - User authentication

### Protected (Require Authentication)
- GET `/api/profile` - Get user profile
- GET `/api/dashboard` - Get dashboard data
- POST `/api/loans/apply` - Apply for loan
- GET `/api/loans` - Get user loans

### Admin Only
- POST `/api/admin/members` - Create member
- GET `/api/admin/members` - List members
- GET `/api/admin/loans/pending` - Pending loans

### Treasurer Only
- GET `/api/treasurer/dashboard` - Treasurer dashboard
- POST `/api/treasurer/salary-deductions/generate` - Generate deductions
- GET `/api/treasurer/salary-deductions/lists` - List deduction lists

See Swagger UI for complete API reference.

## 🤝 Contributing

1. Follow existing code structure
2. Add tests for new features
3. Update Swagger documentation
4. Run tests before committing

## 📄 License

MIT License

## 🆘 Support

### Documentation
- [QUICK_START.md](QUICK_START.md) - Quick setup
- [IMPROVEMENTS.md](IMPROVEMENTS.md) - Feature details
- [SWAGGER_SETUP.md](SWAGGER_SETUP.md) - API docs guide

### Common Issues

**Swagger not loading?**
- Run: `swag init && go build -o tctssf main.go`
- Check: Route registration order in `main.go`

**Redis not connecting?**
- Application continues with in-memory sessions
- Check: `REDIS_URL` in `.env`

**Database connection failed?**
- Verify: MySQL is running
- Check: Credentials in `.env`

## 🎯 Key Achievements

- ✅ 9/9 security improvements implemented
- ✅ Production-ready with HTTPS support
- ✅ Comprehensive API documentation
- ✅ Test coverage for core functionality
- ✅ Zero breaking changes from original
- ✅ Fully backward compatible

---

**Version**: 1.0.0 (Improved)
**Last Updated**: 2025-10-31
**Status**: Production Ready ✅
