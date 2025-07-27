package middleware

import (
	"log"
	"tctssf/config"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
)

// Simple session store (in production, use Redis or proper session management)
var Sessions = make(map[string]models.SessionData)

// AuthMiddleware validates user authentication
func AuthMiddleware(c *fiber.Ctx) error {
	sessionToken := c.Get("Authorization")
	if sessionToken == "" {
		// Try to get token from query parameter (fallback)
		sessionToken = c.Query("token")
	}
	
	if sessionToken == "" {
		log.Printf("No authorization token provided for path: %s", c.Path())
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - No token provided"})
	}
	
	// Check if session exists
	sessionData, exists := Sessions[sessionToken]
	if !exists {
		log.Printf("Invalid session token: %s for path: %s", sessionToken, c.Path())
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}
	
	// Verify user still exists and is active
	var isActive bool
	db := config.GetDB()
	err := db.QueryRow("SELECT is_active FROM users WHERE id = ?", sessionData.UserID).Scan(&isActive)
	if err != nil || !isActive {
		// Clean up invalid session
		delete(Sessions, sessionToken)
		log.Printf("User not found or inactive for session: %s", sessionToken)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}
	
	// Store user data in context
	c.Locals("userID", sessionData.UserID)
	c.Locals("userRole", sessionData.Role)
	return c.Next()
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole").(string)
		
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}
		
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}
}