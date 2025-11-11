package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBDriver   string
	DBSource   string
	JWTSecret  string
	ServerPort string
}

// Load reads configuration from environment variables
func Load() *Config {
	config := &Config{
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBSource:   getEnv("DB_SOURCE", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	// Validate required fields
	if config.DBSource == "" {
		panic("DB_SOURCE environment variable is required")
	}
	if config.JWTSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	return config
}

// getEnv retrieves an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// String returns a formatted string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf("Config{Driver: %s, Port: %s}", c.DBDriver, c.ServerPort)
}