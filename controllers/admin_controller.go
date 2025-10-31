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

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct {
	db *sql.DB
}

// NewAdminController creates a new admin controller
func NewAdminController() *AdminController {
	return &AdminController{
		db: config.GetDB(),
	}
}

// CreateMember creates a new member
func (ac *AdminController) CreateMember(c *fiber.Ctx) error {
	var req models.CreateMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	accountNumber := config.GenerateAccountNumber()
	defaultPassword := "member123" // In production, generate random password and send via email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	
	tx, err := ac.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()
	
	result, err := tx.Exec(`
		INSERT INTO users (account_number, first_name, last_name, email, phone, password_hash, role, joined_date)
		VALUES (?, ?, ?, ?, ?, ?, 'member', ?)`,
		accountNumber, req.FirstName, req.LastName, req.Email, req.Phone, string(hashedPassword), req.JoinedDate)
	
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to create member"})
	}
	
	userID, _ := result.LastInsertId()
	
	// Create savings account
	tx.Exec(`
		INSERT INTO savings_accounts (user_id, monthly_commitment) 
		VALUES (?, ?)`, userID, 5000)
	
	tx.Commit()
	
	return c.JSON(fiber.Map{
		"message":        "Member created successfully",
		"account_number": accountNumber,
		"password":       defaultPassword,
	})
}

// GetMembers returns all members
func (ac *AdminController) GetMembers(c *fiber.Ctx) error {
	rows, err := ac.db.Query(`SELECT id, account_number, first_name, last_name, email, phone, created_at, joined_date, is_active FROM users WHERE role = 'member'`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var members []models.User
	for rows.Next() {
		var member models.User
		var phone sql.NullString
		var joinedDate sql.NullString

		err := rows.Scan(&member.ID, &member.AccountNumber, &member.FirstName, &member.LastName, &member.Email, &phone, &member.CreatedAt, &joinedDate, &member.IsActive)
		if err != nil {
			continue // Skip rows with errors
		}
		
		// Convert NullString to regular string
		if phone.Valid {
			member.Phone = phone.String
		} else {
			member.Phone = ""
		}

		// Convert NullString to pointer for joined_date
		if joinedDate.Valid {
			member.JoinedDate = &joinedDate.String
		}

		members = append(members, member)
	}
	
	return c.JSON(members)
}

// GetAllUsers returns all users
func (ac *AdminController) GetAllUsers(c *fiber.Ctx) error {
	rows, err := ac.db.Query(`SELECT id, account_number, first_name, last_name, email, phone, role, created_at, is_active FROM users ORDER BY created_at DESC`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		var phone sql.NullString
		
		err := rows.Scan(&user.ID, &user.AccountNumber, &user.FirstName, &user.LastName, &user.Email, &phone, &user.Role, &user.CreatedAt, &user.IsActive)
		if err != nil {
			continue // Skip rows with errors
		}
		
		// Convert NullString to regular string
		if phone.Valid {
			user.Phone = phone.String
		} else {
			user.Phone = ""
		}
		
		users = append(users, user)
	}
	
	return c.JSON(users)
}

// CreateUser creates any type of user
func (ac *AdminController) CreateUser(c *fiber.Ctx) error {
	currentUserRole := c.Locals("userRole").(string)
	
	var req models.CreateUserRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Validate role
	if req.Role != "member" && req.Role != "admin" && req.Role != "superadmin" && req.Role != "treasurer" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid role"})
	}
	
	// Role-based restrictions
	if currentUserRole == "admin" && req.Role != "member" {
		return c.Status(403).JSON(fiber.Map{"error": "Admins can only create member accounts"})
	}
	
	accountNumber := config.GenerateAccountNumber()
	defaultPassword := "user123" // In production, generate random password and send via email
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	
	tx, err := ac.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()
	
	result, err := tx.Exec(`
		INSERT INTO users (account_number, first_name, last_name, email, phone, password_hash, role, joined_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		accountNumber, req.FirstName, req.LastName, req.Email, req.Phone, string(hashedPassword), req.Role, req.JoinedDate)
	
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to create user"})
	}
	
	// Create savings account only for members
	if req.Role == "member" {
		userID, _ := result.LastInsertId()
		tx.Exec(`
			INSERT INTO savings_accounts (user_id, monthly_commitment) 
			VALUES (?, ?)`, userID, 5000)
	}
	
	tx.Commit()
	
	return c.JSON(fiber.Map{
		"message":        "User created successfully",
		"account_number": accountNumber,
		"password":       defaultPassword,
	})
}

// GetUserByID returns a user by ID
func (ac *AdminController) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")
	
	var user models.User
	var phone sql.NullString
	
	err := ac.db.QueryRow(`
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
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	
	return c.JSON(user)
}

// UpdateUser updates user information
func (ac *AdminController) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	currentUserRole := c.Locals("userRole").(string)
	
	// Get the user being edited to check permissions
	var existingRole string
	err := ac.db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&existingRole)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	
	var req models.UpdateUserRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	
	// Role-based restrictions for editing
	if currentUserRole == "admin" {
		// Admins can only edit members
		if existingRole != "member" {
			return c.Status(403).JSON(fiber.Map{"error": "Admins can only edit member accounts"})
		}
		// Admins cannot change role to admin/superadmin/treasurer
		if req.Role != "member" {
			return c.Status(403).JSON(fiber.Map{"error": "Admins can only set role to member"})
		}
	}
	
	_, err = ac.db.Exec(`
		UPDATE users 
		SET first_name = ?, last_name = ?, email = ?, phone = ?, role = ?, is_active = ?
		WHERE id = ?`,
		req.FirstName, req.LastName, req.Email, req.Phone, req.Role, req.IsActive, userID)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}
	
	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

// DeleteUser deletes a user
func (ac *AdminController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	
	// Check if user is a superadmin (prevent deletion)
	var role string
	err := ac.db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	
	if role == "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Cannot delete superadmin user"})
	}
	
	_, err = ac.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// BulkImportMembers handles bulk member import
func (ac *AdminController) BulkImportMembers(c *fiber.Ctx) error {
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
		return c.Status(400).JSON(fiber.Map{"error": "CSV file must contain at least header row and one data row"})
	}
	
	// Validate headers (case-insensitive)
	expectedHeaders := []string{"first_name", "last_name", "email", "phone", "monthly_commitment"}
	header := records[0]
	
	// Normalize headers to lowercase
	for i := range header {
		header[i] = strings.ToLower(strings.TrimSpace(header[i]))
	}
	
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[h] = i
	}
	
	// Check required headers
	for _, required := range expectedHeaders {
		if _, exists := headerMap[required]; !exists {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("Missing required column: %s", required),
				"expected_headers": expectedHeaders,
				"found_headers": header,
			})
		}
	}
	
	// Process data rows
	var members []models.BulkMemberData
	var parseErrors []models.BulkImportError
	
	for i, record := range records[1:] { // Skip header row
		rowNum := i + 2 // Excel row number (1-based + header)
		
		if len(record) < len(expectedHeaders) {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "row",
				Message: "Insufficient columns in row",
				Data:    fmt.Sprintf("Expected %d columns, got %d", len(expectedHeaders), len(record)),
			})
			continue
		}
		
		member := models.BulkMemberData{Row: rowNum}
		
		// Extract data based on header positions
		member.FirstName = strings.TrimSpace(record[headerMap["first_name"]])
		member.LastName = strings.TrimSpace(record[headerMap["last_name"]])
		member.Email = strings.TrimSpace(strings.ToLower(record[headerMap["email"]]))
		member.Phone = strings.TrimSpace(record[headerMap["phone"]])
		
		// Parse monthly commitment
		commitmentStr := strings.TrimSpace(record[headerMap["monthly_commitment"]])
		if commitmentStr == "" {
			member.MonthlyCommitment = 0
		} else {
			commitment, err := strconv.ParseFloat(commitmentStr, 64)
			if err != nil {
				parseErrors = append(parseErrors, models.BulkImportError{
					Row:     rowNum,
					Field:   "monthly_commitment",
					Message: "Invalid number format",
					Data:    commitmentStr,
				})
				continue
			}
			member.MonthlyCommitment = commitment
		}
		
		// Validate required fields
		if member.FirstName == "" {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "first_name",
				Message: "First name is required",
				Data:    member.FirstName,
			})
		}
		
		if member.LastName == "" {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "last_name",
				Message: "Last name is required",
				Data:    member.LastName,
			})
		}
		
		if member.Email == "" {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "email",
				Message: "Email is required",
				Data:    member.Email,
			})
		} else if !ac.isValidEmail(member.Email) {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "email",
				Message: "Invalid email format",
				Data:    member.Email,
			})
		}
		
		if member.MonthlyCommitment < 0 {
			parseErrors = append(parseErrors, models.BulkImportError{
				Row:     rowNum,
				Field:   "monthly_commitment",
				Message: "Monthly commitment cannot be negative",
				Data:    fmt.Sprintf("%.2f", member.MonthlyCommitment),
			})
		}
		
		// Only add member if no validation errors for this row
		hasRowError := false
		for _, err := range parseErrors {
			if err.Row == rowNum {
				hasRowError = true
				break
			}
		}
		
		if !hasRowError {
			members = append(members, member)
		}
	}
	
	// If there are parsing errors, return them without processing
	if len(parseErrors) > 0 {
		return c.Status(400).JSON(models.BulkImportResult{
			TotalProcessed: len(records) - 1,
			SuccessCount:   0,
			ErrorCount:     len(parseErrors),
			Errors:         parseErrors,
			SuccessMembers: []models.BulkImportSuccess{},
		})
	}
	
	// Process members
	return ac.processBulkMembers(c, members)
}

// processBulkMembers processes the validated bulk member data
func (ac *AdminController) processBulkMembers(c *fiber.Ctx, members []models.BulkMemberData) error {
	var result models.BulkImportResult
	result.TotalProcessed = len(members)
	
	// Start database transaction
	tx, err := ac.db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database transaction failed"})
	}
	defer tx.Rollback()
	
	for _, member := range members {
		// Check if email already exists
		var existingID int
		err := tx.QueryRow("SELECT id FROM users WHERE email = ?", member.Email).Scan(&existingID)
		if err == nil {
			result.Errors = append(result.Errors, models.BulkImportError{
				Row:     member.Row,
				Field:   "email",
				Message: "Email already exists",
				Data:    member.Email,
			})
			result.ErrorCount++
			continue
		}
		
		// Generate account number and password
		accountNumber := config.GenerateAccountNumber()
		defaultPassword := "member123" // In production, generate random password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		
		// Insert user
		userResult, err := tx.Exec(`
			INSERT INTO users (account_number, first_name, last_name, email, phone, password_hash, role) 
			VALUES (?, ?, ?, ?, ?, ?, 'member')`,
			accountNumber, member.FirstName, member.LastName, member.Email, member.Phone, string(hashedPassword))
		
		if err != nil {
			result.Errors = append(result.Errors, models.BulkImportError{
				Row:     member.Row,
				Field:   "database",
				Message: "Failed to create user account",
				Data:    err.Error(),
			})
			result.ErrorCount++
			continue
		}
		
		userID, _ := userResult.LastInsertId()
		
		// Create savings account
		_, err = tx.Exec(`
			INSERT INTO savings_accounts (user_id, monthly_commitment) 
			VALUES (?, ?)`, userID, 5000)
		
		if err != nil {
			result.Errors = append(result.Errors, models.BulkImportError{
				Row:     member.Row,
				Field:   "database",
				Message: "Failed to create savings account",
				Data:    err.Error(),
			})
			result.ErrorCount++
			continue
		}
		
		// Success
		result.SuccessMembers = append(result.SuccessMembers, models.BulkImportSuccess{
			Row:           member.Row,
			AccountNumber: accountNumber,
			Name:          member.FirstName + " " + member.LastName,
			Email:         member.Email,
			Password:      defaultPassword,
		})
		result.SuccessCount++
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit database transaction"})
	}
	
	return c.JSON(result)
}

// DownloadBulkTemplate provides a CSV template for bulk import
func (ac *AdminController) DownloadBulkTemplate(c *fiber.Ctx) error {
	// Create CSV template with proper escaping
	var csvData [][]string
	csvData = append(csvData, []string{"first_name", "last_name", "email", "phone", "monthly_commitment"})
	csvData = append(csvData, []string{"John", "Doe", "john.doe@example.com", "+250788123456", "10000"})
	csvData = append(csvData, []string{"Jane", "Smith", "jane.smith@example.com", "+250788654321", "15000"})
	csvData = append(csvData, []string{"Robert", "Johnson", "robert.johnson@example.com", "+250788111222", "12000"})
	csvData = append(csvData, []string{"Alice", "Brown", "alice.brown@example.com", "+250788333444", "8000"})
	csvData = append(csvData, []string{"Michael", "Wilson", "michael.wilson@example.com", "+250788555666", "20000"})
	
	// Set headers for file download
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=\"bulk_members_template.csv\"")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	
	// Create CSV content using proper CSV formatting
	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)
	
	// Write all rows
	for _, row := range csvData {
		if err := writer.Write(row); err != nil {
			log.Printf("Error writing CSV row: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate CSV template"})
		}
	}
	
	// Flush the writer to ensure all data is written
	writer.Flush()
	
	if err := writer.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate CSV template"})
	}
	
	log.Printf("Template download requested, CSV size: %d bytes", csvContent.Len())
	return c.SendString(csvContent.String())
}

// Replaced the GetReports function in controllers/admin_controller.go
// The section that calculates TotalSocialFunds and replaced it with this corrected version:

// GetReports returns system reports and statistics
func (ac *AdminController) GetReports(c *fiber.Ctx) error {
	// Get system statistics
	var stats struct {
		TotalUsers        int     `json:"total_users"`
		TotalMembers      int     `json:"total_members"`
		TotalAdmins       int     `json:"total_admins"`
		TotalSavings      float64 `json:"total_savings"`
		TotalSocialFunds  float64 `json:"total_social_funds"`
		TotalLoans        int     `json:"total_loans"`
		PendingLoans      int     `json:"pending_loans"`
		ApprovedLoans     int     `json:"approved_loans"`
		DisbursedLoans    int     `json:"disbursed_loans"`
		TotalLoanAmount   float64 `json:"total_loan_amount"`
		ActiveUsers       int     `json:"active_users"`
		InactiveUsers     int     `json:"inactive_users"`
		ThisYearFunds     float64 `json:"this_year_funds"`
		LastYearFunds     float64 `json:"last_year_funds"`
	}
	
	// Get user counts
	ac.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'member'").Scan(&stats.TotalMembers)
	ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE role IN ('admin', 'superadmin', 'treasurer')").Scan(&stats.TotalAdmins)
	ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = true").Scan(&stats.ActiveUsers)
	ac.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = false").Scan(&stats.InactiveUsers)
	
	// Get savings data - CORRECTED: Personal savings only (current_balance)
	ac.db.QueryRow("SELECT COALESCE(SUM(current_balance), 0) FROM savings_accounts").Scan(&stats.TotalSavings)
	
	// Get social funds - CORRECTED: Social contributions only (social_contributions)
	ac.db.QueryRow("SELECT COALESCE(SUM(social_contributions), 0) FROM savings_accounts").Scan(&stats.TotalSocialFunds)
	
	// Get loan data
	ac.db.QueryRow("SELECT COUNT(*) FROM loans").Scan(&stats.TotalLoans)
	ac.db.QueryRow("SELECT COUNT(*) FROM loans WHERE status = 'pending'").Scan(&stats.PendingLoans)
	ac.db.QueryRow("SELECT COUNT(*) FROM loans WHERE status = 'approved'").Scan(&stats.ApprovedLoans)
	ac.db.QueryRow("SELECT COUNT(*) FROM loans WHERE status = 'disbursed'").Scan(&stats.DisbursedLoans)
	ac.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM loans WHERE status IN ('approved', 'disbursed')").Scan(&stats.TotalLoanAmount)
	
	// Get yearly fund comparison - CORRECTED: Calculate based on actual transactions
	currentYear := time.Now().Year()
	lastYear := currentYear - 1
	
	// This year's total funds - CORRECTED: Sum savings and social contribution transactions separately
	var thisYearSavings, thisYearSocial float64
	ac.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE YEAR(created_at) = ? AND type = 'savings'`,
		currentYear).Scan(&thisYearSavings)
	
	ac.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE YEAR(created_at) = ? AND type = 'social_contribution'`,
		currentYear).Scan(&thisYearSocial)
	
	stats.ThisYearFunds = thisYearSavings + thisYearSocial
	
	// Last year's total funds - CORRECTED: Same calculation for last year
	var lastYearSavings, lastYearSocial float64
	ac.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE YEAR(created_at) = ? AND type = 'savings'`,
		lastYear).Scan(&lastYearSavings)
	
	ac.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE YEAR(created_at) = ? AND type = 'social_contribution'`,
		lastYear).Scan(&lastYearSocial)
	
	stats.LastYearFunds = lastYearSavings + lastYearSocial
	
	// If no transaction data, use current balances for this year estimate
	if stats.ThisYearFunds == 0 {
		stats.ThisYearFunds = stats.TotalSavings + stats.TotalSocialFunds
	}
	
	// Get recent activities
	recentActivities := []map[string]interface{}{}
	
	// Recent user registrations
	rows, err := ac.db.Query(`
		SELECT first_name, last_name, role, created_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT 5`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var firstName, lastName, role string
			var createdAt time.Time
			rows.Scan(&firstName, &lastName, &role, &createdAt)
			recentActivities = append(recentActivities, map[string]interface{}{
				"type":        "user_created",
				"description": fmt.Sprintf("%s %s (%s) joined the system", firstName, lastName, role),
				"date":        createdAt,
			})
		}
	}
	
	// Recent loan applications
	rows2, err := ac.db.Query(`
		SELECT l.amount, l.status, l.created_at, u.first_name, u.last_name
		FROM loans l
		JOIN users u ON l.user_id = u.id
		ORDER BY l.created_at DESC
		LIMIT 5`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var amount float64
			var status, firstName, lastName string
			var createdAt time.Time
			rows2.Scan(&amount, &status, &createdAt, &firstName, &lastName)
			recentActivities = append(recentActivities, map[string]interface{}{
				"type":        "loan_activity",
				"description": fmt.Sprintf("%s %s applied for RWF %.0f loan (%s)", firstName, lastName, amount, status),
				"date":        createdAt,
			})
		}
	}
	
	// Monthly commitment summary - CORRECTED: Show proper breakdown
	monthlyData := []map[string]interface{}{}
	rows3, err := ac.db.Query(`
		SELECT 
			u.first_name, 
			u.last_name, 
			COALESCE(s.monthly_commitment, 0) as commitment,
			COALESCE(s.current_balance, 0) as balance,
			COALESCE(s.social_contributions, 0) as social
		FROM users u
		LEFT JOIN savings_accounts s ON u.id = s.user_id
		WHERE u.role = 'member' AND u.is_active = true
		ORDER BY s.monthly_commitment DESC`)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var firstName, lastName string
			var commitment, balance, social float64
			rows3.Scan(&firstName, &lastName, &commitment, &balance, &social)
			monthlyData = append(monthlyData, map[string]interface{}{
				"name":               firstName + " " + lastName,
				"commitment":         commitment,
				"personal_savings":   balance,          // CORRECTED: Label as personal savings
				"social_contributions": social,         // CORRECTED: Separate social contributions
				"total_savings":      balance + social, // CORRECTED: Total is sum of both
				"monthly_deduction":  commitment + 5000, // 5000 is fixed social contribution
			})
		}
	}
	
	return c.JSON(fiber.Map{
		"stats":             stats,
		"recent_activities": recentActivities,
		"monthly_data":      monthlyData,
	})
}

// CreateAdmin creates a new admin user
func (ac *AdminController) CreateAdmin(c *fiber.Ctx) error {
	// Implementation for creating admin users
	return c.JSON(fiber.Map{"message": "Admin created"})
}

// GetSystemStats returns system statistics
func (ac *AdminController) GetSystemStats(c *fiber.Ctx) error {
	// Implementation for system statistics
	return c.JSON(fiber.Map{"message": "System stats"})
}

// isValidEmail performs simple email validation
func (ac *AdminController) isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}