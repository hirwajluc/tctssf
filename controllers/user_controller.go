package controllers

import (
	"database/sql"
	"fmt"
	"time"
	"tctssf/config"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	db *sql.DB
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{
		db: config.GetDB(),
	}
}

// GetProfile returns user profile information
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	var user models.User
	var phone sql.NullString
	
	err := uc.db.QueryRow(`
		SELECT id, account_number, first_name, last_name, email, phone, role, created_at, is_active 
		FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.AccountNumber, &user.FirstName, &user.LastName,
		&user.Email, &phone, &user.Role, &user.CreatedAt, &user.IsActive)
	
	// Convert NullString to regular string
	if phone.Valid {
		user.Phone = phone.String
	} else {
		user.Phone = ""
	}
	
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	
	return c.JSON(user)
}

// GetDashboard returns dashboard data for the user
func (uc *UserController) GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	var savings models.SavingsAccount
	uc.db.QueryRow(`
		SELECT COALESCE(monthly_commitment, 0), COALESCE(current_balance, 0), 
		COALESCE(social_contributions, 0) 
		FROM savings_accounts WHERE user_id = ?`, userID).Scan(
		&savings.MonthlyCommitment, &savings.CurrentBalance, &savings.SocialContributions)
	
	// Get recent transactions
	rows, _ := uc.db.Query(`
		SELECT type, amount, description, created_at 
		FROM transactions WHERE user_id = ? 
		ORDER BY created_at DESC LIMIT 5`, userID)
	defer rows.Close()
	
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		rows.Scan(&t.Type, &t.Amount, &t.Description, &t.CreatedAt)
		transactions = append(transactions, t)
	}
	
	// Get active loans
	var activeLoans int
	uc.db.QueryRow(`
		SELECT COUNT(*) FROM loans 
		WHERE user_id = ? AND status IN ('approved', 'disbursed')`, userID).Scan(&activeLoans)
	
	return c.JSON(fiber.Map{
		"savings":      savings,
		"transactions": transactions,
		"activeLoans":  activeLoans,
	})
}

// UpdateCommitment updates user's monthly savings commitment
func (uc *UserController) UpdateCommitment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	var req models.UpdateCommitmentRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Validate minimum commitment
	if req.NewCommitment < 5000 {
		return c.Status(400).JSON(fiber.Map{"error": "Minimum commitment is RWF 5,000"})
	}
	
	// Parse effective date
	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format"})
	}
	
	// Ensure effective date is in the future
	if effectiveDate.Before(time.Now().AddDate(0, 0, 1)) {
		return c.Status(400).JSON(fiber.Map{"error": "Effective date must be at least tomorrow"})
	}
	
	// Start transaction
	tx, err := uc.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()
	
	// Update savings account commitment
	_, err = tx.Exec(`
		UPDATE savings_accounts 
		SET monthly_commitment = ? 
		WHERE user_id = ?`,
		req.NewCommitment, userID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update commitment"})
	}
	
	// Record the commitment change as a transaction for audit trail
	_, err = tx.Exec(`
		INSERT INTO transactions (user_id, type, amount, description) 
		VALUES (?, ?, ?, ?)`,
		userID, "commitment_change", req.NewCommitment, 
		fmt.Sprintf("Monthly commitment updated to %v RWF, effective from %s", req.NewCommitment, req.EffectiveDate))
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to record commitment change"})
	}
	
	tx.Commit()
	
	return c.JSON(fiber.Map{
		"message": "Monthly commitment updated successfully",
		"new_commitment": req.NewCommitment,
		"effective_date": req.EffectiveDate,
	})
}

// GetTransactions returns user's transaction history
func (uc *UserController) GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	rows, err := uc.db.Query(`SELECT type, amount, description, created_at FROM transactions WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		rows.Scan(&t.Type, &t.Amount, &t.Description, &t.CreatedAt)
		transactions = append(transactions, t)
	}
	
	return c.JSON(transactions)
}