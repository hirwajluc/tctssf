package models

import "time"

// Loan represents loan applications and status with 3-stage approval
type Loan struct {
	ID                    int        `json:"id"`
	UserID                int        `json:"user_id"`
	Amount                float64    `json:"amount"`
	InterestRate          float64    `json:"interest_rate"`
	RepaymentPeriod       int        `json:"repayment_period"`
	Status                string     `json:"status"` // pending, treasurer_approved, vice_president_approved, president_approved, rejected, disbursed, completed
	TreasurerApprovedBy   *int       `json:"treasurer_approved_by"`
	TreasurerApprovedAt   *time.Time `json:"treasurer_approved_at"`
	VicePresidentApprovedBy *int     `json:"vice_president_approved_by"`
	VicePresidentApprovedAt *time.Time `json:"vice_president_approved_at"`
	PresidentApprovedBy   *int       `json:"president_approved_by"`
	PresidentApprovedAt   *time.Time `json:"president_approved_at"`
	RejectedBy            *int       `json:"rejected_by"`
	RejectedAt            *time.Time `json:"rejected_at"`
	RejectionReason       string     `json:"rejection_reason"`
	CreatedAt             time.Time  `json:"created_at"`
	DisbursedAt           *time.Time `json:"disbursed_at"`
	MonthlyPayment        float64    `json:"monthly_payment"`
	RemainingBalance      float64    `json:"remaining_balance"`
}

// LoanRepayment tracks loan repayments
type LoanRepayment struct {
	ID        int        `json:"id"`
	LoanID    int        `json:"loan_id"`
	Amount    float64    `json:"amount"`
	DueDate   time.Time  `json:"due_date"`
	PaidDate  *time.Time `json:"paid_date"`
	Status    string     `json:"status"` // pending, paid, overdue
}

// LoanApplicationRequest for applying for loans
type LoanApplicationRequest struct {
	Amount          float64 `json:"amount"`
	RepaymentPeriod int     `json:"repayment_period"`
}

// LoanApplicationResponse for loan application result
type LoanApplicationResponse struct {
	Message         string  `json:"message"`
	Amount          float64 `json:"amount"`
	RepaymentPeriod int     `json:"repayment_period"`
	MonthlyPayment  float64 `json:"monthly_payment"`
	MaxEligible     float64 `json:"max_eligible"`
}

// LoanApprovalRequest for approval actions
type LoanApprovalRequest struct {
	Action string `json:"action"` // approve, reject
	Reason string `json:"reason,omitempty"` // required for rejection
}

// LoanWithDetails represents loan with member and approval details
type LoanWithDetails struct {
	ID                      int        `json:"id"`
	UserID                  int        `json:"user_id"`
	Amount                  float64    `json:"amount"`
	RepaymentPeriod         int        `json:"repayment_period"`
	MonthlyPayment          float64    `json:"monthly_payment"`
	Status                  string     `json:"status"`
	CreatedAt               time.Time  `json:"created_at"`
	MemberName              string     `json:"member_name"`
	AccountNumber           string     `json:"account_number"`
	TreasurerApprovedBy     *string    `json:"treasurer_approved_by"`
	TreasurerApprovedAt     *time.Time `json:"treasurer_approved_at"`
	VicePresidentApprovedBy *string    `json:"vice_president_approved_by"`
	VicePresidentApprovedAt *time.Time `json:"vice_president_approved_at"`
	PresidentApprovedBy     *string    `json:"president_approved_by"`
	PresidentApprovedAt     *time.Time `json:"president_approved_at"`
	RejectedBy              *string    `json:"rejected_by"`
	RejectedAt              *time.Time `json:"rejected_at"`
	RejectionReason         string     `json:"rejection_reason"`
	NextApprover            string     `json:"next_approver"` // treasurer, vice_president, president, completed
	CanCurrentUserApprove   bool       `json:"can_current_user_approve"`
}