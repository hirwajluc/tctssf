package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	policy := DefaultPasswordPolicy()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid strong password",
			password: "MyP@ssw0rd123",
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Pass1!",
			wantErr:  true,
		},
		{
			name:     "No uppercase",
			password: "mypassword123!",
			wantErr:  true,
		},
		{
			name:     "No lowercase",
			password: "MYPASSWORD123!",
			wantErr:  true,
		},
		{
			name:     "No number",
			password: "MyPassword!",
			wantErr:  true,
		},
		{
			name:     "No special character",
			password: "MyPassword123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password, policy)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for password: %s", tt.password)
			} else {
				assert.NoError(t, err, "Expected no error for password: %s", tt.password)
			}
		})
	}
}

func TestIsCommonPassword(t *testing.T) {
	tests := []struct {
		password string
		isCommon bool
	}{
		{"password", true},
		{"123456", true},
		{"qwerty", true},
		{"MySecureP@ssw0rd", false},
		{"letmein", true},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := IsCommonPassword(tt.password)
			assert.Equal(t, tt.isCommon, result, "Password: %s", tt.password)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid.email", false},
		{"@domain.com", false},
		{"user@", false},
		{"user name@domain.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			assert.Equal(t, tt.valid, result, "Email: %s", tt.email)
		})
	}
}

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		password      string
		minStrength   int
		maxStrength   int
		description   string
	}{
		{"pass", 0, 1, "Very weak - short"},
		{"password", 0, 1, "Weak - common"},
		{"Password1", 2, 3, "Fair"},
		{"MyP@ssw0rd", 3, 4, "Strong"},
		{"C0mpl3x!P@ssw0rd123", 4, 4, "Very strong"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			strength := PasswordStrength(tt.password)
			assert.GreaterOrEqual(t, strength, tt.minStrength, "Password: %s", tt.password)
			assert.LessOrEqual(t, strength, tt.maxStrength, "Password: %s", tt.password)
		})
	}
}
