package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	ServerPort string
	ServerHost string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisURL string

	// CORS
	AllowedOrigins string

	// TLS
	EnableTLS   bool
	TLSCertFile string
	TLSKeyFile  string

	// Session
	SessionExpirationHours int

	// Environment
	Environment string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		ServerPort:             getEnv("SERVER_PORT", "3000"),
		ServerHost:             getEnv("SERVER_HOST", "localhost"),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "3306"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBName:                 getEnv("DB_NAME", "tctssf"),
		RedisURL:               getEnv("REDIS_URL", "redis://localhost:6379/0"),
		AllowedOrigins:         getEnv("ALLOWED_ORIGINS", "*"),
		EnableTLS:              getEnvBool("ENABLE_TLS", false),
		TLSCertFile:            getEnv("TLS_CERT_FILE", "./certs/server.crt"),
		TLSKeyFile:             getEnv("TLS_KEY_FILE", "./certs/server.key"),
		SessionExpirationHours: getEnvInt("SESSION_EXPIRATION_HOURS", 24),
		Environment:            getEnv("ENVIRONMENT", "development"),
	}

	AppConfig = config
	return config
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool retrieves a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid boolean value for %s: %s, using default: %v", key, value, defaultValue)
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

// GetDBConnectionString returns the database connection string
func (c *Config) GetDBConnectionString() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?parseTime=true"
}
