package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"tctssf/config"
	"tctssf/models"

	"github.com/gofiber/fiber/v2"
)

type LoanController struct {
	db *sql.DB
}

// NewLoanController creates a new loan controller
func NewLoanController() *LoanController {
	return &LoanController{
		db: config.GetDB(),
	}
}

// Helper function to determine specific role based on email/name
func (lc *LoanController) getSpecificRole(userID int, userRole string) string {
	if userRole == "treasurer" {
		return "treasurer"
	}
	if userRole == "superadmin" {
		return "president" // SuperAdmin can act as president
	}
	
	// For admin role, check email or specific_role to distinguish
	var email, specificRole string
	var specificRoleNull sql.NullString
	
	err := lc.db.QueryRow(`
		SELECT email, specific_role 
		FROM users 
		WHERE id = ?`, userID).Scan(&email, &specificRoleNull)
	
	if err != nil {
		log.Printf("Error getting user details: %v", err)
		return "admin" // Default fallback
	}
	
	// If specific_role column exists and has value, use it
	if specificRoleNull.Valid {
		specificRole = specificRoleNull.String
		switch specificRole {
		case "vice_president":
			return "vice_president"
		case "president":
			return "president"
		case "treasurer":
			return "treasurer"
		default:
			return "admin"
		}
	}
	
	// Fallback to email-based detection
	emailLower := strings.ToLower(email)
	if strings.Contains(emailLower, "vicepresident") || strings.Contains(emailLower, "vice") {
		return "vice_president"
	}
	if strings.Contains(emailLower, "president") {
		return "president"
	}
	
	return "admin" // Default for other admins
}

// ApplyForLoan handles loan application with updated eligibility rules
func (lc *LoanController) ApplyForLoan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	var req models.LoanApplicationRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Validate input
	if req.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Loan amount must be greater than 0"})
	}
	
	if req.RepaymentPeriod <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Repayment period must be greater than 0"})
	}
	
	// Minimum loan amount validation
	if req.Amount < 10000 {
		return c.Status(400).JSON(fiber.Map{"error": "Minimum loan amount is RWF 10,000"})
	}
	
	// Check for any existing active loans
	var activeLoans int
	lc.db.QueryRow(`
		SELECT COUNT(*) FROM loans 
		WHERE user_id = ? AND status IN ('pending', 'treasurer_approved', 'vice_president_approved', 'president_approved', 'disbursed')`, userID).Scan(&activeLoans)
	
	if activeLoans > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "You already have an active loan application or disbursed loan. Please complete your current loan before applying for a new one.",
		})
	}
	
	// Get user's current savings balance
	var currentBalance float64
	var monthlyCommitment float64
	
	err := lc.db.QueryRow(`
		SELECT COALESCE(current_balance, 0), COALESCE(monthly_commitment, 0) 
		FROM savings_accounts WHERE user_id = ?`, userID).Scan(&currentBalance, &monthlyCommitment)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Create savings account if it doesn't exist
			_, createErr := lc.db.Exec(`
				INSERT INTO savings_accounts (user_id, monthly_commitment, current_balance) 
				VALUES (?, 0, 0)`, userID)
			if createErr != nil {
				log.Printf("Failed to create savings account: %v", createErr)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to process loan application"})
			}
			currentBalance = 0
			monthlyCommitment = 0
		} else {
			log.Printf("Database error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
	}
	
	// Check if user has a monthly commitment set
	if monthlyCommitment <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "You must set a monthly savings commitment before applying for a loan. Please update your commitment first.",
		})
	}
	
	// Check if user has sufficient savings balance
	if currentBalance <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "You must have savings in your account before applying for a loan. Your current balance is RWF 0.",
		})
	}
	
	// New eligibility rule - Max loan = 6x current savings balance
	maxEligibleLoan := currentBalance * 6
	
	if req.Amount > maxEligibleLoan {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Loan amount exceeds eligibility. Maximum eligible: %s (6x your current savings balance of %s)", 
				lc.formatCurrency(maxEligibleLoan), 
				lc.formatCurrency(currentBalance)),
		})
	}
	
	// Calculate monthly payment
	interestRate := 5.0 // 5% annual interest rate
	monthlyPayment := lc.calculateMonthlyPayment(req.Amount, interestRate, req.RepaymentPeriod)
	
	// Create loan application with initial pending status
	_, err = lc.db.Exec(`
		INSERT INTO loans (user_id, amount, repayment_period, monthly_payment, remaining_balance, status) 
		VALUES (?, ?, ?, ?, ?, 'pending')`,
		userID, req.Amount, req.RepaymentPeriod, monthlyPayment, req.Amount)
	
	if err != nil {
		log.Printf("Failed to create loan application: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit application"})
	}
	
	return c.JSON(models.LoanApplicationResponse{
		Message: "Loan application submitted successfully. It will be reviewed by the treasurer first, then vice-president, and finally the president.",
		Amount: req.Amount,
		RepaymentPeriod: req.RepaymentPeriod,
		MonthlyPayment: monthlyPayment,
		MaxEligible: maxEligibleLoan,
	})
}

// GetLoans returns user's loan history
func (lc *LoanController) GetLoans(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	
	rows, err := lc.db.Query(`
		SELECT id, amount, interest_rate, repayment_period, status, created_at, monthly_payment 
		FROM loans WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var loans []models.Loan
	for rows.Next() {
		var l models.Loan
		rows.Scan(&l.ID, &l.Amount, &l.InterestRate, &l.RepaymentPeriod, &l.Status, &l.CreatedAt, &l.MonthlyPayment)
		loans = append(loans, l)
	}
	
	return c.JSON(loans)
}

// GetPendingLoans returns loans pending for current user's approval level
func (lc *LoanController) GetPendingLoans(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	userRole := c.Locals("userRole").(string)
	
	// Get specific role for better filtering
	specificRole := lc.getSpecificRole(userID, userRole)
	
	log.Printf("User ID: %d, Role: %s, Specific Role: %s", userID, userRole, specificRole)
	
	// Determine what loans this user can see based on their specific role
	var statusFilter string
	switch specificRole {
	case "treasurer":
		statusFilter = "status = 'pending'"
	case "vice_president":
		statusFilter = "status = 'treasurer_approved'"
	case "president":
		statusFilter = "status = 'vice_president_approved'"
	case "admin":
		// General admin can see treasurer_approved and vice_president_approved
		statusFilter = "status IN ('treasurer_approved', 'vice_president_approved')"
	default:
		if userRole == "superadmin" {
			// SuperAdmin can see all pending loans
			statusFilter = "status IN ('pending', 'treasurer_approved', 'vice_president_approved')"
		} else {
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}
	}
	
	log.Printf("Status filter: %s", statusFilter)
	
	query := fmt.Sprintf(`
		SELECT 
			l.id, l.user_id, l.amount, l.repayment_period, l.monthly_payment, l.status, l.created_at,
			u.first_name, u.last_name, u.account_number,
			t1.first_name as treasurer_name, l.treasurer_approved_at,
			t2.first_name as vp_name, l.vice_president_approved_at,
			t3.first_name as president_name, l.president_approved_at,
			r.first_name as rejected_by_name, l.rejected_at, l.rejection_reason
		FROM loans l 
		JOIN users u ON l.user_id = u.id 
		LEFT JOIN users t1 ON l.treasurer_approved_by = t1.id
		LEFT JOIN users t2 ON l.vice_president_approved_by = t2.id
		LEFT JOIN users t3 ON l.president_approved_by = t3.id
		LEFT JOIN users r ON l.rejected_by = r.id
		WHERE %s
		ORDER BY l.created_at ASC`, statusFilter)
	
	log.Printf("Executing query: %s", query)
	
	rows, err := lc.db.Query(query)
	if err != nil {
		log.Printf("Error querying loans: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var loans []models.LoanWithDetails
	for rows.Next() {
		var l models.LoanWithDetails
		var firstName, lastName string
		var treasurerName, vpName, presidentialName, rejectedByName sql.NullString
		
		err := rows.Scan(
			&l.ID, &l.UserID, &l.Amount, &l.RepaymentPeriod, &l.MonthlyPayment, &l.Status, &l.CreatedAt,
			&firstName, &lastName, &l.AccountNumber,
			&treasurerName, &l.TreasurerApprovedAt,
			&vpName, &l.VicePresidentApprovedAt,
			&presidentialName, &l.PresidentApprovedAt,
			&rejectedByName, &l.RejectedAt, &l.RejectionReason)
		
		if err != nil {
			log.Printf("Error scanning loan row: %v", err)
			continue
		}
		
		l.MemberName = firstName + " " + lastName
		
		// Set approver names
		if treasurerName.Valid {
			name := treasurerName.String
			l.TreasurerApprovedBy = &name
		}
		if vpName.Valid {
			name := vpName.String
			l.VicePresidentApprovedBy = &name
		}
		if presidentialName.Valid {
			name := presidentialName.String
			l.PresidentApprovedBy = &name
		}
		if rejectedByName.Valid {
			name := rejectedByName.String
			l.RejectedBy = &name
		}
		
		// Determine next approver and if current user can approve
		l.NextApprover, l.CanCurrentUserApprove = lc.getApprovalStatus(l.Status, specificRole)
		
		loans = append(loans, l)
	}
	
	log.Printf("Found %d loans for user", len(loans))
	
	// Ensure we always return an array, even if empty
	if loans == nil {
		loans = []models.LoanWithDetails{}
	}
	
	return c.JSON(loans)
}

// ProcessLoanApproval handles approval or rejection by appropriate role
func (lc *LoanController) ProcessLoanApproval(c *fiber.Ctx) error {
	loanID := c.Params("id")
	userID := c.Locals("userID").(int)
	userRole := c.Locals("userRole").(string)
	
	// Get specific role
	specificRole := lc.getSpecificRole(userID, userRole)
	
	var req models.LoanApprovalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Validate action
	if req.Action != "approve" && req.Action != "reject" {
		return c.Status(400).JSON(fiber.Map{"error": "Action must be 'approve' or 'reject'"})
	}
	
	// If rejecting, reason is required
	if req.Action == "reject" && req.Reason == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Rejection reason is required"})
	}
	
	// Get current loan status
	var currentStatus string
	err := lc.db.QueryRow("SELECT status FROM loans WHERE id = ?", loanID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Loan not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	
	// Check if user can approve at current stage
	canApprove := lc.canUserApproveAtStage(currentStatus, specificRole)
	if !canApprove {
		return c.Status(403).JSON(fiber.Map{
			"error": fmt.Sprintf("You cannot approve loans at stage: %s with role: %s", currentStatus, specificRole),
		})
	}
	
	if req.Action == "reject" {
		return lc.rejectLoan(loanID, userID, req.Reason)
	} else {
		return lc.approveLoan(loanID, userID, specificRole, currentStatus)
	}
}

// Helper function to approve loan at appropriate stage
func (lc *LoanController) approveLoan(loanID string, userID int, specificRole, currentStatus string) error {
	var updateQuery string
	var message string
	
	switch currentStatus {
	case "pending":
		if specificRole != "treasurer" {
			return fiber.NewError(403, "Only treasurer can approve at this stage")
		}
		updateQuery = `
			UPDATE loans 
			SET status = 'treasurer_approved', treasurer_approved_by = ?, treasurer_approved_at = NOW() 
			WHERE id = ?`
		message = "Loan approved by treasurer. Waiting for vice-president approval."
		
	case "treasurer_approved":
		if specificRole != "vice_president" && specificRole != "admin" {
			return fiber.NewError(403, "Only vice-president can approve at this stage")
		}
		updateQuery = `
			UPDATE loans 
			SET status = 'vice_president_approved', vice_president_approved_by = ?, vice_president_approved_at = NOW() 
			WHERE id = ?`
		message = "Loan approved by vice-president. Waiting for president approval."
		
	case "vice_president_approved":
		if specificRole != "president" && specificRole != "admin" {
			return fiber.NewError(403, "Only president can approve at this stage")
		}
		updateQuery = `
			UPDATE loans 
			SET status = 'president_approved', president_approved_by = ?, president_approved_at = NOW() 
			WHERE id = ?`
		message = "Loan fully approved by president. Ready for disbursement."
		
	default:
		return fiber.NewError(400, "Loan cannot be approved at current stage")
	}
	
	_, err := lc.db.Exec(updateQuery, userID, loanID)
	if err != nil {
		log.Printf("Failed to approve loan: %v", err)
		return fiber.NewError(500, "Failed to approve loan")
	}
	
	return fiber.NewError(200, message)
}

// Helper function to reject loan
func (lc *LoanController) rejectLoan(loanID string, userID int, reason string) error {
	_, err := lc.db.Exec(`
		UPDATE loans 
		SET status = 'rejected', rejected_by = ?, rejected_at = NOW(), rejection_reason = ? 
		WHERE id = ?`, userID, reason, loanID)
	
	if err != nil {
		log.Printf("Failed to reject loan: %v", err)
		return fiber.NewError(500, "Failed to reject loan")
	}
	
	return fiber.NewError(200, "Loan rejected successfully")
}

// Helper function to check if user can approve at current stage
func (lc *LoanController) canUserApproveAtStage(status, specificRole string) bool {
	switch status {
	case "pending":
		return specificRole == "treasurer"
	case "treasurer_approved":
		return specificRole == "vice_president" || specificRole == "admin"
	case "vice_president_approved":
		return specificRole == "president" || specificRole == "admin"
	default:
		return false
	}
}

// Helper function to get approval status and next approver
func (lc *LoanController) getApprovalStatus(status, specificRole string) (string, bool) {
	switch status {
	case "pending":
		return "treasurer", specificRole == "treasurer"
	case "treasurer_approved":
		canApprove := specificRole == "vice_president" || specificRole == "admin"
		return "vice_president", canApprove
	case "vice_president_approved":
		canApprove := specificRole == "president" || specificRole == "admin"
		return "president", canApprove
	case "president_approved":
		return "completed", false
	case "rejected":
		return "rejected", false
	default:
		return "unknown", false
	}
}

// RepayLoan handles loan repayment
func (lc *LoanController) RepayLoan(c *fiber.Ctx) error {
	// Implementation for loan repayment
	return c.JSON(fiber.Map{"message": "Repayment processed"})
}

// calculateMonthlyPayment calculates the monthly payment for a loan
func (lc *LoanController) calculateMonthlyPayment(principal, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 100 / 12
	if monthlyRate == 0 {
		return principal / float64(months)
	}
	return principal * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) / (math.Pow(1+monthlyRate, float64(months)) - 1)
}

// formatCurrency formats amount as currency
func (lc *LoanController) formatCurrency(amount float64) string {
	return fmt.Sprintf("RWF %.0f", amount)
}