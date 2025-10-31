package config

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Load configuration
	if AppConfig == nil {
		LoadConfig()
	}

	log.Println("Attempting to connect to database...")
	log.Printf("Database: %s@%s:%s/%s", AppConfig.DBUser, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName)

	// Connect to the existing database using configuration
	DB, err = sql.Open("mysql", AppConfig.GetDBConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Testing database connection...")
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	log.Println("Creating tables...")
	createTables()

	log.Println("Setting up default users...")
	createDefaultUsers()

	log.Println("Database initialization complete")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			account_number VARCHAR(20) UNIQUE,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20),
			password_hash VARCHAR(255) NOT NULL,
			role ENUM('member', 'admin', 'superadmin', 'treasurer') DEFAULT 'member',
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS savings_accounts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			monthly_commitment DECIMAL(10,2) DEFAULT 5000,
			current_balance DECIMAL(10,2) DEFAULT 0,
			social_contributions DECIMAL(10,2) DEFAULT 0,
			last_contribution TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			type ENUM('savings', 'social_contribution', 'loan_disbursement', 'loan_repayment', 'salary_deduction', 'commitment_change') NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			description TEXT,
			reference_id INT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS loans (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			interest_rate DECIMAL(5,2) DEFAULT 5.0,
			repayment_period INT NOT NULL,
			monthly_payment DECIMAL(10,2) DEFAULT 0,
			remaining_balance DECIMAL(10,2) DEFAULT 0,
			status ENUM('pending', 'approved', 'rejected', 'disbursed', 'completed') DEFAULT 'pending',
			approved_by INT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			disbursed_at TIMESTAMP NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS loan_repayments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			loan_id INT NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			due_date DATE NOT NULL,
			paid_date TIMESTAMP NULL,
			status ENUM('pending', 'paid', 'overdue') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS salary_deduction_lists (
			id INT AUTO_INCREMENT PRIMARY KEY,
			month_year VARCHAR(7) NOT NULL,
			generated_by INT NOT NULL,
			status ENUM('generated', 'sent_to_hr', 'processed') DEFAULT 'generated',
			total_members INT DEFAULT 0,
			total_amount DECIMAL(12,2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			processed_at TIMESTAMP NULL,
			FOREIGN KEY (generated_by) REFERENCES users(id),
			UNIQUE KEY unique_month_year (month_year)
		)`,
		`CREATE TABLE IF NOT EXISTS salary_deduction_items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			list_id INT NOT NULL,
			user_id INT NOT NULL,
			account_number VARCHAR(20) NOT NULL,
			member_name VARCHAR(200) NOT NULL,
			monthly_commitment DECIMAL(10,2) DEFAULT 0,
			social_contribution DECIMAL(10,2) DEFAULT 0,
			loan_repayment DECIMAL(10,2) DEFAULT 0,
			total_deduction DECIMAL(10,2) DEFAULT 0,
			status ENUM('pending', 'processed') DEFAULT 'pending',
			FOREIGN KEY (list_id) REFERENCES salary_deduction_lists(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

func createDefaultUsers() {
	createDefaultSuperAdmin()
	createDefaultAdmin()
	createDefaultTreasurer()
}

func createDefaultSuperAdmin() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'superadmin'").Scan(&count)
	
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := DB.Exec(`
			INSERT INTO users (account_number, first_name, last_name, email, password_hash, role) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			"1000000", "Super", "Admin", "superadmin@rdbsls.rw", string(hashedPassword), "superadmin")
		if err != nil {
			log.Printf("Failed to create default superadmin: %v", err)
		} else {
			log.Println("Default superadmin created: superadmin@rdbsls.rw / admin123")
		}
	}
}

func createDefaultAdmin() {
	var adminCount int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount)
	
	if adminCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := DB.Exec(`
			INSERT INTO users (account_number, first_name, last_name, email, password_hash, role) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			"1000002", "Default", "Admin", "admin@rdbsls.rw", string(hashedPassword), "admin")
		if err != nil {
			log.Printf("Failed to create default admin: %v", err)
		} else {
			log.Println("Default admin created: admin@rdbsls.rw / admin123")
		}
	}
}

func createDefaultTreasurer() {
	var treasurerCount int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'treasurer'").Scan(&treasurerCount)
	
	if treasurerCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("treasurer123"), bcrypt.DefaultCost)
		_, err := DB.Exec(`
			INSERT INTO users (account_number, first_name, last_name, email, password_hash, role) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			"1000001", "Default", "Treasurer", "treasurer@rdbsls.rw", string(hashedPassword), "treasurer")
		if err != nil {
			log.Printf("Failed to create default treasurer: %v", err)
		} else {
			log.Println("Default treasurer created: treasurer@rdbsls.rw / treasurer123")
		}
	}
}

// GenerateAccountNumber generates a unique account number
func GenerateAccountNumber() string {
	var accountNumber string
	for {
		// Generate random 4 digits and prepend with 100
		bytes := make([]byte, 2)
		rand.Read(bytes)
		suffix := fmt.Sprintf("%04d", int(bytes[0])<<8+int(bytes[1])%10000)
		accountNumber = "100" + suffix
		
		var exists bool
		DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE account_number = ?)", accountNumber).Scan(&exists)
		if !exists {
			break
		}
	}
	return accountNumber
}