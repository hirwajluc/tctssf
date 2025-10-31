package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"tctssf/config"
	"tctssf/middleware"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *sql.DB
}

// NewAuthController creates a new auth controller
func NewAuthController() *AuthController {
	return &AuthController{
		db: config.GetDB(),
	}
}

// Login handles user authentication
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	log.Printf("Login attempt for email: %s", req.Email)
	
	var user models.User
	var passwordHash string
	var phone sql.NullString // Handle NULL phone values
	
	err := ac.db.QueryRow(`
		SELECT id, account_number, first_name, last_name, email, phone, role, password_hash, is_active 
		FROM users WHERE email = ?`, req.Email).Scan(
		&user.ID, &user.AccountNumber, &user.FirstName, &user.LastName, 
		&user.Email, &phone, &user.Role, &passwordHash, &user.IsActive)
	
	// Convert NullString to regular string
	if phone.Valid {
		user.Phone = phone.String
	} else {
		user.Phone = ""
	}
	
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for email: %s", req.Email)
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		log.Printf("Database error during login: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	
	if !user.IsActive {
		log.Printf("Inactive user attempted login: %s", req.Email)
		return c.Status(401).JSON(fiber.Map{"error": "Account is inactive"})
	}
	
	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		log.Printf("Invalid password for user: %s", req.Email)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	
	// Generate session token
	token := ac.generateSessionToken()

	// Store session (Redis with in-memory fallback)
	middleware.StoreSession(token, user.ID, user.Role)

	return c.JSON(models.LoginResponse{
		User:  user,
		Token: token,
	})
}

// generateSessionToken generates a unique session token
func (ac *AuthController) generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}