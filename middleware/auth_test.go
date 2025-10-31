package middleware

import (
	"net/http/httptest"
	"testing"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NoToken(t *testing.T) {
	app := fiber.New()
	app.Get("/test", AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode, "Should return 401 when no token provided")
}

func TestStoreSession(t *testing.T) {
	token := "test-token-123"
	userID := 1
	role := "member"

	StoreSession(token, userID, role)

	sessionData, exists := Sessions[token]
	assert.True(t, exists, "Session should be stored in memory")
	assert.Equal(t, userID, sessionData.UserID, "UserID should match")
	assert.Equal(t, role, sessionData.Role, "Role should match")

	// Cleanup
	DeleteSession(token)
}

func TestDeleteSession(t *testing.T) {
	token := "test-token-delete"
	Sessions[token] = models.SessionData{
		UserID: 1,
		Role:   "admin",
	}

	DeleteSession(token)

	_, exists := Sessions[token]
	assert.False(t, exists, "Session should be deleted")
}

func TestRoleMiddleware(t *testing.T) {
	app := fiber.New()

	// Mock auth middleware that sets role
	mockAuth := func(c *fiber.Ctx) error {
		c.Locals("userRole", "admin")
		return c.Next()
	}

	app.Get("/admin", mockAuth, RoleMiddleware("admin"), func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/member", mockAuth, RoleMiddleware("member"), func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test allowed role
	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Admin should access admin route")

	// Test forbidden role
	req2 := httptest.NewRequest("GET", "/member", nil)
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, 403, resp2.StatusCode, "Admin should not access member-only route")
}
