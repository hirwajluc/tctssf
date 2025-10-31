package routes

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"tctssf/controllers"
	"tctssf/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	adminController *controllers.AdminController,
	loanController *controllers.LoanController,
	treasurerController *controllers.TreasurerController,
) {
	// API routes
	api := app.Group("/api")
	
	// Public routes
	api.Post("/login", authController.Login)
	
	// Protected API routes
	protected := api.Group("/", middleware.AuthMiddleware)
	
	// User routes
	protected.Get("/profile", userController.GetProfile)
	protected.Get("/dashboard", userController.GetDashboard)
	protected.Post("/savings/update-commitment", userController.UpdateCommitment)
	protected.Get("/transactions", userController.GetTransactions)
	
	// Loan routes (member level)
	protected.Post("/loans/apply", loanController.ApplyForLoan)
	protected.Get("/loans", loanController.GetLoans)
	protected.Post("/loans/:id/repay", loanController.RepayLoan)
	
	// Admin routes (excluding loan management - moved to separate section)
	admin := protected.Group("/admin", middleware.RoleMiddleware("admin", "superadmin", "treasurer"))
	admin.Post("/members", adminController.CreateMember)
	admin.Get("/members", adminController.GetMembers)
	admin.Post("/users", adminController.CreateUser)
	admin.Get("/users", adminController.GetAllUsers)
	admin.Get("/users/:id", adminController.GetUserByID)
	admin.Put("/users/:id", adminController.UpdateUser)
	admin.Delete("/users/:id", adminController.DeleteUser)
	admin.Get("/reports", adminController.GetReports)
	
	// Bulk import routes
	admin.Post("/members/bulk-import", adminController.BulkImportMembers)
	admin.Get("/members/bulk-template", adminController.DownloadBulkTemplate)
	
	// FIXED: Loan Management Routes - Now supports treasurer + admin roles
	// All roles that can participate in loan approval process
	loanManagement := protected.Group("/admin", middleware.RoleMiddleware("admin", "superadmin", "treasurer"))
	
	// Get pending loans (filtered by user role and approval stage)
	loanManagement.Get("/loans/pending", loanController.GetPendingLoans)
	
	// Unified loan processing endpoint (approve/reject)
	loanManagement.Post("/loans/:id/process", loanController.ProcessLoanApproval)
	
	// DEPRECATED: Keep old endpoints for backward compatibility
	loanManagement.Post("/loans/:id/approve", func(c *fiber.Ctx) error {
		// Set action to approve and call ProcessLoanApproval
		c.Request().SetBody([]byte(`{"action":"approve"}`))
		return loanController.ProcessLoanApproval(c)
	})
	
	loanManagement.Post("/loans/:id/reject", func(c *fiber.Ctx) error {
		// For backward compatibility, set default rejection reason
		c.Request().SetBody([]byte(`{"action":"reject","reason":"Rejected via legacy endpoint"}`))
		return loanController.ProcessLoanApproval(c)
	})
	
	// SuperAdmin routes
	superAdmin := protected.Group("/superadmin", middleware.RoleMiddleware("superadmin"))
	superAdmin.Post("/admins", adminController.CreateAdmin)
	superAdmin.Get("/system-stats", adminController.GetSystemStats)
	
	// Treasurer routes
	treasurer := protected.Group("/treasurer", middleware.RoleMiddleware("treasurer", "superadmin"))
	treasurer.Get("/dashboard", treasurerController.GetDashboard)
	treasurer.Post("/salary-deductions/generate", treasurerController.GenerateSalaryDeductionList)
	treasurer.Get("/salary-deductions/lists", treasurerController.GetSalaryDeductionLists)
	treasurer.Get("/salary-deductions/:id", treasurerController.GetSalaryDeductionDetails)
	treasurer.Get("/salary-deductions/:id/download", treasurerController.DownloadSalaryDeductionCSV)
	treasurer.Post("/salary-deductions/:id/upload-processed", treasurerController.UploadProcessedSalaries)
	treasurer.Delete("/salary-deductions/:id", treasurerController.DeleteSalaryDeductionList)
	treasurer.Get("/members-summary", treasurerController.GetMembersSummary)

	// ADD THESE NEW TREASURER LOAN ROUTES:
	treasurer.Get("/loans/pending", treasurerController.GetPendingLoansForTreasurer)
	treasurer.Post("/loans/:id/approve", loanController.ProcessLoanApproval) // Reuse existing controller
	treasurer.Post("/loans/:id/reject", loanController.ProcessLoanApproval)  // Reuse existing controller
	
	// Add explicit route for testing static files
	app.Get("/test-static", func(c *fiber.Ctx) error {
		return c.SendString("Static file handler is working!")
	})
	
	// Static files handler (must be last) - NO AUTH REQUIRED FOR HTML FILES
	app.Use(serveStaticFiles)
	
	log.Println("Routes configured successfully")
	log.Println("Static file directory: ./frontend")
	log.Println("Loan approval system: 3-stage approval enabled")
	log.Println("  - Stage 1: Treasurer approval")
	log.Println("  - Stage 2: Vice-President (admin) approval") 
	log.Println("  - Stage 3: President (admin/superadmin) approval")
	log.Println("FIXED: Treasurer role now has access to loan management endpoints")
}

// serveStaticFiles handles static file serving with SPA support
func serveStaticFiles(c *fiber.Ctx) error {
	path := c.Path()

	log.Printf("Static file request for path: %s", path)

	// Skip Swagger routes - they should be handled by Swagger middleware
	if strings.HasPrefix(path, "/swagger") {
		log.Printf("Swagger route detected, passing to Swagger handler: %s", path)
		return c.Next()
	}

	// Handle root path
	if path == "/" {
		log.Printf("Serving index.html for root path")
		return c.SendFile("./frontend/index.html")
	}

	// Check if it's an API call
	if strings.HasPrefix(path, "/api/") {
		log.Printf("API call detected, passing to next handler: %s", path)
		return c.Next()
	}
	
	// Try to serve the requested file
	filePath := filepath.Join("./frontend", path)
	log.Printf("Attempting to serve file: %s", filePath)
	
	// Check if file exists
	if _, err := http.Dir("./frontend").Open(path); err == nil {
		log.Printf("File found, serving: %s", filePath)
		return c.SendFile(filePath)
	} else {
		log.Printf("File not found: %s, error: %v", filePath, err)
	}
	
	// For HTML routes that don't exist as files, serve index.html (SPA behavior)
	if strings.HasSuffix(path, ".html") || !strings.Contains(path, ".") {
		log.Printf("HTML file not found, serving index.html for SPA: %s", path)
		return c.SendFile("./frontend/index.html")
	}
	
	// For other missing files, return 404
	log.Printf("File not found, returning 404: %s", path)
	return c.Status(404).SendString("File not found")
}