package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccountNumber(t *testing.T) {
	// Setup test database
	LoadConfig()

	accountNumber := GenerateAccountNumber()

	assert.NotEmpty(t, accountNumber, "Account number should not be empty")
	assert.Len(t, accountNumber, 7, "Account number should be 7 characters long")
	assert.Equal(t, "100", accountNumber[:3], "Account number should start with 100")
}

func TestGetDBConnectionString(t *testing.T) {
	cfg := &Config{
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBHost:     "localhost",
		DBPort:     "3306",
		DBName:     "testdb",
	}

	connStr := cfg.GetDBConnectionString()
	expected := "testuser:testpass@tcp(localhost:3306)/testdb?parseTime=true"

	assert.Equal(t, expected, connStr, "Connection string should match expected format")
}
