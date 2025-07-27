package controllers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"tctssf/config"
	"tctssf/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
)

type TreasurerController struct {
	db *sql.DB
}

// NewTreasurerController creates a new treasurer controller
func NewTreasurerController() *TreasurerController {
	return &TreasurerController{
		db: config.GetDB(),
	}
}

// Replaced the GetDashboard function in controllers/treasurer_controller.go
// The section that calculates TotalSocialFunds and replaced it with this corrected version:

// GetDashboard returns treasurer dashboard statistics
func (tc *TreasurerController) GetDashboard(c *fiber.Ctx) error {
	var stats models.TreasurerDashboardStats
	
	// Get total active members
	tc.db.QueryRow(`
		SELECT COUNT(*) FROM users 
		WHERE role = 'member' AND is_active = 1`).Scan(&stats.TotalActiveMembers)
	
	// Get total monthly commitments (personal savings commitments)
	tc.db.QueryRow(`
		SELECT COALESCE(SUM(s.monthly_commitment), 0)
		FROM savings_accounts s
		JOIN users u ON s.user_id = u.id
		WHERE u.role = 'member' AND u.is_active = 1`).Scan(&stats.TotalMonthlySavings)
	
	// CORRECTED: Get total social contributions (social_contributions only, not current_balance)
	tc.db.QueryRow(`
		SELECT COALESCE(SUM(s.social_contributions), 0)
		FROM savings_accounts s
		JOIN users u ON s.user_id = u.id
		WHERE u.role = 'member' AND u.is_active = 1`).Scan(&stats.TotalSocialFunds)
	
	// Get pending loan repayments for this month
	currentMonth := time.Now().Format("2006-01")
	tc.db.QueryRow(`
		SELECT COALESCE(SUM(l.monthly_payment), 0)
		FROM loans l
		JOIN users u ON l.user_id = u.id
		WHERE l.status = 'disbursed' AND u.is_active = 1`).Scan(&stats.PendingLoanRepayments)
	
	// Check if current month is processed
	tc.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM salary_deduction_lists 
			WHERE month_year = ? AND status = 'processed'
		)`, currentMonth).Scan(&stats.CurrentMonthProcessed)
	
	// Get last processed month
	tc.db.QueryRow(`
		SELECT COALESCE(MAX(month_year), 'None') 
		FROM salary_deduction_lists 
		WHERE status = 'processed'`).Scan(&stats.LastProcessedMonth)
	
	return c.JSON(stats)
}

// GenerateSalaryDeductionList generates a salary deduction list for HR
func (tc *TreasurerController) GenerateSalaryDeductionList(c *fiber.Ctx) error {
	treasurerID := c.Locals("userID").(int)
	
	var req models.GenerateDeductionListRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Validate month format
	if _, err := time.Parse("2006-01", req.MonthYear); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid month format. Use YYYY-MM"})
	}
	
	// Check if list already exists for this month
	var existingID int
	err := tc.db.QueryRow(`
		SELECT id FROM salary_deduction_lists 
		WHERE month_year = ?`, req.MonthYear).Scan(&existingID)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Salary deduction list already exists for this month"})
	}
	
	// Debug: Check member counts BEFORE starting transaction
	var debugUserCount int
	tc.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&debugUserCount)
	log.Printf("Total users in database: %d", debugUserCount)
	
	var debugMemberCount int
	tc.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'member'").Scan(&debugMemberCount)
	log.Printf("Total members in database: %d", debugMemberCount)
	
	var debugActiveMemberCount int
	tc.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'member' AND is_active = 1").Scan(&debugActiveMemberCount)
	log.Printf("Total active members in database: %d", debugActiveMemberCount)
	
	// If no active members, return error immediately
	if debugActiveMemberCount == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("No active members found to process. Debug info - Total users: %d, Total members: %d, Active members: %d", 
				debugUserCount, debugMemberCount, debugActiveMemberCount),
		})
	}
	
	// Start transaction
	tx, err := tc.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database transaction failed"})
	}
	defer tx.Rollback()
	
	// Create the list
	result, err := tx.Exec(`
		INSERT INTO salary_deduction_lists (month_year, generated_by, status) 
		VALUES (?, ?, 'generated')`, req.MonthYear, treasurerID)
	if err != nil {
		log.Printf("Failed to create deduction list: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create deduction list"})
	}
	
	listID, _ := result.LastInsertId()
	
	// Get all active members with their commitments and loans (using main db connection, not transaction)
	rows, err := tc.db.Query(`
		SELECT 
			u.id, u.account_number, u.first_name, u.last_name,
			COALESCE(s.monthly_commitment, 0) as monthly_commitment
		FROM users u
		LEFT JOIN savings_accounts s ON u.id = s.user_id
		WHERE u.role = 'member' AND u.is_active = 1
		ORDER BY u.first_name, u.last_name`)
	
	if err != nil {
		log.Printf("Failed to fetch member data: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch member data"})
	}
	defer rows.Close()
	
	var totalMembers int
	var totalAmount float64
	socialContribution := 5000.0 // Fixed social contribution
	
	for rows.Next() {
		var userID int
		var accountNumber, firstName, lastName string
		var monthlyCommitment float64
		
		err := rows.Scan(&userID, &accountNumber, &firstName, &lastName, &monthlyCommitment)
		if err != nil {
			log.Printf("Failed to scan member row: %v", err)
			continue
		}
		
		// Get loan repayment for this member (using main db connection)
		var loanRepayment float64
		err = tc.db.QueryRow(`
			SELECT COALESCE(SUM(monthly_payment), 0) 
			FROM loans 
			WHERE user_id = ? AND status = 'disbursed'`, userID).Scan(&loanRepayment)
		if err != nil {
			log.Printf("Failed to get loan repayment for user %d: %v", userID, err)
			loanRepayment = 0 // Continue with 0 if there's an error
		}
		
		memberName := firstName + " " + lastName
		totalDeduction := monthlyCommitment + socialContribution + loanRepayment
		
		// Insert deduction item (using transaction)
		_, err = tx.Exec(`
			INSERT INTO salary_deduction_items (
				list_id, user_id, account_number, member_name, 
				monthly_commitment, social_contribution, loan_repayment, total_deduction
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			listID, userID, accountNumber, memberName, 
			monthlyCommitment, socialContribution, loanRepayment, totalDeduction)
		
		if err != nil {
			log.Printf("Failed to insert deduction item for user %d: %v", userID, err)
			continue
		}
		
		totalMembers++
		totalAmount += totalDeduction
		
		log.Printf("Added member: %s (%s) - Commitment: %.2f, Social: %.2f, Loan: %.2f, Total: %.2f", 
			memberName, accountNumber, monthlyCommitment, socialContribution, loanRepayment, totalDeduction)
	}
	
	log.Printf("Processed %d members with total amount: %.2f", totalMembers, totalAmount)
	
	// Check if we processed any members
	if totalMembers == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to process any members despite finding active members in database",
		})
	}
	
	// Update the list with totals
	_, err = tx.Exec(`
		UPDATE salary_deduction_lists 
		SET total_members = ?, total_amount = ? 
		WHERE id = ?`, totalMembers, totalAmount, listID)
	
	if err != nil {
		log.Printf("Failed to update list totals. Members: %d, Amount: %f, ListID: %d, Error: %v", 
			totalMembers, totalAmount, listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update list totals"})
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	
	return c.JSON(fiber.Map{
		"message":       "Salary deduction list generated successfully",
		"list_id":       listID,
		"month_year":    req.MonthYear,
		"total_members": totalMembers,
		"total_amount":  totalAmount,
	})
}

// GetSalaryDeductionLists returns all salary deduction lists
func (tc *TreasurerController) GetSalaryDeductionLists(c *fiber.Ctx) error {
	rows, err := tc.db.Query(`
		SELECT 
			sdl.id, sdl.month_year, sdl.status, sdl.total_members, 
			sdl.total_amount, sdl.created_at, sdl.processed_at,
			u.first_name, u.last_name
		FROM salary_deduction_lists sdl
		JOIN users u ON sdl.generated_by = u.id
		ORDER BY sdl.created_at DESC`)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var lists []map[string]interface{}
	for rows.Next() {
		var id, totalMembers int
		var monthYear, status, firstName, lastName string
		var totalAmount float64
		var createdAt time.Time
		var processedAt sql.NullTime
		
		rows.Scan(&id, &monthYear, &status, &totalMembers, &totalAmount, 
			&createdAt, &processedAt, &firstName, &lastName)
		
		list := map[string]interface{}{
			"id":            id,
			"month_year":    monthYear,
			"status":        status,
			"total_members": totalMembers,
			"total_amount":  totalAmount,
			"created_at":    createdAt,
			"generated_by":  firstName + " " + lastName,
		}
		
		if processedAt.Valid {
			list["processed_at"] = processedAt.Time
		}
		
		lists = append(lists, list)
	}
	
	return c.JSON(lists)
}

// GetSalaryDeductionDetails returns details of a specific salary deduction list
func (tc *TreasurerController) GetSalaryDeductionDetails(c *fiber.Ctx) error {
	listID := c.Params("id")
	
	// Get list details
	var list map[string]interface{} = make(map[string]interface{})
	var id, totalMembers int
	var monthYear, status string
	var totalAmount float64
	var createdAt time.Time
	var processedAt sql.NullTime
	
	err := tc.db.QueryRow(`
		SELECT id, month_year, status, total_members, total_amount, created_at, processed_at
		FROM salary_deduction_lists WHERE id = ?`, listID).Scan(
		&id, &monthYear, &status, &totalMembers, &totalAmount, &createdAt, &processedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Salary deduction list not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	
	list["id"] = id
	list["month_year"] = monthYear
	list["status"] = status
	list["total_members"] = totalMembers
	list["total_amount"] = totalAmount
	list["created_at"] = createdAt
	
	if processedAt.Valid {
		list["processed_at"] = processedAt.Time
	}
	
	// Get deduction items
	rows, err := tc.db.Query(`
		SELECT 
			user_id, account_number, member_name, monthly_commitment,
			social_contribution, loan_repayment, total_deduction, status
		FROM salary_deduction_items
		WHERE list_id = ?
		ORDER BY member_name`, listID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch deduction items"})
	}
	defer rows.Close()
	
	var items []map[string]interface{}
	for rows.Next() {
		var userID int
		var accountNumber, memberName, itemStatus string
		var monthlyCommitment, socialContribution, loanRepayment, totalDeduction float64
		
		rows.Scan(&userID, &accountNumber, &memberName, &monthlyCommitment,
			&socialContribution, &loanRepayment, &totalDeduction, &itemStatus)
		
		items = append(items, map[string]interface{}{
			"user_id":            userID,
			"account_number":     accountNumber,
			"member_name":        memberName,
			"monthly_commitment": monthlyCommitment,
			"social_contribution": socialContribution,
			"loan_repayment":     loanRepayment,
			"total_deduction":    totalDeduction,
			"status":             itemStatus,
		})
	}
	
	list["items"] = items
	return c.JSON(list)
}

// DownloadSalaryDeductionCSV generates and downloads a CSV file for HR
func (tc *TreasurerController) DownloadSalaryDeductionCSV(c *fiber.Ctx) error {
	listID := c.Params("id")
	
	// Verify list exists
	var monthYear string
	err := tc.db.QueryRow("SELECT month_year FROM salary_deduction_lists WHERE id = ?", listID).Scan(&monthYear)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Salary deduction list not found"})
	}
	
	// Get deduction items
	rows, err := tc.db.Query(`
		SELECT 
			account_number, member_name, monthly_commitment,
			social_contribution, loan_repayment, total_deduction
		FROM salary_deduction_items
		WHERE list_id = ?
		ORDER BY member_name`, listID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch deduction items"})
	}
	defer rows.Close()
	
	// Create CSV data
	var csvData [][]string
	csvData = append(csvData, []string{
		"Account Number", "Member Name", "Monthly Savings", 
		"Social Contribution", "Loan Repayment", "Total Deduction", 
		"Deducted Amount", "Status"})
	
	for rows.Next() {
		var accountNumber, memberName string
		var monthlyCommitment, socialContribution, loanRepayment, totalDeduction float64
		
		rows.Scan(&accountNumber, &memberName, &monthlyCommitment,
			&socialContribution, &loanRepayment, &totalDeduction)
		
		csvData = append(csvData, []string{
			accountNumber,
			memberName,
			fmt.Sprintf("%.2f", monthlyCommitment),
			fmt.Sprintf("%.2f", socialContribution),
			fmt.Sprintf("%.2f", loanRepayment),
			fmt.Sprintf("%.2f", totalDeduction),
			"", // Empty column for HR to fill deducted amount
			"", // Empty column for HR to fill status
		})
	}
	
	// Set headers for file download
	filename := fmt.Sprintf("salary_deductions_%s.csv", monthYear)
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	
	// Create CSV content
	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)
	
	for _, row := range csvData {
		if err := writer.Write(row); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate CSV"})
		}
	}
	
	writer.Flush()
	if err := writer.Error(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate CSV"})
	}
	
	// Update status to sent_to_hr
	tc.db.Exec("UPDATE salary_deduction_lists SET status = 'sent_to_hr' WHERE id = ?", listID)
	
	return c.SendString(csvContent.String())
}

// UploadProcessedSalaries processes the uploaded CSV from HR
func (tc *TreasurerController) UploadProcessedSalaries(c *fiber.Ctx) error {
	listID := c.Params("id")
	
	// Verify list exists and is in correct status
	var status string
	err := tc.db.QueryRow("SELECT status FROM salary_deduction_lists WHERE id = ?", listID).Scan(&status)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Salary deduction list not found"})
	}
	
	if status == "processed" {
		return c.Status(400).JSON(fiber.Map{"error": "This list has already been processed"})
	}
	
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}
	
	// Validate file type
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		return c.Status(400).JSON(fiber.Map{"error": "Only CSV files are supported"})
	}
	
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open uploaded file"})
	}
	defer src.Close()
	
	// Parse CSV
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse CSV file"})
	}
	
	if len(records) < 2 {
		return c.Status(400).JSON(fiber.Map{"error": "CSV file must contain header and data rows"})
	}
	
	// Validate headers
	header := records[0]
	expectedHeaders := []string{"account number", "member name", "monthly savings", 
		"social contribution", "loan repayment", "total deduction", "deducted amount", "status"}
	
	// Normalize headers
	for i := range header {
		header[i] = strings.ToLower(strings.TrimSpace(header[i]))
	}
	
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[h] = i
	}
	
	// Check required headers
	requiredIndices := make(map[string]int)
	for _, required := range []string{"account number", "deducted amount", "status"} {
		if idx, exists := headerMap[required]; exists {
			requiredIndices[required] = idx
		} else {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("Missing required column: %s", required),
			})
		}
	}
	
	// Start transaction
	tx, err := tc.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database transaction failed"})
	}
	defer tx.Rollback()
	
	var processedCount int
	var errors []map[string]interface{}
	
	for i, record := range records[1:] { // Skip header
		rowNum := i + 2
		
		if len(record) < len(expectedHeaders) {
			errors = append(errors, map[string]interface{}{
				"row":     rowNum,
				"error":   "Insufficient columns",
				"account": "",
			})
			continue
		}
		
		accountNumber := strings.TrimSpace(record[requiredIndices["account number"]])
		deductedAmountStr := strings.TrimSpace(record[requiredIndices["deducted amount"]])
		status := strings.ToLower(strings.TrimSpace(record[requiredIndices["status"]]))
		
		// Skip rows where status is not "processed" or "paid"
		if status != "processed" && status != "paid" && status != "deducted" {
			continue
		}
		
		// Parse deducted amount
		deductedAmount, err := strconv.ParseFloat(deductedAmountStr, 64)
		if err != nil || deductedAmount <= 0 {
			errors = append(errors, map[string]interface{}{
				"row":     rowNum,
				"error":   "Invalid deducted amount",
				"account": accountNumber,
			})
			continue
		}
		
		// Get deduction item details
		var userID int
		var monthlyCommitment, socialContribution, loanRepayment float64
		err = tx.QueryRow(`
			SELECT user_id, monthly_commitment, social_contribution, loan_repayment
			FROM salary_deduction_items
			WHERE list_id = ? AND account_number = ?`, listID, accountNumber).Scan(
			&userID, &monthlyCommitment, &socialContribution, &loanRepayment)
		
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"row":     rowNum,
				"error":   "Account not found in deduction list",
				"account": accountNumber,
			})
			continue
		}
		
		// Update savings account
		_, err = tx.Exec(`
			UPDATE savings_accounts 
			SET current_balance = current_balance + ?,
				social_contributions = social_contributions + ?,
				last_contribution = NOW()
			WHERE user_id = ?`,
			monthlyCommitment, socialContribution, userID)
		
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"row":     rowNum,
				"error":   "Failed to update savings",
				"account": accountNumber,
			})
			continue
		}
		
		// Process loan repayment if any
		if loanRepayment > 0 {
			// Update loan balance
			_, err = tx.Exec(`
				UPDATE loans 
				SET remaining_balance = remaining_balance - ?
				WHERE user_id = ? AND status = 'disbursed'`,
				loanRepayment, userID)
			
			if err != nil {
				log.Printf("Failed to update loan balance for user %d: %v", userID, err)
			}
			
			// Record loan repayment transaction
			_, err = tx.Exec(`
				INSERT INTO transactions (user_id, type, amount, description, reference_id)
				VALUES (?, 'loan_repayment', ?, ?, ?)`,
				userID, loanRepayment, fmt.Sprintf("Salary deduction loan repayment"), userID)
		}
		
		// Record savings transaction
		_, err = tx.Exec(`
			INSERT INTO transactions (user_id, type, amount, description, reference_id)
			VALUES (?, 'savings', ?, ?, ?)`,
			userID, monthlyCommitment, fmt.Sprintf("Monthly savings via salary deduction"), userID)
		
		// Record social contribution transaction
		_, err = tx.Exec(`
			INSERT INTO transactions (user_id, type, amount, description, reference_id)
			VALUES (?, 'social_contribution', ?, ?, ?)`,
			userID, socialContribution, fmt.Sprintf("Social contribution via salary deduction"), userID)
		
		// Update deduction item status
		_, err = tx.Exec(`
			UPDATE salary_deduction_items 
			SET status = 'processed' 
			WHERE list_id = ? AND user_id = ?`, listID, userID)
		
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"row":     rowNum,
				"error":   "Failed to update item status",
				"account": accountNumber,
			})
			continue
		}
		
		processedCount++
	}
	
	// Update list status
	_, err = tx.Exec(`
		UPDATE salary_deduction_lists 
		SET status = 'processed', processed_at = NOW() 
		WHERE id = ?`, listID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update list status"})
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	
	return c.JSON(fiber.Map{
		"message":         "Salary deductions processed successfully",
		"processed_count": processedCount,
		"error_count":     len(errors),
		"errors":          errors,
	})
}

// GetMembersSummary returns member summary for treasurer
func (tc *TreasurerController) GetMembersSummary(c *fiber.Ctx) error {
	rows, err := tc.db.Query(`
		SELECT 
			u.id, u.account_number, u.first_name, u.last_name,
			COALESCE(s.monthly_commitment, 0) as monthly_commitment,
			COALESCE(s.current_balance, 0) as current_balance,
			COALESCE(s.social_contributions, 0) as social_contributions,
			COALESCE(SUM(l.remaining_balance), 0) as total_loan_balance,
			COALESCE(SUM(l.monthly_payment), 0) as monthly_loan_payment
		FROM users u
		LEFT JOIN savings_accounts s ON u.id = s.user_id
		LEFT JOIN loans l ON u.id = l.user_id AND l.status = 'disbursed'
		WHERE u.role = 'member' AND u.is_active = true
		GROUP BY u.id, u.account_number, u.first_name, u.last_name, 
				 s.monthly_commitment, s.current_balance, s.social_contributions
		ORDER BY u.first_name, u.last_name`)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var members []models.MemberSummary
	for rows.Next() {
		var member models.MemberSummary
		var firstName, lastName string
		
		rows.Scan(&member.ID, &member.AccountNumber, &firstName, &lastName,
			&member.MonthlyCommitment, &member.CurrentBalance, &member.SocialContributions,
			&member.TotalLoanBalance, &member.MonthlyLoanPayment)
		
		member.Name = firstName + " " + lastName
		member.TotalMonthlyDeduction = member.MonthlyCommitment + 5000 + member.MonthlyLoanPayment // 5000 is fixed social contribution
		
		members = append(members, member)
	}
	
	return c.JSON(members)
}

// DeleteSalaryDeductionList handles the deletion of salary deduction lists with password confirmation
func (tc *TreasurerController) DeleteSalaryDeductionList(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid list ID"})
	}
	
	treasurerID := c.Locals("userID").(int)
	
	var req models.DeleteSalaryDeductionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}
	
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}
	
	// Verify the current user's password
	var storedPasswordHash string
	err = tc.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", treasurerID).Scan(&storedPasswordHash)
	if err != nil {
		log.Printf("Failed to get user password hash for user %d: %v", treasurerID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Authentication failed"})
	}
	
	// Compare the provided password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(req.Password)); err != nil {
		log.Printf("Invalid password provided for salary deduction deletion by user %d", treasurerID)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password"})
	}
	
	// Check if list exists and get its status and month_year for response
	var status string
	var monthYear string
	err = tc.db.QueryRow("SELECT status, month_year FROM salary_deduction_lists WHERE id = ?", listID).Scan(&status, &monthYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Salary deduction list not found"})
		}
		log.Printf("Database error while checking list %d: %v", listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	
	// Prevent deletion of processed lists
	if status == "processed" {
		log.Printf("Attempt to delete processed salary deduction list %d by user %d", listID, treasurerID)
		return c.Status(400).JSON(fiber.Map{"error": "Cannot delete processed salary deduction lists"})
	}
	
	// Start transaction to delete the list and its items
	tx, err := tc.db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction for deleting list %d: %v", listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Database transaction failed"})
	}
	defer tx.Rollback()
	
	// Delete salary deduction items first (due to foreign key constraint)
	result, err := tx.Exec("DELETE FROM salary_deduction_items WHERE list_id = ?", listID)
	if err != nil {
		log.Printf("Failed to delete salary deduction items for list %d: %v", listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete salary deduction items"})
	}
	
	itemsDeleted, _ := result.RowsAffected()
	log.Printf("Deleted %d salary deduction items for list %d", itemsDeleted, listID)
	
	// Delete the salary deduction list
	result, err = tx.Exec("DELETE FROM salary_deduction_lists WHERE id = ?", listID)
	if err != nil {
		log.Printf("Failed to delete salary deduction list %d: %v", listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete salary deduction list"})
	}
	
	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("No salary deduction list found with ID %d for deletion", listID)
		return c.Status(404).JSON(fiber.Map{"error": "Salary deduction list not found"})
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit deletion transaction for list %d: %v", listID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit deletion"})
	}
	
	log.Printf("Salary deduction list %d (%s) deleted successfully by treasurer %d", listID, monthYear, treasurerID)
	
	return c.JSON(fiber.Map{
		"message":         fmt.Sprintf("Salary deduction list for %s deleted successfully", monthYear),
		"deleted_list_id": listID,
		"month_year":      monthYear,
	})
}

// Add this method to your TreasurerController

// GetPendingLoansForTreasurer returns loans pending treasurer approval
func (tc *TreasurerController) GetPendingLoansForTreasurer(c *fiber.Ctx) error {
    //userID := c.Locals("userID").(int)
    userRole := c.Locals("userRole").(string)
    
    // Verify treasurer role
    if userRole != "treasurer" && userRole != "superadmin" {
        return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
    }
    
    // Get loans pending treasurer approval
    rows, err := tc.db.Query(`
        SELECT 
            l.id, l.user_id, l.amount, l.repayment_period, l.monthly_payment, l.status, l.created_at,
            u.first_name, u.last_name, u.account_number
        FROM loans l 
        JOIN users u ON l.user_id = u.id 
        WHERE l.status = 'pending'
        ORDER BY l.created_at ASC`)
    
    if err != nil {
        log.Printf("Error querying loans for treasurer: %v", err)
        return c.Status(500).JSON(fiber.Map{"error": "Database error"})
    }
    defer rows.Close()
    
    var loans []map[string]interface{}
    for rows.Next() {
        var id, userID, repaymentPeriod int
        var amount, monthlyPayment float64
        var status, firstName, lastName, accountNumber string
        var createdAt time.Time
        
        err := rows.Scan(&id, &userID, &amount, &repaymentPeriod, &monthlyPayment, &status, 
                        &createdAt, &firstName, &lastName, &accountNumber)
        if err != nil {
            log.Printf("Error scanning loan row: %v", err)
            continue
        }
        
        loans = append(loans, map[string]interface{}{
            "id":               id,
            "user_id":          userID,
            "amount":           amount,
            "repayment_period": repaymentPeriod,
            "monthly_payment":  monthlyPayment,
            "status":           status,
            "created_at":       createdAt,
            "member_name":      firstName + " " + lastName,
            "account_number":   accountNumber,
            "can_approve":      true, // Treasurer can approve pending loans
            "next_approver":    "treasurer",
        })
    }
    
    // Ensure we always return an array
    if loans == nil {
        loans = []map[string]interface{}{}
    }
    
    return c.JSON(loans)
}