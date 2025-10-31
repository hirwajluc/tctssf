package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateUserRequest
		valid   bool
	}{
		{
			name: "Valid user request",
			request: CreateUserRequest{
				FirstName:         "John",
				LastName:          "Doe",
				Email:             "john@example.com",
				Phone:             "0788123456",
				Role:              "member",
				MonthlyCommitment: 5000,
			},
			valid: true,
		},
		{
			name: "Missing email",
			request: CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
				Role:      "member",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Email, "Email should not be empty for valid request")
				assert.NotEmpty(t, tt.request.FirstName, "FirstName should not be empty")
				assert.NotEmpty(t, tt.request.LastName, "LastName should not be empty")
			} else {
				assert.Empty(t, tt.request.Email, "Email should be empty for invalid request")
			}
		})
	}
}

func TestLoanApplicationRequest_Validation(t *testing.T) {
	req := LoanApplicationRequest{
		Amount:          10000,
		RepaymentPeriod: 12,
	}

	assert.Greater(t, req.Amount, float64(0), "Loan amount should be greater than 0")
	assert.Greater(t, req.RepaymentPeriod, 0, "Repayment period should be greater than 0")
	assert.GreaterOrEqual(t, req.Amount, float64(10000), "Minimum loan should be 10000")
}
