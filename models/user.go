package models

import "time"

// User represents system users
type User struct {
	ID            int       `json:"id"`
	AccountNumber string    `json:"account_number"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	IsActive      bool      `json:"is_active"`
}

// SavingsAccount represents member savings
type SavingsAccount struct {
	ID                int     `json:"id"`
	UserID           int     `json:"user_id"`
	MonthlyCommitment float64 `json:"monthly_commitment"`
	CurrentBalance   float64 `json:"current_balance"`
	SocialContributions float64 `json:"social_contributions"`
	LastContribution time.Time `json:"last_contribution"`
}

// Transaction represents all financial transactions
type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Type        string    `json:"type"` // savings, social_contribution, loan_disbursement, loan_repayment
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateUserRequest for creating users
type CreateUserRequest struct {
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	Role              string  `json:"role"`
	MonthlyCommitment float64 `json:"monthly_commitment,omitempty"`
}

// UpdateUserRequest for updating users
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}

// UpdateCommitmentRequest for updating monthly commitment
type UpdateCommitmentRequest struct {
	NewCommitment float64 `json:"new_commitment"`
	EffectiveDate string  `json:"effective_date"`
}