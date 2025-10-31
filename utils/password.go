package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"unicode"
)

// PasswordPolicy defines password requirements
type PasswordPolicy struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// DefaultPasswordPolicy returns the default password policy
func DefaultPasswordPolicy() PasswordPolicy {
	return PasswordPolicy{
		MinLength:      8,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: true,
	}
}

// ValidatePassword checks if a password meets the policy requirements
func ValidatePassword(password string, policy PasswordPolicy) error {
	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if policy.RequireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if policy.RequireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if policy.RequireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// IsCommonPassword checks if password is in a list of common weak passwords
func IsCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "iloveyou", "master", "sunshine", "ashley",
		"bailey", "passw0rd", "shadow", "123123", "654321",
		"superman", "qazwsx", "michael", "football",
	}

	lowerPassword := string(unicode.ToLower(rune(password[0]))) + password[1:]
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	return false
}

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// PasswordStrength calculates password strength (0-4)
// 0: Very Weak, 1: Weak, 2: Fair, 3: Strong, 4: Very Strong
func PasswordStrength(password string) int {
	score := 0

	// Length bonus
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// Character variety
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasUpper && hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	// Penalize common passwords
	if IsCommonPassword(password) {
		score = max(0, score-2)
	}

	return min(score, 4)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateSecurePassword generates a random secure password
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 12
	}

	const (
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		numbers          = "0123456789"
		specialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		allChars         = uppercaseLetters + lowercaseLetters + numbers + specialChars
	)

	password := make([]byte, length)

	// Ensure at least one character from each category
	password[0] = uppercaseLetters[cryptoRandInt(len(uppercaseLetters))]
	password[1] = lowercaseLetters[cryptoRandInt(len(lowercaseLetters))]
	password[2] = numbers[cryptoRandInt(len(numbers))]
	password[3] = specialChars[cryptoRandInt(len(specialChars))]

	// Fill the rest randomly
	for i := 4; i < length; i++ {
		password[i] = allChars[cryptoRandInt(len(allChars))]
	}

	// Shuffle the password
	for i := range password {
		j := cryptoRandInt(len(password))
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// cryptoRandInt generates a cryptographically secure random integer between 0 and max-1
func cryptoRandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}
