package middleware

import (
	"log"
	"tctssf/config"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
)

// In-memory session store (fallback when Redis is unavailable)
var Sessions = make(map[string]models.SessionData)

// AuthMiddleware validates user authentication with Redis (fallback to in-memory)
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

	var userID int
	var role string
	var exists bool

	// Try Redis first
	if config.RedisClient != nil {
		var err error
		userID, role, exists, err = config.GetSession(sessionToken)
		if err != nil {
			log.Printf("Redis error, falling back to in-memory: %v", err)
			// Fall back to in-memory on error
			sessionData, memExists := Sessions[sessionToken]
			if !memExists {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
			}
			userID = sessionData.UserID
			role = sessionData.Role
			exists = true
		}
	} else {
		// Use in-memory if Redis not available
		sessionData, memExists := Sessions[sessionToken]
		if !memExists {
			log.Printf("Invalid session token: %s for path: %s", sessionToken, c.Path())
			return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
		}
		userID = sessionData.UserID
		role = sessionData.Role
		exists = true
	}

	if !exists {
		log.Printf("Invalid session token: %s for path: %s", sessionToken, c.Path())
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}

	// Verify user still exists and is active
	var isActive bool
	db := config.GetDB()
	err := db.QueryRow("SELECT is_active FROM users WHERE id = ?", userID).Scan(&isActive)
	if err != nil || !isActive {
		// Clean up invalid session
		if config.RedisClient != nil {
			config.DeleteSession(sessionToken)
		}
		delete(Sessions, sessionToken)
		log.Printf("User not found or inactive for session: %s", sessionToken)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}

	// Store user data in context
	c.Locals("userID", userID)
	c.Locals("userRole", role)
	return c.Next()
}

// StoreSession stores a session in Redis with in-memory fallback
func StoreSession(token string, userID int, role string) {
	// Store in Redis if available
	if config.RedisClient != nil {
		err := config.SetSession(token, userID, role, 0) // 0 for no expiration
		if err != nil {
			log.Printf("Failed to store session in Redis: %v", err)
		}
	}

	// Always store in-memory as fallback
	Sessions[token] = models.SessionData{
		UserID: userID,
		Role:   role,
	}
}

// DeleteSession removes a session from both Redis and in-memory
func DeleteSession(token string) {
	if config.RedisClient != nil {
		config.DeleteSession(token)
	}
	delete(Sessions, token)
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