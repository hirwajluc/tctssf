# Swagger API Documentation Setup

## ✅ Swagger is Now Configured!

The Swagger/OpenAPI documentation has been successfully set up for the RDB SLS API.

## Accessing Swagger UI

### 1. Start the Server

```bash
# Using the binary
./tctssf

# Or run with Go
go run main.go
```

### 2. Access Swagger UI

Open your browser and navigate to:

```
http://localhost:3000/swagger/index.html
```

### Alternative Endpoints

- **Swagger JSON**: `http://localhost:3000/swagger/doc.json`
- **Swagger YAML**: Available in `docs/swagger.yaml`
- **Raw JSON**: `docs/swagger.json`

## What's Documented

Currently documented endpoints:

### Authentication
- **POST** `/api/login` - User authentication

More endpoints can be added by following the same pattern.

## How to Add More API Documentation

### 1. Add Swagger Annotations to Controllers

Example from `controllers/auth_controller.go`:

```go
// Login handles user authentication
// @Summary User login
// @Description Authenticate user and receive session token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
    // ... implementation
}
```

### 2. Regenerate Swagger Docs

After adding annotations, run:

```bash
swag init
```

This will update:
- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

### 3. Rebuild and Run

```bash
go build -o tctssf main.go
./tctssf
```

## Swagger Annotation Reference

### Common Annotations

```go
// @Summary      Short description
// @Description  Detailed description
// @Tags         Category name
// @Accept       json
// @Produce      json
// @Param        name  location  type  required  description
// @Success      200   {object}  ModelName
// @Failure      400   {object}  map[string]string
// @Router       /path [method]
// @Security     BearerAuth
```

### Parameter Locations

- `path` - URL path parameter (e.g., `/users/{id}`)
- `query` - Query string parameter (e.g., `?page=1`)
- `body` - Request body
- `header` - HTTP header

### Example for Protected Endpoint

```go
// GetProfile retrieves user profile
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /profile [get]
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
    // ... implementation
}
```

## Testing the API with Swagger UI

### 1. Access Swagger UI

Navigate to `http://localhost:3000/swagger/index.html`

### 2. Test Login

1. Click on **POST /api/login**
2. Click **Try it out**
3. Enter credentials:
   ```json
   {
     "email": "admin@rdbsls.rw",
     "password": "admin123"
   }
   ```
4. Click **Execute**
5. Copy the `token` from the response

### 3. Authorize for Protected Endpoints

1. Click the **Authorize** button at the top
2. Enter the token you copied
3. Click **Authorize**
4. Now you can test protected endpoints

## File Structure

```
tctssf/
├── docs/
│   ├── docs.go          # Generated Go code
│   ├── swagger.json     # Generated JSON spec
│   └── swagger.yaml     # Generated YAML spec
├── main.go              # Swagger annotations here
└── controllers/
    └── auth_controller.go  # Endpoint annotations here
```

## Troubleshooting

### Swagger UI Shows "Failed to load API definition"

**Solution**: Make sure you've imported the docs package in `main.go`:

```go
import (
    _ "tctssf/docs"  // Required!
    // ... other imports
)
```

### Changes Not Reflected

**Solution**: Regenerate docs after making changes:

```bash
swag init
go build -o tctssf main.go
./tctssf
```

### swag Command Not Found

**Solution**: Install swag CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Endpoint Not Showing

**Checklist**:
1. ✅ Added `@Router` annotation with correct path
2. ✅ Ran `swag init`
3. ✅ Rebuilt the application
4. ✅ Refreshed the Swagger UI page

## Advanced Configuration

### Custom Swagger Host

Update in `main.go`:

```go
// @host api.yourdomain.com
// @BasePath /api
```

Then regenerate:

```bash
swag init
```

### HTTPS in Swagger

```go
// @schemes https http
```

### Additional Response Types

```go
// @Success 200 {array} models.User "List of users"
// @Success 200 {string} string "Success message"
// @Success 200 {object} map[string]interface{} "Custom response"
```

## Production Considerations

### Security

- Disable Swagger in production (or protect with authentication)
- Use environment variable to control:

```go
if os.Getenv("ENABLE_SWAGGER") == "true" {
    app.Get("/swagger/*", fiberSwagger.WrapHandler)
}
```

### Performance

- Swagger UI is served from the binary (embedded)
- No external dependencies needed
- Minimal performance impact

## Next Steps

### Recommended Documentation

1. **User Endpoints**
   - GET `/api/profile`
   - GET `/api/dashboard`
   - POST `/api/savings/update-commitment`

2. **Loan Endpoints**
   - POST `/api/loans/apply`
   - GET `/api/loans`
   - POST `/api/loans/:id/repay`

3. **Admin Endpoints**
   - POST `/api/admin/members`
   - GET `/api/admin/members`
   - GET `/api/admin/loans/pending`

4. **Treasurer Endpoints**
   - GET `/api/treasurer/dashboard`
   - POST `/api/treasurer/salary-deductions/generate`

### Template for New Endpoints

```go
// MethodName does something
// @Summary Short description
// @Description Longer description
// @Tags CategoryName
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paramName paramLocation paramType required "description"
// @Success 200 {object} ResponseType "Success message"
// @Failure 400 {object} map[string]string "Error message"
// @Router /path [method]
func (ctrl *Controller) MethodName(c *fiber.Ctx) error {
    // Implementation
}
```

## Resources

- [Swag Documentation](https://github.com/swaggo/swag)
- [Swagger Specification](https://swagger.io/specification/)
- [Fiber Swagger](https://github.com/gofiber/swagger)

---

**Status**: ✅ Ready to use
**Last Updated**: 2025-10-31
