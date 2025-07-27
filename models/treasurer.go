package models

import "time"

// SalaryDeductionList represents the monthly list for HR
type SalaryDeductionList struct {
	ID          int       `json:"id"`
	MonthYear   string    `json:"month_year"`
	GeneratedBy int       `json:"generated_by"`
	Status      string    `json:"status"` // generated, sent_to_hr, processed
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	TotalMembers int      `json:"total_members"`
	TotalAmount float64   `json:"total_amount"`
}

// SalaryDeductionItem represents individual member deductions
type SalaryDeductionItem struct {
	ID                int     `json:"id"`
	ListID           int     `json:"list_id"`
	UserID           int     `json:"user_id"`
	AccountNumber    string  `json:"account_number"`
	MemberName       string  `json:"member_name"`
	MonthlyCommitment float64 `json:"monthly_commitment"`
	SocialContribution float64 `json:"social_contribution"`
	LoanRepayment    float64 `json:"loan_repayment"`
	TotalDeduction   float64 `json:"total_deduction"`
	Status           string  `json:"status"` // pending, processed
}

// ProcessedSalaryData for CSV upload processing
type ProcessedSalaryData struct {
	AccountNumber      string  `json:"account_number"`
	DeductedAmount     float64 `json:"deducted_amount"`
	SavingsAmount      float64 `json:"savings_amount"`
	SocialAmount       float64 `json:"social_amount"`
	LoanRepaymentAmount float64 `json:"loan_repayment_amount"`
	Row                int     `json:"row"`
}

// GenerateDeductionListRequest for generating salary deduction lists
type GenerateDeductionListRequest struct {
	MonthYear string `json:"month_year"` // Format: YYYY-MM
}

// TreasurerDashboardStats for treasurer dashboard
type TreasurerDashboardStats struct {
	TotalActiveMembers    int     `json:"total_active_members"`
	TotalMonthlySavings   float64 `json:"total_monthly_savings"`
	TotalSocialFunds      float64 `json:"total_social_funds"`
	PendingLoanRepayments float64 `json:"pending_loan_repayments"`
	CurrentMonthProcessed bool    `json:"current_month_processed"`
	LastProcessedMonth    string  `json:"last_processed_month"`
}

// MemberSummary for member summary reports
type MemberSummary struct {
	ID                     int     `json:"id"`
	AccountNumber          string  `json:"account_number"`
	Name                   string  `json:"name"`
	MonthlyCommitment      float64 `json:"monthly_commitment"`
	CurrentBalance         float64 `json:"current_balance"`
	SocialContributions    float64 `json:"social_contributions"`
	TotalLoanBalance       float64 `json:"total_loan_balance"`
	MonthlyLoanPayment     float64 `json:"monthly_loan_payment"`
	TotalMonthlyDeduction  float64 `json:"total_monthly_deduction"`
}

// DeleteSalaryDeductionRequest for deleting salary deduction lists with password confirmation
type DeleteSalaryDeductionRequest struct {
	Password string `json:"password" validate:"required"`
}